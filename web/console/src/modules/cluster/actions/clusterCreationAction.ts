import { deepClone, ReduxAction, uuid } from '@tencent/ff-redux';

import * as ActionType from '../constants/ActionType';
import { initClusterCreationState } from '../constants/initState';
import { RootState } from '../models';

type GetState = () => RootState;

export const clusterCreationAction = {
  /** 更新cluser的名称 */
  updateClusterCreationState: obj => {
    return (dispatch: Redux.Dispatch, getState: GetState) => {
      let { clusterCreationState } = getState();
      dispatch({
        type: ActionType.UpdateclusterCreationState,
        payload: Object.assign({}, clusterCreationState, obj)
      });
    };
  },

  /** 离开创建页面，清除 Creation当中的内容 */
  clearClusterCreationState: (): ReduxAction<any> => {
    return {
      type: ActionType.UpdateclusterCreationState,
      payload: initClusterCreationState
    };
  }
};
