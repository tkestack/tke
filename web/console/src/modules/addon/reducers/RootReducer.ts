import { combineReducers } from 'redux';

import { createFFListReducer, generateWorkflowReducer, reduceToPayload } from '@tencent/ff-redux';

import * as ActionType from '../constants/ActionType';
import { FFReduxActionName } from '../constants/Config';
import { router } from '../router';
import { AddonEditReducer } from './AddonEditReducer';

export const RootReducer = combineReducers({
  route: router.getReducer(),

  region: createFFListReducer(FFReduxActionName.REGION),

  cluster: createFFListReducer(FFReduxActionName.CLUSTER),

  clusterVersion: reduceToPayload(ActionType.ClusterVersion, '1.16'),

  openAddon: createFFListReducer(FFReduxActionName.OPENADDON),

  addon: createFFListReducer(FFReduxActionName.ADDON),

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
