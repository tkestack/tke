import { combineReducers } from 'redux';
import { router } from '../router';
import { FFReduxActionName } from '../constants/Config';
import * as ActionType from '../constants/ActionType';
import { reduceToPayload } from '@tencent/qcloud-lib';
import { AddonEditReducer } from './AddonEditReducer';
import { generateWorkflowReducer, createFFListReducer } from '@tencent/ff-redux';

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
