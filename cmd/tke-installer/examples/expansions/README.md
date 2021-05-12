# TKE expansion framework demo

## how to use

* copy directory `expansions` into a base directory, for example: `/tmp/`
* with this environment `export EXPANSION_PATH=/tmp/expansions`, you can do `go run cmd/tke-installer/installer.go`
* when you run tke-installer within a container, just put expansion files into `/opt/tke-installer/data/expansions/`, which is the default path for tke-installer running in container

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
