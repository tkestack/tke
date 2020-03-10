import { WorkflowState, FetcherState, FFListModel } from '@tencent/ff-redux';
import { ApiKey, ApiKeyFilter, ApiKeyCreation } from './ApiKey';
import { Repo, RepoFilter, RepoCreation } from './Repo';
import { Image, ImageFilter, ImageCreation } from './Image';
import { Chart, ChartFilter, ChartIns, ChartInsFilter, ChartCreation } from './Chart';
import { RouteState } from '../../../../helpers';

type ApiKeyWorkflow = WorkflowState<ApiKey, void>;
type ApiKeyCreateWorkflow = WorkflowState<ApiKeyCreation, void>;
type RepoWorkflow = WorkflowState<Repo, void>;
type RepoCreateWorkflow = WorkflowState<RepoCreation, void>;
type ImageWorkflow = WorkflowState<Image, void>;
type ImageCreateWorkflow = WorkflowState<ImageCreation, void>;
type ChartWorkflow = WorkflowState<Chart, void>;
type ChartCreateWorkflow = WorkflowState<ChartCreation, void>;

export interface RootState {
  /** 路由 */
  route?: RouteState;

  /** -------- 访问凭证相关 --------- */

  apiKey?: FFListModel<ApiKey, ApiKeyFilter>;

  /** ApiKey 创建编辑参数 */
  apiKeyCreation?: ApiKeyCreation;

  /** 创建业务工作流 */
  createApiKey?: ApiKeyCreateWorkflow;

  /** 删除工作流 */
  deleteApiKey?: ApiKeyWorkflow;

  /** 启用禁用工作流 */
  toggleKeyStatus?: ApiKeyWorkflow;

  /** --------- 仓库管理相关 --------- */

  repo?: FFListModel<Repo, RepoFilter>;

  /** 创建仓库表单参数 */
  repoCreation?: RepoCreation;

  /** 创建仓库工作流 */
  createRepo?: RepoCreateWorkflow;

  /** 删除仓库工作流 */
  deleteRepo?: RepoWorkflow;

  /** --------- 镜像相关 --------- */

  image?: FFListModel<Image, ImageFilter>;

  /** 创建仓库表单参数 */
  imageCreation?: ImageCreation;

  /** 创建仓库工作流 */
  createImage?: ImageCreateWorkflow;

  /** 删除仓库工作流 */
  deleteImage?: ImageWorkflow;

  /** docker registry */
  dockerRegistryUrl?: FetcherState<string>;

  /** -------- chart group ----- */

  chart?: FFListModel<Chart, ChartFilter>;

  chartIns?: FFListModel<ChartIns, ChartInsFilter>;

  /** 创建仓库表单参数 */
  chartCreation?: ChartCreation;

  /** 创建仓库工作流 */
  createChart?: ChartCreateWorkflow;

  /** 删除仓库工作流 */
  deleteChart?: ChartWorkflow;
}
