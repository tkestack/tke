import { FetcherState, FetchState } from '@tencent/qcloud-redux-fetcher';
import { QueryState } from '@tencent/qcloud-redux-query';
import { WorkflowState } from '@tencent/qcloud-redux-workflow';
import { RecordSet } from '@tencent/qcloud-lib';
import { Region, RegionFilter, ResourceInfo, Resource, ResourceFilter } from '../../common';
import { RouteState } from '../../../../helpers';
import { PeEdit, CreateResource } from './';
import { ListModel } from '@tencent/redux-list';

type PeModifyWorkflow = WorkflowState<CreateResource, number>;

export interface RootState {
  /** 路由 */
  route?: RouteState;

  /** 地域列表 */
  region?: ListModel<Region, RegionFilter>;

  /** 集群列表 */
  cluster?: ListModel<Resource, ResourceFilter>;

  /** PersistentEvent的列表 */
  peList?: FetcherState<RecordSet<Resource>>;

  /** PersistentEvent的查询 */
  peQuery?: QueryState<ResourceFilter>;

  /** PersistentEvent的选择 */
  peSelection?: Resource[];

  /** peEdit */
  peEdit?: PeEdit;

  /** resourceInfo */
  resourceInfo?: ResourceInfo;

  /** 设置持久化事件 创建的操作流 */
  modifyPeFlow?: PeModifyWorkflow;

  /** 删除持久化事件 Delete的操作流 */
  deletePeFlow?: PeModifyWorkflow;
}

export interface FetchPeList {
  /** clusterId */
  clusterId?: string;

  /** peName */
  peName?: string;
}
