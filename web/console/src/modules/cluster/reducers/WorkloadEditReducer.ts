import { combineReducers } from 'redux';

import { RecordSet, reduceToPayload } from '@tencent/ff-redux';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { generateQueryReducer } from '@tencent/qcloud-redux-query';

import { initValidator } from '../../common/models';
import * as ActionType from '../constants/ActionType';
import {
  initAffinityRule,
  initContainer,
  initCronMetrics,
  initHpaMetrics,
  initSpecificLabel
} from '../constants/initState';
import { Resource } from '../models';

/** ========= configMap 和 secret 弹窗的相关数据的reducer ======= */
const TempConfigEditReducer = combineReducers({
  configQuery: generateQueryReducer({
    actionType: ActionType.W_QueryConfig
  }),

  configList: generateFetcherReducer<RecordSet<Resource>>({
    actionType: ActionType.W_FetchConfig,
    initialData: {
      recordCount: 0,
      records: [] as Resource[]
    }
  }),

  secretQuery: generateQueryReducer({
    actionType: ActionType.W_QuerySecret
  }),

  secretList: generateFetcherReducer<RecordSet<Resource>>({
    actionType: ActionType.W_FetchSecret,
    initialData: {
      recordCount: 0,
      records: [] as Resource[]
    }
  }),

  configItems: reduceToPayload(ActionType.W_UpdateConfigItems, []),

  configSelection: reduceToPayload(ActionType.W_ConfigSelect, []),

  keyType: reduceToPayload(ActionType.W_ChangeKeyType, 'all'),

  configKeys: reduceToPayload(ActionType.W_UpdateConfigKeys, [])
});

const WorkloadConfigEditReducer = (state, action) => {
  let newState = state;
  // 销毁创建workload页面
  if (action.type === ActionType.ClearWorkloadConfig) {
    newState = undefined;
  }
  return TempConfigEditReducer(newState, action);
};

/** ========= end--- configMap 和 secret 弹窗的相关数据的reducer ======= */

