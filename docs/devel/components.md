# TKEStack Components

[`/cmd`](../../cmd) directory includes every TKEStack components and is where all binaries and container images are built. For detail about how to launch the TKEStack cluster locally see the guide [here](running-locally.md).

## Overview

TKEStack contains 12 core components belonging to 6 services, a dependency list generator and a customized installer.

## Core Components
To bootstrap properly, TKEStack core components need to be run in the order as shown below.

- [`tke-auth-api`](../../cmd/tke-auth-api) integrates [dex](https://github.com/dexidp/dex) to provide an [OpenID Connect](https://en.wikipedia.org/wiki/OpenID_Connect) server, which can provide access to third-party authentication systems, and also provides a default local identify.

- [`tke-auth-controller`](../../cmd/tke-auth-controller) watches the state of the auth API objects through the `tke-auth-api` and configures TKEStack auth resources.

- [`tke-platform-api`](../../cmd/tke-platform-api) is the most important service of TKEStack . It listens to and validates requests to TKEStack platform API then configures its API objects.

- [`tke-platform-controller`](../../cmd/tke-platform-controller) watches the state of the platform API objects through the `tke-platform-api` and configures TKEStack platform.

- [`tke-registry-api`](../../cmd/tke-registry-api) enables a build-in registry and chart repository of helm inside TKEStack.

- [`tke-business-api`](../../cmd/tke-business-api) enables TKEStack project management by business labels.

- [`tke-business-controller`](../../cmd/tke-business-controller) watches the state of the business API objects through the `tke-business-api` and configures TKEStack business resources.

- [`tke-monitor-api`](../../cmd/tke-monitor-api) enables TKEStack monitoring and provides a web UI to configure and view monitoring data.

- [`tke-monitor-controller`](../../cmd/tke-monitor-contoller) watches the state of the monitor API objects through the `tke-monitor-api` and configures TKEStack monitoring.

- [`tke-notify-api`](../../cmd/tke-notify-api) enables TKEStack alert notification and provides a web UI for you to configure alerts and view their status.

- [`tke-notify-controller`](cmd/tke-notify-contoller) watches the state of the notify API objects through the `tke-notify-api` and configures TKEStack notification.

-  provides a web UI to interact with TKEStack.

  > You can refer to [TKEStack architecture](../guide/zh-CN/installation/installation-architecture.md) to know more information.

## Dependency List Generator
- [`generate-images`](../../cmd/generate-images) reads from all the dependencies and generates a list of image dependencies.

## Installer
- [`tke-installer`](../../cmd/tke-installer) provides an easy way to install and launch your own TKEStack. 
- You can refer to [here](../user/tke-installer/introduction.md) for more information about tke-installer。

## Help

If you have any questions, please submit your [issue](https://github.com/tkestack/tke/issues/new/choose) or [PR](https://github.com/tkestack/tke/pulls), or you can directly contact the [components maintainers](../../MAINTAINERS.md).