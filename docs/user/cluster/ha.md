# Creating Highly Available clusters

## 1.TKE HA

The underlying implementation is keepalived.

### Require

- VIP which must be available in the node network

### Usage

You can create a cluster through kubectl.

In your global cluster or setup kubeconfig for tke:

`kubectl create -f cluster.yaml`

cluster.yaml:
```yaml
{
    "apiVersion": "platform.tkestack.io/v1",
    "kind": "Cluster",
    "spec": {
        "clusterCIDR": "10.244.0.0/16",
        "type": "Baremetal",
        "version": "1.14.10",
        "machines": [
            {
                "ip": "192.168.1.2",
                "port": 22,
                "username": "root",
                "password": "MTIzNDU2" // echo -n "123456" | base64
            }
        ],
        "features": {
            "ha": {
                "tke": {
                    "vip": "192.168.1.3" // vip is required!
                }
            }
        }
    }
}
```

## 2.ThirdParty(External) HA

User can use third party HA implementation such as Tencent Load Balancer, Tencent Gateway and etc.

### Require

- Must bind real server before create cluster.

For example:

M1: 192.168.1.2
M2: 192.168.1.3
M3: 192.168.1.4

VIP: 192.168.1.5

Bind rules:

- apiserver 192.168.1.5:6443 => 192.168.1.2:6443,192.168.1.3:6443,192.168.1.4:6443

For tke web console in global cluster which configure in installer:

- tke console http 192.168.1.5:80 => 192.168.1.2:80,192.168.1.3:80,192.168.1.4:80
- tke console https 192.168.1.5:443 => 192.168.1.2:443,192.168.1.3:443,192.168.1.4:443

### Usage

cluster.yaml:
```yaml
{
    "apiVersion": "platform.tkestack.io/v1",
    "kind": "Cluster",
    "spec": {
        "clusterCIDR": "10.244.0.0/16",
        "type": "Baremetal",
        "version": "1.14.10",
        "machines": [
            {
                "ip": "192.168.1.2",
                "port": 22,
                "username": "root",
                "password": "MTIzNDU2" // echo -n "123456" | base64
            },
            {
                "ip": "192.168.1.3",
                "port": 22,
                "username": "root",
                "password": "MTIzNDU2" // echo -n "123456" | base64
            },
            {
                "ip": "192.168.1.4",
                "port": 22,
                "username": "root",
                "password": "MTIzNDU2" // echo -n "123456" | base64
            }
        ],
        "features": {
            "ha": {
                "thirdParty": {
                    "vip": "192.168.1.5", // vip is required!
                    "vport": 6443 // vport is required!
                }
            }
        }
    }
}
```

### Known problems

1. Tencent Load Balancer and Tencent Gateway can't connect self though vip.

        M1 => VIP =can't=> M1 // restricted by implementation

    Means if you use these implementations, your must bind at least **two** servers!