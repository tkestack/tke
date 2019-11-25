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

package logger

import (
	"fmt"

	"tkestack.io/tke/pkg/util/log"
)

// WrapLogger is the implementation for a Logger.
type WrapLogger struct {
	enable bool
}

// EnableLog controls whether print the message.
func (l *WrapLogger) EnableLog(enable bool) {
	l.enable = enable
}

// IsEnabled returns if logger is enabled.
func (l *WrapLogger) IsEnabled() bool {
	return l.enable
}

// Print formats using the default formats for its operands and logs the message.
func (l *WrapLogger) Print(v ...interface{}) {
	if l.enable {
		log.Info(fmt.Sprint(v...))
	}
}

// Printf formats according to a format specifier and logs the message.
func (l *WrapLogger) Printf(format string, v ...interface{}) {
	if l.enable {
		log.Infof(format, v...)
	}
}
