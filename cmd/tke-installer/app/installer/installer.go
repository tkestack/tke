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
	"os"
	"os/exec"
	"path"
	goruntime "runtime"
	"sort"
	"strconv"
	"strings"
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
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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
	"tkestack.io/tke/cmd/tke-installer/app/installer/types"
	baremetalcluster "tkestack.io/tke/pkg/platform/provider/baremetal/cluster"
	baremetalconfig "tkestack.io/tke/pkg/platform/provider/baremetal/config"
	baremetalconstants "tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
	clusterstrategy "tkestack.io/tke/pkg/platform/registry/cluster"
	v1 "tkestack.io/tke/pkg/platform/types/v1"
	"tkestack.io/tke/pkg/spec"
	"tkestack.io/tke/pkg/util/apiclient"
	"tkestack.io/tke/pkg/util/containerregistry"
	"tkestack.io/tke/pkg/util/docker"
	"tkestack.io/tke/pkg/util/hosts"
	"tkestack.io/tke/pkg/util/kubeconfig"
	"tkestack.io/tke/pkg/util/log"
	utilnet "tkestack.io/tke/pkg/util/net"
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

	log             *stdlog.Logger
	steps           []types.Handler
	progress        *types.ClusterProgress
	strategy        *clusterstrategy.Strategy
	clusterProvider clusterprovider.Provider
	isFromRestore   bool

	docker *docker.Docker

	globalClient kubernetes.Interface
	servers      []string
	namespace    string
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
	f, err := os.OpenFile(constants.ClusterLogFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0744)
	if err != nil {
		log.Fatal(err.Error())
	}
	c.log = stdlog.New(f, "", stdlog.LstdFlags)

	c.docker = new(docker.Docker)
	c.docker.Stdout = c.log.Writer()
	c.docker.Stderr = c.log.Writer()

	if !config.Force {
		data, err := ioutil.ReadFile(constants.ClusterFile)
		if err == nil {
			log.Infof("read %q success", constants.ClusterFile)
			err = json.Unmarshal(data, c)
			if err != nil {
				log.Warnf("load tke data error:%s", err)
			}
			log.Infof("load tke data success")
			c.isFromRestore = true
			c.progress.Status = types.StatusDoing
		}
	}

	return c
}

