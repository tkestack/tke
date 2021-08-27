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

import { Cluster } from '../../common/models';

export interface Ckafka extends Identifiable {
  /** ckafka的 instanceId */
  instanceId?: string;

  /** ckafka instanceName */
  instanceName?: string;

  /**实例状态 0: 创建中，1：运行中， 2：删除中 */
  status?: number;

  bandwith?: number;

  diskSize?: number;

  /** vpc网络id */
  vpcId?: string;

  /** 子网id */
  subnitId?: string;

  /** zoneId */
  zoneId?: number;

  /** topic 的数量 */
  topicNum?: number;

  /** vipList */
  vipList?: any[];
}

export interface CkafkaFilter {
  /** ckafka的状态 */
  status?: number;

  /** 集群的信息 */
  cluster?: Cluster;

  /** 当前的地域ID */
  regionId?: number;

  /** 是否能够拉取ckafka的列表 */
  isCanFetchCkafka?: boolean;
}

export interface CTopic extends Identifiable {
  /** topicId */
  topicId?: string;

  /** topicName */
  topicName?: string;
}

export interface CTopicFilter {
  /** instanceId */
  instanceId?: string;

  /** 当前的地域ID */
  regionId?: number;

  /** 是否能够拉取CTopic的列表 */
  isCanFetchCTopic?: boolean;
}
