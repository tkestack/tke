/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package api

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/apis/audit"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"tkestack.io/tke/api/auth"
	"tkestack.io/tke/api/business"
	"tkestack.io/tke/api/monitor"
	"tkestack.io/tke/api/notify"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/api/registry"
	auditconfig "tkestack.io/tke/pkg/audit/apis/config"
	auditconfigv1 "tkestack.io/tke/pkg/audit/apis/config/v1"
	"tkestack.io/tke/pkg/audit/config/codec"
	"tkestack.io/tke/pkg/audit/config/configfiles"
	"tkestack.io/tke/pkg/audit/storage"
	"tkestack.io/tke/pkg/audit/storage/es"
	"tkestack.io/tke/pkg/audit/storage/types"
	utilfs "tkestack.io/tke/pkg/util/filesystem"
	"tkestack.io/tke/pkg/util/log"
)

// GroupName is the api group name for audit.
const GroupName = "audit.tkestack.io"

// Version is the api version for audit.
const Version = "v1"

const blockKey = "block-clusters.txt"

// ClusterControlPlane is the cluster name the tkestack control-planes like tke-platform-api will use to report audit events
const ClusterControlPlane = "control-plane"

var controlPlaneGroups sets.String
var k8sClient kubernetes.Interface

var (
	l             sync.RWMutex
	storeCli      storage.AuditStorage
	blockClusters sets.String
	storeConf     auditconfig.Storage
)

func init() {
	controlPlaneGroups = sets.NewString(
		platform.GroupName,
		registry.GroupName,
		notify.GroupName,
		monitor.GroupName,
		business.GroupName,
		auth.GroupName,
	)
	k8sClient = initK8sClient()
	initWatcher()

	blockClusters = sets.NewString(loadBlockClusters()...)
}

func initK8sClient() kubernetes.Interface {
	clientConfig, err := rest.InClusterConfig()
	if err != nil {
		klog.Fatalf("failed load k8s inCluster config: %v", err)
	}
	kubeClient, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		klog.Fatalf("failed init k8s client: %v", err)
	}
	return kubeClient
}

func initWatcher() {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		klog.Fatalf("failed init fsnotify watcher: %v", err)
	}
	err = w.Add("/app/conf/")
	if err != nil {
		klog.Fatalf("failed add fsnotify watcher: %v", err)
	}
	go watchEvent(w)
}

func watchEvent(w *fsnotify.Watcher) {
	for {
		select {
		case <-w.Events:
			klog.Infof("configmap changed")

			kc := loadConfig()
			if kc != nil {
				var changed = func(a, b *auditconfig.ElasticSearchStorage) bool {
					if a.Username == b.Username && a.Password == b.Password &&
						a.ReserveDays == b.ReserveDays && a.Address == b.Address &&
						a.Indices == b.Indices {
						return false
					}
					return true
				}
				if changed(kc.Storage.ElasticSearch, storeConf.ElasticSearch) {
					klog.Infof("store config changed: %v", kc.Storage)
					esCli, err := es.NewStorage(kc.Storage.ElasticSearch)
					if err != nil {
						klog.Errorf("failed init es client: %v", err)
						continue
					}
					storeConf = kc.Storage
					storeCli.Stop()
					storeCli = esCli
					storeCli.Start()
				} else {
					klog.Infof("store config not changed")
				}
			} else {
				klog.Errorf("load store config failed")
			}

			clusters := loadBlockClusters()
			klog.Infof("reload block clusters: %v", clusters)
			l.Lock()
			blockClusters = sets.NewString(clusters...)
			l.Unlock()

		case err := <-w.Errors:
			klog.Errorf("fsnotify error: %v", err)
		}
	}
}

func loadConfig() *auditconfig.AuditConfiguration {
	loader, err := configfiles.NewFsLoader(utilfs.DefaultFs{}, "/app/conf/tke-audit-api-config.yaml")
	if err != nil {
		klog.Errorf("failed init loader: %v", err)
		return nil
	}
	kc, err := loader.Load()
	if err != nil {
		klog.Errorf("failed load audit config: %v", err)
		return nil
	}
	return kc
}

func loadBlockClusters() []string {
	data, err := ioutil.ReadFile(fmt.Sprintf("/app/conf/%s", blockKey))
	if err != nil {
		return nil
	}
	return strings.Split(string(data), ",")
}

