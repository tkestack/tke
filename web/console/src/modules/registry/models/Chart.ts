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

export interface Chart extends Identifiable {
  metadata?: {
    creationTimestamp?: string;
    name?: string;
    namespace?: string;
  };

  spec?: {
    chartGroupName?: string;
    displayName?: string;
    name?: string;
    tenantID?: string;
    visibility?: string;
  };

  status?: {
    locked?: boolean;
    pullCount?: number;
    versions?: ChartVersion[];
  };

  // custom: store last version data
  lastVersion?: ChartVersion;
  sortedVersions?: ChartVersion[];
}

export interface ChartVersion extends Identifiable {
  chartSize?: number;
  description?: string;
  timeCreated?: string;
  version?: string;
  icon?: string;
  appVersion?: string;
}

export interface RemovedChartVersions {
  versions?: RemovedChartVersion[];
}

export interface RemovedChartVersion {
  name?: string;
  namespace?: string;
  version?: string;
}

export interface ChartFilter {
  repoType?: string;
  projectID?: string;
}

export interface ChartDetailFilter {
  namespace: string;
  name: string;
  projectID: string;
}

export interface ChartVersionFilter {
  chartGroupName?: string;
  chartName?: string;
  chartVersion?: string;
  chartDetailFilter?: ChartDetailFilter;
}

export interface ChartEditor extends Identifiable {
  metadata?: {
    creationTimestamp?: string;
    name?: string;
    namespace?: string;
  };

  spec?: {
    chartGroupName?: string;
    displayName?: string;
    name?: string;
    tenantID?: string;
    visibility?: string;
  };

  status?: {
    pullCount?: number;
    versions?: ChartVersion[];
  };

  /** 是否正在编辑 */
  v_editing?: boolean;
  sortedVersions?: ChartVersion[];
  selectedVersion?: ChartVersion;
}

export interface ChartInfo {
  metadata?: {
    name?: string;
    namespace?: string;
  };

  spec: {
    /** readme */
    readme: {
      [props: string]: string;
    };
    /** values */
    values: {
      [props: string]: string;
    };
    /** files */
    rawFiles: {
      [props: string]: string;
    };
  };

  fileTree?: ChartTreeFile;
}

export interface ChartInfoFilter {
  cluster: string;
  namespace: string;
  metadata: {
    namespace: string;
    name: string;
  };
  chartVersion: string;
  projectID: string;
}

export interface ChartTreeFile {
  name: string;
  data: string;
  fullPath: string;
  children?: ChartTreeFile[];
}

/** 以下废弃 */
export interface ChartCreation extends Identifiable {
  /** 描述 */
  displayName?: 'string';
  /** 仓库名称 */
  name?: 'string';
  v_name?: Validation;
  /** 公开或私有 */
  visibility?: 'Public' | 'Private';
}

export interface ChartIns extends Identifiable {
  apiVersion?: string;

  kind?: boolean;

  /** 元数据 */
  metadata?: {
    annotations?: any;
    clusterName?: string;
    creationTimestamp?: string;
    deletionGracePeriodSeconds?: number;
    deletionTimestamp?: string;
    finalizers?: string[];
    generateName?: string;
    generation?: string;
    labels?: any;
    managedFields?: any[];
    /** 键值 name */
    name?: string;
    namespace?: string;
    ownerReferences?: any;
    resourceVersion?: string;
    selfLink?: string;
    uid?: string;
  };

  spec?: {
    displayName?: string;
    name?: string;
    tenantID?: string;
    visibility?: string;
  };

  status?: {
    locked?: boolean;
    pullCount?: number;
  };
}

export interface ChartInsFilter {
  chartgroup?: string;
}
