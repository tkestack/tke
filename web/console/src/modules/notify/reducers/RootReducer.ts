import { combineReducers } from 'redux';

import { createFFListReducer, generateWorkflowReducer, reduceToPayload } from '@tencent/ff-redux';

import * as ActionType from '../constants/ActionType';
import { Resource } from '../models/Resource';
import { router } from '../router';

interface ResourceFilter {}
export const RootReducer = combineReducers({
  route: router.getReducer(),
  channel: createFFListReducer<Resource, ResourceFilter>('channel'),
  template: createFFListReducer<Resource, ResourceFilter>('template'),
  receiver: createFFListReducer<Resource, ResourceFilter>('receiver'),
  receiverGroup: createFFListReducer<Resource, ResourceFilter>('receiverGroup'),
  resourceDeleteWorkflow: generateWorkflowReducer({ actionType: ActionType.DeleteResource }),
  modifyResourceFlow: generateWorkflowReducer({ actionType: ActionType.ModifyResource }),

  /**
   * 判断是否为国际版
   */
  isI18n: reduceToPayload(ActionType.isI18n, false)
});
