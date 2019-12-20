# How to run TKE locally
This guide will walk you through deploying the full TKE stack on you local machine and allow you to play with the core components. It is highly recommended if you want to develop TKE and contribute regularly.

## Table of Contents

- [Prerequisites](#prerequisites)
  - [OS Requirements](#os-requirements)
  - [Docker](#docker)
  - [etcd](#etcd)
  - [Go](#go)
  - [Node.js and NPM](#node-js-and-npm)
- [Building TKE Components](#building-tke-components)
- [Create Self-signed Certificates](#create-self-signed-certificates)
- [Create Static Token](#create-static-token)
- [Bootstrap TKE Core Components](#bootstrap-tke-core-components)
  - [tke-auth](#tke-auth)
  - [tke-platform-api](#tke-platform-api)
  - [tke-platform-controller](#tke-platform-controller)
  - [tke-registry-api(Optional)](#tke-registry-apioptional)
  - [tke-business-api(Optional)](#tke-business-apioptional)
  - [tke-business-controller(Optional)](#tke-business-controlleroptional)
  - [tke-monitor-api(Optional)](#tke-monitor-apioptional)
  - [tke-monitor-controller(Optional)](#tke-monitor-controlleroptional)
  - [tke-notify-api(Optional)](#tke-notify-apioptional)
  - [tke-notify-controller(Optional)](#tke-notify-controlleroptional)
  - [tke-gateway](#tke-gateway)
- [Access TKE Web UI](#access-tke-web-ui)
- [FAQ](#faq)

## Prerequisites

### OS Requirements
TKE supports running on `Linux`, `Windows` or `macOS` operating systems.

### Docker
TKE requires [Docker](https://docs.docker.com/installation/#installation) version 1.12+ to
run its underlying services as docker containers. Ensure the Docker daemon is working by running `docker ps` and check its version by running `docker --version`.

To install Docker,
  * **macOS:** Use either "Docker for Mac" or “docker-machine”. See instructions [here](https://docs.docker.com/docker-for-mac/).
  * **Linux:**  Find instructions to install Docker for your Linux OS [here](https://docs.docker.com/installation/#installation).

### etcd

[etcd](https://github.com/coreos/etcd/releases) is a persistent non-sql
database. TKE services share a running etcd as backend.

To install etcd,
  * **macOS:** Install and start etcd as a local service
  ```sh
  brew install etcd
  brew services start etcd
  ```

  * **Linux:** Run a single node etcd using docker. See instructions [here](https://github.com/etcd-io/etcd/blob/master/Documentation/op-guide/container.md#running-a-single-node-etcd-1).

### Go

TKE is written in [Go](https://golang.org). See supported version [here](development.md#go).

To install go,
- For macOS users,
  ```sh
  brew install go
  ```
- For other users, see instructions [here](https://golang.org/doc/install).

To configure go,

- Make sure your `$GOPATH`, `$GORROT` and `$PATH` are configured correctly
- Add `tkestack.io` to your Go env as below.
  ```sh
  go env -w GOPRIVATE="tkestack.io"
  go env -w GONOPROXY="tkestack.io"
  ```

### Node.js and NPM

TKE requires Node.js and NPM. See [here](development.md#nodejs) for supported
versions.

- For macOS users,
  ```sh
  brew install nodejs
  ```
- For other users, see instructions
 [here](https://nodejs.org/en/download/package-manager/).

## Building TKE Components
TKE contains 11 core components, a dependency list generator and a customized installer. For detail see [here](/cmd/README.md).

- Clone TKE Repository

  ```
  git clone --depth=1 https://github.com/tkestack/tke.git
  ```

  `--depth=1` parameter is optional and will ensure a smaller download.

- Build binaries

  Once all the dependencies and requirements have been installed and configured,
  you can start compiling TKE on your local machine. Make sure to run it at the TKE root path.
  ```sh
  cd tke
  make build
  ```

  After the compilation is complete, you can get all the binary executables in
the `_output/${host_os}/${host_arch}` directory.

## Create Self-signed Certificates

For security reasons, all TKE core components don't support insecure
HTTP protocol. To enable SSL, you need to make a self-signed root
certificate and a server certificate.

It is highly recommended to use the [mkcert](https://github.com/FiloSottile/mkcert) to
generate certificates for developing and testing TKE, which simplifies the process to create certificates.
see [here](https://github.com/FiloSottile/mkcert#installation) for installation guide.

To create cert using `mkcert`,

```sh
cd tke
mkdir -p _debug/certificates
cd _debug/certificates
# Make a CA and install it to local trusted certificate store.
mkcert -install
# Make server certificate.
mkcert localhost 127.0.0.1 ::1
```

You can find your certificates at

```
_debug/certificates/
├── localhost+2-key.pem
└── localhost+2.pem

0 directories, 2 files
```

## Create Static Token
Create a static token to authenticate all TKE API services.

```
cd tke
mkdir -p _debug
touch _debug/token.csv
echo 'token,admin,1,"administrator"' > _debug/token.csv
  ```

## Bootstrap TKE Core Components
This section will walk you through how to bootstrap TKE on your local machine.

TKE contains 11 core components. For detail see [here](/tke/cmd/README.md). In order for all the
services to run properly, please make sure to follow the guide below to bootstrap them in order.
You could skip the optional components if it is not needed.

For your convenient,
- Run the following command in the TKE root directory
- Export `${host_os}` and `${host_arch}` to your environment variables according to your
 machine. You can find it in your `tke/_output/${host_os}/${host_arch}` path.
- Export `${root_store}` to reference the path of your root certificate created by mkcert in the
previous step. For macOS, the path is usually /Users/${username}/Library/Application Support/mkcert.

### tke-auth

- Create `_debug/auth.json`

  <details>
  <summary>Click to show sample confi </summary>
  <br>

  **_debug/auth.json**
  ```json
  {
    "secure_serving": {
      "tls_cert_file": "_debug/certificates/localhost+2.pem",
      "tls_private_key_file": "_debug/certificates/localhost+2-key.pem"
    },
    "etcd": {
      "servers": ["http://127.0.0.1:2379"]
    },
    "generic": {
      "external_hostname": "localhost",
      "external_port": 9451
    },
    "auth": {
      "assets_path": "./pkg/auth/web",
      "tenant_admin": "admin",
      "tenant_admin_secret": "secret",
      "init_client_id": "client",
      "init_client_secret": "secret",
      "init_client_redirect_uris": [
        "http://localhost:9442/callback",
        "http://127.0.0.1:9442/callback",
        "https://localhost:9441/callback",
        "https://127.0.0.1:9441/callback"
      ]
    }
  }
  ```
  </details>


- Run `tke-auth`

  ```sh
  $ _output/${host_os}/${host_arch}/tke-auth -C _debug/auth.json
  ```

### tke-platform-api

- Create `_debug/platform-api.json`

  <details>
  <summary>Click to view sample config</summary>
  <br>

  ***_debug/platform-api.json***

  ```json
  {
    "authentication": {
      "oidc": {
        "client_id": "client",
        "issuer_url": "https://localhost:9451/oidc",
        "ca_file": "${root_store}/mkcert/rootCA.pem",
        "username_prefix": "-",
        "username_claim": "name",
        "tenantid_claim": "federated_claims"
      },
      "token_auth_file": "_debug/token.csv"
    },
    "secure_serving": {
      "tls_cert_file": "_debug/certificates/localhost+2.pem",
      "tls_private_key_file": "_debug/certificates/localhost+2-key.pem"
    },
    "etcd": {
      "servers": ["http://127.0.0.1:2379"]
    }
  }
  ```
  </details>


- Run `tke-platform-api`

  ```sh
  $ _output/${host_os}/${host_arch}/tke-platform-api -C _debug/platform-api.json
  ```

### tke-platform-controller

- Create `_debug/platform-api-client-config.yaml`

  <details>
  <summary>Click to view sample config</summary>
  <br>

  ***_debug/platform-api-client-config.yaml***

  ```yaml
  apiVersion: v1
  kind: Config
  clusters:
    - name: tke
      cluster:
        certificate-authority: ${root_store}/mkcert/rootCA.pem
        server: https://127.0.0.1:9443
  users:
    - name: admin
      user:
        token: token
  current-context: tke
  contexts:
    - context:
        cluster: tke
        user: admin
      name: tke
  ````
  </details>


- Create `_debug/platform-controller.json`

  <details>
  <summary>Click to view sample config</summary>
  <br>

  ***_debug/platform-controller.json***

  ```json
  {
    "secure_serving": {
      "tls_cert_file": "_debug/certificates/localhost+2.pem",
      "tls_private_key_file": "_debug/certificates/localhost+2-key.pem"
    },
    "client": {
      "platform": {
        "api_server_client_config": "_debug/platform-api-client-config.yaml"
      }
    }
  }
  ```
  </details>


- Run `tke-platform-controller`

  ```sh
  $ _output/${host_os}/${host_arch}/tke-platform-controller -C _debug/platform-controller.json
  ```

### tke-registry-api(Optional)

- Create `_debug/registry-api.json`

  <details>
  <summary>Click to view sample config</summary>
  <br>

  ***_debug/registry-api.json***

  ```json
  {
    "authentication": {
      "oidc": {
        "client_id": "client",
        "issuer_url": "https://localhost:9451/oidc",
        "ca_file": "${root_store}/mkcert/rootCA.pem",
        "token_review_path": "/auth/authn",
        "username_prefix": "-",
        "username_claim": "name",
        "tenantid_claim": "federated_claims"
      },
      "requestheader": {
        "username_headers": "X-Remote-User",
        "group_headers": "X-Remote-Groups",
        "extra_headers_prefix": "X-Remote-Extra-",
        "client_ca_file": "${root_store}/mkcert/rootCA.pem"
      },
      "token_auth_file": "_debug/token.csv"
    },
    "secure_serving": {
      "tls_cert_file": "_debug/certificates/localhost+2.pem",
      "tls_private_key_file": "_debug/certificates/localhost+2-key.pem"
    },
    "etcd": {
      "servers": [
        "http://127.0.0.1:2379"
      ]
    },
    "registry_config": "_debug/registry-config.yaml"
  }
  ```
  </details>


- Create `registry-config.yaml`

  <details>
  <summary>Click to view sample config</summary>
  <br>

  ***registry-config.yaml***

  ```yaml
  apiVersion: registry.config.tkestack.io/v1
  kind: RegistryConfiguration
  storage:
    fileSystem:
      rootDirectory: _debug/registry
  security:
    tokenPrivateKeyFile: keys/private_key.pem
    tokenPublicKeyFile: keys/public.crt
    adminPassword: secret
    adminUsername: admin
    httpSecret: secret
  defaultTenant: default
  ```
  </details>


- Run `tke-registry-api`

  ```sh
  $ _output/${host_os}/${host_arch}/tke-registry-api -C _debug/registry-api.json
  ```

### tke-business-api(Optional)

- Create `_debug/business-api.json`

  <details>
  <summary>Click to view sample config</summary>
  <br>

  ***_debug/business-api.json***

  ```json
  {
    "authentication": {
      "oidc": {
        "client_id": "client",
        "issuer_url": "https://localhost:9451/oidc",
        "ca_file": "${root_store}/mkcert/rootCA.pem",
        "username_prefix": "-",
        "username_claim": "name",
        "tenantid_claim": "federated_claims"
      },
      "token_auth_file": "_debug/token.csv"
    },
    "secure_serving": {
      "tls_cert_file": "_debug/certificates/localhost+2.pem",
      "tls_private_key_file": "_debug/certificates/localhost+2-key.pem"
    },
    "etcd": {
      "servers": ["http://127.0.0.1:2379"]
    },
    "client": {
      "platform": {
        "api_server_client_config": "_debug/platform-api-client-config.yaml"
      }
    }
  }
  ```
  </details>


- Run `tke-business-api`

  ```sh
  $ _output/${host_os}/${host_arch}/tke-business-api -C _debug/business-api.json
  ```

### tke-business-controller(Optional)

- Create `_debug/business-api-client-config.yaml`

  <details>
  <summary>Click to view sample config</summary>
  <br>

  ***_debug/business-api-client-config.yaml***

  ```yaml
  apiVersion: v1
  kind: Config
  clusters:
    - name: tke
      cluster:
        certificate-authority: ${root_store}/mkcert/rootCA.pem
        server: https://127.0.0.1:9447
  users:
    - name: admin
      user:
        token: token
  current-context: tke
  contexts:
    - context:
        cluster: tke
        user: admin
      name: tke
  ```
  </details>


- Create `_debug/business-controller.json`

  <details>
  <summary>Click to view sample config</summary>
  <br>

  ***_debug/business-controller.json***

  ```json
  {
    "secure_serving": {
      "tls_cert_file": "_debug/certificates/localhost+2.pem",
      "tls_private_key_file": "_debug/certificates/localhost+2-key.pem"
    },
    "client": {
      "platform": {
        "api_server_client_config": "_debug/platform-api-client-config.yaml"
      },
      "business": {
        "api_server_client_config": "_debug/business-api-client-config.yaml"
      }
    }
  }
  ```
  </details>


- Run `tke-business-controller`

  ```sh
  $ _output/${host_os}/${host_arch}/tke-business-controller -C _debug/business-controller.json
  ```

### tke-monitor-api(Optional)

- Run influxDB docker container

  `tke-monitor-controller` requires a influxDB with database name "projects" as backend to store the monitoring data.

  ```
  sudo docker volume create influxdb
  sudo docker run -d -p 8086:8086  --volume=influxdb:/var/lib/influxdb  --name influxdb influxdb:latest
  curl -XPOST 'http://localhost:8086/query' --data-urlencode 'q=CREATE DATABASE "projects"'
  ```

- Create `_debug/monitor-config.yaml`

  <details>
  <summary>Click to view sample config</summary>
  <br>

  ***_debug/monitor-config.yaml***

  ```yaml
  apiVersion: monitor.config.tkestack.io/v1
  kind: MonitorConfiguration
  storage:
    influxDB:
      servers:
        - address: http://localhost:8086
  ```
  </details>


- Cerate `_debug/monitor-api-client-config.yaml`

  <details>
  <summary>Click to view sample config</summary>
  <br>

  ***_debug/monitor-api-client-config.yaml***

  ```yaml
  apiVersion: v1
  kind: Config
  clusters:
    - name: tke
      cluster:
        certificate-authority: ${root_store}/mkcert/rootCA.pem
        server: https://127.0.0.1:9455
  users:
    - name: admin
      user:
        token: token
  current-context: tke
  contexts:
    - context:
        cluster: tke
        user: admin
      name: tke

  ```
  </details>


- Cerate `_debug/monitor-api.json`

  <details>
  <summary>Click to view sample config</summary>
  <br>

  ***_debug/monitor-api.json***

  ```json
  {
    "authentication": {
      "oidc": {
        "client_id": "client",
        "issuer_url": "https://localhost:9451/oidc",
        "ca_file": "${root_store}/mkcert/rootCA.pem",
        "username_prefix": "-",
        "username_claim": "name",
        "tenantid_claim": "federated_claims"
      },
      "token_auth_file": "_debug/token.csv"
    },
    "secure_serving": {
      "tls_cert_file": "_debug/certificates/localhost+2.pem",
      "tls_private_key_file": "_debug/certificates/localhost+2-key.pem"
    },
    "etcd": {
      "servers": ["http://127.0.0.1:2379"]
    },
    "client": {
      "platform": {
        "api_server_client_config": "_debug/platform-api-client-config.yaml"
      }
    },
    "monitor_config": "_debug/monitor-config.yaml"
  }

  ```
  </details>


- Run `tke-monitor-api`

  ```sh
  $ _output/${host_os}/${host_arch}/tke-monitor-api -C _debug/monitor-api.json
  ```

### tke-monitor-controller(Optional)

- Cerate `_debug/monitor-controller.json`

  <details>
  <summary>Click to view sample config</summary>
  <br>

  ***_debug/monitor-controller.json***

  Delete the business block if you didn't enable the TKE Business Service previously.

  ```json
  {
    "secure_serving": {
      "tls_cert_file": "_debug/certificates/localhost+2.pem",
      "tls_private_key_file": "_debug/certificates/localhost+2-key.pem"
    },
    "client": {
      "monitor": {
        "api_server_client_config": "_debug/monitor-api-client-config.yaml"
      },
      "business": {
        "api_server_client_config": "_debug/business-api-client-config.yaml"
      }
    },
    "monitor_config": "_debug/monitor-config.yaml"
  }

  ```
  </details>


- Run `tke-monitor-controller`

  ```sh
  $ _output/${host_os}/${host_arch}/tke-monitor-controller -C _debug/monitor-controller.json
  ```

### tke-notify-api(Optional)

- Create `_debug/notify-api.json`

  <details>
  <summary>Click to view sample config</summary>
  <br>

  ***_debug/notify-api.json***

  ```json
  {
    "authentication": {
      "oidc": {
        "client_id": "client",
        "issuer_url": "https://localhost:9451/oidc",
        "ca_file": "${root_store}/mkcert/rootCA.pem",
        "username_prefix": "-",
        "username_claim": "name",
        "tenantid_claim": "federated_claims"
      },
      "requestheader": {
        "username_headers": "X-Remote-User",
        "group_headers": "X-Remote-Groups",
        "extra_headers_prefix": "X-Remote-Extra-",
        "client_ca_file": "${root_store}/mkcert/rootCA.pem"
      },
      "token_auth_file": "_debug/token.csv"
    },
    "secure_serving": {
      "tls_cert_file": "_debug/certificates/localhost+2.pem",
      "tls_private_key_file": "_debug/certificates/localhost+2-key.pem"
    },
    "etcd": {
      "servers": ["http://127.0.0.1:2379"]
    },
    "client": {
      "platform": {
        "api_server_client_config": "_debug/platform-api-client-config.yaml"
      }
    }
  }

  ```
  </details>


### tke-notify-controller(Optional)

- Cerate `_debug/notify-controller.json`

  <details>
  <summary>Click to view sample config</summary>
  <br>

  ***_debug/notify-controller.json***

  ```json
  {
    "secure_serving": {
      "tls_cert_file": "_debug/certificates/localhost+2.pem",
      "tls_private_key_file": "_debug/certificates/localhost+2-key.pem"
    },
    "client": {
      "notify": {
        "api_server_client_config": "_debug/platform-api-client-config.yaml"
      }
    }
  }

  ```
  </details>


- Run `tke-notify-api`

  ```sh
  $ _output/${host_os}/${host_arch}/tke-notify-api -C _debug/notify-api.json
  ```

### tke-gateway

- Cerate `_debug/gateway-config.yaml`

  <details>
  <summary>Click to view sample config</summary>
  <br>

  ***_debug/gateway-config.yaml***

  Depending on what TKE optional services you have started, uncomment the corresponding code to allow tke-gateway to discover optional services.

  ```yaml
  apiVersion: gateway.config.tkestack.io/v1
  kind: GatewayConfiguration
  components:
    auth:
      address: https://127.0.0.1:9451
      passthrough:
        caFile: ${root_store}/mkcert/rootCA.pem
    platform:
      address: https://127.0.0.1:9443
      passthrough:
        caFile: ${root_store}/mkcert/rootCA.pem
    ### Optional Services ###
    # TKE Registry
    # registry:
    #   address: https://127.0.0.1:9453
    #   passthrough:
    #     caFile: ${root_store}/mkcert/rootCA.pem
    # TKE Business
    # business:
    #   address: https://127.0.0.1:9447
    #   frontProxy:
    #     caFile: ${root_store}/mkcert/rootCA.pem
    #     clientCertFile: certificates/localhost+2-client.pem
    #     clientKeyFile: certificates/localhost+2-client-key.pem
    # TKE Monitor
    # monitor:
    #   address: https://127.0.0.1:9455
    #   passthrough:
    #     caFile: ${root_store}/mkcert/rootCA.pem
    # TKE Notify
    # notify:
    #   address: https://127.0.0.1:9457
    #   passthrough:
    #         caFile: ${root_store}/mkcert/rootCA.pem

  ```
  </details>


- Cerate `_debug/gateway.json`

  <details>
  <summary>Click to view sample config</summary>
  <br>

  ***_debug/gateway.json***

  ```json
  {
    "authentication": {
      "oidc": {
        "client_secret": "secret",
        "client_id": "client",
        "issuer_url": "https://localhost:9451/oidc",
        "ca_file": "${root_store}/mkcert/rootCA.pem",
        "username_prefix": "-",
        "username_claim": "name",
        "tenantid_claim": "federated_claims"
      }
    },
    "secure_serving": {
      "tls_cert_file": "_debug/certificates/localhost+2.pem",
      "tls_private_key_file": "_debug/certificates/localhost+2-key.pem"
    },
    "gateway_config": "_debug/gateway-config.yaml"
  }
  ```
  </details>


- Run `tke-gateway`

  ```sh
  $ _output/${host_os}/${host_arch}/tke-gateway -C _debug/gateway.json
  ```

## Access TKE Web UI

Once all the TKE services are up and running, you can access TKE Web UI from your browser:
  * [http://localhost:9442](http://localhost:9442)
  * [https://localhost:9441](https://localhost:9441)

The username and password are specified in the launch configuration of
the `tke-auth` component:
  * ***Username:*** admin
  * ***Password:*** secret

## FAQ
**> Question:** How do I get the `DEBUG` log?

**Answer:** By default, all the core components have `INFO` level log. You can add the following block to your json config to enable `DEBUG` log.
```
"log": {
  "level": "debug"
}
```

**> Question:** How do I find the config options of TKE services?

**Answer:** Instead of using `-C` to pass the configuration file to run TKE services, you can simply use `-h` to get a full list of options.
