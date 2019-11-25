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

package glog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

var (
	logger *zap.Logger
	mu     sync.RWMutex
)

// Init to initialize global logger by given zap.Logger object.
func Init(l *zap.Logger) {
	mu.Lock()
	defer mu.Unlock()
	logger = l
}

// Level is a shim
type Level int32

// Verbose is a shim
type Verbose bool

// Flush is a shim
func Flush() {
	if logger != nil {
		_ = logger.Sync()
	}
}

// V is a shim
func V(level Level) Verbose {
	var lvl zapcore.Level
	if level > 9 {
		lvl = zapcore.InfoLevel
	} else {
		lvl = zapcore.DebugLevel
	}
	if logger == nil {
		return false
	}
	checkEntry := logger.Check(lvl, "")
	return checkEntry != nil
}

// Set is a shim
func (l Level) Set(_ string) error {
	return nil
}

// Info is equivalent to the global Info function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbose) Info(args ...interface{}) {
	if v {
		if logger != nil {
			logger.Info(fmt.Sprint(args...))
		}
	}
}

// Infoln is equivalent to the global Infoln function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbose) Infoln(args ...interface{}) {
	if v {
		if logger != nil {
			logger.Info(fmt.Sprintf("%s\n", fmt.Sprint(args...)))
		}
	}
}

// Infof is equivalent to the global Infof function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbose) Infof(format string, args ...interface{}) {
	if v {
		if logger != nil {
			logger.Info(fmt.Sprintf(format, args...))
		}
	}
}

// Info logs to the INFO log.
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Info(args ...interface{}) {
	if logger != nil {
		logger.Info(fmt.Sprint(args...))
	}
}

// InfoDepth acts as Info but uses depth to determine which call frame to log.
// InfoDepth(0, "msg") is the same as Info("msg").
func InfoDepth(_ int, args ...interface{}) {
	if logger != nil {
		logger.Info(fmt.Sprint(args...))
	}
}

// Infoln logs to the INFO log.
// Arguments are handled in the manner of fmt.Println; a newline is appended if missing.
func Infoln(args ...interface{}) {
	if logger != nil {
		logger.Info(fmt.Sprintf("%s\n", fmt.Sprint(args...)))
	}
}

// Infof logs to the INFO log.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Infof(format string, args ...interface{}) {
	if logger != nil {
		logger.Info(fmt.Sprintf(format, args...))
	}
}

// Warning logs to the WARNING and INFO logs.
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Warning(args ...interface{}) {
	if logger != nil {
		logger.Warn(fmt.Sprint(args...))
	}
}

// WarningDepth acts as Warning but uses depth to determine which call frame to log.
// WarningDepth(0, "msg") is the same as Warning("msg").
func WarningDepth(_ int, args ...interface{}) {
	if logger != nil {
		logger.Warn(fmt.Sprint(args...))
	}
}

// Warningln logs to the WARNING and INFO logs.
// Arguments are handled in the manner of fmt.Println; a newline is appended if missing.
func Warningln(args ...interface{}) {
	if logger != nil {
		logger.Warn(fmt.Sprintf("%s\n", fmt.Sprint(args...)))
	}
}

// Warningf logs to the WARNING and INFO logs.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Warningf(format string, args ...interface{}) {
	if logger != nil {
		logger.Warn(fmt.Sprintf(format, args...))
	}
}

// Error logs to the ERROR, WARNING, and INFO logs.
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Error(args ...interface{}) {
	if logger != nil {
		logger.Error(fmt.Sprint(args...))
	}
}

// ErrorDepth acts as Error but uses depth to determine which call frame to log.
// ErrorDepth(0, "msg") is the same as Error("msg").
func ErrorDepth(_ int, args ...interface{}) {
	if logger != nil {
		logger.Error(fmt.Sprint(args...))
	}
}

// Errorln logs to the ERROR, WARNING, and INFO logs.
// Arguments are handled in the manner of fmt.Println; a newline is appended if missing.
func Errorln(args ...interface{}) {
	if logger != nil {
		logger.Error(fmt.Sprintf("%s\n", fmt.Sprint(args...)))
	}
}

// Errorf logs to the ERROR, WARNING, and INFO logs.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Errorf(format string, args ...interface{}) {
	if logger != nil {
		logger.Error(fmt.Sprintf(format, args...))
	}
}

// Fatal logs to the FATAL, ERROR, WARNING, and INFO logs,
// including a stack trace of all running goroutines, then calls os.Exit(255).
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Fatal(args ...interface{}) {
	if logger != nil {
		logger.Error(fmt.Sprint(args...))
	}
	os.Exit(255)
}

// FatalDepth acts as Fatal but uses depth to determine which call frame to log.
// FatalDepth(0, "msg") is the same as Fatal("msg").
func FatalDepth(_ int, args ...interface{}) {
	if logger != nil {
		logger.Error(fmt.Sprint(args...))
	}
	os.Exit(255)
}

// Fatalln logs to the FATAL, ERROR, WARNING, and INFO logs,
// including a stack trace of all running goroutines, then calls os.Exit(255).
// Arguments are handled in the manner of fmt.Println; a newline is appended if missing.
func Fatalln(args ...interface{}) {
	if logger != nil {
		logger.Error(fmt.Sprintf("%s\n", fmt.Sprint(args...)))
	}
	os.Exit(255)
}

// Fatalf logs to the FATAL, ERROR, WARNING, and INFO logs,
// including a stack trace of all running goroutines, then calls os.Exit(255).
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Fatalf(format string, args ...interface{}) {
	if logger != nil {
		logger.Error(fmt.Sprintf(format, args...))
	}
	os.Exit(255)
}

// Exit logs to the FATAL, ERROR, WARNING, and INFO logs, then calls os.Exit(1).
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Exit(args ...interface{}) {
	if logger != nil {
		logger.Error(fmt.Sprint(args...))
	}
	os.Exit(1)
}

// ExitDepth acts as Exit but uses depth to determine which call frame to log.
// ExitDepth(0, "msg") is the same as Exit("msg").
func ExitDepth(_ int, args ...interface{}) {
	if logger != nil {
		logger.Error(fmt.Sprint(args...))
	}
	os.Exit(1)
}

// Exitln logs to the FATAL, ERROR, WARNING, and INFO logs, then calls os.Exit(1).
func Exitln(args ...interface{}) {
	if logger != nil {
		logger.Error(fmt.Sprintf("%s\n", fmt.Sprint(args...)))
	}
	os.Exit(1)
}

// Exitf logs to the FATAL, ERROR, WARNING, and INFO logs, then calls os.Exit(1).
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Exitf(format string, args ...interface{}) {
	if logger != nil {
		logger.Error(fmt.Sprintf(format, args...))
	}
	os.Exit(1)
}
