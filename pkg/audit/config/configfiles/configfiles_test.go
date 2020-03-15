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
	auditconfig "tkestack.io/tke/pkg/audit/apis/config"
	auditscheme "tkestack.io/tke/pkg/audit/apis/config/scheme"
	auditconfigv1 "tkestack.io/tke/pkg/audit/apis/config/v1"
	utilfiles "tkestack.io/tke/pkg/util/files"
	utilfs "tkestack.io/tke/pkg/util/filesystem"
	utiltest "tkestack.io/tke/pkg/util/test"
)

const configDir = "/test-config-dir"
const auditFile = "audit"

func TestLoad(t *testing.T) {
	cases := []struct {
		desc   string
		file   *string
		expect *auditconfig.AuditConfiguration
		err    string
	}{
		{
			"default from yamlddd",
			newString(`kind: AuditConfiguration
apiVersion: audit.config.tkestack.io/v1
storage:
  elasticSearch:
    address: "123"
`),
			func() *auditconfig.AuditConfiguration {
				dd := newConfig(t)
				dd.Storage.ElasticSearch = &auditconfig.ElasticSearchStorage{
					Address: "123",
				}
				return dd
			}(),
			"",
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			fs := utilfs.NewFakeFs()
			path := filepath.Join(configDir, auditFile)
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

func newConfig(t *testing.T) *auditconfig.AuditConfiguration {
	auditScheme, _, err := auditscheme.NewSchemeAndCodecs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// get the built-in default configuration
	external := &auditconfigv1.AuditConfiguration{}
	auditScheme.Default(external)
	kc := &auditconfig.AuditConfiguration{}
	err = auditScheme.Convert(external, kc, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return kc
}
