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

export interface Resource extends Identifiable {
  /** metadata */
  metadata: {
    [props: string]: any;
  };

  /** spec */
  spec: {
    [props: string]: any;
  };

  /** data */
  data?: {
    [props: string]: any;
  };

  /** status */
  status: {
    [props: string]: any;
  };

  value?: any;

  text?: any;

  /** other */
  [props: string]: any;
}

export interface ResourceFilter {
  /** 命名空间 */
  namespace?: string;

  /** 集群id */
  clusterId?: string;

  /** 集群日志组件 */
  logAgentName?: string;

  /** 地域id */
  regionId?: number;

  /** name */
  specificName?: string;
}
