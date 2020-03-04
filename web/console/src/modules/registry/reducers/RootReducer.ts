import { createFFListReducer, generateWorkflowReducer } from '@tencent/ff-redux';
import { reduceToPayload } from '@tencent/qcloud-lib';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { combineReducers } from 'redux';
import * as ActionType from '../constants/ActionType';
import { Default_D_URL, InitApiKey, InitImage, InitRepo } from '../constants/Config';
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
  })
});
