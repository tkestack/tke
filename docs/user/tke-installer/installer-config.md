# Introduction

The tke-installer configuration use the json format, so data submitted in ui mode can be used as command line configuration file.
You can open the browser console through F12 and view the submitted creation request.(`POST /api/cluster`)
Or you can view by installer logs.(`docker logs tke-installer|grep 'POST /api/cluster HTTP/1.1' -A7`)

# Explanation

The configuration file is divided into two parts.

- `cluster` control the global cluster, [Reference](https://github.com/tkestack/tke/blob/master/api/platform/v1/types.go#L31)
- `config` control the tke [Reference](https://github.com/tkestack/tke/blob/master/cmd/tke-installer/app/installer/installer.go#L131)

## Sample

- one machine `sshpass -p123456 ssh root@172.19.0.2 -p22`
- login credendial provided by tke default values
  - username: admin
  - passowrd: admin
- enable tke components
  - auth
  - registry
  - business
  - monitor
  - gateway

Note: **Copy this sample should remove comment!!!**

```json
{
  "cluster": {
    "spec": {
      "machines": [ // the machines for create global cluster and tke
        {
          "ip": "172.19.0.2",
          "port": 22,
          "username": "root",
          "password": "MTIzNDU2" // password need base64. Such as, echo -n '123456'|base64
        },
        {
          "ip": "172.19.0.3",
          "port": 22,
          "username": "root",
          "password": "MTIzNDU2"
        }
      ]
    }
  },
  "Config": {
    "basic": {
        "username": "admin", // username for login tke console and registry
        "password": "MTIzNDU2" // password for login tke console and registry
    },
    "auth": { // enable auth
      "tke": { // use the tke auth
          
      }
    },
    "registry": { // enable registry
      "tke": { // use the tke registry
        "domain": "registry.tke.com",
      }
    },
    "business": { // enable business

    },
    "monitor": { // enable monitor
      "influxDB": { // use influxdb
        "local": { // use the local influxdb which provided by tke
          
        }
      }
    },
    "gateway": { // enable gateway
      "domain": "console.tke.com" // for https login
    }
  }
}
```