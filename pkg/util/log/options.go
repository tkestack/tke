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
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

const (
	// FlagLevel means the log level.
	FlagLevel            = "log-level"
	flagFormat           = "log-format"
	flagDisableColor     = "log-disable-color"
	flagEnableCaller     = "log-enable-caller"
	flagOutputPaths      = "log-output-paths"
	flagErrorOutputPaths = "log-error-output-paths"
)

const (
	configLevel            = "log.level"
	configFormat           = "log.format"
	configDisableColor     = "log.disable_color"
	configEnableCaller     = "log.enable_caller"
	configOutputPaths      = "log.output_paths"
	configErrorOutputPaths = "log.error_output_paths"
)

// Options contains configuration items related to log.
type Options struct {
	Level            zapcore.Level
	Format           Format
	DisableColor     bool
	EnableCaller     bool
	OutputPaths      []string
	ErrorOutputPaths []string
}

// NewOptions creates a Options object with default parameters.
func NewOptions() *Options {
	return &Options{
		Level:            zapcore.InfoLevel,
		Format:           consoleFormat,
		DisableColor:     false,
		EnableCaller:     false,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

// AddFlags adds flags for log to the specified FlagSet object.
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.String(FlagLevel, o.Level.String(), "Minimum log output `LEVEL`.")
	_ = viper.BindPFlag(configLevel, fs.Lookup(FlagLevel))

	fs.String(flagFormat, o.Format.String(), "Log output `FORMAT`, support plain or json format.")
	_ = viper.BindPFlag(configFormat, fs.Lookup(flagFormat))

	fs.Bool(flagDisableColor, o.DisableColor, "Disable output ansi colors in plain format logs.")
	_ = viper.BindPFlag(configDisableColor, fs.Lookup(flagDisableColor))

	fs.Bool(flagEnableCaller, o.EnableCaller, "Enable output of caller information in the log.")
	_ = viper.BindPFlag(configEnableCaller, fs.Lookup(flagEnableCaller))

	fs.StringSlice(flagOutputPaths, o.OutputPaths, "Output paths of log")
	_ = viper.BindPFlag(configOutputPaths, fs.Lookup(flagOutputPaths))

	fs.StringSlice(flagErrorOutputPaths, o.ErrorOutputPaths, "Error output paths of log")
	_ = viper.BindPFlag(configErrorOutputPaths, fs.Lookup(flagErrorOutputPaths))
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *Options) ApplyFlags() []error {
	var errs []error

	level := viper.GetString(configLevel)
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		errs = append(errs, err)
	} else {
		o.Level = zapLevel
	}

	if format, err := ParseFormat(viper.GetString(configFormat)); err != nil {
		errs = append(errs, err)
	} else {
		o.Format = format
	}

	o.DisableColor = viper.GetBool(configDisableColor)
	o.EnableCaller = viper.GetBool(configEnableCaller)
	o.OutputPaths = viper.GetStringSlice(configOutputPaths)
	o.ErrorOutputPaths = viper.GetStringSlice(configErrorOutputPaths)

	return errs
}
