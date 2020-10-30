import { WorkflowState, FetcherState, FFListModel, FFObjectModel } from '@tencent/ff-redux';
import {
  App,
  AppFilter,
  AppCreation,
  AppEditor,
  AppResource,
  AppResourceFilter,
  ResourceList,
  HistoryList,
  AppHistoryFilter,
  AppHistory,
  History
} from './App';
import { Cluster } from './Cluster';
import { Chart, ChartFilter, ChartInfo, ChartInfoFilter } from './Chart';
import { ChartGroup, ChartGroupFilter } from './ChartGroup';
import { Namespace, NamespaceFilter, ProjectNamespace, ProjectNamespaceFilter } from './Namespace';
import { Project } from './Project';
import { RouteState } from '../../../../helpers';
import { Validation, ValidatorModel } from '@tencent/ff-validator';

type AppWorkflow = WorkflowState<App, void>;

export interface RootState {
  /** 路由 */
  route?: RouteState;

  /** 集群 */
  clusterList?: FFListModel<Cluster>;
  /** 命名空间 */
  namespaceList?: FFListModel<Namespace, NamespaceFilter>;
  projectNamespaceList?: FFListModel<ProjectNamespace, ProjectNamespaceFilter>;

  /** 模板 */
  chartList?: FFListModel<Chart, ChartFilter>;
  chartInfo?: FFObjectModel<ChartInfo, ChartInfoFilter>;
  chartGroupList?: FFListModel<ChartGroup, ChartGroupFilter>;

  /** 业务 */
  projectList?: FFListModel<Project>;

  /** 应用 */
  appList?: FFListModel<App, AppFilter>;
  appCreation?: AppCreation;
  appEditor?: AppEditor;
  appDryRun?: App;
  appValidator?: ValidatorModel;
  appAddWorkflow?: WorkflowState<App, any>;
  appUpdateWorkflow?: WorkflowState<App, any>;
  appRemoveWorkflow?: WorkflowState<App, any>;
  appResource?: FFObjectModel<AppResource, AppResourceFilter>;
  resourceList?: ResourceList;
  appHistory?: FFObjectModel<AppHistory, AppHistoryFilter>;
  historyList?: HistoryList;
  appRollbackWorkflow?: WorkflowState<History, any>;
}
