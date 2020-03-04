import { Identifiable, RecordSet } from '@tencent/ff-redux';

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

  specificName?: string;

  /** type */
  workloadType?: string;
}
