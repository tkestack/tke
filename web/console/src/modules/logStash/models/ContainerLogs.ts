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
