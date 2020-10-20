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
