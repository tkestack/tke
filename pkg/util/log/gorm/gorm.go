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

package gorm

import (
	"fmt"
	"go.uber.org/zap"
	"time"
	"unicode"
)

// Logger is an alternative implementation of *gorm.Logger
type Logger struct {
	logger *zap.Logger
}

// NewLogger create the gorm Logger object by initialized zap logger.
func NewLogger(logger *zap.Logger) *Logger {
	return &Logger{logger}
}

// Print passes arguments to Println
func (l *Logger) Print(values ...interface{}) {
	l.Println(values)
}

// Println format & print log
func (l *Logger) Println(values []interface{}) {
	if len(values) > 1 {
		level := values[0]
		if level == "sql" {
			var fields []zap.Field
			if len(values) > 4 {
				fields = append(fields, zap.String("system", "gorm"), zap.String("gorm.sql", values[3].(string)), zap.Strings("gorm.values", formattedValues(values[4].([]interface{}))), zap.Duration("gorm.time_ms", values[2].(time.Duration)))
			}
			l.logger.Debug("finished SQL query", fields...)
		} else {
			l.logger.Debug("finished SQL event", zap.Any("message", values[1:]))
		}
	}
}

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

func formattedValues(rawValues []interface{}) []string {
	formattedValues := make([]string, 0, len(rawValues))
	for _, value := range rawValues {
		switch v := value.(type) {
		case time.Time:
			formattedValues = append(formattedValues, fmt.Sprint(v))
		case []byte:
			if str := string(v); isPrintable(str) {
				formattedValues = append(formattedValues, fmt.Sprint(str))
			} else {
				formattedValues = append(formattedValues, "<binary>")
			}
		default:
			str := "NULL"
			if v != nil {
				str = fmt.Sprint(v)
			}
			formattedValues = append(formattedValues, str)
		}
	}
	return formattedValues
}
