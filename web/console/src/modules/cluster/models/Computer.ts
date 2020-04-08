import { FetcherState, FFListModel, Identifiable, QueryState, RecordSet, WorkflowState } from '@tencent/ff-redux';

import { CreateResource } from '../../common';
import { Validation } from '../../common/models';
import { Resource, ResourceFilter } from './ResourceOption';

type ComputerWorkflow = WorkflowState<Computer, ComputerOperator>;
type ComputerLabelWorkflow = WorkflowState<ComputerLabelEdition, ComputerOperator>;
type ResourceModifyFlow = WorkflowState<CreateResource, number>;
type ComputerTaintWorkflow = WorkflowState<ComputerTaintEdition, ComputerOperator>;
export interface ComputerOperator {
  /**
   * 集群Id
   */
  clusterId?: string;

  /**
   * 地域
   */
  regionId?: number;

  /**
   * 移出操作方式
   */
  nodeDeleteMode?: string;

  /**
   * node是否进行 Unschedule 的操作
   */
  isUnSchedule?: boolean;
}

export interface ComputerState {
  /** computer的相关配置 */
  computer: FFListModel<Computer, ComputerFilter>;

  /** computer的相关配置 */
  machine: FFListModel<Computer, ComputerFilter>;

  /**创建com工作流 */
  createComputerWorkflow?: ResourceModifyFlow;

  /** unschedule 节点 */
  batchUnScheduleComputer?: ResourceModifyFlow;

  /** turn on scheduling 节点 */
  batchTurnOnSchedulingComputer?: ResourceModifyFlow;

  /** drain the node 驱逐节点 */
  drainComputer?: ComputerWorkflow;

  /** 驱逐节点所包含的pod列表 */
  computerPodList?: FetcherState<RecordSet<Resource>>;

  computerPodQuery?: QueryState<ResourceFilter>;

  /**批量删除 Computer 操作流 */
  deleteComputer?: ResourceModifyFlow;

  /**编辑computer标签 操作流 */
  updateNodeLabel?: ComputerLabelWorkflow;

  /**编辑computer */
  labelEdition?: ComputerLabelEdition;

  /**编辑computerTaint 操作流 */
  updateNodeTaint?: ComputerTaintWorkflow;

  taintEdition?: ComputerTaintEdition;

  isShowMachine?: boolean;

  deleteMachineResouceIns?: Resource;
}

export interface ComputerFilter extends ResourceFilter {}

export interface Computer extends Identifiable {
  /** metadata */
  metadata?: ComputerMetadata;

  /** spec */
  spec?: ComputerSpec;

  /** status */
  status?: ComputerStatus;
}

interface ComputerMetadata {
  /** annotations */
  annotations: {
    [props: string]: any;
  };

  /** 节点的创建时间 */
  creationTimestamp?: string;

  /** 节点的label */
  labels?: {
    [props: string]: any;
  };

  /** 节点的id */
  name: string;

  /**role */
  role: string;

  [props: string]: any;
}

interface ComputerSpec {
  /** 节点的名称 */
  externalID: string;

  /** podCIDR */
  podCIDR: string;

  /**对应mathine的名称master&etcd没有 */
  machineName?: string;

  /**是否可以封锁*/
  unschedulable?: boolean;

  taints: { [props: string]: string }[];
}

interface ComputerStatus {
  /** addresses */
  addresses: Address[];

  /** 可分配的资源 */
  allocatable: Allocatable;

  /** 节点的配置 */
  capacity: Allocatable;

  /** 节点的状态 */
  conditions: Conditions[];

  /** daemonEndPoint */
  daemonEndpoints: {
    kubeletEndpoint: DaemonEndpoints;
  };

  /** 节点当中的镜像 */
  images: Images[];

  /** 节点的信息 */
  nodeInfo: NodeInfo;

  [props: string]: any;
}

interface NodeInfo {
  /** cpu芯片 */
  architecture: string;

  /** bootID */
  bootID: string;

  /** 运行中docker的版本 */
  containerRuntimeVersion: string;

  /** 内核版本 */
  kernelVersion: string;

  /** kube-proxy版本 */
  kubeProxyVersion: string;

  /** kuberlet版本 */
  kubeletVersion: string;

  /** machineID */
  machineID: string;

  /** operatingSystem */
  operatingSystem: string;

  /** osImage */
  osImage: string;

  /** systemUUID */
  systemUUID: string;
}

interface Images {
  names: string[];

  sizeBytes: number;
}

interface DaemonEndpoints {
  /** kube的endpoint */
  kubeletEndpoint: {
    Port: number;
  };
}

interface Conditions {
  /** 最后心跳时间 */
  lastHeartbeatTime: string;

  /** 最后变化时间 */
  lastTransitionTime: string;

  /** 错误信息 */
  message: string;

  /** 原因 */
  reason: string;

  /** 状态 */
  status: string;

  /** 错误类型 */
  type: string;
}

interface Allocatable {
  /** 可分配的cpu */
  cpu: string;

  /** 可分配的mem */
  memory: string;

  /** 可分配的pod */
  pods: string;

  /** 闪存 */
  'ephemeral-storage': string;

  /** 最大 */
  'hugepages-2Mi': string;
}

interface Address {
  /** ip地址 */
  address: string;

  /** ip的类型 */
  type: string;
}

export interface ComputerLabel extends Identifiable {
  key?: string;
  value?: string;
  v_key?: Validation;
  v_value?: Validation;
  disabled?: boolean;
}

export interface ComputerLabelEdition extends Identifiable {
  labels?: ComputerLabel[];

  originLabel?: Object;

  computerName?: string;
}
export interface ComputerTaint extends Identifiable {
  key?: string;
  value?: string;
  v_key?: Validation;
  v_value?: Validation;
  effect?: string;
  disabled?: boolean;
}

export interface ComputerTaintEdition extends Identifiable {
  taints?: ComputerTaint[];
  computerName?: string;
}
