import { FFListModel, WorkflowState } from '@tencent/ff-redux';

import { RouteState } from '../../../../helpers';
import { CreateResource, Region, RegionFilter, Resource, ResourceFilter } from '../../common';
import { Addon } from './';
import { AddonEdit } from './AddonEdit';

type ResourceModifyWorkflow = WorkflowState<CreateResource, number>;

export interface RootState {
  /** 路由 */
  route?: RouteState;

  /** 地域列表 */
  region?: FFListModel<Region, RegionFilter>;

  /** 集群列表 */
  cluster?: FFListModel<Resource, ResourceFilter>;

  /** 集群的版本 */
  clusterVersion?: string;

  /** 集群的下的addon列表 */
  openAddon?: FFListModel<Resource, ResourceFilter>;

  /** 所有的add的列表 */
  addon?: FFListModel<Addon, ResourceFilter>;

  /** 开通addon组件 */
  editAddon?: AddonEdit;

  /** 创建resource资源的操作流程 */
  modifyResourceFlow?: ResourceModifyWorkflow;

  /** 创建多种resource资源的操作流程 */
  applyResourceFlow?: ResourceModifyWorkflow;

  /** 删除resource资源的操作流程 */
  deleteResourceFlow?: ResourceModifyWorkflow;
}
