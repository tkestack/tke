import { Identifiable } from '@tencent/ff-redux';

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

  value?: any;

  text?: any;

  /** other */
  [props: string]: any;
}

export interface ResourceFilter {
  /** 命名空间 */
  namespace?: string;

  /** 集群id */
  clusterId?: string;

  /** 集群日志组件 */
  logAgentName?: string;

  /** 地域id */
  regionId?: number;

  /** name */
  specificName?: string;
}
