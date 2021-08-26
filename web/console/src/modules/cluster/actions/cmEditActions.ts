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

import { ReduxAction, uuid } from '@tencent/ff-redux';

import { cloneDeep } from '../../common/utils';
import * as ActionType from '../constants/ActionType';
import { initVariable, RootState, Variable } from '../models';

type GetState = () => RootState;

export const cmEditActions = {
  /** 输入名称 */
  inputCMName: (name: string): ReduxAction<string> => {
    return {
      type: ActionType.CM_Name,
      payload: name
    };
  },

  /** 选择命名空间 */
  selectNamespace: (namespace: string): ReduxAction<string> => {
    return {
      type: ActionType.CM_Namespace,
      payload: namespace
    };
  },

  /** 新增变量 */
  addVariable: () => {
    return async (dispatch, getState: GetState) => {
      let variables = cloneDeep(getState().subRoot.cmEdit.variables);

      variables.push(Object.assign({}, initVariable, { id: uuid() }));
      dispatch({
        type: ActionType.CM_AddVariable,
        payload: variables
      });
    };
  },

  /** 编辑变量 */
  eidtVariable: (id: string | number, obj: any) => {
    return async (dispatch, getState: GetState) => {
      let variables: Array<Variable> = cloneDeep(getState().subRoot.cmEdit.variables),
        vIndex = variables.findIndex(i => i.id === id);

      if (vIndex > -1) {
        variables[vIndex] = Object.assign(variables[vIndex], obj);
      }

      dispatch({
        type: ActionType.CM_EditVariable,
        payload: variables
      });
    };
  },

  /** 删除变量 */
  deleteVariable: (id: string | number) => {
    return async (dispatch, getState: GetState) => {
      let variables: Array<Variable> = cloneDeep(getState().subRoot.cmEdit.variables),
        vIndex = variables.findIndex(i => i.id === id);

      if (vIndex > -1) {
        variables.splice(vIndex, 1);
      }

      dispatch({
        type: ActionType.CM_DeleteVariable,
        payload: variables
      });
    };
  },

  /** 清除pv的编辑项 */
  clearConfigMapEdit: (): ReduxAction<any> => {
    return {
      type: ActionType.ClearConfigMapEdit
    };
  }
};
