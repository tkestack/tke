import { Identifiable, FFListModel } from '@tencent/ff-redux';

export interface ResourceOption {
  /** 具体的resource列表 */
  ffResourceList?: FFListModel<Resource, ResourceFilter>;

  /** resource的多选选择 */
  resourceMultipleSelection?: Resource[];

  /** resourceDeleteSelection */
  resourceDeleteSelection?: Resource[];
}

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

  /** name */
  specificName?: string;

  meshId?: string;
}

export interface DifferentInterfaceResourceOperation {
  query?: {
    [props: string]: any;
  };
  extraResource?: string;
}
