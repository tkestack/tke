import * as ActionType from '../constants/ActionType';
import { RootState } from '../models';
import { ReduxAction } from '@tencent/qcloud-lib';
import { includes } from '../../common';

type GetState = () => RootState;

export const peEditActions = {
  /** 输入 es的地址 */
  inputEsAddress: (address: string): ReduxAction<string> => {
    return {
      type: ActionType.EsAddress,
      payload: address
    };
  },

  /** 输入当前的索引 */
  inputIndexName: (indexName: string): ReduxAction<string> => {
    return {
      type: ActionType.IndexName,
      payload: indexName
    };
  }
};
