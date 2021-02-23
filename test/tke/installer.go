package tke

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog"
	"net/http"
	"time"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/cmd/tke-installer/app/installer/types"
	"tkestack.io/tke/pkg/util/ssh"
	"tkestack.io/tke/test/util/cloudprovider"
)

type Installer struct {
	node cloudprovider.Instance
}

func InitInstaller(provider cloudprovider.Provider) *Installer {
	nodes, err := provider.CreateInstances(1)
	if err != nil {
		panic(fmt.Errorf("create instance failed. %v", err))
	}
	return &Installer{node: nodes[0]}
}

func (installer *Installer) RunCMD(cmd string) (string, error) {
	klog.Info("Run CMD: ", cmd)
	s, err := ssh.New(&ssh.Config{
		User:     installer.node.Username,
		Password: installer.node.Password,
		Host:     installer.node.PublicIP,
		Port:     int(installer.node.Port),
	})
	if err != nil {
		return "", err
	}
	out, err := s.CombinedOutput(cmd)
	klog.Info("CMD output: ", string(out))
	return string(out), err
}

func (installer *Installer) InstallInstaller(os, arch, version string) error {
	klog.Info("Download and install installer")
	name := fmt.Sprintf("tke-installer-%s-%s-%s", os, arch, version)
	cmd := fmt.Sprintf("wget https://tke-release-1251707795.cos.ap-guangzhou.myqcloud.com/%s.run{,.sha256} && sha256sum --check --status %s.run.sha256 && chmod +x %s.run && ./%s.run",
		name, name, name, name)
	_, err := installer.RunCMD(cmd)
	return err
}

func (installer *Installer) CreateClusterParaTemplate(nodes []cloudprovider.Instance) *types.CreateClusterPara {
	para := new(types.CreateClusterPara)
	for _, one := range nodes {
		para.Cluster.Spec.Machines = append(para.Cluster.Spec.Machines, platformv1.ClusterMachine{
			IP:       one.InternalIP,
			Port:     one.Port,
			Username: one.Username,
			Password: []byte(one.Password),
		})
	}
	para.Config = types.Config{
		Basic: &types.Basic{
			Username: "admin",
			Password: []byte("admin"),
		},
		Auth: types.Auth{
			TKEAuth: &types.TKEAuth{},
		},
		Registry: types.Registry{
			TKERegistry: &types.TKERegistry{
				Domain: "registry.tke.com",
			},
		},
		Business: nil,
		Monitor:  nil,
		Logagent: &types.Logagent{},
		HA:       nil,
		Gateway: &types.Gateway{
			Domain: "console.tke.com",
			Cert: &types.Cert{
				SelfSignedCert: &types.SelfSignedCert{},
			},
		},
	}
	return para
}

func (installer *Installer) Install(createClusterPara *types.CreateClusterPara) error {
	klog.Info("Start installing")
	body, err := json.Marshal(createClusterPara)
	if err != nil {
		return fmt.Errorf("unmarshal CreateClusterPara failed. %v", err)
	}
	url := fmt.Sprintf("http://%s:8080/api/cluster", installer.node.InternalIP)
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("post data to install failed. %v", err)
	}
	resp.Body.Close()

	klog.Info("Wait install finish")
	err = wait.Poll(30*time.Second, 2*time.Hour, func() (bool, error) {
		progress, err := installer.GetInstallProgress()
		if err != nil {
			return false, err
		}
		klog.Info(progress.Status)
		switch progress.Status {
		case types.StatusUnknown, types.StatusDoing:
			return false, nil
		case types.StatusFailed:
			return false, fmt.Errorf("install failed:\n%s", progress.Data)
		case types.StatusSuccess:
			return true, nil
		default:
			klog.Infof("Install again as we got an unknown install progress status: %s", progress.Status)
			resp, err = http.Post(url, "application/json", bytes.NewReader(body))
			if err != nil {
				return false, fmt.Errorf("post data to install failed. %v", err)
			}
			return false, nil
		}
	})
	if err == nil {
		klog.Info("Install finished")
	}
	return err
}

func (installer *Installer) GetInstallProgress() (*types.ClusterProgress, error) {
	url := fmt.Sprintf("http://%s:8080/api/cluster/global/progress", installer.node.InternalIP)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get install progress failed. %v", err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read http response (install progress) body failed. %v", err)
	}
	progress := new(types.ClusterProgress)
	err = json.Unmarshal(data, progress)
	if err != nil {
		return nil, fmt.Errorf("unmarshal ClusterProgress failed. %v", err)
	}
	return progress, nil
}
