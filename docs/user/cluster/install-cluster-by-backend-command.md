## Creating clusters by backend commands cURL and kubectl

### Situations

When users want to create clusters through backend commands instead of console，this guide can help.
- create global cluster
- create business clusters
- delete business clusters
- add nodes into clusters
- remove nodes from clusters

### Requirement

Please refer to [ installation requirements](../../guide/zh-CN/installation/installation-requirement.md) for information.


### How to do

- Download the installation package and run the package in the installer machine.For example download v1.9.0 amd package.

```
arch=amd64 version=v1.9.0 && wget https://tke-release-1251707795.cos.ap-guangzhou.myqcloud.com/tke-installer-linux-$arch-$version.run{,.sha256}
```
```shell
arch=amd64 version=v1.9.0 && sha256sum --check --status tke-installer-linux-$arch-$version.run.sha256 && chmod +x tke-installer-linux-$arch-$version.run && ./tke-installer-linux-$arch-$version.run

```
- Create TKEStack global cluster
```
curl -v  'http://127.00.1:8080/api/cluster'   -H 'Connection: keep-alive'   -H 'Accept: application/json, text/plain, */*' -H 'Content-Type: application/json;charset=UTF-8'   -H 'Origin: http://127.0.0.1:8080'   -H 'Referer: http://127.0.0.1:8080/'   -H 'Accept-Language: zh-CN,zh;q=0.9' -d @global-cluster.json
```
Refer to the `Configuration files`chapter at the end for `global-cluster.json`.
- Create kubeconfig to access the cluster

    1、Create service account
    ```
    kubectl create serviceaccount k8sadmin -n kube-system
    ```
    2、Create clusterrolebinding use k8sadmin and clusterrole-admin
    ```
    kubectl create clusterrolebinding k8sadmin --clusterrole=cluster-admin --serviceaccount=kube-system:k8sadmin
    ```
    3、Get the token
    ```
    kubectl -n kube-system describe secret $(sudo kubectl -n kube-system get secret | (grep k8sadmin || echo "$_") | awk '{print $1}') | grep token: | awk '{print $2}'
    ```
    4、Create kubeconfig file named user.config.Replace the token with you created in steps 3.
    Other field like `certificate-authority-data`,`server`you can get from your global cluster kubeconfig.
    ```yaml
    apiVersion: v1
    clusters:
    - cluster:
        certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM2VENDQWRHZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQ0FYRFRJeU1EVXdOVEEwTlRJeE5sb1lEekl3TlRJd05ESTNNRFExTWpFMldqQVZNUk13RVFZRApWUVFERXdwcmRXSmxjbTVsZEdWek1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBCnp2dmhPd1dCTGJqTFl2TVd6U0FWdTN6bGt3dUdkMnNMZ2FSVmxUcXhQMktGWW9PM1ZQN2tYb0VDV3A1MmFZUzcKSCtxcHhBbzJNYkFiOHlOZmNoa2M2YnBwNmhkUklTUzRKSWZlQXhPSXRWVnpzcjRYM2I4a2R2Nk9PYjNXczZQcwp4anBmSmdyWWp3TG9BNCtuVFIzdHB5VElWYXEzb1BidmtVMHk5eTViRmlZOEMyYmNnSGU1TmRkcEFGYVY3cE04ClJIc3lwcndIdjdmWUFDVkVDd21xQUIyUjMyNE5wamd4am5tVmp2TFhRYmpzbjlxMVdPLy9wbEh3ZGh1U2dYcEoKd1FTdXZ4YnhkQSt4V1JBRUlXOTdPSW80ai9PQkpWc0EvZzREL1l3ZytSNExJUWM5M3hnZXBldDY2TnpxRDZKcApyMkV4UHBVeUNkNmV0WkdTS2IxOGRRSURBUUFCbzBJd1FEQU9CZ05WSFE4QkFmOEVCQU1DQXFRd0R3WURWUjBUCkFRSC9CQVV3QXdFQi96QWRCZ05WSFE0RUZnUVVKRFBPTlhrcFUzb01GWWtBWkhSZUp6cHNzVWd3RFFZSktvWkkKaHZjTkFRRUxCUUFEZ2dFQkFDQnhENis3UUw2aWpmY0NHTUhtZkVqNHh5dFVpdWJNM0RVMmJZUU82SklyK0FjQgpHdVVEc0lhdDhzUUtCb0YzeHZobVNhYzNHYmFURGVmVW85K2dhM0Q2ZEZ0UTVFMmFoR21JenRaaURUZ29aVVJVCkIrTklpQVdZVlhydWhGUUpONHREK3A3MVNTNXhPV1UvR2tXNmRxRXc3c3p6OVdobFNSOTg4eURIUWFlV2ZqSXIKbGcvQTN2U0YxckZRd2IwYVdUVGhWaUVjWUovTkhQSldkdDRSRG0rRlcrV0xVUFcvVzdNOW5jUGRERGFFaWt2NQo0YXFEMUFnclZSWkNrV1VsZExWQ3NEdVJROEhHWU1lVitFWUxtNDVwY0ZFa3pxTUtnVXUvdTNhUC8zY3ZYQllSClg3SU8wa2Iwd2lPajNTOGZudHNCS2RIWTBhZUFZdnVjUDIrY3M5QT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
        server: https://10.0.71.117:6443
      name: global
    contexts:
    - context:
        cluster: global
        user: global-admin
      name: global-context-default
    current-context: global-context-default
    kind: Config
    preferences: {}
    users:
    - name: global-admin
      user:
        token: eyJhbGciOiJSUzI1NiIsImtpZCI6Ikd0VTY2TVpzd1NNeUI1TkRrYU5KYU04aXlvZ1hnWVdleG5jUXlQVnREa1EifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrOHNhZG1pbi10b2tlbi13aG50eiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJrOHNhZG1pbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImU0ZjYwZjAwLTY3NjMtNGY1MS1hMjY2LTU5ZWQ1N2RkZWFlNiIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTprOHNhZG1pbiJ9.pL4BIFyyHfm-hprkHJJdn6euVOQ_HdJkcnkxR6WymR_Q1KAe4HA8RUbq_ujIhEfgCBUVhLjcQsXKbUhFYxoKQN3rPfyfNIsO8NbgOi_hI3zIayH99Z9NC9zwCtkKj_6iEA4kxb0DFYvYkpLlUYj8YfqxOu0n-a7XPWOnvqvLgGfEoQD0bATC4CweSJOU-xpFIpKYsYvAMh4IOr-VpA-YTAn8-ddBvfvTqzP7k9x-VJ0wp-10JgPHp0MNT7rZazZL9klesZC-0mEdrdHutyD_DU8U9gRCUCxLMXBqaJrPiNAwgB_HXTNsvUL5vUVOZDjJtNo-rDqU3nIkaW8l0dXdaA
    ```
