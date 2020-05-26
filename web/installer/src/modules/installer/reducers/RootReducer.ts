import { combineReducers } from 'redux';

import {
    generateFetcherReducer, generateWorkflowReducer, reduceToPayload, ReduxAction, uuid
} from '@tencent/ff-redux';

import { Record } from '../../common/models';
import * as ActionType from '../constants/ActionType';
import { initEdit } from './initState';

export const RootReducer = combineReducers({
  step: reduceToPayload(ActionType.StepNext, 'step1'),

  cluster: generateFetcherReducer<Record<any>>({
    actionType: ActionType.FetchCluster,
    initialData: {
      record: {
        config: {},
        progress: {}
      },
      auth: {
        isAuthorized: true
      }
    }
  }),

  isVerified: reduceToPayload(ActionType.VerifyLicense, -1),

  licenseConfig: reduceToPayload(ActionType.GetLicenseConfig, {}),

  clusterProgress: generateFetcherReducer<Record<any>>({
    actionType: ActionType.FetchProgress,
    initialData: {
      record: {},
      auth: {
        isAuthorized: true
      }
    }
  }),

  editState: (state = Object.assign({}, initEdit, { id: uuid() }), action: any) => {
    if (action.type === ActionType.UpdateEdit) {
      return Object.assign({}, state, action.payload);
    } else {
      return state;
    }
  },

  createCluster: generateWorkflowReducer({
    actionType: ActionType.CreateCluster
  })
});