// RegisterRoute is used to register prefix path routing matches for all
// configured backend components.
func RegisterRoute(container *restful.Container, cfg *auditconfig.AuditConfiguration) error {
	return registerAuditRoute(container, cfg)
}

func registerAuditRoute(container *restful.Container, cfg *auditconfig.AuditConfiguration) error {
	ws := new(restful.WebService)
	ws.Path(fmt.Sprintf("/apis/%s/%s/events", GroupName, Version))
	ws.Produces(restful.MIME_JSON)
	ws.Consumes(restful.MIME_JSON, "text/csv")
	var err error
	storeConf = cfg.Storage
	storeCli, err = es.NewStorage(cfg.Storage.ElasticSearch)
	if err != nil {
		return err
	}
	storeCli.Start()
	ws.Route(ws.POST("/sink/{clusterName}").To(sinkEvents).
		Operation("createEventsByCluster").
		Doc("Create new audit events").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON))
	ws.Route(ws.GET("/list").To(listEvents).
		Operation("listEvents").
		Doc("Create new audit events").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON))
	ws.Route(ws.GET("/listFieldValues").To(listFieldValues).
		Operation("listFieldValues2").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON))
	ws.Route(ws.GET("/listBlockClusters").To(listBlockClusters).
		Operation("listBlockClusters").
		Doc("list blocked clusters").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON))
	ws.Route(ws.POST("/blockCluster").To(blockClusterAudit).
		Operation("patchBlockCluster").
		Doc("block or un-block cluster to audit").
		Param(restful.QueryParameter("clustername", "clustername that will block or un-block")).
		Param(restful.QueryParameter("block", "true or false to indicate block or un-block")).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON))
	ws.Route(ws.POST("/updateStoreConfig").To(configStore).
		Operation("replaceStoreConfig").
		Doc("update store config for audit events").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON))
	ws.Route(ws.GET("/getStoreConfig").To(getStoreConfig).
		Operation("getStoreConfig").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON))
	ws.Route(ws.POST("/configTest").To(testStoreConfig).
		Operation("getConfigTest").
		Doc("test if store config is valid").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON))
	ws.Route(ws.GET("/filedownload").To(downloadEvents).
		Operation("getAuditEvents").
		Doc("download audit events").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON, "text/csv"))
	container.Add(ws)
	return nil
}

func sinkEvents(request *restful.Request, response *restful.Response) {
	var eventList audit.EventList
	err := request.ReadEntity(&eventList)
	if err != nil {
		log.Infof("failed read events: %v", err)
		response.Write([]byte("failed"))
		return
	}
	clusterName := request.PathParameter("clusterName")
	events := types.ConvertEvents(eventList.Items)
	for _, event := range events {
		event.ClusterName = clusterName
	}
	events = eventsFilter(events)
	err = storeCli.Save(events)
	if err != nil {
		log.Errorf("failed save events: %v", err)
	}
	response.Write([]byte("success"))

}

type filterFunc func(e *types.Event) bool

func controlPlaneFilter(e *types.Event) bool {
	if e.ClusterName == ClusterControlPlane {
		if !controlPlaneGroups.Has(e.APIGroup) {
			return false
		}
	}
	return true
}

func userKubeletFilter(e *types.Event) bool {
	if strings.HasPrefix(e.UserName, "system:node:") {
		return false
	}
	return true
}

func kubesystemServiceAccountFilter(e *types.Event) bool {
	if strings.HasPrefix(e.UserName, "system:serviceaccount:kube-system:") {
		return false
	}
	return true
}

func tkePlatformControllerFilter(e *types.Event) bool {
	if e.UserName == "admin" && e.Verb == "update" && e.Resource == "clusters" && strings.HasPrefix(e.UserAgent, "tke-platform-controller") {
		return false
	}
	return true
}

func blockClustersFilter(e *types.Event) bool {
	l.RLock()
	defer l.RUnlock()
	if blockClusters.Has(e.ClusterName) {
		return false
	}
	return true
}

var eventFilters = []filterFunc{
	controlPlaneFilter,
	userKubeletFilter,
	kubesystemServiceAccountFilter,
	tkePlatformControllerFilter,
	blockClustersFilter,
}

func eventFilter(e *types.Event) bool {
	for _, f := range eventFilters {
		if !f(e) {
			return false
		}
	}
	return true
}

