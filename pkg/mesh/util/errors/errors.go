/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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
 *
 */

package errors

import (
	"strings"
	"sync"
)

// MultiError is a collection of errors that implements the error interface.
type MultiError struct {
	mu   sync.RWMutex
	errs []error
}

// NewMultiError returns a new MultiError.
func NewMultiError() *MultiError {
	return &MultiError{}
}

// Error returns the list of errors separated by newlines.
func (merr *MultiError) Error() string {
	merr.mu.RLock()
	defer merr.mu.RUnlock()

	var errs []string
	for _, err := range merr.errs {
		errs = append(errs, err.Error())
	}

	return strings.Join(errs, "\n")
}

// Errors returns the error slice containing the error collection.
func (merr *MultiError) Errors() []error {
	merr.mu.RLock()
	defer merr.mu.RUnlock()

	return merr.errs
}

// Add appends an error to the error collection.
func (merr *MultiError) Add(err error) {
	merr.mu.Lock()
	defer merr.mu.Unlock()

	// Unwrap *MultiError to ensure that depth never exceeds 1.
	if me, ok := err.(*MultiError); ok {
		merr.errs = append(merr.errs, me.Errors()...)
		return
	}

	merr.errs = append(merr.errs, err)
}

// Empty returns whether the *MultiError contains any errors.
func (merr *MultiError) Empty() bool {
	return len(merr.errs) == 0
}
