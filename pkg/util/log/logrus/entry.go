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
	"fmt"
	"go.uber.org/zap"
)

// Entry is the final or intermediate Logrus logging entry.
type Entry struct {
	logger *zap.Logger
	Logger *Logger
	// Level the log entry was logged at: Trace, Debug, Info, Warn, Error, Fatal or Panic
	// This field will be set on entry firing and the value will be equal to the one in Logger struct field.
	Level Level
	// Message passed to Trace, Debug, Info, Warn, Error, Fatal or Panic
	Message string
}

// Print logs a message at level Print on the compatibleLogger.
func (l *Entry) Print(args ...interface{}) {
	l.logger.Info(fmt.Sprint(args...))
}

// Println logs a message at level Print on the compatibleLogger.
func (l *Entry) Println(args ...interface{}) {
	l.logger.Info(fmt.Sprint(args...))
}

// Printf logs a message at level Print on the compatibleLogger.
func (l *Entry) Printf(format string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

// Trace logs a message at level Trace on the compatibleLogger.
func (l *Entry) Trace(args ...interface{}) {
	l.logger.Debug(fmt.Sprint(args...))
}

// Traceln logs a message at level Trace on the compatibleLogger.
func (l *Entry) Traceln(args ...interface{}) {
	l.logger.Debug(fmt.Sprint(args...))
}

// Tracef logs a message at level Trace on the compatibleLogger.
func (l *Entry) Tracef(format string, args ...interface{}) {
	l.logger.Debug(fmt.Sprintf(format, args...))
}

// Debug logs a message at level Debug on the compatibleLogger.
func (l *Entry) Debug(args ...interface{}) {
	l.logger.Debug(fmt.Sprint(args...))
}

// Debugln logs a message at level Debug on the compatibleLogger.
func (l *Entry) Debugln(args ...interface{}) {
	l.logger.Debug(fmt.Sprint(args...))
}

// Debugf logs a message at level Debug on the compatibleLogger.
func (l *Entry) Debugf(format string, args ...interface{}) {
	l.logger.Debug(fmt.Sprintf(format, args...))
}

// Info logs a message at level Info on the compatibleLogger.
func (l *Entry) Info(args ...interface{}) {
	l.logger.Info(fmt.Sprint(args...))
}

// Infoln logs a message at level Info on the compatibleLogger.
func (l *Entry) Infoln(args ...interface{}) {
	l.logger.Info(fmt.Sprint(args...))
}

// Infof logs a message at level Info on the compatibleLogger.
func (l *Entry) Infof(format string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

// Warn logs a message at level Warn on the compatibleLogger.
func (l *Entry) Warn(args ...interface{}) {
	l.logger.Warn(fmt.Sprint(args...))
}

// Warnln logs a message at level Warn on the compatibleLogger.
func (l *Entry) Warnln(args ...interface{}) {
	l.logger.Warn(fmt.Sprint(args...))
}

// Warnf logs a message at level Warn on the compatibleLogger.
func (l *Entry) Warnf(format string, args ...interface{}) {
	l.logger.Warn(fmt.Sprintf(format, args...))
}

// Warning logs a message at level Warn on the compatibleLogger.
func (l *Entry) Warning(args ...interface{}) {
	l.logger.Warn(fmt.Sprint(args...))
}

// Warningln logs a message at level Warning on the compatibleLogger.
func (l *Entry) Warningln(args ...interface{}) {
	l.logger.Warn(fmt.Sprint(args...))
}

// Warningf logs a message at level Warning on the compatibleLogger.
func (l *Entry) Warningf(format string, args ...interface{}) {
	l.logger.Warn(fmt.Sprintf(format, args...))
}

// Error logs a message at level Error on the compatibleLogger.
func (l *Entry) Error(args ...interface{}) {
	l.logger.Error(fmt.Sprint(args...))
}

// Errorln logs a message at level Error on the compatibleLogger.
func (l *Entry) Errorln(args ...interface{}) {
	l.logger.Error(fmt.Sprint(args...))
}

// Errorf logs a message at level Error on the compatibleLogger.
func (l *Entry) Errorf(format string, args ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, args...))
}

// Fatal logs a message at level Fatal on the compatibleLogger.
func (l *Entry) Fatal(args ...interface{}) {
	l.logger.Fatal(fmt.Sprint(args...))
}

// Fatalln logs a message at level Fatal on the compatibleLogger.
func (l *Entry) Fatalln(args ...interface{}) {
	l.logger.Fatal(fmt.Sprint(args...))
}

// Fatalf logs a message at level Fatal on the compatibleLogger.
func (l *Entry) Fatalf(format string, args ...interface{}) {
	l.logger.Fatal(fmt.Sprintf(format, args...))
}

// Panic logs a message at level Painc on the compatibleLogger.
func (l *Entry) Panic(args ...interface{}) {
	l.logger.Panic(fmt.Sprint(args...))
}

// Panicln logs a message at level Painc on the compatibleLogger.
func (l *Entry) Panicln(args ...interface{}) {
	l.logger.Panic(fmt.Sprint(args...))
}

// Panicf logs a message at level Painc on the compatibleLogger.
func (l *Entry) Panicf(format string, args ...interface{}) {
	l.logger.Panic(fmt.Sprintf(format, args...))
}

// WithError return a logger with an error field.
func (l *Entry) WithError(err error) *Entry {
	return &Entry{
		logger:  l.logger.With(zap.Error(err)),
		Logger:  l.Logger,
		Level:   l.Level,
		Message: l.Message,
	}
}

// WithField return a logger with an extra field.
func (l *Entry) WithField(key string, value interface{}) *Entry {
	return &Entry{
		logger:  l.logger.With(zap.Any(key, value)),
		Logger:  l.Logger,
		Level:   l.Level,
		Message: l.Message,
	}
}

// WithFields return a logger with extra fields.
func (l *Entry) WithFields(fields Fields) *Entry {
	if len(fields) == 0 {
		return l
	}
	i := 0
	var clog *Entry
	for k, v := range fields {
		if i == 0 {
			clog = l.WithField(k, v)
		} else {
			if clog != nil {
				clog = clog.WithField(k, v)
			}
		}
		i++
	}
	return clog
}
