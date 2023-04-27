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
/** ========================= start FFRedux的相关配置 ======================== */
export const FFReduxActionName = {
  OPENADDON: 'openAddon',
  CLUSTER: 'cluster'
};
/** ========================= end FFRedux的相关配置 ======================== */

/** 判断当前集群状态下能否创建日志收集器 */
export const canCreateLogStash = ['Running', 'Scaling'];

/** 判断当前logdaemonset状态下能否启动创建按钮*/
export const canCreateLogStashInLogDaemonset = ['Running', '404'];

/**判断当前状态下能否获取loglist */
export const canFetchLogList = ['Running'];

/** 日志采集器的状态 */
export const collectorStatus = {
  Initializing: {
    text: 'Initializing',
    classname: 'text-weak'
  },
  Running: {
    text: 'Running',
    classname: 'text-success'
  },
  Checking: {
    text: 'Checking',
    classname: 'text-danger'
  },
  Reinitializing: {
    text: 'Reinitializing',
    classname: 'text-warning'
  },
  Fail: {
    text: 'Fail',
    classname: 'text-danger'
  }
};

/** 日志类型映射 */
export const logModeMap = {
  'container-log': '容器标准输出',
  'host-log': '指定主机文件',
  'pod-log': '指定容器文件'
};

/** 输入源类型映射 */
export const inputTypeMap = {
  'container-log': 'container',
  'host-log': 'node',
  'pod-log': 'containerFile'
};

/** 日志采集的类型 */
export const logModeList = {
  container: {
    value: 'container',
    name: '容器标准输出'
  },
  // containerFile: {
  //   value: 'containerFile',
  //   name: '容器文件路径'
  // },
  node: {
    value: 'node',
    name: '节点文件路径'
  }
};

/**
 * pre: 指定容器日志
 * 日志源类型的值
 */
export const originModeList = [
  {
    value: 'selectAll',
    name: '所有容器'
  },
  {
    value: 'selectOne',
    name: '指定容器'
  }
];

/**
 * pre: 消费端类型
 */
export const consumerModeList = [
  // {
  //   value: 'kafka',
  //   name: 'Kafka'
  // },
  {
    value: 'es',
    name: 'Elasticsearch'
  }
];
/** 输出源类型映射 */
export const outputTypeMap = {
  kafka: 'kafka',
  ckafka: 'kafka',
  cls: 'cls',
  elasticsearch: 'es'
};

export const ClsLogSetSupportRegionList = [1, 4, 6, 8, 16];
export const clsRegionMap = {
  1: 'ap-guangzhou',
  4: 'ap-shanghai',
  8: 'ap-beijing',
  16: 'ap-chengdu',
  6: 'na-toronto'
};

export const CkafkaSupportRegionList = [1, 4, 5, 7, 8, 9, 11, 15, 16, 19, 21, 25];

export const ResourceListMapForContainerLog = [
  {
    name: 'Deployment',
    value: 'deployment'
  },
  {
    name: 'Daemonset',
    value: 'daemonset'
  },
  {
    name: 'Statefulset',
    value: 'statefulset'
  },
  {
    name: 'Cronjob',
    value: 'cronjob'
  },
  {
    name: 'Job',
    value: 'job'
  },
  {
    name: 'TApp',
    value: 'tapp'
  }
];

export const ResourceListMapForPodLog = [
  {
    name: 'Deployment',
    value: 'deployment'
  },
  {
    name: 'Daemonset',
    value: 'daemonset'
  },
  {
    name: 'Statefulset',
    value: 'statefulset'
  },
  {
    name: 'Cronjob',
    value: 'cronjob'
  }
];

export const HOST_LOG_INPUT_PATH_PREFIX = '/run/containerd/io.containerd.runtime.v2.task/k8s.io/*/rootfs';
