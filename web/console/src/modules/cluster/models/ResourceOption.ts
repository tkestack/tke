import { FetcherState } from '@tencent/qcloud-redux-fetcher';
import { QueryState } from '@tencent/qcloud-redux-query';
import { Identifiable, RecordSet } from '@tencent/qcloud-lib';

export interface ResourceOption {
  /** resource的查询 */
  resourceQuery?: QueryState<ResourceFilter>;

  /** resource的列表 */
  resourceList?: FetcherState<RecordSet<Resource>>;

  /** resource的选择 */
  resourceSelection?: Resource[];

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
