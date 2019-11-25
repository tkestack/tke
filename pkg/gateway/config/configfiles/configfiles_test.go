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

package configfiles

import (
	"fmt"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"path/filepath"
	"testing"
	gatewayconfig "tkestack.io/tke/pkg/gateway/apis/config"
	gatewayscheme "tkestack.io/tke/pkg/gateway/apis/config/scheme"
	gatewayconfigv1 "tkestack.io/tke/pkg/gateway/apis/config/v1"
	utilfiles "tkestack.io/tke/pkg/util/files"
	utilfs "tkestack.io/tke/pkg/util/filesystem"
	utiltest "tkestack.io/tke/pkg/util/test"
)

const configDir = "/test-config-dir"
const relativePath = "relative/path/test"
const gatewayFile = "gateway"

func TestLoad(t *testing.T) {
	cases := []struct {
		desc   string
		file   *string
		expect *gatewayconfig.GatewayConfiguration
		err    string
	}{
		// missing file
		{
			"missing file",
			nil,
			nil,
			"failed to read",
		},
		// empty file
		{
			"empty file",
			newString(``),
			nil,
			"was empty",
		},
		// invalid format
		{
			"invalid yaml",
			newString(`*`),
			nil,
			"failed to decode",
		},
		{
			"invalid json",
			newString(`{*`),
			nil,
			"failed to decode",
		},
		// invalid object
		{
			"missing kind",
			newString(`{"apiVersion":"gateway.config.tkestack.io/v1"}`),
			nil,
			"failed to decode",
		},
		{
			"missing version",
			newString(`{"kind":"GatewayConfiguration"}`),
			nil,
			"failed to decode",
		},
		{
			"unregistered kind",
			newString(`{"kind":"BogusKind","apiVersion":"gateway.config.tkestack.io/v1"}`),
			nil,
			"failed to decode",
		},
		{
			"unregistered version",
			newString(`{"kind":"GatewayConfiguration","apiVersion":"bogusversion"}`),
			nil,
			"failed to decode",
		},

		// empty object with correct kind and version should result in the defaults for that kind and version
		{
			"default from yaml",
			newString(`kind: GatewayConfiguration
apiVersion: gateway.config.tkestack.io/v1`),
			newConfig(t),
			"",
		},
		{
			"default from json",
			newString(`{"kind":"GatewayConfiguration","apiVersion":"gateway.config.tkestack.io/v1"}`),
			newConfig(t),
			"",
		},

		// relative path
		{
			"yaml, relative path is resolved",
			newString(fmt.Sprintf(`kind: GatewayConfiguration
apiVersion: gateway.config.tkestack.io/v1
components:
  platform:
    passthrough:
      caFile: %s`, relativePath)),
			func() *gatewayconfig.GatewayConfiguration {
				gc := newConfig(t)
				gc.Components.Platform = &gatewayconfig.Component{
					Passthrough: &gatewayconfig.PassthroughComponent{
						CAFile: filepath.Join(configDir, relativePath),
					},
				}
				return gc
			}(),
			"",
		},
		{
			"json, relative path is resolved",
			newString(fmt.Sprintf(`{
  "kind":"GatewayConfiguration",
  "apiVersion":"gateway.config.tkestack.io/v1",
  "components": {
    "platform": {
      "passthrough": {
        "caFile": "%s"
      }
    }
  }
}`, relativePath)),
			func() *gatewayconfig.GatewayConfiguration {
				gc := newConfig(t)
				gc.Components.Platform = &gatewayconfig.Component{
					Passthrough: &gatewayconfig.PassthroughComponent{
						CAFile: filepath.Join(configDir, relativePath),
					},
				}
				return gc
			}(),
			"",
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			fs := utilfs.NewFakeFs()
			path := filepath.Join(configDir, gatewayFile)
			if c.file != nil {
				if err := addFile(fs, path, *c.file); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
			loader, err := NewFsLoader(fs, path)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			kc, err := loader.Load()
			if utiltest.SkipRest(t, c.desc, err, c.err) {
				return
			}
			if !apiequality.Semantic.DeepEqual(c.expect, kc) {
				t.Fatalf("expect %#v but got %#v", *c.expect, *kc)
			}
		})
	}
}

func TestResolveRelativePaths(t *testing.T) {
	absolutePath := filepath.Join(configDir, "absolute")
	cases := []struct {
		desc   string
		path   string
		expect string
	}{
		{"empty path", "", ""},
		{"absolute path", absolutePath, absolutePath},
		{"relative path", relativePath, filepath.Join(configDir, relativePath)},
	}

	gc := newConfig(t)
	gc.Components.Platform = &gatewayconfig.Component{
		Passthrough: &gatewayconfig.PassthroughComponent{},
	}
	paths := gatewayconfig.GatewayConfigurationPathRefs(gc)
	if len(paths) == 0 {
		t.Fatalf("requires at least one path field to exist in the GatewayConfiguration type")
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			// set the path, resolve it, and check if it resolved as we would expect
			*(paths[0]) = c.path
			resolveRelativePaths(paths, configDir)
			if *(paths[0]) != c.expect {
				t.Fatalf("expect %s but got %s", c.expect, *(paths[0]))
			}
		})
	}
}

func newString(s string) *string {
	return &s
}

func addFile(fs utilfs.Filesystem, path string, file string) error {
	if err := utilfiles.EnsureDir(fs, filepath.Dir(path)); err != nil {
		return err
	}
	return utilfiles.ReplaceFile(fs, path, []byte(file))
}

func newConfig(t *testing.T) *gatewayconfig.GatewayConfiguration {
	gatewayScheme, _, err := gatewayscheme.NewSchemeAndCodecs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// get the built-in default configuration
	external := &gatewayconfigv1.GatewayConfiguration{}
	gatewayScheme.Default(external)
	kc := &gatewayconfig.GatewayConfiguration{}
	err = gatewayScheme.Convert(external, kc, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return kc
}
