# TKE expansion framework demo

## how to use

* copy directory `expansions` into a base directory, for example: `/tmp/`
* specify `"EnableCustomExpansion": true`, `"CustomExpansionDir": "/tmp/expansions"` in your tke.json, then you can do `go run cmd/tke-installer/installer.go`
* when you run tke-installer within a container, just put expansion files into `/opt/tke-installer/data/expansions/`, which is the default mount path for tke-installer running in container

## expansion layout
* directory:
    * hooks: pre-install, post-cluster-ready, post-install files which TKE-installer will use
    * TODOS:
        * applications: application yaml files
        * charts: chart tars
        * files: `copy file` origin version, then they are going to be copied to `files_generated` directory
        * provider: provider files to override defaults
* file:
    * TODOS:
        * images.tar.gz: expansion images, which will be loaded by TKE-installer and push into registry
        * values.yaml: this file describes two parts of values:
            * global values: which will be used by each chart
            * chart value map: for each will be used by the chart it specifies

## changes of `tke.json`

* EnableCustomExpansion bool
    * if true: will enable expansion. default false
* CustomExpansionDir string
    * path to expansions. default `data/expansions`
```
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
  "CustomExpansionDir": ""
 }
}
```
