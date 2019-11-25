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

import "fmt"

// ExternalAddress to build external http address by given scheme, host and port.
func ExternalAddress(scheme, host string, port int) string {
	switch scheme {
	case "http", "HTTP", "Http":
		if port == 80 {
			return fmt.Sprintf("http://%s", host)
		}
		return fmt.Sprintf("http://%s:%d", host, port)
	case "https", "HTTPS", "Https":
		if port == 443 {
			return fmt.Sprintf("https://%s", host)
		}
		return fmt.Sprintf("https://%s:%d", host, port)
	default:
		return fmt.Sprintf("%s://%s:%d", scheme, host, port)
	}
}

// ExternalEndpoint to build external http endpoint by given address and path.
func ExternalEndpoint(scheme, host string, port int, path string) string {
	if path == "" {
		path = "/"
	}
	return fmt.Sprintf("%s%s", ExternalAddress(scheme, host, port), path)
}
