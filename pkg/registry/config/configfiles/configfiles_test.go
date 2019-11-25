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
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"path/filepath"
	"testing"
	registryconfig "tkestack.io/tke/pkg/registry/apis/config"
	registryscheme "tkestack.io/tke/pkg/registry/apis/config/scheme"
	registryconfigv1 "tkestack.io/tke/pkg/registry/apis/config/v1"
	utilfiles "tkestack.io/tke/pkg/util/files"
	utilfs "tkestack.io/tke/pkg/util/filesystem"
	utiltest "tkestack.io/tke/pkg/util/test"
)

const configDir = "/test-config-dir"
const registryFile = "registry"

func TestLoad(t *testing.T) {
	cases := []struct {
		desc   string
		file   *string
		expect *registryconfig.RegistryConfiguration
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
			newString(`{"apiVersion":"registry.config.tkestack.io/v1"}`),
			nil,
			"failed to decode",
		},
		{
			"missing version",
			newString(`{"kind":"RegistryConfiguration"}`),
			nil,
			"failed to decode",
		},
		{
			"unregistered kind",
			newString(`{"kind":"BogusKind","apiVersion":"registry.config.tkestack.io/v1"}`),
			nil,
			"failed to decode",
		},
		{
			"unregistered version",
			newString(`{"kind":"RegistryConfiguration","apiVersion":"bogusversion"}`),
			nil,
			"failed to decode",
		},

		// empty object with correct kind and version should result in the defaults for that kind and version
		{
			"default from yaml",
			newString(`kind: RegistryConfiguration
apiVersion: registry.config.tkestack.io/v1`),
			newConfig(t),
			"",
		},
		{
			"default from json",
			newString(`{"kind":"RegistryConfiguration","apiVersion":"registry.config.tkestack.io/v1"}`),
			newConfig(t),
			"",
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			fs := utilfs.NewFakeFs()
			path := filepath.Join(configDir, registryFile)
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

func newString(s string) *string {
	return &s
}

func addFile(fs utilfs.Filesystem, path string, file string) error {
	if err := utilfiles.EnsureDir(fs, filepath.Dir(path)); err != nil {
		return err
	}
	return utilfiles.ReplaceFile(fs, path, []byte(file))
}

func newConfig(t *testing.T) *registryconfig.RegistryConfiguration {
	registryScheme, _, err := registryscheme.NewSchemeAndCodecs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// get the built-in default configuration
	external := &registryconfigv1.RegistryConfiguration{}
	registryScheme.Default(external)
	kc := &registryconfig.RegistryConfiguration{}
	err = registryScheme.Convert(external, kc, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return kc
}
