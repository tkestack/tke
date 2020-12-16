import { affinityType } from './../constants/Config';
import * as ActionType from '../constants/ActionType';
import {
  RootState,
  WorkloadEdit,
  MountItem,
  LimitItem,
  HealthCheckItem,
  WorkloadLabel,
  HpaMetrics,
  VolumeItem,
  ContainerItem,
  ImagePullSecrets,
  ServiceEdit,
  ContainerEnv
} from '../models';
import { cloneDeep } from '../../common/utils';
import { workloadEditActions } from './workloadEditActions';
import { validateServiceActions } from './validateServiceActions';
import { AffinityRule, CronMetrics } from '../models/WorkloadEdit';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { TappGrayUpdateEditItem } from '../models/ResourceDetailState';

type GetState = () => RootState;

export const validateWorkloadActions = {
  /** 校验 workload名称是否正确 */
  _validateWorkloadName(name: string) {
    let reg = /^[a-z]([-a-z0-9]*[a-z0-9])?$/,
      status = 0,
      message = '';

    // 验证workload名称
    if (!name) {
      status = 2;
      message = t('Workload名称不能为空');
    } else if (name.length > 63) {
      status = 2;
      message = t('Workload名称不能超过63个字符');
    } else if (!reg.test(name)) {
      status = 2;
      message = t('Workload名称格式不正确');
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateWorkloadName() {
    return async (dispatch, getState: GetState) => {
      const { workloadName } = getState().subRoot.workloadEdit;

      const result = await validateWorkloadActions._validateWorkloadName(workloadName);

      dispatch({
        type: ActionType.WV_WorkloadName,
        payload: result
      });
    };
  },

  /** 校验执行策略是否正确 */
  _validateCronSchedule(schedule: string) {
    // cronjob的规则 * * * * * 分别代表 minute(0-59)、 hour(0-23)、 day(1-31)、 month(1-12)、 day of week(0-6)
    let status = 0,
      message = '',
      reg = /^(\*|([0-9]|1[0-9]|2[0-9]|3[0-9]|4[0-9]|5[0-9])|\*\/([0-9]|1[0-9]|2[0-9]|3[0-9]|4[0-9]|5[0-9])) (\*|([0-9]|1[0-9]|2[0-3])|\*\/([0-9]|1[0-9]|2[0-3])) (\*|([1-9]|1[0-9]|2[0-9]|3[0-1])|\*\/([1-9]|1[0-9]|2[0-9]|3[0-1])) (\*|([1-9]|1[0-2])|\*\/([1-9]|1[0-2])) (\*|([0-6])|\*\/([0-6]))$/;

    if (!schedule) {
      status = 2;
      message = t('执行策略不能为空');
    } else if (!reg.test(schedule)) {
      status = 2;
      message = t('执行策略格式不正确');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateCronSchedule() {
    return async (dispatch, getState: GetState) => {
      const { cronSchedule } = getState().subRoot.workloadEdit;

      const result = validateWorkloadActions._validateCronSchedule(cronSchedule);
      dispatch({
        type: ActionType.WV_CronSchedule,
        payload: result
      });
    };
  },

  /** 校验重复次数是否正确 */
  _validateJobCompletion(completion: number) {
    let status = 0,
      message = '',
      reg = /^\d+$/;

    if (completion < 1) {
      status = 2;
      message = t('执行次数至少为1');
    } else if (!reg.test(completion + '')) {
      status = 2;
      message = t('格式不正确，只能为数值');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateJobCompletion() {
    return async (dispatch, getState: GetState) => {
      const { completion } = getState().subRoot.workloadEdit;

      const result = validateWorkloadActions._validateJobCompletion(+completion);
      dispatch({
        type: ActionType.WV_Completion,
        payload: result
      });
    };
  },

  /** 校验job并行度 */
  _validateJobParallel(parallel: number) {
    let status = 0,
      message = '',
      reg = /^\d+$/;

    if (parallel < 1) {
      status = 2;
      message = t('Job并行度至少为1');
    } else if (!reg.test(parallel + '')) {
      status = 2;
      message = t('格式不正确，只能为数值');
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateJobParallel() {
    return async (dispatch, getState: GetState) => {
      const { parallelism } = getState().subRoot.workloadEdit;

      const result = validateWorkloadActions._validateJobParallel(+parallelism);
      dispatch({
        type: ActionType.WV_Parallelism,
        payload: result
      });
    };
  },

  /** 校验 workload 描述 */
  _validateWorkloadDesp(desp: string) {
    let status = 0,
      message = '';

    // 验证ingress描述
    if (desp && desp.length > 1000) {
      status = 2;
      message = t('Workload描述不能超过1000个字符');
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateWorkloadDesp() {
    return async (dispatch, getState: GetState) => {
      const { description } = getState().subRoot.workloadEdit;

      const result = validateWorkloadActions._validateWorkloadDesp(description);

      dispatch({
        type: ActionType.WV_Description,
        payload: result
      });
    };
  },

  /** 校验标签是否正确 */
  _validateWorkloadLabelKey(name: string, labels: WorkloadLabel[]) {
    // label 支持 [A-Z0-9a-z]开头和结尾，中间还可以有 -_.
    let reg = /^[A-Za-z0-9][-A-Za-z0-9_\./]*?[A-Za-z0-9]$/,
      status = 0,
      message = '';

    // 验证label
    if (!name) {
      status = 2;
      message = t('Key不能为空');
    } else if (!reg.test(name)) {
      status = 2;
      message = t('Key格式不正确');
    } else if (labels.filter(label => label.labelKey === name).length > 1) {
      status = 2;
      message = t('Key名称不能重复');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateWorkloadLabelKey(name: string, labelId: string) {
    return async (dispatch, getState: GetState) => {
      const labels: WorkloadLabel[] = cloneDeep(getState().subRoot.workloadEdit.workloadLabels),
        labelIndex = labels.findIndex(x => x.id === labelId),
        result = validateWorkloadActions._validateWorkloadLabelKey(name, labels);

      labels[labelIndex]['v_labelKey'] = result;
      dispatch({
        type: ActionType.W_WorkloadLabels,
        payload: labels
      });
    };
  },

  _validateAllWorkloadLabelKey(labels: WorkloadLabel[]) {
    let result = true;
    labels.forEach(label => {
      result = result && validateWorkloadActions._validateWorkloadLabelKey(label.labelKey, labels).status === 1;
    });
    return result;
  },

  validateAllWorkloadLabelKey() {
    return async (dispatch, getState: GetState) => {
      getState().subRoot.workloadEdit.workloadLabels.forEach(label => {
        dispatch(validateWorkloadActions.validateWorkloadLabelKey(label.labelKey, label.id + ''));
      });
    };
  },

  /** 校验标签的key是否正确 */
  _validateWorkloadLabelValue(value: string) {
    let status = 0,
      message = '',
      reg = /^[a-z]([-a-z0-9]*[a-z0-9])?$/;

    // 验证label
    if (!value) {
      status = 2;
      message = t('Value不能为空');
    } else if (!reg.test(value)) {
      status = 2;
      message = t('格式不正确');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateWorkloadLabelValue(value: string, labelId: string) {
    return async (dispatch, getState: GetState) => {
      const labels: WorkloadLabel[] = cloneDeep(getState().subRoot.workloadEdit.workloadLabels),
        labelIndex = labels.findIndex(x => x.id === labelId),
        result = validateWorkloadActions._validateWorkloadLabelValue(value);

      labels[labelIndex]['v_labelValue'] = result;
      dispatch({
        type: ActionType.W_WorkloadLabels,
        payload: labels
      });
    };
  },

  _validateAllWorkloadLabelValue(labels: WorkloadLabel[]) {
    let result = true;
    labels.forEach(label => {
      result = result && validateWorkloadActions._validateWorkloadLabelKey(label.labelKey, labels).status === 1;
    });
    return result;
  },

  validateAllWorkloadLabelValue() {
    return async (dispatch, getState: GetState) => {
      getState().subRoot.workloadEdit.workloadLabels.forEach(label => {
        dispatch(validateWorkloadActions.validateWorkloadLabelValue(label.labelValue, label.id + ''));
      });
    };
  },

  /** 校验annotations的key、value是否正确 */
  _validateWorkloadAnnotationsKey(name: string, annotations: WorkloadLabel[]) {
    let status = 0,
      message = '',
      reg = /^([A-Za-z0-9][-A-Za-z0-9_\.\/]*)?[A-Za-z0-9]$/;

    // 验证key
    if (!name) {
      status = 2;
      message = t('Key不能为空');
    } else if (!reg.test(name)) {
      status = 2;
      message = t('Key格式不正确');
    } else if (annotations.filter(annotataion => annotataion.labelKey === name).length > 1) {
      status = 2;
      message = t('Key名称不能重复');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateWorkloadAnnotationsKey(value: string, aId: string) {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      const annotations: WorkloadLabel[] = cloneDeep(getState().subRoot.workloadEdit.workloadAnnotations),
        aIndex = annotations.findIndex(x => x.id === aId),
        result = validateWorkloadActions._validateWorkloadAnnotationsKey(value, annotations);

      annotations[aIndex]['v_labelKey'] = result;
      dispatch({
        type: ActionType.W_WorkloadAnnotations,
        payload: annotations
      });
    };
  },

  _validateAllWorkloadAnnotationsKey(annotations: WorkloadLabel[]) {
    let result = true;
    annotations.forEach(annotation => {
      result =
        result &&
        validateWorkloadActions._validateWorkloadAnnotationsKey(annotation.labelKey, annotations).status === 1;
    });
    return result;
  },

  validateAllWorkloadAnnotationsKey() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      getState().subRoot.workloadEdit.workloadAnnotations.forEach(annotation => {
        dispatch(validateWorkloadActions.validateWorkloadAnnotationsKey(annotation.labelKey, annotation.id + ''));
      });
    };
  },

  /** 校验annotations的value是否正确 */
  _validateWorkloadAnnotationsValue(value: string) {
    let status = 0,
      message = '';
    // reg = /^([A-Za-z0-9][-A-Za-z0-9_\.]*)?[A-Za-z0-9]$/;

    // 验证key
    if (!value) {
      status = 2;
      message = t('Value不能为空');
      // } else if (!reg.test(value)) {
      //   status = 2;
      //   message = t('Value格式不正确');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateWorkloadAnnotationsValue(value: string, aId: string) {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      const annotations: WorkloadLabel[] = cloneDeep(getState().subRoot.workloadEdit.workloadAnnotations),
        aIndex = annotations.findIndex(x => x.id === aId),
        result = validateWorkloadActions._validateWorkloadAnnotationsValue(value);

      annotations[aIndex]['v_labelValue'] = result;
      dispatch({
        type: ActionType.W_WorkloadAnnotations,
        payload: annotations
      });
    };
  },

  _validateAllWorkloadAnnotationsValue(annotations: WorkloadLabel[]) {
    let result = true;
    annotations.forEach(annotation => {
      result = result && validateWorkloadActions._validateWorkloadAnnotationsValue(annotation.labelValue).status === 1;
    });
    return result;
  },

  validateAllWorkloadAnnotationsValue() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      getState().subRoot.workloadEdit.workloadAnnotations.forEach(annotation => {
        dispatch(validateWorkloadActions.validateWorkloadAnnotationsValue(annotation.labelValue, annotation.id + ''));
      });
    };
  },

  /** 校验命名空间 */
  _validateNamespace(namespace: string) {
    let status = 0,
      message = '';

    // 验证命名空间的选择
    if (!namespace) {
      status = 2;
      message = t('命名空间不能为空');
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateNamespace() {
    return async (dispatch, getState: GetState) => {
      const { namespace } = getState().subRoot.workloadEdit;

      const result = validateWorkloadActions._validateNamespace(namespace);

      dispatch({
        type: ActionType.WV_Namespace,
        payload: result
      });
    };
  },

  /** 校验数据卷的名称 */
  _validateVolumeName(name: string, volumes: VolumeItem[]) {
    let reg = /^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/,
      status = 0,
      message = '';

    if (!name) {
      status = 2;
      message = t('请填写数据卷名称');
    } else if (name.length > 63) {
      status = 2;
      message = t('数据卷名称不能超过63个字符');
    } else if (!reg.test(name)) {
      status = 2;
      message = t('只能包含小写字母、数字及分隔符("-")，且不能以分隔符开头或结尾');
    } else if (name && volumes.filter(vol => vol.name === name).length > 1) {
      status = 2;
      message = t('数据卷名称不能重复');
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateVolumeName(name: string, vId: string) {
    return async (dispatch, getState: GetState) => {
      const volumes: VolumeItem[] = cloneDeep(getState().subRoot.workloadEdit.volumes),
        vIndex = volumes.findIndex(vol => vol.id === vId),
        result = validateWorkloadActions._validateVolumeName(name, volumes);

      volumes[vIndex]['v_name'] = result;

      dispatch({
        type: ActionType.W_UpdateVolumes,
        payload: volumes
      });
    };
  },

  _validateAllVolumeName(volumes: VolumeItem[]) {
    let result = true;
    volumes.forEach(volume => {
      result = result && validateWorkloadActions._validateVolumeName(volume.name, volumes).status === 1;
    });

    return result;
  },

  validateAllVolumeName() {
    return async (dispatch, getState: GetState) => {
      const volumes = getState().subRoot.workloadEdit.volumes;
      volumes.forEach(volume => {
        dispatch(validateWorkloadActions.validateVolumeName(volume.name, volume.id + ''));
      });
    };
  },

  /** 校验pvc的相关选择 */
  _validateVolumePvc(pvcName: string, volumes: VolumeItem[]) {
    let status = 0,
      message = '';

    if (!pvcName) {
      status = 2;
      message = t('请选择pvc');
    } else if (volumes.filter(vol => vol.volumeType === 'pvc' && vol.pvcSelection === pvcName).length > 1) {
      status = 2;
      message = t('pvc不可重复');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateVolumePvc(pvcName: string, vId: string) {
    return async (dispatch, getState: GetState) => {
      const volumes: VolumeItem[] = cloneDeep(getState().subRoot.workloadEdit.volumes),
        vIndex = volumes.findIndex(vol => vol.id === vId),
        result = validateWorkloadActions._validateVolumePvc(pvcName, volumes);

      volumes[vIndex]['v_pvcSelection'] = result;
      dispatch({
        type: ActionType.W_UpdateVolumes,
        payload: volumes
      });
    };
  },

  _validateAllPvcSelection(volumes: VolumeItem[]) {
    let result = true;
    volumes.forEach(volume => {
      if (volume.volumeType === 'pvc') {
        result = result && validateWorkloadActions._validateVolumePvc(volume.pvcSelection, volumes).status === 1;
      }
    });
    return result;
  },

  validateAllPvcSelection() {
    return async (dispatch, getState: GetState) => {
      const volumes = getState().subRoot.workloadEdit.volumes;
      volumes.forEach(volume => {
        volume.volumeType === 'pvc' &&
          dispatch(validateWorkloadActions.validateVolumePvc(volume.pvcSelection, volume.id + ''));
      });
    };
  },

  /** 校验数据卷 nfs路径 */
  _validateNfsPath(path: string, vId: string, volumes: VolumeItem[]) {
    const volume = volumes.find(vol => vol.id === vId);
    let reg = /^(((ht|f)tps?):\/\/)?[\w-]+(\.[\w-]+)+([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?$/,
      status = 0,
      message = '';

    if (volume.volumeType === 'nfsDisk') {
      if (path) {
        if (!reg.test(path)) {
          status = 2;
          message = t('NFS路径格式不正确');
        } else if (volumes.filter(vol => vol.nfsPath === path && vol.volumeType === 'nfsDisk').length > 1) {
          status = 2;
          message = t('NFS路径称不可重复');
        } else {
          status = 1;
          message = '';
        }
      } else {
        status = 2;
        message = t('NFS路径不能为空');
      }
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateNfsPath(path: string, vId: string) {
    return async (dispatch, getState: GetState) => {
      const volumes: VolumeItem[] = cloneDeep(getState().subRoot.workloadEdit.volumes),
        vIndex = volumes.findIndex(v => v.id === vId),
        result = validateWorkloadActions._validateNfsPath(path, vId, volumes);

      volumes[vIndex]['v_nfsPath'] = result;

      dispatch({
        type: ActionType.W_UpdateVolumes,
        payload: volumes
      });
    };
  },

  _validateAllNfsPath(volumes: VolumeItem[]) {
    let result = true;
    volumes.forEach(volume => {
      result = result && validateWorkloadActions._validateNfsPath(volume.nfsPath, volume.id + '', volumes).status === 1;
    });

    return result;
  },

  validateAllNfsPath() {
    return async (dispatch, getState: GetState) => {
      const volumes = getState().subRoot.workloadEdit.volumes;
      volumes.forEach(volume => {
        dispatch(validateWorkloadActions.validateNfsPath(volume.nfsPath, volume.id + ''));
      });
    };
  },

  /** 校验hostPath是否正确 */
  _validateHostPath(path: string) {
    let status = 0,
      message = '',
      reg = /^\/[\/\w-_\.]*$/;

    if (!path) {
      status = 2;
      message = t('HostPath路径不能为空');
    } else if (!reg.test(path)) {
      status = 2;
      message = t('HostPath格式不正确');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateHostPath(path: string, vId: string) {
    return async (dispatch, getState: GetState) => {
      const volumes: VolumeItem[] = cloneDeep(getState().subRoot.workloadEdit.volumes),
        vIndex = volumes.findIndex(v => v.id === vId),
        result = validateWorkloadActions._validateHostPath(path);

      volumes[vIndex]['v_hostPath'] = result;
      dispatch({
        type: ActionType.W_UpdateVolumes,
        payload: volumes
      });
    };
  },

  _validateAllHostPath(volumes: VolumeItem[]) {
    let result = true;
    volumes.forEach(volume => {
      if (volume.volumeType === 'hostPath') {
        result = result && validateWorkloadActions._validateHostPath(volume.hostPath).status === 1;
      }
    });
    return result;
  },

  validateAllHosPath() {
    return async (dispatch, getState: GetState) => {
      const volumes = getState().subRoot.workloadEdit.volumes;
      volumes.forEach(volume => {
        volume.volumeType === 'hostPath' &&
          dispatch(validateWorkloadActions.validateHostPath(volume.hostPath, volume.id + ''));
      });
    };
  },

  _validateAllVolumeIsMounted(volumes: VolumeItem[], containers: ContainerItem[]) {
    let allIsMounted = true,
      validateVolumes = {};

    // 已经挂载的数据卷的名称数组
    const mounts: MountItem[] = [];
    containers.forEach(c => {
      mounts.push(...c.mounts);
    });
    const mountsName = mounts.map(m => m.volume);

    volumes.forEach(v => {
      const isMounted = !!mountsName.find(m => m === v.name);
      allIsMounted = allIsMounted && isMounted;

      validateVolumes[v.id + ''] = isMounted;
    });

    return { allIsMounted, validateVolumes };
  },

  validateAllVolumeIsMounted() {
    return async (dispatch, getState: GetState) => {
      const { volumes, containers } = getState().subRoot.workloadEdit;

      const result = validateWorkloadActions._validateAllVolumeIsMounted(volumes, containers);

      // 更新全部是否ok的标志
      dispatch({
        type: ActionType.W_IsAllVolumeIsMounted,
        payload: result.allIsMounted
      });

      // 更新每一个挂在项的校验状态
      const newVolumes: VolumeItem[] = cloneDeep(volumes),
        validateVolumes = result.validateVolumes;
      newVolumes.forEach(v => {
        v['isMounted'] = validateVolumes[v.id + ''];
      });
      dispatch({
        type: ActionType.W_UpdateVolumes,
        payload: newVolumes
      });
    };
  },

  /** 判断是否可以新增数据卷，对数据卷中的属性状态进行相关的判断 */
  _canAddVolume(volumes: VolumeItem[]) {
    let result = true;
    volumes.forEach(volume => {
      result = result && validateWorkloadActions._validateVolumeName(volume.name, volumes).status === 1;
      if (volume.volumeType === 'pvc') {
        result = result && validateWorkloadActions._validateVolumePvc(volume.pvcSelection, volumes).status === 1;
      } else if (volume.volumeType === 'hostPath') {
        result = result && validateWorkloadActions._validateHostPath(volume.hostPath).status === 1;
      }
    });
    return result;
  },

  /** 校验容器名称是否正确 */
  _validateContainerName(name: string, containers: ContainerItem[]) {
    let reg = /^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/,
      status = 0,
      message = '';

    // 验证容器服务名称
    if (!name) {
      status = 2;
      message = t('容器名称不能为空');
    } else if (containers.filter(c => c.name === name).length > 1) {
      status = 2;
      message = t('容器名称不能重复');
    } else if (name.length > 63) {
      status = 2;
      message = t('容器名称不能超过63个字符');
    } else if (!reg.test(name)) {
      status = 2;
      message = t('容器名称格式不正确');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateContainerName(name: string, cKey: string) {
    return async (dispatch, getState: GetState) => {
      const containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(item => item.id === cKey),
        result = validateWorkloadActions._validateContainerName(name, containers);

      containers[cIndex]['v_name'] = result;
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },

  validateAllContainerName() {
    return async (dispatch, getState: GetState) => {
      const containers: ContainerItem[] = getState().subRoot.workloadEdit.containers;
      containers.forEach(container => {
        dispatch(validateWorkloadActions.validateContainerName(container.name, container.id + ''));
      });
    };
  },

  /** 校验镜像名称是否合法 */
  _validateRegistrySelection(registry: string) {
    let status = 2,
      message = '';

    // 验证容器的镜像
    if (!registry) {
      status = 2;
      message = t('镜像不能为空');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateRegistrySelection(registry: string, cKey: string) {
    return async (dispatch, getState: GetState) => {
      const containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cKey),
        result = validateWorkloadActions._validateRegistrySelection(registry);

      containers[cIndex]['v_registry'] = result;
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },

  validateGrayUpdateRegistrySelection(index_in: number, imageName: string) {
    return async (dispatch, getState: GetState) => {
      const { editTappGrayUpdate } = getState().subRoot.resourceDetailState;
      const target: TappGrayUpdateEditItem = cloneDeep(editTappGrayUpdate);
      target.containers[index_in].v_imageName = validateWorkloadActions._validateRegistrySelection(
        imageName
      );

      dispatch({
        type: ActionType.W_TappGrayUpdate,
        payload: target
      });
    };
  },

  validateGrayUpdate() {
    return async (dispatch, getState: GetState) => {
      const { editTappGrayUpdate } = getState().subRoot.resourceDetailState;
      editTappGrayUpdate.containers.forEach((container, index_in) => {
        dispatch(
          validateWorkloadActions.validateGrayUpdateRegistrySelection(index_in, container.imageName)
        );
      });
    };
  },

  /** 校验数据卷挂在是否正确 */
  _validateVolumeMount(volumeName: string, mId: string, mounts: MountItem[], volumes: VolumeItem[]) {
    let mount = mounts.find(m => m.id === mId),
      status = 0,
      message = '';

    if (!volumeName && mount.mountPath) {
      status = 2;
      message = t('请选择数据卷');
    } else if (!volumes.find(v => v.name === volumeName)) {
      status = 2;
      message = t('数据卷不存在');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateVolumeMount(volumeName: string, cId: string, mId: string) {
    return async (dispatch, getState: GetState) => {
      const { containers, volumes } = getState().subRoot.workloadEdit;

      const newContainers: ContainerItem[] = cloneDeep(containers),
        cIndex = newContainers.findIndex(c => c.id === cId),
        mIndex = newContainers[cIndex]['mounts'].findIndex(m => m.id === mId),
        result = validateWorkloadActions._validateVolumeMount(
          volumeName,
          mId,
          newContainers[cIndex]['mounts'],
          volumes
        );

      newContainers[cIndex]['mounts'][mIndex]['v_volume'] = result;
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: newContainers
      });
    };
  },

  _validateAllVolumeMount(mounts: MountItem[], volumes: VolumeItem[]) {
    let result = true;
    mounts.forEach(m => {
      if (m.mountPath) {
        result =
          result && validateWorkloadActions._validateVolumeMount(m.volume, m.id + '', mounts, volumes).status === 1;
      }
    });
    return result;
  },

  validateAllVolumeMount(container: ContainerItem) {
    return async (dispatch, getState: GetState) => {
      container.mounts.forEach(mount => {
        if (mount.mountPath) {
          dispatch(validateWorkloadActions.validateVolumeMount(mount.volume, container.id + '', mount.id + ''));
        }
      });
    };
  },

  /** 校验挂载的挂载路径 */
  _validateVolumeMountPath(mountPath: string, mId: string, mounts: MountItem[]) {
    const mount = mounts.find(m => m.id === mId);
    let status = 0,
      message = '',
      reg = /:/;

    if (!mountPath && mount.volume) {
      status = 2;
      message = t('请输入目标路径');
    } else if (reg.test(mountPath)) {
      status = 2;
      message = t('目标路径不可包含 ":" ');
    } else if (mountPath && mounts.filter(m => m.mountPath === mountPath).length > 1) {
      status = 2;
      message = t('目标路径不可重复');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateVolumeMountPath(mountPath: string, cId: string, mId: string) {
    return async (dispatch, getState: GetState) => {
      const { containers, volumes } = getState().subRoot.workloadEdit;

      const newContainers: ContainerItem[] = cloneDeep(containers),
        cIndex = newContainers.findIndex(c => c.id === cId),
        mIndex = newContainers[cIndex]['mounts'].findIndex(m => m.id === mId),
        result = validateWorkloadActions._validateVolumeMountPath(mountPath, mId, newContainers[cIndex]['mounts']);

      newContainers[cIndex]['mounts'][mIndex]['v_mountPath'] = result;
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: newContainers
      });
    };
  },

  _validateAllVolumeMountPath(mounts: MountItem[]) {
    let result = true;
    mounts.forEach(mount => {
      if (mount.volume) {
        result =
          result &&
          validateWorkloadActions._validateVolumeMountPath(mount.mountPath, mount.id + '', mounts).status === 1;
      }
    });
    return result;
  },

  validateAllVolumeMountPath(container: ContainerItem) {
    return async (dispatch, getState: GetState) => {
      container.mounts.forEach(item => {
        if (item.volume) {
          dispatch(validateWorkloadActions.validateVolumeMountPath(item.mountPath, container.id + '', item.id + ''));
        }
      });
    };
  },

  /** 校验容器 CPU限制 */
  _validateCpuLimit(cpu: number, type: string, cpuLimit: LimitItem[]) {
    let reg = /^\d+(\.\d{1,3})?$/,
      status = 0,
      message = '';

    // 验证CPU限制
    if (isNaN(cpu)) {
      status = 2;
      message = t('数据格式不正确，CPU限制只能是小数，且只能精确到0.01');
    } else if (!cpu) {
      status = 1;
      message = '';
    } else if (!reg.test(cpu + '')) {
      status = 2;
      message = t('数据格式不正确，CPU限制只能是小数，且只能精确到0.01');
    } else if (cpu < 0.01) {
      status = 2;
      message = t('CPU限制最小为0.01');
    } else {
      if (type === 'limit') {
        status = 1;
        message = '';
      } else {
        const limit = cpuLimit.find(cpu => cpu.type === 'limit');
        if (!limit || (limit && cpu - +limit.value <= 0) || !limit.value) {
          status = 1;
          message = '';
        } else {
          status = 2;
          message = t('request限制不能超过limit限制');
        }
      }
    }

    return { status, message };
  },

  validateCpuLimit(cpu: string, cId: string, cpuId: string) {
    return async (dispatch, getState: GetState) => {
      const containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cId),
        cpuLimit = containers[cIndex].cpuLimit,
        cpuIndex = cpuLimit.findIndex(c => c.id === cpuId),
        result = validateWorkloadActions._validateCpuLimit(+cpu, cpuLimit[cpuIndex].type, cpuLimit);

      containers[cIndex]['cpuLimit'][cpuIndex]['v_value'] = result;
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },

  _validateAllCpuLimit(cpuLimit: LimitItem[]) {
    let result = true;
    cpuLimit.forEach(item => {
      result = result && validateWorkloadActions._validateCpuLimit(+item.value, item.type, cpuLimit).status === 1;
    });
    return result;
  },

  validateAllCpuLimit(container: ContainerItem) {
    return async (dispatch, getState: GetState) => {
      container.cpuLimit.forEach(item => {
        dispatch(validateWorkloadActions.validateCpuLimit(item.value, container.id + '', item.id + ''));
      });
    };
  },

  /** 校验容器 mem限制 */
  _validateMemLimit(mem: number, type: string, memLimit: LimitItem[]) {
    let reg = /^\d+?$/,
      status = 0,
      message = '';

    // 验证内存限制
    if (isNaN(mem)) {
      status = 2;
      message = t('只能输入正整数');
    } else if (!mem) {
      status = 1;
      message = '';
    } else if (!reg.test(mem + '')) {
      status = 2;
      message = t('只能输入正整数');
    } else {
      if (type === 'limit') {
        if (mem >= 4) {
          status = 1;
          message = '';
        } else {
          status = 2;
          message = t('limit限制要大于等于4Mib');
        }
      } else {
        const limit = memLimit[1];
        if (mem > 0) {
          if ((limit && limit.value === '') || (limit && mem - +limit.value <= 0)) {
            status = 1;
            message = '';
          } else {
            status = 2;
            message = t('request限制不能超过limit限制');
          }
        } else {
          status = 2;
          message = t('request限制要大于等于1Mib');
        }
      }
    }
    return { status, message };
  },
  /** 校验容器 Gpu限制 */
  _validateGpuCoreLimit(gpuCore: number) {
    let reg1 = /^\d+(\.\d{1,1})?$/,
      reg2 = /^\d+?$/,
      status = 0,
      message = '';

    // 验证内存限制
    if (isNaN(gpuCore)) {
      status = 2;
      message = t('gpu限制只能填写0.1-1或者1的正整数倍');
    } else if (!gpuCore && gpuCore !== 0) {
      status = 2;
      message = t('gpu限制不能为空');
    } else if (!reg1.test(gpuCore + '')) {
      status = 2;
      message = t('gpu限制只能填写0.1-1或者1的正整数倍');
    } else {
      if (gpuCore >= 1 && !reg2.test(gpuCore + '')) {
        status = 2;
        message = t('gpu限制只能填写0.1-1或者1的正整数倍');
      } else {
        status = 1;
        message = '';
      }
    }
    return { status, message };
  },

  validateGpuCoreLimit(gpuCore: string, cId: string) {
    return async (dispatch, getState: GetState) => {
      const containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cId),
        result = validateWorkloadActions._validateGpuCoreLimit(+gpuCore);

      containers[cIndex].v_gpuCore = result;
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },
  /** 校验容器 GpuMem限制 */
  _validateGpuMemLimit(gpuMem: number) {
    let reg = /^\d+?$/,
      status = 0,
      message = '';

    // 验证内存限制
    if (isNaN(gpuMem)) {
      status = 2;
      message = t('只能输入正整数');
    } else if (!gpuMem && gpuMem !== 0) {
      status = 2;
      message = t('GPU显存不能为空');
    } else if (!reg.test(gpuMem + '')) {
      status = 2;
      message = t('只能输入正整数');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateGpuMemLimit(gpuMem: string, cId: string) {
    return async (dispatch, getState: GetState) => {
      const containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cId),
        result = validateWorkloadActions._validateGpuMemLimit(+gpuMem);

      containers[cIndex].v_gpuMem = result;
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },
  validateMemLimit(mem: string, cId: string, mId: string) {
    return async (dispatch, getState: GetState) => {
      const containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cId),
        memLimit = containers[cIndex].memLimit,
        mIndex = memLimit.findIndex(m => m.id === mId),
        result = validateWorkloadActions._validateMemLimit(+mem, memLimit[mIndex].type, memLimit);

      containers[cIndex]['memLimit'][mIndex]['v_value'] = result;
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },

  _validateAllMemLimit(memLimit: LimitItem[]) {
    let result = true;
    memLimit.forEach(item => {
      result = result && validateWorkloadActions._validateMemLimit(+item.value, item.type, memLimit).status === 1;
    });
    return result;
  },

  validateAllMemLimit(container: ContainerItem) {
    return async (dispatch, getState: GetState) => {
      container.memLimit.forEach(item => {
        dispatch(validateWorkloadActions.validateMemLimit(item.value, container.id + '', item.id + ''));
      });
    };
  },

  /** 校验容器的 环境变量是否正确 */
  _validateEnvName(name: string, envs: ContainerEnv.ItemWithId[]) {
    let reg = /^[A-Za-z][-A-Za-z0-9_\.]*$/,
      status = 0,
      message = '';

    if (!name) {
      status = 2;
      message = t('环境变量名称不能为空');
    } else if (!reg.test(name)) {
      status = 2;
      message = t('环境变量名称格式不对');
    } else if (name && envs.filter(e => e.name === name).length > 1) {
      status = 2;
      message = t('环境变量不可重复');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateEnvName(name: string, cId: string, eId: string) {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      const containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cId),
        envs: ContainerEnv.ItemWithId[] = containers[cIndex].envItems,
        eIndex = envs.findIndex(e => e.id === eId),
        result = validateWorkloadActions._validateEnvName(name, envs);

      containers[cIndex]['envItems'][eIndex]['v_name'] = result;
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },

  _validateAllEnvName(envs: ContainerEnv.ItemWithId[]) {
    let result = true;
    envs.forEach(e => {
      result = result && validateWorkloadActions._validateEnvName(e.name, envs).status === 1;
    });
    return result;
  },

  validateAllEnvName(container: ContainerItem) {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      container.envItems.forEach(e => {
        dispatch(validateWorkloadActions.validateEnvName(e.name, container.id + '', e.id + ''));
      });
    };
  },

  /** 校验envItem的内容是否正确 */
  _validateNewEnvItemValue(value: string) {
    let status = 1,
      message = '';

    if (!value) {
      status = 2;
      message = t('请选择下拉项');
    }
    return { status, message };
  },

  validateNewEnvItemValue(keyNames: (keyof ContainerEnv.Item)[], cId: string, eId: string) {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      const containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cId),
        newEnvItems: ContainerEnv.ItemWithId[] = containers[cIndex]['envItems'],
        vIndex = newEnvItems.findIndex(e => e.id === eId);

      keyNames.forEach(keyName => {
        const result = validateWorkloadActions._validateNewEnvItemValue(newEnvItems[vIndex][keyName] as string);
        containers[cIndex]['envItems'][vIndex][`v_${keyName}`] = result;
      });
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },

  /** 校验所有的EnvItem */
  _validateAllNewEnvItemValue(envItems: ContainerEnv.ItemWithId[]) {
    let result = true;
    envItems.forEach(envItem => {
      if (envItem.type === ContainerEnv.EnvTypeEnum.ConfigMapRef) {
        result =
          result &&
          validateWorkloadActions._validateNewEnvItemValue(envItem.configMapName).status === 1 &&
          validateWorkloadActions._validateNewEnvItemValue(envItem.configMapDataKey).status === 1;
      } else if (envItem.type === ContainerEnv.EnvTypeEnum.SecretKeyRef) {
        result =
          result &&
          validateWorkloadActions._validateNewEnvItemValue(envItem.secretName).status === 1 &&
          validateWorkloadActions._validateNewEnvItemValue(envItem.secretDataKey).status === 1;
      }
    });
    return result;
  },

  validateAllNewEnvItemValue(container: ContainerItem) {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      container.envItems.forEach(envItem => {
        const cId = container.id + '';
        const eId = envItem.id + '';
        if (envItem.type === ContainerEnv.EnvTypeEnum.ConfigMapRef) {
          dispatch(validateWorkloadActions.validateNewEnvItemValue(['configMapName', 'configMapDataKey'], cId, eId));
        } else if (envItem.type === ContainerEnv.EnvTypeEnum.SecretKeyRef) {
          dispatch(validateWorkloadActions.validateNewEnvItemValue(['secretName', 'secretDataKey'], cId, eId));
        }
      });
    };
  },

  /** 校验容器的工作/日志目录设置 */
  _validateDir(dir: string, errorMessage: string) {
    let reg = /^\/(?:[a-zA-Z0-9-_.]+\/?)*$/,
      status = 0,
      message = '';

    if (dir && !reg.test(dir)) {
      status = 2;
      message = errorMessage;
      // message = t('工作目录格式不正确');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  _validateWorkingDir(dir: string) {
    return this._validateDir(dir, t('工作目录格式不正确'));
  },

  _validateLogDir(dir: string) {
    return this._validateDir(dir, t('日志目录格式不正确'));
  },

  validateDir(dir: string, cId: string, v_key: string, errorMessage: string) {
    return async (dispatch, getState: GetState) => {
      const containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cId),
        result = validateWorkloadActions._validateDir(dir, errorMessage);

      containers[cIndex][v_key] = result;
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },

  validateWorkingDir(dir: string, cId: string) {
    return this.validateDir(dir, cId, 'v_workingDir', t('工作目录格式不正确'));
  },

  validateLogDir(dir: string, cId: string) {
    return this.validateDir(dir, cId, 'v_logDir', t('日志目录格式不正确'));
  },

  /** 校验容器的健康检查的端口 */
  _validateHealthPort(port: string) {
    let reg = /^\d+$/,
      status = 0,
      message = '';

    if (!port) {
      status = 2;
      message = t('检查端口不能为空');
    } else if (!reg.test(port)) {
      status = 2;
      message = t('检查端口格式有误');
    } else if (+port - 0 < 1 || +port - 0 > 65535) {
      status = 2;
      message = t('检查端口范围必须在1~65535之间');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateHealthPort(cKey: string, hType: string) {
    return async (dispatch, getState: GetState) => {
      const containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cKey),
        result = validateWorkloadActions._validateHealthPort(containers[cIndex]['healthCheck'][hType]['port']);

      containers[cIndex]['healthCheck'][hType]['v_port'] = result;
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },

  /** 校验容器的健康执行命令 */
  _validateHealthCmd(cmd: string) {
    let status = 0,
      message = '';

    if (!cmd) {
      status = 2;
      message = t('执行命令不能为空');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateHealthCmd(cKey: string, hType: string) {
    return async (dispatch, getState: GetState) => {
      const containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cKey),
        result = validateWorkloadActions._validateHealthCmd(containers[cIndex]['healthCheck'][hType]['cmd']);

      containers[cIndex]['healthCheck'][hType]['v_cmd'] = result;
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },

  /** 校验容器启动延时时间 */
  _validateHealthDelayTime(delayTime: number) {
    let reg = /^\d+$/,
      status = 0,
      message = '';

    if (!delayTime && delayTime !== 0) {
      status = 2;
      message = t('启动延时不能为空');
    } else if (!reg.test(delayTime + '')) {
      status = 2;
      message = t('启动延时格式有误');
    } else if (delayTime - 0 < 0 || delayTime - 0 > 60) {
      status = 2;
      message = t('启动延时范围必须在0~60s之间');
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  validateHealthDelayTime(cKey: string, hType: string) {
    return async (dispatch, getState: GetState) => {
      const containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cKey),
        result = validateWorkloadActions._validateHealthDelayTime(
          +containers[cIndex]['healthCheck'][hType]['delayTime']
        );

      containers[cIndex]['healthCheck'][hType]['v_delayTime'] = result;
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },

  /** 校验容器的响应延迟时间 */
  _validateHealthTimeOut(timeOut: number) {
    let reg = /^\d+$/,
      status = 0,
      message = '';

    if (!timeOut && timeOut !== 0) {
      status = 2;
      message = t('响应超时不能为空');
    } else if (!reg.test(timeOut + '')) {
      status = 2;
      message = t('响应超时格式有误');
    } else if (timeOut - 0 < 2 || timeOut - 0 > 60) {
      status = 2;
      message = t('响应超时范围必须在2~60s之间');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateHealthTimeOut(cKey: string, hType: string) {
    return async (dispatch, getState: GetState) => {
      const containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cKey),
        result = validateWorkloadActions._validateHealthTimeOut(+containers[cIndex]['healthCheck'][hType]['timeOut']);

      containers[cIndex]['healthCheck'][hType]['v_timeOut'] = result;
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },

  /** 容器的间隔时间 */
  _validateHealthIntervalTime(intervalTime: number, timeOut: number) {
    let reg = /^\d+$/,
      status = 0,
      message = '';

    if (!intervalTime && intervalTime !== 0) {
      status = 2;
      message = t('间隔时间不能为空');
    } else if (!reg.test(intervalTime + '')) {
      status = 2;
      message = t('间隔时间格式有误');
    } else if (intervalTime - 0 <= timeOut || intervalTime - 0 > 300) {
      status = 2;
      message = t('间隔时间范围必须大于响应超时，并不超过300s');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateHealthIntervalTime(cKey: string, hType: string) {
    return async (dispatch, getState: GetState) => {
      const containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cKey),
        result = validateWorkloadActions._validateHealthIntervalTime(
          +containers[cIndex]['healthCheck'][hType]['intervalTime'],
          +containers[cIndex]['healthCheck'][hType]['timeOut']
        );

      containers[cIndex]['healthCheck'][hType]['v_intervalTime'] = result;
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },

  /** 校验容器的健康阈值 */
  _validateHealthThreshold(healthThreshold: number) {
    let reg = /^\d+$/,
      status = 0,
      message = '';

    if (!healthThreshold && healthThreshold !== 0) {
      status = 2;
      message = t('健康阈值不能为空');
    } else if (!reg.test(healthThreshold + '')) {
      status = 2;
      message = t('健康阈值格式有误');
    } else if (healthThreshold - 0 < 1 || healthThreshold - 0 > 10) {
      status = 2;
      message = t('健康阈值范围必须在1~10之间');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateHealthThreshold(cKey: string, hType: string) {
    return async (dispatch, getState: GetState) => {
      const containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cKey),
        result = validateWorkloadActions._validateHealthThreshold(
          +containers[cIndex]['healthCheck'][hType]['healthThreshold']
        );

      containers[cIndex]['healthCheck'][hType]['v_healthThreshold'] = result;
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },

  /** 校验容器的不健康阈值 */
  _validateUnhealthThreshold(unhealthThreshold: number) {
    let reg = /^\d+$/,
      status = 0,
      message = '';

    if (!unhealthThreshold && unhealthThreshold !== 0) {
      status = 2;
      message = t('不健康阈值不能为空');
    } else if (!reg.test(unhealthThreshold + '')) {
      status = 2;
      message = t('不健康阈值格式有误');
    } else if (unhealthThreshold - 0 < 1 || unhealthThreshold - 0 > 10) {
      status = 2;
      message = t('不健康阈值范围必须在1~10之间');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateUnHealthThreshold(cKey: string, hType: string) {
    return async (dispatch, getState: GetState) => {
      const containers: ContainerItem[] = cloneDeep(getState().subRoot.workloadEdit.containers),
        cIndex = containers.findIndex(c => c.id === cKey),
        result = validateWorkloadActions._validateUnhealthThreshold(
          +containers[cIndex]['healthCheck'][hType]['unhealthThreshold']
        );

      containers[cIndex]['healthCheck'][hType]['v_unhealthThreshold'] = result;
      dispatch({
        type: ActionType.W_UpdateContainers,
        payload: containers
      });
    };
  },

  /** 校验健康检查 */
  _validateHealthCheck(check: HealthCheckItem) {
    let result = true;
    result =
      result &&
      validateWorkloadActions._validateHealthDelayTime(check.delayTime).status === 1 &&
      validateWorkloadActions._validateHealthTimeOut(check.timeOut).status === 1 &&
      validateWorkloadActions._validateHealthIntervalTime(check.intervalTime, check.timeOut).status === 1 &&
      validateWorkloadActions._validateHealthThreshold(check.healthThreshold).status === 1 &&
      validateWorkloadActions._validateUnhealthThreshold(check.unhealthThreshold).status === 1;

    if (check.checkMethod === 'methodCmd') {
      result = result && validateWorkloadActions._validateHealthCmd(check.cmd).status === 1;
    } else {
      result = result && validateWorkloadActions._validateHealthPort(check.port).status === 1;
    }
    return result;
  },

  validateHealthCheck(container: ContainerItem, hType: string) {
    return async (dispatch, getState: GetState) => {
      const cKey = container.id + '';
      dispatch(validateWorkloadActions.validateHealthPort(cKey, hType));
      dispatch(validateWorkloadActions.validateHealthDelayTime(cKey, hType));
      dispatch(validateWorkloadActions.validateHealthTimeOut(cKey, hType));
      dispatch(validateWorkloadActions.validateHealthIntervalTime(cKey, hType));
      dispatch(validateWorkloadActions.validateHealthCmd(cKey, hType));
      dispatch(validateWorkloadActions.validateHealthThreshold(cKey, hType));
      dispatch(validateWorkloadActions.validateUnHealthThreshold(cKey, hType));
    };
  },

  /** 判断是否能新增容器 */
  _canAddContainer(container: ContainerItem, volumes: VolumeItem[]) {
    let state =
      container.v_name.status === 2 ||
      container.v_registry.status === 2 ||
      container.v_workingDir.status === 2 ||
      container.v_arg.status === 2 ||
      container.v_cmd.status === 2;

    // 如果有数据卷，则判断是有已经都挂载了
    if (volumes.length) {
      container.mounts.forEach(mount => {
        state = state || mount.v_mountPath.status === 2 || mount.v_volume.status === 2;
      });
    }

    // 判断cpu/内存限制的推荐值
    container.cpuLimit.forEach(limit => {
      state = state || limit.v_value.status === 2;
    });
    container.memLimit.forEach(limit => {
      state = state || limit.v_value.status === 2;
    });

    // 如果有环境变量，则判断环境变量是否ok
    if (container.envItems.length) {
      container.envItems.forEach(env => {
        state =
          state ||
          env.v_name.status === 2 ||
          env.v_configMapName.status === 2 ||
          env.v_configMapDataKey.status === 2 ||
          env.v_secretName.status === 2 ||
          env.v_secretDataKey.status === 2;
      });
    }

    // 如果有健康检查，则校验健康检查
    if (container.healthCheck.isOpenLiveCheck) {
      const liveCheck = container.healthCheck.liveCheck;
      state =
        state ||
        liveCheck.v_port.status === 2 ||
        liveCheck.v_delayTime.status === 2 ||
        liveCheck.v_timeOut.status === 2 ||
        liveCheck.v_intervalTime.status === 2 ||
        liveCheck.v_healthThreshold.status === 2 ||
        liveCheck.v_unhealthThreshold.status === 2;
    }

    if (container.healthCheck.isOpenReadyCheck) {
      const readyCheck = container.healthCheck.readyCheck;
      state =
        state ||
        readyCheck.v_port.status === 2 ||
        readyCheck.v_delayTime.status === 2 ||
        readyCheck.v_timeOut.status === 2 ||
        readyCheck.v_intervalTime.status === 2 ||
        readyCheck.v_healthThreshold.status === 2 ||
        readyCheck.v_unhealthThreshold.status === 2;
    }

    return !state && !!container.name && !!container.registry;
  },

  /** 校验容器的所有参数是否已经全部ok */
  _validateContainer(container: ContainerItem, volumes: VolumeItem[], containers: ContainerItem[]) {
    let result = true;

    result =
      result &&
      validateWorkloadActions._validateContainerName(container.name, containers).status === 1 &&
      validateWorkloadActions._validateRegistrySelection(container.registry).status === 1 &&
      validateWorkloadActions._validateGpuCoreLimit(+container.gpuCore).status === 1 &&
      validateWorkloadActions._validateGpuMemLimit(+container.gpuMem).status === 1 &&
      validateWorkloadActions._validateAllCpuLimit(container.cpuLimit) &&
      validateWorkloadActions._validateAllMemLimit(container.memLimit) &&
      validateWorkloadActions._validateWorkingDir(container.workingDir).status === 1;

    if (container.envItems.length) {
      result =
        result &&
        validateWorkloadActions._validateAllEnvName(container.envItems) &&
        validateWorkloadActions._validateAllNewEnvItemValue(container.envItems);
    }

    // 校验挂载业务
    const filters = volumes.filter(v => {
      // 判断当前的挂载项是否还存在，是否已经被删除
      return v.name;
    });
    if (filters.length) {
      if (container.mounts.length) {
        result =
          result &&
          validateWorkloadActions._validateAllVolumeMountPath(container.mounts) &&
          validateWorkloadActions._validateAllVolumeMount(container.mounts, volumes);
      }
    }

    // 校验健康检查
    if (container.healthCheck.isOpenLiveCheck) {
      result = result && validateWorkloadActions._validateHealthCheck(container.healthCheck.liveCheck);
    }

    if (container.healthCheck.isOpenReadyCheck) {
      result = result && validateWorkloadActions._validateHealthCheck(container.healthCheck.readyCheck);
    }

    return result;
  },

  validateContainer(container: ContainerItem) {
    return async (dispatch, getState: GetState) => {
      if (container) {
        let isHealthCheckOk = true;

        const cKey = container.id + '';
        dispatch(validateWorkloadActions.validateContainerName(container.name, cKey));
        dispatch(validateWorkloadActions.validateRegistrySelection(container.registry, cKey));
        dispatch(validateWorkloadActions.validateAllCpuLimit(container));
        dispatch(validateWorkloadActions.validateGpuCoreLimit(container.gpuCore + '', cKey));
        dispatch(validateWorkloadActions.validateGpuMemLimit(container.gpuMem + '', cKey));
        dispatch(validateWorkloadActions.validateAllMemLimit(container));
        dispatch(validateWorkloadActions.validateWorkingDir(container.workingDir, cKey));
        dispatch(validateWorkloadActions.validateLogDir(container.logDir, cKey));

        // 校验环境变量是否都ok
        if (container.envItems.length) {
          dispatch(validateWorkloadActions.validateAllEnvName(container));
          dispatch(validateWorkloadActions.validateAllNewEnvItemValue(container));
        }

        // 校验挂载项是否都ok
        if (container.mounts.length) {
          dispatch(validateWorkloadActions.validateAllVolumeMount(container));
          dispatch(validateWorkloadActions.validateAllVolumeMountPath(container));
        }

        // 校验健康检查
        if (container.healthCheck.isOpenLiveCheck) {
          dispatch(validateWorkloadActions.validateHealthCheck(container, 'liveCheck'));
          isHealthCheckOk =
            isHealthCheckOk && validateWorkloadActions._validateHealthCheck(container.healthCheck.liveCheck);
        }

        if (container.healthCheck.isOpenReadyCheck) {
          dispatch(validateWorkloadActions.validateHealthCheck(container, 'readyCheck'));
          isHealthCheckOk =
            isHealthCheckOk && validateWorkloadActions._validateHealthCheck(container.healthCheck.readyCheck);
        }

        // 这里对高级设置的内容进行检查，如果有错误，则需要展开告诉用户，高级设置当中存在错误
        if (
          validateWorkloadActions._validateWorkingDir(container.workingDir).status === 2 ||
          isHealthCheckOk === false
        ) {
          dispatch(workloadEditActions.modifyAdvancedSettingValidate(true, cKey));
        } else {
          dispatch(workloadEditActions.modifyAdvancedSettingValidate(false, cKey));
        }
      }
    };
  },

  /** 校验hpa触发策略的类型是否正确 */
  _validateHpaType(type: string, metricArr: HpaMetrics[]) {
    let status = 0,
      message = '';

    if (metricArr.filter(item => item.type === type).length > 1) {
      status = 2;
      message = t('触发策略不可重复');
    } else if (type === 'cpuUtilization' && metricArr.find(item => item.type === 'cpuAverage')) {
      status = 2;
      message = t('CPU利用率和CPU使用量不可同时设置');
    } else if (type === 'cpuAverage' && metricArr.find(item => item.type === 'cpuUtilization')) {
      status = 2;
      message = t('CPU利用率和CPU使用量不可同时设置');
    } else if (type === 'memoryUtilization' && metricArr.find(item => item.type === 'memoryAverage')) {
      status = 2;
      message = t('内存利用率和内存使用量不可同时设置');
    } else if (type === 'memoryAverage' && metricArr.find(item => item.type === 'memoryUtilization')) {
      status = 2;
      message = t('内存利用率和内存使用量不可同时设置');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateHpaType(type: string, mId: string) {
    return async (dispatch, getState: GetState) => {
      const metricArr: HpaMetrics[] = cloneDeep(getState().subRoot.workloadEdit.metrics),
        result = validateWorkloadActions._validateHpaType(type, metricArr),
        mIndex = metricArr.findIndex(item => item.id === mId);

      metricArr[mIndex]['v_type'] = result;
      dispatch({
        type: ActionType.W_UpdateMetrics,
        payload: metricArr
      });
    };
  },

  _valdiateAllHpaType(metricArr: HpaMetrics[]) {
    let result = true;
    metricArr.forEach(item => {
      result = result && validateWorkloadActions._validateHpaType(item.type, metricArr).status === 1;
    });
    return result;
  },

  validateAllHpaType() {
    return async (dispatch, getState: GetState) => {
      getState().subRoot.workloadEdit.metrics.forEach(item => {
        dispatch(validateWorkloadActions.validateHpaType(item.type, item.id + ''));
      });
    };
  },

  /** 校验阈值是否正确 */
  _validateHpaValue(value: string, type: string, containers: ContainerItem[]) {
    let status = 0,
      message = '',
      reg = /^\d+|\d+\.\d+$/;

    if (!value) {
      status = 2;
      message = t('阈值不能为空');
    } else if (!reg.test(value)) {
      status = 2;
      message = t('阈值必须要数值');
    } else if (+value < 0) {
      status = 2;
      message = t('阈值须大于等于0');
    } else if (
      type === 'cpuUtilization' &&
      containers.filter(container => container.cpuLimit.find(item => item.type === 'request').value === '').length > 0
    ) {
      status = 2;
      message = t('设置CPU利用率需要设置CPU Request');
    } else if (
      type === 'memoryUtilization' &&
      containers.filter(container => container.memLimit.find(item => item.type === 'request').value === '').length > 0
    ) {
      status = 2;
      message = t('设置内存利用率需要设置内存 Request');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateHpaValue(value: string, mId: string) {
    return async (dispatch, getState: GetState) => {
      const { containers, metrics } = getState().subRoot.workloadEdit;
      const metricArr: HpaMetrics[] = cloneDeep(metrics),
        mIndex = metricArr.findIndex(item => item.id === mId),
        result = validateWorkloadActions._validateHpaValue(value, metricArr[mIndex].type, containers);

      metricArr[mIndex]['v_value'] = result;
      dispatch({
        type: ActionType.W_UpdateMetrics,
        payload: metricArr
      });
    };
  },

  _validateAllHpaValue(metricArr: HpaMetrics[], containers: ContainerItem[]) {
    let result = true;
    metricArr.forEach(item => {
      result = result && validateWorkloadActions._validateHpaValue(item.value, item.type, containers).status === 1;
    });
    return result;
  },

  validateAllHpaValue() {
    return async (dispatch, getState: GetState) => {
      getState().subRoot.workloadEdit.metrics.forEach(item => {
        dispatch(validateWorkloadActions.validateHpaValue(item.value, item.id + ''));
      });
    };
  },

  /** 校验最小实例数 */
  _validateMinReplicas(replicas: string) {
    let status = 0,
      message = '',
      reg = /^\d+$/;

    if (!replicas) {
      status = 2;
      message = t('实例数不能为空');
    } else if (!reg.test(replicas)) {
      status = 2;
      message = t('实例数必须为整数');
    } else if (+replicas < 1) {
      status = 2;
      message = t('实例数最小为1');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateMinReplicas() {
    return async (dispatch, getState: GetState) => {
      const { minReplicas } = getState().subRoot.workloadEdit,
        result = validateWorkloadActions._validateMinReplicas(minReplicas);

      dispatch({
        type: ActionType.WV_MinReplicas,
        payload: result
      });
    };
  },

  /** 校验最大实例数 */
  _validateMaxReplicas(replicas: string, min: string) {
    let status = 0,
      message = '',
      reg = /^\d+$/;

    if (!replicas) {
      status = 2;
      message = t('实例数不能为空');
    } else if (!reg.test(replicas)) {
      status = 2;
      message = t('实例数必须为整数');
    } else if (min && +replicas < +min) {
      status = 2;
      message = t('最大实例数需大于最小实例数');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateMaxReplicas() {
    return async (dispatch, getState: GetState) => {
      const { maxReplicas, minReplicas } = getState().subRoot.workloadEdit,
        result = validateWorkloadActions._validateMaxReplicas(maxReplicas, minReplicas);

      dispatch({
        type: ActionType.WV_MaxReplicas,
        payload: result
      });
    };
  },

  /** 校验cronHpa的crontab是否正确 */
  validateCronTab(mId: string) {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      const metricArr: CronMetrics[] = cloneDeep(getState().subRoot.workloadEdit.cronMetrics),
        mIndex = metricArr.findIndex(item => item.id === mId),
        result = validateWorkloadActions._validateCronSchedule(metricArr[mIndex].crontab);

      metricArr[mIndex]['v_crontab'] = result;
      dispatch({
        type: ActionType.W_UpdateCronMetrics,
        payload: metricArr
      });
    };
  },

  validateAllCronTab() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      getState().subRoot.workloadEdit.cronMetrics.forEach(item => {
        dispatch(validateWorkloadActions.validateCronTab(item.id + ''));
      });
    };
  },

  /** 校验cronHpa的目标实例数是否正确 */
  _validateCronTargetReplicas(replicas: string) {
    let status = 1,
      message = '',
      reg = /\d+/;

    if (!replicas) {
      status = 2;
      message = t('目标实例数不能为空');
    } else if (!reg.test(replicas)) {
      status = 2;
      message = t('目标实例数只能为整数');
    } else if (+replicas < 0) {
      status = 2;
      message = t('目标实例数需大于等于0');
    }

    return { status, message };
  },

  validateCronTargetReplicas(mId: string) {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      const metricArr: CronMetrics[] = cloneDeep(getState().subRoot.workloadEdit.cronMetrics),
        mIndex = metricArr.findIndex(item => item.id === mId),
        result = validateWorkloadActions._validateCronTargetReplicas(metricArr[mIndex].targetReplicas);

      metricArr[mIndex]['v_targetReplicas'] = result;
      dispatch({
        type: ActionType.W_UpdateCronMetrics,
        payload: metricArr
      });
    };
  },

  validateAllCronTargetReplicas() {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      getState().subRoot.workloadEdit.cronMetrics.forEach(item => {
        dispatch(validateWorkloadActions.validateCronTargetReplicas(item.id + ''));
      });
    };
  },

  /** 校验当前的imagePullSecrets */
  _validateImagePullSecret(secretName, list: ImagePullSecrets[]) {
    let status = 0,
      message = '';

    if (!secretName) {
      status = 2;
      message = t('请选择dockercfg类型的Secret');
    } else if (list.filter(item => item.secretName === secretName).length > 1) {
      status = 2;
      message = t('不可选择相同的Secret');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateImagePullSecret(secretName: string, sId: string) {
    return async (dispatch, getState: GetState) => {
      const newList: ImagePullSecrets[] = cloneDeep(getState().subRoot.workloadEdit.imagePullSecrets),
        sIndex = newList.findIndex(item => item.id === sId),
        result = validateWorkloadActions._validateImagePullSecret(secretName, newList);

      newList[sIndex]['v_secretName'] = result;

      dispatch({
        type: ActionType.ImagePullSecrets,
        payload: newList
      });
    };
  },

  /** 校验整个workload表单是否正确 */
  _validateWorkloadEdit(workloadEdit: WorkloadEdit, serviceEdit: ServiceEdit) {
    let result = true;

    result =
      result &&
      validateWorkloadActions._validateWorkloadName(workloadEdit.workloadName).status === 1 &&
      validateWorkloadActions._validateWorkloadDesp(workloadEdit.description).status === 1 &&
      validateWorkloadActions._validateAllWorkloadLabelValue(workloadEdit.workloadLabels) &&
      validateWorkloadActions._validateAllWorkloadLabelKey(workloadEdit.workloadLabels) &&
      validateWorkloadActions._validateNamespace(workloadEdit.namespace).status === 1 &&
      validateWorkloadActions._validateAllVolumeName(workloadEdit.volumes) &&
      validateWorkloadActions._validateAllPvcSelection(workloadEdit.volumes) &&
      validateWorkloadActions._validateAllNfsPath(workloadEdit.volumes) &&
      validateWorkloadActions._validateAllHostPath(workloadEdit.volumes) &&
      validateWorkloadActions._validateAllVolumeIsMounted(workloadEdit.volumes, workloadEdit.containers).allIsMounted;

    if (workloadEdit.workloadAnnotations.length) {
      result =
        result &&
        validateWorkloadActions._validateAllWorkloadAnnotationsKey(workloadEdit.workloadAnnotations) &&
        validateWorkloadActions._validateAllWorkloadAnnotationsValue(workloadEdit.workloadAnnotations);
    }

    workloadEdit.containers.forEach(c => {
      result = result && validateWorkloadActions._validateContainer(c, workloadEdit.volumes, workloadEdit.containers);
    });

    // 判断当前的workload的类型
    const isCronJob = workloadEdit.workloadType === 'cronjob',
      isJob = workloadEdit.workloadType === 'job',
      isDeployment = workloadEdit.workloadType === 'deployment',
      isStatefulset = workloadEdit.workloadType === 'statefulset',
      isTapp = workloadEdit.workloadType === 'tapp';
    // workloadType是 cronjob，则需要校验执行策略
    if (isCronJob) {
      result = result && validateWorkloadActions._validateCronSchedule(workloadEdit.cronSchedule).status === 1;
    }

    // workloadType为 cronjob 或者 job的时候，需要校验 重复次数 和 并行数
    if (isCronJob || isJob) {
      result =
        result &&
        validateWorkloadActions._validateJobCompletion(+workloadEdit.completion).status === 1 &&
        validateWorkloadActions._validateJobParallel(+workloadEdit.parallelism).status === 1;
    }

    // 如果当前是deployment 并且实例更新方式为autoScale
    if ((isDeployment || isTapp) && workloadEdit.scaleType === 'autoScale') {
      result =
        result &&
        validateWorkloadActions._valdiateAllHpaType(workloadEdit.metrics) &&
        validateWorkloadActions._validateAllHpaValue(workloadEdit.metrics, workloadEdit.containers) &&
        validateWorkloadActions._validateMinReplicas(workloadEdit.minReplicas).status === 1 &&
        validateWorkloadActions._validateMaxReplicas(workloadEdit.maxReplicas, workloadEdit.minReplicas).status === 1;
    }

    // 这里是同时创建服务的时候，需要校验Service相关的信息
    if ((isDeployment || isStatefulset) && workloadEdit.isCreateService) {
      result = result && validateServiceActions._validateUpdateServiceAccessEdit(serviceEdit);
    }

    //如果设置了节点亲和性
    if (workloadEdit.nodeAffinityType === affinityType.rule) {
      result = result && validateWorkloadActions._validateAllNodeAffinityRule(workloadEdit.nodeAffinityRule);
    } else if (workloadEdit.nodeAffinityType === affinityType.node) {
      result = result && workloadEdit.computer.selections.length !== 0;
    }

    return result;
  },

  validateWorkloadEdit() {
    return async (dispatch, getState: GetState) => {
      let { subRoot } = getState(),
        { workloadEdit } = subRoot,
        { containers, workloadType, isCreateService, workloadAnnotations } = workloadEdit;

      dispatch(validateWorkloadActions.validateWorkloadName());
      dispatch(validateWorkloadActions.validateWorkloadDesp());
      dispatch(validateWorkloadActions.validateAllWorkloadLabelKey());
      dispatch(validateWorkloadActions.validateAllWorkloadLabelValue());
      dispatch(validateWorkloadActions.validateNamespace());

      // 校验workloadAnnotataions
      if (workloadAnnotations.length) {
        dispatch(validateWorkloadActions.validateAllWorkloadAnnotationsKey());
        dispatch(validateWorkloadActions.validateAllWorkloadAnnotationsValue());
      }

      // 数据卷的相关校验
      dispatch(validateWorkloadActions.validateAllVolumeName());
      dispatch(validateWorkloadActions.validateAllNfsPath());
      dispatch(validateWorkloadActions.validateAllHosPath());
      dispatch(validateWorkloadActions.validateAllPvcSelection());
      dispatch(validateWorkloadActions.validateAllVolumeIsMounted());

      // 校验容器的编辑是否都正确
      containers.forEach(c => {
        dispatch(validateWorkloadActions.validateContainer(c));
      });

      // 判断当前的workload的类型
      const isCronJob = workloadType === 'cronjob',
        isJob = workloadType === 'job',
        isDeployment = workloadEdit.workloadType === 'deployment',
        isStatefulset = workloadEdit.workloadType === 'statefulset',
        isTapp = workloadEdit.workloadType === 'tapp';

      // workloadType是 cronjob，则需要校验执行策略
      if (isCronJob) {
        dispatch(validateWorkloadActions.validateCronSchedule());
      }

      // workloadType为 cronjob 或者 job的时候，需要校验 重复次数 和 并行数
      if (isCronJob || isJob) {
        dispatch(validateWorkloadActions.validateJobCompletion());
        dispatch(validateWorkloadActions.validateJobParallel());
      }

      // 如果当前是deployment或者Tapp 并且实例更新方式为autoScale
      if ((isDeployment || isTapp) && workloadEdit.scaleType === 'autoScale') {
        dispatch(validateWorkloadActions.validateAllHpaType());
        dispatch(validateWorkloadActions.validateAllHpaValue());
        dispatch(validateWorkloadActions.validateMinReplicas());
        dispatch(validateWorkloadActions.validateMaxReplicas());
      }

      // 这里是同时创建服务的时候，需要校验Service相关的信息
      if ((isDeployment || isStatefulset || isTapp) && isCreateService) {
        dispatch(validateServiceActions.validateUpdateServiceAccessEdit());
      }

      if (workloadEdit.nodeAffinityType === affinityType.rule) {
        dispatch(validateWorkloadActions.validateAllNodeAffinityRule());
      } else if (workloadEdit.nodeAffinityType === affinityType.node) {
        dispatch(validateWorkloadActions.validateNodeAffinitySelector());
      }
    };
  },

  /** ========================== start 这里是 更新镜像相关的校验 =================================== */
  _validateMinReadySeconds(seconds: string) {
    let status = 0,
      message = '',
      reg = /^\d+$/;

    if (!seconds) {
      status = 2;
      message = t('更新间隔不能为空值');
    } else if (!reg.test(seconds)) {
      status = 2;
      message = t('更新间隔必须为自然数');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateMinReadySeconds() {
    return async (dispatch, getState: GetState) => {
      const { minReadySeconds } = getState().subRoot.workloadEdit,
        result = validateWorkloadActions._validateMinReadySeconds(minReadySeconds);

      dispatch({
        type: ActionType.WV_MinReadySeconds,
        payload: result
      });
    };
  },

  /** 校验批量的大小 */
  _validateBatchSize(size: string) {
    let status = 0,
      message = '',
      reg = /^(100|[1-9]?\d(\.\d\d?\d?)?)%$|^\d+$/;

    if (!reg.test(size)) {
      status = 2;
      message = t('数值格式不正确，必须为0、正整数或者正百分数');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  validateBatchSize() {
    return async (dispatch, getState: GetState) => {
      const { batchSize } = getState().subRoot.workloadEdit,
        result = validateWorkloadActions._validateBatchSize(batchSize);

      dispatch({
        type: ActionType.WV_BatchSize,
        payload: result
      });
    };
  },

  /** 校验maxSurge */
  validateMaxSurge() {
    return async (dispatch, getState: GetState) => {
      const { maxSurge } = getState().subRoot.workloadEdit,
        result = validateWorkloadActions._validateBatchSize(maxSurge);

      dispatch({
        type: ActionType.WV_MaxSurge,
        payload: result
      });
    };
  },

  /** 校验 maxUnavailable */
  validateMaxUnavaiable(noZero?: boolean) {
    return async (dispatch, getState: GetState) => {
      const { maxUnavailable, workloadType } = getState().subRoot.workloadEdit;
      const isTapp = workloadType === 'tapp';
      let result;
      if (isTapp) {
        result = validateWorkloadActions._validateMaxUnavaiableForTapp(maxUnavailable, noZero);
      } else {
        result = validateWorkloadActions._validateBatchSize(maxUnavailable);
      }

      dispatch({
        type: ActionType.WV_MaxUnavailable,
        payload: result
      });
    };
  },

  _validateMaxUnavaiableForTapp(size: string, noZero?: boolean) {
    let reg = /^\d+$/,
      status = 0,
      message = '';
    if (!size) {
      status = 2;
      message = '数值不能为空';
    } else if (noZero && (size === '0' || !reg.test(size))) {
      status = 2;
      message = '数值不正确，必须为正整数';
    } else if (!reg.test(size)) {
      status = 2;
      message = '数值格式不正确，必须为0或者正整数';
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },

  /** 校验 partition的合法性 */
  validatePartition() {
    return async (dispatch, getState: GetState) => {
      const { partition } = getState().subRoot.workloadEdit,
        result = validateWorkloadActions._validateBatchSize(partition);

      dispatch({
        type: ActionType.WV_Partition,
        payload: result
      });
    };
  },

  /** 校验更新的所有选项是否合法 */
  _validateUpdateRegistryEdit(workloadEdit: WorkloadEdit) {
    const {
      minReadySeconds,
      resourceUpdateType,
      workloadType,
      partition,
      rollingUpdateStrategy,
      maxSurge,
      maxUnavailable,
      batchSize,
      containers
    } = workloadEdit;

    let result = true;
    if (resourceUpdateType === 'RollingUpdate') {
      // 只有statefulset没有 minReadySeconds
      if (workloadType !== 'statefulset') {
        result = result && validateWorkloadActions._validateMinReadySeconds(minReadySeconds).status === 1;
      }

      if (workloadType === 'deployment') {
        if (rollingUpdateStrategy === 'userDefined') {
          result =
            result &&
            validateWorkloadActions._validateBatchSize(maxSurge).status === 1 &&
            validateWorkloadActions._validateBatchSize(maxUnavailable).status === 1;
        } else {
          result = result && validateWorkloadActions._validateBatchSize(batchSize).status === 1;
        }
      } else if (workloadType === 'statefulset') {
        result = result && validateWorkloadActions._validateBatchSize(partition).status === 1;
      }
    }
    if (workloadType === 'tapp') {
      result = result && validateWorkloadActions._validateMaxUnavaiableForTapp(maxUnavailable, true).status === 1;
    }
    containers.forEach(container => {
      result = result && validateWorkloadActions._validateRegistrySelection(container.registry).status === 1;
    });

    return result;
  },

  validateUpdateRegistryEdit() {
    return async (dispatch, getState: GetState) => {
      const { rollingUpdateStrategy, resourceUpdateType, containers, workloadType } = getState().subRoot.workloadEdit;

      const isStatefulset = workloadType === 'statefulset';
      const isTapp = workloadType === 'tapp';

      if (resourceUpdateType === 'RollingUpdate') {
        !isStatefulset && dispatch(validateWorkloadActions.validateMinReadySeconds());
        // 如果当前滚动更新资源为deployment
        if (workloadType === 'deployment') {
          if (rollingUpdateStrategy === 'userDefined') {
            dispatch(validateWorkloadActions.validateMaxSurge());
            dispatch(validateWorkloadActions.validateMaxUnavaiable());
          } else {
            dispatch(validateWorkloadActions.validateBatchSize());
          }
        } else if (workloadType === 'statefulset') {
          dispatch(validateWorkloadActions.validatePartition());
        }
      }

      if (isTapp) {
        dispatch(validateWorkloadActions.validateMaxUnavaiable(true));
      }

      containers.forEach(container => {
        dispatch(validateWorkloadActions.validateRegistrySelection(container.registry, container.id + ''));
      });
    };
  },

  /** ========================== end 这里是 更新镜像相关的校验 =================================== */

  /** ========================== start 这里是 更新实例数量相关 的校验 =================================== */
  _validatePodNumEdit(workloadEdit: WorkloadEdit) {
    let result = true;

    if (workloadEdit.scaleType === 'autoScale') {
      result =
        result &&
        validateWorkloadActions._valdiateAllHpaType(workloadEdit.metrics) &&
        validateWorkloadActions._validateAllHpaValue(workloadEdit.metrics, workloadEdit.containers) &&
        validateWorkloadActions._validateMinReplicas(workloadEdit.minReplicas).status === 1 &&
        validateWorkloadActions._validateMaxReplicas(workloadEdit.maxReplicas, workloadEdit.minReplicas).status === 1;
    }
    return result;
  },

  validatePodNumEdit() {
    return async (dispatch, getState: GetState) => {
      const { scaleType } = getState().subRoot.workloadEdit;

      // 只有hpa的东西需要校验
      if (scaleType === 'autoScale') {
        dispatch(validateWorkloadActions.validateAllHpaType());
        dispatch(validateWorkloadActions.validateAllHpaValue());
        dispatch(validateWorkloadActions.validateMinReplicas());
        dispatch(validateWorkloadActions.validateMaxReplicas());
      }
    };
  },

  /** ========================== end 这里是 更新实例数量 的校验 =================================== */

  /** ========================== start 这里是 校验node节点亲和性规则相关的校验 =================================== */
  _validateNodeAffinityRuleKey(key: string) {
    let reg = /^([A-Za-z0-9][-A-Za-z0-9_./]*)?[A-Za-z0-9]$/,
      status = 0,
      message = '';
    if (!key) {
      status = 2;
      message = t('标签名不能为空');
    } else if (key.length > 63) {
      status = 2;
      message = t('标签名长度不能超过63个字符');
    } else if (!reg.test(key)) {
      status = 2;
      message = t('标签格式不正确');
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },
  _validateNodeAffinityRuleValue(value: string, operator: string) {
    let reg = /^([A-Za-z0-9][-A-Za-z0-9_./]*)?[A-Za-z0-9]$/,
      regNum = /^[0-9]*$/,
      status = 0,
      message = '';
    if (operator !== 'Exists' && operator !== 'DoesNotExist') {
      if (!value) {
        status = 2;
        message = t('自定义规则不能为空');
      } else {
        const valueArray = value.split(';');
        if (operator === 'Lt' || operator === 'Gt') {
          if (valueArray.length !== 1) {
            status = 2;
            message = t('Gt和Lt操作符只支持一个value值');
          } else {
            if (!regNum.test(valueArray[0])) {
              status = 2;
              message = t('Gt和Lt操作符value值格式必须为数字');
            } else {
              status = 1;
              message = '';
            }
          }
        } else {
          for (let i = 0; i < valueArray.length; ++i) {
            if (!valueArray[i]) {
              status = 2;
              message = t('标签名不能为空');
            } else if (valueArray[i].length > 63) {
              status = 2;
              message = t('标签名长度不能超过63个字符');
            } else if (!reg.test(valueArray[i])) {
              status = 2;
              message = t('标签格式不正确');
            } else {
              status = 1;
              message = '';
            }
            if (status === 2) {
              break;
            }
          }
        }
      }
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  },
  validateNodeAffinityRuleKey(type: string, id: string) {
    return async (dispatch, getState: GetState) => {
      let { nodeAffinityRule } = getState().subRoot.workloadEdit,
        requiredExecution = cloneDeep(nodeAffinityRule.requiredExecution),
        preferredExecution = cloneDeep(nodeAffinityRule.preferredExecution),
        result;
      if (type === 'preferred') {
        const preferredMatchExpressions = preferredExecution[0].preference.matchExpressions,
          index = preferredMatchExpressions.findIndex(e => e.id === id);
        result = validateWorkloadActions._validateNodeAffinityRuleKey(preferredMatchExpressions[index].key);
        preferredExecution[0].preference.matchExpressions[index].v_key = result;
      } else if (type === 'required') {
        const requiredMatchExpressions = requiredExecution[0].matchExpressions,
          index = requiredMatchExpressions.findIndex(e => e.id === id);
        result = validateWorkloadActions._validateNodeAffinityRuleKey(requiredMatchExpressions[index].key);
        requiredExecution[0].matchExpressions[index].v_key = result;
      }
      dispatch({
        type: ActionType.W_UpdateNodeAffinityRule,
        payload: Object.assign({}, nodeAffinityRule, {
          requiredExecution,
          preferredExecution
        })
      });
    };
  },
  validateNodeAffinityRuleValue(type: string, id: string) {
    return async (dispatch, getState: GetState) => {
      let { nodeAffinityRule } = getState().subRoot.workloadEdit,
        requiredExecution = cloneDeep(nodeAffinityRule.requiredExecution),
        preferredExecution = cloneDeep(nodeAffinityRule.preferredExecution),
        result;
      if (type === 'preferred') {
        const preferredMatchExpressions = preferredExecution[0].preference.matchExpressions,
          index = preferredMatchExpressions.findIndex(e => e.id === id);
        result = validateWorkloadActions._validateNodeAffinityRuleValue(
          preferredMatchExpressions[index].values,
          preferredMatchExpressions[index].operator
        );
        preferredExecution[0].preference.matchExpressions[index].v_values = result;
      } else if (type === 'required') {
        const requiredMatchExpressions = requiredExecution[0].matchExpressions,
          index = requiredMatchExpressions.findIndex(e => e.id === id);
        result = validateWorkloadActions._validateNodeAffinityRuleValue(
          requiredMatchExpressions[index].values,
          requiredMatchExpressions[index].operator
        );
        requiredExecution[0].matchExpressions[index].v_values = result;
      }
      dispatch({
        type: ActionType.W_UpdateNodeAffinityRule,
        payload: Object.assign({}, nodeAffinityRule, {
          requiredExecution,
          preferredExecution
        })
      });
    };
  },

  validateAllNodeAffinityRule() {
    return async (dispatch, getState: GetState) => {
      let { nodeAffinityRule } = getState().subRoot.workloadEdit,
        { requiredExecution, preferredExecution } = nodeAffinityRule;
      requiredExecution[0].matchExpressions.forEach(rule => {
        dispatch(validateWorkloadActions.validateNodeAffinityRuleKey('required', rule.id + ''));
        dispatch(validateWorkloadActions.validateNodeAffinityRuleValue('required', rule.id + ''));
      });
      preferredExecution[0].preference.matchExpressions.forEach(rule => {
        dispatch(validateWorkloadActions.validateNodeAffinityRuleKey('preferred', rule.id + ''));
        dispatch(validateWorkloadActions.validateNodeAffinityRuleValue('preferred', rule.id + ''));
      });
    };
  },

  _validateAllNodeAffinityRule(nodeAffinityRule: AffinityRule) {
    const { requiredExecution, preferredExecution } = nodeAffinityRule;
    let result = true;
    requiredExecution[0].matchExpressions.forEach(rule => {
      result =
        result &&
        validateWorkloadActions._validateNodeAffinityRuleKey(rule.key).status === 1 &&
        validateWorkloadActions._validateNodeAffinityRuleValue(rule.values, rule.operator).status === 1;
    });
    preferredExecution[0].preference.matchExpressions.forEach(rule => {
      result =
        result &&
        validateWorkloadActions._validateNodeAffinityRuleKey(rule.key).status === 1 &&
        validateWorkloadActions._validateNodeAffinityRuleValue(rule.values, rule.operator).status === 1;
    });
    return result;
  },

  validateNodeAffinitySelector() {
    return async (dispatch, getState: GetState) => {
      const { computer } = getState().subRoot.workloadEdit;
      let status = 0,
        message = '';
      if (computer.selections.length === 0) {
        status = 2;
        message = t('选择节点不能为空');
      } else {
        status = 1;
        message = '';
      }
      dispatch({
        type: ActionType.WV_NodeSelector,
        payload: {
          status,
          message
        }
      });
    };
  }
  /** ========================== end 这里是 校验node节点亲和性规则相关的校验 =================================== */
};
