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

package logrus

import (
	"context"
	"go.uber.org/zap"
	"sync"
	"time"
)

var (
	logger *Logger
	lock   sync.RWMutex
)

func Init(zapLogger *zap.Logger) {
	lock.Lock()
	defer lock.Unlock()
	logger = NewLogger(zapLogger)
}

func StandardLogger() *Logger {
	return logger
}

// WithError creates an entry from the standard logger and adds an error to it, using the value defined in ErrorKey as key.
func WithError(err error) *Entry {
	if logger != nil {
		return logger.WithError(err)
	}
	return nil
}

// WithContext creates an entry from the standard logger and adds a context to it.
func WithContext(ctx context.Context) *Entry {
	if logger != nil {
		return logger.WithContext(ctx)
	}
	return nil
}

// WithField creates an entry from the standard logger and adds a field to
// it. If you want multiple fields, use `WithFields`.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithField(key string, value interface{}) *Entry {
	if logger != nil {
		return logger.WithField(key, value)
	}
	return nil
}

// WithFields creates an entry from the standard logger and adds multiple
// fields to it. This is simply a helper for `WithField`, invoking it
// once for each field.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithFields(fields Fields) *Entry {
	if logger != nil {
		return logger.WithFields(fields)
	}
	return nil
}

// WithTime creats an entry from the standard logger and overrides the time of
// logs generated with it.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithTime(t time.Time) *Entry {
	if logger != nil {
		return logger.WithTime(t)
	}
	return nil
}

// Trace logs a message at level Trace on the standard logger.
func Trace(args ...interface{}) {
	if logger != nil {
		logger.Trace(args)
	}
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	if logger != nil {
		logger.Debug(args)
	}
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {
	if logger != nil {
		logger.Print(args)
	}
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	if logger != nil {
		logger.Info(args)
	}
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	if logger != nil {
		logger.Warn(args)
	}
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	if logger != nil {
		logger.Warning(args)
	}
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	if logger != nil {
		logger.Error(args)
	}
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	if logger != nil {
		logger.Panic(args)
	}
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(args ...interface{}) {
	if logger != nil {
		logger.Fatal(args)
	}
}

// Tracef logs a message at level Trace on the standard logger.
func Tracef(format string, args ...interface{}) {
	if logger != nil {
		logger.Tracef(format, args)
	}
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	if logger != nil {
		logger.Debugf(format, args)
	}
}

// Printf logs a message at level Info on the standard logger.
func Printf(format string, args ...interface{}) {
	if logger != nil {
		logger.Printf(format, args)
	}
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	if logger != nil {
		logger.Infof(format, args)
	}
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	if logger != nil {
		logger.Warnf(format, args)
	}
}

// Warningf logs a message at level Warn on the standard logger.
func Warningf(format string, args ...interface{}) {
	if logger != nil {
		logger.Warningf(format, args)
	}
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	if logger != nil {
		logger.Errorf(format, args)
	}
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(format string, args ...interface{}) {
	if logger != nil {
		logger.Panicf(format, args)
	}
}

// Fatalf logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalf(format string, args ...interface{}) {
	if logger != nil {
		logger.Fatalf(format, args)
	}
}

// Traceln logs a message at level Trace on the standard logger.
func Traceln(args ...interface{}) {
	if logger != nil {
		logger.Traceln(args)
	}
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	if logger != nil {
		logger.Debugln(args)
	}
}

// Println logs a message at level Info on the standard logger.
func Println(args ...interface{}) {
	if logger != nil {
		logger.Println(args)
	}
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	if logger != nil {
		logger.Infoln(args)
	}
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) {
	if logger != nil {
		logger.Warnln(args)
	}
}

// Warningln logs a message at level Warn on the standard logger.
func Warningln(args ...interface{}) {
	if logger != nil {
		logger.Warningln(args)
	}
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	if logger != nil {
		logger.Errorln(args)
	}
}

// Panicln logs a message at level Panic on the standard logger.
func Panicln(args ...interface{}) {
	if logger != nil {
		logger.Panicln(args)
	}
}

// Fatalln logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalln(args ...interface{}) {
	if logger != nil {
		logger.Fatalln(args)
	}
}
