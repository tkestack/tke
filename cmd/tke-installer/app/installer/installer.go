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

package installer

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"github.com/thoas/go-funk"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8stypes "k8s.io/apimachinery/pkg/types"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd/api"
	clientcmdlatest "k8s.io/client-go/tools/clientcmd/api/latest"
	certutil "k8s.io/client-go/util/cert"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	kubeaggregatorclientset "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
	applicationv1 "tkestack.io/tke/api/application/v1"
	applicationclientset "tkestack.io/tke/api/client/clientset/versioned/typed/application/v1"
	tkeclientset "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	registryclientset "tkestack.io/tke/api/client/clientset/versioned/typed/registry/v1"
	"tkestack.io/tke/api/platform"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/cmd/tke-installer/app/config"
	"tkestack.io/tke/cmd/tke-installer/app/installer/certs"
	"tkestack.io/tke/cmd/tke-installer/app/installer/constants"
	"tkestack.io/tke/cmd/tke-installer/app/installer/images"
	"tkestack.io/tke/cmd/tke-installer/app/installer/types"
	helmaction "tkestack.io/tke/pkg/application/helm/action"
	baremetalcluster "tkestack.io/tke/pkg/platform/provider/baremetal/cluster"
	baremetalconfig "tkestack.io/tke/pkg/platform/provider/baremetal/config"
	baremetalconstants "tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	baremetal "tkestack.io/tke/pkg/platform/provider/baremetal/images"
	galaxy "tkestack.io/tke/pkg/platform/provider/baremetal/phases/galaxy/images"
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
	clusterstrategy "tkestack.io/tke/pkg/platform/registry/cluster"
	v1 "tkestack.io/tke/pkg/platform/types/v1"
	platformutil "tkestack.io/tke/pkg/platform/util"

	"tkestack.io/tke/pkg/spec"
	"tkestack.io/tke/pkg/util/apiclient"
	"tkestack.io/tke/pkg/util/containerregistry"
	"tkestack.io/tke/pkg/util/docker"
	"tkestack.io/tke/pkg/util/hosts"
	utilhttp "tkestack.io/tke/pkg/util/http"
	"tkestack.io/tke/pkg/util/kubeconfig"
	"tkestack.io/tke/pkg/util/log"
	utilnet "tkestack.io/tke/pkg/util/net"
	"tkestack.io/tke/pkg/util/pkiutil"
	"tkestack.io/tke/pkg/util/ssh"

	// import platform schema
	_ "tkestack.io/tke/api/platform/install"
)

const namespace = "tke"

type TKE struct {
	Config  *config.Config           `json:"config"`
	Para    *types.CreateClusterPara `json:"para"`
	Cluster *v1.Cluster              `json:"cluster"`
	Step    int                      `json:"step"`
	// IncludeSelf means installer is using one of cluster's machines
	IncludeSelf bool `json:"includeSelf"`

	log             log.Logger
	steps           []types.Handler
	progress        *types.ClusterProgress
	strategy        *clusterstrategy.Strategy
	clusterProvider clusterprovider.Provider
	isFromRestore   bool

	docker *docker.Docker

	globalClient      kubernetes.Interface
	helmClient        *helmaction.Client
	platformClient    tkeclientset.PlatformV1Interface
	registryClient    registryclientset.RegistryV1Interface
	applicationClient applicationclientset.ApplicationV1Interface
	servers           []string
	namespace         string
}

func New(config *config.Config) *TKE {
	c := new(TKE)

	c.Config = config
	c.Para = new(types.CreateClusterPara)
	c.Cluster = new(v1.Cluster)
	c.progress = new(types.ClusterProgress)
	c.progress.Status = types.StatusUnknown

	clusterProvider, err := clusterprovider.GetProvider("Baremetal")
	if err != nil {
		panic(err)
	}
	c.clusterProvider = clusterProvider

	_ = os.MkdirAll(path.Dir(constants.ClusterLogFile), 0755)
	logOptions := log.NewOptions()
	logOptions.DisableColor = true
	logOptions.OutputPaths = []string{constants.ClusterLogFile}
	logOptions.ErrorOutputPaths = logOptions.OutputPaths
	log.Init(logOptions)
	c.log = log.WithName("tke-installer")

	c.docker = new(docker.Docker)
	c.docker.Stdout = c.log
	c.docker.Stderr = c.log

	if !config.Force {
		err = c.loadTKEData()
		if err == nil {
			c.isFromRestore = true
			c.progress.Status = types.StatusDoing
		}
	}

	return c
}

func (t *TKE) loadTKEData() error {
	data, err := ioutil.ReadFile(constants.ClusterFile)
	if err == nil {
		t.log.Infof("read %q success", constants.ClusterFile)
		err = json.Unmarshal(data, t)
		if err != nil {
			t.log.Infof("load tke data error:%s", err)
		} else {
			log.Infof("load tke data success")
		}
	}
	return err
}

func (t *TKE) initSteps() {

	if t.Config.EnableCustomExpansion {
		t.steps = append(t.steps, []types.Handler{
			{
				Name: "Init expansion",
				Func: t.initExpansion,
			},
		}...)
		t.steps = append(t.steps, []types.Handler{
			{
				Name: "Prepare expansion files",
				Func: t.prepareExpansionFiles,
			},
		}...)
	}

	t.steps = append(t.steps, []types.Handler{
		{
			Name: "Execute pre install hook",
			Func: t.preInstallHook,
		},
	}...)

	// UseDockerHub, no need load images, start local tcr and push images
	// TKERegistry load images && start local registry && push images to local registry
	// && deploy tke-registry-api && push images to tke-registry
	// ThirdPartyRegistry load images && push images
	if !t.Para.Config.Registry.IsOfficial() {
		t.steps = append(t.steps, []types.Handler{
			{
				Name: "Load images",
				Func: t.loadImages,
			},
			{
				Name: "Tag images",
				Func: t.tagImages,
			},
		}...)
	}

	// if both set, don't setup local registry
	if t.Para.Config.Registry.ThirdPartyRegistry == nil &&
		t.Para.Config.Registry.TKERegistry != nil {
		t.steps = append(t.steps, []types.Handler{
			{
				Name: "Setup local registry",
				Func: t.setupLocalRegistry,
			},
		}...)
	}

	if !t.Para.Config.Registry.IsOfficial() {
		t.steps = append(t.steps, []types.Handler{
			{
				Name: "Push base components images",
				Func: t.pushBaseComImages,
			},
		}...)
	}

	t.steps = append(t.steps, []types.Handler{
		{
			Name: "Generate certificates for TKE components",
			Func: t.generateCertificates,
		},
		{
			Name: "Create global cluster",
			Func: t.createGlobalCluster,
		},
		{
			Name: "Write kubeconfig",
			Func: t.writeKubeconfig,
		},
		{
			Name: "Execute post cluster ready hook",
			Func: t.postClusterReadyHook,
		},
		{
			Name: "Prepare front proxy certificates",
			Func: t.prepareFrontProxyCertificates,
		},
		{
			Name: "Create namespace for install TKE",
			Func: t.createNamespace,
		},
		{
			Name: "Prepare certificates",
			Func: t.prepareCertificates,
		},
		{
			Name: "Prepare baremetal provider config",
			Func: t.prepareBaremetalProviderConfig,
		},
		{
			Name: "Install etcd",
			Func: t.installETCD,
		},
		{
			Name: "Patch platform versions in cluster info",
			Func: t.patchPlatformVersion,
		},
	}...)

	t.steps = append(t.steps, []types.Handler{
		{
			Name: "Init Platform Applications",
			Func: t.initPlatformApps,
		},
		{
			Name: "Preprocess Platform Applications",
			Func: t.preprocessPlatformApps,
		},
		{
			Name: "Install Platform Applications",
			Func: t.installPlatformApps,
		},
	}...)

	if t.Para.Config.Registry.TKERegistry != nil {
		t.steps = append(t.steps, []types.Handler{
			{
				Name: "Prepare pull images",
				Func: t.prepareImages,
			},
			{
				Name: "Install tke-registry chart",
				Func: t.installTKERegistryChart,
			},
		}...)
	}

	if t.Para.Config.Gateway != nil {
		if t.IncludeSelf {
			t.steps = append(t.steps, []types.Handler{
				{
					Name: "Stop local registry to give up 80/443 for ingress",
					Func: t.stopLocalRegistry,
				},
			}...)
		}
		t.steps = append(t.steps, []types.Handler{
			{
				Name: "Install tke-gateway chart",
				Func: t.installTKEGatewayChart,
			},
			{
				Name: "Install ingress-nginx chart",
				Func: t.installIngressChart,
			},
		}...)
	}

	if t.Para.Config.Registry.ThirdPartyRegistry == nil &&
		t.Para.Config.Registry.TKERegistry != nil {
		t.steps = append(t.steps, []types.Handler{
			{
				Name: "Prepare push images to TKE registry",
				Func: t.preparePushImagesToTKERegistry,
			},
			{
				Name: "Push images to registry",
				Func: t.pushImages,
			},
			{
				Name: "Set global cluster hosts",
				Func: t.setGlobalClusterHosts,
			},
		}...)
	}

	if t.auditEnabled() {
		t.steps = append(t.steps, []types.Handler{
			{
				Name: "Install tke audit",
				Func: t.installTKEAudit,
			},
		}...)
	}

	if t.businessEnabled() {
		t.steps = append(t.steps, []types.Handler{
			{
				Name: "Install tke-business-api",
				Func: t.installTKEBusinessAPI,
			},
			{
				Name: "Install tke-business-controller",
				Func: t.installTKEBusinessController,
			},
		}...)
	}

	if t.Para.Config.Monitor != nil {
		if t.Para.Config.Monitor.InfluxDBMonitor != nil &&
			t.Para.Config.Monitor.InfluxDBMonitor.LocalInfluxDBMonitor != nil {
			t.steps = append(t.steps, []types.Handler{
				{
					Name: "Install InfluxDB chart",
					Func: t.installInfluxDBChart,
				},
			}...)
		}
		if t.Para.Config.Monitor.ThanosMonitor != nil {
			t.steps = append(t.steps, []types.Handler{
				{
					Name: "Install Thanos",
					Func: t.installThanos,
				},
			}...)
		}

		t.steps = append(t.steps, []types.Handler{
			{
				Name: "Install tke-monitor-api",
				Func: t.installTKEMonitorAPI,
			},
			{
				Name: "Install tke-monitor-controller",
				Func: t.installTKEMonitorController,
			},
			{
				Name: "Install tke-notify-api",
				Func: t.installTKENotifyAPI,
			},
			{
				Name: "Install tke-notify-controller",
				Func: t.installTKENotifyController,
			},
		}...)
	}

	if t.Para.Config.Logagent != nil {
		t.steps = append(t.steps, []types.Handler{
			{
				Name: "Install tke-logagent-api",
				Func: t.installTKELogagentAPI,
			},
			{
				Name: "Install tke-logagent-controller",
				Func: t.installTKELogagentController,
			},
		}...)
	}

	if t.Para.Config.Application != nil {
		t.steps = append(t.steps, []types.Handler{
			{
				Name: "Install tke-application-api",
				Func: t.installTKEApplicationAPI,
			},
			{
				Name: "Install tke-application-controller",
				Func: t.installTKEApplicationController,
			},
		}...)
	}

	if t.Para.Config.Mesh != nil {
		t.steps = append(t.steps, []types.Handler{
			{
				Name: "Install tke-mesh-api",
				Func: t.installTKEMeshAPI,
			},
			{
				Name: "Install tke-mesh-controller",
				Func: t.installTKEMeshController,
			},
		}...)
	}

	// others

	// Add more tke component before THIS!!!
	t.steps = append(t.steps, []types.Handler{
		{
			Name: "Register tke api into global cluster",
			Func: t.registerAPI,
		},
		{
			Name: "Import resource to TKE platform",
			Func: t.importResource,
		},
	}...)

	// if t.Para.Config.Registry.ThirdPartyRegistry == nil &&
	// 	t.Para.Config.Registry.TKERegistry != nil {
	// 	t.steps = append(t.steps, []types.Handler{
	// 		{
	// 			Name: "Import charts",
	// 			Func: t.importCharts,
	// 		},
	// 		{
	// 			Name: "Import Expansion Charts",
	// 			Func: t.importExpansionCharts,
	// 		},
	// 	}...)
	// }

	if len(t.Para.Config.ExpansionApps) > 0 {
		t.steps = append(t.steps, []types.Handler{
			{
				Name: "Install Applications",
				Func: t.installApplications,
			},
		}...)
	}

	t.steps = append(t.steps, []types.Handler{
		{
			Name: "Execute post install hook",
			Func: t.postInstallHook,
		},
	}...)

	t.steps = funk.Filter(t.steps, func(step types.Handler) bool {
		return !funk.ContainsString(t.Para.Config.SkipSteps, step.Name)
	}).([]types.Handler)

	t.log.Info("Steps:")
	for i, step := range t.steps {
		t.log.Infof("%d %s", i, step.Name)
	}
}

