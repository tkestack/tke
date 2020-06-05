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

package keepalived

import (
	"bytes"
	"fmt"

	//"strings"

	"github.com/pkg/errors"
	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/platform/provider/baremetal/images"
	"tkestack.io/tke/pkg/util/ssh"
	"tkestack.io/tke/pkg/util/template"
)

type Option struct {
	IP  string
	VIP string
}

func Install(s ssh.Interface, option *Option) error {
	networkInterface := ssh.GetNetworkInterface(s, option.IP)
	if networkInterface == "" {
		return fmt.Errorf("can't get network interface by %s", option.IP)
	}

	data, err := template.ParseFile(constants.ManifestsDir+"keepalived/keepalived.conf", map[string]interface{}{
		"Interface": networkInterface,
		"VIP":       option.VIP,
	})
	if err != nil {
		return errors.Wrap(err, option.IP)
	}
	err = s.WriteFile(bytes.NewReader(data), constants.KeepavliedConfigFile)
	if err != nil {
		return errors.Wrap(err, option.IP)
	}

	data, err = template.ParseFile(constants.ManifestsDir+"keepalived/keepalived.yaml", map[string]interface{}{
		"Image": images.Get().Keepalived.FullName(),
	})
	if err != nil {
		return errors.Wrap(err, option.IP)
	}
	err = s.WriteFile(bytes.NewReader(data), constants.KeepavlivedManifestFile)
	if err != nil {
		return errors.Wrap(err, option.IP)
	}

	return nil
}
