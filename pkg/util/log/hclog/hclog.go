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

package hclog

import (
	"fmt"
	hcLogLib "github.com/hashicorp/go-hclog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"log"
)

// HCLogger implement the go-hclog interface.
type HCLogger struct {
	logger *zap.Logger
}

// NewHCLogger creates the HCLogger object by given logger.
func NewHCLogger(logger *zap.Logger) *HCLogger {
	return &HCLogger{
		logger: logger,
	}
}

// Trace emit a message and key/value pairs at the Trace level.
func (h *HCLogger) Trace(msg string, args ...interface{}) {
	h.logger.Debug(fmt.Sprintf(msg, args...))
}

// Debug emit a message and key/value pairs at the Debug level.
func (h *HCLogger) Debug(msg string, args ...interface{}) {
	h.logger.Debug(fmt.Sprintf(msg, args...))
}

// Info emit a message and key/value pairs at the Info level.
func (h *HCLogger) Info(msg string, args ...interface{}) {
	h.logger.Info(fmt.Sprintf(msg, args...))
}

// Warn emit a message and key/value pairs at the Warn level.
func (h *HCLogger) Warn(msg string, args ...interface{}) {
	h.logger.Warn(fmt.Sprintf(msg, args...))
}

// Error emit a message and key/value pairs at the Error level.
func (h *HCLogger) Error(msg string, args ...interface{}) {
	h.logger.Error(fmt.Sprintf(msg, args...))
}

// IsTrace indicate if TRACE logs would be emitted.
func (h *HCLogger) IsTrace() bool {
	entry := h.logger.Check(zap.DebugLevel, "")
	return entry != nil
}

// IsDebug indicate if DEBUG logs would be emitted.
func (h *HCLogger) IsDebug() bool {
	entry := h.logger.Check(zap.DebugLevel, "")
	return entry != nil
}

// IsInfo indicate if INFO logs would be emitted.
func (h *HCLogger) IsInfo() bool {
	entry := h.logger.Check(zap.InfoLevel, "")
	return entry != nil
}

// IsWarn indicate if WARN logs would be emitted.
func (h *HCLogger) IsWarn() bool {
	entry := h.logger.Check(zap.WarnLevel, "")
	return entry != nil
}

// IsError indicate if ERROR logs would be emitted.
func (h *HCLogger) IsError() bool {
	entry := h.logger.Check(zap.ErrorLevel, "")
	return entry != nil
}

// With creates a sub logger that will always have the given key/value pairs
func (h *HCLogger) With(args ...interface{}) hcLogLib.Logger {
	if len(args) != 2 {
		return h
	}
	return NewHCLogger(h.logger.With(zap.Any(args[0].(string), args[1])))
}

// Named create a logger that will prepend the name string on the front of all
// messages.
// If the logger already has a name, the new value will be appended to the current
// name. That way, a major subsystem can use this to decorate all it's own logs
// without losing context.
func (h *HCLogger) Named(name string) hcLogLib.Logger {
	return NewHCLogger(h.logger.Named(name))
}

// ResetNamed create a logger that will prepend the name string on the front of
// all messages.
// This sets the name of the logger to the value directly, unlike Named which honor
// the current name as well.
func (h *HCLogger) ResetNamed(name string) hcLogLib.Logger {
	return h
}

// SetLevel updates the level. This should affect all sub-loggers as well. If an
// implementation cannot update the level on the fly, it should no-op.
func (h *HCLogger) SetLevel(level hcLogLib.Level) {

}

// StandardLogger return a value that conforms to the stdlib log.Logger interface
func (h *HCLogger) StandardLogger(opts *hcLogLib.StandardLoggerOptions) *log.Logger {
	l, _ := zap.NewStdLogAt(h.logger, zapcore.InfoLevel)
	return l
}

// StandardWriter return a value that conforms to io.Writer, which can be passed into log.SetOutput()
func (h *HCLogger) StandardWriter(opts *hcLogLib.StandardLoggerOptions) io.Writer {
	return h
}

// Write writes len(p) bytes from p to the underlying data stream.
func (h *HCLogger) Write(p []byte) (n int, err error) {
	h.Info(string(p))
	return len(p), nil
}
