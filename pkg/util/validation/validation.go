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

package validation

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"
)

// IsHTTPSReachle tests that https://host:port is reachble in timeout.
func IsHTTPSReachle(host string, port int32, timeout time.Duration) error {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout: timeout,
			}).DialContext,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	url := fmt.Sprintf("https://%s:%d", host, port)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
