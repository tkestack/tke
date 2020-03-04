import { WorkflowState } from '@tencent/ff-redux';
import { RouteState } from '../../../../helpers/Router';
import { Resource, ResourceFilter } from './Resource';
import { FFListModel } from '@tencent/ff-redux';

type ResourceOpWorkflow = WorkflowState<Resource, {}>;

export interface RootState {
  /**
   * 路由
   */
  route?: RouteState;

  channel?: FFListModel<Resource, ResourceFilter>;
  template?: FFListModel<Resource, ResourceFilter>;
  receiver?: FFListModel<Resource, ResourceFilter>;
  receiverGroup?: FFListModel<Resource, ResourceFilter>;
  resourceDeleteWorkflow?: ResourceOpWorkflow;
  modifyResourceFlow?: ResourceOpWorkflow;
  /** 是否为国际版 */
  isI18n?: boolean;
}
