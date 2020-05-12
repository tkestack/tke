import {
    createFFListActions, extend, generateWorkflowActionCreator, isSuccessWorkflow, OperationTrigger
} from '@tencent/ff-redux';

import * as ActionType from '../constants/ActionType';
import { Manager, ManagerFilter, RootState } from '../models';
import { ProjectEdition } from '../models/Project';
import * as WebAPI from '../WebAPI';
import { projectActions } from './projectActions';

type GetState = () => RootState;

const FFModelManagerActions = createFFListActions<Manager, ManagerFilter>({
  actionName: 'manager',
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchUser(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().manager;
  }
});
const restActions = {
  selectManager: (manager: Manager[]) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch(FFModelManagerActions.selects(manager));
    };
  },

  fetchAdminstratorInfo: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let response = await WebAPI.fetchAdminstratorInfo();
      dispatch({
        type: ActionType.FetchAdminstratorInfo,
        payload: response
      });
    };
  },

  modifyAdminstrator: generateWorkflowActionCreator<ProjectEdition, void>({
    actionType: ActionType.ModifyAdminstrator,
    workflowStateLocator: (state: RootState) => state.modifyAdminstrator,
    operationExecutor: WebAPI.modifyAdminstrator,
    after: {
      [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState) => {
        let { modifyAdminstrator, route } = getState();
        if (isSuccessWorkflow(modifyAdminstrator)) {
          dispatch(restActions.modifyAdminstrator.reset());
          dispatch(restActions.fetchAdminstratorInfo());
          dispatch(projectActions.clearEdition());
        }
      }
    }
  }),

  initAdminstrator: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let {
        adminstratorInfo,
        manager: { list }
      } = getState();
      if (adminstratorInfo.spec) {
        let members = adminstratorInfo.spec.administrators
          ? adminstratorInfo.spec.administrators.map(item => {
              let finder = list.data.records.find(manager => manager.name === item);
              if (finder) {
                return finder;
              } else {
                return {
                  name: item,
                  displayName: '用户不存在'
                };
              }
            })
          : [];
        dispatch({
          type: ActionType.UpdateProjectEdition,
          payload: Object.assign({}, getState().projectEdition, {
            members,
            resourceVersion: adminstratorInfo.metadata.resourceVersion,
            id: adminstratorInfo.metadata.name
          })
        });
      }
    };
  }
};

export const managerActions = extend({}, FFModelManagerActions, restActions);
