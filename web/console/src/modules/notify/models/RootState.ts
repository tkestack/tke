import { WorkflowState } from '@tencent/qcloud-redux-workflow';
import { RouteState } from '../../../../helpers/Router';
import { Resource, ResourceFilter } from './Resource';
import { ListModel } from '@tencent/redux-list';

type ResourceOpWorkflow = WorkflowState<Resource, {}>;

export interface RootState {
  /**
   * 路由
   */
  route?: RouteState;

  channel?: ListModel<Resource, ResourceFilter>;
  template?: ListModel<Resource, ResourceFilter>;
  receiver?: ListModel<Resource, ResourceFilter>;
  receiverGroup?: ListModel<Resource, ResourceFilter>;
  resourceDeleteWorkflow?: ResourceOpWorkflow;
  modifyResourceFlow?: ResourceOpWorkflow;
  /** 是否为国际版 */
  isI18n?: boolean;
}
