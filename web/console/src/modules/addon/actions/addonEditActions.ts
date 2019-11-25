import * as ActionType from '../constants/ActionType';
import { ReduxAction } from '@tencent/qcloud-lib';
import { peEditActions } from './peEditActions';

export const addonEditActions = {
  pe: peEditActions,

  /** 需要开通的扩展组件的名称 */
  selectAddonName: (name: string) => {
    return async (dispatch: Redux.Dispatch) => {
      dispatch({
        type: ActionType.AddonName,
        payload: name
      });
    };
  },

  /** 清除开通addon的相关信息 */
  clearCreateAddon: (): ReduxAction<void> => {
    return { type: ActionType.ClearAddonEdit };
  }
};
