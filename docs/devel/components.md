# TKE Components

`/cmd` directory includes every TKE components and is where all binaries and container images are built. For detail about how to launch the TKE cluster see the guide [here](/docs/running-locally.md).

## Overview

TKE contains 12 core components belonging to 6 services, a dependency list generator and a customized installer.

## Core Components
To bootstrap properly, TKE core components need to be run in the order as shown below.

- [`tke-auth-api`](/cmd/tke-auth-api) integrates [dex](https://github.com/dexidp/dex) to provide an OpenID Connect server, which can provide access to third-party authentication systems, and also provides a default local identify.
- [`tke-auth-controller`](/cmd/tke-auth-controller) watches the state of the auth api objects through the `tke-auth-api` and configures TKE auth resources.
- [`tke-platform-api`](/cmd/platform-api) is the most important service of TKE . It listens to and validates requests to TKE platform api then configures its api objects.
- [`tke-platform-controller`](/cmd/tke-platform-controller) watches the state of the platform api objects through the `tke-platform-api` and configures TKE platform.
- [`tke-registry-api`](/cmd/tke-registry-api) enables a build-in registry and chart repository of helm inside TKE.
- [`tke-business-api`](/cmd/tke-business-api) enables TKE project management by business labels.
- [`tke-business-controller`](/cmd/tke-business-controller) watches the state of the business api objects through the `tke-business-api` and configures TKE business resources.
- [`tke-monitor-api`](/cmd/tke-monitor-api) enables TKE monitoring and provides a web UI to configure and view monitoring data.
- [`tke-monitor-controller`](/cmd/tke-monitor-contoller) watches the state of the monitor api objects through the `tke-monitor-api` and configures TKE monitoring.
- [`tke-notify-api`](/cmd/tke-notify-api) enables TKE alert notification and provides a web UI for you to configure alerts and view their status.
- [`tke-notify-controller`](cmd/tke-notify-contoller) watches the state of the notify api objects through the `tke-notify-api` and configures TKE notification.
- [`tke-gateway`](/cmd/tke-gateway) provides a web UI to interact with TKE .

## Dependency List Generator
- [`generate-images`](/cmd/generate-images) reads from all the dependencies and generates a list of image dependencies.

## Installer
- [`tke-installer`](/cmd/tke-installer) provides an easy way to install and launch your own TKE.
