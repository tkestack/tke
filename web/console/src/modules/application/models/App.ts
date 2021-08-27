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

export interface App extends Identifiable {
  metadata?: {
    creationTimestamp?: string;
    name?: string;
    namespace?: string;
    generation?: number;
  };

  spec?: {
    dryRun?: boolean;
    chart?: {
      chartGroupName?: string;
      chartName?: string;
      chartVersion?: string;
      tenantID?: string;
    };
    name?: string;
    targetCluster?: string;
    tenantID?: string;
    type?: string;
    values?: {
      rawValues?: string;
      rawValuesType?: string;
      values?: string[];
    };
  };
  status?: {
    lastTransitionTime?: string;
    phase?: string;
    releaseLastUpdated?: string;
    releaseStatus?: string;
    revision?: number;
    observedGeneration?: number;
    manifest?: string;
  };
}

export interface AppFilter {
  cluster?: string;
  namespace?: string;
}

export interface AppDetailFilter {
  cluster?: string;
  namespace?: string;
  name?: string;
}

export interface AppCreation extends Identifiable {
  metadata: {
    namespace: string;
  };

  spec?: {
    dryRun?: boolean;
    chart?: {
      chartGroupName?: string;
      chartName?: string;
      chartVersion?: string;
      tenantID?: string;
    };
    tenantID?: string;
    name?: string;
    targetCluster?: string;
    type?: string;
    values?: {
      rawValues?: string;
      rawValuesType?: string;
      values?: string[];
    };
  };
}

export interface AppEditor extends Identifiable {
  metadata?: {
    namespace?: string;
    name?: string;
    creationTimestamp?: string;
    generation?: number;
  };

  spec?: {
    dryRun?: boolean;
    chart?: {
      chartGroupName?: string;
      chartName?: string;
      chartVersion?: string;
      tenantID?: string;
    };
    tenantID?: string;
    name?: string;
    targetCluster?: string;
    type?: string;
    values?: {
      rawValues?: string;
      rawValuesType?: string;
      values?: string[];
    };
  };

  status?: {
    observedGeneration?: number;
  };

  /** 是否正在编辑 */
  v_editing?: boolean;
}

export interface AppResource extends Identifiable {
  metadata?: {
    namespace?: string;
    name?: string;
  };
  spec?: {
    type?: string;
    tenantID?: string;
    name?: string;
    resources?: any;
    targetCluster?: string;
  };
}

export interface AppResourceFilter {
  cluster?: string;
  namespace?: string;
  name?: string;
}

export interface Resource extends Identifiable {
  metadata?: {
    namespace?: string;
    name?: string;
  };
  kind?: string;
  cluster?: string;
  yaml?: string;
  object?: any;
}

export interface ResourceList extends Identifiable {
  resources?: Map<string, Resource[]>;
}

export interface History extends Identifiable {
  revision?: number;
  updated?: string;
  status?: string;
  chart?: string;
  appVersion?: string;
  description?: string;
  manifest?: string;
  involvedObject?: App;
}

export interface HistoryList extends Identifiable {
  histories?: History[];
}

export interface AppHistory extends Identifiable {
  metadata?: {
    namespace?: string;
    name?: string;
  };
  spec?: {
    type?: string;
    tenantID?: string;
    name?: string;
    targetCluster?: string;
    histories?: History[];
  };
}

export interface AppHistoryFilter {
  cluster?: string;
  namespace?: string;
  name?: string;
}
