package installer

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/segmentio/ksuid"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"tkestack.io/tke/cmd/tke-installer/app/installer/certs"
	"tkestack.io/tke/cmd/tke-installer/app/installer/constants"
	images "tkestack.io/tke/pkg/platform/provider/baremetal/images"
	"tkestack.io/tke/pkg/spec"
	"tkestack.io/tke/pkg/util/pkiutil"
)

const (
	DataDir = "data/"

	AuthChartValuePath     = "auth-chart-values.yaml"
	PlatformChartValuePath = "platform-chart-values.yaml"
	GatewayChartValuePath  = "gateway-chart-values.yaml"
)

func GenerateValueChart() error {
	customConfig := CustomConfig{}
	data, err := ioutil.ReadFile("customConfig.yaml")
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, &customConfig); err != nil {
		return err
	}

	if isExist, _ := PathExists(DataDir); !isExist {
		os.Mkdir(DataDir, 0777)
	}

	if err := PatchPlatformVersion(); err != nil {
		return err
	}

	if err := customConfig.GenerateCertificates(); err != nil {
		return err
	}
	oIDCClientSecret := ksuid.New().String()
	if err := customConfig.GenerateAuthChartValuesYaml(oIDCClientSecret); err != nil {
		return err
	}
	if err := customConfig.GeneratePlatformChartValuesYaml(); err != nil {
		return err
	}
	if err := customConfig.GenerateGatewayChartValuesYaml(oIDCClientSecret); err != nil {
		return err
	}
	return nil
}

