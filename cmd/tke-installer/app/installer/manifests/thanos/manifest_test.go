/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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
 *
 */

package thanos

import (
	"encoding/base64"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"tkestack.io/tke/cmd/tke-installer/app/installer/types"
	"tkestack.io/tke/pkg/kubectl"
	"tkestack.io/tke/pkg/util/template"
)

func TestManifest(t *testing.T) {
	matches, err := filepath.Glob("thanos.yaml")
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	if len(matches) == 0 {
		t.Fatalf("matches zero: %v", matches)
	}
	s := `type: S3
config:
  bucket: ""
  endpoint: ""
  region: ""
  access_key: ""
  insecure: false
  signature_version2: false
  secret_key: ""
  put_user_metadata: {}
  http_config:
    idle_conn_timeout: 1m30s
    response_header_timeout: 2m
    insecure_skip_verify: false
    tls_handshake_timeout: 10s
    expect_continue_timeout: 1s
    max_idle_conns: 100
    max_idle_conns_per_host: 100
    max_conns_per_host: 0
  trace:
    enable: false
  list_objects_version: ""
  part_size: 134217728
  sse_config:
    type: ""
    kms_key_id: ""
    kms_encryption_context: {}
    encryption_key: ""
`
	bucketConfig := &types.ThanosBucketConfig{}
	err = yaml.Unmarshal([]byte(s), bucketConfig)
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	thanosYamlBytes, err := yaml.Marshal(bucketConfig)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	thanosYaml := base64.StdEncoding.EncodeToString(thanosYamlBytes)

	params := map[string]interface{}{
		"Image":      "thanos:v0.15.0",
		"ThanosYaml": thanosYaml,
	}
	for _, filename := range matches {
		data, err := template.ParseFile(filename, params)
		if !assert.Nil(t, err) {
			t.FailNow()
		}
		/*err = ioutil.WriteFile("out.yaml", data, os.ModePerm)
		if !assert.Nil(t, err) {
			t.FailNow()
		}*/
		data, err = kubectl.Validate(data)
		if !assert.Nil(t, err) {
			t.Fatal(string(data))
		}
	}
}
