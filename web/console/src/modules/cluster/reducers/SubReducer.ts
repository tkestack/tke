import { initAllcationRatioEdition } from './../constants/initState';
import { combineReducers } from 'redux';
import { reduceToPayload, RecordSet } from '@tencent/qcloud-lib';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';
import { generateWorkflowReducer } from '@tencent/ff-redux';
import { SubRouter } from '../models';
import * as ActionType from '../constants/ActionType';
import { ComputerReducer } from './ComputerReducer';
import { ResourceReducer } from './ResourceReducer';
import { ServiceEditReducer } from './ServiceEditReducer';
import { NamespaceEditReducer } from './NamespaceEditReducer';
import { ResourceDetailReducer } from './ResourceDetailReducer';
import { WorkloadEditReducer } from './WorkloadEditReducer';
import { ResourceLogReducer } from './ResourceLogReducer';
import { ResourceEventReducer } from './ResourceEventReducer';
import { SecretEditReducer } from './SecretEditReducer';
import { ConfigMapEditReducer } from './ConfigMapEditReducer';
import { LbcfEditReducer } from './LbcfEditReducer';
import { DetailResourceReducer } from './DetailResourceReducer';

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
