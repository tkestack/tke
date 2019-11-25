import { FetcherState } from '@tencent/qcloud-redux-fetcher';
import { QueryState } from '@tencent/qcloud-redux-query';
import { WorkflowState } from '@tencent/qcloud-redux-workflow';
import { RecordSet } from '@tencent/qcloud-lib';
import { ResourceInfo } from '../../common/models';
import {
  WorkloadEdit,
  SecretEdit,
  ConfigMapEdit,
  ResourceEventOption,
  ResourceLogOption,
  ComputerState,
  ResourceOption,
  CreateResource,
  SubRouter,
  SubRouterFilter,
  ResourceDetailState,
  ServiceEdit,
  NamespaceEdit,
  LbcfEdit
} from './index';
import { DetailResourceOption } from './DetailResourceOption';
import { AllocationRatioEdition } from './AllocationRatioEdition';
import { AddonStatus } from './Addon';

type ResourceModifyWorkflow = WorkflowState<CreateResource, number | any>;

export interface SubRootState {
  /** 节点列表 */
  computerState?: ComputerState;

  /** 超售比 */
  clusterAllocationRatioEdition?: AllocationRatioEdition;

  updateClusterAllocationRatio?: ResourceModifyWorkflow;

  /** 二级菜单栏配置列表查询 */
  subRouterQuery?: QueryState<SubRouterFilter>;

  /** 二级菜单栏配置列表 */
  subRouterList?: FetcherState<RecordSet<SubRouter>>;

  /** 当前的模式 create | update | resource */
  mode?: string;

  /** 创建多种resource资源的操作流程 */
  applyResourceFlow?: ResourceModifyWorkflow;

  /**创建多种resource (使用不同的接口)的操作流程 */
  applyDifferentInterfaceResourceFlow?: ResourceModifyWorkflow;
  /** 创建resource资源的操作流 */
  modifyResourceFlow?: ResourceModifyWorkflow;

  modifyMultiResourceWorkflow?: ResourceModifyWorkflow;

  /** 删除resource资源的操作流 */
  deleteResourceFlow?: ResourceModifyWorkflow;

  /** 更新Service的访问方式、更新Ingress的转发配置等的操作流 */
  updateResourcePart?: ResourceModifyWorkflow;

  updateMultiResource?: ResourceModifyWorkflow;

  /** 当前的请求资源名称 */
  resourceName?: string;

  /** resourcrInfo */
  resourceInfo?: ResourceInfo;

  detailResourceOption?: DetailResourceOption;

  /** 通用resource 数据结构 */
  resourceOption?: ResourceOption;

  /** resource detail详情 */
  resourceDetailState?: ResourceDetailState;

  /** editService */
  serviceEdit?: ServiceEdit;

  /** editNamespace */
  namespaceEdit?: NamespaceEdit;

  /** editResource */
  workloadEdit?: WorkloadEdit;

  /** editSecret */
  secretEdit?: SecretEdit;

  lbcfEdit?: LbcfEdit;

  /** editConfigMap */
  cmEdit?: ConfigMapEdit;

  /** resourcelog的相关配置 */
  resourceLogOption?: ResourceLogOption;

  /** resourceEvent的相关配置 */
  resourceEventOption?: ResourceEventOption;

  /** 是否需要进行命名空间的拉取 */
  isNeedFetchNamespace?: boolean;

  /** 使用已有lb白名单 */
  isNeedExistedLb?: boolean;

  addons?: AddonStatus;
}
