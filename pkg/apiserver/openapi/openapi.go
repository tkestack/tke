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

package openapi

import (
	"fmt"
	"github.com/go-openapi/spec"
	"k8s.io/apimachinery/pkg/version"
	openapinamer "k8s.io/apiserver/pkg/endpoints/openapi"
	genericapiserver "k8s.io/apiserver/pkg/server"
	openapicommon "k8s.io/kube-openapi/pkg/common"
	"os"
	"strings"
	"tkestack.io/tke/api/platform"
	appversion "tkestack.io/tke/pkg/app/version"
)

// SetupOpenAPI to setup the generic api server openapi configuration.
func SetupOpenAPI(genericAPIServerConfig *genericapiserver.Config, getDefinitions openapicommon.GetOpenAPIDefinitions, title string, license string, host string, port int) {
	appVersion := appversion.Get()

	// openAPI
	genericAPIServerConfig.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(getDefinitions, openapinamer.NewDefinitionNamer(platform.Scheme))
	genericAPIServerConfig.OpenAPIConfig.Info.Title = title
	genericAPIServerConfig.OpenAPIConfig.Info.License = &spec.License{Name: license}
	genericAPIServerConfig.OpenAPIConfig.Info.Version = appVersion.GitVersion
	genericAPIServerConfig.OpenAPIConfig.PostProcessSpec = postProcessOpenAPISpec(host, port)

	// version
	genericAPIServerConfig.Version = &version.Info{
		GitVersion: appVersion.GitVersion,
		BuildDate:  appVersion.BuildDate,
		GoVersion:  appVersion.GoVersion,
		Compiler:   appVersion.Compiler,
		Platform:   appVersion.Platform,
	}
}

func postProcessOpenAPISpec(host string, port int) func(*spec.Swagger) (*spec.Swagger, error) {
	return func(swagger *spec.Swagger) (*spec.Swagger, error) {
		swagger.Schemes = []string{"https"}
		debug := os.Getenv("TKE_DEBUG")
		if debug != "" && strings.ToLower(debug) == "true" {
			swagger.Host = fmt.Sprintf("localhost:%d", port)
		} else {
			swagger.Host = fmt.Sprintf("%s:%d", host, port)
		}
		for _, path := range swagger.Paths.Paths {
			if path.Get != nil {
				path.Get.Summary = path.Get.Description
			}
			if path.Delete != nil {
				path.Delete.Summary = path.Delete.Description
			}
			if path.Options != nil {
				path.Options.Summary = path.Options.Description
			}
			if path.Head != nil {
				path.Head.Summary = path.Head.Description
			}
			if path.Patch != nil {
				path.Patch.Summary = path.Patch.Description
			}
			if path.Post != nil {
				path.Post.Summary = path.Post.Description
			}
			if path.Put != nil {
				path.Put.Summary = path.Put.Description
			}
		}
		return swagger, nil
	}
}
