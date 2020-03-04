import { WorkflowState, OperationResult } from '@tencent/ff-redux';
import { User, UserFilter, Strategy, StrategyFilter, Category } from './index';
import { RouteState } from '../../../../helpers';
import { RecordSet } from '@tencent/qcloud-lib';
import { FFListModel, FetcherState } from '@tencent/ff-redux';
import { ResourceFilter } from '@src/modules/common';

type userWorkflow = WorkflowState<User, any>;
type strategyWorkflow = WorkflowState<Strategy, any>;
type associateWorkflow = WorkflowState<{ id: string; userNames: [] }, any>;

export interface RootState {
  /** 用户信息 */
  userList?: FFListModel<User, UserFilter>;
  addUserWorkflow?: userWorkflow;
  removeUserWorkflow?: userWorkflow;
  user?: User;
  filterUsers?: User[];
  getUser?: OperationResult<User>;
  updateUser?: FetcherState<RecordSet<any>>;
  userStrategyList?: FFListModel<Strategy, ResourceFilter>;

  /** 策略相关 */
  strategyList?: FFListModel<Strategy, StrategyFilter>;
  addStrategyWorkflow?: strategyWorkflow;
  removeStrategyWorkflow?: strategyWorkflow;
  associatedUsersList?: FFListModel<User, UserFilter>;
  removeAssociatedUser?: associateWorkflow;
  addAssociatedUser?: associateWorkflow;
  getStrategy?: OperationResult<Strategy>;
  updateStrategy?: FetcherState<RecordSet<any>>;

  /** 类别 */
  categoryList?: FetcherState<RecordSet<Category>>;

  /** 路由 */
  route?: RouteState;
}
