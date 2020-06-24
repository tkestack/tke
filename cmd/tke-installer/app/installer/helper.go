/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package installer

import (
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/emicklei/go-restful"
	"tkestack.io/tke/pkg/util/log"
)

func globalLogging(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	now := time.Now()

	reqBytes, err := httputil.DumpRequest(req.Request, true)
	if err != nil {
		_ = resp.WriteError(http.StatusInternalServerError, err)
		return
	}

	log.Infof("raw http request:\n%s", reqBytes)
	chain.ProcessFilter(req, resp)
	log.Infof("%s %s %v", req.Request.Method, req.Request.URL, time.Since(now))
}
