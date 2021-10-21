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
  extend,
  generateWorkflowActionCreator,
  OperationTrigger,
  isSuccessWorkflow,
  createFFListActions
} from '@tencent/ff-redux';
import { RootState, CommonUserAssociation, UserPlain } from '../../models';
import * as ActionTypes from '../../constants/ActionType';
import * as WebAPI from '../../WebAPI';
import { initCommonUserAssociationState } from '../../constants/initState';

type GetState = () => RootState;

/**
 * 列表操作
 */
const fetchUserActions = createFFListActions<UserPlain, void>({
  actionName: ActionTypes.UserPlainList,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchCommonUserList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().userPlainList;
  },
  onFinish: (record, dispatch: Redux.Dispatch) => {}
});

const restActions = {
  /** 选中用户，根据原始数据计算将添加的用户和将删除的用户 */
  selectUser: (users: UserPlain[]) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      /** 选中关联用户，则更新关联用户面板WorkflowDialog会用到的CommonUserAssociation状态数据 */
      /** 比对计算出新增和删除的用户，originUsers是指原先绑定的用户 */
      const { originUsers } = getState().commonUserAssociation;
      const getDifferenceSet = (arr1, arr2) => {
        let a1 = arr1.map(JSON.stringify);
        let a2 = arr2.map(JSON.stringify);
        return a1
          .concat(a2)
          .filter(v => !a1.includes(v) || !a2.includes(v))
          .map(JSON.parse);
      };
      let allUsers = users.concat(originUsers);
      let removeUsers = getDifferenceSet(users, allUsers);
      let addUsers = getDifferenceSet(originUsers, allUsers);
      dispatch({
        type: ActionTypes.UpdateCommonUserAssociation,
        payload: Object.assign({}, getState().commonUserAssociation, {
          users: users,
          addUsers: addUsers,
          removeUsers: removeUsers
        })
      });
    };
  },

  /** 清除用户关联状态数据 */
  clearUserAssociation: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionTypes.UpdateCommonUserAssociation,
        payload: initCommonUserAssociationState
      });
    };
  }
};

export const associateActions = extend(
  {},
  {
    userList: fetchUserActions
  },
  restActions
);
