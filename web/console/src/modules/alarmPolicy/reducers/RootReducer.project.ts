import { AlarmPolicyEditReducer } from './AlarmPolicyEditReducer';
import { combineReducers } from 'redux';
import { reduceToPayload, RecordSet } from '@tencent/qcloud-lib';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';
import { generateWorkflowReducer } from '@tencent/qcloud-redux-workflow';
import { router } from '../router';
import * as ActionType from '../constants/ActionType';
import { Region, ClusterFilter, Group, Namespace, Resource } from '../models';
import { Cluster } from '../../common';
import { createListReducer } from '@tencent/redux-list';
import { FFReduxActionName } from '../../cluster/constants/Config';

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

  cluster: createListReducer<Cluster, ClusterFilter>(FFReduxActionName.CLUSTER),
  clusterVersion: reduceToPayload(ActionType.InitClusterVersion, '1.8'),

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

  alarmPolicy: createListReducer('AlarmPolicy'),

  userList: createListReducer('UserList'),

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

  channel: createListReducer<Resource, ClusterFilter>('channel'),
  template: createListReducer<Resource, ClusterFilter>('template'),
  receiver: createListReducer<Resource, ClusterFilter>('receiver'),
  receiverGroup: createListReducer<Resource, ClusterFilter>('receiverGroup'),

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
