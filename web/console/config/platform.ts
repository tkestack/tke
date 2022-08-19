/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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
/** 获取当前控制台modules的 域名映射表 */
export enum ConsoleModuleEnum {
  /** tke-apiserver 版本 */
  PLATFORM = 'platform',

  /** 业务的版本详情 */
  Business = 'business',

  /** 通知模块 */
  Notify = 'notify',

  /** 告警模块 */
  Monitor = 'monitor',

  /** 镜像仓库 */
  Registry = 'registry',

  /** 日志模块 */
  LogAgent = 'logagent',

  /** 认证模块 */
  Auth = 'auth',

  /** 审计模块 */
  Audit = 'audit',

  /** Helm应用模块 */
  Application = 'application',

  /** 中间件列表模块 */
  Middleware = 'middleware'
}

export enum PlatformTypeEnum {
  /** 平台 */
  Manager = 'manager',

  /** 业务 */
  Business = 'business'
}
