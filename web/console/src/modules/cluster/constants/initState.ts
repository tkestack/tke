/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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
import { uuid } from '@tencent/ff-redux';

import { initValidator } from '../../common/models';
import {
  ConfigItems,
  ContainerItem,
  DialogNameEnum,
  DialogState,
  HealthCheck,
  HealthCheckItem,
  HpaMetrics,
  ImagePullSecrets,
  LimitItem,
  MountItem,
  PortMap,
  SecretData,
  Selector,
  VolumeItem,
  WorkloadLabel
} from '../models';
import { CronMetrics } from '../models/WorkloadEdit';
import { BackendType } from './Config';

/** 创建服务，端口映射的初始值 */
export const initPortsMap: PortMap = {
  id: uuid(),
  protocol: 'TCP',
  v_protocol: initValidator,
  targetPort: '',
  v_targetPort: initValidator,
  port: '',
  v_port: initValidator,
  nodePort: '',
  v_nodePort: initValidator
};

/** 创建服务，selector的初始值 */
export const initSelector: Selector = {
  id: uuid(),
  key: '',
  value: '',
  v_key: initValidator,
  v_value: initValidator
};

export const initLbcfBGPort = {
  id: uuid(),
  portNumber: '',
  protocol: 'TCP',
  v_portNumber: initValidator
};

export const initStringArray = {
  id: uuid(),
  value: '',
  v_value: initValidator
};

export const initLbcfBackGroupEdition = {
  onEdit: true,
  id: uuid(),
  name: '',
  v_name: initValidator,
  backgroupType: BackendType.Pods,
  staticAddress: [initStringArray],
  byName: [],
  serviceName: '',
  v_serviceName: initValidator,
  ports: [initLbcfBGPort],
  labels: [initSelector]
};

/** 创建workload，数据卷的初始化 */
export const initVolume: VolumeItem = {
  id: uuid(),
  volumeType: 'emptyDir',
  v_volumeType: initValidator,
  name: '',
  v_name: initValidator,
  hostPathType: 'DirectoryOrCreate',
  hostPath: '',
  v_hostPath: initValidator,
  nfsPath: '',
  v_nfsPath: initValidator,
  configKey: [],
  configName: '',
  secretKey: [],
  secretName: '',
  pvcSelection: '',
  v_pvcSelection: initValidator,
  newPvcName: '',
  pvcEditInfo: {
    accessMode: '',
    storageClassName: '',
    storage: ''
  },
  isMounted: true
};

/** 创建workload的时候，hpa的初始化指标 */
export const initHpaMetrics: HpaMetrics = {
  id: uuid(),
  type: 'cpuUtilization',
  v_type: initValidator,
  value: '',
  v_value: initValidator
};

/** 创建workloa的时候，定时调节的初始化值 */
export const initCronMetrics: CronMetrics = {
  id: uuid(),
  crontab: '',
  v_crontab: initValidator,
  targetReplicas: '',
  v_targetReplicas: initValidator
};

/** 创建 secret */
export const initSecretData: SecretData = {
  id: uuid(),
  keyName: '',
  v_keyName: initValidator,
  value: '',
  v_value: initValidator
};

/** 创建workload，label的初始化 */
export const initSpecificLabel: WorkloadLabel = {
  id: uuid(),
  labelKey: 'k8s-app',
  v_labelKey: initValidator,
  labelValue: '',
  v_labelValue: initValidator
};

export const initWorkloadLabel: WorkloadLabel = {
  id: uuid(),
  labelKey: '',
  v_labelKey: initValidator,
  labelValue: '',
  v_labelValue: initValidator
};

/** 创建workload，annotataionos的初始化 */
export const initWorkloadAnnotataions: WorkloadLabel = {
  id: uuid(),
  labelKey: '',
  v_labelKey: initValidator,
  labelValue: '',
  v_labelValue: initValidator
};

/** 创建workload，configMap的items的初始化 */
export const initConfigMapItem: ConfigItems = {
  id: uuid(),
  configKey: '',
  v_configKey: initValidator,
  path: '',
  v_path: initValidator,
  mode: '0644',
  v_mode: initValidator
};

/** 创建workload，containers的初始化 */
export const initMount: MountItem[] = [
  {
    id: uuid(),
    volume: '',
    v_volume: initValidator,
    mountPath: '',
    v_mountPath: initValidator,
    mountSubPath: '',
    v_mountSubPath: initValidator,
    mode: 'rw',
    v_mode: initValidator
  }
];

