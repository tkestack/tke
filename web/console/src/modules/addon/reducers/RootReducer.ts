import { combineReducers } from 'redux';
import { router } from '../router';
import { FFReduxActionName } from '../constants/Config';
import * as ActionType from '../constants/ActionType';
import { reduceToPayload } from '@tencent/qcloud-lib';
import { AddonEditReducer } from './AddonEditReducer';
import { generateWorkflowReducer } from '@tencent/qcloud-redux-workflow';
import { createListReducer } from '@tencent/redux-list';

export const RootReducer = combineReducers({
  route: router.getReducer(),

  region: createListReducer(FFReduxActionName.REGION),

  cluster: createListReducer(FFReduxActionName.CLUSTER),

  clusterVersion: reduceToPayload(ActionType.ClusterVersion, '1.16'),

  openAddon: createListReducer(FFReduxActionName.OPENADDON),

  addon: createListReducer(FFReduxActionName.ADDON),

  editAddon: AddonEditReducer,

  modifyResourceFlow: generateWorkflowReducer({
    actionType: ActionType.ModifyResource
  }),

  applyResourceFlow: generateWorkflowReducer({
    actionType: ActionType.ApplyResource
  }),

  deleteResourceFlow: generateWorkflowReducer({
    actionType: ActionType.DeleteResource
  })
});
