import { WorkflowState, OperationResult } from '@tencent/qcloud-redux-workflow';
import { User, UserFilter, Strategy, StrategyFilter, Category } from './index';
import { RouteState } from '../../../../helpers';
import { FetcherState } from '@tencent/qcloud-redux-fetcher';
import { RecordSet } from '@tencent/qcloud-lib';
import { ListModel } from '@tencent/redux-list';
import { ResourceFilter } from '@src/modules/common';

type userWorkflow = WorkflowState<User, any>;
type strategyWorkflow = WorkflowState<Strategy, any>;
type associateWorkflow = WorkflowState<{ id: string; userNames: [] }, any>;

export interface RootState {
  /** 用户信息 */
  userList?: ListModel<User, UserFilter>;
  addUserWorkflow?: userWorkflow;
  removeUserWorkflow?: userWorkflow;
  user?: User;
  filterUsers?: User[];
  getUser?: OperationResult<User>;
  updateUser?: FetcherState<RecordSet<any>>;
  userStrategyList?: ListModel<Strategy, ResourceFilter>;

  /** 策略相关 */
  strategyList?: ListModel<Strategy, StrategyFilter>;
  addStrategyWorkflow?: strategyWorkflow;
  removeStrategyWorkflow?: strategyWorkflow;
  associatedUsersList?: ListModel<User, UserFilter>;
  removeAssociatedUser?: associateWorkflow;
  addAssociatedUser?: associateWorkflow;
  getStrategy?: OperationResult<Strategy>;
  updateStrategy?: FetcherState<RecordSet<any>>;

  /** 类别 */
  categoryList?: FetcherState<RecordSet<Category>>;

  /** 路由 */
  route?: RouteState;
}
