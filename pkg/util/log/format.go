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

package log

import (
	"fmt"
	"strings"
)

// Format type
type Format string

const (
	consoleFormat Format = "console"
	jsonFormat    Format = "json"
)

// String convert the log format to a string.
func (format Format) String() string {
	if b, err := format.MarshalText(); err == nil {
		return string(b)
	}
	return ""
}

// ParseFormat takes a string format and returns the log format constant.
func ParseFormat(f string) (Format, error) {
	switch strings.ToLower(f) {
	case "console":
		return consoleFormat, nil
	case "json":
		return jsonFormat, nil
	default:
		return "", fmt.Errorf("not a valid log format: %q", f)
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (format *Format) UnmarshalText(text []byte) error {
	l, err := ParseFormat(string(text))
	if err != nil {
		return err
	}

	*format = l

	return nil
}

// MarshalText encodes the log format into UTF-8-encoded text and returns the
// result.
func (format Format) MarshalText() ([]byte, error) {
	switch format {
	case consoleFormat, jsonFormat:
		return []byte(format), nil
	default:
		return nil, fmt.Errorf("not a valid log format")
	}
}
