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

package http

import "testing"

func TestExternalAddress(t *testing.T) {
	cases := []struct {
		scheme string
		host   string
		port   int
		expect string
		desc   string
	}{
		{
			scheme: "http",
			host:   "example.com",
			port:   80,
			expect: "http://example.com",
			desc:   "http without port",
		}, {
			scheme: "http",
			host:   "example.com",
			port:   81,
			expect: "http://example.com:81",
			desc:   "http with port",
		}, {
			scheme: "https",
			host:   "example.com",
			port:   443,
			expect: "https://example.com",
			desc:   "https without port",
		}, {
			scheme: "https",
			host:   "example.com",
			port:   8443,
			expect: "https://example.com:8443",
			desc:   "https with port",
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			address := ExternalAddress(c.scheme, c.host, c.port)
			if address != c.expect {
				t.Fatalf("expect %#v but got %#v", c.expect, address)
			}
		})
	}
}

func TestExternalEndpoint(t *testing.T) {
	cases := []struct {
		scheme string
		host   string
		port   int
		path   string
		expect string
		desc   string
	}{
		{
			scheme: "http",
			host:   "example.com",
			port:   80,
			expect: "http://example.com/",
			desc:   "http without path",
		}, {
			scheme: "http",
			host:   "example.com",
			port:   81,
			path:   "/test",
			expect: "http://example.com:81/test",
			desc:   "http with path",
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			address := ExternalEndpoint(c.scheme, c.host, c.port, c.path)
			if address != c.expect {
				t.Fatalf("expect %#v but got %#v", c.expect, address)
			}
		})
	}
}
