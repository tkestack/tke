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

package gateway

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"regexp"

	"golang.org/x/oauth2"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"tkestack.io/tke/pkg/apiserver/authentication/authenticator/oidc"
	"tkestack.io/tke/pkg/gateway/api"
	gatewayconfig "tkestack.io/tke/pkg/gateway/apis/config"
	"tkestack.io/tke/pkg/gateway/assets"
	"tkestack.io/tke/pkg/gateway/proxy"
	"tkestack.io/tke/pkg/gateway/websocket"
	"tkestack.io/tke/pkg/gateway/webtty"

	"html/template"
)

const (
	DefaultTitle   = "TKEStack"
	DefaultLogoDir = "default"
)

// ExtraConfig contains the additional configuration of apiserver.
type ExtraConfig struct {
	ServerName        string
	OAuthConfig       *oauth2.Config
	OIDCHttpClient    *http.Client
	OIDCAuthenticator *oidc.Authenticator
	GatewayConfig     *gatewayconfig.GatewayConfiguration
	HeaderRequest     bool
}

// Config contains the core configuration instance of server and additional
// configuration.
type Config struct {
	GenericConfig *genericapiserver.RecommendedConfig
	ExtraConfig   ExtraConfig
}

type completedConfig struct {
	GenericConfig genericapiserver.CompletedConfig
	ExtraConfig   *ExtraConfig
}

// CompletedConfig embed a private pointer of Config.
type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

// Gateway contains state for TKE gateway server.
type Gateway struct {
	GenericAPIServer *genericapiserver.GenericAPIServer
}

// Complete fills in any fields not set that are required to have valid data.
// It's mutating the receiver.
func (cfg *Config) Complete() CompletedConfig {
	c := completedConfig{
		cfg.GenericConfig.Complete(),
		&cfg.ExtraConfig,
	}

	return CompletedConfig{&c}
}

// New returns a new instance of APIServer from the given config.
func (c completedConfig) New(delegationTarget genericapiserver.DelegationTarget) (*Gateway, error) {
	s, err := c.GenericConfig.New(c.ExtraConfig.ServerName, delegationTarget)
	if err != nil {
		return nil, err
	}

	registerCallbackRoute(s.Handler.NonGoRestfulMux, c.ExtraConfig.OAuthConfig, c.ExtraConfig.OIDCHttpClient, c.ExtraConfig.GatewayConfig.DisableOIDCProxy)

	if !c.ExtraConfig.GatewayConfig.DisableOIDCProxy {
		if err := registerAuthRoute(s.Handler.NonGoRestfulMux, c.ExtraConfig.OIDCHttpClient, c.ExtraConfig.OIDCAuthenticator); err != nil {
			return nil, err
		}
	}
	consoleConfig := new(gatewayconfig.ConsoleConfig)

	if c.ExtraConfig.GatewayConfig.ConsoleConfig != nil {
		consoleConfig = c.ExtraConfig.GatewayConfig.ConsoleConfig
	} else {
		consoleConfig.Title = DefaultTitle
		consoleConfig.LogoDir = DefaultLogoDir

	}

	files, err := ioutil.ReadDir(assets.RootDir)
	if err != nil {
		return nil, err
	}

	sourceRe := regexp.MustCompile(`\.tmpl\.html$`)
	targetRe := regexp.MustCompile(`\.tmpl`)

	for _, file := range files {
		if !sourceRe.MatchString(file.Name()) {
			continue
		}
		var buf bytes.Buffer
		t, err := template.New(file.Name()).Delims("{%", "%}").ParseFiles(assets.RootDir + file.Name())
		if err != nil {
			return nil, err
		}
		if err = t.Execute(&buf, consoleConfig); err != nil {
			return nil, err
		}
		// // remove .tmpl in file name
		targetFileName := targetRe.ReplaceAllString(file.Name(), "")
		if err = ioutil.WriteFile(assets.RootDir+targetFileName, buf.Bytes(), 0644); err != nil {
			return nil, err
		}
	}

	if err := proxy.RegisterRoute(s.Handler.NonGoRestfulMux, c.ExtraConfig.GatewayConfig, c.ExtraConfig.OIDCAuthenticator); err != nil {
		return nil, err
	}

	if err := api.RegisterRoute(s.Handler.GoRestfulContainer, c.ExtraConfig.GatewayConfig, c.ExtraConfig.OAuthConfig, c.ExtraConfig.OIDCHttpClient, c.ExtraConfig.OIDCAuthenticator, c.ExtraConfig.HeaderRequest); err != nil {
		return nil, err
	}

	if err := webtty.RegisterRoute(s.Handler.NonGoRestfulMux, c.ExtraConfig.GatewayConfig); err != nil {
		return nil, err
	}
	if err := websocket.RegisterRoute(s.Handler.NonGoRestfulMux, c.ExtraConfig.GatewayConfig); err != nil {
		return nil, err
	}

	if !c.ExtraConfig.HeaderRequest {
		assets.RegisterRoute(s.Handler.NonGoRestfulMux, c.ExtraConfig.OAuthConfig, c.ExtraConfig.GatewayConfig.DisableOIDCProxy)
	} else {
		assets.RegisterRoute(s.Handler.NonGoRestfulMux, nil, c.ExtraConfig.GatewayConfig.DisableOIDCProxy)
	}

	m := &Gateway{
		GenericAPIServer: s,
	}

	return m, nil
}
