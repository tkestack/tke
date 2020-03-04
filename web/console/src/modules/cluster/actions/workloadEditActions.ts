import { FetchOptions, generateFetcherActionCreator } from '@tencent/ff-redux';
import { extend, ReduxAction, uuid } from '@tencent/qcloud-lib';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';
import { resourceConfig } from '../../../../config';
import { cloneDeep } from '../../common/utils';
import * as ActionType from '../constants/ActionType';
import {
  initContainer,
  initEnv,
  initHpaMetrics,
  initMount,
  initValueFrom,
  initVolume,
  initWorkloadLabel
} from '../constants/initState';
import {
  ContainerItem,
  EnvItem,
  HpaMetrics,
  MetricOption,
  Resource,
  ResourceFilter,
  RootState,
  ValueFrom,
  VolumeItem,
  WorkloadLabel
} from '../models';
import { router } from '../router';
import * as WebAPI from '../WebAPI';
import { initCronMetrics, initmatchExpressions, initWorkloadAnnotataions } from './../constants/initState';
import { CronMetrics, MatchExpressions } from './../models/WorkloadEdit';
import { validateWorkloadActions } from './validateWorkloadActions';
import { workloadConfigActions } from './workloadConfigActions';
import { workloadPvcActions } from './workloadPvcActions';
import { workloadSecretActions } from './workloadSecretActions';

type GetState = () => RootState;
const fetchOptions: FetchOptions = {
  noCache: false
};

/** ========== start 更新实例数量，拉取hpa的相关信息 ================= */

const fetchHpaActions = generateFetcherActionCreator({
  actionType: ActionType.W_FetchHpaList,
  fetcher: async (getState: GetState, fetchOptions, dispatch) => {
    let { clusterVersion } = getState(),
      { hpaQuery } = getState().subRoot.workloadEdit;

    let isClearData = fetchOptions && fetchOptions.noCache ? true : false;
    let hpaResourceInfo = resourceConfig(clusterVersion)['hpa'];
    let response = await WebAPI.fetchSpecificResourceList(hpaQuery, hpaResourceInfo, isClearData, true);
    return response;
  },
  finish: (dispatch, getState: GetState) => {
    let { hpaList } = getState().subRoot.workloadEdit;

    if (hpaList.data.recordCount) {
      // 如果拉取到了hpa，那就说明该deployment创建了相对应的Hpa资源，那么需要更新hpa的相关信息
      dispatch(workloadEditActions.updateScaleType('autoScale'));

      let hpaInfo = hpaList.data.records[0];
      dispatch(workloadEditActions.inputMinReplicas(hpaInfo.spec.minReplicas || ''));
      dispatch(workloadEditActions.inputMaxReplicas(hpaInfo.spec.maxReplicas || ''));

      let metrics = hpaInfo.spec.metrics || [];
      dispatch(hpaActions.initMetricForUpdate(metrics));
    }
  }
});

const queryHpaActions = generateQueryActionCreator<ResourceFilter>({
  actionType: ActionType.W_QueryHpaList,
  bindFetcher: fetchHpaActions
});

const hpaRestActions = {
  /** 初始化autoscale的相关配置策略 */
  initMetricForUpdate: (metrics: MetricOption[]) => {
    return async (dispatch, getState: GetState) => {
      let newMetrics = [];

      const reduceCpu = cpu => {
        if (isNaN(cpu)) {
          return parseInt(cpu) / 1000;
        } else {
          return cpu;
        }
      };

      metrics.forEach(item => {
        let metricType = '',
          value;

        if (item.type === 'Resource') {
          if (item.resource.name === 'cpu') {
            metricType = item.resource.targetAverageUtilization ? 'cpuUtilization' : 'cpuAverage';
            value = item.resource.targetAverageUtilization
              ? item.resource.targetAverageUtilization
              : reduceCpu(item.resource.targetAverageValue);
          } else {
            metricType = item.resource.targetAverageUtilization ? 'memoryUtilization' : 'memoryAverage';
            value = item.resource.targetAverageUtilization
              ? item.resource.targetAverageUtilization
              : parseInt(item.resource.targetAverageValue + '');
          }
        } else {
          metricType = item.pods.metricName === 'pod_in_bandwidth' ? 'inBandwidth' : 'outBandwidth';
          value = item.pods.targetAverageValue;
        }

        let metric: HpaMetrics = Object.assign({}, initHpaMetrics, {
          id: uuid(),
          type: metricType,
          value
        });

        newMetrics.push(metric);
      });

      dispatch({
        type: ActionType.W_UpdateMetrics,
        payload: newMetrics
      });
    };
  }
};

