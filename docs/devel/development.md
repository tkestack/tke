# Development Guide

**Table of Contents**

- [Requirements](#requirements)
    - [Go](#go)
    - [Node.js](#node.js)
- [Clone source code](#clone-source-code)
- [Building binary](#building-binary)
- [Building docker image](#building-docker-image)

This document is the canonical source of truth for things like supported
toolchain versions for building TKE.

Please submit an [issue] on Github if you
* Notice a requirement that this doc does not capture.
* Find a different doc that specifies requirements (the doc should instead link
  here).

## Requirements

TKE development helper scripts assume an up-to-date GNU tools environment. 
Recent Linux distros should work out-of-the-box.

macOS ships with outdated BSD-based tools. We recommend installing [macOS GNU
tools].

Note that Mingw64/Cygwin on Windows is not supported yet.

### Go

TKE's backend is written in [Go](http://golang.org). If you don't have a Go
development environment, please [set one up](http://golang.org/doc/code.html).

| TKE             | requires Go       |
|-----------------|-------------------|
| 0.8-0.12        | 1.12.5            |
| 1.0+            | 1.12.5 - 1.13.3   |

> TKE uses [go modules](https://github.com/golang/go/wiki/Modules) to manage 
>dependencies. If your go version is `1.12.x`, please do not place the code 
>under `$GOPATH/src`.

### Node.js

TKE's frontend is written in [Typescript](https://www.typescriptlang.org/). 
To bundle TKE's frontend code, you need a Node.js and NPM execution environment, 
please [set one up](https://nodejs.org/en/download/package-manager/).

| TKE             | requires Node.js  | requires NPM  |
|-----------------|-------------------|---------------|
| 0.8-0.12        | 9.4+              | 5.6+          |
| 1.0+            | 10.3+             | 6.1+          |

## Clone source code

```sh
# Clone the repository on your machine
git clone git@github.com:tkestack/tke.git

# If you don't have a SSH key, feel free to clone using HTTPS instead
# git clone https://github.com/tkestack/tke.git
```

## Building binary

The following section is a quick start on how to build TKE on a local OS/shell 
environment.

```sh
make build
```

The best way to validate your current setup is to build a small part of TKE. 
This way you can address issues without waiting for the full build to complete. 
To build a specific part of TKE use the `BINS` environment variable to let the 
build scripts know you want to build only a certain package/executable.

```sh
make build BINS=${package_you_want}
make build BINS="${package_you_want_1} ${package_you_want_2}"
```

*Note:* This applies to all top level folders under tke/cmd.

So for the tke-gateway, you can run:

```sh
make build BINS=tke-gateway
```

If everything checks out you will have an executable in the `_output/{platform}` 
directory to play around with.

*Note:* If you are using `CDPATH`, you must either start it with a leading 
colon, or unset the variable. The make rules and scripts to build require the 
current directory to come first on the CD search path in order to properly 
navigate between directories.

```sh
cd ${working_dir}/tke
make
```

To build binaries for all platforms:

```sh
make build.all
```

To build a specific os/arch of TKE use the `PLATFORMS` environment variable to 
let the build scripts know you want to build only for os/arch.

```sh
make build.all PLATFORMS="linux_amd64 windows_amd64 darwin_amd64"
```

## Building docker image

In a production environment, it is recommended to run the TKE components in the 
[Kubernetes](https://kubernetes.io/) cluster. TKE will build and push the image 
to the [Docker Hub](https://cloud.docker.com/u/tkestack/repository/list) after 
each release.

If you don't have docker installed, see [here](running-locally.md#docker) for 
instructions on how to install on different operating systems.

If you need to build a container image locally, you can simply execute the 
instructions:

```sh
make image
```

If you want to build a container image of just one or more components, you can 
use `IMAGES` variables to control:

```sh
make image IMAGES=${package_you_want}
make image IMAGES="${package_you_want_1} ${package_you_want_2}"
```

*Note:* This applies to all top level folders under tke/cmd.

So for the tke-platform-api, you can run:

```sh
make image IMAGES=tke-platform-api
```