func (t *TKE) Run() {
	var err error
	if t.Config.NoUI {
		err = t.run()
	} else {
		err = t.runWithUI()
	}
	if err != nil {
		log.Error(err.Error())
	}
}

func (t *TKE) run() error {
	if !t.isFromRestore {
		if err := t.loadPara(); err != nil {
			return err
		}
		err := t.prepare()
		if err != nil {
			statusErr := err.Status()
			return errors.New(statusErr.String())
		}
	}

	t.do()

	return nil
}

func (t *TKE) runWithUI() error {
	a := NewAssertsResource()
	restful.Add(a.WebService())
	restful.Add(t.WebService())
	s := NewSSHResource()
	restful.Add(s.WebService())

	restful.Filter(globalLogging)

	switch {
	case t.Config.PrepareCustomK8sImages:
		err := t.prepareForPrepareCustomImages(context.Background())
		if err != nil {
			return err
		}
		go t.doPrepareCustomImages()
	case t.Config.PrepareCustomCharts:
		err := t.prepareForPrepareCustomCharts(context.Background())
		if err != nil {
			return err
		}
		go t.doPrepareCustomCharts()
	case t.Config.Upgrade:
		err := t.prepareForUpgrade(context.Background())
		if err != nil {
			return err
		}
		go t.do()
	default:
		if t.isFromRestore {
			go t.do()
		}
	}

	log.Infof("Starting %s at http://%s", t.Config.ServerName, t.Config.ListenAddr)
	return http.ListenAndServe(t.Config.ListenAddr, nil)
}

func (t *TKE) loadPara() error {
	data, err := ioutil.ReadFile(t.Config.Config)
	if err != nil {
		return fmt.Errorf("read config error:%s", err)
	}
	err = json.Unmarshal(data, t.Para)
	if err != nil {
		return fmt.Errorf("parse config error:%s", err)
	}
	log.Infow("read success", "Para", t.Para)

	return nil
}

// WebService creates a new service that can handle REST requests for ClusterResource.
func (t *TKE) WebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/api/cluster")
	ws.Produces(restful.MIME_JSON)
	ws.Consumes(restful.MIME_JSON)

	ws.Route(ws.POST("").To(t.createCluster).
		Reads(types.CreateClusterPara{}).Writes(types.CreateClusterPara{}))

	ws.Route(ws.PUT("{name}/retry").To(t.retryCreateCluster))

	ws.Route(ws.GET("{name}").To(t.findCluster).
		Writes(types.CreateClusterPara{}))

	ws.Route(ws.GET("{name}/progress").To(t.findClusterProgress))

	return ws
}

func (t *TKE) completeWithProvider() {
	clusterProvider, err := baremetalcluster.NewProvider()
	if err != nil {
		panic(err)
	}
	t.clusterProvider = clusterProvider
}

func (t *TKE) prepare() apierrors.APIStatus {
	t.setConfigDefault(&t.Para.Config)
	statusError := t.validateConfig(t.Para.Config)
	if statusError != nil {
		return statusError
	}

	platform.Scheme.Default(&t.Para.Cluster)
	t.setClusterDefault(&t.Para.Cluster, &t.Para.Config)

	// mock platform api
	t.completeWithProvider()
	t.strategy = clusterstrategy.NewStrategy(nil)

	ctx := request.WithUser(context.Background(), &user.DefaultInfo{Name: constants.DefaultTeantID})

	v1Cluster := &t.Para.Cluster

	// PrepareForCreate
	platformCluster := new(platform.Cluster)
	err := platform.Scheme.Convert(v1Cluster, platformCluster, nil)
	if err != nil {
		return apierrors.NewInternalError(err)
	}
	t.strategy.PrepareForCreate(ctx, platformCluster)

	// Validate
	kinds, _, err := platform.Scheme.ObjectKinds(v1Cluster)
	if err != nil {
		return apierrors.NewInternalError(err)
	}

	allErrs := t.strategy.Validate(ctx, platformCluster)
	if len(allErrs) > 0 {
		return apierrors.NewInvalid(kinds[0].GroupKind(), v1Cluster.GetName(), allErrs)
	}

	err = platform.Scheme.Convert(platformCluster, v1Cluster, nil)
	if err != nil {
		return apierrors.NewInternalError(err)
	}

	statusError = t.validateResource(v1Cluster)
	if statusError != nil {
		return statusError
	}

	t.Cluster.Cluster = v1Cluster
	for _, one := range t.Cluster.Spec.Machines {
		ok, err := utilnet.InterfaceHasAddr(one.IP)
		if err != nil {
			return apierrors.NewInternalError(err)
		}
		if ok {
			t.IncludeSelf = true
			break
		}
	}

	err = t.completeExpansionApps()
	if err != nil {
		return apierrors.NewInternalError(err)
	}
	t.backup()

	return nil
}

func (t *TKE) setConfigDefault(config *types.Config) {
	if config.Basic == nil {
		config.Basic = &types.Basic{
			Username: "admin",
			Password: []byte("admin"),
		}
	}

	if config.Registry.TKERegistry != nil {
		if config.Registry.TKERegistry.Domain == "" {
			config.Registry.TKERegistry.Domain = "registry.tke.com"
		}
		config.Registry.TKERegistry.Namespace = "library"
		config.Registry.TKERegistry.Username = config.Basic.Username
		config.Registry.TKERegistry.Password = config.Basic.Password
	}
	if config.Auth.TKEAuth != nil {
		config.Auth.TKEAuth.TenantID = constants.DefaultTeantID
		config.Auth.TKEAuth.Username = config.Basic.Username
		config.Auth.TKEAuth.Password = config.Basic.Password
	}

	if config.Gateway != nil {
		if config.Gateway.Domain == "" {
			config.Gateway.Domain = "console.tke.com"
		}
		if config.Gateway.Cert == nil {
			config.Gateway.Cert = &types.Cert{
				SelfSignedCert: &types.SelfSignedCert{},
			}
		}
	}

	if config.HA != nil {
		if config.HA.ThirdPartyHA != nil {
			if config.HA.ThirdPartyHA.VPort == 0 {
				config.HA.ThirdPartyHA.VPort = 6443
			}
		}
	}

	config.Logagent = new(types.Logagent)

	if config.Application != nil {
		config.Application.RegistryDomain = config.Registry.Domain()
		config.Application.RegistryUsername = config.Registry.Username()
		config.Application.RegistryPassword = config.Registry.Password()
	}
}

func (t *TKE) setClusterDefault(cluster *platformv1.Cluster, config *types.Config) {
	if cluster.APIVersion == "" {
		cluster.APIVersion = platformv1.SchemeGroupVersion.String()
	}
	if cluster.Kind == "" {
		cluster.Kind = "Cluster"
	}
	cluster.Name = "global"
	cluster.Spec.DisplayName = "TKE"
	cluster.Spec.TenantID = constants.DefaultTeantID
	if t.Para.Config.Auth.TKEAuth != nil {
		cluster.Spec.TenantID = t.Para.Config.Auth.TKEAuth.TenantID
	}
	if cluster.Spec.Version == "" {
		cluster.Spec.Version = spec.K8sVersions[0] // use newest version
	}
	if cluster.Spec.ClusterCIDR == "" {
		cluster.Spec.ClusterCIDR = "10.244.0.0/16"
	}
	if cluster.Spec.Type == "" {
		cluster.Spec.Type = t.clusterProvider.Name()
	}
	if cluster.Spec.NetworkDevice == "" {
		cluster.Spec.NetworkDevice = "eth0"
	}
	cluster.Spec.Features.EnableMasterSchedule = true

	cluster.Spec.PublicAlternativeNames = append(cluster.Spec.PublicAlternativeNames, t.Para.Config.Gateway.Domain)
	if config.HA != nil {
		if t.Para.Config.HA.TKEHA != nil {
			cluster.Spec.Features.HA = &platformv1.HA{
				TKEHA: &platformv1.TKEHA{
					VIP:  t.Para.Config.HA.TKEHA.VIP,
					VRID: t.Para.Config.HA.TKEHA.VRID,
				},
			}
		}
		if t.Para.Config.HA.ThirdPartyHA != nil {
			cluster.Spec.Features.HA = &platformv1.HA{
				ThirdPartyHA: &platformv1.ThirdPartyHA{
					VIP:   t.Para.Config.HA.ThirdPartyHA.VIP,
					VPort: t.Para.Config.HA.ThirdPartyHA.VPort,
				},
			}
		}
	}
	if config.Business != nil {
		cluster.Spec.Features.AuthzWebhookAddr = &platformv1.AuthzWebhookAddr{
			Builtin: &platformv1.BuiltinAuthzWebhookAddr{},
		}
	}
}

func (t *TKE) validateConfig(config types.Config) *apierrors.StatusError {
	validate := validator.New()
	err := validate.Struct(config)
	if err != nil {
		return apierrors.NewBadRequest(err.Error())
	}

	if config.Gateway != nil || config.Registry.TKERegistry != nil {
		if config.Basic.Username == "" || config.Basic.Password == nil {
			return apierrors.NewBadRequest("username or password required when enabled gateway or registry")
		}
	}

	if config.Auth.TKEAuth == nil && config.Auth.OIDCAuth == nil {
		return apierrors.NewBadRequest("tke auth or oidc auth required")
	}

	if config.Registry.TKERegistry == nil && config.Registry.ThirdPartyRegistry == nil {
		return apierrors.NewBadRequest("tke registry or third party registry required")
	}

	if config.Registry.ThirdPartyRegistry != nil && config.Registry.ThirdPartyRegistry.Username != "" {
		cmd := exec.Command("docker", "login",
			"--username", config.Registry.ThirdPartyRegistry.Username,
			"--password", string(config.Registry.ThirdPartyRegistry.Password),
			config.Registry.ThirdPartyRegistry.Domain,
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			if _, ok := err.(*exec.ExitError); ok {
				return apierrors.NewBadRequest(string(out))
			}
			return apierrors.NewInternalError(err)
		}
	}

	if config.Monitor != nil {
		if config.Monitor.InfluxDBMonitor != nil && config.Monitor.ESMonitor != nil {
			return apierrors.NewBadRequest("influxdb or es only had one")
		}
	}

	if config.Gateway != nil && config.Gateway.Cert.ThirdPartyCert != nil {
		statusError := t.validateCertAndKey(config.Gateway.Cert.ThirdPartyCert.Certificate,
			config.Gateway.Cert.ThirdPartyCert.PrivateKey, config.Gateway.Domain)
		if statusError != nil {
			return statusError
		}
	}

	return nil
}

func (t *TKE) validateCertAndKey(certificate []byte, privateKey []byte, dnsName string) *apierrors.StatusError {
	if (certificate != nil && privateKey == nil) || (certificate == nil && privateKey != nil) {
		return apierrors.NewBadRequest("certificate and privateKey must offer together")
	}

	if certificate != nil {
		_, err := tls.X509KeyPair(certificate, privateKey)
		if err != nil {
			return apierrors.NewBadRequest(err.Error())
		}

		certs1, err := certutil.ParseCertsPEM(certificate)
		if err != nil {
			return apierrors.NewBadRequest(err.Error())
		}
		err = certs1[0].VerifyHostname(dnsName)
		if err != nil {
			return apierrors.NewBadRequest(err.Error())
		}
	}
	return nil
}

