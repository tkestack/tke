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
import { Identifiable } from '@tencent/ff-redux';

import { Validation } from '../../common/models';

export interface PortMap extends Identifiable {
  /**
   * 协议（TCP或UDP）
   */
  protocol?: string;
  v_protocol?: Validation;

  /**
   * 容器监听的端口
   */
  targetPort?: string;
  v_targetPort?: Validation;

  /**
   * 服务端口的端口
   */
  port?: string;
  v_port?: Validation;

  /**
   * 主机的端口
   */
  nodePort?: string;
  v_nodePort?: Validation;
}
