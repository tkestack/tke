import { ReduxAction } from '@tencent/ff-redux';

import { includes } from '../../common';
import * as ActionType from '../constants/ActionType';
import { RootState } from '../models';

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
