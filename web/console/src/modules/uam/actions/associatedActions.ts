import { extend, ReduxAction } from '@tencent/qcloud-lib';
import { generateFetcherActionCreator, FetchOptions, FetchState } from '@tencent/qcloud-redux-fetcher';
import { generateWorkflowActionCreator, OperationResult, OperationTrigger } from '@tencent/qcloud-redux-workflow';
import { generateQueryActionCreator, QueryState } from '@tencent/qcloud-redux-query';
import { RootState, User, UserFilter } from '../models';
import { ComputerFilter } from '../../cluster/models';
import * as ActionTypes from '../constants/ActionTypes';
import * as WebAPI from '../WebAPI';
import { router } from '../router';
import { createListAction } from '@tencent/redux-list';
type GetState = () => RootState;
export interface AssociatedUserParams {
  id: string;
  userNames: string[];
}

const associateUser = generateWorkflowActionCreator<AssociatedUserParams, void>({
  actionType: ActionTypes.AddAssociatedUser,
  workflowStateLocator: (state: RootState) => state.addAssociatedUser,
  operationExecutor: WebAPI.associateUser,
  after: {
    [OperationTrigger.Done]: (dispatch, getState) => {
      const { route } = getState();
      let urlParam = router.resolve(route);
      const { sub } = urlParam;
      if (sub) {
        dispatch(FFModelAssociatedUsersActions.applyFilter({}));
      }
    }
  }
});

/**
 * 获取策略关联的用户
 */
const FFModelAssociatedUsersActions = createListAction<User, UserFilter>({
  actionName: ActionTypes.GetStrategyAssociatedUsers,
  fetcher: async (query, getState: GetState) => {
    let { route } = getState();
    const urlParams = router.resolve(route);
    let result = await WebAPI.fetchStrategyAssociatedUsers(query.filter.search || urlParams.sub);
    return result;
  },
  getRecord: (getState: GetState) => {
    return getState().associatedUsersList;
  }
});

/**
 * 删除策略关联的用户
 */
const removeAssociatedUser = generateWorkflowActionCreator<any, void>({
  actionType: ActionTypes.RemoveAssociatedUser,
  workflowStateLocator: (state: RootState) => state.removeAssociatedUser,
  operationExecutor: WebAPI.removeAssociatedUser,
  after: {
    [OperationTrigger.Done]: dispatch => {
      dispatch(FFModelAssociatedUsersActions.applyFilter({}));
    }
  }
});
const restActions = { associateUser, removeAssociatedUser };

export const associateActions = extend({}, FFModelAssociatedUsersActions, restActions);
