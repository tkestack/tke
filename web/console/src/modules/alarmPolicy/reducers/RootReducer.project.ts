import { combineReducers } from 'redux';

import {
    createFFListReducer, generateWorkflowReducer, RecordSet, reduceToPayload
} from '@tencent/ff-redux';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';

import { FFReduxActionName } from '../../cluster/constants/Config';
import { Cluster } from '../../common';
import * as ActionType from '../constants/ActionType';
import { ClusterFilter, Group, Namespace, Region, Resource } from '../models';
import { router } from '../router';
import { AlarmPolicyEditReducer } from './AlarmPolicyEditReducer';

export const RootReducer = combineReducers({
  route: router.getReducer(),

  regionQuery: generateQueryReducer({
    actionType: ActionType.QueryRegion
  }),

  regionList: generateFetcherReducer<RecordSet<Region>>({
    actionType: ActionType.FetchRegion,
    initialData: {
      recordCount: 0,
      records: [] as Region[]
    }
  }),

  regionSelection: reduceToPayload(ActionType.SelectRegion, { value: 1 }),

  cluster: createFFListReducer<Cluster, ClusterFilter>(FFReduxActionName.CLUSTER),
  clusterVersion: reduceToPayload(ActionType.InitClusterVersion, '1.16'),

  /**当前集群命名空间 */
  namespaceList: generateFetcherReducer({
    actionType: ActionType.FetchNamespaceList,
    initialData: {
      recordCount: 0,
      records: [] as Namespace[]
    }
  }),
  namespaceQuery: generateQueryReducer({
    actionType: ActionType.QueryNamespaceList
  }),

  /**当前命名空间下pod列表 */
  workloadList: generateFetcherReducer({
    actionType: ActionType.FetchWorkloadList,
    initialData: {
      recordCount: 0,
      records: [] as Resource[]
    }
  }),

  workloadQuery: generateQueryReducer({
    actionType: ActionType.QueryWorkloadList
  }),

  alarmPolicy: createFFListReducer('AlarmPolicy'),

  userList: createFFListReducer('UserList'),

  /** 当前新建告警 */
  alarmPolicyEdition: AlarmPolicyEditReducer,

  /** 创建告警workflow */
  alarmPolicyCreateWorkflow: generateWorkflowReducer({ actionType: ActionType.CreateAlarmPolicy }),

  /** 复制告警workflow */
  alarmPolicyUpdateWorkflow: generateWorkflowReducer({ actionType: ActionType.UpdateAlarmPolicy }),

  /** 删除告警workflow */
  alarmPolicyDeleteWorkflow: generateWorkflowReducer({ actionType: ActionType.DeleteAlarmPolicy }),

  /**详情 */
  alarmPolicyDetail: reduceToPayload(ActionType.FetchalarmPolicyDetail, {}),

  channel: createFFListReducer<Resource, ClusterFilter>('channel'),
  template: createFFListReducer<Resource, ClusterFilter>('template'),
  receiver: createFFListReducer<Resource, ClusterFilter>('receiver'),
  receiverGroup: createFFListReducer<Resource, ClusterFilter>('receiverGroup'),

  /**
   * 判断是否为国际版
   */
  isI18n: reduceToPayload(ActionType.isI18n, false),

  namespaceSelection: reduceToPayload(ActionType.SelectNamespace, ''),
  projectList: reduceToPayload(ActionType.InitProjectList, []),
  projectSelection: reduceToPayload(ActionType.ProjectSelection, ''),
  projectNamespaceQuery: generateQueryReducer({
    actionType: ActionType.QueryProjectNamespace
  }),
  projectNamespaceList: generateFetcherReducer<RecordSet<Resource>>({
    actionType: ActionType.FetchProjectNamespace,
    initialData: {
      recordCount: 0,
      records: [] as Resource[]
    }
  })
});
