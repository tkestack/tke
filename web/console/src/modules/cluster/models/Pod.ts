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

export interface Pod extends Identifiable {
  /** metadata */
  metadata?: Metadata;

  /** spec */
  spec?: Spec;

  /** status */
  status?: Status;
}

interface Metadata {
  annotations?: {
    'kubernetes.io/created-by'?: string;

    [props: string]: string;
  };

  creationTimestamp?: string;

  deletionTimestamp?: string;

  name?: string;

  namespace?: string;

  [props: string]: any;
}

interface Spec {
  containers?: PodContainer[];

  [props: string]: any;
}

interface Status {
  containerStatuses?: any[];

  conditions?: any[];

  phase?: string;

  reason?: string;

  qosClass?: string;

  /** pod所在node 的ip */
  hostIP?: string;

  /** pod的ip */
  podIP?: string;

  /** pod启动时间 */
  startTime?: string;
}

export interface PodContainer extends Identifiable {
  env?: Env[];

  image?: string;

  imagePullPolicy?: string;

  name?: string;

  resources?: any;

  terminationMessagePath?: string;

  terminationMessagePolicy?: string;

  [props: string]: any;
}

interface Env {
  name?: string;

  value?: string;
}

/** Node详情页内的pod列表的筛选框的项 */
export interface PodFilterInNode {
  /** 命名空间 */
  namespace?: string;

  /** podName */
  podName?: string;

  /** pod的状态值 */
  phase?: string;

  ip?: string;
}
