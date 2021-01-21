/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
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

package simplerest

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	restclient "k8s.io/client-go/rest"
)

// Interface rest client interface
type Interface interface {
	Verb(verb string) *Request
	Post() *Request
	Put() *Request
	Get() *Request
	Delete() *Request
}

// simple rest client , inspired by k8s client-go restClient
type Client struct {
	// base is the root URL for all invocations of the client
	base *url.URL

	// content describes how a RESTClient encodes and decodes responses.
	content ClientContentConfig

	// Set specific behavior of the client.  If not set http.DefaultClient will be used.
	Client *http.Client
}

type ClientContentConfig struct {
	// AcceptContentTypes specifies the types the client will accept and is optional.
	// If not set, ContentType will be used to define the Accept header
	AcceptContentTypes string
	// ContentType specifies the wire format used to communicate with the server.
	// This value will be set as the Accept header on requests made to the server if
	// AcceptContentTypes is not set, and as the default content type on any object
	// sent to the server. If not set, "application/json" is used.
	ContentType string
}

// RESTClientFor new simple rest client for config
func RESTClientFor(config *restclient.Config) (*Client, error) {

	baseURL, err := defaultServerURLFor(config)
	if err != nil {
		return nil, err
	}

	// TODO do not coupling with client-go
	transport, err := restclient.TransportFor(config)
	if err != nil {
		return nil, err
	}

	var httpClient *http.Client
	if transport != http.DefaultTransport {
		httpClient = &http.Client{Transport: transport}
		if config.Timeout > 0 {
			httpClient.Timeout = config.Timeout
		}
	}

	clientContent := ClientContentConfig{
		AcceptContentTypes: config.AcceptContentTypes,
		ContentType:        config.ContentType,
	}

	return NewRESTClient(baseURL, clientContent, httpClient)
}

// NewRESTClient creates a new RESTClient. This client performs generic REST functions
// such as Get, Put, Post, and Delete on specified paths.
func NewRESTClient(baseURL *url.URL, config ClientContentConfig, client *http.Client) (*Client, error) {
	if len(config.ContentType) == 0 {
		config.ContentType = "application/json"
	}

	base := *baseURL
	if !strings.HasSuffix(base.Path, "/") {
		base.Path += "/"
	}
	base.RawQuery = ""
	base.Fragment = ""

	return &Client{
		base:    &base,
		content: config,
		Client:  client,
	}, nil
}

func defaultServerURLFor(config *restclient.Config) (*url.URL, error) {
	hasCA := len(config.CAFile) != 0 || len(config.CAData) != 0
	hasCert := len(config.CertFile) != 0 || len(config.CertData) != 0
	defaultTLS := hasCA || hasCert || config.Insecure
	base := config.Host
	if base == "" {
		base = "localhost"
	}

	hostURL, err := url.Parse(base)
	if err != nil || hostURL.Scheme == "" || hostURL.Host == "" {
		scheme := "http://"
		if defaultTLS {
			scheme = "https://"
		}
		hostURL, err = url.Parse(scheme + base)
		if err != nil {
			return nil, err
		}
		if hostURL.Path != "" && hostURL.Path != "/" {
			return nil, fmt.Errorf("host must be a URL or a host:port pair: %q", base)
		}
	}
	return hostURL, nil
}

func (c *Client) Verb(verb string) *Request {
	return NewRequest(c).Verb(verb)
}

// Post begins a POST request. Short for c.Verb("POST").
func (c *Client) Post() *Request {
	return c.Verb("POST")
}

// Put begins a PUT request. Short for c.Verb("PUT").
func (c *Client) Put() *Request {
	return c.Verb("PUT")
}

// Get begins a GET request. Short for c.Verb("GET").
func (c *Client) Get() *Request {
	return c.Verb("GET")
}

// Delete begins a DELETE request. Short for c.Verb("DELETE").
func (c *Client) Delete() *Request {
	return c.Verb("DELETE")
}
