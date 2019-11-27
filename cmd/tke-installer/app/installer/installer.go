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
	"encoding/json"
	"fmt"
	"io/ioutil"
	stdlog "log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/emicklei/go-restful"
	pkgerrors "github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"github.com/thoas/go-funk"
	"gopkg.in/go-playground/validator.v9"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd/api"
	clientcmdlatest "k8s.io/client-go/tools/clientcmd/api/latest"
	certutil "k8s.io/client-go/util/cert"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	kubeaggregatorclientset "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
	tkeclientset "tkestack.io/tke/api/client/clientset/versioned"
	"tkestack.io/tke/api/platform"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/cmd/tke-installer/app/config"
	"tkestack.io/tke/cmd/tke-installer/app/installer/certs"
	"tkestack.io/tke/cmd/tke-installer/app/installer/constants"
	"tkestack.io/tke/cmd/tke-installer/app/installer/images"
	baremetalcluster "tkestack.io/tke/pkg/platform/provider/baremetal/cluster"
	baremetalconfig "tkestack.io/tke/pkg/platform/provider/baremetal/config"
	baremetalconstants "tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
	clusterstrategy "tkestack.io/tke/pkg/platform/registry/cluster"
	platformutil "tkestack.io/tke/pkg/platform/util"
	"tkestack.io/tke/pkg/util"
	"tkestack.io/tke/pkg/util/apiclient"
	"tkestack.io/tke/pkg/util/containerregistry"
	"tkestack.io/tke/pkg/util/hosts"
	"tkestack.io/tke/pkg/util/kubeconfig"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/ssh"

	// import platform schema
	_ "tkestack.io/tke/api/platform/install"
)

const (
	dataDir        = "data"
	clusterFile    = dataDir + "/tke.json"
	clusterLogFile = dataDir + "/tke.log"

	hooksDir             = "hooks"
	preInstallHook       = hooksDir + "/pre-install"
	postClusterReadyHook = hooksDir + "/post-cluster-ready"
	postInstallHook      = hooksDir + "/post-install"

	registryDomain             = "docker.io"
	registryNamespace          = "tkestack"
	imagesFile                 = "images.tar.gz"
	imagesPattern              = registryNamespace + "/*"
	localRegistryImage         = registryNamespace + "/local-tcr:v1.0.0"
	localRegistryPort          = 5000
	localRegistryContainerName = "tcr"
	defaultTeantID             = "default"

	k8sVersion = "1.14.6"
)

// ClusterResource is the REST layer to the Cluster domain
type TKE struct {
	Config  *config.Config           `json:"config"`
	Para    *CreateClusterPara       `json:"para"`
	Cluster *clusterprovider.Cluster `json:"cluster"`
	Step    int                      `json:"step"`

	log              *stdlog.Logger
	steps            []handler
	clusterProviders *sync.Map
	strategy         *clusterstrategy.Strategy
	clusterProvider  clusterprovider.Provider
	process          *ClusterProgress
	isFromRestore    bool

	globalClient kubernetes.Interface
	servers      []string
	namespace    string
}

// CreateClusterPara for create cluster parameter
type CreateClusterPara struct {
	Cluster platformv1.Cluster `json:"cluster"`
	Config  Config             `json:"Config"`
}

// Config is the installer config
type Config struct {
	Basic    Basic     `json:"basic"`
	Auth     Auth      `json:"auth"`
	Registry Registry  `json:"registry"`
	Business *Business `json:"business,omitempty"`
	Monitor  *Monitor  `json:"monitor,omitempty"`
	HA       *HA       `json:"ha,omitempty"`
	Gateway  *Gateway  `json:"gateway,omitempty"`
}

type Basic struct {
	Username string `json:"username"`
	Password []byte `json:"password"`
}

type Auth struct {
	TKEAuth  *TKEAuth  `json:"tke,omitempty"`
	OIDCAuth *OIDCAuth `json:"oidc,omitempty"`
}

type TKEAuth struct {
	TenantID string `json:"tenantID"`
	Username string `json:"username"`
	Password []byte `json:"password"`
}

type OIDCAuth struct {
	IssuerURL string `json:"issuerURL" validate:"required"`
	ClientID  string `json:"clientID" validate:"required"`
	CACert    []byte `json:"caCert"`
}

// Registry for remote registry
type Registry struct {
	TKERegistry        *TKERegistry        `json:"tke,omitempty"`
	ThirdPartyRegistry *ThirdPartyRegistry `json:"thirdParty,omitempty"`
}

func (r *Registry) UseDevRegistry() bool {
	return r.ThirdPartyRegistry != nil &&
		r.ThirdPartyRegistry.Domain == registryDomain &&
		r.ThirdPartyRegistry.Namespace == registryNamespace
}

func (r *Registry) Domain() string {
	if r.ThirdPartyRegistry != nil { // first use third party when both set
		return r.ThirdPartyRegistry.Domain
	}
	return r.TKERegistry.Domain
}

func (r *Registry) Namespace() string {
	if r.ThirdPartyRegistry != nil {
		return r.ThirdPartyRegistry.Namespace
	}
	return r.TKERegistry.Namespace
}

type TKERegistry struct {
	Domain    string `json:"domain" validate:"hostname_rfc1123"`
	Namespace string `json:"namespace"`
	Username  string `json:"username"`
	Password  []byte `json:"password"`
}

type ThirdPartyRegistry struct {
	Domain    string `json:"domain" validate:"hostname_rfc1123"`
	Namespace string `json:"namespace" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Password  []byte `json:"password" validate:"required"`
}

type Business struct {
}

type Monitor struct {
	ESMonitor       *ESMonitor       `json:"es,omitempty"`
	InfluxDBMonitor *InfluxDBMonitor `json:"influxDB,omitempty"`
}

type ESMonitor struct {
	URL      string `json:"url" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password []byte `json:"password" validate:"required"`
}

type InfluxDBMonitor struct {
	LocalInfluxDBMonitor    *LocalInfluxDBMonitor    `json:"local,omitempty"`
	ExternalInfluxDBMonitor *ExternalInfluxDBMonitor `json:"external,omitempty"`
}

type LocalInfluxDBMonitor struct {
}

