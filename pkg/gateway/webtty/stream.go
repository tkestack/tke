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
	"fmt"
	"net/http"
	"net/url"

	"k8s.io/apimachinery/pkg/util/httpstream"
	apimachinerycommand "k8s.io/apimachinery/pkg/util/remotecommand"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/transport/spdy"
	"tkestack.io/tke/pkg/platform/apiserver/filter"
	"tkestack.io/tke/pkg/util/log"
)

type streamCreator interface {
	CreateStream(headers http.Header) (httpstream.Stream, error)
}

type streamProtocolHandler interface {
	stream(conn streamCreator) error
}

// streamExecutor handles transporting standard shell streams over an httpstream connection.
type streamExecutor struct {
	upgrader  spdy.Upgrader
	transport http.RoundTripper

	method      string
	url         *url.URL
	protocols   []string
	token       string
	clusterName string
}

// NewSPDYExecutorForTransports connects to the provided server using the given transport,
// upgrades the response using the given upgrader to multiplexed bidirectional streams.
func NewSPDYExecutorForTransports(transport http.RoundTripper, upgrader spdy.Upgrader, method string, url *url.URL, clusterName string, token string) (remotecommand.Executor, error) {
	return NewSPDYExecutorForProtocols(
		transport, upgrader, method, url, clusterName, token,
		apimachinerycommand.StreamProtocolV4Name,
		apimachinerycommand.StreamProtocolV3Name,
		apimachinerycommand.StreamProtocolV2Name,
		apimachinerycommand.StreamProtocolV1Name,
	)
}

// NewSPDYExecutorForProtocols connects to the provided server and upgrades the connection to
// multiplexed bidirectional streams using only the provided protocols. Exposed for testing, most
// callers should use NewSPDYExecutor or NewSPDYExecutorForTransports.
func NewSPDYExecutorForProtocols(transport http.RoundTripper, upgrader spdy.Upgrader, method string, url *url.URL, clusterName string, token string, protocols ...string) (remotecommand.Executor, error) {
	return &streamExecutor{
		upgrader:    upgrader,
		transport:   transport,
		method:      method,
		url:         url,
		protocols:   protocols,
		token:       token,
		clusterName: clusterName,
	}, nil
}

// Stream opens a protocol streamer to the server and streams until a client closes
// the connection or the server disconnects.
func (e *streamExecutor) Stream(options remotecommand.StreamOptions) error {
	req, err := http.NewRequest(e.method, e.url.String(), nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", e.token))
	req.Header.Add(filter.ClusterNameHeaderKey, e.clusterName)

	conn, protocol, err := spdy.Negotiate(
		e.upgrader,
		&http.Client{Transport: e.transport},
		req,
		e.protocols...,
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	var streamer streamProtocolHandler

	switch protocol {
	case apimachinerycommand.StreamProtocolV4Name:
		streamer = newStreamProtocolV4(options)
	case apimachinerycommand.StreamProtocolV3Name:
		streamer = newStreamProtocolV3(options)
	case apimachinerycommand.StreamProtocolV2Name:
		streamer = newStreamProtocolV2(options)
	case "":
		log.Infof("The server did not negotiate a streaming protocol version. Falling back to %s", apimachinerycommand.StreamProtocolV1Name)
		fallthrough
	case apimachinerycommand.StreamProtocolV1Name:
		streamer = newStreamProtocolV1(options)
	}

	return streamer.stream(conn)
}
