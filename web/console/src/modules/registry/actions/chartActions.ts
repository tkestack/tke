import { initValidator } from './../../common/models/Validation';
import { extend, deepClone, uuid } from '@tencent/qcloud-lib';
import { generateWorkflowActionCreator, OperationTrigger, isSuccessWorkflow } from '@tencent/qcloud-redux-workflow';
import { RootState, Chart, ChartFilter } from '../models';
import * as ActionType from '../constants/ActionType';
import { InitChart } from '../constants/Config';
import * as WebAPI from '../WebAPI';
import { router } from '../router';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { createListAction } from '@tencent/redux-list';
import { ChartCreation } from '../models/Chart';

type GetState = () => RootState;

const FFModelChartActions = createListAction<Chart, ChartFilter>({
  actionName: 'chart',
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchChartList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().chart;
  }
});

const restActions = {
  /** 创建 Chart */
  createChart: generateWorkflowActionCreator<ChartCreation, void>({
    actionType: ActionType.CreateChart,
    workflowStateLocator: (state: RootState) => state.createChart,
    operationExecutor: WebAPI.createChart,
    after: {
      [OperationTrigger.Done]: (dispatch, getState) => {
        let { createChart, route } = getState();
        if (isSuccessWorkflow(createChart)) {
          dispatch(restActions.createChart.reset());
          dispatch(restActions.clearEdition());
          dispatch(chartActions.fetch());
          let urlParams = router.resolve(route);
          router.navigate(Object.assign({}, urlParams, { sub: 'chart', mode: 'list' }), {});
        }
      }
    }
  }),

  /** 删除 Chart */
  deleteChart: generateWorkflowActionCreator<Chart, void>({
    actionType: ActionType.DeleteChart,
    workflowStateLocator: (state: RootState) => state.deleteChart,
    operationExecutor: WebAPI.deleteChart,
    after: {
      [OperationTrigger.Done]: (dispatch, getState) => {
        let { deleteChart, route } = getState();
        if (isSuccessWorkflow(deleteChart)) {
          dispatch(restActions.deleteChart.reset());
          dispatch(chartActions.fetch());
        }
      }
    }
  }),

  /** --begin编辑action */
  inputChartDesc: (value: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.UpdateChartCreation,
        payload: Object.assign({}, getState().chartCreation, { displayName: value })
      });
    };
  },

  inputChartName: (value: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.UpdateChartCreation,
        payload: Object.assign({}, getState().chartCreation, { name: value })
      });
      dispatch(chartActions.validateChartName(value));
    };
  },

  selectChartVisibility: (value: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.UpdateChartCreation,
        payload: Object.assign({}, getState().chartCreation, { visibility: value })
      });
    };
  },

  validateChartName(value: string) {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let result = chartActions._validateChartName(value);
      dispatch({
        type: ActionType.UpdateChartCreation,
        payload: Object.assign({}, getState().chartCreation, { v_name: result })
      });
    };
  },

  _validateChartName(name: string) {
    let reg = /^[a-z]([-a-z0-9]*[a-z0-9])?$/,
      status = 0,
      message = '';

    if (!name) {
      status = 2;
      message = t('Chart Group 不能为空');
    } else if (name.length > 63) {
      status = 2;
      message = t('Chart Group 不能超过63个字符');
    } else if (!reg.test(name)) {
      status = 2;
      message = t('Chart Group 格式不正确');
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  clearEdition: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.UpdateChartCreation,
        payload: InitChart
      });
    };
  }
  /** --end编辑action */
};

export const chartActions = extend({}, FFModelChartActions, restActions);
