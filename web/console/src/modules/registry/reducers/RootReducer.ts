import { combineReducers } from 'redux';

import {
  createFFListReducer,
  generateWorkflowReducer,
  reduceToPayload,
  generateFetcherReducer
} from '@tencent/ff-redux';

import * as ActionType from '../constants/ActionType';
import { InitApiKey, InitRepo, InitChart, InitImage, Default_D_URL } from '../constants/Config';
import { router } from '../router';

export const RootReducer = combineReducers({
  route: router.getReducer(),

  /** 访问凭证相关 */
  apiKey: createFFListReducer('apiKey'),

  createApiKey: generateWorkflowReducer({
    actionType: ActionType.CreateApiKey
  }),

  apiKeyCreation: reduceToPayload(ActionType.UpdateApiKeyCreation, InitApiKey),

  deleteApiKey: generateWorkflowReducer({
    actionType: ActionType.DeleteApiKey
  }),

  toggleKeyStatus: generateWorkflowReducer({
    actionType: ActionType.ToggleKeyStatus
  }),

  /** 镜像仓库相关 */
  repo: createFFListReducer('repo'),

  createRepo: generateWorkflowReducer({
    actionType: ActionType.CreateRepo
  }),

  repoCreation: reduceToPayload(ActionType.UpdateRepoCreation, InitRepo),

  deleteRepo: generateWorkflowReducer({
    actionType: ActionType.DeleteRepo
  }),

  /** 镜像相关 */
  image: createFFListReducer('image'),

  createImage: generateWorkflowReducer({
    actionType: ActionType.CreateImage
  }),

  imageCreation: reduceToPayload(ActionType.UpdateImageCreation, InitImage),

  deleteImage: generateWorkflowReducer({
    actionType: ActionType.DeleteImage
  }),

  dockerRegistryUrl: generateFetcherReducer({
    actionType: ActionType.FetchDockerRegUrl,
    initialData: Default_D_URL
  }),

  /** chart group */
  chart: createFFListReducer('chart'),

  chartIns: createFFListReducer('chartIns'),

  createChart: generateWorkflowReducer({
    actionType: ActionType.CreateChart
  }),

  chartCreation: reduceToPayload(ActionType.UpdateChartCreation, InitChart),

  deleteChart: generateWorkflowReducer({
    actionType: ActionType.DeleteChart
  })
});
