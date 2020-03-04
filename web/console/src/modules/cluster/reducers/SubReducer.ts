import { combineReducers } from 'redux';

import { generateWorkflowReducer, RecordSet, reduceToPayload } from '@tencent/ff-redux';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';

import * as ActionType from '../constants/ActionType';
import { initAllcationRatioEdition } from '../constants/initState';
import { SubRouter } from '../models';
import { ComputerReducer } from './ComputerReducer';
import { ConfigMapEditReducer } from './ConfigMapEditReducer';
import { DetailResourceReducer } from './DetailResourceReducer';
import { LbcfEditReducer } from './LbcfEditReducer';
import { NamespaceEditReducer } from './NamespaceEditReducer';
import { ResourceDetailReducer } from './ResourceDetailReducer';
import { ResourceEventReducer } from './ResourceEventReducer';
import { ResourceLogReducer } from './ResourceLogReducer';
import { ResourceReducer } from './ResourceReducer';
import { SecretEditReducer } from './SecretEditReducer';
import { ServiceEditReducer } from './ServiceEditReducer';
import { WorkloadEditReducer } from './WorkloadEditReducer';

const TempReducer = combineReducers({
  computerState: ComputerReducer,

  clusterAllocationRatioEdition: reduceToPayload(
    ActionType.UpdateClusterAllocationRatioEdition,
    initAllcationRatioEdition
  ),

  updateClusterAllocationRatio: generateWorkflowReducer({
    actionType: ActionType.UpdateClusterAllocationRatio
  }),

  subRouterQuery: generateQueryReducer({
    actionType: ActionType.QuerySubRouterList
  }),

  subRouterList: generateFetcherReducer<RecordSet<SubRouter>>({
    actionType: ActionType.FetchSubRouterList,
    initialData: {
      recordCount: 0,
      records: [] as SubRouter[]
    }
  }),

  mode: reduceToPayload(ActionType.SelectMode, 'list'),

  applyResourceFlow: generateWorkflowReducer({
    actionType: ActionType.ApplyResource
  }),

  applyDifferentInterfaceResourceFlow: generateWorkflowReducer({
    actionType: ActionType.ApplyDifferentInterfaceResource
  }),

  modifyResourceFlow: generateWorkflowReducer({
    actionType: ActionType.ModifyResource
  }),

  modifyMultiResourceWorkflow: generateWorkflowReducer({
    actionType: ActionType.ModifyMultiResource
  }),

  deleteResourceFlow: generateWorkflowReducer({
    actionType: ActionType.DeleteResource
  }),

  updateResourcePart: generateWorkflowReducer({
    actionType: ActionType.UpdateResourcePart
  }),

  updateMultiResource: generateWorkflowReducer({
    actionType: ActionType.UpdateMultiResource
  }),

  resourceName: reduceToPayload(ActionType.InitResourceName, ''),

  resourceInfo: reduceToPayload(ActionType.InitResourceInfo, {}),

  detailResourceOption: DetailResourceReducer,

  resourceOption: ResourceReducer,

  resourceDetailState: ResourceDetailReducer,

  serviceEdit: ServiceEditReducer,

  namespaceEdit: NamespaceEditReducer,

  workloadEdit: WorkloadEditReducer,

  secretEdit: SecretEditReducer,

  cmEdit: ConfigMapEditReducer,

  lbcfEdit: LbcfEditReducer,

  resourceLogOption: ResourceLogReducer,

  resourceEventOption: ResourceEventReducer,

  isNeedFetchNamespace: reduceToPayload(ActionType.IsNeedFetchNamespace, true),

  isNeedExistedLb: reduceToPayload(ActionType.IsNeedExistedLb, false),
  addons: reduceToPayload(ActionType.FetchClusterAddons, {})
});

export const SubReducer = (state, action) => {
  let newState = state;
  // 销毁详情页面
  if (action.type === ActionType.ClearSubRoot) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};
