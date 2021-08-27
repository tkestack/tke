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

import { ResourceFilter } from './ResourceOption';

export interface Event extends Identifiable {
  /** 事件出现的次数 */
  count?: string;

  /** 首次出现的时间 */
  firstTimestamp?: string;

  /** 最后出现的时间 */
  lastTimestamp?: string;

  /** message */
  message?: string;

  /** metadata */
  metadata?: Metadata;

  /** involvedObject */
  involvedObject?: InvolvedObject;

  /** reason */
  reason?: string;

  /** source */
  source?: any;

  /** 事件的级别 */
  type?: string;
}

interface InvolvedObject {
  /** apiVersion */
  apiVersion?: string;

  /** kind */
  kind?: string;

  /** name */
  name?: string;

  /** namespace */
  namespace?: string;

  /** resourceVersion */
  resourceVersion?: string;

  /** uid */
  uid?: string;
}

interface Metadata {
  /** creationTimestamp */
  creationTimestamp?: string;

  /** name */
  name?: string;

  /** namespace */
  namespace?: string;

  /** resourceVersion */
  resourceVersion?: string;

  /** selfLink */
  selfLink?: string;

  /** uid */
  uid?: string;
}

/** 资源详情页当中的 事件filter数据类型EventFilter */
export interface EventFilter extends ResourceFilter {
  /** kind */
  kind?: string;

  /** workload的名称 */
  name?: string;
}