/** cpu推荐的初始值 */
const initCpuLimit: LimitItem[] = [
  {
    id: uuid(),
    type: 'request',
    value: '0.25',
    v_value: initValidator
  },
  {
    id: uuid(),
    type: 'limit',
    value: '0.5',
    v_value: initValidator
  }
];

/** 内存推荐值的初始值 */
const initMemLimit: LimitItem[] = [
  {
    id: uuid(),
    type: 'request',
    value: '256',
    v_value: initValidator
  },
  {
    id: uuid(),
    type: 'limit',
    value: '1024',
    v_value: initValidator
  }
];

/** 健康检查具体项初始值 */
const initHealthCheckItem: HealthCheckItem = {
  checkMethod: 'methodTcp',
  port: '',
  v_port: initValidator,
  protocol: 'HTTP',
  path: '/',
  v_path: initValidator,
  cmd: '',
  v_cmd: initValidator,
  delayTime: 0,
  v_delayTime: initValidator,
  timeOut: 2,
  v_timeOut: initValidator,
  intervalTime: 3,
  v_intervalTime: initValidator,
  healthThreshold: 1,
  v_healthThreshold: initValidator,
  unhealthThreshold: 1,
  v_unhealthThreshold: initValidator
};

/** 初始化健康值检查 */
const initHealthCheck: HealthCheck = {
  isOpenLiveCheck: false,
  isOpenReadyCheck: false,
  liveCheck: initHealthCheckItem,
  readyCheck: initHealthCheckItem
};

/** 初始化imagePullSecrets */
export const initImagePullSecrets: ImagePullSecrets = {
  id: uuid(),
  secretName: '',
  v_secretName: initValidator
};

/** 初始化容器的选项 */
export const initContainer: ContainerItem = {
  id: uuid(),
  status: 'editing',
  name: '',
  v_name: initValidator,
  registry: '',
  v_registry: initValidator,
  tag: '',
  mounts: initMount,
  memLimit: initMemLimit,
  cpuLimit: initCpuLimit,
  envItems: [],
  isOpenAdvancedSetting: false,
  isAdvancedError: false,
  gpu: 0,
  gpuCore: '0',
  v_gpuCore: initValidator,
  gpuMem: '0',
  v_gpuMem: initValidator,
  workingDir: '',
  v_workingDir: initValidator,
  cmd: '',
  v_cmd: initValidator,
  arg: '',
  v_arg: initValidator,
  healthCheck: initHealthCheck,
  privileged: false,
  addCapabilities: [],
  dropCapabilities: [],
  imagePullPolicy: 'Always'
};

/**初始化node节点亲和性 */
export const initmatchExpressions = {
  key: '',
  operator: 'In',
  values: '',
  v_key: initValidator,
  v_values: initValidator
};

export const initAffinityRule = {
  requiredExecution: [{ matchExpressions: [Object.assign({}, initmatchExpressions, { id: uuid() })] }],
  preferredExecution: [
    {
      preference: { matchExpressions: [Object.assign({}, initmatchExpressions, { id: uuid() })] },
      weight: 1
    }
  ]
};

/** 各弹窗的状态显示 */
export const initDialogState: DialogState = {
  [DialogNameEnum.clusterStatusDialog]: false,
  [DialogNameEnum.kuberctlDialog]: false,
  [DialogNameEnum.computerStatusDialog]: false
};

export const initClusterCreationState = {
  /**链接集群名字 */
  name: '',
  v_name: initValidator,

  /**apiServer地址 */
  apiServer: '',
  v_apiServer: initValidator,

  /** port */
  port: '',
  v_port: initValidator,

  /**证书 */
  certFile: '',
  v_certFile: initValidator,

  token: '',
  v_token: initValidator,

  jsonData: {},

  currentStep: 1,

  clientCert: '',
  clientKey: '',
  username: '',
  as: '',
  clusternetCertificate: '',
  clusternetPrivatekey: ''
};

export const initAllcationRatioEdition = {
  id: uuid(),
  isUseCpu: false,
  isUseMemory: false,
  cpuRatio: '',
  v_cpuRatio: initValidator,
  memoryRatio: '',
  v_memoryRatio: initValidator
};
