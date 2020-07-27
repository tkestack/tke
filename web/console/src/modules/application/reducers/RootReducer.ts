import { combineReducers } from 'redux';

import {
  createFFListReducer,
  generateWorkflowReducer,
  reduceToPayload,
  generateFetcherReducer,
  createFFObjectReducer
} from '@tencent/ff-redux';

import { Cluster, Namespace, ProjectNamespace, Project } from '../models';
import * as ActionTypes from '../constants/ActionTypes';
import { router } from '../router';
import { createValidatorReducer } from '@tencent/ff-validator';
import { AppValidateSchema } from '../constants/AppValidateConfig';
import { initAppCreationState, initAppEditorState, initResourceList, initHistoryList } from '../constants/initState';

export const RootReducer = combineReducers({
  route: router.getReducer(),

  /** 集群 */
  /** listActions.selectByValue依赖于valueField */
  clusterList: createFFListReducer(
    ActionTypes.ClusterList,
    '',
    (x: Cluster) => x.spec.displayName,
    (x: Cluster) => x.metadata.name
  ),
  /** 命名空间 */
  namespaceList: createFFListReducer(
    ActionTypes.NamespaceList,
    '',
    (x: Namespace) => x.metadata.name,
    (x: Namespace) => x.metadata.name
  ),
  projectNamespaceList: createFFListReducer(
    ActionTypes.ProjectNamespaceList,
    '',
    (x: ProjectNamespace) => x.metadata.name,
    (x: ProjectNamespace) => x.spec.clusterName + '/' + x.spec.namespace
  ),

  /** 模板 */
  chartList: createFFListReducer(ActionTypes.ChartList, null, null, null, {
    query: {
      paging: {
        pageSize: 9
      }
    }
  }),
  chartInfo: createFFObjectReducer(ActionTypes.ChartInfo),
  chartGroupList: createFFListReducer(ActionTypes.ChartGroupList),

  /** 业务 */
  projectList: createFFListReducer(
    ActionTypes.ProjectList,
    '',
    (x: Project) => x.metadata.name,
    (x: Project) => x.metadata.name
  ),

  /** 应用 */
  appList: createFFListReducer(ActionTypes.AppList),
  appCreation: reduceToPayload(ActionTypes.UpdateAppCreationState, initAppCreationState),
  appEditor: reduceToPayload(ActionTypes.UpdateAppEditorState, initAppEditorState),
  appValidator: createValidatorReducer(AppValidateSchema),
  appAddWorkflow: generateWorkflowReducer({
    actionType: ActionTypes.AddApp
  }),
  appUpdateWorkflow: generateWorkflowReducer({
    actionType: ActionTypes.UpdateApp
  }),
  appRemoveWorkflow: generateWorkflowReducer({
    actionType: ActionTypes.RemoveApp
  }),
  appResource: createFFObjectReducer(ActionTypes.AppResource),
  resourceList: reduceToPayload(ActionTypes.ResourceList, initResourceList),
  appHistory: createFFObjectReducer(ActionTypes.AppHistory),
  historyList: reduceToPayload(ActionTypes.HistoryList, initHistoryList),
  appRollbackWorkflow: generateWorkflowReducer({
    actionType: ActionTypes.RollbackApp
  })
});