func eventsFilter(events []*types.Event) []*types.Event {
	var result []*types.Event
	for i := range events {
		if eventFilter(events[i]) {
			result = append(result, events[i])
		}
	}
	return result
}

func parseQueryParam(request *restful.Request) *storage.QueryParameter {
	params := storage.QueryParameter{
		ClusterName: request.QueryParameter("cluster"),
		Namespace:   request.QueryParameter("namespace"),
		Resource:    request.QueryParameter("resource"),
		Name:        request.QueryParameter("name"),
		Query:       request.QueryParameter("query"),
		UserName:    request.QueryParameter("user"),
	}
	startTime := request.QueryParameter("startTime")
	if startTime != "" {
		if stime, err := strconv.ParseInt(startTime, 10, 64); err == nil {
			params.StartTime = stime
		}
	}
	endTime := request.QueryParameter("endTime")
	if endTime != "" {
		if etime, err := strconv.ParseInt(endTime, 10, 64); err == nil {
			params.EndTime = etime
		}
	}
	page := request.QueryParameter("pageIndex")
	size := request.QueryParameter("pageSize")
	if size != "" {
		if s, err := strconv.Atoi(size); err == nil && s > 0 {
			params.Size = s
		}
	}
	if params.Size == 0 {
		params.Size = 10
	}
	if page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			params.Offset = (p - 1) * params.Size
		}
	}
	return &params
}

func listEvents(request *restful.Request, response *restful.Response) {
	params := parseQueryParam(request)
	events, total, err := storeCli.Query(params)
	if err != nil {
		log.Errorf("failed to query events: %v", err)
		writeStatusResponse(response, err)
	} else {
		response.WriteEntity(Pagination{ResultStatus: ResultStatus{Code: 0, Message: ""}, Total: total, Items: events})
	}
}

func blockClusterAudit(request *restful.Request, response *restful.Response) {
	clusterName := request.QueryParameter("clustername")
	if clusterName == "" {
		writeStatusResponse(response, nil)
		return
	}
	block, _ := strconv.ParseBool(request.QueryParameter("block"))
	var blocked bool
	l.RLock()
	blocked = blockClusters.Has(clusterName)
	l.RUnlock()
	if block {
		if !blocked {
			cm, err := k8sClient.CoreV1().ConfigMaps("tke").Get(context.Background(), "tke-audit-api", metav1.GetOptions{})
			if err != nil {
				klog.Errorf("error get tke-audit-api configmap: %v", err)
				writeStatusResponse(response, err)
				return
			}
			str := cm.Data[blockKey]
			clusters := sets.NewString(strings.Split(str, ",")...)
			if !clusters.Has(clusterName) {
				clusters.Insert(clusterName)
				str = strings.Join(clusters.List(), ",")
				cm.Data[blockKey] = str
				_, err = k8sClient.CoreV1().ConfigMaps("tke").Update(context.Background(), cm, metav1.UpdateOptions{})
				if err != nil {
					klog.Errorf("failed update blockClusters: %v", err)
					writeStatusResponse(response, err)
					return
				}
			}
		}
	} else {
		if blocked {
			cm, err := k8sClient.CoreV1().ConfigMaps("tke").Get(context.Background(), "tke-audit-api", metav1.GetOptions{})
			if err != nil {
				klog.Errorf("error get tke-audit-api configmap: %v", err)
				writeStatusResponse(response, err)
				return
			}
			str := cm.Data[blockKey]
			clusters := sets.NewString(strings.Split(str, ",")...)
			if clusters.Has(clusterName) {
				clusters.Delete(clusterName)
				str = strings.Join(clusters.List(), ",")
				cm.Data[blockKey] = str
				_, err = k8sClient.CoreV1().ConfigMaps("tke").Update(context.Background(), cm, metav1.UpdateOptions{})
				if err != nil {
					klog.Errorf("failed update blockClusters: %v", err)
					writeStatusResponse(response, err)
					return
				}
			}
		}
	}
	writeStatusResponse(response, nil)
}

func listBlockClusters(request *restful.Request, response *restful.Response) {
	var clusters []string
	l.RLock()
	clusters = blockClusters.Delete("").List()
	l.RUnlock()
	response.WriteAsJson(clusters)
}

