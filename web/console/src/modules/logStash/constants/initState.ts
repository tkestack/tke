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
import { ContainerFilePathItem, ContainerLogs, MetadataItem } from '../models';
import { LogDaemonSetStatus } from '../models/LogDaemonset';
import { ResourceTarget } from '../models/Resource';
import { ResourceListMapForContainerLog, ResourceListMapForPodLog } from './Config';

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
    cronjob: initData,
    tapp: initData
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