// validateResource validate the cpu and memory of cluster machines whether meets the requirements.
func (t *TKE) validateResource(cluster *platformv1.Cluster) *apierrors.StatusError {
	var (
		errs             []error
		cpuSum           int
		memoryInBytesSum uint64
		firstNodeDisk    int
	)
	for i, machine := range cluster.Spec.Machines {
		sshConfig := &ssh.Config{
			User:       machine.Username,
			Host:       machine.IP,
			Port:       int(machine.Port),
			Password:   string(machine.Password),
			PrivateKey: machine.PrivateKey,
			PassPhrase: machine.PassPhrase,
		}
		s, err := ssh.New(sshConfig)
		if err != nil {
			return apierrors.NewInternalError(err)
		}

		if i == 0 {
			firstNodeDisk, err = ssh.DiskAvail(s, constants.PathForDiskSpaceRequest)
			if err != nil {
				return apierrors.NewInternalError(fmt.Errorf("get disk availble space error: %w", err))
			}
		}

		cpu, err := ssh.NumCPU(s)
		if err != nil {
			return apierrors.NewInternalError(fmt.Errorf("get cpu error: %w", err))
		}
		cpuSum += cpu

		memInBytes, err := ssh.MemoryCapacity(s)
		if err != nil {
			return apierrors.NewInternalError(fmt.Errorf("get memory error: %w", err))
		}
		memoryInBytesSum += memInBytes
	}
	if cpuSum < constants.CPURequest {
		errs = append(errs, fmt.Errorf("sum of cpu in all nodes needs to be greater than %d",
			constants.CPURequest))
	}
	if math.Ceil(float64(memoryInBytesSum)/1024/1024/1024) < constants.MemoryRequest {
		errs = append(errs, fmt.Errorf("sum of memory in all nodes needs to be greater than %d GiB",
			constants.MemoryRequest))
	}
	if firstNodeDisk < constants.FirstNodeDiskSpaceRequest {
		errs = append(errs, fmt.Errorf("availble space of disk for %s in the first node which you input in cluster needs to be greater than %d GiB",
			constants.PathForDiskSpaceRequest, constants.FirstNodeDiskSpaceRequest))
	}

	if len(errs) != 0 {
		return apierrors.NewBadRequest(utilerrors.NewAggregate(errs).Error())
	}

	return nil
}

func (t *TKE) completeProviderConfigForRegistry() error {
	c, err := baremetalconfig.New(constants.ProviderConfigFile)
	if err != nil {
		return err
	}
	c.Registry.Prefix = t.Para.Config.Registry.Prefix()

	if t.Para.Config.Registry.TKERegistry != nil {
		ip, err := utilnet.GetSourceIP(t.Para.Cluster.Spec.Machines[0].IP)
		if err != nil {
			return errors.Wrap(err, "get ip for registry error")
		}
		c.Registry.IP = ip
	}
	if t.auditEnabled() {
		c.Audit.Address = t.determineGatewayHTTPSAddress()
	}

	return c.Save(constants.ProviderConfigFile)
}

func (t *TKE) determineGatewayHTTPSAddress() string {
	var host string
	if t.Para.Config.Gateway.Domain != "" {
		host = t.Para.Config.Gateway.Domain
	} else if t.Para.Config.HA != nil {
		host = t.Para.Config.HA.VIP()
	} else {
		host = t.Para.Cluster.Spec.Machines[0].IP
	}
	return fmt.Sprintf("https://%s", host)
}

func (t *TKE) auditEnabled() bool {
	return t.Para.Config.Audit != nil &&
		t.Para.Config.Audit.ElasticSearch != nil &&
		t.Para.Config.Audit.ElasticSearch.Address != ""
}

func (t *TKE) businessEnabled() bool {
	return t.Para.Config.Business != nil
}

func (t *TKE) createCluster(req *restful.Request, rsp *restful.Response) {
	apiStatus := func() apierrors.APIStatus {
		if t.Step != 0 {
			return apierrors.NewAlreadyExists(platformv1.Resource("Cluster"), "global")
		}
		para := new(types.CreateClusterPara)
		err := req.ReadEntity(para)
		if err != nil {
			return apierrors.NewBadRequest(err.Error())
		}
		t.Para = para
		if err := t.prepare(); err != nil {
			return err
		}
		go t.do()

		return nil
	}()

	if apiStatus != nil {
		_ = rsp.WriteHeaderAndJson(int(apiStatus.Status().Code), apiStatus.Status(), restful.MIME_JSON)
	} else {
		_ = rsp.WriteHeaderAndEntity(http.StatusCreated, t.Para)
	}
}

func (t *TKE) retryCreateCluster(req *restful.Request, rsp *restful.Response) {
	go t.do()
	_ = rsp.WriteEntity(nil)
}

func (t *TKE) findCluster(request *restful.Request, response *restful.Response) {
	apiStatus := func() apierrors.APIStatus {
		clusterName := request.PathParameter("name")
		if t.Cluster.Cluster == nil || t.Cluster.Name != clusterName {
			return apierrors.NewNotFound(platform.Resource("Cluster"), clusterName)
		}

		return nil
	}()

	if apiStatus != nil {
		_ = response.WriteHeaderAndJson(int(apiStatus.Status().Code), apiStatus.Status(), restful.MIME_JSON)
	} else {
		_ = response.WriteEntity(&types.CreateClusterPara{
			Cluster: *t.Cluster.Cluster,
			Config:  t.Para.Config,
		})
	}
}

func (t *TKE) findClusterProgress(request *restful.Request, response *restful.Response) {
	var err error
	var data []byte
	apiStatus := func() apierrors.APIStatus {
		clusterName := request.PathParameter("name")
		if t.Cluster.Cluster == nil {
			return apierrors.NewBadRequest("no cluater available")
		}
		if t.Cluster.Name != clusterName {
			return apierrors.NewNotFound(platform.Resource("Cluster"), clusterName)
		}
		data, err = ioutil.ReadFile(constants.ClusterLogFile)
		if err != nil {
			return apierrors.NewInternalError(err)
		}
		t.progress.Data = string(data)

		return nil
	}()

	if apiStatus != nil {
		response.WriteHeaderAndJson(int(apiStatus.Status().Code), apiStatus.Status(), restful.MIME_JSON)
	} else {
		response.WriteEntity(t.progress)
	}
}

func (t *TKE) do() {
	ctx := t.log.WithContext(context.Background())

	var taskType string
	if t.Config.Upgrade {
		taskType = "upgrade"
		t.upgradeSteps()
	} else {
		taskType = "install"
		containerregistry.Init(t.Para.Config.Registry.Domain(), t.Para.Config.Registry.Namespace())
		t.initSteps()
	}

	if !t.Config.Upgrade && t.runAfterClusterReady() {
		t.initDataForDeployTKE()
	}

	t.doSteps(ctx, taskType)

	t.progress.Status = types.StatusSuccess
	if t.Para.Config.Gateway != nil {
		var host string
		if t.Para.Config.Gateway.Domain != "" {
			host = t.Para.Config.Gateway.Domain
		} else if t.Para.Config.HA != nil {
			host = t.Para.Config.HA.VIP()
		} else {
			host = t.Para.Cluster.Spec.Machines[0].IP
		}
		t.progress.URL = fmt.Sprintf("http://%s", host)

		t.progress.Username = t.Para.Config.Basic.Username
		t.progress.Password = t.Para.Config.Basic.Password

		if t.Para.Config.Gateway.Cert.SelfSignedCert != nil {
			t.progress.CACert, _ = ioutil.ReadFile(constants.CACrtFile)
		}

		if t.Para.Config.Gateway.Domain != "" {
			t.progress.Hosts = append(t.progress.Hosts, t.Para.Config.Gateway.Domain)
		}

		cfg, _ := t.getKubeconfig()
		t.progress.Kubeconfig, _ = runtime.Encode(clientcmdlatest.Codec, cfg)
	}

	if t.Para.Config.Registry.TKERegistry != nil {
		t.progress.Hosts = append(t.progress.Hosts, t.Para.Config.Registry.TKERegistry.Domain)
	}

	if t.Para.Config.HA != nil {
		t.progress.Servers = append(t.progress.Servers, t.Para.Config.HA.VIP())
	}
	t.progress.Servers = append(t.progress.Servers, t.servers...)

}

func (t *TKE) doSteps(ctx context.Context, taskType string) {
	start := time.Now()
	if t.Step == 0 {
		t.log.Infof("===>starting %s task", taskType)
		t.progress.Status = types.StatusDoing
	}

	for t.Step < len(t.steps) {
		wait.PollInfinite(10*time.Second, func() (bool, error) {
			t.log.Infof("%d.%s doing", t.Step, t.steps[t.Step].Name)
			start := time.Now()
			err := t.steps[t.Step].Func(ctx)
			if err != nil {
				t.progress.Status = types.StatusRetrying
				t.log.Errorf("%d.%s [Failed] [%fs] error %s", t.Step, t.steps[t.Step].Name, time.Since(start).Seconds(), err)
				return false, nil
			}
			t.log.Infof("%d.%s [Success] [%fs]", t.Step, t.steps[t.Step].Name, time.Since(start).Seconds())

			t.Step++
			t.backup()
			t.progress.Status = types.StatusDoing
			return true, nil
		})
	}

	t.log.Infof("===>%s task [Sucesss] [%fs]", taskType, time.Since(start).Seconds())
}

func (t *TKE) runAfterClusterReady() bool {
	return t.Cluster.Status.Phase == platformv1.ClusterRunning
}

func (t *TKE) generateCertificates(ctx context.Context) error {
	var dnsNames []string
	ips := []net.IP{net.ParseIP("127.0.0.1")}
	if t.Para.Config.Gateway != nil && t.Para.Config.Gateway.Domain != "" {
		if ip := net.ParseIP(t.Para.Config.Gateway.Domain); ip != nil {
			ips = append(ips, ip)
		} else {
			dnsNames = append(dnsNames, t.Para.Config.Gateway.Domain)
		}
	}
	if t.Para.Config.Registry.TKERegistry != nil {
		dnsNames = append(dnsNames, t.Para.Config.Registry.TKERegistry.Domain, "*."+t.Para.Config.Registry.TKERegistry.Domain)
	}

	for _, one := range t.Cluster.Spec.Machines {
		ips = append(ips, net.ParseIP(one.IP))
	}
	if t.Para.Config.HA != nil {
		if t.Para.Config.HA.TKEHA != nil {
			ips = append(ips, net.ParseIP(t.Para.Config.HA.TKEHA.VIP))
		}
		if t.Para.Config.HA.ThirdPartyHA != nil {
			ips = append(ips, net.ParseIP(t.Para.Config.HA.ThirdPartyHA.VIP))
		}
	}
	return certs.Generate(dnsNames, ips, constants.DataDir)
}

