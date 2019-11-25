# Running Locally

**Table of Contents**

- [Requirements](#requirements)
    - [OS](#os)
    - [Docker](#docker)
    - [etcd](#etcd)
    - [Go](#go)
    - [Node.js and NPM](#nodejs-and-npm)
- [Clone the repository](#clone-the-repository)
- [Build the binary](#build-the-binary)
- [Make self-signed certificate](#make-self-signed-certificate)
- [Running on the machine](#running-on-the-machine)
    - [Create static token auth file](#create-static-token-auth-file)
    - [tke-auth](#tke-auth)
    - [tke-platform-api](#tke-platform-api)
    - [tke-platform-controller](#tke-platform-controller)
    - [tke-business-api](#tke-business-api)
    - [tke-business-controller](#tke-business-controller)
    - [tke-gateway](#tke-gateway)
- [Open UI console](#open-ui-console)

## Requirements

### OS

TKE support running on `linux`, `Windows` or `macOS` system.

### Docker

At least [Docker](https://docs.docker.com/installation/#installation) 1.12+. 
Ensure the Docker daemon is running and can be contacted (try `docker ps`).  
In order for the TKE component to work properly locally, the underlying services
 it depends on will run as a docker container.

Docker, using one of the following configurations:
  * **macOS** You can either use Docker for Mac or docker-machine. See 
  installation instructions [here](https://docs.docker.com/docker-for-mac/).
  * **Linux with local Docker**  Install Docker according to the 
  [instructions](https://docs.docker.com/installation/#installation) for your OS.

### etcd

[etcd](https://github.com/coreos/etcd/releases) is a backend persistent non-sql 
database that TKE requires for almost all components.

If etcd is not installed on the machine, in addition to the installation of the 
package management tool corresponding to the OS.
  * **macOS** You can use the following command to install and start the etcd 
  service.
  
  ```sh
  $ brew install etcd
  $ brew service start etcd
  ```

  * **Linux** You can use docker to start a single-node etcd to run in the 
  [official documentation of etcd](https://github.com/etcd-io/etcd/blob/master/Documentation/op-guide/container.md#running-a-single-node-etcd-1).

### Go

You need [go](https://golang.org/doc/install) in your path 
(see [here](development.md#go) for supported versions), please make sure it is 
installed and in your ``$PATH``.

If you use the macOS system, you can use the following command to install:

```sh
$ brew install go
```

Make sure that the `tkestack.io` domain name is added to the `GOPRIVATE` and 
`GONOPROXY` environment variables after you have an eligible go version. If not, 
you can simply execute the following command:

```sh
$ go env -w GOPRIVATE="tkestack.io"
$ go env -w GONOPROXY="tkestack.io"
```

### Node.js and NPM

You need a Node.js and NPM (see [here](development.md#nodejs) for supported 
versions) execution environment, please [set one up](https://nodejs.org/en/download/package-manager/).

If you use the macOS system, you can use the following command to install:

```sh
$ brew install nodejs
```

## Clone the repository

In order to run TKE you must have the code on the local machine. Cloning this 
repository is sufficient.

```$ git clone --depth=1 https://github.com/tkestack/tke.git```

The `--depth=1` parameter is optional and will ensure a smaller download.

In the subsequent documentation, you assume that all operations are performed in 
the code root path, so you need to switch the working path to the directory 
first.

```sh
$ cd tke
```

## Build the binary

Once all the dependencies and requirements have been installed and configured, 
you can execute `make build` in the root of the code to compile all the 
components of TKE. 

After the compilation is complete, you can get all the binary executables in 
the `_output/${host_os}/${host_arch}` directory.

## Make self-signed certificate

For security reasons, all service components of tke do not support the insecure 
HTTP protocol. In order to enable SSL, you need to make a self-signed root 
certificate, a server certificate.

To generate the certificate used to develop the test, it is highly recommended 
to use the [mkcert](https://github.com/FiloSottile/mkcert) tool, which 
simplifies the process and configuration of certificate generation. See 
[here](https://github.com/FiloSottile/mkcert#installation) for installation.

```sh
$ mkdir -p _debug/certificates
$ cd _debug/certificates
$ # Make a CA and install it to local trusted certificate store.
$ mkcert -install
$ # Make server certificate.
$ mkcert localhost 127.0.0.1 ::1
```

Then, you can get:

```
.
├── localhost+2-key.pem
└── localhost+2.pem

0 directories, 2 files
```

## Running on the machine

### Create static token auth file

First you need to create a directory `_debug` to hold all the configuration files.

```
$ mkdir -p _debug
```

Then you need to create a static token authentication file to provide static 
token authentication for all API type services.

***_debug/token.csv***

```csv
token,admin,1,"administrator"
```

### tke-auth

Generate the configuration files needed to run the `tke-auth` component.

***_debug/auth.json***

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

Running it:

```sh
$ _output/${host_os}/${host_arch}/tke-auth -C _debug/auth.json
```

### tke-platform-api

Generate the configuration files needed to run the `tke-platform-api` component.

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

> The path represented by `${root_store}` is the storage path of the root 
> certificate created by using the `mkcert` tool. 
> If it is the `macOS` operating system, the path is generally 
> `/Users/${username}/Library/Application Support/mkcert`.

Running it:

```sh
$ _output/${host_os}/${host_arch}/tke-platform-api -C _debug/platform-api.json
```

### tke-platform-controller

Generate the configuration files needed to run the `tke-platform-controller` 
component.

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
```

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

Running it:

```sh
$ _output/${host_os}/${host_arch}/tke-platform-controller -C _debug/platform-controller.json
```

### tke-business-api

Generate the configuration files needed to run the `tke-business-api` 
component.

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

Running it:

```sh
$ _output/${host_os}/${host_arch}/tke-business-api -C _debug/business-api.json
```

### tke-business-controller

Generate the configuration files needed to run the `tke-business-controller` 
component.

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

Running it:

```sh
$ _output/${host_os}/${host_arch}/tke-business-controller -C _debug/business-controller.json
```

### tke-gateway

Generate the configuration files needed to run the `tke-gateway` component.

***_debug/gateway-config.yaml***

```yaml
apiVersion: gateway.config.tkestack.io/v1
kind: GatewayConfiguration
components:
  platform:
    address: https://127.0.0.1:9443
    passthrough:
      caFile: ${root_store}/mkcert/rootCA.pem
  business:
    address: https://127.0.0.1:9447
    passthrough:
      caFile: ${root_store}/mkcert/rootCA.pem
```

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

Running it:

```sh
$ _output/${host_os}/${host_arch}/tke-gateway -C _debug/gateway.json
```

## Open UI console

Once all the components are working, you can open a browser to access:
  * [http://localhost:9442](http://localhost:9442)
  * [https://localhost:9441](https://localhost:9441)

The login username and password are specified in the launch configuration of 
the previous `tke-auth` component:
  * ***Username***: admin
  * ***Password***: secret
