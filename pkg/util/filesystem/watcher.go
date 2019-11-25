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

// Portions Copyright 2014 The Kubernetes Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package filesystem

import (
	"github.com/fsnotify/fsnotify"
)

// FSWatcher is a callback-based filesystem watcher abstraction for fsnotify.
type FSWatcher interface {
	// Initializes the watcher with the given watch handlers.
	// Called before all other methods.
	Init(FSEventHandler, FSErrorHandler) error

	// Starts listening for events and errors.
	// When an event or error occurs, the corresponding handler is called.
	Run()

	// Add a filesystem path to watch
	AddWatch(path string) error
}

// FSEventHandler is called when a fsnotify event occurs.
type FSEventHandler func(event fsnotify.Event)

// FSErrorHandler is called when a fsnotify error occurs.
type FSErrorHandler func(err error)

type fsNotifyWatcher struct {
	watcher      *fsnotify.Watcher
	eventHandler FSEventHandler
	errorHandler FSErrorHandler
}

var _ FSWatcher = &fsNotifyWatcher{}

// NewFSNotifyWatcher returns an implementation of FSWatcher that continuously listens for
// fsNotify events and calls the event handler as soon as an event is received.
func NewFSNotifyWatcher() FSWatcher {
	return &fsNotifyWatcher{}
}

func (w *fsNotifyWatcher) AddWatch(path string) error {
	return w.watcher.Add(path)
}

func (w *fsNotifyWatcher) Init(eventHandler FSEventHandler, errorHandler FSErrorHandler) error {
	var err error
	w.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	w.eventHandler = eventHandler
	w.errorHandler = errorHandler
	return nil
}

func (w *fsNotifyWatcher) Run() {
	go func() {
		defer func() {
			_ = w.watcher.Close()
		}()
		for {
			select {
			case event := <-w.watcher.Events:
				if w.eventHandler != nil {
					w.eventHandler(event)
				}
			case err := <-w.watcher.Errors:
				if w.errorHandler != nil {
					w.errorHandler(err)
				}
			}
		}
	}()
}
