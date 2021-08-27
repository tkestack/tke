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
import { any } from 'prop-types';

import { Identifiable } from '@tencent/ff-redux';

import { LogStashSpec } from './LogStashEdit';

export interface Log extends Identifiable {
  apiVersion?: string;
  kind?: string;
  metadata?: {
    creationTimestamp: string;
    name: string;
    namespace: string;
    resourceVersion?: string;
  };
  spec: {
    input: {
      type?: string;
      [props: string]: any;
    };
    output: {
      type?: string;
      [props: string]: any;
    };
  };
  [props: string]: any;
}

export interface LogFilter {
  /** 地域的id */
  regionId?: number;

  /** 日志收集器的ID*/
  collectorId?: string;

  /** 日志收集器所属的集群ID*/
  clusterId?: string;

  /** 日志组件名称 */
  logAgentName?: string;

  /** 根据状态进行复选 */
  status?: string;

  /** 是否清除*/
  isClear?: boolean;

  /**命名空间 */
  namespace?: string;

  /** specificName */
  specificName?: string;
}

export interface LogOperator {
  /**
   * 地域
   */
  regionId?: number;

  /**
   * 集群id
   */
  clusterId?: string;

  /**
   * 当前的编辑类型 create | update
   */
  mode?: string;

  /**
   * 日志收集器id
   */
  collectorName?: string;
}