func PatchPlatformVersion() error {
	versionsByte, err := json.Marshal(spec.K8sVersions)
	if err != nil {
		return err
	}
	patchData := ClusterInfoPatch{}
	patchData.Data.TkeVersion = "466b0576c4b2b04979dfce9f3ac10177a8afbfc5"
	// get出来的值格式不对
	// patchData.Data.TkeVersion = version.Get().GitVersion
	patchData.Data.K8sValidVersions = string(versionsByte)

	bytes, err := yaml.Marshal(patchData)
	if err != nil {
		return err
	}
	patchFile := "patch.yaml"
	if err := ioutil.WriteFile(patchFile, bytes, 0644); err != nil {
		return err
	}

	commandStr := fmt.Sprintf("kubectl patch configmap cluster-info -n kube-public --patch \"$(cat %s)\"", patchFile)
	cmd := exec.Command("/bin/bash", "-c", commandStr)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (customConfig *CustomConfig) GenerateCertificates() error {
	serverIPs := customConfig.ServerIPs
	dnsNames := customConfig.DNSNames
	frontProxyCaCrtAbsPath := customConfig.FrontProxyCaCrtAbsPath
	etcdCrtAbsPath := customConfig.EtcdCrtAbsPath
	etcdKeyAbsPath := customConfig.EtcdKeyAbsPath

	ips := []net.IP{net.ParseIP("127.0.0.1")}
	if len(serverIPs) != 0 {
		for _, serverIPStr := range serverIPs {
			certIP := net.ParseIP(serverIPStr)
			if certIP == nil {
				return fmt.Errorf(fmt.Sprintf("SERVER IP: %s FORMAT error", serverIPStr))
			}
			ips = append(ips, certIP)
		}
	} else {
		return errors.New("SERVER_IPs should not be null")
	}
	if err := certs.Generate(dnsNames, ips, DataDir); err != nil {
		return err
	}

	frontProxyCaCrt, err := ioutil.ReadFile(frontProxyCaCrtAbsPath)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(constants.FrontProxyCACrtFile, frontProxyCaCrt, 0644); err != nil {
		return err
	}

	etcdCrt, err := ioutil.ReadFile(etcdCrtAbsPath)
	if err != nil {
		return err
	}
	etcdKey, err := ioutil.ReadFile(etcdKeyAbsPath)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(DataDir+"etcd-ca.crt", etcdCrt, 0644); err != nil {
		return err
	}
	etcdClientCertData, etcdClientKeyData, err := pkiutil.GenerateClientCertAndKey("tke", nil, etcdCrt, etcdKey)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(DataDir+"etcd.crt", etcdClientCertData, 0644); err != nil {
		return err
	}
	if err := ioutil.WriteFile(DataDir+"etcd.key", etcdClientKeyData, 0644); err != nil {
		return err
	}

	passwordCsv := fmt.Sprintf("%s,admin,1,administrator", ksuid.New().String())
	tokenCsv := fmt.Sprintf("%s,admin,1,administrator", ksuid.New().String())
	if err := ioutil.WriteFile(DataDir+"password.csv", []byte(passwordCsv), 0644); err != nil {
		return err
	}
	if err := ioutil.WriteFile(DataDir+"token.csv", []byte(tokenCsv), 0644); err != nil {
		return err
	}
	return nil
}

func (customConfig *CustomConfig) GenerateAuthChartValuesYaml(oIDCClientSecret string) error {
	cacrt, err := ioutil.ReadFile(DataDir + "ca.crt")
	if err != nil {
		return err
	}
	caKey, err := ioutil.ReadFile(DataDir + "ca.key")
	if err != nil {
		return err
	}
	adminCrt, err := ioutil.ReadFile(DataDir + "admin.crt")
	if err != nil {
		return err
	}
	adminKey, err := ioutil.ReadFile(DataDir + "admin.key")
	if err != nil {
		return err
	}
	serverCrt, err := ioutil.ReadFile(DataDir + "server.crt")
	if err != nil {
		return err
	}
	serverKey, err := ioutil.ReadFile(DataDir + "server.key")
	if err != nil {
		return err
	}
	webhookCrt, err := ioutil.ReadFile(DataDir + "webhook.crt")
	if err != nil {
		return err
	}
	webhookKey, err := ioutil.ReadFile(DataDir + "webhook.key")
	if err != nil {
		return err
	}
	etcdCaCrt, err := ioutil.ReadFile(DataDir + "etcd-ca.crt")
	if err != nil {
		return err
	}
	etcdCrt, err := ioutil.ReadFile(DataDir + "etcd.crt")
	if err != nil {
		return err
	}
	etcdKey, err := ioutil.ReadFile(DataDir + "etcd.key")
	if err != nil {
		return err
	}
	frontProxyCaCrt, err := ioutil.ReadFile(DataDir + "front-proxy-ca.crt")
	if err != nil {
		return err
	}
	tokenCsv, err := ioutil.ReadFile(DataDir + "token.csv")
	if err != nil {
		return err
	}
	passwordCsv, err := ioutil.ReadFile(DataDir + "password.csv")
	if err != nil {
		return err
	}
	obj := AuthChartValue{
		Etcd:            customConfig.Etcd,
		Cacrt:           string(cacrt),
		CaKey:           string(caKey),
		AdminCrt:        string(adminCrt),
		AdminKey:        string(adminKey),
		ServerCrt:       string(serverCrt),
		ServerKey:       string(serverKey),
		WebhookCrt:      string(webhookCrt),
		WebHookKey:      string(webhookKey),
		EtcdCaCrt:       string(etcdCaCrt),
		EtcdCrt:         string(etcdCrt),
		EtcdKey:         string(etcdKey),
		FrontProxyCaCrt: string(frontProxyCaCrt),
		PasswordCsv:     string(passwordCsv),
		TokenCsv:        string(tokenCsv),
	}

	originAuthCustomConfig := customConfig.TkeAuth
	if originAuthCustomConfig == nil {
		return errors.New("auth custom config content error")
	}

	if originAuthCustomConfig.API.Replicas <= 0 {
		return errors.New("auth custom config api replicas le 0")
	}
	if len(originAuthCustomConfig.API.Image) == 0 {
		return errors.New("auth custom config api image nil")
	}
	if originAuthCustomConfig.API.RedirectHosts == nil {
		originAuthCustomConfig.API.RedirectHosts = make([]string, len(customConfig.ServerIPs))
	}
	originAuthCustomConfig.API.RedirectHosts = append(originAuthCustomConfig.API.RedirectHosts, customConfig.ServerIPs...)
	originAuthCustomConfig.API.OIDCClientSecret = oIDCClientSecret

	if originAuthCustomConfig.Controller.Replicas <= 0 {
		return errors.New("auth custom config api replicas le 0")
	}
	if len(originAuthCustomConfig.Controller.Image) == 0 {
		return errors.New("auth custom config api image nil")
	}
	if len(originAuthCustomConfig.Controller.AdminUsername) == 0 && len(originAuthCustomConfig.Controller.AdminPassword) == 0 {
		originAuthCustomConfig.Controller.AdminUsername = "admin"
		originAuthCustomConfig.Controller.AdminPassword = base64.StdEncoding.EncodeToString([]byte(originAuthCustomConfig.Controller.AdminUsername))
	} else if len(originAuthCustomConfig.Controller.AdminUsername) == 0 {
		return errors.New("auth controller admin userName empty when admin password exists")
	} else if len(originAuthCustomConfig.Controller.AdminPassword) == 0 {
		return errors.New("auth controller admin password empty when admin userName exists")
	}

	obj.TKEAuth = *originAuthCustomConfig

	bytes, errMarshal := yaml.Marshal(obj)
	if errMarshal != nil {
		return errMarshal
	}
	if err := ioutil.WriteFile(AuthChartValuePath, bytes, 0644); err != nil {
		return err
	}
	return nil
}

func (customConfig *CustomConfig) GeneratePlatformChartValuesYaml() error {
	cacrt, err := ioutil.ReadFile(DataDir + "ca.crt")
	if err != nil {
		return err
	}

	obj := PlatformChartValue{
		Etcd:  customConfig.Etcd,
		Cacrt: string(cacrt),
	}

	originPlatformCustomConfig := customConfig.TkePlatform
	if originPlatformCustomConfig == nil {
		return errors.New("platform custom config content error")
	}

	if len(customConfig.TkePlatform.PublicIP) <= 0 {
		return errors.New("platform custom config publicIP nil")
	}
	if len(originPlatformCustomConfig.MetricsServerImage) <= 0 {
		originPlatformCustomConfig.MetricsServerImage = images.Get().MetricsServer.BaseName()
	}
	if len(originPlatformCustomConfig.AddonResizerImage) <= 0 {
		originPlatformCustomConfig.AddonResizerImage = images.Get().AddonResizer.BaseName()
	}

	if originPlatformCustomConfig.API.Replicas <= 0 {
		return errors.New("platform custom config api replicas le 0")
	}
	if len(originPlatformCustomConfig.API.Image) == 0 {
		return errors.New("platform custom config api image nil")
	}

	if originPlatformCustomConfig.Controller.Replicas <= 0 {
		return errors.New("platform custom config controller replicas le 0")
	}
	if len(originPlatformCustomConfig.Controller.Image) == 0 {
		return errors.New("platform custom config controller image nil")
	}
	if len(originPlatformCustomConfig.Controller.ProviderResImage) == 0 {
		return errors.New("platform custom config controller provider res image nil")
	}
	if len(originPlatformCustomConfig.Controller.MonitorStorageAddresses) == 0 {
		originPlatformCustomConfig.Controller.MonitorStorageAddresses = fmt.Sprintf("http://%s:8086", customConfig.ServerIPs[0])
	}
	obj.TKEPlatform = *originPlatformCustomConfig

	bytes, errMarshal := yaml.Marshal(obj)
	if errMarshal != nil {
		return errMarshal
	}
	if err := ioutil.WriteFile(PlatformChartValuePath, bytes, 0644); err != nil {
		return err
	}
	return nil
}

func (customConfig *CustomConfig) GenerateGatewayChartValuesYaml(oIDCClientSecret string) error {
	obj := GatewayChartValue{
		Etcd: customConfig.Etcd,
	}

	originGatewayCustomConfig := customConfig.TkeGateway
	if originGatewayCustomConfig == nil {
		return errors.New("gateway custom config content error")
	}

	if len(originGatewayCustomConfig.Image) == 0 {
		return errors.New("gateway custom config image nil")
	}
	if len(originGatewayCustomConfig.OIDCClientSecret) == 0 {
		originGatewayCustomConfig.OIDCClientSecret = oIDCClientSecret
	}
	obj.TKEGateway = *originGatewayCustomConfig

	bytes, errMarshal := yaml.Marshal(obj)
	if errMarshal != nil {
		return errMarshal
	}
	if err := ioutil.WriteFile(GatewayChartValuePath, bytes, 0644); err != nil {
		return err
	}
	return nil
}
