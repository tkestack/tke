# Introduction

The tke-installer runs in docker mode, all dependent resources and configuration files are included in the image.

Enter the conatiner: `docker exec -it tke-installer bash`

# Help

```
bash-5.0# bin/tke-installer --help
The TKE Installer is used to setup the first kubernetes cluster.

Usage:
  tke-installer [flags]

Flags:
  -C, --config FILE                      Read configuration from specified FILE, support JSON, TOML, YAML, HCL, or Java properties formats.
      --log-level LEVEL                  Minimum log output LEVEL. (default "info")
      --log-format FORMAT                Log output FORMAT, support plain or json format. (default "console")
      --log-disable-color                Disable output ansi colors in plain format logs.
      --log-enable-caller                Enable output of caller information in the log.
      --log-output-paths strings         Output paths of log (default [stdout])
      --log-error-output-paths strings   Error output paths of log (default [stderr])
      --listen-addr string               listen addr (default "172.19.0.215:8080")
      --no-ui                            run without web
      --input string                     specify input file (default "conf/tke.json")
      --force                            force run as clean
      --sync-projects-with-namespaces    Enable creating/deleting the corresponding namespace when creating/deleting a project.
      --replicas int                     tke components replicas (default 2)
  -V, --version version[=true]           Print version information and quit.
  -H, --help                             Help for Tencent Kubernetes Engine Installer.
      --log-flush-frequency duration     Maximum number of seconds between log flushes (default 5s)
```

# Package layout

```
/app
├── bin
│   ├── kubectl # extra tools for debug or hooks
│   └── tke-installer # main program
├── conf # for console mode
├── data # runtime data
│   ├── admin.crt # x509 auth for tke
│   ├── admin.key # x509 auth for tke
│   ├── admin.kubeconfig # kubeconfig for tke auth
│   ├── ca.crt # ca for tke
│   ├── ca.key # ca for tke
│   ├── etcd-ca.crt # etcd ca for tke
│   ├── etcd.crt # etcd client cert for tke
│   ├── etcd.key  # etcd client cert for tke
│   ├── oidc_client_secret # for oidc auth
│   ├── password.csv # password auth for tke
│   ├── server.crt # x509 cert for tke components
│   ├── server.key # x509 cert for tke components
│   ├── tke.json # status data of tke-installer
│   ├── tke.log # installation log
│   └── token.csv # token auth for tke
├── hooks # hooks for installation
│   ├── post-cluster-ready
│   ├── post-install
│   └── pre-install
├── hosts # for control the docker of host machine
├── images.tar.gz # all dependent images include k8s related and tke components
├── manifests # manifests of tke components
│   ├── etcd
│   │   └── etcd.yaml
│   ├── influxdb
│   │   └── influxdb.yaml
│   ├── keepalived
│   │   └── keepalived.yaml
│   ├── nginx
│   │   └── nginx.yaml
│   ├── tke-auth
│   │   └── tke-auth.yaml
│   ├── tke-business-api
│   │   └── tke-business-api.yaml
│   ├── tke-business-controller
│   │   └── tke-business-controller.yaml
│   ├── tke-gateway
│   │   └── tke-gateway.yaml
│   ├── tke-monitor-api
│   │   └── tke-monitor-api.yaml
│   ├── tke-monitor-controller
│   │   └── tke-monitor-controller.yaml
│   ├── tke-notify-api
│   │   └── tke-notify-api.yaml
│   ├── tke-notify-controller
│   │   └── tke-notify-controller.yaml
│   ├── tke-platform-api
│   │   └── tke-platform-api.yaml
│   ├── tke-platform-controller
│   │   └── tke-platform-controller.yaml
│   └── tke-registry-api
│       └── tke-registry-api.yaml
└── provider # cluster provider
    └── baremetal
        ├── conf
        │   ├── config.yaml # baremetal cluster provider config
        │   ├── docker
        │   │   ├── daemon.json # docker config for /etc/docker/daemon.json
        │   │   └── docker.service # docker systemd service
        │   ├── kubeadm
        │   │   ├── 10-kubeadm.conf # kubelet systemd service which is customized by kubeadm
        │   │   └── kubeadm-config.yaml # kubeadm config
        │   ├── kubelet
        │   │   ├── kubelet-config.conf # kubelet config
        │   │   └── kubelet.service # kubelet systemd service
        │   ├── oidc-ca.crt # oidc relative
        │   └── sysctl.conf # linux sysctl
        ├── manifests
        │   └── gpu
        │       └── nvidia-device-plugin.yaml # nvidia-device-plugin manifests
        └── res cluster's resource
            ├── NVIDIA-Linux-x86_64-440.31.run
            ├── cni-plugins-amd64-v0.7.5.tgz
            ├── docker-18.09.9.tgz
            ├── kubeadm-v1.15.1.tar.gz
            ├── kubernetes-node-linux-amd64-v1.14.6.tar.gz
            └── nvidia-container-runtime-3.1.4.tgz 
```

