import {
  extend,
  generateWorkflowActionCreator,
  OperationResult,
  OperationTrigger,
  createFFListActions
} from '@tencent/ff-redux';
import { RootState, App, AppFilter } from '../../models';
import * as ActionTypes from '../../constants/ActionTypes';
import * as WebAPI from '../../WebAPI';
type GetState = () => RootState;

/**
 * 列表操作
 */
const fetchAppActions = createFFListActions<App, AppFilter>({
  actionName: ActionTypes.AppList,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchAppList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().appList;
  },
  onFinish: (record, dispatch: Redux.Dispatch) => {
    let isNotNeedPoll = true;
    if (record.data.recordCount) {
      isNotNeedPoll =
        record.data.records.filter(
          item =>
            !item.status ||
            (item.status['phase'] && item.status['phase'] !== 'Succeeded') ||
            !item.status['releaseStatus'] ||
            item.status['releaseStatus'] !== 'deployed' ||
            (item.status['observedGeneration'] && item.status['observedGeneration'] < item.metadata['generation'])
        ).length === 0;
    }
    if (isNotNeedPoll) {
      dispatch(fetchAppActions.clearPolling());
    }
  }
});

/**
 * 删除应用
 */
const removeAppWorkflow = generateWorkflowActionCreator<App, void>({
  actionType: ActionTypes.RemoveApp,
  workflowStateLocator: (state: RootState) => state.appRemoveWorkflow,
  operationExecutor: WebAPI.deleteApp,
  after: {
    [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
      //在列表页删除的动作，因此直接重新拉取一次数据
      dispatch(listActions.poll({ namespace: getState().appRemoveWorkflow.targets[0].metadata.namespace }));
      /** 结束工作流 */
      dispatch(listActions.removeAppWorkflow.reset());
    }
  }
});

const restActions = {
  removeAppWorkflow,

  /** 轮询操作 */
  poll: (filter: AppFilter) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch(
        listActions.polling({
          delayTime: 3000,
          filter: filter
        })
      );
    };
  }
};

export const listActions = extend({}, fetchAppActions, restActions);
