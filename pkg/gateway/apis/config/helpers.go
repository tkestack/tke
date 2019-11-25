/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package config

// GatewayConfigurationPathRefs returns pointers to all of the GatewayConfiguration fields that contain filepaths.
// You might use this, for example, to resolve all relative paths against some common root before
// passing the configuration to the application. This method must be kept up to date as new fields are added.
func GatewayConfigurationPathRefs(gc *GatewayConfiguration) []*string {
	var paths []*string
	if gc.Components.Platform != nil {
		if gc.Components.Platform.Passthrough != nil {
			paths = append(paths, &gc.Components.Platform.Passthrough.CAFile)
		}
		if gc.Components.Platform.FrontProxy != nil {
			paths = append(paths, &gc.Components.Platform.FrontProxy.CAFile)
			paths = append(paths, &gc.Components.Platform.FrontProxy.ClientCertFile)
			paths = append(paths, &gc.Components.Platform.FrontProxy.ClientKeyFile)
		}
	}
	if gc.Components.Business != nil {
		if gc.Components.Business.Passthrough != nil {
			paths = append(paths, &gc.Components.Business.Passthrough.CAFile)
		}
		if gc.Components.Business.FrontProxy != nil {
			paths = append(paths, &gc.Components.Business.FrontProxy.CAFile)
			paths = append(paths, &gc.Components.Business.FrontProxy.ClientCertFile)
			paths = append(paths, &gc.Components.Business.FrontProxy.ClientKeyFile)
		}
	}
	if gc.Components.Notify != nil {
		if gc.Components.Notify.Passthrough != nil {
			paths = append(paths, &gc.Components.Notify.Passthrough.CAFile)
		}
		if gc.Components.Notify.FrontProxy != nil {
			paths = append(paths, &gc.Components.Notify.FrontProxy.CAFile)
			paths = append(paths, &gc.Components.Notify.FrontProxy.ClientCertFile)
			paths = append(paths, &gc.Components.Notify.FrontProxy.ClientKeyFile)
		}
	}
	if gc.Components.Monitor != nil {
		if gc.Components.Monitor.Passthrough != nil {
			paths = append(paths, &gc.Components.Monitor.Passthrough.CAFile)
		}
		if gc.Components.Monitor.FrontProxy != nil {
			paths = append(paths, &gc.Components.Monitor.FrontProxy.CAFile)
			paths = append(paths, &gc.Components.Monitor.FrontProxy.ClientCertFile)
			paths = append(paths, &gc.Components.Monitor.FrontProxy.ClientKeyFile)
		}
	}
	if gc.Components.Auth != nil {
		if gc.Components.Auth.Passthrough != nil {
			paths = append(paths, &gc.Components.Auth.Passthrough.CAFile)
		}
		if gc.Components.Auth.FrontProxy != nil {
			paths = append(paths, &gc.Components.Auth.FrontProxy.CAFile)
			paths = append(paths, &gc.Components.Auth.FrontProxy.ClientCertFile)
			paths = append(paths, &gc.Components.Auth.FrontProxy.ClientKeyFile)
		}
	}
	if gc.Components.Registry != nil {
		if gc.Components.Registry.Passthrough != nil {
			paths = append(paths, &gc.Components.Registry.Passthrough.CAFile)
		}
		if gc.Components.Registry.FrontProxy != nil {
			paths = append(paths, &gc.Components.Registry.FrontProxy.CAFile)
			paths = append(paths, &gc.Components.Registry.FrontProxy.ClientCertFile)
			paths = append(paths, &gc.Components.Registry.FrontProxy.ClientKeyFile)
		}
	}
	return paths
}
