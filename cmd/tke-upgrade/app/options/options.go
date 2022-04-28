package options

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	platformv1 "tkestack.io/tke/api/platform/v1"
	v1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/cmd/tke-installer/app/config"
	"tkestack.io/tke/cmd/tke-installer/app/installer/constants"
	"tkestack.io/tke/cmd/tke-installer/app/installer/images"
	"tkestack.io/tke/cmd/tke-installer/app/installer/types"
	"tkestack.io/tke/pkg/util/containerregistry"
)

type Options map[string]interface{}

type TKE struct {
	Dir           string
	Version       string
	Servers       []string
	RedirectHosts []string

	Config  *config.Config           `json:"config"`
	Para    *types.CreateClusterPara `json:"para"`
	Cluster struct {
		v1.Cluster
		ClusterCredential v1.ClusterCredential
	} `json:"cluster"`
}

func New(dir string, version string) *TKE {
	return &TKE{
		Dir:     dir,
		Version: version,
	}
}

func (t *TKE) read(path string) (data []byte) {
	data, err := ioutil.ReadFile(t.Dir + "/" + path)
	if err != nil {
		panic(err)
	}
	return
}

func (t *TKE) loadTKE() {
	data := t.read(constants.ClusterFile)
	err := json.Unmarshal(data, &t)
	if err != nil {
		panic(err)
	}
}

func (t *TKE) Init() {
	t.loadTKE()

	containerregistry.Init(t.Para.Config.Registry.Domain(), t.Para.Config.Registry.Namespace())
	for _, address := range t.Cluster.Status.Addresses {
		if address.Type == platformv1.AddressReal {
			t.Servers = append(t.Servers, address.Host)
		}
	}
	t.RedirectHosts = t.Servers
	t.RedirectHosts = append(t.RedirectHosts, "tke-gateway")
	if t.Para.Config.Gateway != nil && t.Para.Config.Gateway.Domain != "" {
		t.RedirectHosts = append(t.RedirectHosts, t.Para.Config.Gateway.Domain)
	}
	if t.Para.Config.HA != nil {
		t.RedirectHosts = append(t.RedirectHosts, t.Para.Config.HA.VIP())
	}
	if t.Para.Cluster.Spec.PublicAlternativeNames != nil {
		t.RedirectHosts = append(t.RedirectHosts, t.Para.Cluster.Spec.PublicAlternativeNames...)
	}
}

func (t *TKE) GetFullName(name string) string {
	return containerregistry.GetImagePrefix(name) + "-amd64:" + t.Version
}

