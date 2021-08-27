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

export interface Replicaset extends Identifiable {
  /** metadata */
  metadata?: Metadata;

  /** spec */
  spec?: Spec;

  /** status */
  status?: Status;
}

interface Metadata {
  /** annotations */
  annotations?: {
    'deployment.kubernetes.io/desired-replicas'?: string;

    'deployment.kubernetes.io/max-replicas'?: string;

    'deployment.kubernetes.io/revision'?: string;

    [props: string]: string;
  };

  /** creationTimestamp */
  creationTimestamp?: string;

  /** generation */
  generation?: string;

  /** labels */
  labels?: {
    [props: string]: string;
  };

  /**name */
  name?: string;

  namespace?: string;

  [props: string]: any;
}

interface Spec {
  /** replicas */
  replicas?: string;

  /** selector */
  selector?: {
    matchLabels: {
      [props: string]: string;
    };
  };

  /** template */
  template?: {
    metadata?: {
      [props: string]: any;
    };

    spec?: {
      [props: string]: any;
    };
  };
}

interface Status {
  /** availableReplicas */
  availableReplicas?: string;

  /** conditions */
  conditions?: any;

  /** fullyLabeledReplicas */
  fullyLabeledReplicas?: string;

  /** observedGeneration */
  observedGeneration?: string;

  /** readyReplicas */
  readyReplicas?: string;

  /** replicas */
  replicas?: string;
}