func (t *TKE) prepareFrontProxyCertificates(ctx context.Context) error {
	machine := t.Cluster.Spec.Machines[0]
	sshConfig := &ssh.Config{
		User:       machine.Username,
		Host:       machine.IP,
		Port:       int(machine.Port),
		Password:   string(machine.Password),
		PrivateKey: machine.PrivateKey,
		PassPhrase: machine.PassPhrase,
	}
	s, err := ssh.New(sshConfig)
	if err != nil {
		return err
	}
	data, err := s.ReadFile("/etc/kubernetes/pki/front-proxy-ca.crt")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(constants.FrontProxyCACrtFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (t *TKE) createGlobalCluster(ctx context.Context) error {
	// update provider config and recreate
	err := t.completeProviderConfigForRegistry()
	if err != nil {
		return err
	}
	t.completeWithProvider()

	t.Cluster.Spec.Features.ContainerRuntime = platformv1.Containerd

	if t.Cluster.Spec.ClusterCredentialRef == nil {
		credential := &platformv1.ClusterCredential{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("cc-%s", t.Cluster.Name),
			},
			TenantID:    t.Cluster.Spec.TenantID,
			ClusterName: t.Cluster.Name,
		}
		t.Cluster.ClusterCredential = credential
		t.Cluster.Spec.ClusterCredentialRef = &corev1.LocalObjectReference{Name: credential.Name}
	}

	for t.Cluster.Status.Phase == platformv1.ClusterInitializing {
		err := t.clusterProvider.OnCreate(ctx, t.Cluster)
		if err != nil {
			return err
		}
		t.backup()
	}

	err = t.initDataForDeployTKE()
	if err != nil {
		return fmt.Errorf("init data for deploy tke error: %w", err)
	}

	return nil
}

func (t *TKE) backup() error {
	data, err := json.MarshalIndent(t, "", " ")
	if err != nil {
		t.log.Infof("json marshal tke failed, err = %s", err.Error())
		return err
	}
	return ioutil.WriteFile(constants.ClusterFile, data, 0777)
}
func (t *TKE) loadImages(ctx context.Context) error {
	if _, err := os.Stat(constants.ImagesFile); err != nil {
		return err
	}
	return t.docker.LoadImages(constants.ImagesFile)
}

func (t *TKE) tagImages(ctx context.Context) error {
	tkeImages, err := t.docker.GetImages(constants.ImagesPattern)
	if err != nil {
		return err
	}

	for _, image := range tkeImages {
		imageNames := strings.Split(image, "/")
		if len(imageNames) != 2 {
			t.log.Infof("invalid image name:name=%s", image)
			continue
		}
		name, _, _, err := t.docker.GetNameArchTag(imageNames[1])
		if err != nil {
			t.log.Infof("skip invalid image: %s", image)
			continue
		}
		if name == "tke-installer" { // no need to push installer image for speed up
			continue
		}

		target := fmt.Sprintf("%s/%s/%s", t.Para.Config.Registry.Domain(), t.Para.Config.Registry.Namespace(), imageNames[1])

		err = t.docker.TagImage(image, target)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *TKE) setupLocalRegistry(ctx context.Context) error {
	server := t.Para.Config.Registry.Domain()

	// for push image to local registry
	localHosts := hosts.LocalHosts{Host: server, File: "hosts"}
	err := localHosts.Set("127.0.0.1")
	if err != nil {
		return err
	}
	localHosts.File = "/app/hosts"
	err = localHosts.Set("127.0.0.1")
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile("hosts")
	if err != nil {
		return err
	}
	t.log.Info(string(data))

	return nil
}

func (t *TKE) readOrGenerateString(filename string) string {
	var (
		data []byte
		err  error
	)
	data, err = ioutil.ReadFile(filename)
	if err != nil {
		data = []byte(ksuid.New().String())
		ioutil.WriteFile(filename, data, 0644)
	}

	return string(data)
}

func (t *TKE) initDataForDeployTKE() error {
	var err error
	t.globalClient, err = t.Cluster.ClientsetForBootstrap()
	if err != nil {
		return err
	}

	t.helmClient, err = t.Cluster.HelmClientsetForBootstrap(t.namespace)
	if err != nil {
		return err
	}

	t.platformClient, err = t.Cluster.PlatformClientsetForBootstrap()
	if err != nil {
		return err
	}

	t.registryClient, err = t.Cluster.RegistryClientsetForBootstrap()
	if err != nil {
		return err
	}

	t.applicationClient, err = t.Cluster.RegistryApplicationForBootstrap()
	if err != nil {
		return err
	}

	for _, address := range t.Cluster.Status.Addresses {
		if address.Type == platformv1.AddressReal {
			t.servers = append(t.servers, address.Host)
		}
	}

	t.namespace = namespace

	return nil
}

func (t *TKE) createNamespace(ctx context.Context) error {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: t.namespace,
		},
	}

	return apiclient.CreateOrUpdateNamespace(ctx, t.globalClient, ns)
}

func (t *TKE) prepareCertificates(ctx context.Context) error {
	caCrt, err := ioutil.ReadFile(constants.CACrtFile)
	if err != nil {
		return err
	}
	caKey, err := ioutil.ReadFile(constants.CAKeyFile)
	if err != nil {
		return err
	}
	frontProxyCACrt, err := ioutil.ReadFile(constants.FrontProxyCACrtFile)
	if err != nil {
		return err
	}
	serverCrt, err := ioutil.ReadFile(constants.ServerCrtFile)
	if err != nil {
		return err
	}
	serverKey, err := ioutil.ReadFile(constants.ServerKeyFile)
	if err != nil {
		return err
	}
	adminCrt, err := ioutil.ReadFile(constants.AdminCrtFile)
	if err != nil {
		return err
	}
	adminKey, err := ioutil.ReadFile(constants.AdminKeyFile)
	if err != nil {
		return err
	}
	webhookCrt, err := ioutil.ReadFile(constants.WebhookCrtFile)
	if err != nil {
		return err
	}
	webhookKey, err := ioutil.ReadFile(constants.WebhookKeyFile)
	if err != nil {
		return err
	}

	if t.Cluster.Spec.Etcd.External != nil {
		return fmt.Errorf("external etcd specified, but ca key is not provided yet")
	}

	etcdClientCertData, etcdClientKeyData, err := pkiutil.GenerateClientCertAndKey(namespace, nil,
		t.Cluster.ClusterCredential.ETCDCACert, t.Cluster.ClusterCredential.ETCDCAKey)
	if err != nil {
		return fmt.Errorf("prepareCertificates fail:%w", err)
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "certs",
			Namespace: t.namespace,
		},
		Data: map[string]string{
			"etcd-ca.crt":        string(t.Cluster.ClusterCredential.ETCDCACert),
			"etcd.crt":           string(etcdClientCertData),
			"etcd.key":           string(etcdClientKeyData),
			"ca.crt":             string(caCrt),
			"ca.key":             string(caKey),
			"front-proxy-ca.crt": string(frontProxyCACrt),
			"server.crt":         string(serverCrt),
			"server.key":         string(serverKey),
			"admin.crt":          string(adminCrt),
			"admin.key":          string(adminKey),
			"webhook.crt":        string(webhookCrt),
			"webhook.key":        string(webhookKey),
		},
	}

	if t.Para.Config.Auth.OIDCAuth != nil {
		cm.Data["oidc-ca.crt"] = string(t.Para.Config.Auth.OIDCAuth.CACert)
	}

	if t.Para.Config.Registry.TKERegistry != nil && t.Para.Config.Registry.TKERegistry.HarborCAFile != "" {
		cm.Data["harbor-ca.crt"] = t.Para.Config.Registry.TKERegistry.HarborCAFile
	}

	cm.Data["password.csv"] = fmt.Sprintf("%s,admin,1,administrator", ksuid.New().String())
	cm.Data["token.csv"] = fmt.Sprintf("%s,admin,1,administrator", ksuid.New().String())

	for k, v := range cm.Data {
		err := ioutil.WriteFile(path.Join(constants.DataDir, k), []byte(v), 0644)
		if err != nil {
			return err
		}
	}

	return apiclient.CreateOrUpdateConfigMap(ctx, t.globalClient, cm)
}

func (t *TKE) authzWebhookBuiltinEndpoint() string {
	endPointHost := t.Para.Cluster.Spec.Machines[0].IP

	// use VIP in HA situation
	if t.Para.Cluster.Spec.Features.HA != nil {
		if t.Para.Cluster.Spec.Features.HA.TKEHA != nil {
			endPointHost = t.Para.Cluster.Spec.Features.HA.TKEHA.VIP
		}
		if t.Para.Cluster.Spec.Features.HA.ThirdPartyHA != nil {
			endPointHost = t.Para.Cluster.Spec.Features.HA.ThirdPartyHA.VIP
		}
	}

	return utilhttp.MakeEndpoint("https", endPointHost,
		constants.AuthzWebhookNodePort, "/auth/authz")
}

func (t *TKE) prepareBaremetalProviderConfig(ctx context.Context) error {
	providerConfig, err := baremetalconfig.New(constants.ProviderConfigFile)
	if err != nil {
		return err
	}
	if t.Para.Config.Registry.ThirdPartyRegistry == nil &&
		t.Para.Config.Registry.TKERegistry != nil {
		providerConfig.Registry.IP = t.Para.RegistryIP()
	}
	if t.auditEnabled() {
		providerConfig.Audit.Address = t.determineGatewayHTTPSAddress()
	}
	if t.businessEnabled() {
		providerConfig.Business.Enabled = true
	}
	providerConfig.PlatformAPIClientConfig = "conf/tke-platform-config.yaml"
	providerConfig.ApplicationAPIClientConfig = "conf/tke-application-config.yaml"
	// todo using ingress to expose authz service for ha.(
	//  users do not known nodeport when assigned vport in third party loadbalance)
	providerConfig.AuthzWebhook.Endpoint = t.authzWebhookBuiltinEndpoint()

	err = providerConfig.Save(constants.ProviderConfigFile)
	if err != nil {
		return err
	}

	configMaps := []struct {
		Name string
		File string
	}{
		{
			Name: "provider-config",
			File: baremetalconstants.ConfDir + "*.yaml",
		},
		{
			Name: "provider-config",
			File: baremetalconstants.ConfDir + "*.conf",
		},
		{
			Name: "docker",
			File: baremetalconstants.ConfDir + "docker/*",
		},
		{
			Name: "kubelet",
			File: baremetalconstants.ConfDir + "kubelet/*",
		},
		{
			Name: "kubeadm",
			File: baremetalconstants.ConfDir + "kubeadm/*",
		},
		{
			Name: "gpu-manifests",
			File: baremetalconstants.ManifestsDir + "/gpu/*",
		},
		{
			Name: "gpu-manager-manifests",
			File: baremetalconstants.ManifestsDir + "/gpu-manager/*",
		},
		{
			Name: "csi-operator-manifests",
			File: baremetalconstants.ManifestsDir + "/csi-operator/*",
		},
		{
			Name: "keepalived-manifests",
			File: baremetalconstants.ManifestsDir + "/keepalived/*",
		},
		{
			Name: "metrics-server-manifests",
			File: baremetalconstants.ManifestsDir + "/metrics-server/*",
		},
		{
			Name: "cilium-manifests",
			File: baremetalconstants.ManifestsDir + "/cilium/*",
		},
	}
	for _, one := range configMaps {
		err := apiclient.CreateOrUpdateConfigMapFromFile(ctx, t.globalClient,
			&corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      one.Name,
					Namespace: t.namespace,
				},
			},
			one.File)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *TKE) prepareImages(ctx context.Context) error {
	for _, machine := range t.Cluster.Spec.Machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}
		needPushImages := []string{images.Get().TKEGateway.FullName(),
			images.Get().TKERegistryAPI.FullName(),
			images.Get().TKERegistryController.FullName(),
			images.Get().TKEAuthAPI.FullName(),
			images.Get().TKEAuthController.FullName(),
			images.Get().TKEPlatformAPI.FullName(),
			images.Get().TKEPlatformController.FullName(),
			images.Get().NginxIngress.FullName(),
			images.Get().KebeWebhookCertgen.FullName()}
		for _, name := range needPushImages {
			cmdString := fmt.Sprintf("nerdctl --insecure-registry --namespace k8s.io pull %s", name)
			_, err = machineSSH.CombinedOutput(cmdString)
			if err != nil {
				return errors.Wrap(err, machine.IP)
			}
		}
	}
	return nil
}

func (t *TKE) stopLocalRegistry(ctx context.Context) error {
	if !t.docker.Healthz() {
		t.log.Info("Actively exit in order to reconnect to the docker service")
		os.Exit(1)
	}
	err := t.docker.RemoveContainers("registry-http", "registry-https")
	if err != nil {
		return err
	}
	return nil
}