func (t *TKE) TKEGateway() (option Options) {
	option = map[string]interface{}{
		"Image":            t.GetFullName(images.Get().TKEGateway.Name),
		"OIDCClientSecret": string(t.read(constants.OIDCClientSecretFile)),
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
	return
}

func (t *TKE) ETCD() (option Options) {
	option = map[string]interface{}{
		"servers": t.Servers,
	}
	return
}

func (t *TKE) TKEAuthAPI() (option Options) {
	option = map[string]interface{}{
		"Replicas":         t.Config.Replicas,
		"Image":            t.GetFullName(images.Get().TKEAuthAPI.Name),
		"OIDCClientSecret": string(t.read(constants.OIDCClientSecretFile)),
		"AdminUsername":    t.Para.Config.Auth.TKEAuth.Username,
		"TenantID":         t.Para.Config.Auth.TKEAuth.TenantID,
		"RedirectHosts":    t.RedirectHosts,
	}
	return
}

func (t *TKE) TKEAuthController() (option Options) {
	option = map[string]interface{}{
		"Replicas":      t.Config.Replicas,
		"Image":         t.GetFullName(images.Get().TKEAuthController.Name),
		"AdminUsername": t.Para.Config.Auth.TKEAuth.Username,
		"AdminPassword": string(t.Para.Config.Auth.TKEAuth.Password),
	}
	return
}

func (t *TKE) TKEPlatformAPI() (option Options) {
	option = map[string]interface{}{
		"Replicas":   t.Config.Replicas,
		"Image":      t.GetFullName(images.Get().TKEPlatformAPI.Name),
		"EnableAuth": t.Para.Config.Auth.TKEAuth != nil,
	}
	if t.Para.Config.Auth.OIDCAuth != nil {
		option["OIDCClientID"] = t.Para.Config.Auth.OIDCAuth.ClientID
		option["OIDCIssuerURL"] = t.Para.Config.Auth.OIDCAuth.IssuerURL
		option["UseOIDCCA"] = t.Para.Config.Auth.OIDCAuth.CACert != nil
	}
	return
}

func (t *TKE) TKEPlatformController() (option Options) {
	option = map[string]interface{}{
		"Replicas":                t.Config.Replicas,
		"Image":                   t.GetFullName(images.Get().TKEPlatformController.Name),
		"ProviderResImage":        images.Get().ProviderRes.FullName(),
		"RegistryDomain":          t.Para.Config.Registry.Domain(),
		"RegistryNamespace":       t.Para.Config.Registry.Namespace(),
		"MonitorStorageType":      "",
		"MonitorStorageAddresses": "",
	}
	if t.Para.Config.Monitor != nil {
		if t.Para.Config.Monitor.InfluxDBMonitor != nil {
			option["MonitorStorageType"] = "influxdb"
			if t.Para.Config.Monitor.InfluxDBMonitor.LocalInfluxDBMonitor != nil {
				option["MonitorStorageAddresses"] = t.getLocalInfluxdbAddress()
			} else if t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor != nil {
				address := t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.URL
				if t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.Username != "" {
					address = address + "&u=" + t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.Username
				}
				if t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.Password != nil {
					address = address + "&p=" + string(t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.Password)
				}
				option["MonitorStorageAddresses"] = address
			}
		} else if t.Para.Config.Monitor.ESMonitor != nil {
			option["MonitorStorageType"] = "elasticsearch"
			address := t.Para.Config.Monitor.ESMonitor.URL
			if t.Para.Config.Monitor.ESMonitor.Username != "" {
				address = address + "&u=" + t.Para.Config.Monitor.ESMonitor.Username
			}
			if t.Para.Config.Monitor.ESMonitor.Password != nil {
				address = address + "&p=" + string(t.Para.Config.Monitor.ESMonitor.Password)
			}
			option["MonitorStorageAddresses"] = address
		}
	}
	return
}

func (t *TKE) TKEBusinessAPI() (option Options) {
	option = map[string]interface{}{
		"Replicas":                   t.Config.Replicas,
		"Image":                      t.GetFullName(images.Get().TKEBusinessAPI.Name),
		"TenantID":                   t.Para.Config.Auth.TKEAuth.TenantID,
		"Username":                   t.Para.Config.Auth.TKEAuth.Username,
		"SyncProjectsWithNamespaces": t.Config.SyncProjectsWithNamespaces,
		"EnableAuth":                 t.Para.Config.Auth.TKEAuth != nil,
		"EnableRegistry":             t.Para.Config.Registry.TKERegistry != nil,
	}
	if t.Para.Config.Auth.OIDCAuth != nil {
		option["OIDCClientID"] = t.Para.Config.Auth.OIDCAuth.ClientID
		option["OIDCIssuerURL"] = t.Para.Config.Auth.OIDCAuth.IssuerURL
		option["UseOIDCCA"] = t.Para.Config.Auth.OIDCAuth.CACert != nil
	}
	return
}

func (t *TKE) TKEBusinessController() (option Options) {
	option = map[string]interface{}{
		"Replicas":       t.Config.Replicas,
		"Image":          t.GetFullName(images.Get().TKEBusinessController.Name),
		"EnableAuth":     t.Para.Config.Auth.TKEAuth != nil,
		"EnableRegistry": t.Para.Config.Registry.TKERegistry != nil,
	}
	return
}

func (t *TKE) InfluxDB() (option Options) {
	option = map[string]interface{}{
		"Image":    images.Get().InfluxDB.FullName(),
		"NodeName": t.Servers[0],
	}
	return
}

func (t *TKE) TKEMonitorAPI() (option Options) {
	option = map[string]interface{}{
		"Replicas":   t.Config.Replicas,
		"Image":      t.GetFullName(images.Get().TKEMonitorAPI.Name),
		"EnableAuth": t.Para.Config.Auth.TKEAuth != nil,
	}
	if t.Para.Config.Auth.OIDCAuth != nil {
		option["OIDCClientID"] = t.Para.Config.Auth.OIDCAuth.ClientID
		option["OIDCIssuerURL"] = t.Para.Config.Auth.OIDCAuth.IssuerURL
		option["UseOIDCCA"] = t.Para.Config.Auth.OIDCAuth.CACert != nil
	}
	if t.Para.Config.Monitor != nil {
		if t.Para.Config.Monitor.ESMonitor != nil {
			option["StorageType"] = "es"
			option["StorageAddress"] = t.Para.Config.Monitor.ESMonitor.URL
			option["StorageUsername"] = t.Para.Config.Monitor.ESMonitor.Username
			option["StoragePassword"] = string(t.Para.Config.Monitor.ESMonitor.Password)
		} else if t.Para.Config.Monitor.InfluxDBMonitor != nil {
			option["StorageType"] = "influxDB"

			if t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor != nil {
				option["StorageAddress"] = t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.URL
				option["StorageUsername"] = t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.Username
				option["StoragePassword"] = string(t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.Password)
			} else if t.Para.Config.Monitor.InfluxDBMonitor.LocalInfluxDBMonitor != nil {
				// todo
				option["StorageAddress"] = t.getLocalInfluxdbAddress()
			}
		}
	}

	return
}

func (t *TKE) TKEMonitorController() (option Options) {
	option = map[string]interface{}{
		"Replicas":       t.Config.Replicas,
		"Image":          t.GetFullName(images.Get().TKEMonitorController.Name),
		"EnableBusiness": t.Para.Config.Business != nil,
	}
	if t.Para.Config.Monitor != nil {
		if t.Para.Config.Monitor.ESMonitor != nil {
			option["StorageType"] = "es"
			option["StorageAddress"] = t.Para.Config.Monitor.ESMonitor.URL
			option["StorageUsername"] = t.Para.Config.Monitor.ESMonitor.Username
			option["StoragePassword"] = string(t.Para.Config.Monitor.ESMonitor.Password)
		} else if t.Para.Config.Monitor.InfluxDBMonitor != nil {
			option["StorageType"] = "influxDB"

			if t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor != nil {
				option["StorageAddress"] = t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.URL
				option["StorageUsername"] = t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.Username
				option["StoragePassword"] = string(t.Para.Config.Monitor.InfluxDBMonitor.ExternalInfluxDBMonitor.Password)
			} else if t.Para.Config.Monitor.InfluxDBMonitor.LocalInfluxDBMonitor != nil {
				option["StorageAddress"] = t.getLocalInfluxdbAddress()
			}
		}
	}

	return
}

func (t *TKE) TKENotifyAPI() (option Options) {
	option = map[string]interface{}{
		"Replicas":   t.Config.Replicas,
		"Image":      t.GetFullName(images.Get().TKENotifyAPI.Name),
		"EnableAuth": t.Para.Config.Auth.TKEAuth != nil,
	}
	if t.Para.Config.Auth.OIDCAuth != nil {
		option["OIDCClientID"] = t.Para.Config.Auth.OIDCAuth.ClientID
		option["OIDCIssuerURL"] = t.Para.Config.Auth.OIDCAuth.IssuerURL
		option["UseOIDCCA"] = t.Para.Config.Auth.OIDCAuth.CACert != nil
	}
	return
}

func (t *TKE) TKENotifyController() (option Options) {
	option = map[string]interface{}{
		"Replicas": t.Config.Replicas,
		"Image":    t.GetFullName(images.Get().TKENotifyController.Name),
	}
	return
}

func (t *TKE) TKERegistryAPI() (option Options) {
	option = map[string]interface{}{
		"Image":         t.GetFullName(images.Get().TKERegistryAPI.Name),
		"NodeName":      t.Servers[0],
		"AdminUsername": t.Para.Config.Registry.TKERegistry.Username,
		"AdminPassword": string(t.Para.Config.Registry.TKERegistry.Password),
		"EnableAuth":    t.Para.Config.Auth.TKEAuth != nil,
		"DomainSuffix":  t.Para.Config.Registry.TKERegistry.Domain,
	}
	if t.Para.Config.Auth.OIDCAuth != nil {
		option["OIDCClientID"] = t.Para.Config.Auth.OIDCAuth.ClientID
		option["OIDCIssuerURL"] = t.Para.Config.Auth.OIDCAuth.IssuerURL
		option["UseOIDCCA"] = t.Para.Config.Auth.OIDCAuth.CACert != nil
	}
	return
}

func (t *TKE) getLocalInfluxdbAddress() string {
	var influxdbAddress string = fmt.Sprintf("http://%s:30086", t.Servers[0])
	if t.Para.Config.HA != nil && len(t.Para.Config.HA.VIP()) > 0 {
		vip := t.Para.Config.HA.VIP()
		influxdbAddress = fmt.Sprintf("http://%s:30086", vip) // influxdb svc must be set as NodePort type, and the nodePort is 30086
	}
	return influxdbAddress
}
