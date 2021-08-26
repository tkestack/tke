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

import { ReduxAction, extend } from '@tencent/ff-redux';
import { RootState, PolicyInfoFilter, PolicyEditor, Policy } from '../../models';
import * as ActionTypes from '../../constants/ActionTypes';
import * as WebAPI from '../../WebAPI';
import { initPolicyEditorState } from '../../constants/initState';
import { router } from '../../router';
type GetState = () => RootState;

const restActions = {

  fetchPolicy: (filter: PolicyInfoFilter) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let response = await WebAPI.fetchPolicy(filter);
      let editor: PolicyEditor = response;
      dispatch({
        type: ActionTypes.UpdatePolicyEditorState,
        payload: editor
      });
    };
  },

  /** 更新状态 */
  updateEditorState: (obj) => {
    return (dispatch: Redux.Dispatch, getState: GetState) => {
      let { policyEditor } = getState();
      dispatch({
        type: ActionTypes.UpdatePolicyEditorState,
        payload: Object.assign({}, policyEditor, obj)
      });
    };
  },

  /** 离开更新页面，清除Editor当中的内容 */
  clearEditorState: (): ReduxAction<any> => {
    return {
      type: ActionTypes.UpdatePolicyEditorState,
      payload: initPolicyEditorState
    };
  }
};
export const detailActions = extend({}, restActions);