const hpaActions = extend({}, queryHpaActions, fetchHpaActions, hpaRestActions);
/** ========== end 更新实例数量，拉取hpa的相关信息 ================= */

/** ============================== start cronhpa的相关操作 ============================== */

const cronHpaActions = {
  /** 新增metric */
  addMetric: () => {
    return async (dispatc: Redux.Dispatch, getState: GetState) => {
      let metricsArr: CronMetrics[] = cloneDeep(getState().subRoot.workloadEdit.cronMetrics);

      metricsArr.push(Object.assign({}, initCronMetrics, { id: uuid() }));
      dispatc({
        type: ActionType.W_UpdateCronMetrics,
        payload: metricsArr
      });
    };
  },

  /** 删除metric */
  deleteMetric: (mId: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let metricsArr: HpaMetrics[] = cloneDeep(getState().subRoot.workloadEdit.cronMetrics),
        mIndex = metricsArr.findIndex(item => item.id === mId);
      metricsArr.splice(mIndex, 1);
      dispatch({
        type: ActionType.W_UpdateCronMetrics,
        payload: metricsArr
      });
    };
  },

  /** 更新metric */
  updateMetric: (obj: any, mId: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let metricsArr: HpaMetrics[] = cloneDeep(getState().subRoot.workloadEdit.cronMetrics),
        mIndex = metricsArr.findIndex(item => item.id === mId),
        objKeys = Object.keys(obj);

      objKeys.forEach(item => {
        metricsArr[mIndex][item] = obj[item];
      });

      dispatch({
        type: ActionType.W_UpdateCronMetrics,
        payload: metricsArr
      });
    };
  }
};
/** ============================== end cronhpa的相关操作 ============================== */

