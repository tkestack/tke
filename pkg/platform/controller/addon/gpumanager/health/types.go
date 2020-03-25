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

package health

import (
	v1 "tkestack.io/tke/api/platform/v1"
)

// Prober is used for probing the GPUManager instance status
type Prober interface {
	// Run means starts to run prober
	Run(ch <-chan struct{})
	// Exist tells you whether GPUManager's key is in the prober
	Exist(key string) bool
	// Set means sets a prober for the given GPUManager
	Set(key string, v *v1.GPUManager)
	// DeleteDaemonSet means remove a prober of the given GPUManager's key
	Del(key string)
}
