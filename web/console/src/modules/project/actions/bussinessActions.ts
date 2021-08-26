/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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

import { UserManagedProject, UserManagedProjectFilter } from './../models/Project';
import { User, UserFilter, UserInfo } from './../models/User';
import { FFReduxActionName } from './../constants/Config';
import { createFFListActions, extend, createFFObjectActions } from '@tencent/ff-redux';
import * as ActionType from '../constants/ActionType';
import { Cluster, ClusterFilter, RootState } from '../models';
import * as WebAPI from '../WebAPI';

type GetState = () => RootState;

/** 集群列表的Actions */
const FFModelUserManagedProjectActions = createFFListActions<UserManagedProject, UserManagedProjectFilter>({
  actionName: FFReduxActionName.UserManagedProjects,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchUserManagedProjects(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().userManagedProjects;
  },
  onFinish: (record, dispatch, getState: GetState) => {}
});

const FFObjectNamespaceCertInfoActions = createFFObjectActions<UserInfo, string>({
  actionName: FFReduxActionName.UserInfo,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchUserId(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().userInfo;
  },
  onFinish: (record, dispatch, getState: GetState) => {
    dispatch(FFModelUserManagedProjectActions.applyFilter({ userId: record.data.uid }));
  }
});

export const bussinessActions = {
  userManagedProject: FFModelUserManagedProjectActions,

  userInfo: FFObjectNamespaceCertInfoActions,

  initPlatformType: (platformType: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.PlatformType,
        payload: platformType
      });
    };
  }
};
