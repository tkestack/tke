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

package cluster

const (
	tokenFileTemplate = `%s,admin,admin,system:masters
`

	schedulerPolicyConfig = `
{
   "apiVersion" : "v1",
   "extenders" : [
      {
         "apiVersion" : "v1beta1",
         "enableHttps" : false,
         "filterVerb" : "predicates",
         "managedResources" : [
            {
               "ignoredByScheduler" : false,
               "name" : "tencent.com/vcuda-core"
            }
         ],
         "nodeCacheCapable" : false,
         "urlPrefix" : "http://gpu-quota-admission:3456/scheduler"
      },
      {
         "apiVersion" : "v1beta1",
         "enableHttps" : false,
         "filterVerb" : "filter",
         "BindVerb": "bind",
         "weight": 1,
         "enableHttps": false,
         "managedResources" : [
            {
               "ignoredByScheduler" : true,
               "name" : "tkestack.io/eni-ip"
            }
         ],
         "nodeCacheCapable" : false,
         "urlPrefix" : "http://127.0.0.1:32760/v1"
      }
   ],
   "kind" : "Policy"
}
`
)
