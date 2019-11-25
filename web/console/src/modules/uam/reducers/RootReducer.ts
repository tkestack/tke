import { combineReducers } from 'redux';
import { generateFetcherReducer } from '@tencent/qcloud-redux-fetcher';
import { generateWorkflowReducer } from '@tencent/qcloud-redux-workflow';
import * as ActionTypes from '../constants/ActionTypes';
import { router } from '../router';
import { reduceToPayload } from '@tencent/qcloud-lib';
import { createListReducer } from '@tencent/redux-list';

export const RootReducer = combineReducers({
  route: router.getReducer(),
  userList: createListReducer(ActionTypes.UserList),
  addUserWorkflow: generateWorkflowReducer({
    actionType: ActionTypes.AddUser
  }),
  removeUserWorkflow: generateWorkflowReducer({
    actionType: ActionTypes.RemoveUser
  }),
  filterUsers: reduceToPayload(ActionTypes.FetchUserByName, []),
  getUser: generateFetcherReducer<Object>({
    actionType: ActionTypes.GetUser,
    initialData: {}
  }),
  updateUser: generateFetcherReducer<Object>({
    actionType: ActionTypes.UpdateUser,
    initialData: {}
  }),

  userStrategyList: createListReducer(ActionTypes.UserStrategyList),

  strategyList: createListReducer(ActionTypes.StrategyList),
  addStrategyWorkflow: generateWorkflowReducer({
    actionType: ActionTypes.AddStrategy
  }),
  removeStrategyWorkflow: generateWorkflowReducer({
    actionType: ActionTypes.RemoveStrategy
  }),
  getStrategy: generateFetcherReducer<Object>({
    actionType: ActionTypes.GetStrategy,
    initialData: {}
  }),
  updateStrategy: generateFetcherReducer<Object>({
    actionType: ActionTypes.UpdateStrategy,
    initialData: {}
  }),

  categoryList: generateFetcherReducer<Object>({
    actionType: ActionTypes.GetCategories,
    initialData: {}
  }),
  associatedUsersList: createListReducer(ActionTypes.GetStrategyAssociatedUsers),
  removeAssociatedUser: generateWorkflowReducer({
    actionType: ActionTypes.RemoveAssociatedUser
  }),
  addAssociatedUser: generateWorkflowReducer({
    actionType: ActionTypes.AddAssociatedUser
  })
});
