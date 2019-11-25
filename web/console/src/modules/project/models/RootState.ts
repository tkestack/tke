import { Resource } from './../../common/models/Resource';
import { WorkflowState } from '@tencent/qcloud-redux-workflow';
import { Region, RegionFilter } from '../../common/models';
import {
  Project,
  ProjectFilter,
  ProjectEdition,
  Manager,
  Namespace,
  NamespaceFilter,
  NamespaceEdition,
  NamespaceOperator,
  Cluster,
  ClusterFilter,
  ManagerFilter
} from './';
import { RouteState } from '../../../../helpers';
import { ListModel } from '@tencent/redux-list';

type ProjectWorkflow = WorkflowState<Project, void>;
type ProjectEditWorkflow = WorkflowState<ProjectEdition, void>;
type NamespaceWorkflow = WorkflowState<Namespace, NamespaceOperator>;
type NamespaceEditWorkflow = WorkflowState<NamespaceEdition, NamespaceOperator>;

export interface RootState {
  /** 路由 */
  route?: RouteState;

  project?: ListModel<Project, ProjectFilter>;

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

  namespace?: ListModel<Namespace, NamespaceFilter>;

  /** Namespace编辑参数 */
  namespaceEdition?: NamespaceEdition;

  /** 创建业务工作流 */
  createNamespace?: NamespaceEditWorkflow;

  /** 创建业务工作流 */
  editNamespaceResourceLimit?: NamespaceEditWorkflow;

  /** 删除业务工作流 */
  deleteNamespace?: NamespaceWorkflow;

  /** 地域列表 */
  region?: ListModel<Region, RegionFilter>;

  /** 集群列表*/
  cluster?: ListModel<Cluster, ClusterFilter>;

  /** 负责人列表 */
  manager?: ListModel<Manager, ManagerFilter>;

  /** 设置管理员*/
  modifyAdminstrator?: ProjectEditWorkflow;

  /**当前管理员 */
  adminstratorInfo?: Resource;
}
