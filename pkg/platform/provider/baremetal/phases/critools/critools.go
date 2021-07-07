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
 */

package critools

import (
	"bytes"
	"fmt"
	"path"

	"tkestack.io/tke/pkg/util/template"

	"github.com/pkg/errors"
	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/platform/provider/baremetal/res"
	"tkestack.io/tke/pkg/util/ssh"
)

type Option struct {
	RuntimeEndPoint string
	ImageEndpoint   string
	TimeOut         string
}

const (
	containerdConfigFile = "/etc/crictl.yaml"
)

func Install(s ssh.Interface, option *Option) error {
	dstFile, err := res.CriTools.CopyToNodeWithDefault(s)
	if err != nil {
		return err
	}

	cmd := "tar xvaf %s -C %s --strip-components=1"
	_, stderr, exit, err := s.Execf(cmd, dstFile, constants.DstBinDir)
	if err != nil {
		return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", cmd, exit, stderr, err)
	}

	data, err := template.ParseFile(path.Join(constants.ConfDir, "cri-tools/crictl.yaml"), option)
	if err != nil {
		return err
	}
	err = s.WriteFile(bytes.NewReader(data), containerdConfigFile)
	if err != nil {
		return errors.Wrapf(err, "write %s error", containerdConfigFile)
	}

	return nil
}
