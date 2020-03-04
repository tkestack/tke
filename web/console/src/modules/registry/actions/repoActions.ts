import {
    createFFListActions, extend, generateWorkflowActionCreator, isSuccessWorkflow, OperationTrigger
} from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';

import * as ActionType from '../constants/ActionType';
import { InitRepo } from '../constants/Config';
import { Repo, RepoFilter, RootState } from '../models';
import { RepoCreation } from '../models/Repo';
import { router } from '../router';
import * as WebAPI from '../WebAPI';

type GetState = () => RootState;

const FFModelRepoActions = createFFListActions<Repo, RepoFilter>({
  actionName: 'repo',
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchRepoList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().repo;
  }
});

const restActions = {
  /** 创建 Repo */
  createRepo: generateWorkflowActionCreator<RepoCreation, void>({
    actionType: ActionType.CreateRepo,
    workflowStateLocator: (state: RootState) => state.createRepo,
    operationExecutor: WebAPI.createRepo,
    after: {
      [OperationTrigger.Done]: (dispatch, getState) => {
        let { createRepo, route } = getState();
        if (isSuccessWorkflow(createRepo)) {
          dispatch(restActions.createRepo.reset());
          dispatch(restActions.clearEdition());
          dispatch(repoActions.fetch());
          let urlParams = router.resolve(route);
          router.navigate(Object.assign({}, urlParams, { sub: 'repo', mode: 'list' }), {});
        }
      }
    }
  }),

  /** 删除 Repo */
  deleteRepo: generateWorkflowActionCreator<Repo, void>({
    actionType: ActionType.DeleteRepo,
    workflowStateLocator: (state: RootState) => state.deleteRepo,
    operationExecutor: WebAPI.deleteRepo,
    after: {
      [OperationTrigger.Done]: (dispatch, getState) => {
        let { deleteRepo, route } = getState();
        if (isSuccessWorkflow(deleteRepo)) {
          dispatch(restActions.deleteRepo.reset());
          dispatch(repoActions.fetch());
        }
      }
    }
  }),

  /** --begin编辑action */
  inputRepoDesc: (value: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.UpdateRepoCreation,
        payload: Object.assign({}, getState().repoCreation, { displayName: value })
      });
    };
  },

  inputRepoName: (value: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.UpdateRepoCreation,
        payload: Object.assign({}, getState().repoCreation, { name: value })
      });
      dispatch(repoActions.validateRepoName(value));
    };
  },

  selectRepoVisibility: (value: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.UpdateRepoCreation,
        payload: Object.assign({}, getState().repoCreation, { visibility: value })
      });
    };
  },

  validateRepoName(value: string) {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let result = repoActions._validateRepoName(value);
      dispatch({
        type: ActionType.UpdateRepoCreation,
        payload: Object.assign({}, getState().repoCreation, { v_name: result })
      });
    };
  },

  _validateRepoName(name: string) {
    let reg = /^[a-z]([-a-z0-9]*[a-z0-9])?$/,
      status = 0,
      message = '';

    if (!name) {
      status = 2;
      message = t('命名空间不能为空');
    } else if (name.length > 63) {
      status = 2;
      message = t('命名空间不能超过63个字符');
    } else if (!reg.test(name)) {
      status = 2;
      message = t('命名空间格式不正确');
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  clearEdition: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.UpdateRepoCreation,
        payload: InitRepo
      });
    };
  }
  /** --end编辑action */
};

export const repoActions = extend({}, FFModelRepoActions, restActions);
