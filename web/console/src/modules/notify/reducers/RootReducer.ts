import { combineReducers } from 'redux';
import { reduceToPayload } from '@tencent/qcloud-lib';
import { generateWorkflowReducer } from '@tencent/qcloud-redux-workflow';
import { router } from '../router';
import * as ActionType from '../constants/ActionType';
import { Resource } from '../models/Resource';
import { createListReducer } from '@tencent/redux-list';

interface ResourceFilter {}
export const RootReducer = combineReducers({
  route: router.getReducer(),
  channel: createListReducer<Resource, ResourceFilter>('channel'),
  template: createListReducer<Resource, ResourceFilter>('template'),
  receiver: createListReducer<Resource, ResourceFilter>('receiver'),
  receiverGroup: createListReducer<Resource, ResourceFilter>('receiverGroup'),
  resourceDeleteWorkflow: generateWorkflowReducer({ actionType: ActionType.DeleteResource }),
  modifyResourceFlow: generateWorkflowReducer({ actionType: ActionType.ModifyResource }),

  /**
   * 判断是否为国际版
   */
  isI18n: reduceToPayload(ActionType.isI18n, false)
});
