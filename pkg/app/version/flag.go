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

package version

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/pflag"
)

type value int

const (
	boolFalse value = 0
	boolTrue  value = 1
	raw       value = 2
)

const strRawVersion string = "raw"

func (v *value) IsBoolFlag() bool {
	return true
}

func (v *value) Get() interface{} {
	return *v
}

func (v *value) Set(s string) error {
	if s == strRawVersion {
		*v = raw
		return nil
	}
	boolVal, err := strconv.ParseBool(s)
	if boolVal {
		*v = boolTrue
	} else {
		*v = boolFalse
	}
	return err
}

func (v *value) String() string {
	if *v == raw {
		return strRawVersion
	}
	return fmt.Sprintf("%v", bool(*v == boolTrue))
}

// The type of the flag as required by the pflag.value interface
func (v *value) Type() string {
	return "version"
}

const flagName = "version"
const flagShortHand = "V"

var (
	v = boolFalse
)

// AddFlags registers this package's flags on arbitrary FlagSets, such that they
// point to the same value as the global flags.
func AddFlags(fs *pflag.FlagSet) {
	fs.VarP(&v, flagName, flagShortHand, "Print version information and quit.")
	// "--version" will be treated as "--version=true"
	fs.Lookup(flagName).NoOptDefVal = "true"
}

// PrintAndExitIfRequested will check if the -version flag was passed and, if so,
// print the version and exit.
func PrintAndExitIfRequested(appName string) {
	if v == raw {
		fmt.Printf("%s\n", Get())
		os.Exit(0)
	} else if v == boolTrue {
		fmt.Printf("%s %s\n", appName, Get().GitVersion)
		os.Exit(0)
	}
}
