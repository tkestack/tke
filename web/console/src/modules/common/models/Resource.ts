import { Identifiable } from '@tencent/qcloud-lib';

export interface Resource extends Identifiable {
  /** metadata */
  metadata: {
    [props: string]: any;
  };

  /** spec */
  spec: {
    [props: string]: any;
  };

  /** data */
  data?: {
    [props: string]: any;
  };

  /** status */
  status: {
    [props: string]: any;
  };

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

  /** name */
  specificName?: string;
}