- Create business cluster
```
kubectl --kubeconfig=./user.config create -f business-cluster.json -o json
```
Refer to the `Configuration files`chapter at the end for `business-cluster.json`.
- Delete business cluster
```
kubectl --kubeconfig=./user.config delete -f business-cluster.json
```
- Add node to cluster
```
kubectl --kubeconfig=./user.config apply -f machines.json -o yaml
```
Refer to the `Configuration files`chapter at the end for `machines.json`.
- Remove node from cluster
```
kubectl --kubeconfig=./user.config delete -f machines.json
```

### Configuration files
Make sure fill all the parameter in annotation based on your hosts and your configuration.
Make sure remove all the annotations on the json files before you use it.
-  golbal-cluster.json
```
{
    "cluster":{
        "apiVersion":"platform.tkestack.io/v1",
        "kind":"Cluster",
        "spec":{
            "networkDevice":"eth0",  // choose the network device, default eth0
            "features":{
                "enableMetricsServer":true, // enable MetricsServer
                "enableCilium":false
            },
            "dockerExtraArgs":{

            },
            "kubeletExtraArgs":{

            },
            "apiServerExtraArgs":{

            },
            "controllerManagerExtraArgs":{

            },
            "schedulerExtraArgs":{

            },
            "clusterCIDR":"192.168.0.0/16", // cluster CIDR: 192.168.0.0/16~19;172.16.0.0/16~19;10.0.0.0/14~19
            "properties":{
                "maxClusterServiceNum":256, // can use 128，256，512，1024，2048，4096，8192，16384，32768
                "maxNodePodNum":256 // can use 32，64，128，256
            },
            "type":"Baremetal",
            "machines":[
                {
                    "ip":"10.0.0.999", // node IP address,when you have more than one node,make sure all the nodes can ssh each other
                    "port":22,  // ssh port
                    "username":"ubuntu", // username for the node to ssh
                    "privateKey":"xxxxx" // private key when ssh the node
                },
                {  // if you only have one node,remove it
                    "ip":"10.0.0.104",
                    "port":22,
                    "username":"ubuntu",
                    "privateKey":"xxxxx"
                },
                { // if you only have one node,remove it
                    "ip":"10.0.0.218",
                    "port":22,
                    "username":"ubuntu",
                    "privateKey":"xxxxx"
                }
            ]
        }
    },
    "config":{
        "basic":{
            "username":"admin", // cluster username you login
            "password":"YWRtaW4=" // password you login the cluster use echo -n "<yourpassword>" | base64
        },
        "auth":{
            "tke":{

            }
        },
        "registry":{
            "tke":{
                "domain":"registry.tke.com" // registry domain name
            }
        },
        "application":{ // enable application

        },
        "business":{ //enable business

        },
        "audit":{ //enable audit with fake username and password.Make sure to config tke-audit configmap with correct values after install.
            "elasticSearch":{
                "address":"http://10.0.0.1:9200",
                "username":"skipTKEAuditHealthCheck",
                "password":"MTIzNDU2",
                "reserveDays":7
            }
        },
        "monitor":{
            "influxDB":{ // use influxDB
                "local":{

                }
            }
        },
        "logagent":{ // enable logagent

        },
        "ha":{ // enable ha
            "tke":{
                "vip":"xxxxx" // vip address tke provided
            }
        },
        "gateway":{
            "domain":"xxxxx", // gateway domain name
            "cert":{
                "selfSigned":{

                }
            }
        }
    }
}
```
- business-cluster.json
```
{
    "apiVersion":"platform.tkestack.io/v1",
    "kind":"Cluster",
    "metadata":{
        "generateName":"cls",
	    "name": "cls-hm" // business cluster name
    },
    "spec":{
        "displayName":"test",
        "clusterCIDR":"10.244.0.0/16", // business cluster CIDR
        "networkDevice":"eth0", // network device
        "features":{
            "containerRuntime":"containerd", //container runtime
            "enableMetricsServer":true, // enable Metrics Server
            "enableCilium":false
        },
        "properties":{
            "maxClusterServiceNum":256,
            "maxNodePodNum":256
        },
        "type":"Baremetal",
	    "tenantID": "default",
        "version":"1.21.4-tke.3",
        "machines":[
            {
                "ip":"10.0.71.210", //node ip address
                "port":22, // ssh port
                "username":"ubuntu", // username when ssh the node
                "privateKey":"xxxxx", //private key when ssh the node
                "labels":{

                }
            }
        ]
    }
}
```
- machines.json
```
{
    "kind":"Machine",
    "apiVersion":"platform.tkestack.io/v1",
    "metadata": {
    "name": "cls-hm-node1" // add node name
  },
    "spec":{
        "clusterName":"cls-hm", // the cluster name you want to add node
        "ip":"10.0.71.201", // node ip address
        "port":22, // ssh port
        "username":"ubuntu", // ssh username when login the node
        "privateKey":"xxxxxx", // private key when ssh the node
        "labels":{

        },
        "type":"Baremetal"
    }
}
```





