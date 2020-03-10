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