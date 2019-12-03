import { GroupFilter } from './Group';
import { RecordSet } from '@tencent/qcloud-lib';
import { FetcherState, FetchState } from '@tencent/qcloud-redux-fetcher';
import { QueryState } from '@tencent/qcloud-redux-query';
import { WorkflowState } from '@tencent/qcloud-redux-workflow';
//import { RouteState } from "@tencent/qcloud-nmc";
import { RouteState } from '../../../../helpers/Router';
import {
  Region,
  RegionFilter,
  ClusterFilter,
  AlarmPolicyFilter,
  AlarmPolicyEdition,
  AlarmPolicy,
  Group,
  Namespace,
  NamespaceFilter,
  Resource,
  ResourceFilter,
  AlarmPolicyOperator
} from './';
import { Cluster } from '../../common';
import { User, UserFilter } from '../../uam/models';
import { ListModel } from '@tencent/redux-list';

type AlarmPolicyOpWorkflow = WorkflowState<AlarmPolicy, AlarmPolicyOperator>;
type AlarmPolicyCreateWorkflow = WorkflowState<AlarmPolicyEdition, AlarmPolicyOperator>;

export interface RootState {
  /**
   * 路由
   */
  route?: RouteState;

  /**
   * 地域查询
   */
  regionQuery?: QueryState<RegionFilter>;

  /**
   * 地域列表
   */
  regionList?: FetcherState<RecordSet<Region>>;

  /**
   * 选择的地域
   */
  regionSelection?: Region;

  cluster?: ListModel<Cluster, ClusterFilter>;

  /**当前集群命名空间 */
  namespaceList?: FetcherState<RecordSet<Namespace>>;

  namespaceQuery?: QueryState<NamespaceFilter>;

  /**当前命名空间下pod列表 */
  workloadList?: FetcherState<RecordSet<Resource>>;

  workloadQuery?: QueryState<ResourceFilter>;

  clusterVersion?: string;

  alarmPolicy?: ListModel<AlarmPolicy, AlarmPolicyFilter>;

  userList?: ListModel<User, UserFilter>;

  /** 当前新建告警 */
  alarmPolicyEdition?: AlarmPolicyEdition;

  /** 创建告警workflow */
  alarmPolicyCreateWorkflow?: AlarmPolicyCreateWorkflow;

  /** 更新告警workflow */
  alarmPolicyUpdateWorkflow?: AlarmPolicyCreateWorkflow;

  /** 删除告警workflow */
  alarmPolicyDeleteWorkflow?: AlarmPolicyOpWorkflow;

  /**详情 */
  alarmPolicyDetail?: AlarmPolicy;

  /**组列表 */
  // groupList?: FetcherState<RecordSet<Group>>;

  channel?: ListModel<Resource, ResourceFilter>;
  template?: ListModel<Resource, ResourceFilter>;
  receiver?: ListModel<Resource, ResourceFilter>;
  receiverGroup?: ListModel<Resource, ResourceFilter>;

  groupQuery?: QueryState<GroupFilter>;

  /** 是否为国际版 */
  isI18n?: boolean;

  // /** namespace列表 */
  // namespaceList?: FetcherState<RecordSet<Namespace>>;

  // /** namespace查询条件 */
  // namespaceQuery?: QueryState<ResourceFilter>;

  /** namespace selection */
  namespaceSelection?: string;

  /** namespacesetQuery */
  projectNamespaceQuery?: QueryState<ResourceFilter>;

  /** namespaceset */
  projectNamespaceList?: FetcherState<RecordSet<Resource>>;

  /** projectList */
  projectList?: any[];

  /** projectSelection */
  projectSelection?: string;
}