func (t *TKE) installTKEGatewayChart(ctx context.Context) error {
	values := t.getTKEGatewayOptions(ctx)
	chartPathOptions := &helmaction.ChartPathOptions{}
	installOptions := &helmaction.InstallOptions{
		Namespace:        t.namespace,
		ReleaseName:      "tke-gateway",
		DependencyUpdate: false,
		Values:           values,
		Timeout:          10 * time.Minute,
		ChartPathOptions: *chartPathOptions,
	}

	chartFilePath := constants.ChartDirName + "tke-gateway/"
	if _, err := t.helmClient.InstallWithLocal(ctx, installOptions, chartFilePath); err != nil {
		// uninstallOptions := helmaction.UninstallOptions{
		// 	Timeout:     10 * time.Minute,
		// 	ReleaseName: "tke-gateway",
		// 	Namespace:   t.namespace,
		// }
		// reponse, err := t.helmClient.Uninstall(&uninstallOptions)
		// if err != nil {
		// 	return fmt.Errorf("%s uninstall fail, err = %s", reponse.Release.Name, err.Error())
		// }
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDaemonset(ctx, t.globalClient, t.namespace, "tke-gateway")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) getTKEGatewayOptions(ctx context.Context) map[string]interface{} {
	option := map[string]interface{}{
		"image":             images.Get().TKEGateway.FullName(),
		"oIDCClientSecret":  t.readOrGenerateString(constants.OIDCClientSecretFile),
		"selfSigned":        t.Para.Config.Gateway.Cert.SelfSignedCert != nil,
		"enableRegistry":    t.Para.Config.Registry.TKERegistry != nil,
		"enableAuth":        t.Para.Config.Auth.TKEAuth != nil,
		"enableMonitor":     t.Para.Config.Monitor != nil,
		"enableBusiness":    t.businessEnabled(),
		"enableLogagent":    t.Para.Config.Logagent != nil,
		"enableAudit":       t.auditEnabled(),
		"enableApplication": t.Para.Config.Application != nil,
		"enableMesh":        t.Para.Config.Mesh != nil,
	}
	if t.Para.Config.Registry.TKERegistry != nil {
		option["registryDomainSuffix"] = t.Para.Config.Registry.TKERegistry.Domain
	}
	if t.Para.Config.Auth.TKEAuth != nil {
		option["tenantID"] = t.Para.Config.Auth.TKEAuth.TenantID
	}
	if t.Para.Config.Gateway.Cert.ThirdPartyCert != nil {
		option["serverCrt"] = string(t.Para.Config.Gateway.Cert.ThirdPartyCert.Certificate)
		option["serverKey"] = string(t.Para.Config.Gateway.Cert.ThirdPartyCert.PrivateKey)
	}
	return option
}

func (t *TKE) installTKELogagentAPI(ctx context.Context) error {
	options := map[string]interface{}{
		"Replicas":       t.Config.Replicas,
		"Image":          images.Get().TKELogagentAPI.FullName(),
		"TenantID":       t.Para.Config.Auth.TKEAuth.TenantID,
		"Username":       t.Para.Config.Auth.TKEAuth.Username,
		"EnableAuth":     t.Para.Config.Auth.TKEAuth != nil,
		"EnableRegistry": t.Para.Config.Registry.TKERegistry != nil,
		"EnableAudit":    t.auditEnabled(),
	}
	if t.Para.Config.Auth.OIDCAuth != nil {
		options["OIDCClientID"] = t.Para.Config.Auth.OIDCAuth.ClientID
		options["OIDCIssuerURL"] = t.Para.Config.Auth.OIDCAuth.IssuerURL
		options["UseOIDCCA"] = t.Para.Config.Auth.OIDCAuth.CACert != nil
	}
	err := apiclient.CreateResourceWithDir(ctx, t.globalClient, "manifests/tke-logagent-api/*.yaml", options)
	if err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(ctx, t.globalClient, t.namespace, "tke-logagent-api")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) installTKELogagentController(ctx context.Context) error {
	options := map[string]interface{}{
		"Replicas":          t.Config.Replicas,
		"Image":             images.Get().TKELogagentController.FullName(),
		"EnableAuth":        t.Para.Config.Auth.TKEAuth != nil,
		"EnableRegistry":    t.Para.Config.Registry.TKERegistry != nil,
		"RegistryDomain":    t.Para.Config.Registry.Domain(),
		"RegistryNamespace": t.Para.Config.Registry.Namespace(),
	}
	err := apiclient.CreateResourceWithDir(ctx, t.globalClient, "manifests/tke-logagent-controller/*.yaml", options)
	if err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(ctx, t.globalClient, t.namespace, "tke-logagent-controller")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) installETCD(ctx context.Context) error {
	return apiclient.CreateResourceWithDir(ctx, t.globalClient, "manifests/etcd/*.yaml", nil)
}

func (t *TKE) getTKEAuthAPIOptions(ctx context.Context) (map[string]interface{}, error) {
	redirectHosts := t.servers
	redirectHosts = append(redirectHosts, "tke-gateway")
	if t.Para.Config.Gateway != nil && t.Para.Config.Gateway.Domain != "" {
		redirectHosts = append(redirectHosts, t.Para.Config.Gateway.Domain)
	}
	if t.Para.Config.HA != nil {
		redirectHosts = append(redirectHosts, t.Para.Config.HA.VIP())
		redirectHosts = append(redirectHosts, t.Para.Config.HA.VIP()+":31443")
		redirectHosts = append(redirectHosts, t.Para.Config.HA.VIP()+":31180")
	}
	if t.Para.Cluster.Spec.PublicAlternativeNames != nil {
		redirectHosts = append(redirectHosts, t.Para.Cluster.Spec.PublicAlternativeNames...)
	}
	cacrt, err := ioutil.ReadFile(constants.DataDir + "ca.crt")
	if err != nil {
		return nil, err
	}

	option := map[string]interface{}{
		"replicas":         t.Config.Replicas,
		"image":            images.Get().TKEAuthAPI.FullName(),
		"oIDCClientSecret": t.readOrGenerateString(constants.OIDCClientSecretFile),
		"adminUsername":    t.Para.Config.Auth.TKEAuth.Username,
		"tenantID":         t.Para.Config.Auth.TKEAuth.TenantID,
		"redirectHosts":    redirectHosts,
		"nodePort":         constants.AuthzWebhookNodePort,
		"enableAudit":      t.auditEnabled(),
		"caCrt":            string(cacrt),
	}
	return option, nil
}

func (t *TKE) getTKEAuthControllerOptions(ctx context.Context) map[string]interface{} {
	option := map[string]interface{}{
		"replicas":      t.Config.Replicas,
		"image":         images.Get().TKEAuthController.FullName(),
		"adminUsername": t.Para.Config.Auth.TKEAuth.Username,
		"adminPassword": string(t.Para.Config.Auth.TKEAuth.Password),
	}
	return option
}

func (t *TKE) installTKEAudit(ctx context.Context) error {
	options := map[string]interface{}{
		"Replicas":   t.Config.Replicas,
		"Image":      images.Get().TKEAudit.FullName(),
		"EnableAuth": t.Para.Config.Auth.TKEAuth != nil,
	}
	if t.Para.Config.Auth.OIDCAuth != nil {
		options["OIDCClientID"] = t.Para.Config.Auth.OIDCAuth.ClientID
		options["OIDCIssuerURL"] = t.Para.Config.Auth.OIDCAuth.IssuerURL
		options["UseOIDCCA"] = t.Para.Config.Auth.OIDCAuth.CACert != nil
	}

	if t.Para.Config.Audit.ElasticSearch != nil {
		options["StorageType"] = "es"
		options["StorageAddress"] = t.Para.Config.Audit.ElasticSearch.Address
		options["ReserveDays"] = t.Para.Config.Audit.ElasticSearch.ReserveDays
		options["Username"] = t.Para.Config.Audit.ElasticSearch.Username
		options["Password"] = t.Para.Config.Audit.ElasticSearch.Password
		options["Index"] = t.Para.Config.Audit.ElasticSearch.Index
	}

	if err := apiclient.CreateResourceWithDir(ctx, t.globalClient, "manifests/tke-audit-api/*.yaml", options); err != nil {
		return err
	}

	if t.Para.Config.Audit.ElasticSearch != nil && strings.Compare(t.Para.Config.Audit.ElasticSearch.Username, "skipTKEAuditHealthCheck") == 0 {
		return nil
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(ctx, t.globalClient, t.namespace, "tke-audit-api")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) getTKEPlatformAPIOptions(ctx context.Context) (map[string]interface{}, error) {
	cacrt, err := ioutil.ReadFile(constants.DataDir + "ca.crt")
	if err != nil {
		return nil, err
	}
	options := map[string]interface{}{
		"replicas":    t.Config.Replicas,
		"image":       images.Get().TKEPlatformAPI.FullName(),
		"enableAuth":  t.Para.Config.Auth.TKEAuth != nil,
		"enableAudit": t.auditEnabled(),
		"caCrt":       string(cacrt),
	}
	if t.Para.Config.Auth.OIDCAuth != nil {
		options["oIDCClientID"] = t.Para.Config.Auth.OIDCAuth.ClientID
		options["oIDCIssuerURL"] = t.Para.Config.Auth.OIDCAuth.IssuerURL
		options["useOIDCCA"] = t.Para.Config.Auth.OIDCAuth.CACert != nil
	}
	return options, nil
}

func (t *TKE) getTKEPlatformControllerOptions(ctx context.Context) map[string]interface{} {
	options := map[string]interface{}{
		"replicas":                t.Config.Replicas,
		"image":                   images.Get().TKEPlatformController.FullName(),
		"providerResImage":        images.Get().ProviderRes.FullName(),
		"registryDomain":          t.Para.Config.Registry.Domain(),
		"registryNamespace":       t.Para.Config.Registry.Namespace(),
		"monitorStorageType":      "",
		"monitorStorageAddresses": "",
	}
	if t.Para.Config.Monitor != nil {
		if t.Para.Config.Monitor.InfluxDBMonitor != nil {
			options["monitorStorageType"] = "influxdb"
			if t.Para.Config.Monitor.InfluxDBMonitor.LocalInfluxDBMonitor != nil {
				options["monitorStorageAddresses"] = t.getLocalInfluxdbAddress()
			} else if t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor != nil {
				address := t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.URL
				if t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.Username != "" {
					address = address + "&u=" + t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.Username
				}
				if t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.Password != nil {
					address = address + "&p=" + string(t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.Password)
				}
				options["monitorStorageAddresses"] = address
			}
		} else if t.Para.Config.Monitor.ESMonitor != nil {
			options["monitorStorageType"] = "elasticsearch"
			address := t.Para.Config.Monitor.ESMonitor.URL
			if t.Para.Config.Monitor.ESMonitor.Username != "" {
				address = address + "&u=" + t.Para.Config.Monitor.ESMonitor.Username
			}
			if t.Para.Config.Monitor.ESMonitor.Password != nil {
				address = address + "&p=" + string(t.Para.Config.Monitor.ESMonitor.Password)
			}
			options["monitorStorageAddresses"] = address
		} else if t.Para.Config.Monitor.ThanosMonitor != nil {
			options["monitorStorageType"] = "thanos"
			// thanos receive remote-write node-port address
			options["monitorStorageAddresses"] = fmt.Sprintf("http://%s:31141", t.servers[0])
		}
	}
	return options
}

func (t *TKE) installTKEBusinessAPI(ctx context.Context) error {
	options := map[string]interface{}{
		"Replicas":                   t.Config.Replicas,
		"Image":                      images.Get().TKEBusinessAPI.FullName(),
		"TenantID":                   t.Para.Config.Auth.TKEAuth.TenantID,
		"Username":                   t.Para.Config.Auth.TKEAuth.Username,
		"SyncProjectsWithNamespaces": t.Config.SyncProjectsWithNamespaces,
		"EnableAuth":                 t.Para.Config.Auth.TKEAuth != nil,
		"EnableRegistry":             t.Para.Config.Registry.TKERegistry != nil,
		"EnableAudit":                t.auditEnabled(),
	}
	if t.Para.Config.Auth.OIDCAuth != nil {
		options["OIDCClientID"] = t.Para.Config.Auth.OIDCAuth.ClientID
		options["OIDCIssuerURL"] = t.Para.Config.Auth.OIDCAuth.IssuerURL
		options["UseOIDCCA"] = t.Para.Config.Auth.OIDCAuth.CACert != nil
	}
	err := apiclient.CreateResourceWithDir(ctx, t.globalClient, "manifests/tke-business-api/*.yaml", options)
	if err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(ctx, t.globalClient, t.namespace, "tke-business-api")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) installTKEBusinessController(ctx context.Context) error {
	err := apiclient.CreateResourceWithDir(ctx, t.globalClient, "manifests/tke-business-controller/*.yaml",
		map[string]interface{}{
			"Replicas":       t.Config.Replicas,
			"Image":          images.Get().TKEBusinessController.FullName(),
			"EnableAuth":     t.Para.Config.Auth.TKEAuth != nil,
			"EnableRegistry": t.Para.Config.Registry.TKERegistry != nil,
		})
	if err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(ctx, t.globalClient, t.namespace, "tke-business-controller")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) installInfluxDBChart(ctx context.Context) error {

	options, err := t.getInfluxDBOptions(ctx)
	if err != nil {
		return fmt.Errorf("get influxdb options failed: %v", err)
	}

	influxDB := &types.PlatformApp{
		HelmInstallOptions: &helmaction.InstallOptions{
			Namespace:        t.namespace,
			ReleaseName:      "influxdb",
			Values:           options,
			DependencyUpdate: false,
			ChartPathOptions: helmaction.ChartPathOptions{},
		},
		LocalChartPath: constants.ChartDirName + "influxdb/",
		Enable:         true,
		ConditionFunc: func() (bool, error) {
			ok, err := apiclient.CheckDeployment(ctx, t.globalClient, t.namespace, "influxdb")
			if err != nil || !ok {
				return false, nil
			}
			return ok, nil
		},
	}
	return t.installPlatformApp(ctx, influxDB)
}

func (t *TKE) getInfluxDBOptions(ctx context.Context) (map[string]interface{}, error) {

	options := map[string]interface{}{
		"image": images.Get().InfluxDB.FullName(),
	}

	useCephRbd, useCephFS, useNFS := false, false, false
	for _, platformApp := range t.Para.Config.PlatformApps {
		if !platformApp.Enable || !platformApp.Installed {
			continue
		}
		if strings.EqualFold(platformApp.HelmInstallOptions.ReleaseName, constants.CephRBDChartReleaseName) {
			useCephRbd = true
			options["cephRbd"] = true
			options["cephRbdPVCName"] = "ceph-rbd-influxdb-pvc"
			options["cephRbdStorageClassName"] = constants.CephRBDStorageClassName
			break
		}
		if strings.EqualFold(platformApp.HelmInstallOptions.ReleaseName, constants.CephFSChartReleaseName) {
			useCephFS = true
			options["cephFS"] = true
			options["cephFSPVCName"] = "ceph-fs-influxdb-pvc"
			options["cephFSStorageClassName"] = constants.CephFSStorageClassName
			break
		}
		if strings.EqualFold(platformApp.HelmInstallOptions.ReleaseName, constants.NFSChartReleaseName) {
			useNFS = true
			options["nfs"] = true
			options["nfsPVCName"] = "nfs-influxdb-pvc"
			options["nfsStorageClassName"] = constants.NFSStorageClassName
			break
		}
	}

	if !(useCephRbd || useNFS || useCephFS) {
		options["baremetalStorage"] = true
		node, err := apiclient.GetNodeByMachineIP(ctx, t.globalClient, t.servers[0])
		if err != nil {
			return nil, err
		}
		options["nodeName"] = node.Name
	}
	return options, nil
}

func (t *TKE) installThanos(ctx context.Context) error {
	// TODO:2021-02-23 deploy thanos
	/*node, err := apiclient.GetNodeByMachineIP(ctx, t.globalClient, t.servers[0])
	if err != nil {
		return err
	}*/
	bucketConfig := t.Para.Config.Monitor.ThanosMonitor.BucketConfig
	thanosYamlBytes, err := yaml.Marshal(bucketConfig)
	if err != nil {
		return err
	}
	thanosYaml := base64.StdEncoding.EncodeToString(thanosYamlBytes)
	params := map[string]interface{}{
		"Image":      images.Get().Thanos.FullName(),
		"ThanosYaml": thanosYaml,
	}
	err = apiclient.CreateResourceWithDir(ctx, t.globalClient, "manifests/thanos/*.yaml", params)
	if err != nil {
		return err
	}
	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckStatefulSet(ctx, t.globalClient, t.namespace, "thanos-store")
		if err != nil || !ok {
			return false, nil
		}
		ok, err = apiclient.CheckStatefulSet(ctx, t.globalClient, t.namespace, "thanos-receive")
		if err != nil || !ok {
			return false, nil
		}
		ok, err = apiclient.CheckStatefulSet(ctx, t.globalClient, t.namespace, "thanos-compact")
		if err != nil || !ok {
			return false, nil
		}
		ok, err = apiclient.CheckStatefulSet(ctx, t.globalClient, t.namespace, "thanos-rule")
		if err != nil || !ok {
			return false, nil
		}
		ok, err = apiclient.CheckDeployment(ctx, t.globalClient, t.namespace, "thanos-query")
		if err != nil || !ok {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) installTKEMonitorAPI(ctx context.Context) error {
	options := map[string]interface{}{
		"Replicas":       t.Config.Replicas,
		"Image":          images.Get().TKEMonitorAPI.FullName(),
		"EnableAuth":     t.Para.Config.Auth.TKEAuth != nil,
		"EnableBusiness": t.businessEnabled(),
		"EnableAudit":    t.auditEnabled(),
	}
	if t.Para.Config.Auth.OIDCAuth != nil {
		options["OIDCClientID"] = t.Para.Config.Auth.OIDCAuth.ClientID
		options["OIDCIssuerURL"] = t.Para.Config.Auth.OIDCAuth.IssuerURL
		options["UseOIDCCA"] = t.Para.Config.Auth.OIDCAuth.CACert != nil
	}
	if t.Para.Config.Monitor != nil {
		if t.Para.Config.Monitor.ESMonitor != nil {
			options["StorageType"] = "es"
			options["StorageAddress"] = t.Para.Config.Monitor.ESMonitor.URL
			options["StorageUsername"] = t.Para.Config.Monitor.ESMonitor.Username
			options["StoragePassword"] = string(t.Para.Config.Monitor.ESMonitor.Password)
		} else if t.Para.Config.Monitor.InfluxDBMonitor != nil {
			options["StorageType"] = "influxDB"

			if t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor != nil {
				options["StorageAddress"] = t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.URL
				options["StorageUsername"] = t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.Username
				options["StoragePassword"] = string(t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.Password)
			} else if t.Para.Config.Monitor.InfluxDBMonitor.LocalInfluxDBMonitor != nil {
				// todo
				options["StorageAddress"] = t.getLocalInfluxdbAddress()
			}
		} else if t.Para.Config.Monitor.ThanosMonitor != nil {
			options["StorageType"] = "thanos"
			// thanos-query address
			options["StorageAddresses"] = "http://thanos-query.tke.svc.cluster.local:9090"
		}
	}

	if err := apiclient.CreateResourceWithDir(ctx, t.globalClient, "manifests/tke-monitor-api/*.yaml", options); err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(ctx, t.globalClient, t.namespace, "tke-monitor-api")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) installTKEMonitorController(ctx context.Context) error {
	params := map[string]interface{}{
		"Replicas":                t.Config.Replicas,
		"Image":                   images.Get().TKEMonitorController.FullName(),
		"EnableBusiness":          t.businessEnabled(),
		"RegistryDomain":          t.Para.Config.Registry.Domain(),
		"RegistryNamespace":       t.Para.Config.Registry.Namespace(),
		"MonitorStorageType":      "",
		"MonitorStorageAddresses": "",
	}
	if t.Para.Config.Monitor != nil {
		if t.Para.Config.Monitor.ESMonitor != nil {
			address := t.Para.Config.Monitor.ESMonitor.URL
			params["StorageType"] = "es"
			params["StorageAddress"] = address
			params["StorageUsername"] = t.Para.Config.Monitor.ESMonitor.Username
			params["StoragePassword"] = string(t.Para.Config.Monitor.ESMonitor.Password)
			params["MonitorStorageType"] = "elasticsearch"
			if t.Para.Config.Monitor.ESMonitor.Username != "" {
				address = address + "&u=" + t.Para.Config.Monitor.ESMonitor.Username
			}
			if t.Para.Config.Monitor.ESMonitor.Password != nil {
				address = address + "&p=" + string(t.Para.Config.Monitor.ESMonitor.Password)
			}
			params["MonitorStorageAddresses"] = address
		} else if t.Para.Config.Monitor.InfluxDBMonitor != nil {
			params["StorageType"] = "influxDB"
			params["MonitorStorageType"] = "influxdb"
			if t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor != nil {
				address := t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.URL
				params["StorageAddress"] = address
				params["StorageUsername"] = t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.Username
				params["StoragePassword"] = string(t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.Password)
				if t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.Username != "" {
					address = address + "&u=" + t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.Username
				}
				if t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.Password != nil {
					address = address + "&p=" + string(t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.Password)
				}
				params["MonitorStorageAddresses"] = address
			} else if t.Para.Config.Monitor.InfluxDBMonitor.LocalInfluxDBMonitor != nil {
				params["StorageAddress"] = t.getLocalInfluxdbAddress()
				params["MonitorStorageAddresses"] = t.getLocalInfluxdbAddress()
			}
		} else if t.Para.Config.Monitor.ThanosMonitor != nil {
			params["StorageType"] = "thanos"
			params["MonitorStorageType"] = "thanos"
			// thanos-query address
			params["MonitorStorageAddresses"] = "http://thanos-query.tke.svc.cluster.local:9090"
		}
		params["RetentionDays"] = t.Para.Config.Monitor.RetentionDays // can accept a nil value
	}

	if err := apiclient.CreateResourceWithDir(ctx, t.globalClient, "manifests/tke-monitor-controller/*.yaml", params); err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(ctx, t.globalClient, t.namespace, "tke-monitor-controller")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) installTKENotifyAPI(ctx context.Context) error {
	options := map[string]interface{}{
		"Replicas":    t.Config.Replicas,
		"Image":       images.Get().TKENotifyAPI.FullName(),
		"EnableAuth":  t.Para.Config.Auth.TKEAuth != nil,
		"EnableAudit": t.auditEnabled(),
	}
	if t.Para.Config.Auth.OIDCAuth != nil {
		options["OIDCClientID"] = t.Para.Config.Auth.OIDCAuth.ClientID
		options["OIDCIssuerURL"] = t.Para.Config.Auth.OIDCAuth.IssuerURL
		options["UseOIDCCA"] = t.Para.Config.Auth.OIDCAuth.CACert != nil
	}
	err := apiclient.CreateResourceWithDir(ctx, t.globalClient, "manifests/tke-notify-api/*.yaml", options)
	if err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(ctx, t.globalClient, t.namespace, "tke-notify-api")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) installTKENotifyController(ctx context.Context) error {
	err := apiclient.CreateResourceWithDir(ctx, t.globalClient, "manifests/tke-notify-controller/*.yaml",
		map[string]interface{}{
			"Replicas": t.Config.Replicas,
			"Image":    images.Get().TKENotifyController.FullName(),
		})
	if err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(ctx, t.globalClient, t.namespace, "tke-notify-controller")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) installTKERegistryChart(ctx context.Context) error {
	registryAPIOptions, err := t.getTKERegistryAPIOptions(ctx)
	if err != nil {
		return fmt.Errorf("get tke-registry-api options failed: %v", err)
	}
	registryControllerOptions, err := t.getTKERegistryControllerOptions(ctx)
	if err != nil {
		return fmt.Errorf("get tke-registry-controller options failed: %v", err)
	}
	tkeRegistry := &types.PlatformApp{
		HelmInstallOptions: &helmaction.InstallOptions{
			Namespace:   t.namespace,
			ReleaseName: "tke-registry",
			Values: map[string]interface{}{
				"api":        registryAPIOptions,
				"controller": registryControllerOptions,
			},
			DependencyUpdate: false,
			ChartPathOptions: helmaction.ChartPathOptions{},
		},
		LocalChartPath: constants.ChartDirName + "tke-registry/",
		Enable:         true,
		ConditionFunc: func() (bool, error) {
			apiOk, err := apiclient.CheckDeployment(ctx, t.globalClient, t.namespace, "tke-registry-api")
			if err != nil {
				return false, nil
			}
			controllerOk, err := apiclient.CheckDeployment(ctx, t.globalClient, t.namespace, "tke-registry-controller")
			if err != nil {
				return false, nil
			}
			return apiOk && controllerOk, nil
		},
	}
	return t.installPlatformApp(ctx, tkeRegistry)
}

func (t *TKE) installIngressChart(ctx context.Context) error {
	rawValues := `
 controller:
   name: controller
   image:
     registry: registry.tke.com
     image: library/ingress-nginx-controller
     tag: "v1.1.3"
     digest: ""
     pullPolicy: IfNotPresent
     runAsUser: 101
     allowPrivilegeEscalation: true
   containerName: controller
   containerPort:
     http: 80
     https: 443
   dnsPolicy: ClusterFirstWithHostNet
   hostNetwork: true
   ingressClass: nginx
   kind: DaemonSet
   nodeSelector:
     node-role.kubernetes.io/master: ""
   service:
     enabled: false
   admissionWebhooks:
     patch:
       enabled: true
       image:
         registry: registry.tke.com
         image: library/kube-webhook-certgen
         tag: "v1.1.1"
         digest: ""
 `
	tkeRegistry := &types.PlatformApp{
		HelmInstallOptions: &helmaction.InstallOptions{
			Namespace:        t.namespace,
			ReleaseName:      "ingress-nginx",
			DependencyUpdate: false,
			ChartPathOptions: helmaction.ChartPathOptions{},
		},
		LocalChartPath: constants.ChartDirName + "ingress-nginx/",
		Enable:         true,
		RawValues:      rawValues,
		RawValuesType:  applicationv1.RawValuesTypeYaml,
		ConditionFunc: func() (bool, error) {
			ok, err := apiclient.CheckDaemonset(ctx, t.globalClient, t.namespace, "ingress-nginx-controller")
			if err != nil {
				return false, nil
			}
			return ok, nil
		},
	}
	return t.installPlatformApp(ctx, tkeRegistry)
}

func (t *TKE) getTKERegistryAPIOptions(ctx context.Context) (map[string]interface{}, error) {

	options := map[string]interface{}{
		"replicas":       t.Config.Replicas,
		"namespace":      t.namespace,
		"image":          images.Get().TKERegistryAPI.FullName(),
		"adminUsername":  t.Para.Config.Registry.TKERegistry.Username,
		"adminPassword":  string(t.Para.Config.Registry.TKERegistry.Password),
		"enableAuth":     t.Para.Config.Auth.TKEAuth != nil,
		"enableBusiness": t.businessEnabled(),
		"domainSuffix":   t.Para.Config.Registry.TKERegistry.Domain,
		"enableAudit":    t.auditEnabled(),
		"harborEnabled":  t.Para.Config.Registry.TKERegistry.HarborEnabled,
		"harborCAFile":   t.Para.Config.Registry.TKERegistry.HarborCAFile,
	}
	// check if s3 enabled
	storageConfig := t.Para.Config.Registry.TKERegistry.Storage
	s3Enabled := (storageConfig != nil && storageConfig.S3 != nil)
	options["s3Enabled"] = s3Enabled
	if s3Enabled {
		options["s3Storage"] = storageConfig.S3
	}
	// or enable filesystem by default
	options["filesystemEnabled"] = !s3Enabled
	if options["filesystemEnabled"] == true {
		useCephRbd, useCephFS, useNFS := false, false, false
		for _, platformApp := range t.Para.Config.PlatformApps {
			if !platformApp.Enable || !platformApp.Installed {
				continue
			}
			if strings.EqualFold(platformApp.HelmInstallOptions.ReleaseName, constants.CephRBDChartReleaseName) {
				useCephRbd = true
				options["cephRbd"] = true
				options["cephRbdPVCName"] = "ceph-rbd-registry-pvc"
				options["cephRbdStorageClassName"] = constants.CephRBDStorageClassName
				break
			}
			if strings.EqualFold(platformApp.HelmInstallOptions.ReleaseName, constants.CephFSChartReleaseName) {
				useCephFS = true
				options["cephFS"] = true
				options["cephFSPVCName"] = "ceph-fs-registry-pvc"
				options["cephFSStorageClassName"] = constants.CephFSStorageClassName
				break
			}
			if strings.EqualFold(platformApp.HelmInstallOptions.ReleaseName, constants.NFSChartReleaseName) {
				useNFS = true
				options["nfs"] = true
				options["nfsPVCName"] = "nfs-registry-pvc"
				options["nfsStorageClassName"] = constants.NFSStorageClassName
				break
			}
		}
		if !(useCephRbd || useCephFS || useNFS) {
			options["baremetalStorage"] = true
			node, err := apiclient.GetNodeByMachineIP(ctx, t.globalClient, t.servers[0])
			if err != nil {
				return nil, err
			}
			options["nodeName"] = node.Name
		}
	}

	if t.Para.Config.Auth.OIDCAuth != nil {
		options["oIDCClientID"] = t.Para.Config.Auth.OIDCAuth.ClientID
		options["oIDCIssuerURL"] = t.Para.Config.Auth.OIDCAuth.IssuerURL
		options["useOIDCCA"] = t.Para.Config.Auth.OIDCAuth.CACert != nil
	}

	return options, nil
}

func (t *TKE) getTKERegistryControllerOptions(ctx context.Context) (map[string]interface{}, error) {

	node, err := apiclient.GetNodeByMachineIP(ctx, t.globalClient, t.servers[0])
	if err != nil {
		return nil, err
	}

	options := map[string]interface{}{
		"replicas":           t.Config.Replicas,
		"image":              images.Get().TKERegistryController.FullName(),
		"nodeName":           node.Name,
		"adminUsername":      t.Para.Config.Registry.TKERegistry.Username,
		"adminPassword":      string(t.Para.Config.Registry.TKERegistry.Password),
		"enableAuth":         t.Para.Config.Auth.TKEAuth != nil,
		"enableBusiness":     t.businessEnabled(),
		"domainSuffix":       t.Para.Config.Registry.TKERegistry.Domain,
		"defaultChartGroups": defaultChartGroupsStringConfig,
	}
	// check if s3 enabled
	storageConfig := t.Para.Config.Registry.TKERegistry.Storage
	s3Enabled := (storageConfig != nil && storageConfig.S3 != nil)
	options["s3Enabled"] = s3Enabled
	if s3Enabled {
		options["s3Storage"] = storageConfig.S3
	}
	// or enable filesystem by default
	options["filesystemEnabled"] = !s3Enabled
	if options["filesystemEnabled"] == true {
		useCephRbd, useNFS := false, false
		for _, platformApp := range t.Para.Config.PlatformApps {
			if !platformApp.Enable || !platformApp.Installed {
				continue
			}
			if strings.EqualFold(platformApp.HelmInstallOptions.ReleaseName, constants.CephRBDChartReleaseName) {
				useCephRbd = true
				break
			}
			if strings.EqualFold(platformApp.HelmInstallOptions.ReleaseName, constants.NFSChartReleaseName) {
				useNFS = true
				break
			}
		}
		options["baremetalStorage"] = !(useCephRbd || useNFS)
	}
	return options, nil
}

func (t *TKE) installTKEApplicationAPI(ctx context.Context) error {
	options := map[string]interface{}{
		"Replicas":              t.Config.Replicas,
		"Image":                 images.Get().TKEApplicationAPI.FullName(),
		"EnableAuth":            t.Para.Config.Auth.TKEAuth != nil,
		"EnableRegistry":        t.Para.Config.Registry.TKERegistry != nil,
		"EnableAudit":           t.auditEnabled(),
		"RegistryAdminUsername": t.Para.Config.Application.RegistryUsername,
		"RegistryAdminPassword": string(t.Para.Config.Application.RegistryPassword),
		"RegistryDomainSuffix":  t.Para.Config.Application.RegistryDomain,
	}
	if t.Para.Config.Auth.OIDCAuth != nil {
		options["OIDCClientID"] = t.Para.Config.Auth.OIDCAuth.ClientID
		options["OIDCIssuerURL"] = t.Para.Config.Auth.OIDCAuth.IssuerURL
		options["UseOIDCCA"] = t.Para.Config.Auth.OIDCAuth.CACert != nil
	}
	err := apiclient.CreateResourceWithDir(ctx, t.globalClient, "manifests/tke-application-api/*.yaml", options)
	if err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(ctx, t.globalClient, t.namespace, "tke-application-api")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) installTKEApplicationController(ctx context.Context) error {
	err := apiclient.CreateResourceWithDir(ctx, t.globalClient, "manifests/tke-application-controller/*.yaml",
		map[string]interface{}{
			"Replicas":              t.Config.Replicas,
			"Image":                 images.Get().TKEApplicationController.FullName(),
			"RegistryAdminUsername": t.Para.Config.Application.RegistryUsername,
			"RegistryAdminPassword": string(t.Para.Config.Application.RegistryPassword),
			"RegistryDomainSuffix":  t.Para.Config.Application.RegistryDomain,
		})
	if err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(ctx, t.globalClient, t.namespace, "tke-application-controller")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) installTKEMeshAPI(ctx context.Context) error {
	options := map[string]interface{}{
		"Replicas":    t.Config.Replicas,
		"Image":       images.Get().TKEMeshAPI.FullName(),
		"EnableAuth":  t.Para.Config.Auth.TKEAuth != nil,
		"EnableAudit": t.auditEnabled(),
	}
	if t.Para.Config.Auth.OIDCAuth != nil {
		options["OIDCClientID"] = t.Para.Config.Auth.OIDCAuth.ClientID
		options["OIDCIssuerURL"] = t.Para.Config.Auth.OIDCAuth.IssuerURL
		options["UseOIDCCA"] = t.Para.Config.Auth.OIDCAuth.CACert != nil
	}

	if err := apiclient.CreateResourceWithDir(ctx, t.globalClient, "manifests/tke-mesh-api/*.yaml", options); err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(ctx, t.globalClient, t.namespace, "tke-mesh-api")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) installTKEMeshController(ctx context.Context) error {
	err := apiclient.CreateResourceWithDir(ctx, t.globalClient, "manifests/tke-mesh-controller/*.yaml",
		map[string]interface{}{
			"Replicas":          t.Config.Replicas,
			"Image":             images.Get().TKEMeshController.FullName(),
			"RegistryDomain":    t.Para.Config.Registry.Domain(),
			"RegistryNamespace": t.Para.Config.Registry.Namespace(),
		})
	if err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(ctx, t.globalClient, t.namespace, "tke-mesh-controller")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) preparePushImagesToTKERegistry(ctx context.Context) error {
	if !t.docker.Healthz() {
		t.log.Info("Actively exit in order to reconnect to the docker service")
		os.Exit(1)
	}
	domains := []string{
		t.Para.Config.Registry.Domain(),
		constants.DefaultTeantID + "." + t.Para.Config.Registry.Domain(),
	}

	ip := t.servers[0]
	if t.Para.Config.HA != nil && len(t.Para.Config.HA.VIP()) > 0 {
		ip = t.Para.Config.HA.VIP()
	}

	for _, domain := range domains {
		localHosts := hosts.LocalHosts{Host: domain, File: "hosts"}
		err := localHosts.Set(ip)
		if err != nil {
			return err
		}
		localHosts.File = "/app/hosts"
		err = localHosts.Set(ip)
		if err != nil {
			return err
		}
	}

	dir := path.Join(constants.DockerCertsDir, t.Para.Config.Registry.Domain())
	_ = os.MkdirAll(dir, 0777)
	caCert, _ := ioutil.ReadFile(constants.CACrtFile)
	err := ioutil.WriteFile(path.Join(dir, "ca.crt"), caCert, 0644)
	if err != nil {
		return err
	}

	cmd := exec.Command("docker", "login",
		"--username", t.Para.Config.Registry.TKERegistry.Username,
		"--password", string(t.Para.Config.Registry.TKERegistry.Password),
		t.Para.Config.Registry.Domain(),
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return errors.New(string(out))
		}
		return err
	}

	return nil
}

func (t *TKE) registerAPI(ctx context.Context) error {
	caCert, _ := ioutil.ReadFile(constants.CACrtFile)

	restConfig, err := t.Cluster.RESTConfigForBootstrap()
	if err != nil {
		return err
	}
	client, err := kubeaggregatorclientset.NewForConfig(restConfig)
	if err != nil {
		return err
	}

	svcs := []string{"tke-platform-api"}
	if t.Para.Config.Auth.TKEAuth != nil {
		svcs = append(svcs, "tke-auth-api")
	}
	if t.businessEnabled() {
		svcs = append(svcs, "tke-business-api")
	}
	if t.Para.Config.Monitor != nil {
		svcs = append(svcs, "tke-notify-api", "tke-monitor-api")
	}
	if t.Para.Config.Registry.TKERegistry != nil {
		svcs = append(svcs, "tke-registry-api")
	}
	if t.Para.Config.Application != nil {
		svcs = append(svcs, "tke-application-api")
	}
	if t.Para.Config.Mesh != nil {
		svcs = append(svcs, "tke-mesh-api")
	}
	for _, one := range svcs {
		name := strings.TrimSuffix(one[4:], "-api")

		apiService := &apiregistrationv1.APIService{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("v1.%s.tkestack.io", name),
			},
			Spec: apiregistrationv1.APIServiceSpec{
				Group:                fmt.Sprintf("%s.tkestack.io", name),
				GroupPriorityMinimum: 1000,
				CABundle:             caCert,
				Version:              "v1",
				VersionPriority:      5,
				Service: &apiregistrationv1.ServiceReference{
					Namespace: t.namespace,
					Name:      one,
				},
			},
		}

		err = wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
			_, err := client.ApiregistrationV1().APIServices().Get(ctx, apiService.Name, metav1.GetOptions{})
			if err == nil {
				err := client.ApiregistrationV1().APIServices().Delete(ctx, apiService.Name, metav1.DeleteOptions{})
				if err != nil {
					return false, nil
				}
			}
			if _, err := client.ApiregistrationV1().APIServices().Create(ctx, apiService, metav1.CreateOptions{}); err != nil {
				if !apierrors.IsAlreadyExists(err) {
					return false, nil
				}
			}
			return true, nil
		})
		if err != nil {
			return errors.Wrapf(err, "register apiservice %v error", one)
		}

		err = wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
			a, err := client.ApiregistrationV1().APIServices().Get(ctx, apiService.Name, metav1.GetOptions{})
			if err != nil {
				return false, nil
			}
			for _, one := range a.Status.Conditions {
				return one.Type == apiregistrationv1.Available && one.Status == apiregistrationv1.ConditionTrue, nil
			}
			return false, nil
		})
		if err != nil {
			return errors.Wrapf(err, "check apiservices %v error", one)
		}
	}
	return nil
}

