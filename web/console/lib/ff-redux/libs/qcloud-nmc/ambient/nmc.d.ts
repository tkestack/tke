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

/// <reference path="constants.d.ts" />
/// <reference path="seajs.d.ts" />
/// <reference path="manager.d.ts" />
/// <reference path="appUtil.d.ts" />
/// <reference path="qccomponent.d.ts" />
/// <reference path="net.d.ts" />
/// <reference path="router.d.ts" />
/// <reference path="eventtarget.d.ts" />
/// <reference path="tips.d.ts" />

declare namespace nmc {
  interface Require {
    (id: string): any;
    (id: "router"): nmc.Router;
    (id: "qccomponent"): nmc.Bee;
    (id: "tips"): nmc.Tips;
    (id: "net"): nmc.Net;
    (id: "$"): JQueryStatic;
    (id: "manager"): nmc.Manager;
    (id: "appUtil"): nmc.AppUtil;
    (id: "config/constants"): nmc.Constants;
    async(modules: string[], callback: Function);
  }
  export function render(html: string, moduleName: string): void;
}
