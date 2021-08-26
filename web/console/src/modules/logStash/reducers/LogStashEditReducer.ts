/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

import { combineReducers } from 'redux';

import { RecordSet, reduceToPayload } from '@tencent/ff-redux';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';

import { initValidator, Namespace, NamespaceFilter } from '../../common/models';
import * as ActionType from '../constants/ActionType';
import {
    initContainerFilePath, initContainerFileWorkloadType, initContainerInputOption,
    initResourceTarget
} from '../constants/initState';
import {
    Ckafka, CkafkaFilter, Cls, ClsTopic, CTopic, CTopicFilter, Pod, Resource
} from '../models';

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

  esUsername: reduceToPayload(ActionType.EsUsername, ''),

  esPassword: reduceToPayload(ActionType.EsPassword, ''),

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
