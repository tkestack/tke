# Tkestack addon framework recap


**Author**: huxiaoliang([@huxiaoliang](https://github.com/huxiaoliang))

**Status** (20201220): Suspend

## Summary

In order to extend Tkestack functionality so that support more value-add but don't impact core part too much, Tkestack introduced `addon` framework to address this requirement. There are 2 types of`addon` implementations currently:

1. Manifests based: there are several sub directory under `manifests` to identity individual `addon` instance,  `tke-platform-api` will use `go-templete` render the yaml files with customized parameters, then apply them to business/target cluster to install:

```
root@dev:~/pkg/platform/provider/baremetal# tree manifests/  
manifests/
├── csi-operator
│   └── csi-operator.yaml
├── gpu
│   └── nvidia-device-plugin.yaml
├── gpu-manager
│   └── gpu-manager.yaml
├── keepalived
│   ├── keepalived.conf
│   └── keepalived.yaml
└── metrics-server
    └── metrics-server.yaml
```

2. Controller based: There are several individual `controller` inside `tke-platform-controller`in global cluster to watch specified `addon` CR operation,  then leverage `tke-platform-api` access to business/target cluster manage `addon` instance lief cycle.

```
root@dev:~/pkg# tree platform/controller/addon/ -L 1
platform/controller/addon/
├── cronhpa
├── helm
├── ipam
├── lbcf
├── logcollector
├── persistentevent
├── prometheus
├── storage
└── tappcontroller

root@VM-0-77-ubuntu:~# curl -sk -H "Authorization: Bearer $(cat /etc/kubernetes/known_tokens.csv |cut -d "," -f 1)" -H "Content-Type:application/json" https://127.0.0.1:6443/apis/platform.tkestack.io/v1/clusteraddontypes | jq -r .items[].metadata.name
lbcf
helm
persistentevent
logcollector
csioperator
prometheus
ipam
tappcontroller
volumedecorator
cronhpa
```

After `tke-application` enabled, Tkestack has the ability to use `helm chart` as k8s native approach manage application directly, so `addon` framework should recap according to this new `out-tree` approach, the benefits as bellows:

 - Loose coupled relations between `tke-platform` and individual `addon`: if new addon onboard or old addon update/upgrade, no need rebuild or update core part, individual `addon`chart build and release will out of Tkestack,  no hard dependency for each other
 
 - Unify `addon` instance life cycle management: `tke-application` will take the responsibility for managing all `addon` charts include:
   - Create addon xxx               -->  `helm install xxx`
   - Upgrade addon xxx           --> `helm upgrade xxx`
   - Delete addon xxx              -->  `helm delete xxx`
   - Healthy check addon xxx  -->  `helm get xxx`

 - Decouple the installation and upgrade process of Tkestack: k8s , Tkestack built-in components and Tkestack addons

 - The `helm hook` mechanism allow chart developers to intervene at certain points in a chart release's life cycle to support more scenario

 - Better development experience and easily integration for internal developer and community contributor to extend Tkestack

## Scope

 **In-Scope**: 
 1. (**P1**) Porting 2 types of addons to helm charts
 2. (**P1**) Enable CI to build all addon charts
 3. (**P1**) UI enhance to use new API manage addon
 4. (**P2**) Nice to have: Support hook mechanism to pick up `user-defined` charts and push them to chart repo during Tkestack installation for day 2 install
 3. (**P2**) Define apps in `cluster` object, and create apps during creating business cluster(s)

**Out-Of-Scope**: 

 1. Tkestack built-in component helm chart support
 2. Define apps in `global cluster` object (temporarily)
 3. Transform `tke coms` to `build-in charts` and `build-in apps` (temporarily)

## Limitation

1. Enable `tke-application` when creating `global cluster`

## Main proposal

1. Enable `tke-application` installed as default during Tkestack installation (done)

2. `helm push` sdk will used to push all charts tgz package from `bootstrap` container to chart repo during tkestack (done):
- `https://github.com/tkestack/tke/pull/1182`

3. Label the chart so that distinguish system built-in addon chart and other charts,  `chart list` API will retrieval the chart instead of  `clusteraddontypes`:
- `https://github.com/tkestack/tke/issues/1357`

4. Tkestack `tke-application` controller will handle cross tenant request validation, below pr should get revert (done):
- `https://github.com/tkestack/tke/pull/978`
- `https://github.com/tkestack/tke/pull/1007`

5. `tke-installer` will push build-in/expansion charts to registry (done):
- `https://github.com/tkestack/tke/pull/1284`
- `https://github.com/tkestack/tke/pull/1375`

6. `tke-installer` will install build-in/expansion applications (done):
- `https://github.com/tkestack/tke/pull/1350`

7. Add label `build-in` for `build-in apps` if they are installed by `tke-installer` and their charts are `build-in charts` 
which are default charts in `tke-installer` release package:
- `https://github.com/tkestack/tke/issues/1359`

8. Support upgrade apps with `build-in` label during `tke-installer` upgrading Tkestack:
- `https://github.com/tkestack/tke/issues/1358`

9. Tkestack `platform` will create applications defined in `cluster` object during creating business cluster:
- `https://github.com/tkestack/tke/pull/1372`

10. Transform `addons` to `build-in charts` and install them through `tke-application`

11. Enhance `tke-installer` and `tkestack-gateway` UI

![enter image description here](../../docs/images/addon-charts.png)

## Future work

1. Define apps in `global cluster` object, and install apps during creating `global cluster`
2. Transform `tke coms` to helm charts
3. Define `tke coms` as `build-in apps` in `global cluster` object

## User case

#### Case 1. Installer install build-in apps during creating global cluster

Before UI support tke-installer to set apps in `global cluster` object, hardcode some `build-in apps` in `tke-installer` and use `tke.json` with empty `PlatformApps`:

```json
{
	"config": {
		"ServerName": "tke-installer",
		"ListenAddr": ":8080",
		"NoUI": false,
		"Config": "conf/tke.json",
		"Force": false,
		"SyncProjectsWithNamespaces": false,
		"Replicas": 2,
		"Upgrade": false,
		"PrepareCustomK8sImages": false,
		"PrepareCustomCharts": false,
		"Kubeconfig": "conf/kubeconfig",
		"RegistryUsername": "",
		"RegistryPassword": "",
		"RegistryDomain": "",
		"RegistryNamespace": "",
		"CustomUpgradeResourceDir": "data/custom_upgrade_resource",
		"CustomChartsName": "custom.charts.tar.gz",
		"EnableCustomExpansion": true,
		"CustomExpansionDir": "data/expansions/",
		// empty
		"PlatformApps": []
	},
........
}
```

`build-in apps` will fullfill `PlatformApps` during installing.

Tkestack will manage `build-in` apps life-cycle through `build-in` labels. It means that `build-in apps` will be upgraded if Tkestack platform is upgraded.

#### Case 2. Installer install expansion apps during creating global cluster

Use `tke.json` with expansions apps:

```json
{
	"config": {
		"ServerName": "tke-installer",
		"ListenAddr": ":8080",
		"NoUI": false,
		"Config": "conf/tke.json",
		"Force": false,
		"SyncProjectsWithNamespaces": false,
		"Replicas": 2,
		"Upgrade": false,
		"PrepareCustomK8sImages": false,
		"PrepareCustomCharts": false,
		"Kubeconfig": "conf/kubeconfig",
		"RegistryUsername": "",
		"RegistryPassword": "",
		"RegistryDomain": "",
		"RegistryNamespace": "",
		"CustomUpgradeResourceDir": "data/custom_upgrade_resource",
		"CustomChartsName": "custom.charts.tar.gz",
		"EnableCustomExpansion": true,
		"CustomExpansionDir": "data/expansions/",
		// expansion apps
		"PlatformApps": [
			{
				// app name
				"Name": "demo",
				// if enabled, for UI
				"Enable": true,
				"Chart": {
					// chart name
					"Name": "demo",
					"TenantID": "default",
					// repo
					"ChartGroupName": "public",
					"Version": "1.0.0",
					"TargetCluster": "global",
					"TargetNamespace": "default",
					// helm chart values
					"Values": {
						"key2": "val2-override"
					}
				}
			}
		]
	},
......
}
```

#### Case 3. tke-installer upgrade build-in apps

Download next minor version of current version `tke-installer` and upgrade through `tke-installerxxx --upgrade`.

#### Case 4. Define apps in cluster object and install apps during createing cluster

Define apps in cluster object:

```yaml
---
apiVersion: platform.tkestack.io/v1
kind: Cluster
metadata:
  generateName: cls
spec:
  displayName: test
  tenantID: default
  clusterCIDR: 10.244.0.0/16
  networkDevice: eth0
  features:
    enableMetricsServer: true
    enableCilium: false
    platformApps:     # define apps in cluster
    - name: demo      # app name
      enable: true    # if enabled, for UI
      chart:
        name: demo    # chart name
        tenantID: default
        chartGroupName: public
        version: 1.0.0
        targetNamespace: default
        values: 'key2: val2-override' # helm chart values
  properties:
    maxClusterServiceNum: 256
    maxNodePodNum: 256
  type: Baremetal
  version: 1.20.4-tke.1
  machines:
  - ip: your_ip
    port: 22
    username: root
    privateKey:
    password:
    labels: {}
```


## PR

## Reference