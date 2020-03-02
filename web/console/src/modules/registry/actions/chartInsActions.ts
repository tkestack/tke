import { extend } from '@tencent/qcloud-lib';
import { RootState, ChartIns, ChartInsFilter } from '../models';
import * as ActionType from '../constants/ActionType';
import * as WebAPI from '../WebAPI';
import { router } from '../router';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { createListAction } from '@tencent/redux-list';

type GetState = () => RootState;

const FFModelChartInsActions = createListAction<ChartIns, ChartInsFilter>({
  actionName: 'chartIns',
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchChartInsList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().chartIns;
  }
});

export const chartInsActions = extend({}, FFModelChartInsActions);
