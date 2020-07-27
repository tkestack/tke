import {
  extend,
  generateWorkflowActionCreator,
  OperationResult,
  OperationTrigger,
  createFFListActions
} from '@tencent/ff-redux';
import { RootState, ChartGroup, ChartGroupFilter } from '../../models';
import * as ActionTypes from '../../constants/ActionTypes';
import * as WebAPI from '../../WebAPI';
type GetState = () => RootState;

/**
 * 列表操作
 */
const fetchChartGroupActions = createFFListActions<ChartGroup, ChartGroupFilter>({
  actionName: ActionTypes.ChartGroupList,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchChartGroupList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().chartGroupList;
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
      dispatch(fetchChartGroupActions.clearPolling());
    }
  }
});

const restActions = {};

export const listActions = extend({}, fetchChartGroupActions, restActions);
