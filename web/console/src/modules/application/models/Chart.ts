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
    pullCount?: number;
    versions?: ChartVersion[];
  };

  // custom: store last version data
  lastVersion?: ChartVersion;
  sortedVersions?: ChartVersion[];
  projectID?: string;
}

export interface ChartVersion {
  chartSize?: number;
  description?: string;
  timeCreated?: string;
  version?: string;
  icon?: string;
  appVersion?: string;
}

export interface ChartFilter {
  namespace?: string;
  repoType?: string;
  projectID?: string;
}

export interface ChartInfo {
  metadata?: {
    name?: string;
    namespace?: string;
  };

  spec: {
    files: {
      [props: string]: string;
    };
    values: {
      [props: string]: string;
    };
  };
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
