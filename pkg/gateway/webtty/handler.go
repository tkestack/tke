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
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"k8s.io/client-go/tools/remotecommand"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/transport/spdy"
	"tkestack.io/tke/pkg/gateway/token"
	"tkestack.io/tke/pkg/util/log"
)

const (
	// endOfTransmission defines end transmission flag.
	endOfTransmission = "\u0004"
	// bufferSize specify I/O buffer sizes. If a buffer size is zero, then buffers
	// allocated by the HTTP server are used. The I/O buffer sizes do not limit
	// the size of the messages that can be sent or received.
	bufferSize = 1024
	// writeWait defines time allowed to write a message to the peer.
	writeWait = 10 * time.Second
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  bufferSize,
	WriteBufferSize: bufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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

	projectName := req.URL.Query().Get("projectName")

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

	conn, err := wsUpgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Error("Failed initialize websocket connection", log.Err(err))
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

	executor, err := NewSPDYExecutorForTransports(h.roundTripper, h.upgrader, http.MethodPost, reqURL, clusterName, projectName, strings.TrimSpace(t.ID))
	if err != nil {
		log.Error("Failed to create SPDY executor", log.Err(err))
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	handler := &streamHandler{conn: conn, resizeEvent: make(chan remotecommand.TerminalSize)}
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
	conn        *websocket.Conn
	resizeEvent chan remotecommand.TerminalSize
}

type xtermMessage struct {
	Input   string `json:"input"`
	MsgType string `json:"type"`
	Rows    uint16 `json:"rows"`
	Cols    uint16 `json:"cols"`
}

// Next handles pty->process resize events, called in a loop from remotecommand
// as long as the process is running.
func (handler *streamHandler) Next() *remotecommand.TerminalSize {
	select {
	case size := <-handler.resizeEvent:
		return &size
	}
}

// Read handles pty->process messages (stdin, resize), called in a loop from
// remotecommand as long as the process is running.
func (handler *streamHandler) Read(p []byte) (int, error) {
	var xtermMsg xtermMessage
	if err := handler.conn.ReadJSON(&xtermMsg); err != nil {
		log.Error("Failed to read and unmarshal message from websocket", log.Err(err))
		return copy(p, endOfTransmission), err
	}

	switch xtermMsg.MsgType {
	case "resize":
		handler.resizeEvent <- remotecommand.TerminalSize{
			Width:  xtermMsg.Cols,
			Height: xtermMsg.Rows,
		}
		return 0, nil
	case "input":
		if len(xtermMsg.Input) == 0 {
			return 0, nil
		}
		input, err := decodeInputMessage(xtermMsg.Input)
		if err != nil {
			return copy(p, endOfTransmission), err
		}
		return copy(p, input), nil
	default:
		log.Error("Unknown message type from websocket", log.String("type", xtermMsg.MsgType))
		return copy(p, endOfTransmission), fmt.Errorf("unknown message type %s from websocket", xtermMsg.MsgType)
	}
}

// Write handles process->pty stdout, called from remotecommand whenever there
// is any output.
func (handler *streamHandler) Write(p []byte) (int, error) {
	size := len(p)
	if size == 0 {
		return 0, nil
	}
	output := encodeOutputMessage(p)
	_ = handler.conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err := handler.conn.WriteMessage(websocket.TextMessage, []byte(output)); err != nil {
		log.Error("Failed write message to websocket connection", log.Err(err))
		return size, err
	}
	return size, nil
}

// Close shutdown the websocket connection.
func (handler *streamHandler) Close() error {
	return handler.conn.Close()
}

func decodeInputMessage(input string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(input)
}

func encodeOutputMessage(output []byte) string {
	return base64.StdEncoding.EncodeToString(output)
}
