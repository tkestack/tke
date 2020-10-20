import {
  extend,
  generateWorkflowActionCreator,
  OperationResult,
  OperationTrigger,
  createFFListActions
} from '@tencent/ff-redux';
import { RootState, Chart, ChartFilter } from '../../models';
import * as ActionTypes from '../../constants/ActionType';
import * as WebAPI from '../../WebAPI';
type GetState = () => RootState;

/**
 * 列表操作
 */
const fetchChartActions = createFFListActions<Chart, ChartFilter>({
  actionName: ActionTypes.ChartList,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchChartList(query);
    // 此部分若放置在onFinish，会出现lastVersion异步返回的现象，且ui不会更新
    if (response.recordCount > 0) {
      let records = response.records.map(r => {
        if (r.status.versions) {
          let sorted = r.status.versions.sort((a, b) => {
            let oDate1 = new Date(a.timeCreated);
            let oDate2 = new Date(b.timeCreated);
            return oDate1.getTime() > oDate2.getTime() ? -1 : 1;
          });
          r.lastVersion = sorted[0];
          r.sortedVersions = sorted;
        }
        return r;
      });
      response.records = records;
    }
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().chartList;
  },
  onFinish: (record, dispatch: Redux.Dispatch) => {
    let isNotNeedPoll = true;
    if (record.data.recordCount) {
      isNotNeedPoll =
        record.data.records.filter(
          item => item.status && item.status['phase'] && item.status['phase'] === 'Terminating'
        ).length === 0;
    }
    if (isNotNeedPoll) {
      dispatch(fetchChartActions.clearPolling());
    }
  }
});

/**
 * 删除模板
 */
const removeChartWorkflow = generateWorkflowActionCreator<Chart, ChartFilter>({
  actionType: ActionTypes.RemoveChart,
  workflowStateLocator: (state: RootState) => state.chartRemoveWorkflow,
  operationExecutor: WebAPI.deleteChart,
  after: {
    [OperationTrigger.Done]: (dispatch: Redux.Dispatch, getState: GetState) => {
      //在列表页删除的动作，因此直接重新拉取一次数据
      const params = getState().chartRemoveWorkflow.params;
      dispatch(listActions.poll(params));
      /** 结束工作流 */
      dispatch(listActions.removeChartWorkflow.reset());
    }
  }
});

const restActions = {
  removeChartWorkflow,

  /** 轮询操作 */
  poll: (filter: ChartFilter) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch(
        listActions.polling({
          filter: filter,
          delayTime: 3000
        })
      );
    };
  }
};

export const listActions = extend({}, fetchChartActions, restActions);
