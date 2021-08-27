/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */
import {
    createFFListActions, extend, generateWorkflowActionCreator, OperationTrigger
} from '@tencent/ff-redux';

import * as ActionTypes from '../constants/ActionTypes';
import { RootState, User, UserFilter } from '../models';
import { router } from '../router';
import * as WebAPI from '../WebAPI';

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
const FFModelAssociatedUsersActions = createFFListActions<User, UserFilter>({
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
