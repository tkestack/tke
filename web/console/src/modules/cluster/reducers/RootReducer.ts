import { generateWorkflowReducer } from '@tencent/qcloud-redux-workflow';
import { createListReducer } from '@tencent/redux-list';
import { combineReducers } from 'redux';
import { reduceToPayload, RecordSet } from '@tencent/qcloud-lib';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';
import * as ActionType from '../constants/ActionType';
import { Namespace } from '../models';
import { Cluster } from '../../common/models';
import { SubReducer } from './SubReducer';
import { router } from '../router';
import { FFReduxActionName } from '../constants/Config';
import { initDialogState, initClusterCreationState } from '../constants/initState';
import { CreateICReducer } from './CreateICReducer';

export const RootReducer = combineReducers({
  route: router.getReducer(),

  region: createListReducer(FFReduxActionName.REGION),

  cluster: createListReducer(FFReduxActionName.CLUSTER),

  clustercredential: reduceToPayload(ActionType.FetchClustercredential, {
    name: '',
    clusterName: '',
    caCert: '',
    token: ''
  }),

  updateClusterToken: generateWorkflowReducer({
    actionType: ActionType.UpdateClusterToken
  }),

  clusterVersion: reduceToPayload(ActionType.ClusterVersion, '1.16'),

  clusterInfoQuery: generateQueryReducer({
    actionType: ActionType.QueryClusterInfo
  }),

  clusterInfoList: generateFetcherReducer<RecordSet<Cluster>>({
    actionType: ActionType.FetchClusterInfo,
    initialData: {
      recordCount: 0,
      records: [] as Cluster[]
    }
  }),

  namespaceList: generateFetcherReducer<RecordSet<Namespace>>({
    actionType: ActionType.FetchNamespaceList,
    initialData: {
      recordCount: 0,
      records: [] as Namespace[]
    }
  }),

  namespaceQuery: generateQueryReducer({
    actionType: ActionType.QueryNamespaceList
  }),

  namespaceSelection: reduceToPayload(ActionType.SelectNamespace, 'default'),

  subRoot: SubReducer,

  mode: reduceToPayload(ActionType.ChangeMode, 'create'),

  dialogState: reduceToPayload(ActionType.UpdateDialogState, initDialogState),

  deleteClusterFlow: generateWorkflowReducer({
    actionType: ActionType.DeleteCluster
  }),

  clusterCreationState: reduceToPayload(ActionType.UpdateclusterCreationState, initClusterCreationState),

  createClusterFlow: generateWorkflowReducer({
    actionType: ActionType.CreateCluster
  }),

  createIC: CreateICReducer,
  createICWorkflow: generateWorkflowReducer({
    actionType: ActionType.CreateIC
  }),

  modifyClusterName: generateWorkflowReducer({
    actionType: ActionType.ModifyClusterNameWorkflow
  })
});