func (t *TKE) initSteps() {
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
	if !IsDevRegistry(t.Para.Config.Registry) {
		t.steps = append(t.steps, []types.Handler{
			{
				Name: "Load images",
				Func: t.loadImages,
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

	if !IsDevRegistry(t.Para.Config.Registry) {
		t.steps = append(t.steps, []types.Handler{
			{
				Name: "Push images",
				Func: t.pushImages,
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
			Name: "Execute post deploy hook",
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
	}...)

	if t.Para.Config.Auth.TKEAuth != nil {
		t.steps = append(t.steps, []types.Handler{
			{
				Name: "Install tke-auth-api",
				Func: t.installTKEAuthAPI,
			},
			{
				Name: "Install tke-auth-controller",
				Func: t.installTKEAuthController,
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

	t.steps = append(t.steps, []types.Handler{
		{
			Name: "Install tke-platform-api",
			Func: t.installTKEPlatformAPI,
		},
		{
			Name: "Install tke-platform-controller",
			Func: t.installTKEPlatformController,
		},
	}...)

	if t.Para.Config.Registry.TKERegistry != nil {
		t.steps = append(t.steps, []types.Handler{
			{
				Name: "Install tke-registry-api",
				Func: t.installTKERegistryAPI,
			},
		}...)
	}

	if t.Para.Config.Business != nil {
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
					Name: "Install InfluxDB",
					Func: t.installInfluxDB,
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

	// others

	// Add more tke component before THIS!!!
	if t.Para.Config.Gateway != nil {
		if t.IncludeSelf {
			t.steps = append(t.steps, []types.Handler{
				{
					Name: "Prepare images before stop local registry",
					Func: t.prepareImages,
				},
				{
					Name: "Stop local registry to give up 80/443 for tke-gateway",
					Func: t.stopLocalRegistry,
				},
			}...)
		}
		t.steps = append(t.steps, []types.Handler{
			{
				Name: "Install tke-gateway",
				Func: t.installTKEGateway,
			},
		}...)
	}

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

	t.steps = append(t.steps, []types.Handler{
		{
			Name: "Execute post deploy hook",
			Func: t.postInstallHook,
		},
	}...)

	t.steps = funk.Filter(t.steps, func(step types.Handler) bool {
		return !funk.ContainsString(t.Para.Config.SkipSteps, step.Name)
	}).([]types.Handler)

	t.log.Println("Steps:")
	for i, step := range t.steps {
		t.log.Printf("%d %s", i, step.Name)
	}
}

func (t *TKE) Run() {
	var err error
	if t.Config.NoUI {
		err = t.run(context.Background())
	} else {
		err = t.runWithUI(context.Background())
	}
	if err != nil {
		log.Error(err.Error())
	}
}

func (t *TKE) run(ctx context.Context) error {
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

	t.do(ctx)

	return nil
}

func (t *TKE) runWithUI(ctx context.Context) error {
	a := NewAssertsResource()
	restful.Add(a.WebService())
	restful.Add(t.WebService())
	s := NewSSHResource()
	restful.Add(s.WebService())

	restful.Filter(globalLogging)

	if t.isFromRestore {
		go t.do(ctx)
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

func (t *TKE) prepare() errors.APIStatus {
	t.setConfigDefault(&t.Para.Config)
	statusError := t.validateConfig(t.Para.Config)
	if statusError != nil {
		return statusError
	}

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

	statusError = t.validateResource(v1Cluster)
	if statusError != nil {
		return statusError
	}

	t.Cluster.Cluster = v1Cluster
	for _, one := range t.Cluster.Spec.Machines {
		ok, err := utilnet.InterfaceHasAddr(one.IP)
		if err != nil {
			return errors.NewInternalError(err)
		}
		if ok {
			t.IncludeSelf = true
			break
		}
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

	if config.HA != nil {
		if t.Para.Config.HA.TKEHA != nil {
			cluster.Spec.Features.HA = &platformv1.HA{
				TKEHA: &platformv1.TKEHA{VIP: t.Para.Config.HA.TKEHA.VIP},
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
}

func (t *TKE) validateConfig(config types.Config) *errors.StatusError {
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

	if config.Monitor != nil {
		if config.Monitor.InfluxDBMonitor != nil && config.Monitor.ESMonitor != nil {
			return errors.NewBadRequest("influxdb or es only had one")
		}
	}

	var dnsNames []string
	if config.Gateway != nil && config.Gateway.Domain != "" {
		dnsNames = append(dnsNames, config.Gateway.Domain)
	}
	if config.Registry.TKERegistry != nil {
		dnsNames = append(dnsNames, config.Registry.TKERegistry.Domain, "*."+config.Registry.TKERegistry.Domain)
	}
	if config.Gateway != nil && config.Gateway.Cert.ThirdPartyCert != nil {
		statusError := t.validateCertAndKey(config.Gateway.Cert.ThirdPartyCert.Certificate,
			config.Gateway.Cert.ThirdPartyCert.PrivateKey, dnsNames)
		if statusError != nil {
			return statusError
		}
	}

	return nil
}

func (t *TKE) validateCertAndKey(certificate []byte, privateKey []byte, dnsNames []string) *errors.StatusError {
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

// validateResource validate the cpu and memory of cluster machines whether meets the requirements.
func (t *TKE) validateResource(cluster *platformv1.Cluster) *errors.StatusError {
	var (
		cpuSum    int
		memorySum int
	)
	for _, machine := range cluster.Spec.Machines {
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
			return errors.NewInternalError(err)
		}
		cmd := "nproc --all"
		stdout, err := s.CombinedOutput(cmd)
		if err != nil {
			return errors.NewInternalError(fmt.Errorf("get cpu error: %w", err))
		}
		cpu, err := strconv.Atoi(strings.TrimSpace(string(stdout)))
		if err != nil {
			return errors.NewInternalError(fmt.Errorf("convert cpu value error: %w", err))
		}
		cpuSum += cpu

		cmd = "free -g | grep Mem  | awk '{print $2}'"
		stdout, err = s.CombinedOutput(cmd)
		if err != nil {
			return errors.NewInternalError(fmt.Errorf("get memory error: %w", err))
		}
		memory, err := strconv.Atoi(strings.TrimSpace(string(stdout)))
		if err != nil {
			return errors.NewInternalError(fmt.Errorf("convert memory value error: %w", err))
		}
		memorySum += memory
	}
	if cpuSum < constants.CPURequest {
		return errors.NewBadRequest(fmt.Sprintf("at lease %d cores are required", constants.CPURequest))
	}
	if memorySum < constants.MemoryRequest {
		return errors.NewBadRequest(fmt.Sprintf("at lease %d GiB memory are required", constants.MemoryRequest))
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
			return pkgerrors.Wrap(err, "get ip for registry error")
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

func (t *TKE) createCluster(req *restful.Request, rsp *restful.Response) {
	apiStatus := func() errors.APIStatus {
		if t.Step != 0 {
			return errors.NewAlreadyExists(platformv1.Resource("Cluster"), "global")
		}
		para := new(types.CreateClusterPara)
		err := req.ReadEntity(para)
		if err != nil {
			return errors.NewBadRequest(err.Error())
		}
		t.Para = para
		if err := t.prepare(); err != nil {
			return err
		}
		go t.do(req.Request.Context())

		return nil
	}()

	if apiStatus != nil {
		_ = rsp.WriteHeaderAndJson(int(apiStatus.Status().Code), apiStatus.Status(), restful.MIME_JSON)
	} else {
		_ = rsp.WriteHeaderAndEntity(http.StatusCreated, t.Para)
	}
}

func (t *TKE) retryCreateCluster(req *restful.Request, rsp *restful.Response) {
	go t.do(req.Request.Context())
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
		_ = response.WriteEntity(&types.CreateClusterPara{
			Cluster: *t.Cluster.Cluster,
			Config:  t.Para.Config,
		})
	}
}

func (t *TKE) findClusterProgress(request *restful.Request, response *restful.Response) {
	var err error
	var data []byte
	apiStatus := func() errors.APIStatus {
		clusterName := request.PathParameter("name")
		if t.Cluster.Cluster == nil {
			return errors.NewBadRequest("no cluater available")
		}
		if t.Cluster.Name != clusterName {
			return errors.NewNotFound(platform.Resource("Cluster"), clusterName)
		}
		data, err = ioutil.ReadFile(constants.ClusterLogFile)
		if err != nil {
			return errors.NewInternalError(err)
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

func (t *TKE) do(ctx context.Context) {
	start := time.Now()

	containerregistry.Init(t.Para.Config.Registry.Domain(), t.Para.Config.Registry.Namespace())
	t.initSteps()

	if t.Step == 0 {
		t.log.Print("===>starting install task")
		t.progress.Status = types.StatusDoing
	}

	if t.runAfterClusterReady() {
		t.initDataForDeployTKE()
	}

	for t.Step < len(t.steps) {
		t.log.Printf("%d.%s doing", t.Step, t.steps[t.Step].Name)
		start := time.Now()
		err := t.steps[t.Step].Func(ctx)
		if err != nil {
			t.progress.Status = types.StatusFailed
			t.log.Printf("%d.%s [Failed] [%fs] error %s", t.Step, t.steps[t.Step].Name, time.Since(start).Seconds(), err)
			return
		}
		t.log.Printf("%d.%s [Success] [%fs]", t.Step, t.steps[t.Step].Name, time.Since(start).Seconds())

		t.Step++
		t.backup()
	}

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

	t.progress.Servers = t.servers
	if t.Para.Config.HA != nil {
		t.progress.Servers = append(t.progress.Servers, t.Para.Config.HA.VIP())
	}
	t.log.Printf("===>install task [Sucesss] [%fs]", time.Since(start).Seconds())
}

func (t *TKE) runAfterClusterReady() bool {
	return t.Cluster.Status.Phase == platformv1.ClusterRunning
}

func (t *TKE) generateCertificates(ctx context.Context) error {
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

	var errCount int
	for {
		start := time.Now()
		err := t.clusterProvider.OnCreate(ctx, t.Cluster)
		if err != nil {
			return err
		}
		t.backup()
		condition := t.Cluster.Status.Conditions[len(t.Cluster.Status.Conditions)-1]
		switch condition.Status {
		case platformv1.ConditionFalse: // means current condition run into error
			t.log.Printf("OnCreate.%s [Failed] [%fs] reason: %s message: %s retry: %d",
				condition.Type, time.Since(start).Seconds(),
				condition.Reason, condition.Message, errCount)
			if errCount >= 10 {
				return pkgerrors.New("the retry limit has been reached")
			}

			errCount++
		case platformv1.ConditionUnknown: // means has next condition need to process
			condition = t.Cluster.Status.Conditions[len(t.Cluster.Status.Conditions)-2]
			t.log.Printf("OnCreate.%s [Success] [%fs]", condition.Type, time.Since(start).Seconds())

			t.Cluster.Status.Reason = ""
			t.Cluster.Status.Message = ""
			errCount = 0
		case platformv1.ConditionTrue: // means all condition is done
			if t.Cluster.Status.Phase != platformv1.ClusterRunning {
				return pkgerrors.Errorf("OnCreate.%s, no next condition but cluster is not running!", condition.Type)
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
	return ioutil.WriteFile(constants.ClusterFile, data, 0777)
}

func (t *TKE) loadImages(ctx context.Context) error {
	if _, err := os.Stat(constants.ImagesFile); err != nil {
		return err
	}
	err := t.docker.LoadImages(constants.ImagesFile)
	if err != nil {
		return err
	}

	tkeImages, err := t.docker.GetImages(constants.ImagesPattern)
	if err != nil {
		return err
	}

	for _, image := range tkeImages {
		imageNames := strings.Split(image, "/")
		if len(imageNames) != 2 {
			t.log.Printf("invalid image name:name=%s", image)
			continue
		}
		name, _, _, err := t.docker.GetNameArchTag(imageNames[1])
		if err != nil {
			t.log.Printf("skip invalid image: %s", image)
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
	localHosts.File = "/etc/hosts"
	err = localHosts.Set("127.0.0.1")
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile("hosts")
	if err != nil {
		return err
	}
	t.log.Print(string(data))

	return nil
}

func (t *TKE) startLocalRegistry() error {
	err := t.stopLocalRegistry(context.Background())
	if err != nil {
		return err
	}

	err = t.docker.ClearLocalManifests()
	if err != nil {
		return err
	}

	registryImage := strings.ReplaceAll(images.Get().Registry.FullName(), ":", fmt.Sprintf("-%s:", goruntime.GOARCH))

	err = t.docker.RunImage(registryImage, constants.RegistryHTTPOptions, "")
	if err != nil {
		return err
	}

	// for docker manifest create which --insecure is not working
	err = t.docker.RunImage(registryImage, constants.RegistryHTTPSOptions, "")
	if err != nil {
		return err
	}

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

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "certs",
			Namespace: t.namespace,
		},
		Data: map[string]string{
			"etcd-ca.crt":        string(t.Cluster.ClusterCredential.ETCDCACert),
			"etcd.crt":           string(t.Cluster.ClusterCredential.ETCDAPIClientCert),
			"etcd.key":           string(t.Cluster.ClusterCredential.ETCDAPIClientKey),
			"ca.crt":             string(caCrt),
			"ca.key":             string(caKey),
			"front-proxy-ca.crt": string(frontProxyCACrt),
			"server.crt":         string(serverCrt),
			"server.key":         string(serverKey),
			"admin.crt":          string(adminCrt),
			"admin.key":          string(adminKey),
		},
	}

	if t.Para.Config.Auth.OIDCAuth != nil {
		cm.Data["oidc-ca.crt"] = string(t.Para.Config.Auth.OIDCAuth.CACert)
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

func (t *TKE) prepareBaremetalProviderConfig(ctx context.Context) error {
	c, err := baremetalconfig.New(constants.ProviderConfigFile)
	if err != nil {
		return err
	}
	if t.Para.Config.Registry.ThirdPartyRegistry == nil &&
		t.Para.Config.Registry.TKERegistry != nil {
		ip := t.Cluster.Spec.Machines[0].IP // registry current only run in first node
		c.Registry.IP = ip
	}
	if t.auditEnabled() {
		c.Audit.Address = t.determineGatewayHTTPSAddress()
	}

	err = c.Save(constants.ProviderConfigFile)
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
			Name: "csi-operator-manifests",
			File: baremetalconstants.ManifestsDir + "/csi-operator/*",
		},
		{
			Name: "keepalived-manifests",
			File: baremetalconstants.ManifestsDir + "/keepalived/*",
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
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "tke-gateway",
			Namespace: t.namespace,
		},
		Spec: corev1.PodSpec{
			NodeSelector: map[string]string{
				"node-role.kubernetes.io/master": "",
			},
			Containers: []corev1.Container{
				{
					Name:  "tke-gateway",
					Image: images.Get().TKEGateway.FullName(),
				},
			},
		},
	}
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	err := apiclient.PullImageWithPod(ctx, t.globalClient, pod)
	if err != nil {
		return fmt.Errorf("prepare image error: %w", err)
	}

	return nil
}

func (t *TKE) stopLocalRegistry(ctx context.Context) error {
	err := t.docker.RemoveContainers("registry-http", "registry-https")
	if err != nil {
		return err
	}
	return nil
}

func (t *TKE) installTKEGateway(ctx context.Context) error {
	option := map[string]interface{}{
		"Image":            images.Get().TKEGateway.FullName(),
		"OIDCClientSecret": t.readOrGenerateString(constants.OIDCClientSecretFile),
		"SelfSigned":       t.Para.Config.Gateway.Cert.SelfSignedCert != nil,
		"EnableRegistry":   t.Para.Config.Registry.TKERegistry != nil,
		"EnableAuth":       t.Para.Config.Auth.TKEAuth != nil,
		"EnableMonitor":    t.Para.Config.Monitor != nil,
		"EnableBusiness":   t.Para.Config.Business != nil,
		"EnableLogagent":   t.Para.Config.Logagent != nil,
		"EnableAudit":      t.auditEnabled(),
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
	err := apiclient.CreateResourceWithDir(ctx, t.globalClient, "manifests/tke-gateway/*.yaml", option)
	if err != nil {
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

func (t *TKE) installTKELogagentAPI(ctx context.Context) error {
	options := map[string]interface{}{
		"Replicas":       t.Config.Replicas,
		"Image":          images.Get().TKELogagentAPI.FullName(),
		"TenantID":       t.Para.Config.Auth.TKEAuth.TenantID,
		"Username":       t.Para.Config.Auth.TKEAuth.Username,
		"EnableAuth":     t.Para.Config.Auth.TKEAuth != nil,
		"EnableRegistry": t.Para.Config.Registry.TKERegistry != nil,
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
	return apiclient.CreateResourceWithDir(ctx, t.globalClient, "manifests/etcd/*.yaml",
		map[string]interface{}{
			"Servers": t.servers,
		})
}

func (t *TKE) installTKEAuthAPI(ctx context.Context) error {
	redirectHosts := t.servers
	redirectHosts = append(redirectHosts, "tke-gateway")
	if t.Para.Config.Gateway != nil && t.Para.Config.Gateway.Domain != "" {
		redirectHosts = append(redirectHosts, t.Para.Config.Gateway.Domain)
	}
	if t.Para.Config.HA != nil {
		redirectHosts = append(redirectHosts, t.Para.Config.HA.VIP())
	}
	if t.Para.Cluster.Spec.PublicAlternativeNames != nil {
		redirectHosts = append(redirectHosts, t.Para.Cluster.Spec.PublicAlternativeNames...)
	}

	option := map[string]interface{}{
		"Replicas":         t.Config.Replicas,
		"Image":            images.Get().TKEAuthAPI.FullName(),
		"OIDCClientSecret": t.readOrGenerateString(constants.OIDCClientSecretFile),
		"AdminUsername":    t.Para.Config.Auth.TKEAuth.Username,
		"TenantID":         t.Para.Config.Auth.TKEAuth.TenantID,
		"RedirectHosts":    redirectHosts,
	}
	err := apiclient.CreateResourceWithDir(ctx, t.globalClient, "manifests/tke-auth-api/*.yaml", option)
	if err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(ctx, t.globalClient, t.namespace, "tke-auth-api")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) installTKEAuthController(ctx context.Context) error {
	err := apiclient.CreateResourceWithDir(ctx, t.globalClient, "manifests/tke-auth-controller/*.yaml",
		map[string]interface{}{
			"Replicas":      t.Config.Replicas,
			"Image":         images.Get().TKEAuthController.FullName(),
			"AdminUsername": t.Para.Config.Auth.TKEAuth.Username,
			"AdminPassword": string(t.Para.Config.Auth.TKEAuth.Password),
		})
	if err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(ctx, t.globalClient, t.namespace, "tke-auth-controller")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
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
	}

	if err := apiclient.CreateResourceWithDir(ctx, t.globalClient, "manifests/tke-audit-api/*.yaml", options); err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(ctx, t.globalClient, t.namespace, "tke-audit-api")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) installTKEPlatformAPI(ctx context.Context) error {
	options := map[string]interface{}{
		"Replicas":    t.Config.Replicas,
		"Image":       images.Get().TKEPlatformAPI.FullName(),
		"EnableAuth":  t.Para.Config.Auth.TKEAuth != nil,
		"EnableAudit": t.auditEnabled(),
	}
	if t.Para.Config.Auth.OIDCAuth != nil {
		options["OIDCClientID"] = t.Para.Config.Auth.OIDCAuth.ClientID
		options["OIDCIssuerURL"] = t.Para.Config.Auth.OIDCAuth.IssuerURL
		options["UseOIDCCA"] = t.Para.Config.Auth.OIDCAuth.CACert != nil
	}
	err := apiclient.CreateResourceWithDir(ctx, t.globalClient, "manifests/tke-platform-api/*.yaml", options)
	if err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(ctx, t.globalClient, t.namespace, "tke-platform-api")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) installTKEPlatformController(ctx context.Context) error {
	params := map[string]interface{}{
		"Replicas":                t.Config.Replicas,
		"Image":                   images.Get().TKEPlatformController.FullName(),
		"ProviderResImage":        images.Get().ProviderRes.FullName(),
		"RegistryDomain":          t.Para.Config.Registry.Domain(),
		"RegistryNamespace":       t.Para.Config.Registry.Namespace(),
		"MonitorStorageType":      "",
		"MonitorStorageAddresses": "",
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

	if err := apiclient.CreateResourceWithDir(ctx, t.globalClient, "manifests/tke-platform-controller/*.yaml", params); err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(ctx, t.globalClient, t.namespace, "tke-platform-controller")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
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

func (t *TKE) installInfluxDB(ctx context.Context) error {
	err := apiclient.CreateResourceWithDir(ctx, t.globalClient, "manifests/influxdb/*.yaml",
		map[string]interface{}{
			"Image":    images.Get().InfluxDB.FullName(),
			"NodeName": t.servers[0],
		})
	if err != nil {
		return err
	}
	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckStatefulSet(ctx, t.globalClient, t.namespace, "influxdb")
		if err != nil || !ok {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) installTKEMonitorAPI(ctx context.Context) error {
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
		"Replicas":       t.Config.Replicas,
		"Image":          images.Get().TKEMonitorController.FullName(),
		"EnableBusiness": t.Para.Config.Business != nil,
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
		"Replicas":   t.Config.Replicas,
		"Image":      images.Get().TKENotifyAPI.FullName(),
		"EnableAuth": t.Para.Config.Auth.TKEAuth != nil,
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

func (t *TKE) installTKERegistryAPI(ctx context.Context) error {
	options := map[string]interface{}{
		"Image":         images.Get().TKERegistryAPI.FullName(),
		"NodeName":      t.servers[0],
		"AdminUsername": t.Para.Config.Registry.TKERegistry.Username,
		"AdminPassword": string(t.Para.Config.Registry.TKERegistry.Password),
		"EnableAuth":    t.Para.Config.Auth.TKEAuth != nil,
		"DomainSuffix":  t.Para.Config.Registry.TKERegistry.Domain,
	}
	if t.Para.Config.Auth.OIDCAuth != nil {
		options["OIDCClientID"] = t.Para.Config.Auth.OIDCAuth.ClientID
		options["OIDCIssuerURL"] = t.Para.Config.Auth.OIDCAuth.IssuerURL
		options["UseOIDCCA"] = t.Para.Config.Auth.OIDCAuth.CACert != nil
	}
	err := apiclient.CreateResourceWithDir(ctx, t.globalClient, "manifests/tke-registry-api/*.yaml", options)
	if err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		ok, err := apiclient.CheckDeployment(ctx, t.globalClient, t.namespace, "tke-registry-api")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
}

func (t *TKE) preparePushImagesToTKERegistry(ctx context.Context) error {
	localHosts := hosts.LocalHosts{Host: t.Para.Config.Registry.Domain(), File: "hosts"}
	err := localHosts.Set(t.servers[0])
	if err != nil {
		return err
	}
	localHosts.File = "/etc/hosts"
	err = localHosts.Set(t.servers[0])
	if err != nil {
		return err
	}

	dir := path.Join(constants.DockerCertsDir, t.Para.Config.Registry.Domain())
	_ = os.MkdirAll(dir, 0777)
	caCert, _ := ioutil.ReadFile(constants.CACrtFile)
	err = ioutil.WriteFile(path.Join(dir, "ca.crt"), caCert, 0644)
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

func (t *TKE) registerAPI(ctx context.Context) error {
	caCert, _ := ioutil.ReadFile(constants.CACrtFile)

	restConfig, err := t.Cluster.RESTConfigForBootstrap(&rest.Config{})
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

		err = wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
			_, err := client.ApiregistrationV1().APIServices().Get(ctx, apiService.Name, metav1.GetOptions{})
			if err == nil {
				err := client.ApiregistrationV1().APIServices().Delete(ctx, apiService.Name, metav1.DeleteOptions{})
				if err != nil {
					return false, nil
				}
			}
			if _, err := client.ApiregistrationV1().APIServices().Create(ctx, apiService, metav1.CreateOptions{}); err != nil {
				if !errors.IsAlreadyExists(err) {
					return false, nil
				}
			}
			return true, nil
		})
		if err != nil {
			return pkgerrors.Wrapf(err, "register apiservice %v error", one)
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
			return pkgerrors.Wrapf(err, "check apiservices %v error", one)
		}
	}
	return nil
}

func (t *TKE) importResource(ctx context.Context) error {
	restConfig, err := t.Cluster.RESTConfigForBootstrap(&rest.Config{Timeout: 120 * time.Second})
	if err != nil {
		return err
	}

	client, err := tkeclientset.NewForConfig(restConfig)
	if err != nil {
		return err
	}

	// ensure api ready
	err = wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		_, err = client.PlatformV1().Clusters().List(ctx, metav1.ListOptions{})
		if err != nil {
			return false, nil
		}
		_, err = client.PlatformV1().ClusterCredentials().List(ctx, metav1.ListOptions{})
		if err != nil {
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return err
	}

	_, err = client.PlatformV1().ClusterCredentials().Get(ctx, t.Cluster.ClusterCredential.Name, metav1.GetOptions{})
	if err == nil {
		err := client.PlatformV1().ClusterCredentials().Delete(ctx, t.Cluster.ClusterCredential.Name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	_, err = client.PlatformV1().ClusterCredentials().Create(ctx, t.Cluster.ClusterCredential, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	_, err = client.PlatformV1().Clusters().Get(ctx, t.Cluster.Name, metav1.GetOptions{})
	if err == nil {
		err := client.PlatformV1().Clusters().Delete(ctx, t.Cluster.Name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	_, err = client.PlatformV1().Clusters().Create(ctx, t.Cluster.Cluster, metav1.CreateOptions{})
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
	sort.Strings(tkeImages)
	tkeImagesSet := sets.NewString(tkeImages...)
	manifestSet := sets.NewString()

	// clear all local manifest lists before create any manifest list
	err = t.docker.ClearLocalManifests()
	if err != nil {
		return err
	}

	for i, image := range tkeImages {
		name, arch, tag, err := t.docker.GetNameArchTag(image)
		if err != nil { // skip invalid image
			t.log.Printf("skip invalid image: %s", image)
			continue
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
			manifestName := fmt.Sprintf("%s:%s", name, tag)
			manifestSet.Insert(manifestName) // To speed up, push manifests after all changes have made

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

		t.log.Printf("upload %s to registry success[%d/%d]", image, i+1, len(tkeImages))
	}

	sortedManifests := manifestSet.List()
	for i, manifest := range sortedManifests {
		err = t.docker.PushManifest(manifest, true)
		if err != nil {
			return nil
		}
		t.log.Printf("push manifest %s to registry success[%d/%d]", manifest, i+1, len(sortedManifests))
	}

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
			err := remoteHosts.Set(t.Cluster.Spec.Machines[0].IP)
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
