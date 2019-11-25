import { combineReducers } from 'redux';
import { reduceToPayload, RecordSet } from '@tencent/qcloud-lib';
import * as ActionType from '../constants/ActionType';
import { initValidator } from '../../common/models';
import { initSecretData } from '../constants/initState';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { Namespace } from '../models';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';

const TempReducer = combineReducers({
  name: reduceToPayload(ActionType.Sec_Name, ''),

  v_name: reduceToPayload(ActionType.SecV_Name, initValidator),

  nsList: generateFetcherReducer<RecordSet<Namespace>>({
    actionType: ActionType.Sec_FetchNsList,
    initialData: {
      recordCount: 0,
      records: [] as Namespace[]
    }
  }),

  nsQuery: generateQueryReducer({
    actionType: ActionType.Sec_QueryNsList
  }),

  secretType: reduceToPayload(ActionType.Sec_SecretType, 'Opaque'),

  data: reduceToPayload(ActionType.Sec_UpdateData, [initSecretData]),

  nsType: reduceToPayload(ActionType.Sec_NsType, 'specific'),

  nsListSelection: reduceToPayload(ActionType.Sec_NamespaceSelection, []),

  domain: reduceToPayload(ActionType.Sec_Domain, ''),

  v_domain: reduceToPayload(ActionType.SecV_Domain, initValidator),

  username: reduceToPayload(ActionType.Sec_Username, ''),

  v_username: reduceToPayload(ActionType.SecV_Username, initValidator),

  password: reduceToPayload(ActionType.Sec_Password, ''),

  v_password: reduceToPayload(ActionType.SecV_Password, initValidator)
});

export const SecretEditReducer = (state, action) => {
  let newState = state;
  // 销毁创建 Secret 界面
  if (action.type === ActionType.ClearSecretEdit) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};
