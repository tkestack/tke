import { uuid } from '@tencent/qcloud-lib';
import { ContainerLogs, MetadataItem, ContainerFilePathItem } from '../models';
import { initValidator } from '../../common/models';
import { ResourceListMapForContainerLog, ResourceListMapForPodLog } from './Config';
import { ResourceTarget } from '../models/Resource';
import { LogDaemonSetStatus } from '../models/LogDaemonset';

/** 地域的初始化信息 */
export const initRegionInfo = {
  name: '广州',
  value: 1,
  area: '华南地区'
};

/** 初始化指定容器日志的编辑项 */
export const initWorkloadList = (initData: any) => {
  return {
    deployment: initData,
    statefulset: initData,
    daemonset: initData,
    job: initData,
    cronjob: initData
  };
};

export const initContainerInputOption: ContainerLogs = {
  id: uuid(),
  namespaceSelection: 'default',
  v_namespaceSelection: initValidator,
  collectorWay: 'container',
  workloadType: 'deployment',
  status: 'editing',
  workloadSelection: initWorkloadList([]),
  workloadList: initWorkloadList([]),
  workloadListFetch: initWorkloadList(false),
  v_workloadSelection: initValidator
};

/** metadata初始变量 */
export const initMetadata: MetadataItem = {
  id: uuid(),
  metadataKey: '',
  v_metadataKey: initValidator,
  metadataValue: '',
  v_metadataValue: initValidator
};

/**容器文件路径初始变量 */
export const initContainerFilePath: ContainerFilePathItem = {
  id: uuid(),
  containerName: '',
  containerFilePath: '',
  v_containerName: initValidator,
  v_containerFilePath: initValidator
};

export const initContainerFileWorkloadType: string = ResourceListMapForContainerLog[0].value;

/**默认选择容器标准输出*/
export const initResourceTarget: ResourceTarget = {
  isForContainerFile: false,
  isForContainerLogs: true
};

export const initLogDaemonsetStatus: LogDaemonSetStatus = {
  phase: '',
  reason: ''
};
