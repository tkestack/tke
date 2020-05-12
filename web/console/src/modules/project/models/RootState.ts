import { FetcherState, FFListModel, OperationResult, RecordSet, WorkflowState } from '@tencent/ff-redux';

import { RouteState } from '../../../../helpers';
import { Region, RegionFilter, ResourceFilter } from '../../common/models';
import { Resource } from '../../common/models/Resource';
import {
  User,
  Member,
  UserFilter,
  PolicyPlain,
  PolicyFilter,
  Cluster, ClusterFilter, Manager, ManagerFilter, Namespace, NamespaceEdition, NamespaceFilter,
  NamespaceOperator, Project, ProjectEdition, ProjectFilter
} from './index';

type ProjectWorkflow = WorkflowState<Project, void>;
type ProjectEditWorkflow = WorkflowState<ProjectEdition, void>;
type NamespaceWorkflow = WorkflowState<Namespace, NamespaceOperator>;
type NamespaceEditWorkflow = WorkflowState<NamespaceEdition, NamespaceOperator>;
type userWorkflow = WorkflowState<Member, any>;

export interface RootState {
  /** 路由 */
  route?: RouteState;

  project?: FFListModel<Project, ProjectFilter>;

  /** 业务编辑参数 */
  projectEdition?: ProjectEdition;

  /** 创建业务工作流 */
  createProject?: ProjectEditWorkflow;

  /** 编辑业务名称工作流 */
  editProjectName?: ProjectEditWorkflow;

  /** 编辑业务负责人工作流 */
  editProjectManager?: ProjectEditWorkflow;

  /** 编辑业务描述工作流 */
  editProjecResourceLimit?: ProjectEditWorkflow;

  /** 删除业务工作流 */
  deleteProject?: ProjectWorkflow;

  namespace?: FFListModel<Namespace, NamespaceFilter>;

  /** Namespace编辑参数 */
  namespaceEdition?: NamespaceEdition;

  /** 创建业务工作流 */
  createNamespace?: NamespaceEditWorkflow;

  /** 创建业务工作流 */
  editNamespaceResourceLimit?: NamespaceEditWorkflow;

  /** 删除业务工作流 */
  deleteNamespace?: NamespaceWorkflow;

  /** 地域列表 */
  region?: FFListModel<Region, RegionFilter>;

  /** 集群列表*/
  cluster?: FFListModel<Cluster, ClusterFilter>;

  /** 负责人列表 */
  manager?: FFListModel<Manager, ManagerFilter>;

  /** 设置管理员*/
  modifyAdminstrator?: ProjectEditWorkflow;

  /**当前管理员 */
  adminstratorInfo?: Resource;

  /** 用户信息 */
  userList?: FFListModel<User, UserFilter>;
  addUserWorkflow?: userWorkflow;

  /** 关联策略相关，单独设置，不赋予任何场景相关的命名 */
  policyPlainList?: FFListModel<PolicyPlain, PolicyFilter>;
}
