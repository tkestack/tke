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

import { FFObjectModel } from './../../../../lib/ff-redux/src/object/Model';
import { FetcherState } from '@tencent/ff-redux';

import { RegionCluster } from '../../common/models';

export interface RootState {
  clusterOverview?: FFObjectModel<ClusterOverview, ClusterOverviewFilter>;
}

export interface ClusterOverviewFilter {}

export interface ClusterOverview {
  clusterCount: number;
  clusterAbnormal: number;
  nodeCount: number;
  nodeAbnormal: number;
  workloadCount: number;
  workloadAbnormal: number;
  projectCount: number;
  projectAbnormal: number;
  clusters: ClusterDetail[];
}

export interface ClusterDetail {
  clusterID: string;
  clusterPhase: string;
  nodeCount: number;
  nodeAbnormal: number;
  workloadCount: number;
  workloadAbnormal: number;
  cpuUsage: string;
  cpuRequest: number;
  cpuLimit: number;
  cpuCapacity: number;
  cpuAllocatable: number;
  cpuRequestRate: string;
  cpuAllocatableRate: string;
  memUsage: string;
  memRequest: number;
  memLimit: number;
  memCapacity: number;
  memAllocatable: number;
  memRequestRate: string;
  memAllocatableRate: string;
  schedulerHealthy: boolean;
  controllerManagerHealthy: boolean;
  etcdHealthy: boolean;
  clusterDisplayName?: string;
}
