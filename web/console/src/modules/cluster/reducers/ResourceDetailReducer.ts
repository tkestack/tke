import { reduceToPayload, RecordSet } from '@tencent/qcloud-lib';
import { generateFetcherReducer, FetcherState } from '@tencent/qcloud-redux-fetcher';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';
import { generateWorkflowReducer } from '@tencent/qcloud-redux-workflow';
import { combineReducers } from 'redux';
import * as ActionType from '../constants/ActionType';
import { Event, Replicaset, Pod, Resource, ResourceFilter } from '../models';
import { createListReducer } from '@tencent/redux-list';
import { FFReduxActionName } from '../constants/Config';

/** ==== start 日志的相关处理 ============ */
const logOptionReducer = combineReducers({
  podName: reduceToPayload(ActionType.PodName, ''),

  containerName: reduceToPayload(ActionType.ContainerName, ''),

  tailLines: reduceToPayload(ActionType.TailLines, '100'),

  isAutoRenew: reduceToPayload(ActionType.IsAutoRenewPodLog, false)
});
/** ==== start 日志的相关处理 ============ */

const TempReducer = combineReducers({
  yamlList: generateFetcherReducer<RecordSet<string>>({
    actionType: ActionType.FetchYaml,
    initialData: {
      recordCount: 0,
      records: [] as string[]
    }
  }),

  event: createListReducer<Event, ResourceFilter>(FFReduxActionName.DETAILEVENT),

  rsQuery: generateQueryReducer({
    actionType: ActionType.QueryRsList
  }),

  rsList: generateFetcherReducer<RecordSet<Replicaset>>({
    actionType: ActionType.FetchRsList,
    initialData: {
      recordCount: 0,
      records: [] as Replicaset[]
    }
  }),

  rollbackResourceFlow: generateWorkflowReducer({
    actionType: ActionType.RollBackResource
  }),

  removeTappPodFlow: generateWorkflowReducer({
    actionType: ActionType.RemoveTappPod
  }),

  rsSelection: reduceToPayload(ActionType.RsSelection, []),

  podQuery: generateQueryReducer({
    actionType: ActionType.QueryPodList
  }),

  podList: generateFetcherReducer<RecordSet<Pod>>({
    actionType: ActionType.FetchPodList,
    initialData: {
      recordCount: 0,
      records: [] as Pod[]
    }
  }),

  podFilterInNode: reduceToPayload(ActionType.PodFilterInNode, {}),

  containerList: reduceToPayload(ActionType.FetchContainerList, []),

  podSelection: reduceToPayload(ActionType.PodSelection, []),

  deletePodFlow: generateWorkflowReducer({
    actionType: ActionType.DeletePod
  }),

  updateGrayTappFlow: generateWorkflowReducer({
    actionType: ActionType.UpdateGrayTapp
  }),

  editTappGrayUpdate: reduceToPayload(ActionType.W_TappGrayUpdate, []),

  isShowLoginDialog: reduceToPayload(ActionType.IsShowLoginDialog, false),

  logQuery: generateQueryReducer({
    actionType: ActionType.QueryLogList
  }),

  logList: generateFetcherReducer<RecordSet<string>>({
    actionType: ActionType.FetchLogList,
    initialData: {
      recordCount: 0,
      records: [] as string[]
    }
  }),

  logOption: logOptionReducer,

  secretQuery: generateQueryReducer({
    actionType: ActionType.QuerySecretList
  }),

  secretList: generateFetcherReducer<RecordSet<Resource>>({
    actionType: ActionType.FetchSecretList,
    initialData: {
      recordCount: 0,
      records: [] as Resource[]
    }
  }),

  secretSelection: reduceToPayload(ActionType.SecretSelection, []),

  modifyNamespaceSecretFlow: generateWorkflowReducer({
    actionType: ActionType.ModifyNamespaceSecret
  })
});

export const ResourceDetailReducer = (state, action) => {
  let newState = state;
  // 销毁详情页面
  if (action.type === ActionType.ClearResourceDetail) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};
