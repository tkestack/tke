import { Identifiable } from '@tencent/qcloud-lib';

export interface Resource extends Identifiable {
  /** metadata */
  metadata?: any;

  /** spec */
  spec?: any;

  /** status */
  status?: any;

  /** data */
  data?: any;

  /** other */
  [props: string]: any;
}

export interface ResourceFilter {
  /** 命名空间 */
  namespace?: string;

  /** 集群id */
  clusterId?: string;

  /** 地域id */
  regionId?: number;

  /** workloadType */
  workloadType?: string;

  /** isCanFetchResourceList */
  isCanFetchResourceList?: boolean;
}

export interface ResourceTarget {
  isForContainerFile: boolean;
  isForContainerLogs: boolean;
}

export interface WorkLoadList {
  name: string;
  value: string;
}
