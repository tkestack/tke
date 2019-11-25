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

package route

import (
	"github.com/emicklei/go-restful"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/pkg/monitor/services"
)

// MonitorResource represents the REST resource of monitor.
type MonitorResource struct {
	PlatformClient platformversionedclient.PlatformV1Interface
	RulesOperator  services.BackendConfigProcessor
}

// WebService returns the restful webservice object.
func (r *MonitorResource) WebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/apis/v1/monitor")
	ws.Produces(restful.MIME_JSON)
	ws.Consumes(restful.MIME_JSON)

	// Register rules path
	r.RulesOperator.RegisterWebService(ws)

	// TODO add alarm receivers

	// TODO add send webhook

	return ws
}
