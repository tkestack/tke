import { combineReducers } from 'redux';
import { reduceToPayload, RecordSet } from '@tencent/qcloud-lib';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import * as ActionType from '../constants/ActionType';
import { initValidator, Namespace, NamespaceFilter } from '../../common/models';
import {
  initContainerInputOption,
  initContainerFilePath,
  initContainerFileWorkloadType,
  initResourceTarget
} from '../constants/initState';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';
import { Ckafka, CkafkaFilter, CTopic, CTopicFilter, Cls, ClsTopic, Resource, Pod } from '../models';

const TempReducer = combineReducers({
  logStashName: reduceToPayload(ActionType.LogStashName, ''),

  v_logStashName: reduceToPayload(ActionType.V_LogStashName, initValidator),

  v_clusterSelection: reduceToPayload(ActionType.V_SelectClusterSelection, initValidator),

  logMode: reduceToPayload(ActionType.ChangeLogMode, 'container'),

  isSelectedAllNamespace: reduceToPayload(ActionType.IsSelectedAllNamespace, 'selectAll'),

  containerLogs: reduceToPayload(ActionType.UpdateContainerLogs, [initContainerInputOption]),

  nodeLogPath: reduceToPayload(ActionType.NodeLogPath, ''),

  v_nodeLogPath: reduceToPayload(ActionType.V_NodeLogPath, initValidator),

  metadatas: reduceToPayload(ActionType.UpdateMetadata, []),

  consumerMode: reduceToPayload(ActionType.ChangeConsumerMode, 'kafka'),

  addressIP: reduceToPayload(ActionType.AddressIP, ''),

  v_addressIP: reduceToPayload(ActionType.V_AddressIP, initValidator),

  addressPort: reduceToPayload(ActionType.AddressPort, ''),

  v_addressPort: reduceToPayload(ActionType.V_AddressPort, initValidator),

  topic: reduceToPayload(ActionType.Topic, ''),

  v_topic: reduceToPayload(ActionType.V_Topic, initValidator),

  esAddress: reduceToPayload(ActionType.EsAddress, ''),

  v_esAddress: reduceToPayload(ActionType.V_EsAddress, initValidator),

  indexName: reduceToPayload(ActionType.IndexName, ''),

  v_indexName: reduceToPayload(ActionType.V_IndexName, initValidator),

  resourceList: generateFetcherReducer<RecordSet<Resource>>({
    actionType: ActionType.FetchResourceList,
    initialData: {
      recordCount: 0,
      records: [] as Resource[]
    }
  }),

  resourceQuery: generateQueryReducer({
    actionType: ActionType.QueryResourceList
  }),

  resourceTarget: reduceToPayload(ActionType.UpdateResourceTarget, initResourceTarget),

  isFirstFetchResource: reduceToPayload(ActionType.isFirstFetchResource, true),

  containerFileNamespace: reduceToPayload(ActionType.SelectContainerFileNamespace, ''),
  v_containerFileNamespace: reduceToPayload(ActionType.V_ContainerFileNamespace, initValidator),

  containerFileWorkloadType: reduceToPayload(ActionType.SelectContainerFileWorkloadType, initContainerFileWorkloadType),
  v_containerFileWorkloadType: reduceToPayload(ActionType.V_ContainerFileWorkloadType, initValidator),

  containerFileWorkload: reduceToPayload(ActionType.SelectContainerFileWorkload, ''),
  v_containerFileWorkload: reduceToPayload(ActionType.V_ContainerFileWorkload, initValidator),

  containerFilePaths: reduceToPayload(ActionType.UpdateContainerFilePaths, [initContainerFilePath]),

  podListQuery: generateQueryReducer({
    actionType: ActionType.QueryPodList
  }),

  podList: generateFetcherReducer<RecordSet<Pod>>({
    actionType: ActionType.FetchPodList,
    initialData: {
      recordCount: 0,
      records: [] as Pod[]
    }
  }),

  containerFileWorkloadList: reduceToPayload(ActionType.UpdateContaierFileWorkloadList, [])
});

export const LogStashEditReducer = (state, action) => {
  let newState = state;
  // 销毁页面
  if (action.type === ActionType.ClearLogStashEdit) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};
