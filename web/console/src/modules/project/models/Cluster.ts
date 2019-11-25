import { Identifiable } from '@tencent/qcloud-lib';

export interface Cluster extends Identifiable {
  /** 集群Id */
  clusterId?: string | number;

  /** 名称 */
  clusterName?: string;
}

export interface ClusterFilter {
  regionId?: string | number;
}
