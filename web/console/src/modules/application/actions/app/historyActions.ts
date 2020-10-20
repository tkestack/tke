import {
  extend,
  createFFObjectActions,
  uuid,
  generateWorkflowActionCreator,
  OperationTrigger,
  isSuccessWorkflow
} from '@tencent/ff-redux';
import { RootState, AppHistory, AppHistoryFilter, History } from '../../models';
import * as ActionTypes from '../../constants/ActionTypes';
import * as WebAPI from '../../WebAPI';
import { router } from '../../router';
type GetState = () => RootState;
const tips = seajs.require('tips');

/**
 * 获取历史列表
 */

const fetchHistoryActions = createFFObjectActions<AppHistory, AppHistoryFilter>({
  actionName: ActionTypes.AppHistory,
  fetcher: async (query, getState: GetState, fetchOptions, dispatch: Redux.Dispatch) => {
    let response = await WebAPI.fetchAppHistory(query.filter);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().appHistory;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    let histories = [];
    if (record.data) {
      histories = record.data.spec.histories.map(h => {
        return Object.assign({}, h, { id: uuid(), involvedObject: record.data });
      });
    }
    dispatch({
      type: ActionTypes.HistoryList,
      payload: {
        histories: histories
      }
    });
  }
});

/**
 * 回滚应用
 */
const rollbackAppWorkflow = generateWorkflowActionCreator<History, void>({
  actionType: ActionTypes.RollbackApp,
  workflowStateLocator: (state: RootState) => state.appRollbackWorkflow,
  operationExecutor: WebAPI.rollbackApp,
  after: {
    [OperationTrigger.Done]: (dispatch, getState: GetState) => {
      if (isSuccessWorkflow(getState().appRollbackWorkflow)) {
        let { route } = getState();
        router.navigate({ mode: 'list', sub: 'app' }, route.queries);
      }
      /** 结束工作流 */
      dispatch(historyActions.rollbackAppWorkflow.reset());
    }
  }
});

const restActions = {
  rollbackAppWorkflow
};

export const historyActions = extend({}, fetchHistoryActions, restActions);
