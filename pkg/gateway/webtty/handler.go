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

package webtty

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"

	"k8s.io/client-go/tools/remotecommand"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/transport/spdy"
	"tkestack.io/tke/pkg/gateway/token"
	"tkestack.io/tke/pkg/util/log"
)

type handler struct {
	url          *url.URL
	roundTripper http.RoundTripper
	upgrader     spdy.Upgrader
}

// NewHandler to create a reverse proxy handler and returns it.
// The reverse proxy will parse the requested cookie content, get the token in
// it, and append it as the http request header to the backend service component.
func NewHandler(address string) (http.Handler, error) {
	u, err := url.Parse(address)
	if err != nil {
		log.Error("Failed to parse backend service address", log.String("address", address), log.Err(err))
		return nil, err
	}

	cfg := &rest.Config{
		Host: u.Host,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}
	roundTripper, upgrader, err := spdy.RoundTripperFor(cfg)
	if err != nil {
		return nil, err
	}

	return &handler{
		url:          u,
		roundTripper: roundTripper,
		upgrader:     upgrader,
	}, nil
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// read cookie
	t, err := token.RetrieveToken(req)
	if err != nil {
		log.Error("Failed to retrieve token from webtty", log.Err(err))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	clusterName := req.URL.Query().Get("clusterName")
	if clusterName == "" {
		log.Error("Failed to get cluster name from webtty request")
		http.Error(w, "Invalid cluster name", http.StatusBadRequest)
		return
	}

	namespace := req.URL.Query().Get("namespace")
	if namespace == "" {
		log.Error("Failed to get namespace from webtty request")
		http.Error(w, "Invalid namespace", http.StatusBadRequest)
		return
	}

	podName := req.URL.Query().Get("podName")
	if podName == "" {
		log.Error("Failed to get pod name from webtty request")
		http.Error(w, "Invalid pod name", http.StatusBadRequest)
		return
	}

	containerName := req.URL.Query().Get("containerName")
	if containerName == "" {
		log.Error("Failed to get container name from webtty request")
		http.Error(w, "Invalid container name", http.StatusBadRequest)
		return
	}

	command := req.URL.Query().Get("command")
	if command == "" {
		command = "/bin/sh"
	}

	wsConn, err := InitWebsocket(w, req)
	if err != nil {
		log.Error("Failed tp init websocket", log.Err(err))
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	u := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/exec?container=%s&command=%s&stdin=true&stdout=true&stderr=false&tty=true", namespace, podName, containerName, command)
	reqURL, err := url.Parse(u)
	if err != nil {
		log.Error("Failed to generate pod exec url", log.Err(err))
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	reqURL.Host = h.url.Host
	reqURL.Scheme = h.url.Scheme

	executor, err := NewSPDYExecutorForTransports(h.roundTripper, h.upgrader, http.MethodPost, reqURL, clusterName, strings.TrimSpace(t.ID))
	if err != nil {
		log.Error("Failed to create SPDY executor", log.Err(err))
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	handler := &streamHandler{wsConn: wsConn, resizeEvent: make(chan remotecommand.TerminalSize)}
	if err := executor.Stream(remotecommand.StreamOptions{
		Stdin:             handler,
		Stdout:            handler,
		Stderr:            handler,
		TerminalSizeQueue: handler,
		Tty:               true,
	}); err != nil {
		log.Error("Failed to stream exec command", log.Err(err))
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
}

type streamHandler struct {
	wsConn      *WSConnection
	resizeEvent chan remotecommand.TerminalSize
}

type xtermMessage struct {
	Input   string `json:"input"`
	MsgType string `json:"type"`
	Rows    uint16 `json:"rows"`
	Cols    uint16 `json:"cols"`
}

func (handler *streamHandler) Next() (size *remotecommand.TerminalSize) {
	ret := <-handler.resizeEvent
	size = &ret
	return
}

func (handler *streamHandler) Read(p []byte) (int, error) {
	msg, err := handler.wsConn.Read()
	if err != nil {
		return 0, err
	}

	var xtermMsg xtermMessage
	if err = json.Unmarshal(msg.Data, &xtermMsg); err != nil {
		return 0, nil
	}

	size := 0
	if xtermMsg.MsgType == "resize" {
		handler.resizeEvent <- remotecommand.TerminalSize{
			Width:  xtermMsg.Cols,
			Height: xtermMsg.Rows,
		}
	} else if xtermMsg.MsgType == "input" {
		size = len(xtermMsg.Input)
		copy(p, xtermMsg.Input)
	}
	return size, nil
}

func (handler *streamHandler) Write(p []byte) (size int, err error) {
	size = len(p)
	err = handler.wsConn.Write(websocket.TextMessage, p)
	return
}
