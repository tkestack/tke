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
import { extend, Identifiable, RecordSet } from '@tencent/ff-redux';

export interface LbcfResource extends Identifiable {
  /** metadata */
  metadata?: {
    name: string;
    namespace: string;
  };

  /** spec */
  spec?: {
    lbDriver: string;
    lbSpec: {
      lbID?: string;
      lbVpcID?: string;
    };
    backGroups: BackendGroup[];
  };

  /** status */
  status?: any;

  /** other */
  [props: string]: any;
}

export interface BackendGroup {
  name: string;

  pods?: PodBackend;

  service?: ServiceBackend;

  static?: string[];

  status: {
    backends: number;
    registeredBackends: number;
  };
  backendRecords: BackendRecord[];
}

export interface PodBackend {
  labels: {
    [props: string]: string;
  };
  port: {
    portNumber: number;
    protocol: string;
  };
  byName: string[];
}

export interface ServiceBackend {
  name: string;
  port: {
    portNumber: number;
    protocol: string;
  };
  nodeSelector: {
    [props: string]: string;
  };
}

export interface BackendRecord {
  name: string;
  backendAddr: string;
  conditions: Condition[];
}

interface Condition {
  lastTransitionTime: string;
  message: string;
  status: string;
  type: string;
}
