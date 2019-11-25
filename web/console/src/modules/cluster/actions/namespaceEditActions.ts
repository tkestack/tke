import { ReduxAction } from '@tencent/qcloud-lib';
import * as ActionType from '../constants/ActionType';
import { RootState } from '../models';

type GetState = () => RootState;

export const namespaceEditActions = {
  /** 更新namespace的名称 */
  inputNamespaceName: (name: string): ReduxAction<string> => {
    return {
      type: ActionType.N_Name,
      payload: name
    };
  },

  /** 更新namespace的描述 */
  inputNamespaceDesp: (desp: string): ReduxAction<string> => {
    return {
      type: ActionType.N_Description,
      payload: desp
    };
  },

  /** 离开创建页面，清除 namespaceEdit当中的内容 */
  clearNamespaceEdit: (): ReduxAction<any> => {
    return {
      type: ActionType.ClearNamespaceEdit
    };
  }
};
