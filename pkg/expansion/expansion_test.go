package expansion

import (
	"gopkg.in/yaml.v2"
	"testing"
)

func Test_expansionDriver_readConfig(t *testing.T) {
	a := &ExpansionDriver{}
	err := yaml.Unmarshal([]byte(testContent), a)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("d.extra args %+v", a.CreateClusterExtraArgs)
}

var testContent = `operator: ""
charts: []
files: []
hooks:
  - pre-install
installer_skip_steps:
  - "Prepare push images to TKE registry"
  - "Push images to registry"
  - "Tag images"
  - "Load images"
  - "Push images"
  - "Setup local registry"
create_cluster_skip_conditions:
  - EnsureCopyFiles
  - EnsureDocker
  - EnsureRegistryHosts
  - EnsureKubelet
  - EnsureCNIPlugins
  - EnsureConntrackTools
  - EnsureKubeadm
create_cluster_delegate_conditions:
  - EnsurePreClusterInstallHook
create_cluster_extra_args:
  dockerExtraArgs:
    data-root: "/data/tcnp/data/docker"
  kubeletExtraArgs:
    feature-gates: "ServiceTopology=true,EndpointSlice=true"
    max-pods: "256"
  apiServerExtraArgs:
    enable-admission-plugins: "NodeRestriction"
    service-node-port-range: "80-32767"
    feature-gates: "ServiceTopology=true,EndpointSlice=true"
  controllerManagerExtraArgs:
    feature-gates: "ServiceTopology=true,EndpointSlice=true"
  schedulerExtraArgs:
    feature-gates: "ServiceTopology=true,EndpointSlice=true"
  etcd:
    local:
      dataDir: "/data/tcnp/data/etcd"
      extraArgs:
        data-dir: "/data/tcnp/data/etcd"`