func (t *TKE) importResource(ctx context.Context) error {
	var err error
	// ensure api ready
	err = wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		_, err = t.platformClient.Clusters().List(ctx, metav1.ListOptions{})
		if err != nil {
			return false, nil
		}
		_, err = t.platformClient.ClusterCredentials().List(ctx, metav1.ListOptions{})
		if err != nil {
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return err
	}

	_, err = t.platformClient.ClusterCredentials().Get(ctx, t.Cluster.ClusterCredential.Name, metav1.GetOptions{})
	if err == nil {
		err := t.platformClient.ClusterCredentials().Delete(ctx, t.Cluster.ClusterCredential.Name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	_, err = t.platformClient.ClusterCredentials().Create(ctx, t.Cluster.ClusterCredential, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	_, err = t.platformClient.Clusters().Get(ctx, t.Cluster.Name, metav1.GetOptions{})
	if err == nil {
		err := t.platformClient.Clusters().Delete(ctx, t.Cluster.Name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	_, err = t.platformClient.Clusters().Create(ctx, t.Cluster.Cluster, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (t *TKE) pushImages(ctx context.Context) error {
	imagesFilter := fmt.Sprintf("%s/*", t.Para.Config.Registry.Namespace())
	if t.Para.Config.Registry.Domain() != "docker.io" { // docker images filter ignore docker.io
		imagesFilter = t.Para.Config.Registry.Domain() + "/" + imagesFilter
	}
	tkeImages, err := t.docker.GetImages(imagesFilter)
	if err != nil {
		return err
	}
	return t.dockerPush(tkeImages)
}
func (t *TKE) pushBaseComImages(ctx context.Context) error {
	archsFlag := []string{"amd64", "arm64"}
	supportMultiArchImages := []func() []string{
		baremetal.List,
		images.ListBaseComponents,
		galaxy.List,
	}

	var result []string
	for _, f := range supportMultiArchImages {
		for _, one := range f() {
			one := fmt.Sprintf("%s/%s/%s", t.Para.Config.Registry.Domain(),
				t.Para.Config.Registry.Namespace(), one)
			if isUnsupportMultiArch(one) {
				result = append(result, one)
			} else {
				for _, arch := range archsFlag {
					result = append(result, strings.ReplaceAll(one, ":", "-"+arch+":"))
				}
			}
		}
	}

	result = funk.UniqString(result)
	return t.dockerPush(result)
}

func isUnsupportMultiArch(name string) bool {
	specialUnsupportMultiArch := []string{"nvidia-device-plugin", "gpu"}
	for _, one := range specialUnsupportMultiArch {
		if strings.Contains(name, one) {
			return true
		}
	}

	return false
}

func (t *TKE) dockerPush(tkeImages []string) error {
	sort.Strings(tkeImages)
	tkeImagesSet := sets.NewString(tkeImages...)
	manifestSet := sets.NewString()

	// clear all local manifest lists before create any manifest list
	err := t.docker.ClearLocalManifests()
	if err != nil {
		return err
	}

	manifestsChan := make(chan string, 1)

	for _, image := range tkeImages {
		go func(image string) {
			for {
				err := t.pushTKEImage(image, tkeImagesSet, manifestsChan)
				if err == nil {
					break
				}
				t.log.Errorf("push %s failed: %v", err)
				time.Sleep(5 * time.Second)
			}
		}(image)
	}

	for range tkeImages {
		manifestName := <-manifestsChan
		if manifestName != "" {
			manifestSet.Insert(manifestName)
		}
	}

	sortedManifests := manifestSet.List()
	for _, manifest := range sortedManifests {
		go func(manifest string) {
			for {
				err := t.pushTKEManifest(manifest, manifestsChan)
				if err == nil {
					break
				}
				t.log.Errorf("push manifest %s failed: %v", err)
				time.Sleep(5 * time.Second)
			}
		}(manifest)
	}
	for range sortedManifests {
		<-manifestsChan
	}

	close(manifestsChan)

	return nil
}

func (t *TKE) pushTKEManifest(manifest string, manifestsChan chan string) error {
	err := t.docker.PushManifest(manifest, true)
	if err != nil {
		return err
	}
	manifestsChan <- ""
	t.log.Infof("push manifest %s to registry success", manifest)
	return nil
}

func (t *TKE) pushTKEImage(image string, tkeImagesSet sets.String, manifestsChan chan string) error {
	name, arch, tag, err := t.docker.GetNameArchTag(image)
	var manifestName string
	if err != nil { // skip invalid image
		t.log.Infof("skip invalid image: %s", image)
		manifestsChan <- ""
		return nil
	}

	if arch == "" {
		// ignore image without arch when has image with arch for avoid overwrite manifest when push image without arch
		for _, specArch := range spec.Archs {
			nameWithArch := fmt.Sprintf("%s-%s:%s", name, specArch, tag)
			if tkeImagesSet.Has(nameWithArch) { // check whether has image with any arch
				continue
			}
		}

		// only push image
		err = t.docker.PushImage(image)
		if err != nil {
			return err
		}
	} else {
		// when arch != "", need create manifest list
		manifestName = fmt.Sprintf("%s:%s", name, tag)

		err = t.docker.PushImageWithArch(image, manifestName, arch, "", false)
		if err != nil {
			return err
		}

		if arch == spec.Arm64 {
			err = t.docker.PushArm64Variants(image, name, tag)
			if err != nil {
				return err
			}
		}
	}

	t.log.Infof("upload %s to registry success", image)
	manifestsChan <- manifestName
	return nil
}

func (t *TKE) preInstallHook(ctx context.Context) error {
	return t.execHook(constants.PreInstallHook)
}

func (t *TKE) postClusterReadyHook(ctx context.Context) error {
	return t.execHook(constants.PostClusterReadyHook)
}

func (t *TKE) postInstallHook(ctx context.Context) error {
	return t.execHook(constants.PostInstallHook)
}

func (t *TKE) execHook(filename string) error {
	t.log.Infof("Execute hook script %s", filename)
	cmd := exec.Command(filename)
	cmd.Stdout = t.log
	cmd.Stderr = t.log
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (t *TKE) getKubeconfig() (*api.Config, error) {
	host, err := t.Cluster.Host()
	if err != nil {
		return nil, err
	}

	return kubeconfig.CreateWithToken(host,
		t.Cluster.Name,
		"admin",
		t.Cluster.ClusterCredential.CACert,
		*t.Cluster.ClusterCredential.Token,
	), nil
}

func (t *TKE) setGlobalClusterHosts(ctx context.Context) error {
	domains := []string{
		t.Para.Config.Registry.Domain(),
		t.Cluster.Spec.TenantID + "." + t.Para.Config.Registry.Domain(),
	}

	ip := t.Cluster.Spec.Machines[0].IP
	if t.Para.Config.HA != nil && len(t.Para.Config.HA.VIP()) > 0 {
		ip = t.Para.Config.HA.VIP()
	}

	for _, machine := range t.Cluster.Spec.Machines {
		sshConfig := &ssh.Config{
			User:       machine.Username,
			Host:       machine.IP,
			Port:       int(machine.Port),
			Password:   string(machine.Password),
			PrivateKey: machine.PrivateKey,
			PassPhrase: machine.PassPhrase,
		}
		s, err := ssh.New(sshConfig)
		if err != nil {
			return err
		}
		for _, one := range domains {
			remoteHosts := hosts.RemoteHosts{Host: one, SSH: s}
			err := remoteHosts.Set(ip)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *TKE) writeKubeconfig(ctx context.Context) error {
	cfg, err := t.getKubeconfig()
	if err != nil {
		return err
	}
	data, err := runtime.Encode(clientcmdlatest.Codec, cfg)
	if err != nil {
		return err
	}
	_ = ioutil.WriteFile(constants.KubeconfigFile, data, 0644)
	_ = os.MkdirAll("/root/.kube", 0755)
	return ioutil.WriteFile("/root/.kube/config", data, 0644)
}

func (t *TKE) patchPlatformVersion(ctx context.Context) error {
	if t.globalClient == nil {
		return errors.New("can't get cluster client")
	}

	tkeVersion, _, err := platformutil.GetPlatformVersionsFromClusterInfo(ctx, t.globalClient)
	if err != nil {
		return err
	}
	if len(tkeVersion) == 0 {
		log.Infof("set platform version to %s", spec.TKEVersion)
	} else {
		log.Infof("patch platform version from %s to %s", tkeVersion, spec.TKEVersion)
	}
	if tkeVersion == spec.TKEVersion {
		log.Info("skip patch platform version, current installer version is equal to platform version")
		return nil
	}

	versionsByte, err := json.Marshal(spec.K8sVersions)
	if err != nil {
		return err
	}
	patchData := map[string]interface{}{
		"data": map[string]interface{}{
			"k8sValidVersions": string(versionsByte),
			"tkeVersion":       spec.TKEVersion,
		},
	}
	return t.patchClusterInfo(ctx, patchData)
}

func (t *TKE) patchClusterInfo(ctx context.Context, patchData interface{}) error {
	patchByte, err := json.Marshal(patchData)
	if err != nil {
		return err
	}
	_, err = t.globalClient.CoreV1().ConfigMaps("kube-public").Patch(ctx, "cluster-info", k8stypes.MergePatchType, patchByte, metav1.PatchOptions{})
	return err
}

func (t *TKE) getLocalInfluxdbAddress() string {
	var influxdbAddress string = fmt.Sprintf("http://%s:30086", t.servers[0])
	if t.Para.Config.HA != nil && len(t.Para.Config.HA.VIP()) > 0 {
		vip := t.Para.Config.HA.VIP()
		influxdbAddress = fmt.Sprintf("http://%s:30086", vip) // influxdb svc must be set as NodePort type, and the nodePort is 30086
	}
	return influxdbAddress
}
