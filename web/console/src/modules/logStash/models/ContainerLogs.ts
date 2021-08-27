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
import { Identifiable } from '@tencent/ff-redux';

import { Resource } from '../../cluster/models';
import { Validation } from '../../common/models';

export interface ContainerLogs extends Identifiable {
  /** 当前的namespace */
  namespaceSelection: string;
  v_namespaceSelection: Validation;

  /** 采集的方式：全部容器、指定工作负载、指定Labels */
  collectorWay: string;

  /** 当前的workload的类型 */
  workloadType: string;

  /** 当前的状态 edited 非编辑状态, editing: 编辑状态 */
  status: string;

  /** workloadList */
  workloadList: WorkloadType<Resource>;

  /** 判断workloadList是否已经拉取过 */
  workloadListFetch: WorkloadType<any>;

  /** 选择workload的集合 */
  workloadSelection: WorkloadType<string>;
  v_workloadSelection: Validation;
}

export interface WorkloadSelection {
  value: string;
  label: string;
}

export interface WorkloadType<T> {
  deployment: T[];
  statefulset: T[];
  daemonset: T[];
  job: T[];
  cronjob: T[];
  tapp: T[];
}
