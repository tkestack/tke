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

package cniplugins

import (
	"fmt"

	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/platform/provider/baremetal/res"
	"tkestack.io/tke/pkg/util/ssh"
)

type Option struct {
}

func Install(s ssh.Interface, option *Option) error {
	dstFile, err := res.CNIPlugins.CopyToNodeWithDefault(s)
	if err != nil {
		return err
	}

	_, stderr, exit, err := s.Execf("[ -d %s ] || mkdir -p %s", constants.CNIBinDir, constants.CNIBinDir)
	if exit != 0 || err != nil {
		return fmt.Errorf("clean %s failed:exit %d:stderr %s:error %s", constants.EtcdDataDir, exit, stderr, err)
	}

	cmd := "tar xvaf %s -C %s"
	_, stderr, exit, err = s.Execf(cmd, dstFile, constants.CNIBinDir)
	if err != nil {
		return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", cmd, exit, stderr, err)
	}

	return nil
}
