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

import { Identifiable } from '@tencent/ff-redux';

import { BaseType, Validation } from './';

export interface PortMapItem extends Identifiable {
  /**协议（TCP或UDP）*/
  protocol?: string;
  v_protocol?: Validation;

  /**容器监听的端口 */
  containerPort?: string;
  v_containerPort?: Validation;

  /**外网lb的端口 */
  lbPort?: string;
  v_lbPort?: Validation;

  /**主机端口生成方式 */
  generateType?: string;

  /**主机的端口 */
  nodePort?: string;
  v_nodePort?: Validation;
}

/**协议类型 */
export interface Protocol extends BaseType {}
