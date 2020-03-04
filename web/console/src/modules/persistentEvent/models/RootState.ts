import { WorkflowState, FFListModel, FetcherState, QueryState } from '@tencent/ff-redux';
import { RecordSet } from '@tencent/qcloud-lib';
import { Region, RegionFilter, ResourceInfo, Resource, ResourceFilter } from '../../common';
import { RouteState } from '../../../../helpers';
import { PeEdit, CreateResource } from './';

type PeModifyWorkflow = WorkflowState<CreateResource, number>;

export interface RootState {
  /** 路由 */
  route?: RouteState;

  /** 地域列表 */
  region?: FFListModel<Region, RegionFilter>;

  /** 集群列表 */
  cluster?: FFListModel<Resource, ResourceFilter>;

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
