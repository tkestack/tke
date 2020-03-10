import * as JsYAML from 'js-yaml';

import { ReduxAction } from '@tencent/ff-redux';
import { generateFetcherActionCreator } from '@tencent/qcloud-redux-fetcher';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import * as ActionType from '../constants/ActionType';
import { Helm, RootState } from '../models';
import * as WebAPI from '../WebAPI';

type GetState = () => RootState;
const tips = seajs.require('tips');

const fetchHistory = generateFetcherActionCreator({
  actionType: ActionType.FetchHistory,
  fetcher: async (getState: GetState, fetchOptions, dispatch) => {
    let {
      detailState: { historyQuery }
    } = getState();
    let histories = await WebAPI.fetchHistory(
      { helmName: historyQuery.filter.helmName },
      historyQuery.filter.regionId,
      historyQuery.filter.clusterId
    );
    return histories;
  },
  finish: (dispatch, getState: GetState) => {
    // let {
    //   listState: { helmList },
    //   route
    // } = getState();
    // if (helmList.data.recordCount) {
    //   dispatch(restActions.select(helmList.data.records[0]));
    // }
  }
});
const queryHistory = generateQueryActionCreator({
  actionType: ActionType.QueryHistory,
  bindFetcher: fetchHistory
});

export const detailActions = {
  clear: (): ReduxAction<void> => {
    return { type: ActionType.ClearDetail };
  },
  select: (helm: Helm) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.SelectHelm,
        payload: helm
      });
    };
  },

  timer4refresh: null,
  fetchHelm: (helmName: string) => {
    return async (dispatch, getState: GetState) => {
      let {
        route,
        detailState: { isRefresh }
      } = getState();
      if (isRefresh) {
        clearTimeout(detailActions.timer4refresh);
        detailActions.timer4refresh = setTimeout(() => {
          dispatch(detailActions.fetchHelm(helmName));
        }, 3000);
      } else {
        clearTimeout(detailActions.timer4refresh);
        detailActions.timer4refresh = null;
      }
      let helm = await WebAPI.fetchHelm({ helmName }, +route.queries['rid'], route.queries['clusterId']);
      let yamlResponse = await WebAPI.fetchHelmResourceList(
        { helmName },
        +route.queries['rid'],
        route.queries['clusterId']
      );

      try {
        let yamls = yamlResponse.release.manifest.split('---').slice(1);
        helm.resources = yamls.map(item => {
          let json = JsYAML.safeLoad(item);

          return {
            name: json.metadata.name,
            kind: json.kind,
            yaml: item
          };
        });
      } catch (e) {
        helm.resources = [];
        tips.error(t('资源列表读取失败'), 2000);
      }

      helm.chart_metadata = yamlResponse.release.chart.metadata;

      helm.valueYaml = yamlResponse.release.chart.values.raw;

      if (yamlResponse.release.config.raw) {
        try {
          helm.configYaml = JsYAML.safeDump(JSON.parse(yamlResponse.release.config.raw));
        } catch (e) {}
      }

      dispatch({
        type: ActionType.FetchHelm,
        payload: helm
      });
    };
  },
  fetchHistory,
  queryHistory,
  rollback: (helmName: string, version: number) => {
    return async (dispatch, getState: GetState) => {
      let { route } = getState();

      await WebAPI.rollbackVersion({ helmName, version }, +route.queries['rid'], route.queries['clusterId']);

      dispatch(fetchHistory.fetch());
    };
  },

  setRefresh: (isRefresh: boolean) => {
    return async (dispatch, getState: GetState) => {
      let { route } = getState();
      dispatch({
        type: ActionType.IsRefresh,
        payload: isRefresh
      });

      if (isRefresh) {
        dispatch(detailActions.fetchHelm(getState().detailState.helm.name));
      } else {
        clearTimeout(detailActions.timer4refresh);
        detailActions.timer4refresh = null;
      }
    };
  }
};