## FAQ

1. How to completely reinstall?

    1. In your installer machine, remove installer data.(`rm -rf /opt/tke-installer/data`)And restart installer container.(`docker restart tke-installer`)
    2. In your machines of global cluster, run the clean scripts.(`curl -s https://tke-release-1251707795.cos.ap-guangzhou.myqcloud.com/tools/clean.sh | sh`)

    The reason why we need to do this is because we want to protect the security and autonomy of user data. The user's data is deleted by the user to avoid production accidents caused by data loss

2. How to resume installation when the installation is interrupted abnormally, such as network jitter?

        docker restart tke-installer

3. How to customize docker or k8s parameters?

    Before installation, enter the tke-installer(`docker exec -it tke-installer bash`), modify related configuration files.
    Be careful with the original template variables and formatting specifications, **especially the indentation and spaces,** and be especially careful!
    The changes here will take effect on all clusters, including the **global** cluster!
    Be sure to understand what you are doing. If it is not necessary, please do not modify it. It may cause the entire installation process to fail.

    - docker
      - /app/provider/baremetal/conf/docker/daemon.json
      - /app/provider/baremetal/conf/docker/docker.service
    - kubelet
      - /app/provider/baremetal/conf/kubelet/kubelet-config.conf
      - /app/provider/baremetal/conf/kubelet/kubelet.service
    - kubeadm which control k8s
      - /app/provider/baremetal/conf/kubeadm/10-kubeadm.conf
      - /app/provider/baremetal/conf/kubeadm/kubeadm-config.yaml
    - manifests which control all tke components
      - /app/manifests

4. Whether to support command line mode?

    Installer run in ui mode by default, provide a user-friendly visual installation interface. 
    Command line mode can be enabled through no-ui, which facilitates automated installation for advanced users.
    Config should provide though an json file(`/opt/tke-insatller/conf/tke.json`) and should rerun tke-installer container.
    `version=v1.2.4 && docker run --name tke-installer -d --privileged --net=host -v/etc/hosts:/app/hosts -v/etc/docker:/etc/docker -v/var/run/docker.sock:/var/run/docker.sock -v$DATA_DIR:/app/data -v$INSTALL_DIR/conf:/app/conf --entrypoint="/app/bin/tke-installer --no-ui" tkestack/tke-installer:$version`
    [Configuration file format](./installer-config.md)

5. How to automate custom logic or integrate your own deployment?

    The installer provides a hook mechanism that can execute specific programs at specific stages of the installation process.
    You can put installation scripts in the `/app/hooks/` directory of the installer container. It should be noted that the names of the hooks must be corresponding.

    - pre-install which can be used to modify installer related configurations for customization
    - post-cluster-ready which can be used to customize tke components
    - post-install which can be used to automate your own deployment tasks or adjust global cluster or tke configuration

    If you want to automate integration hooks, you can modify the container creation parameters of tke-installer.
    Assuming your hook has been placed in `/opt/tke-installer/hooks`, You should rerun tke-installer container.
    `version=v1.2.4 && hooks=/opt/tke-installer/hooks && docker run --name tke-installer -d --privileged --net=host -v/etc/hosts:/app/hosts -v/etc/docker:/etc/docker -v/var/run/docker.sock:/var/run/docker.sock -v$DATA_DIR:/app/data -v$INSTALL_DIR/conf:/app/conf -v$hooks:/app/hooks tkestack/tke-installer:$version`