func configStore(request *restful.Request, response *restful.Response) {
	store := auditconfig.Storage{}
	request.ReadEntity(&store)
	if store.ElasticSearch == nil {
		writeStatusResponse(response, fmt.Errorf("store config not set"))
		return
	}
	_, err := es.NewStorage(store.ElasticSearch)
	if err != nil {
		writeStatusResponse(response, err)
		return
	}
	cm, err := k8sClient.CoreV1().ConfigMaps("tke").Get(context.Background(), "tke-audit-api", metav1.GetOptions{})
	if err != nil {
		writeStatusResponse(response, err)
		return
	}
	conf := auditconfig.AuditConfiguration{}
	conf.Storage = store
	data, err := codec.EncodeAuditConfig(&conf, auditconfigv1.SchemeGroupVersion)
	if err != nil {
		writeStatusResponse(response, err)
		return
	}
	cm.Data["tke-audit-api-config.yaml"] = string(data)
	_, err = k8sClient.CoreV1().ConfigMaps("tke").Update(context.Background(), cm, metav1.UpdateOptions{})
	if err != nil {
		klog.Errorf("failed update store config: %v", err)
		writeStatusResponse(response, err)
		return
	}
	writeStatusResponse(response, nil)
}

func getStoreConfig(request *restful.Request, response *restful.Response) {
	conf := storeConf
	conf.ElasticSearch.Username = "***"
	conf.ElasticSearch.Password = "***"
	response.WriteAsJson(conf)
}

func testStoreConfig(request *restful.Request, response *restful.Response) {
	store := auditconfig.Storage{}
	request.ReadEntity(&store)
	if store.ElasticSearch == nil {
		writeStatusResponse(response, fmt.Errorf("store config not set"))
		return
	}
	_, err := es.NewStorage(store.ElasticSearch)
	writeStatusResponse(response, err)
}

func downloadEvents(request *restful.Request, response *restful.Response) {
	params := parseQueryParam(request)
	params.Offset = 0
	if params.Size == 10 {
		params.Size = 10000
	}
	events, _, err := storeCli.Query(params)
	if err != nil {
		log.Errorf("failed to query events: %v", err)
		writeStatusResponse(response, err)
	} else {
		log.Infof("get %d events", len(events))
		b := &bytes.Buffer{}
		b.WriteString("\xEF\xBB\xBF") // UTF-8 BOM
		w := csv.NewWriter(b)
		if err := w.Write([]string{"时间", "操作人", "操作类型", "集群", "命名空间", "资源类型", "操作对象", "操作结果", "详情"}); err != nil {
			writeStatusResponse(response, err)
		}
		for _, event := range events {
			record := []string{
				time.Unix(event.RequestReceivedTimestamp/1000, 0).Format("2006-01-02T15:04:05Z07:00"),
				event.UserName,
				event.Verb,
				event.ClusterName,
				event.Namespace,
				event.Resource,
				event.Name,
				event.Status,
			}
			data, _ := json.Marshal(event)
			record = append(record, string(data))
			w.Write(record)
		}
		w.Flush()
		if err := w.Error(); err != nil {
			writeStatusResponse(response, err)
		}
		response.AddHeader("Content-Description", "File Transfer")
		response.AddHeader("Content-Disposition", fmt.Sprintf("attachment; filename=auditevent_%d.csv", time.Now().Unix()))
		response.AddHeader("Content-Type", "text/csv")
		response.Write(b.Bytes())
	}
}

func listFieldValues(request *restful.Request, response *restful.Response) {
	result := storeCli.FieldValues()
	response.WriteEntity(result)

}

// IgnoredAuthPathPrefixes returns a list of path prefixes that does not need to
// go through the built-in authentication and authorization middleware of apiserver.
func IgnoredAuthPathPrefixes() []string {
	return []string{
		fmt.Sprintf("/apis/%s/%s/events/sink", GroupName, Version),
	}
}

func writeStatusResponse(response *restful.Response, err error) {
	if err != nil {
		response.WriteHeaderAndEntity(http.StatusBadRequest, ResultStatus{Code: -1, Message: err.Error()})
		return
	}
	response.WriteEntity(ResultStatus{Code: 0, Message: ""})
}

type ResultStatus struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Pagination struct {
	ResultStatus `json:",inline"`
	Total        int            `json:"total"`
	Items        []*types.Event `json:"items"`
}