export const workloadEditActions = {
  /** hpa的相关操作 */
  hpa: hpaActions,

  /** configMap的相关操作 */
  config: workloadConfigActions,

  /** secret的相关操作 */
  secret: workloadSecretActions,

  /** pvc的相关操作 */
  pvc: workloadPvcActions,

  /* cronhpa的相关操作 */
  cronhpa: cronHpaActions,

  /** 更新workload名称 */
  inputWorkloadName: (name: string) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.W_WorkloadName,
        payload: name
      });

      // 输入之后，自动同步到标签的第一项
      let labelId = getState().subRoot.workloadEdit.workloadLabels[0].id;
      dispatch(workloadEditActions.updateLabels({ labelValue: name }, labelId + ''));
    };
  },

  /** 更新workload的描述名称 */
  inputWorkloadDesp: (desp: string): ReduxAction<string> => {
    return {
      type: ActionType.W_Description,
      payload: desp
    };
  },

  /** 新增labels */
  addLabels: () => {
    return async (dispatch, getState: GetState) => {
      let labels: WorkloadLabel[] = cloneDeep(getState().subRoot.workloadEdit.workloadLabels);

      labels.push(Object.assign({}, initWorkloadLabel, { id: uuid() }));
      dispatch({
        type: ActionType.W_WorkloadLabels,
        payload: labels
      });
    };
  },

  /** 删除labels */
  deleteLabels: (labelId: string) => {
    return async (dispatch, getState: GetState) => {
      let labels: WorkloadLabel[] = cloneDeep(getState().subRoot.workloadEdit.workloadLabels),
        labelIndex = labels.findIndex(label => label.id === labelId);

      labels.splice(labelIndex, 1);
      dispatch({
        type: ActionType.W_WorkloadLabels,
        payload: labels
      });
    };
  },

  /** 更新labels */
  updateLabels: (obj: any, labelId: string) => {
    return async (dispatch, getState: GetState) => {
      let labels: WorkloadLabel[] = cloneDeep(getState().subRoot.workloadEdit.workloadLabels),
        labelIndex = labels.findIndex(label => label.id === labelId),
        objKeys = Object.keys(obj);

      objKeys.forEach(item => {
        labels[labelIndex][item] = obj[item];
      });
      dispatch({
        type: ActionType.W_WorkloadLabels,
        payload: labels
      });
    };
  },

  /** 新增annotataions */
  addAnnotations: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let annotataions: WorkloadLabel[] = cloneDeep(getState().subRoot.workloadEdit.workloadAnnotations);

      annotataions.push(Object.assign({}, initWorkloadAnnotataions, { id: uuid() }));
      dispatch({
        type: ActionType.W_WorkloadAnnotations,
        payload: annotataions
      });
    };
  },

  /** 删除annotataions */
  deleteAnnotations: (aId: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let annotataions: WorkloadLabel[] = cloneDeep(getState().subRoot.workloadEdit.workloadAnnotations),
        aIndex = annotataions.findIndex(annotation => annotation.id === aId);

      annotataions.splice(aIndex, 1);
      dispatch({
        type: ActionType.W_WorkloadAnnotations,
        payload: annotataions
      });
    };
  },

  /** 更新annotataions */
  updateAnnotations: (obj: any, aId: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let annotataions: WorkloadLabel[] = cloneDeep(getState().subRoot.workloadEdit.workloadAnnotations),
        aIndex = annotataions.findIndex(annotation => annotation.id === aId),
        objKeys = Object.keys(obj);

      objKeys.forEach(item => {
        annotataions[aIndex][item] = obj[item];
      });
      dispatch({
        type: ActionType.W_WorkloadAnnotations,
        payload: annotataions
      });
    };
  },

  /** 选择 workload 命名空间 */
  selectNamespace: (namespace: string) => {
    return async (dispatch, getState: GetState) => {
      let { route } = getState();

      dispatch({
        type: ActionType.W_Namespace,
        payload: namespace
      });

      // 校验命名空间是否合理
      dispatch(validateWorkloadActions.validateNamespace());

      // 选择命名空间之后，去拉取相对应的config : configMap 或者 是 secret 的列表
      dispatch(
        workloadEditActions.config.applyFilter({
          regionId: +route.queries['rid'],
          clusterId: route.queries['clusterId'],
          namespace
        })
      );
      dispatch(
        workloadEditActions.secret.applyFilter({
          regionId: +route.queries['rid'],
          clusterId: route.queries['clusterId'],
          namespace
        })
      );
      dispatch(
        workloadEditActions.pvc.applyFilter({
          regionId: +route.queries['rid'],
          clusterId: route.queries['clusterId'],
          namespace
        })
      );
    };
  },

  /** 选择resource的类型 */
  selectResourceType: (resourceType: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { scaleType } = getState().subRoot.workloadEdit;

      dispatch({
        type: ActionType.W_WorkloadType,
        payload: resourceType
      });

      // 这里去判断是否需要展示实例数量，只有deployment和statefulset以及tapp需要展现实例数量
      dispatch({
        type: ActionType.W_IsNeedContainerNum,
        payload:
          resourceType === 'deployment' || resourceType === 'statefulset' || resourceType === 'tapp' ? true : false
      });

      // 只有deployment有 自动调节，并且只有 deployment 和 statefulset需要使用
      if (resourceType === 'statefulset' && scaleType === 'autoScale') {
        dispatch(workloadEditActions.updateScaleType('manualScale'));
      }
    };
  },

  /** 输入cronjob的策略类型 */
  inputCronjobSchedule: (schedule: string): ReduxAction<string> => {
    return {
      type: ActionType.W_CronSchedule,
      payload: schedule.trim()
    };
  },

  /** 输入job的完成次数 */
  inputJobCompletion: (completion: string): ReduxAction<string> => {
    return {
      type: ActionType.W_Completion,
      payload: completion
    };
  },

  /** 输入job的并行度 */
  inputJobParallelism: (parallel: string): ReduxAction<string> => {
    return {
      type: ActionType.W_Parallelism,
      payload: parallel
    };
  },

  /** 选择失败重启策略 */
  selectRestartPolicy: (policy: string): ReduxAction<string> => {
    return {
      type: ActionType.W_RestartPolicy,
      payload: policy
    };
  },
  /**tapp节点异常重启策略 */
  selectNodeAbnormalMigratePolicy: (policy: string): ReduxAction<string> => {
    return {
      type: ActionType.W_NodeAbnormalMigratePolicy,
      payload: policy
    };
  },

  /** 新增数据卷 */
  addVolume: () => {
    return async (dispatch, getState: GetState) => {
      let volumes = cloneDeep(getState().subRoot.workloadEdit.volumes);
      let newVolume = Object.assign({}, initVolume, { id: uuid() });

      volumes.push(newVolume);
      dispatch({
        type: ActionType.W_UpdateVolumes,
        payload: volumes
      });
    };
  },

  /** 删除数据卷 */
  deleteVolume: (vId: string) => {
    return async (dispatch, getState: GetState) => {
      let volumes: VolumeItem[] = cloneDeep(getState().subRoot.workloadEdit.volumes),
        vIndex = volumes.findIndex(v => v.id === vId);

      volumes.splice(vIndex, 1);
      dispatch({
        type: ActionType.W_UpdateVolumes,
        payload: volumes
      });
    };
  },

  /** 新增容器挂载点 */
  addMount: (cId: string) => {
    return async (dispatch, getState: GetState) => {
      let containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cId),
        newMount = Object.assign({}, initMount[0], { id: uuid() });

      containers[cIndex]['mounts'] = containers[cIndex]['mounts'].concat([newMount]);
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },

  /** 删除容器挂载点 */
  deleteMount: (cId: string, mId: string) => {
    return async (dispatch, getState: GetState) => {
      let containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cId),
        mounts = containers[cIndex]['mounts'],
        mIndex = mounts.findIndex(m => m.id === mId);

      mounts.splice(mIndex, 1);
      containers[cIndex]['mounts'] = mounts;
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },

  /** 更新 volume的相关配置 */
  updateVolume: (obj: any, vId: string) => {
    return async (dispatch, getState: GetState) => {
      let volumes: VolumeItem[] = cloneDeep(getState().subRoot.workloadEdit.volumes),
        vIndex = volumes.findIndex(vol => vol.id === vId);

      let keys = Object.keys(obj);
      keys.forEach(item => {
        volumes[vIndex][item] = obj[item];
      });

      dispatch({
        type: ActionType.W_UpdateVolumes,
        payload: volumes
      });
    };
  },

  /** 更新cbs dialog的状态，是否展示模态框 */
  toggleCbsDialog: () => {
    return async (dispatch, getState: GetState) => {
      let isShow = getState().subRoot.workloadEdit.isShowCbsDialog;

      dispatch({
        type: ActionType.W_IsShowCbsDialog,
        payload: !isShow
      });
    };
  },

  /** 更新configmap dialog的状态，是否展示模态框 */
  toggleConfigDialog: () => {
    return async (dispatch, getState: GetState) => {
      let isShow = getState().subRoot.workloadEdit.isShowConfigDialog;

      dispatch({
        type: ActionType.W_IsShowConfigDialog,
        payload: !isShow
      });
    };
  },

  /** 是否展示pvc设置的的模态框 */
  togglePvcDialog: () => {
    return async (dispatch, getState: GetState) => {
      let isShow = getState().subRoot.workloadEdit.isShowPvcDialog;
      dispatch({
        type: ActionType.W_IsShowPvcDialog,
        payload: !isShow
      });
    };
  },

  /** 是否展示主机路径配置的模态框 */
  toggleHostPathDialog: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let isShow = getState().subRoot.workloadEdit.isShowHostPathDialog;
      dispatch({
        type: ActionType.W_IsShowHostPathDialog,
        payload: !isShow
      });
    };
  },

  /** 更新目前正在操作的volume表单的 volume */
  changeCurrentEditingVolumeId: (vId: string): ReduxAction<string> => {
    return {
      type: ActionType.W_CurrentEditingVolumeId,
      payload: vId
    };
  },

  /** 更新容器的信息 */
  updateContainer: (obj: any, cKey: string) => {
    return async (dispatch, getState: GetState) => {
      let containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cKey);

      let objKeys = Object.keys(obj);

      objKeys.forEach(key => {
        containers[cIndex][key] = obj[key];
      });
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },

  /** 删除容器 */
  deleteContainer: (cKey: string) => {
    return async (dispatch, getState: GetState) => {
      let containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cKey);

      containers.splice(cIndex, 1);
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },

  /** 更新容器的挂载点 */
  updateMount: (obj: any, cKey: string, mId: string) => {
    return async (dispatch, getState: GetState) => {
      let containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cKey),
        mIndex = containers[cIndex].mounts.findIndex(m => m.id === mId),
        objKeys = Object.keys(obj);

      objKeys.forEach(item => {
        containers[cIndex]['mounts'][mIndex][item] = obj[item];
      });
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },

  /** 更新容器的cpulimit的值 */
  updateCpuLimit: (obj: any, cId: string, cpuId: string) => {
    return async (dispatch, getState: GetState) => {
      let containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cId),
        cpuIndex = containers[cIndex].cpuLimit.findIndex(c => c.id === cpuId),
        objKeys = Object.keys(obj);

      objKeys.forEach(item => {
        containers[cIndex]['cpuLimit'][cpuIndex][item] = obj[item];
      });
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },

  /** 更新容器的memlimit的值 */
  updateMemLimit: (obj: any, cId: string, mId: string) => {
    return async (dispatch, getState: GetState) => {
      let containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cId),
        memIndex = containers[cIndex].memLimit.findIndex(m => m.id === mId),
        objKeys = Object.keys(obj);

      objKeys.forEach(item => {
        containers[cIndex]['memLimit'][memIndex][item] = obj[item];
      });
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },

  /** 新增环境变量 */
  addEnv: (cKey: string) => {
    return async (dispatch, getState: GetState) => {
      let containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cKey);

      let newEnv: EnvItem = Object.assign({}, initEnv, { id: uuid() });

      containers[cIndex]['envs'].push(newEnv);
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },

  /** 删除环境变量 */
  deleteEnv: (cKey: string, eId: string) => {
    return async (dispatch, getState: GetState) => {
      let containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cKey),
        envs: EnvItem[] = containers[cIndex]['envs'],
        eIndex = envs.findIndex(e => e.id === eId);

      containers[cIndex]['envs'].splice(eIndex, 1);
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },

  /** 更新环境变量 */
  updateEnv: (obj: any, cKey: string, eId: string) => {
    return async (dispatch, getState: GetState) => {
      let containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cKey),
        envs: EnvItem[] = containers[cIndex]['envs'],
        eIndex = envs.findIndex(e => e.id === eId),
        objKeys = Object.keys(obj);

      objKeys.forEach(item => {
        envs[eIndex][item] = obj[item];
      });
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },

  /** 新增valueFrom */
  addValueFrom: (cKey: string) => {
    return async (dispatch, getState: GetState) => {
      let containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cKey);

      let newValueFrom: ValueFrom = Object.assign({}, initValueFrom, { id: uuid() });

      containers[cIndex]['valueFrom'].push(newValueFrom);
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },

  /** 删除valueFrom */
  deleteValueFrom: (cKey: string, vId: string) => {
    return async (dispatch, getState: GetState) => {
      let containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cKey),
        valueFrom: ValueFrom[] = containers[cIndex]['valueFrom'],
        vIndex = valueFrom.findIndex(v => v.id === vId);

      containers[cIndex]['valueFrom'].splice(vIndex, 1);
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },

  /** 更新valueFrom */
  updateValueFrom: (obj: any, cKey: string, vId: string) => {
    return async (dispatch, getState: GetState) => {
      let containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cKey),
        valueFrom: ValueFrom[] = containers[cIndex]['valueFrom'],
        vIndex = valueFrom.findIndex(v => v.id === vId),
        objKeys = Object.keys(obj);

      objKeys.forEach(item => {
        valueFrom[vIndex][item] = obj[item];
      });
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },

  /** 是否开启高级设置 */
  toggleAdvancedSetting: (cKey: string) => {
    return async (dispatch, getState: GetState) => {
      let containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cKey);

      containers[cIndex]['isOpenAdvancedSetting'] = !containers[cIndex]['isOpenAdvancedSetting'];
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },

  /** 更新健康检查相关配置 */
  updateHealthCheck: (obj: any, cKey: string, hType: string) => {
    return async (dispatch, getState) => {
      let containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cKey),
        objKeys = Object.keys(obj);

      objKeys.forEach(item => {
        if (hType) {
          containers[cIndex]['healthCheck'][hType][item] = obj[item];
        } else {
          // 更新是否开启健康检查等值
          containers[cIndex]['healthCheck'][item] = obj[item];
        }
      });
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },

  /** 新增容器实例 */
  addContainer: () => {
    return async (dispatch, getState: GetState) => {
      let containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        editingIndex = containers.findIndex(c => c.id === 'editing'),
        newContainer = Object.assign({}, initContainer, { id: uuid() });

      if (editingIndex > -1) {
        containers[editingIndex]['status'] = 'edited';
      }
      containers.push(newContainer);
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },

  /** 如果高级设置当中有错误，需要弹开告诉用户 */
  modifyAdvancedSettingValidate: (isNotOk: boolean, cId: string) => {
    return async (dispatch, getState: GetState) => {
      let containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cId),
        container = containers[cIndex];

      containers[cIndex]['isAdvancedError'] = isNotOk;
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });

      // 这里判断如果 高级设置里面 有错误，并且高级设置 是关闭状态下的，则需要展开
      if (container.status === 'editing' && container.isOpenAdvancedSetting === false && isNotOk === true) {
        dispatch(workloadEditActions.toggleAdvancedSetting(container.id + ''));
      }
    };
  },

  /** 变更当前的实例数量的更新类型 */
  updateScaleType: (scaleType: string): ReduxAction<string> => {
    return {
      type: ActionType.W_ChangeScaleType,
      payload: scaleType
    };
  },

  /** 变更实例的数量 */
  updateContainerNum: (num: string): ReduxAction<string> => {
    return {
      type: ActionType.W_ContainerNum,
      payload: num
    };
  },

  /** 是否启用同时创建Service */
  isCreateService: (isNeed: boolean): ReduxAction<boolean> => {
    return {
      type: ActionType.W_IsCreateService,
      payload: !isNeed
    };
  },

  /** 操作metric的相关 */
  updateMetric: (obj: any, mId: string) => {
    return async (dispatch, getState: GetState) => {
      let metricsArr: HpaMetrics[] = cloneDeep(getState().subRoot.workloadEdit.metrics),
        mIndex = metricsArr.findIndex(item => item.id === mId),
        objKeys = Object.keys(obj);

      objKeys.forEach(item => {
        metricsArr[mIndex][item] = obj[item];
      });

      dispatch({
        type: ActionType.W_UpdateMetrics,
        payload: metricsArr
      });
    };
  },

  /** 删除metric */
  deleteMetric: (mId: string) => {
    return async (dispatch, getState: GetState) => {
      let metricsArr: HpaMetrics[] = cloneDeep(getState().subRoot.workloadEdit.metrics),
        mIndex = metricsArr.findIndex(item => item.id === mId);
      metricsArr.splice(mIndex, 1);
      dispatch({
        type: ActionType.W_UpdateMetrics,
        payload: metricsArr
      });
    };
  },

  /** 新增metric */
  addMetric: () => {
    return async (dispatch, getState: GetState) => {
      let metricsArr: HpaMetrics[] = cloneDeep(getState().subRoot.workloadEdit.metrics);

      metricsArr.push(Object.assign({}, initHpaMetrics, { id: uuid() }));
      dispatch({
        type: ActionType.W_UpdateMetrics,
        payload: metricsArr
      });
    };
  },

  /** 修改最小实例数量 */
  inputMinReplicas: (min: string): ReduxAction<string> => {
    return {
      type: ActionType.W_MinReplicas,
      payload: min
    };
  },

  /** 修改最大实例数量 */
  inputMaxReplicas: (max: string): ReduxAction<string> => {
    return {
      type: ActionType.W_MaxReplicas,
      payload: max
    };
  },

  /**选择node亲和性方式 */
  selectNodeSelectType: (type: string) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.W_SelectNodeAffinityType,
        payload: type
      });
    };
  },

  /**选择亲和性调度指定节点调度 */
  selectNodeSelector: computers => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.W_SelectNodeSelector,
        payload: computers
      });
    };
  },
  deleteAffinityRule: (type?: string, id?: string) => {
    return async (dispatch, getState: GetState) => {
      let { nodeAffinityRule } = getState().subRoot.workloadEdit;
      let affinityRule;
      if (type === 'preferred') {
        let preferredMatchExpressions = cloneDeep(nodeAffinityRule.preferredExecution[0].preference.matchExpressions),
          index = preferredMatchExpressions.findIndex(e => e.id === id);
        preferredMatchExpressions.splice(index, 1);
        affinityRule = Object.assign({}, nodeAffinityRule, {
          preferredExecution: [{ preference: { matchExpressions: preferredMatchExpressions }, weight: 1 }]
        });
      } else if (type === 'required') {
        let requiredMatchExpressions = cloneDeep(nodeAffinityRule.requiredExecution[0].matchExpressions),
          index = requiredMatchExpressions.findIndex(e => e.id === id);
        requiredMatchExpressions.splice(index, 1);
        affinityRule = Object.assign({}, nodeAffinityRule, {
          requiredExecution: [{ matchExpressions: requiredMatchExpressions }]
        });
      }
      dispatch({
        type: ActionType.W_UpdateNodeAffinityRule,
        payload: affinityRule
      });
    };
  },
  /**添加node亲和性自定义规则 */
  updateAffinityRule: (type?: string, id?: string, obj?: any) => {
    return async (dispatch, getState: GetState) => {
      let { nodeAffinityRule } = getState().subRoot.workloadEdit;
      let affinityRule;
      if (type === 'preferred') {
        let preferredMatchExpressions = cloneDeep(nodeAffinityRule.preferredExecution[0].preference.matchExpressions),
          index = preferredMatchExpressions.findIndex(e => e.id === id),
          objKeys = Object.keys(obj);
        objKeys.forEach(item => {
          preferredMatchExpressions[index][item] = obj[item];
        });
        affinityRule = Object.assign({}, nodeAffinityRule, {
          preferredExecution: [{ preference: { matchExpressions: preferredMatchExpressions }, weight: 1 }]
        });
      } else if (type === 'required') {
        let requiredMatchExpressions = cloneDeep(nodeAffinityRule.requiredExecution[0].matchExpressions),
          index = requiredMatchExpressions.findIndex(e => e.id === id),
          objKeys = Object.keys(obj);
        objKeys.forEach(item => {
          requiredMatchExpressions[index][item] = obj[item];
        });
        affinityRule = Object.assign({}, nodeAffinityRule, {
          requiredExecution: [{ matchExpressions: requiredMatchExpressions }]
        });
      }
      dispatch({
        type: ActionType.W_UpdateNodeAffinityRule,
        payload: affinityRule
      });
    };
  },
  /**更新node亲和性自定义规则 */
  addAffinityRule: (type?: string) => {
    return async (dispatch, getState: GetState) => {
      let { nodeAffinityRule } = getState().subRoot.workloadEdit;
      let newRule: MatchExpressions = Object.assign({}, initmatchExpressions, { id: uuid() });
      let affinityRule;
      if (type === 'preferred') {
        let preferredExecution = cloneDeep(nodeAffinityRule.preferredExecution);
        preferredExecution[0].preference.matchExpressions.push(newRule);
        affinityRule = Object.assign({}, nodeAffinityRule, { preferredExecution });
      } else if (type === 'required') {
        let requiredExecution = cloneDeep(nodeAffinityRule.requiredExecution);
        requiredExecution[0].matchExpressions.push(newRule);
        affinityRule = Object.assign({}, nodeAffinityRule, { requiredExecution });
      }
      dispatch({
        type: ActionType.W_UpdateNodeAffinityRule,
        payload: affinityRule
      });
    };
  },
  /** 判断是否是gpu的白名单 */
  isCanUseGpu: () => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.W_IsCanUseGpu,
        payload: true
      });
    };
  },

  /**判断是否可以创建tapp */
  isCanUseTapp: () => {
    return async (dispatch, getState: GetState) => {
      let { addons } = getState().subRoot;
      //如果创建了tappConto
      let result =
        addons['TappController'] && addons['TappController'].status.toLocaleLowerCase() === 'running' ? true : false;

      dispatch({
        type: ActionType.W_IsCanUseTapp,
        payload: result
      });
    };
  },
  /** 选择网络模式 */
  selectNetworkType: (networkType: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.W_NetworkType,
        payload: networkType
      });
    };
  },

  /** 选择浮动ip 回收策略**/
  selectFloatingIPReleasePolicy: (floatingIPReleasePolicy: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.W_FloatingIPReleasePolicy,
        payload: floatingIPReleasePolicy
      });
    };
  },

  /** ================ start 下面是滚动更新镜像的相关操作 ======================== */

  /** 改变资源的更新方式 */
  changeResourceUpdateType: (type: string): ReduxAction<string> => {
    return {
      type: ActionType.W_ResourceUpdateType,
      payload: type
    };
  },

  /** 修改更新间隔 */
  inputMinReadySeconds: (seconds: string): ReduxAction<string> => {
    return {
      type: ActionType.W_MinReadySeconds,
      payload: seconds
    };
  },

  /** 修改批量大小 */
  inputBatchSize: (size: string): ReduxAction<string> => {
    return {
      type: ActionType.W_BatchSize,
      payload: size
    };
  },

  /** 修改maxSurge的大小 */
  inputMaxSurge: (size: string): ReduxAction<string> => {
    return {
      type: ActionType.W_MaxSurge,
      payload: size
    };
  },

  /** 修改maxUnavaiable的大小 */
  inputMaxUnavaiable: (size: string): ReduxAction<string> => {
    return {
      type: ActionType.W_MaxUnavailable,
      payload: size
    };
  },

  /** 修改partition的大小 */
  inputPartition: (size: string): ReduxAction<string> => {
    return {
      type: ActionType.W_Partition,
      payload: size
    };
  },

  /** 选择更新策略 */
  changeRollingUpdateStrategy: (type: string): ReduxAction<string> => {
    return {
      type: ActionType.W_RollingUpdateStrategy,
      payload: type
    };
  },

  /** 初始化container的信息 */
  initContainersForUpdate: (containers: any[]) => {
    return async (dispatch, getState: GetState) => {
      let newContainers = [];

      newContainers = containers.map(item => {
        let [registry = '', tag = ''] = item.image.split(':');

        let tmp: ContainerItem = Object.assign({}, initContainer, {
          id: uuid(),
          registry,
          tag,
          name: item.name
        });
        return tmp;
      });

      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: newContainers
      });
    };
  },
  // W_UpdateOversoldRatio
  initOversoldRatio: oversoldRatio => {
    return async (dispatch, getState: GetState) => {
      let containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers);
      let cpuRatio = oversoldRatio.cpu ? true : false,
        memoryRatio = oversoldRatio.memory ? true : false;
      containers.forEach((item, index) => {
        if (cpuRatio) {
          item.cpuLimit[0].value = '';
        }
        if (memoryRatio) {
          item.memLimit[0].value = '';
        }
      });
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
      dispatch({
        type: ActionType.W_UpdateOversoldRatio,
        payload: oversoldRatio
      });
    };
  },

  /** 滚动更新镜像的时候，初始化一些数据 */
  initWorkloadEditForUpdateRegistry: (resource: Resource) => {
    return async (dispatch, getState: GetState) => {
      let { route } = getState(),
        urlParams = router.resolve(route);

      // 判断是否为deployment, deployment和其他两种资源你的 更新策略放在不同字段当中
      dispatch(workloadEditActions.selectResourceType(urlParams['resourceName']));
      let isDeployment = urlParams['resourceName'] === 'deployment',
        isStatefulset = urlParams['resourceName'] === 'statefulset',
        isDaemonset = urlParams['resourceName'] === 'daemonset',
        isTapp = urlParams['resourceName'] === 'tapp';

      // 当前的资源的更新策略类型
      let strategyType = isDeployment ? resource.spec.strategy.type : resource.spec.updateStrategy.type;

      let containers = resource.spec.template && resource.spec.template.spec.containers;

      // 初始化容器的相关信息
      dispatch(workloadEditActions.initContainersForUpdate(containers));
      // 初始化更新的方式
      strategyType && dispatch(workloadEditActions.changeResourceUpdateType(strategyType));

      // 这里是去更新策略配置的信息: Deployment
      if (strategyType === 'RollingUpdate') {
        if (isDeployment) {
          // 如果只有maxSurge 有值，而maxUnavailable为0，则为启动新的pod，停止旧的pod
          let maxSurge = resource.spec.strategy.rollingUpdate.maxSurge,
            maxUnavailable = resource.spec.strategy.rollingUpdate.maxUnavailable,
            rollingUpdateStrategy = 'createPod',
            minReadySeconds = resource.spec.minReadySeconds;

          // 更新当前的最小更新间隔时间
          minReadySeconds && dispatch(workloadEditActions.inputMinReadySeconds(minReadySeconds));

          if (maxSurge === 0) {
            dispatch(workloadEditActions.inputBatchSize(maxUnavailable));
            rollingUpdateStrategy = 'destroyPod';
          } else if (maxUnavailable === 0) {
            dispatch(workloadEditActions.inputBatchSize(maxSurge));
            rollingUpdateStrategy = 'createPod';
          } else {
            dispatch(workloadEditActions.inputMaxSurge(maxSurge));
            dispatch(workloadEditActions.inputMaxUnavaiable(maxUnavailable));
            rollingUpdateStrategy = 'userDefined';
          }

          dispatch(workloadEditActions.changeRollingUpdateStrategy(rollingUpdateStrategy));
        }

        // 这里是去更新策略配置的信息：Statefulset
        if (isStatefulset) {
          let partition = resource.spec.updateStrategy.rollingUpdate.partition || '0';
          dispatch(workloadEditActions.inputPartition(partition));
        }

        // 这里是去更新策略配置的信息：Daemonset
        if (isDaemonset) {
          let maxUnavailable = resource.spec.updateStrategy.rollingUpdate.maxUnavailable;
          dispatch(workloadEditActions.inputMaxUnavaiable(maxUnavailable));
        }
      }

      // 这里是去更新策略配置的信息：Tapp
      if (isTapp) {
        let maxUnavailable = resource.spec.updateStrategy.maxUnavailable;
        dispatch(workloadEditActions.inputMaxUnavaiable(maxUnavailable));
      }
    };
  },

  /** ================ end 下面是滚动更新镜像的相关操作 ======================== */

  /** 离开创建页面，清除workloadEdit当中的内容 */
  clearWorkloadEdit: (): ReduxAction<any> => {
    return {
      type: ActionType.ClearWorkloadEdit
    };
  }
};
