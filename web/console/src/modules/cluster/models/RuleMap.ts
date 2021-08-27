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

export interface RuleMap extends Identifiable {
  /** host 域名 */
  host?: string;
  v_host?: Validation;

  /** url路径 */
  path?: string;
  v_path?: Validation;

  /** 后端服务名称 */
  serviceName?: string;
  v_serviceName?: Validation;

  /** 后端服务端口 */
  servicePort?: string | number;
  v_servicePort?: Validation;

  /** protocol */
  protocol?: string;
  v_protocol?: Validation;
}
