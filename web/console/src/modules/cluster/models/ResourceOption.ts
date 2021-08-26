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

import { Identifiable, FFListModel } from '@tencent/ff-redux';

export interface ResourceOption {
  /** 具体的resource列表 */
  ffResourceList?: FFListModel<Resource, ResourceFilter>;

  /** resource的多选选择 */
  resourceMultipleSelection?: Resource[];

  /** resourceDeleteSelection */
  resourceDeleteSelection?: Resource[];
}

export interface Resource extends Identifiable {
  /** metadata */
  metadata?: any;

  /** spec */
  spec?: any;

  /** status */
  status?: any;

  /** data */
  data?: any;

  /** other */
  [props: string]: any;
}

export interface ResourceFilter {
  /** 命名空间 */
  namespace?: string;

  /** 集群id */
  clusterId?: string;

  /** 地域id */
  regionId?: number;

  /** name */
  specificName?: string;

  meshId?: string;

  labelSelector?: string;
}

export interface DifferentInterfaceResourceOperation {
  query?: {
    [props: string]: any;
  };
  extraResource?: string;
}
