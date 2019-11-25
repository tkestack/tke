import { ExternalTrafficPolicy, SessionAffinity } from './../constants/Config';
import { combineReducers } from 'redux';
import { reduceToPayload, RecordSet } from '@tencent/qcloud-lib';
import * as ActionType from '../constants/ActionType';
import { initValidator } from '../../common/models';
import { initPortsMap } from '../constants/initState';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { Resource } from '../models';

const TempReducer = combineReducers({
  serviceName: reduceToPayload(ActionType.S_ServiceName, ''),

  v_serviceName: reduceToPayload(ActionType.SV_ServiceName, initValidator),

  description: reduceToPayload(ActionType.S_Description, ''),

  v_description: reduceToPayload(ActionType.SV_Description, initValidator),

  namespace: reduceToPayload(ActionType.S_Namespace, ''),

  v_namespace: reduceToPayload(ActionType.SV_Namespace, initValidator),

  communicationType: reduceToPayload(ActionType.S_CommunicationType, 'ClusterIP'),

  portsMap: reduceToPayload(ActionType.S_UpdatePortsMap, [initPortsMap]),

  isOpenHeadless: reduceToPayload(ActionType.S_IsOpenHeadless, false),

  selector: reduceToPayload(ActionType.S_Selector, []),

  isShowWorkloadDialog: reduceToPayload(ActionType.S_IsShowWorkloadDialog, false),

  workloadType: reduceToPayload(ActionType.S_WorkloadType, 'deployment'),

  workloadQuery: generateQueryReducer({
    actionType: ActionType.S_QueryWorkloadList
  }),

  workloadList: generateFetcherReducer<RecordSet<Resource>>({
    actionType: ActionType.S_FetchWorkloadList,
    initialData: {
      recordCount: 0,
      records: [] as Resource[]
    }
  }),

  workloadSelection: reduceToPayload(ActionType.S_WorkloadSelection, []),

  externalTrafficPolicy: reduceToPayload(ActionType.S_ChooseExternalTrafficPolicy, ExternalTrafficPolicy.Cluster),

  sessionAffinity: reduceToPayload(ActionType.S_ChoosesessionAffinity, SessionAffinity.None),

  sessionAffinityTimeout: reduceToPayload(ActionType.S_InputsessionAffinityTimeout, 30),

  v_sessionAffinityTimeout: reduceToPayload(ActionType.SV_sessionAffinityTimeout, initValidator)
});

export const ServiceEditReducer = (state, action) => {
  let newState = state;
  // 销毁创建服务页面
  if (action.type === ActionType.ClearServiceEdit) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};
