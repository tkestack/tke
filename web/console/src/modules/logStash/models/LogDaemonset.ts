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

import { cluster } from 'config/resource/k8sConfig';

import { Identifiable } from '@tencent/ff-redux';

export interface LogDaemonset extends Identifiable {
  /** kind */
  kind?: string;

  /**apiVersion */
  apiVersion?: string;

  /** metadata */
  metadata?: Metadata;

  /** spec */
  spec?: Spec;

  /** status */
  status?: Status;
}

interface Metadata {
  creationTimestamp?: string;

  name?: string;

  [props: string]: any;
}

interface Spec {
  clusterName?: string;

  [props: string]: any;
}

interface Status {
  phase?: string;

  reason?: string;

  retryCount?: number;

  [props: string]: any;
}

export interface LogDaemonSetFliter {
  clusterId?: string;

  specificName?: string;
}

export interface LogDaemonSetStatus {
  phase?: string;
  reason?: string;
}