const TempReducer = combineReducers({
  workloadName: reduceToPayload(ActionType.W_WorkloadName, ''),

  v_workloadName: reduceToPayload(ActionType.WV_WorkloadName, initValidator),

  description: reduceToPayload(ActionType.W_Description, ''),

  v_description: reduceToPayload(ActionType.WV_Description, initValidator),

  workloadLabels: reduceToPayload(ActionType.W_WorkloadLabels, [initSpecificLabel]),

  workloadAnnotations: reduceToPayload(ActionType.W_WorkloadAnnotations, []),

  namespace: reduceToPayload(ActionType.W_Namespace, ''),

  v_namespace: reduceToPayload(ActionType.WV_Namespace, initValidator),

  workloadType: reduceToPayload(ActionType.W_WorkloadType, 'deployment'),

  cronSchedule: reduceToPayload(ActionType.W_CronSchedule, ''),

  v_cronSchedule: reduceToPayload(ActionType.WV_CronSchedule, initValidator),

  completion: reduceToPayload(ActionType.W_Completion, '1'),

  v_completion: reduceToPayload(ActionType.WV_Completion, initValidator),

  parallelism: reduceToPayload(ActionType.W_Parallelism, '1'),

  v_parallelism: reduceToPayload(ActionType.WV_Parallelism, initValidator),

  restartPolicy: reduceToPayload(ActionType.W_RestartPolicy, 'OnFailure'),

  nodeAbnormalMigratePolicy: reduceToPayload(ActionType.W_NodeAbnormalMigratePolicy, 'true'),

  volumes: reduceToPayload(ActionType.W_UpdateVolumes, []),

  isAllVolumeIsMounted: reduceToPayload(ActionType.W_IsAllVolumeIsMounted, false),

  isShowCbsDialog: reduceToPayload(ActionType.W_IsShowCbsDialog, false),

  isShowConfigDialog: reduceToPayload(ActionType.W_IsShowConfigDialog, false),

  isShowPvcDialog: reduceToPayload(ActionType.W_IsShowPvcDialog, false),

  isShowHostPathDialog: reduceToPayload(ActionType.W_IsShowHostPathDialog, false),

  currentEditingVolumeId: reduceToPayload(ActionType.W_CurrentEditingVolumeId, ''),

  pvcQuery: generateQueryReducer({
    actionType: ActionType.W_QueryPvcList
  }),

  pvcList: generateFetcherReducer<RecordSet<Resource>>({
    actionType: ActionType.W_FetchPvcList,
    initialData: {
      recordCount: 0,
      records: [] as Resource[]
    }
  }),

  configEdit: WorkloadConfigEditReducer,

  canAddContainer: reduceToPayload(ActionType.W_CanAddContainer, false),

  containers: reduceToPayload(ActionType.W_UpdateContainers, [initContainer]),

  scaleType: reduceToPayload(ActionType.W_ChangeScaleType, 'manualScale'),

  isOpenCronHpa: reduceToPayload(ActionType.W_IsOpenCronHpa, false),

  containerNum: reduceToPayload(ActionType.W_ContainerNum, '1'),

  isNeedContainerNum: reduceToPayload(ActionType.W_IsNeedContainerNum, true),

  minReplicas: reduceToPayload(ActionType.W_MinReplicas, ''),

  v_minReplicas: reduceToPayload(ActionType.WV_MinReplicas, initValidator),

  maxReplicas: reduceToPayload(ActionType.W_MaxReplicas, ''),

  v_maxReplicas: reduceToPayload(ActionType.WV_MaxReplicas, initValidator),

  metrics: reduceToPayload(ActionType.W_UpdateMetrics, [initHpaMetrics]),

  cronMetrics: reduceToPayload(ActionType.W_UpdateCronMetrics, [initCronMetrics]),

  isCreateService: reduceToPayload(ActionType.W_IsCreateService, true),

  imagePullSecrets: reduceToPayload(ActionType.ImagePullSecrets, []),

  isCanUseGpu: reduceToPayload(ActionType.W_IsCanUseGpu, false),

  isCanUseTapp: reduceToPayload(ActionType.W_IsCanUseTapp, false),

  networkType: reduceToPayload(ActionType.W_NetworkType, 'Overlay'),

  hpaQuery: generateQueryReducer({
    actionType: ActionType.W_QueryHpaList
  }),

  hpaList: generateFetcherReducer<RecordSet<Resource>>({
    actionType: ActionType.W_FetchHpaList,
    initialData: {
      recordCount: 0,
      records: [] as Resource[]
    }
  }),

  resourceUpdateType: reduceToPayload(ActionType.W_ResourceUpdateType, 'RollingUpdate'),

  minReadySeconds: reduceToPayload(ActionType.W_MinReadySeconds, '10'),

  v_minReadySeconds: reduceToPayload(ActionType.WV_MinReadySeconds, initValidator),

  rollingUpdateStrategy: reduceToPayload(ActionType.W_RollingUpdateStrategy, 'createPod'),

  batchSize: reduceToPayload(ActionType.W_BatchSize, '1'),

  v_batchSize: reduceToPayload(ActionType.WV_BatchSize, initValidator),

  maxSurge: reduceToPayload(ActionType.W_MaxSurge, ''),

  v_maxSurge: reduceToPayload(ActionType.WV_MaxSurge, initValidator),

  maxUnavailable: reduceToPayload(ActionType.W_MaxUnavailable, ''),

  v_maxUnavailable: reduceToPayload(ActionType.WV_MaxUnavailable, initValidator),

  partition: reduceToPayload(ActionType.W_Partition, ''),

  v_partition: reduceToPayload(ActionType.WV_Partition, initValidator),

  nodeSelection: reduceToPayload(ActionType.W_SelectNodeSelector, []),

  v_nodeSelection: reduceToPayload(ActionType.WV_NodeSelector, initValidator),

  nodeAffinityType: reduceToPayload(ActionType.W_SelectNodeAffinityType, 'unset'),

  nodeAffinityRule: reduceToPayload(ActionType.W_UpdateNodeAffinityRule, initAffinityRule),

  oversoldRatio: reduceToPayload(ActionType.W_UpdateOversoldRatio, {}),

  floatingIPReleasePolicy: reduceToPayload(ActionType.W_FloatingIPReleasePolicy, 'always')
});

export const WorkloadEditReducer = (state, action) => {
  let newState = state;
  // 销毁创建workload页面
  if (action.type === ActionType.ClearWorkloadEdit) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};