type ExternalInfluxDBMonitor struct {
	URL      string `json:"url" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password []byte `json:"password" validate:"required"`
}

type HA struct {
	TKEHA        *TKEHA        `json:"tke,omitempty"`
	ThirdPartyHA *ThirdPartyHA `json:"thirdParty,omitempty"`
}

func (ha *HA) VIP() string {
	if ha.TKEHA != nil {
		return ha.TKEHA.VIP
	}
	return ha.ThirdPartyHA.VIP
}

type TKEHA struct {
	VIP string `json:"vip" validate:"required"`
}

type ThirdPartyHA struct {
	VIP string `json:"vip" validate:"required"`
}

type Gateway struct {
	Domain string `json:"domain"`
	Cert   Cert   `json:"cert"`
}

type Cert struct {
	SelfSignedCert *SelfSignedCert `json:"selfSigned,omitempty"`
	ThirdPartyCert *ThirdPartyCert `json:"thirdParty,omitempty"`
}

type SelfSignedCert struct {
}

type ThirdPartyCert struct {
	Certificate []byte `json:"certificate" validate:"required"`
	PrivateKey  []byte `json:"privateKey" validate:"required"`
}

type Keepalived struct {
	VIP string `json:"vip,omitempty"`
}

// ClusterProgress use for findClusterProgress
type ClusterProgress struct {
	Status     ClusterProgressStatus `json:"status"`
	Data       string                `json:"data"`
	URL        string                `json:"url,omitempty"`
	Username   string                `json:"username,omitempty"`
	Password   []byte                `json:"password,omitempty"`
	CACert     []byte                `json:"caCert,omitempty"`
	Hosts      []string              `json:"hosts,omitempty"`
	Servers    []string              `json:"servers,omitempty"`
	Kubeconfig []byte                `json:"kubeconfig,omitempty"`
}

type handler struct {
	Name string
	Func func() error
}

// ClusterProgressStatus use for ClusterProgress
type ClusterProgressStatus string

const (
	statusUnknown = "Unknown"
	statusDoing   = "Doing"
	statusSuccess = "Success"
	statusFailed  = "Failed"
)

const (
	pluginConfigFile = "provider/baremetal/conf/config.yaml"
)

func NewTKE(config *config.Config) *TKE {
	c := new(TKE)

	c.Config = config
	c.Para = new(CreateClusterPara)
	c.Cluster = new(clusterprovider.Cluster)

	_ = os.MkdirAll(path.Dir(clusterLogFile), 0755)
	f, err := os.OpenFile(clusterLogFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0744)
	if err != nil {
		log.Fatal(err.Error())
	}
	c.log = stdlog.New(f, "", stdlog.LstdFlags)

	if !config.Force {
		data, err := ioutil.ReadFile(clusterFile)
		if err == nil {
			log.Infof("read %q success", clusterFile)
			err = json.Unmarshal(data, c)
			if err != nil {
				log.Warnf("load tke data error:%s", err)
			}
			log.Infof("load tke data success")
			c.isFromRestore = true
		}
	}

	c.clusterProviders = new(sync.Map)
	c.process = new(ClusterProgress)
	c.process.Status = statusUnknown

	return c
}

func (t *TKE) initSteps() {
	t.steps = append(t.steps, []handler{
		{
			Name: "Execute pre install hook",
			Func: t.preInstallHook,
		},
	}...)

	// UseDockerHub, no need load images, start local tcr and push images
	// TKERegistry load images && start local registry && push images to local registry
	// && deploy tke-registry-api && push images to tke-registry
	// ThirdPartyRegistry load images && push images
	if !t.Para.Config.Registry.UseDevRegistry() {
		t.steps = append(t.steps, []handler{
			{
				Name: "Load images",
				Func: t.loadImages,
			},
		}...)
	}

	// if both set, don't setup local registry
	if t.Para.Config.Registry.ThirdPartyRegistry == nil &&
		t.Para.Config.Registry.TKERegistry != nil {
		t.steps = append(t.steps, []handler{
			{
				Name: "Setup local registry",
				Func: t.setupLocalRegistry,
			},
		}...)
	}

	if !t.Para.Config.Registry.UseDevRegistry() {
		t.steps = append(t.steps, []handler{
			{
				Name: fmt.Sprintf("Push images to %s/%s", t.Para.Config.Registry.Domain(), t.Para.Config.Registry.Namespace()),
				Func: t.pushImages,
			},
		}...)
	}

	t.steps = append(t.steps, []handler{
		{
			Name: "Generate certificates for TKE components",
			Func: t.generateCertificates,
		},
		{
			Name: "Prepare front proxy certificates",
			Func: t.prepareFrontProxyCertificates,
		},
		{
			Name: "Create global cluster",
			Func: t.createGlobalCluster,
		},
		{
			Name: "Execute post deploy hook",
			Func: t.postClusterReadyHook,
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
	}...)

	if t.Para.Config.Auth.TKEAuth != nil {
		t.steps = append(t.steps, []handler{
			{
				Name: "Install tke-auth",
				Func: t.installTKEAuth,
			},
		}...)
	}

	t.steps = append(t.steps, []handler{
		{
			Name: "Install tke-platform-api",
			Func: t.installTKEPlatformAPI,
		},
		{
			Name: "Install tke-platform-controller",
			Func: t.installTKEPlatformController,
		},
	}...)

	if t.Para.Config.Business != nil {
		t.steps = append(t.steps, []handler{
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
			t.steps = append(t.steps, []handler{
				{
					Name: "Install InfluxDB",
					Func: t.installInfluxDB,
				},
			}...)
		}
		t.steps = append(t.steps, []handler{
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

	if t.Para.Config.Registry.TKERegistry != nil {
		t.steps = append(t.steps, []handler{
			{
				Name: "Install tke-registry-api",
				Func: t.installTKERegistryAPI,
			},
		}...)
	}

	if t.Para.Config.Gateway != nil {
		t.steps = append(t.steps, []handler{
			{
				Name: "Install tke-gateway",
				Func: t.installTKEGateway,
			},
		}...)
	}

	if t.Para.Config.HA != nil && t.Para.Config.HA.TKEHA != nil {
		t.steps = append(t.steps, []handler{
			{
				Name: "Install keepalived",
				Func: t.installKeepalived,
			},
		}...)
	}

	t.steps = append(t.steps, []handler{
		{
			Name: "Register tke api into global cluster",
			Func: t.registerAPI,
		},
		{
			Name: "Import resource to TKE platform",
			Func: t.importResource,
		},
	}...)

	if t.Para.Config.Registry.ThirdPartyRegistry == nil &&
		t.Para.Config.Registry.TKERegistry != nil {
		t.steps = append(t.steps, []handler{
			{
				Name: "Prepare push images to TKE registry",
				Func: t.preparePushImagesToTKERegistry,
			},
			{
				Name: "Push images to registry",
				Func: t.pushImages,
			},
		}...)
	}

	t.steps = append(t.steps, []handler{
		{
			Name: "Write kubeconfig",
			Func: t.writeKubeconfig,
		},
		{
			Name: "Execute post deploy hook",
			Func: t.postInstallHook,
		},
	}...)
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
			return pkgerrors.New(statusErr.String())
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
		Reads(CreateClusterPara{}).Writes(CreateClusterPara{}))

	ws.Route(ws.PUT("{name}/retry").To(t.retryCreateCluster))

	ws.Route(ws.GET("{name}").To(t.findCluster).
		Writes(CreateClusterPara{}))

	ws.Route(ws.GET("{name}/progress").To(t.findClusterProgress))

	return ws
}

func (t *TKE) completeWithProvider() {
	clusterProvider, err := baremetalcluster.NewProvider()
	if err != nil {
		panic(err)
	}
	t.clusterProvider = clusterProvider
	t.clusterProviders.Store(clusterProvider.Name(), clusterProvider)
}

// Global Filter
func globalLogging(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	now := time.Now()

	reqBytes, err := httputil.DumpRequest(req.Request, true)
	if err != nil {
		_ = resp.WriteError(http.StatusInternalServerError, err)
		return
	}

	log.Infof("raw http request:\n%s", reqBytes)
	chain.ProcessFilter(req, resp)
	log.Infof("%s %s %v", req.Request.Method, req.Request.URL, time.Since(now))
}

func (t *TKE) prepare() errors.APIStatus {
	t.SetConfigDefault(&t.Para.Config)
	statusError := t.ValidateConfig(t.Para.Config)
	if statusError != nil {
		return statusError
	}

	t.SetClusterDefault(&t.Para.Cluster, &t.Para.Config)

	// mock platform api
	t.completeWithProvider()
	t.strategy = clusterstrategy.NewStrategy(t.clusterProviders, nil)

	ctx := request.WithUser(context.Background(), &user.DefaultInfo{Name: defaultTeantID})

	v1Cluster := &t.Para.Cluster

	// PrepareForCreate
	platformCluster := new(platform.Cluster)
	err := platform.Scheme.Convert(v1Cluster, platformCluster, nil)
	if err != nil {
		return errors.NewInternalError(err)
	}
	t.strategy.PrepareForCreate(ctx, platformCluster)

	// Validate
	kinds, _, err := platform.Scheme.ObjectKinds(v1Cluster)
	if err != nil {
		return errors.NewInternalError(err)
	}
	allErrs := t.strategy.Validate(ctx, platformCluster)
	if len(allErrs) > 0 {
		return errors.NewInvalid(kinds[0].GroupKind(), v1Cluster.GetName(), allErrs)
	}

	err = platform.Scheme.Convert(platformCluster, v1Cluster, nil)
	if err != nil {
		return errors.NewInternalError(err)
	}

	t.Cluster.Cluster = *v1Cluster
	t.backup()

	return nil
}

func (t *TKE) SetConfigDefault(config *Config) {
	if config.Registry.TKERegistry != nil {
		config.Registry.TKERegistry.Namespace = "library"
		config.Registry.TKERegistry.Username = config.Basic.Username
		config.Registry.TKERegistry.Password = config.Basic.Password
	}
	if config.Auth.TKEAuth != nil {
		config.Auth.TKEAuth.TenantID = defaultTeantID
		config.Auth.TKEAuth.Username = config.Basic.Username
		config.Auth.TKEAuth.Password = config.Basic.Password
	}

}
func (t *TKE) SetClusterDefault(cluster *platformv1.Cluster, config *Config) {
	cluster.Name = "global"
	cluster.Spec.DisplayName = "TKE"

	cluster.Spec.TenantID = defaultTeantID
	if t.Para.Config.Auth.TKEAuth != nil {
		cluster.Spec.TenantID = t.Para.Config.Auth.TKEAuth.TenantID
	}
	cluster.Spec.Version = k8sVersion

	if config.HA != nil && config.HA.ThirdPartyHA != nil {
		cluster.Status.Addresses = append(cluster.Status.Addresses, platformv1.ClusterAddress{
			Type: platformv1.AddressAdvertise,
			Host: config.HA.VIP(),
			Port: 6443,
		})
	}
}

func (t *TKE) ValidateConfig(config Config) *errors.StatusError {
	validate := validator.New()
	err := validate.Struct(config)
	if err != nil {
		return errors.NewBadRequest(err.Error())
	}

	if config.Gateway != nil || config.Registry.TKERegistry != nil {
		if config.Basic.Username == "" || config.Basic.Password == nil {
			return errors.NewBadRequest("username or password required when enabled gateway or registry")
		}
	}

	if config.Auth.TKEAuth == nil && config.Auth.OIDCAuth == nil {
		return errors.NewBadRequest("tke auth or oidc auth required")
	}

	if config.Registry.TKERegistry == nil && config.Registry.ThirdPartyRegistry == nil {
		return errors.NewBadRequest("tke registry or third party registry required")
	}

	if config.Registry.ThirdPartyRegistry != nil {
		cmd := exec.Command("docker", "login",
			"--username", config.Registry.ThirdPartyRegistry.Username,
			"--password", string(config.Registry.ThirdPartyRegistry.Password),
			config.Registry.ThirdPartyRegistry.Domain,
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			if _, ok := err.(*exec.ExitError); ok {
				return errors.NewBadRequest(string(out))
			}
			return errors.NewInternalError(err)
		}
	}

	if config.Monitor.InfluxDBMonitor != nil && config.Monitor.ESMonitor != nil {
		return errors.NewBadRequest("influxdb or es only had one")
	}

	var dnsNames []string
	if config.Gateway != nil && config.Gateway.Domain != "" {
		dnsNames = append(dnsNames, config.Gateway.Domain)
	}
	if config.Registry.TKERegistry != nil {
		dnsNames = append(dnsNames, config.Registry.TKERegistry.Domain, "*."+config.Registry.TKERegistry.Domain)
	}
	if config.Gateway.Cert.ThirdPartyCert != nil {
		statusError := t.ValidateCertAndKey(config.Gateway.Cert.ThirdPartyCert.Certificate,
			config.Gateway.Cert.ThirdPartyCert.PrivateKey, dnsNames)
		if statusError != nil {
			return statusError
		}
	}

	return nil
}

func (t *TKE) ValidateCertAndKey(certificate []byte, privateKey []byte, dnsNames []string) *errors.StatusError {
	if (certificate != nil && privateKey == nil) || (certificate == nil && privateKey != nil) {
		return errors.NewBadRequest("certificate and privateKey must offer together")
	}

	if certificate != nil {
		cert, err := tls.X509KeyPair(certificate, privateKey)
		if err != nil {
			return errors.NewBadRequest(err.Error())
		}
		if len(cert.Certificate) != 1 {
			return errors.NewBadRequest("certificate must only has one cert")
		}
		certs1, err := certutil.ParseCertsPEM(certificate)
		if err != nil {
			return errors.NewBadRequest(err.Error())
		}
		for _, one := range dnsNames {
			if !funk.Contains(certs1[0].DNSNames, one) {
				return errors.NewBadRequest(fmt.Sprintf("certificate DNSNames must contains %v", one))
			}
		}
	}

	return nil
}

func (t *TKE) initProviderConfig() error {
	c, err := baremetalconfig.New(pluginConfigFile)
	if err != nil {
		return err
	}
	c.Registry.Domain = t.Para.Config.Registry.Domain()
	if t.Para.Config.Registry.ThirdPartyRegistry == nil &&
		t.Para.Config.Registry.TKERegistry != nil {
		ip := t.Cluster.Spec.Machines[0].IP
		if t.Para.Config.HA != nil {
			ip = t.Para.Config.HA.VIP()
		}
		c.Registry.IP = ip
	}

	return c.Save(pluginConfigFile)
}

func (t *TKE) createCluster(req *restful.Request, rsp *restful.Response) {
	apiStatus := func() errors.APIStatus {
		err := req.ReadEntity(t.Para)
		if err != nil {
			return errors.NewBadRequest(err.Error())
		}
		if err := t.prepare(); err != nil {
			return err
		}
		go t.do()

		return nil
	}()

	if apiStatus != nil {
		_ = rsp.WriteHeaderAndJson(int(apiStatus.Status().Code), apiStatus.Status(), restful.MIME_JSON)
	} else {
		_ = rsp.WriteEntity(t.Para)
	}
}

func (t *TKE) retryCreateCluster(req *restful.Request, rsp *restful.Response) {
	go t.do()
	_ = rsp.WriteEntity(nil)
}

func (t *TKE) findCluster(request *restful.Request, response *restful.Response) {
	apiStatus := func() errors.APIStatus {
		clusterName := request.PathParameter("name")
		if t.Cluster == nil {
			return errors.NewBadRequest("no cluater available")
		}
		if t.Cluster.Name != clusterName {
			return errors.NewNotFound(platform.Resource("Cluster"), clusterName)
		}

		return nil
	}()

	if apiStatus != nil {
		_ = response.WriteHeaderAndJson(int(apiStatus.Status().Code), apiStatus.Status(), restful.MIME_JSON)
	} else {
		_ = response.WriteEntity(&CreateClusterPara{
			Cluster: t.Cluster.Cluster,
			Config:  t.Para.Config,
		})
	}
}

func (t *TKE) findClusterProgress(request *restful.Request, response *restful.Response) {
	var err error
	var data []byte
	apiStatus := func() errors.APIStatus {
		clusterName := request.PathParameter("name")
		if t.Cluster == nil {
			return errors.NewBadRequest("no cluater available")
		}
		if t.Cluster.Name != clusterName {
			return errors.NewNotFound(platform.Resource("Cluster"), clusterName)
		}
		data, err = ioutil.ReadFile(clusterLogFile)
		if err != nil {
			return errors.NewInternalError(err)
		}
		t.process.Data = string(data)

		return nil
	}()

	if apiStatus != nil {
		response.WriteHeaderAndJson(int(apiStatus.Status().Code), apiStatus.Status(), restful.MIME_JSON)
	} else {
		if t.process.Status == statusSuccess {
			if t.Para.Config.Gateway != nil {
				var host string
				if t.Para.Config.Gateway.Domain != "" {
					host = t.Para.Config.Gateway.Domain
				} else if t.Para.Config.HA != nil {
					host = t.Para.Config.HA.VIP()
				} else {
					host = t.Para.Cluster.Spec.Machines[0].IP
				}
				t.process.URL = fmt.Sprintf("http://%s", host)

				t.process.Username = t.Para.Config.Basic.Username
				t.process.Password = t.Para.Config.Basic.Password

				if t.Para.Config.Gateway.Cert.SelfSignedCert != nil {
					t.process.CACert, _ = ioutil.ReadFile(constants.CACrtFile)
				}

				if t.Para.Config.Gateway.Domain != "" {
					t.process.Hosts = append(t.process.Hosts, t.Para.Config.Gateway.Domain)
				}

				cfg, _ := t.getKubeconfig()
				t.process.Kubeconfig, _ = runtime.Encode(clientcmdlatest.Codec, cfg)
			}

			if t.Para.Config.Registry.TKERegistry != nil {
				t.process.Hosts = append(t.process.Hosts, t.Para.Config.Registry.TKERegistry.Domain)
			}

			t.process.Servers = t.servers
			if t.Para.Config.HA != nil {
				t.process.Servers = append(t.process.Servers, t.Para.Config.HA.VIP())
			}
		}
		response.WriteEntity(t.process)
	}
}

func (t *TKE) do() {
	start := time.Now()

	containerregistry.Init(t.Para.Config.Registry.Domain(), t.Para.Config.Registry.Namespace())
	t.initSteps()

	if t.Step == 0 {
		t.log.Print("===>starting install task")
		t.process.Status = statusDoing
	}

	if t.runAfterClusterReady() {
		t.initDataForDeployTKE()
	}

	for t.Step < len(t.steps) {
		t.log.Printf("%d.%s doing", t.Step, t.steps[t.Step].Name)
		start := time.Now()
		err := t.steps[t.Step].Func()
		if err != nil {
			t.process.Status = statusFailed
			t.log.Printf("%d.%s [Failed] [%fs] error %s", t.Step, t.steps[t.Step].Name, time.Since(start).Seconds(), err)
			return
		}
		t.log.Printf("%d.%s [Success] [%fs]", t.Step, t.steps[t.Step].Name, time.Since(start).Seconds())

		t.Step++
		t.backup()
	}

	t.process.Status = statusSuccess
	t.log.Printf("===>install task [Sucesss] [%fs]", time.Since(start).Seconds())
}

func (t *TKE) runAfterClusterReady() bool {
	return t.Cluster.Status.Phase == platformv1.ClusterRunning
}

func (t *TKE) generateCertificates() error {
	var dnsNames []string
	if t.Para.Config.Gateway != nil && t.Para.Config.Gateway.Domain != "" {
		dnsNames = append(dnsNames, t.Para.Config.Gateway.Domain)
	}
	if t.Para.Config.Registry.TKERegistry != nil {
		dnsNames = append(dnsNames, t.Para.Config.Registry.TKERegistry.Domain, "*."+t.Para.Config.Registry.TKERegistry.Domain)
	}

	ips := []net.IP{net.ParseIP("127.0.0.1")}
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
	return certs.Generate(dnsNames, ips)
}

func (t *TKE) prepareFrontProxyCertificates() error {
	if t.Cluster.Spec.APIServerExtraArgs == nil {
		t.Cluster.Spec.APIServerExtraArgs = make(map[string]string)
	}
	t.Cluster.Spec.APIServerExtraArgs["proxy-client-cert-file"] = "/etc/kubernetes/admin.crt"
	t.Cluster.Spec.APIServerExtraArgs["proxy-client-key-file"] = "/etc/kubernetes/admin.key"
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
		err = s.CopyFile(constants.AdminCrtFile, "/etc/kubernetes/admin.crt")
		if err != nil {
			return err
		}
		err = s.CopyFile(constants.AdminKeyFile, "/etc/kubernetes/admin.key")
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *TKE) createGlobalCluster() error {
	err := t.initProviderConfig()
	if err != nil {
		return err
	}
	// do again like platform controller
	t.completeWithProvider()

	if t.Cluster.ClusterCredential.Name == "" { // set ClusterCredential default value
		t.Cluster.ClusterCredential = platformv1.ClusterCredential{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("cc-%s", t.Cluster.Name),
			},
			TenantID:    t.Cluster.Spec.TenantID,
			ClusterName: t.Cluster.Name,
		}
	}

	var errCount int
	for {
		start := time.Now()
		resp, err := t.clusterProvider.OnInitialize(*t.Cluster)
		if err != nil {
			return err
		}
		t.Cluster = &resp
		t.backup()
		condition := resp.Status.Conditions[len(resp.Status.Conditions)-1]
		switch condition.Status {
		case platformv1.ConditionFalse: // means current condition run into error
			t.log.Printf("OnInitialize.%s [Failed] [%fs] reason: %s message: %s retry: %d",
				condition.Type, time.Since(start).Seconds(),
				condition.Reason, condition.Message, errCount)
			if errCount >= 10 {
				return pkgerrors.New("the retry limit has been reached")
			}

			errCount++
		case platformv1.ConditionUnknown: // means has next condition need to process
			condition = resp.Status.Conditions[len(resp.Status.Conditions)-2]
			t.log.Printf("OnInitialize.%s [Success] [%fs]", condition.Type, time.Since(start).Seconds())

			t.Cluster.Status.Reason = ""
			t.Cluster.Status.Message = ""
			errCount = 0
		case platformv1.ConditionTrue: // means all condition is done
			if resp.Status.Phase != platformv1.ClusterRunning {
				return pkgerrors.Errorf("OnInitialize.%s, no next condition but cluster is not running!", condition.Type)
			}
			return t.initDataForDeployTKE()
		default:
			return pkgerrors.Errorf("unknown condition status %s", condition.Status)
		}
		time.Sleep(5 * time.Second)
	}
}

func (t *TKE) backup() error {
	data, _ := json.MarshalIndent(t, "", " ")
	return ioutil.WriteFile(clusterFile, data, 0777)
}

func (t *TKE) loadImages() error {
	if _, err := os.Stat(imagesFile); err != nil {
		return err
	}
	cmd := exec.Command("docker", "load", "-i", imagesFile)
	cmd.Stdout = t.log.Writer()
	cmd.Stderr = t.log.Writer()
	err := cmd.Run()
	if err != nil {
		return pkgerrors.Wrap(err, "docker load error")
	}

	cmd = exec.Command("sh", "-c",
		fmt.Sprintf("docker images --format='{{.Repository}}:{{.Tag}}' --filter='reference=%s'", imagesPattern),
	)
	out, err := cmd.Output()
	if err != nil {
		return pkgerrors.Wrap(err, "docker images error")
	}
	tkeImages := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, image := range tkeImages {
		imageNames := strings.Split(image, "/")
		if len(imageNames) != 2 {
			t.log.Printf("invalid image name:name=%s", image)
			continue
		}
		nameAndTag := strings.Split(imageNames[1], ":")
		if nameAndTag[1] == "<none>" {
			t.log.Printf("skip invalid tag:name=%s", image)
			continue
		}
		target := fmt.Sprintf("%s/%s/%s", t.Para.Config.Registry.Domain(), t.Para.Config.Registry.Namespace(), imageNames[1])

		cmd := exec.Command("docker", "tag", image, target)
		cmd.Stdout = t.log.Writer()
		cmd.Stderr = t.log.Writer()
		err = cmd.Run()
		if err != nil {
			return pkgerrors.Wrap(err, "docker tag error:%s")
		}
	}

	return nil
}

func (t *TKE) setupLocalRegistry() error {
	server := t.Para.Config.Registry.Domain()

	err := t.startLocalRegistry()
	if err != nil {
		return pkgerrors.Wrap(err, "start local registry error")
	}

	// for push image to local registry
	localHosts := hosts.LocalHosts{Host: server, File: "hosts"}
	err = localHosts.Set("127.0.0.1")
	if err != nil {
		return err
	}

	// for pull image from local registry on node
	ip, err := util.GetExternalIP()
	if err != nil {
		return pkgerrors.Wrap(err, "get external ip error")
	}
	err = t.injectRemoteHosts([]string{t.Para.Config.Registry.Domain()}, ip)
	if err != nil {
		return pkgerrors.Wrap(err, "inject remote hosts error")
	}

	return nil
}

func (t *TKE) injectRemoteHosts(host []string, ip string) error {
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
		for _, one := range host {
			remoteHosts := hosts.RemoteHosts{Host: one, SSH: s}
			err = remoteHosts.Set(ip)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *TKE) startLocalRegistry() error {
	cmd := exec.Command("sh", "-c",
		fmt.Sprintf("docker rm -f %s", localRegistryContainerName),
	)
	cmd.Run()

	cmd = exec.Command("sh", "-c",
		fmt.Sprintf("docker run -d -v `pwd`/registry:/var/lib/registry -p 80:%d --restart always --name %s %s",
			localRegistryPort, localRegistryContainerName, localRegistryImage),
	)
	cmd.Stdout = t.log.Writer()
	cmd.Stderr = t.log.Writer()
	err := cmd.Run()
	if err != nil {
		return pkgerrors.Wrap(err, "docker run error")
	}

	return nil
}

// GetClientset return the clientset of component
func (t *TKE) GetClientset(name string, dnsDomain string) (tkeclientset.Interface, error) {
	token, _ := ioutil.ReadFile(constants.TokenFile)
	caCert, _ := ioutil.ReadFile(constants.CACrtFile)

	return apiclient.GetPlatformClientset(fmt.Sprintf("https://%s.%s", name, dnsDomain), string(token), caCert)
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
	t.globalClient, err = platformutil.BuildExternalClientSetNoStatus(&t.Cluster.Cluster, &t.Cluster.ClusterCredential)
	if err != nil {
		return err
	}

	for _, address := range t.Cluster.Status.Addresses {
		if address.Type == platformv1.AddressReal {
			t.servers = append(t.servers, address.Host)
		}
	}

	t.namespace = "tke"

	return nil
}

func (t *TKE) createNamespace() error {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: t.namespace,
		},
	}

	return apiclient.CreateOrUpdateNamespace(t.globalClient, ns)
}

func (t *TKE) prepareCertificates() error {
	caCrt, err := ioutil.ReadFile(constants.CACrtFile)
	if err != nil {
		return err
	}
	caKey, err := ioutil.ReadFile(constants.CAKeyFile)
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

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "certs",
			Namespace: t.namespace,
		},
		Data: map[string]string{
			"etcd-ca.crt": string(t.Cluster.ClusterCredential.ETCDCACert),
			"etcd.crt":    string(t.Cluster.ClusterCredential.ETCDAPIClientCert),
			"etcd.key":    string(t.Cluster.ClusterCredential.ETCDAPIClientKey),
			"ca.crt":      string(caCrt),
			"ca.key":      string(caKey),
			"server.crt":  string(serverCrt),
			"server.key":  string(serverKey),
			"admin.crt":   string(adminCrt),
			"admin.key":   string(adminKey),
		},
	}

	if t.Para.Config.Auth.OIDCAuth != nil {
		cm.Data["oidc-ca.crt"] = string(t.Para.Config.Auth.OIDCAuth.CACert)
	}

	cm.Data["password.csv"] = fmt.Sprintf("%s,admin,1,administrator", ksuid.New().String())
	cm.Data["token.csv"] = fmt.Sprintf("%s,admin,1,administrator", ksuid.New().String())

	for k, v := range cm.Data {
		err := ioutil.WriteFile(path.Join(dataDir, k), []byte(v), 0644)
		if err != nil {
			return err
		}
	}

	return apiclient.CreateOrUpdateConfigMap(t.globalClient, cm)
}

func (t *TKE) prepareBaremetalProviderConfig() error {
	configMaps := []struct {
		Name string
		File string
	}{
		{
			Name: "provider-config",
			File: baremetalconstants.ConfDir + "config.yaml",
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
	}
	for _, one := range configMaps {
		err := apiclient.CreateOrUpdateConfigMapFromFile(t.globalClient,
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

func (t *TKE) installKeepalived() error {
	err := apiclient.CreateResourceWithDir(t.globalClient, "manifests/keepalived/*.yaml",
		map[string]interface{}{
			"Image": images.Get().Keepalived.FullName(),
			"VIP":   t.Para.Config.HA.TKEHA.VIP,
		})
	if err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDaemonset(t.globalClient, t.namespace, "tke-gateway")
		if err != nil {
			return false, nil
		}
		if t.Para.Config.HA != nil && t.Para.Config.HA.TKEHA != nil {
			t.Cluster.Status.Addresses = append(t.Cluster.Status.Addresses, platformv1.ClusterAddress{
				Type: platformv1.AddressAdvertise,
				Host: t.Para.Config.HA.TKEHA.VIP,
				Port: 6443,
			})
		}
		return ok, nil
	})
}

func (t *TKE) installTKEGateway() error {
	option := map[string]interface{}{
		"Image":            images.Get().TKEGateway.FullName(),
		"OIDCClientSecret": t.readOrGenerateString(constants.OIDCClientSecretFile),
		"SelfSigned":       t.Para.Config.Gateway.Cert.SelfSignedCert != nil,
		"EnableRegistry":   t.Para.Config.Registry.TKERegistry != nil,
		"EnableAuth":       t.Para.Config.Auth.TKEAuth != nil,
		"EnableMonitor":    t.Para.Config.Monitor != nil,
		"EnableBusiness":   t.Para.Config.Business != nil,
	}
	if t.Para.Config.Registry.TKERegistry != nil {
		option["RegistryDomainSuffix"] = t.Para.Config.Registry.TKERegistry.Domain
	}
	if t.Para.Config.Auth.TKEAuth != nil {
		option["TenantID"] = t.Para.Config.Auth.TKEAuth.TenantID
	}
	if t.Para.Config.Gateway.Cert.ThirdPartyCert != nil {
		option["ServerCrt"] = t.Para.Config.Gateway.Cert.ThirdPartyCert.Certificate
		option["ServerKey"] = t.Para.Config.Gateway.Cert.ThirdPartyCert.PrivateKey
	}
	err := apiclient.CreateResourceWithDir(t.globalClient, "manifests/tke-gateway/*.yaml", option)
	if err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDaemonset(t.globalClient, t.namespace, "tke-gateway")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) installETCD() error {
	return apiclient.CreateResourceWithDir(t.globalClient, "manifests/etcd/*.yaml",
		map[string]interface{}{
			"Servers": t.servers,
		})
}

func (t *TKE) installTKEAuth() error {
	redirectHosts := t.servers
	redirectHosts = append(redirectHosts, "tke-gateway")
	if t.Para.Config.Gateway != nil && t.Para.Config.Gateway.Domain != "" {
		redirectHosts = append(redirectHosts, t.Para.Config.Gateway.Domain)
	}
	if t.Para.Config.HA != nil {
		redirectHosts = append(redirectHosts, t.Para.Config.HA.VIP())
	}

	option := map[string]interface{}{
		"Replicas":         t.Config.Replicas,
		"Image":            images.Get().TKEAuth.FullName(),
		"OIDCClientSecret": t.readOrGenerateString(constants.OIDCClientSecretFile),
		"AdminUsername":    t.Para.Config.Auth.TKEAuth.Username,
		"AdminPassword":    string(t.Para.Config.Auth.TKEAuth.Password),
		"TenantID":         t.Para.Config.Auth.TKEAuth.TenantID,
		"RedirectHosts":    redirectHosts,
	}
	err := apiclient.CreateResourceWithDir(t.globalClient, "manifests/tke-auth/*.yaml", option)
	if err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(t.globalClient, t.namespace, "tke-auth")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) installTKEPlatformAPI() error {
	options := map[string]interface{}{
		"Replicas":                      t.Config.Replicas,
		"Image":                         images.Get().TKEPlatformAPI.FullName(),
		"BaremetalClusterProviderImage": images.Get().BaremetalClusterProvider.FullName(),
		"BaremetalMachineProviderImage": images.Get().BaremetalMachineProvider.FullName(),
		"EnableAuth":                    t.Para.Config.Auth.TKEAuth != nil,
	}
	if t.Para.Config.Auth.OIDCAuth != nil {
		options["OIDCClientID"] = t.Para.Config.Auth.OIDCAuth.ClientID
		options["OIDCIssuerURL"] = t.Para.Config.Auth.OIDCAuth.IssuerURL
		options["UseOIDCCA"] = t.Para.Config.Auth.OIDCAuth.CACert != nil
	}
	err := apiclient.CreateResourceWithDir(t.globalClient, "manifests/tke-platform-api/*.yaml", options)
	if err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(t.globalClient, t.namespace, "tke-platform-api")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) installTKEPlatformController() error {
	params := map[string]interface{}{
		"Replicas":                      t.Config.Replicas,
		"Image":                         images.Get().TKEPlatformController.FullName(),
		"BaremetalClusterProviderImage": images.Get().BaremetalClusterProvider.FullName(),
		"BaremetalMachineProviderImage": images.Get().BaremetalMachineProvider.FullName(),
		"ProviderResImage":              images.Get().ProviderRes.FullName(),
		"RegistryDomain":                t.Para.Config.Registry.Domain(),
		"MonitorStorageType":            "",
		"MonitorStorageAddresses":       "",
	}
	if t.Para.Config.Monitor != nil {
		if t.Para.Config.Monitor.InfluxDBMonitor != nil {
			params["MonitorStorageType"] = "influxdb"
			if t.Para.Config.Monitor.InfluxDBMonitor.LocalInfluxDBMonitor != nil {
				params["MonitorStorageAddresses"] = fmt.Sprintf("http://%s:8086", t.servers[0])
			} else if t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor != nil {
				address := t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.URL
				if t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.Username != "" {
					address = address + "&u=" + t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.Username
				}
				if t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.Password != nil {
					address = address + "&p=" + string(t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.Password)
				}
				params["MonitorStorageAddresses"] = address
			}
		} else if t.Para.Config.Monitor.ESMonitor != nil {
			params["MonitorStorageType"] = "elasticsearch"
			address := t.Para.Config.Monitor.ESMonitor.URL
			if t.Para.Config.Monitor.ESMonitor.Username != "" {
				address = address + "&u=" + t.Para.Config.Monitor.ESMonitor.Username
			}
			if t.Para.Config.Monitor.ESMonitor.Password != nil {
				address = address + "&p=" + string(t.Para.Config.Monitor.ESMonitor.Password)
			}
			params["MonitorStorageAddresses"] = address
		}
	}

	if err := apiclient.CreateResourceWithDir(t.globalClient, "manifests/tke-platform-controller/*.yaml", params); err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(t.globalClient, t.namespace, "tke-platform-controller")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) installTKEBusinessAPI() error {
	options := map[string]interface{}{
		"Replicas":                   t.Config.Replicas,
		"Image":                      images.Get().TKEBusinessAPI.FullName(),
		"TenantID":                   t.Para.Config.Auth.TKEAuth.TenantID,
		"Username":                   t.Para.Config.Auth.TKEAuth.Username,
		"SyncProjectsWithNamespaces": t.Config.SyncProjectsWithNamespaces,
		"EnableAuth":                 t.Para.Config.Auth.TKEAuth != nil,
	}
	if t.Para.Config.Auth.OIDCAuth != nil {
		options["OIDCClientID"] = t.Para.Config.Auth.OIDCAuth.ClientID
		options["OIDCIssuerURL"] = t.Para.Config.Auth.OIDCAuth.IssuerURL
		options["UseOIDCCA"] = t.Para.Config.Auth.OIDCAuth.CACert != nil
	}
	err := apiclient.CreateResourceWithDir(t.globalClient, "manifests/tke-business-api/*.yaml", options)
	if err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(t.globalClient, t.namespace, "tke-business-api")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) installTKEBusinessController() error {
	err := apiclient.CreateResourceWithDir(t.globalClient, "manifests/tke-business-controller/*.yaml",
		map[string]interface{}{
			"Replicas": t.Config.Replicas,
			"Image":    images.Get().TKEBusinessController.FullName(),
		})
	if err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(t.globalClient, t.namespace, "tke-business-controller")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) installInfluxDB() error {
	err := apiclient.CreateResourceWithDir(t.globalClient, "manifests/influxdb/*.yaml",
		map[string]interface{}{
			"Image":    images.Get().InfluxDB.FullName(),
			"NodeName": t.servers[0],
		})
	if err != nil {
		return err
	}
	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckStatefulSet(t.globalClient, t.namespace, "influxdb")
		if err != nil || !ok {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) installTKEMonitorAPI() error {
	options := map[string]interface{}{
		"Replicas":   t.Config.Replicas,
		"Image":      images.Get().TKEMonitorAPI.FullName(),
		"EnableAuth": t.Para.Config.Auth.TKEAuth != nil,
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
			options["StoragePassword"] = t.Para.Config.Monitor.ESMonitor.Password
		} else if t.Para.Config.Monitor.InfluxDBMonitor != nil {
			options["StorageType"] = "influxDB"

			if t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor != nil {
				options["StorageAddress"] = t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.URL
				options["StorageUsername"] = t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.Username
				options["StoragePassword"] = t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.Password
			} else if t.Para.Config.Monitor.InfluxDBMonitor.LocalInfluxDBMonitor != nil {
				// todo
				options["StorageAddress"] = fmt.Sprintf("http://%s:8086", t.servers[0])
			}
		}
	}

	if err := apiclient.CreateResourceWithDir(t.globalClient, "manifests/tke-monitor-api/*.yaml", options); err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(t.globalClient, t.namespace, "tke-monitor-api")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) installTKEMonitorController() error {
	params := map[string]interface{}{
		"Replicas": t.Config.Replicas,
		"Image":    images.Get().TKEMonitorController.FullName(),
	}
	if t.Para.Config.Monitor != nil {
		if t.Para.Config.Monitor.ESMonitor != nil {
			params["StorageType"] = "es"
			params["StorageAddress"] = t.Para.Config.Monitor.ESMonitor.URL
			params["StorageUsername"] = t.Para.Config.Monitor.ESMonitor.Username
			params["StoragePassword"] = t.Para.Config.Monitor.ESMonitor.Password
		} else if t.Para.Config.Monitor.InfluxDBMonitor != nil {
			params["StorageType"] = "influxDB"

			if t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor != nil {
				params["StorageAddress"] = t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.URL
				params["StorageUsername"] = t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.Username
				params["StoragePassword"] = t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.Password
			} else if t.Para.Config.Monitor.InfluxDBMonitor.LocalInfluxDBMonitor != nil {
				params["StorageAddress"] = fmt.Sprintf("http://%s:8086", t.servers[0])
			}
		}
	}

	if err := apiclient.CreateResourceWithDir(t.globalClient, "manifests/tke-monitor-controller/*.yaml", params); err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(t.globalClient, t.namespace, "tke-monitor-controller")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) installTKENotifyAPI() error {
	options := map[string]interface{}{
		"Replicas":   t.Config.Replicas,
		"Image":      images.Get().TKENotifyAPI.FullName(),
		"EnableAuth": t.Para.Config.Auth.TKEAuth != nil,
	}
	if t.Para.Config.Auth.OIDCAuth != nil {
		options["OIDCClientID"] = t.Para.Config.Auth.OIDCAuth.ClientID
		options["OIDCIssuerURL"] = t.Para.Config.Auth.OIDCAuth.IssuerURL
		options["UseOIDCCA"] = t.Para.Config.Auth.OIDCAuth.CACert != nil
	}
	err := apiclient.CreateResourceWithDir(t.globalClient, "manifests/tke-notify-api/*.yaml", options)
	if err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(t.globalClient, t.namespace, "tke-notify-api")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) installTKENotifyController() error {
	err := apiclient.CreateResourceWithDir(t.globalClient, "manifests/tke-notify-controller/*.yaml",
		map[string]interface{}{
			"Replicas": t.Config.Replicas,
			"Image":    images.Get().TKENotifyController.FullName(),
		})
	if err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(t.globalClient, t.namespace, "tke-notify-controller")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) installTKERegistryAPI() error {
	options := map[string]interface{}{
		"Image":         images.Get().TKERegistryAPI.FullName(),
		"NodeName":      t.servers[0],
		"AdminUsername": t.Para.Config.Registry.TKERegistry.Username,
		"AdminPassword": string(t.Para.Config.Registry.TKERegistry.Password),
		"EnableAuth":    t.Para.Config.Auth.TKEAuth != nil,
	}
	if t.Para.Config.Auth.OIDCAuth != nil {
		options["OIDCClientID"] = t.Para.Config.Auth.OIDCAuth.ClientID
		options["OIDCIssuerURL"] = t.Para.Config.Auth.OIDCAuth.IssuerURL
		options["UseOIDCCA"] = t.Para.Config.Auth.OIDCAuth.CACert != nil
	}
	err := apiclient.CreateResourceWithDir(t.globalClient, "manifests/tke-registry-api/*.yaml", options)
	if err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(t.globalClient, t.namespace, "tke-registry-api")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) preparePushImagesToTKERegistry() error {
	localHosts := hosts.LocalHosts{Host: t.Para.Config.Registry.Domain()}
	err := localHosts.Set(t.servers[0])
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
			return pkgerrors.New(string(out))
		}
		return err
	}

	return nil
}

func (t *TKE) registerAPI() error {
	caCert, _ := ioutil.ReadFile(constants.CACrtFile)

	restConfig, err := platformutil.GetRestConfig(&t.Cluster.Cluster, &t.Cluster.ClusterCredential)
	if err != nil {
		return err
	}
	client, err := kubeaggregatorclientset.NewForConfig(restConfig)
	if err != nil {
		return err
	}

	svcs := []string{"tke-platform-api"}
	if t.Para.Config.Business != nil {
		svcs = append(svcs, "tke-business-api")
	}
	if t.Para.Config.Monitor != nil {
		svcs = append(svcs, "tke-notify-api", "tke-monitor-api")
	}
	if t.Para.Config.Registry.TKERegistry != nil {
		svcs = append(svcs, "tke-registry-api")
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

		_, err := client.ApiregistrationV1().APIServices().Get(apiService.Name, metav1.GetOptions{})
		if err == nil {
			err := client.ApiregistrationV1().APIServices().Delete(apiService.Name, &metav1.DeleteOptions{})
			if err != nil {
				return err
			}
		}
		if _, err := client.ApiregistrationV1().APIServices().Create(apiService); err != nil {
			if !errors.IsAlreadyExists(err) {
				return err
			}
		}

		err = wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
			a, err := client.ApiregistrationV1().APIServices().Get(apiService.Name, metav1.GetOptions{})
			if err != nil {
				return false, nil
			}
			for _, one := range a.Status.Conditions {
				return one.Type == apiregistrationv1.Available && one.Status == apiregistrationv1.ConditionTrue, nil
			}
			return false, nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TKE) importResource() error {
	restConfig, err := platformutil.GetRestConfig(&t.Cluster.Cluster, &t.Cluster.ClusterCredential)
	if err != nil {
		return err
	}
	client, err := tkeclientset.NewForConfig(restConfig)
	if err != nil {
		return err
	}

	// ensure api ready
	err = wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		_, err = client.PlatformV1().Clusters().List(metav1.ListOptions{})
		if err != nil {
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return err
	}

	_, err = client.PlatformV1().Clusters().Get(t.Cluster.Name, metav1.GetOptions{})
	if err == nil {
		err := client.PlatformV1().Clusters().Delete(t.Cluster.Name, &metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	_, err = client.PlatformV1().Clusters().Create(&t.Cluster.Cluster)
	if err != nil {
		return err
	}

	_, err = client.PlatformV1().ClusterCredentials().Get(t.Cluster.ClusterCredential.Name, metav1.GetOptions{})
	if err == nil {
		err := client.PlatformV1().ClusterCredentials().Delete(t.Cluster.ClusterCredential.Name, &metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	_, err = client.PlatformV1().ClusterCredentials().Create(&t.Cluster.ClusterCredential)
	if err != nil {
		return err
	}

	return nil
}

func (t *TKE) pushImages() error {
	imagesFilter := fmt.Sprintf("%s/*", t.Para.Config.Registry.Namespace())
	if t.Para.Config.Registry.Domain() != "docker.io" { // docker images filter ignore docker.io
		imagesFilter = t.Para.Config.Registry.Domain() + "/" + imagesFilter
	}
	cmd := exec.Command("sh", "-c",
		fmt.Sprintf("docker images --format='{{.Repository}}:{{.Tag}}' --filter='reference=%s'", imagesFilter),
	)
	out, err := cmd.Output()
	if err != nil {
		return pkgerrors.Wrap(err, "docker images error")
	}
	tkeImages := strings.Split(strings.TrimSpace(string(out)), "\n")
	for i, image := range tkeImages {
		nameAndTag := strings.Split(image, ":")
		if nameAndTag[1] == "<none>" {
			t.log.Printf("skip invalid tag:name=%s", image)
			continue
		}

		cmd = exec.Command("docker", "push", image)
		cmd.Stdout = t.log.Writer()
		cmd.Stderr = t.log.Writer()
		err = cmd.Run()
		if err != nil {
			return pkgerrors.Wrap(err, "docker push error")
		}

		t.log.Printf("upload %s to registry success[%d/%d]", image, i+1, len(tkeImages))
	}

	return nil
}

func (t *TKE) preInstallHook() error {
	return t.execHook(preInstallHook)
}

func (t *TKE) postClusterReadyHook() error {
	return t.execHook(postClusterReadyHook)
}

func (t *TKE) postInstallHook() error {
	return t.execHook(postInstallHook)
}

func (t *TKE) execHook(filename string) error {
	t.log.Printf("Execute hook script %s", filename)
	cmd := exec.Command(filename)
	cmd.Stdout = t.log.Writer()
	cmd.Stderr = t.log.Writer()
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (t *TKE) getKubeconfig() (*api.Config, error) {
	addr, err := platformutil.ClusterV1Address(&t.Cluster.Cluster)
	if err != nil {
		return nil, err
	}

	return kubeconfig.CreateWithToken(addr,
		t.Cluster.Name,
		"admin",
		t.Cluster.ClusterCredential.CACert,
		*t.Cluster.ClusterCredential.Token,
	), nil
}

func (t *TKE) writeKubeconfig() error {
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
