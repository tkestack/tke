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
import { Validation, EnvItem, BaseType } from './';

export interface HealthCheck {
  /**是否开启存活检查 */
  isOpenLiveCheck?: boolean;

  /**存活检查参数 */
  liveCheck?: HealthCheckItem;

  /**是否开启就绪检查 */
  isOpenReadyCheck?: boolean;

  /**就绪检查参数*/
  readyCheck?: HealthCheckItem;
}

export interface HealthCheckItem {
  /**检查方法 */
  checkMethod?: string;
  v_checkMethod?: Validation;

  /**检查端口 */
  port?: number;
  v_port?: Validation;

  /**检查协议 */
  protocol?: string;
  v_protocol?: Validation;

  /**检查路径 */
  path?: string;
  v_path?: Validation;

  /**命令 */
  cmd?: string;
  v_cmd?: Validation;

  /**启动延时 */
  delayTime?: number;
  v_delayTime?: Validation;

  /**响应超时 */
  timeOut?: number;
  v_timeOut?: Validation;

  /**间隔时间 */
  intervalTime?: number;
  v_intervalTime?: Validation;

  /**健康阈值 */
  healthNum?: number;
  v_healthNum?: Validation;

  /**不健康阈值 */
  unhealthNum?: number;
  v_unhealthNum?: Validation;
}

/**健康检查协议 */
export interface HttpType extends BaseType {}

/**健康检查类型 */
export interface CheckType extends BaseType {}

/**健康检查方法 */
export interface CheckMethod extends BaseType {}
