import { ReduxAction, uuid } from '@tencent/qcloud-lib';
import * as ActionType from '../constants/ActionType';
import { RootState, Variable, initVariable } from '../models';
import { cloneDeep } from '../../common/utils';

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
