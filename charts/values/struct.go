package main

type ClusterInfoPatch struct {
	Data struct {
		K8sValidVersions string `yaml:"k8sValidVersions"`
		TkeVersion       string `yaml:"tkeVersion"`
	}
}

type EtcdConfig struct {
	Host string `yaml:"host"`
}

type AuthChartValue struct {
	Cacrt           string           `yaml:"caCrt"`
	CaKey           string           `yaml:"caKey"`
	AdminCrt        string           `yaml:"adminCrt"`
	AdminKey        string           `yaml:"adminKey"`
	ServerCrt       string           `yaml:"serverCrt"`
	ServerKey       string           `yaml:"serverKey"`
	WebhookCrt      string           `yaml:"webhookCrt"`
	WebHookKey      string           `yaml:"webhookKey"`
	EtcdCaCrt       string           `yaml:"etcdCaCrt"`
	EtcdCrt         string           `yaml:"etcdCrt"`
	EtcdKey         string           `yaml:"etcdKey"`
	FrontProxyCaCrt string           `yaml:"frontProxyCaCrt"`
	PasswordCsv     string           `yaml:"passwordCsv"`
	TokenCsv        string           `yaml:"tokenCsv"`
	Etcd            *EtcdConfig      `yaml:"etcd"`
	TKEAuth         AuthCustomConfig `yaml:",inline"`
}

type PlatformChartValue struct {
	Cacrt       string               `yaml:"caCrt"`
	Etcd        *EtcdConfig          `yaml:"etcd"`
	TKEPlatform PlatformCustomConfig `yaml:",inline"`
}

type GatewayChartValue struct {
	Etcd       *EtcdConfig         `yaml:"etcd"`
	TKEGateway GatewayCustomConfig `yaml:",inline"`
}

type CustomConfig struct {
	// TODO 目前只允许namespace为tke，后续开放支持配置namespace
	// ChartNamespace				string		`yaml:"chartNamespace"`
	ServerIPs              []string              `yaml:"serverIPs,flow"`
	DNSNames               []string              `yaml:"dnsNames,flow"`
	FrontProxyCaCrtAbsPath string                `yaml:"frontProxyCaCrtAbsPath"`
	EtcdCrtAbsPath         string                `yaml:"etcdCrtAbsPath"`
	EtcdKeyAbsPath         string                `yaml:"etcdKeyAbsPath"`
	Etcd                   *EtcdConfig           `yaml:"etcd"`
	TkeAuth                *AuthCustomConfig     `yaml:"tke-auth"`
	TkePlatform            *PlatformCustomConfig `yaml:"tke-platform"`
	TkeGateway             *GatewayCustomConfig  `yaml:"tke-gateway"`
}

type BaseCustomConfig struct {
	Replicas int    `yaml:"replicas"`
	Image    string `yaml:"image"`
}

type AuthCustomConfig struct {
	API struct {
		BaseCustomConfig `yaml:",inline"`
		RedirectHosts    []string `yaml:"redirectHosts,flow"`
		NodePort         string   `yaml:"nodePort,omitempty"`
		EnableAudit      string   `yaml:"enableAudit,omitempty"`
		TenantID         string   `yaml:"tenantID,omitempty"`
		OIDCClientSecret string   `yaml:"oIDCClientSecret,omitempty"`
		AdminUsername    string   `yaml:"adminUsername,omitempty"`
	}
	Controller struct {
		BaseCustomConfig `yaml:",inline"`
		AdminUsername    string `yaml:"adminUsername,omitempty"`
		AdminPassword    string `yaml:"adminPassword,omitempty"`
	}
}

type PlatformCustomConfig struct {
	PublicIP           string `yaml:"publicIP"`
	MetricsServerImage string `yaml:"metricsServerImage"`
	AddonResizerImage  string `yaml:"addonResizerImage"`
	API                struct {
		BaseCustomConfig `yaml:",inline"`
		EnableAuth       string `yaml:"enableAuth,omitempty"`
		EnableAudit      string `yaml:"enableAudit,omitempty"`
		OIDCClientID     string `yaml:"oIDCClientID,omitempty"`
		OIDCIssuerURL    string `yaml:"oIDCIssuerURL,omitempty"`
		UseOIDCCA        string `yaml:"useOIDCCA,omitempty"`
	}
	Controller struct {
		BaseCustomConfig        `yaml:",inline"`
		ProviderResImage        string `yaml:"providerResImage"`
		RegistryDomain          string `yaml:"registryDomain,omitempty"`
		RegistryNamespace       string `yaml:"registryNamespace,omitempty"`
		MonitorStorageType      string `yaml:"monitorStorageType,omitempty"`
		MonitorStorageAddresses string `yaml:"monitorStorageAddresses,omitempty"`
	}
}

type GatewayCustomConfig struct {
	Image                string `yaml:"image"`
	RegistryDomainSuffix string `yaml:"registryDomainSuffix,omitempty"`
	TenantID             string `yaml:"tenantID,omitempty"`
	OIDCClientSecret     string `yaml:"oIDCClientSecret,omitempty"`
	SelfSigned           string `yaml:"selfSigned,omitempty"`
	EnableAuth           string `yaml:"enableAuth,omitempty"`
	EnableRegistry       string `yaml:"enableRegistry,omitempty"`
	EnableBusiness       string `yaml:"enableBusiness,omitempty"`
	EnableMonitor        string `yaml:"enableMonitor,omitempty"`
	EnableLogagent       string `yaml:"enableLogagent,omitempty"`
	EnableAudit          string `yaml:"enableAudit,omitempty"`
	EnableApplication    string `yaml:"enableApplication,omitempty"`
	EnableMesh           string `yaml:"enableMesh,omitempty"`
	ServerKey            string `yaml:"serverKey,omitempty"`
	ServerCrt            string `yaml:"serverCrt,omitempty"`
}
