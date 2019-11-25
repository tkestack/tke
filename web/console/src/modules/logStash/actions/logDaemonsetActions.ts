import { extend } from '@tencent/qcloud-lib';
import { generateFetcherActionCreator } from '@tencent/qcloud-redux-fetcher';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';
import * as ActionType from '../constants/ActionType';
import { RootState } from '../models';
import * as WebAPI from '../WebAPI';
import { Validation } from '../../common/models';
import { resourceConfig } from '../../../../config';
import { LogDaemonSetFliter } from '../models/LogDaemonset';
import { CommonAPI, includes } from '../../common';
import { logActions } from './logActions';
import { canFetchLogList } from '../constants/Config';
import { router } from '../router';
type GetState = () => RootState;

/** 获取Log采集器的列表的Action */
const fetchLogDaemonsetActions = generateFetcherActionCreator({
  actionType: ActionType.FetchLogDaemonset,
  fetcher: async (getState: GetState, fetchOptions, dispatch) => {
    let { clusterVersion, logDaemonsetQuery } = getState();
    let resourceInfo = resourceConfig(clusterVersion)['addon_logcollector'];
    let isClearData = fetchOptions && fetchOptions.noCache ? true : false;
    let response = await CommonAPI.fetchResourceList({
      query: logDaemonsetQuery,
      resourceInfo,
      isClearData
    });
    return response;
  },
  finish: (dispatch: Redux.Dispatch, getState: GetState) => {
    dispatch(logDaemonsetActions.isOpenLogStash());
    dispatch(logDaemonsetActions.isDaemonsetNormal());
    let { route, namespaceSelection, isOpenLogStash, isDaemonsetNormal } = getState();
    let { clusterId } = route.queries;
    let urlParams = router.resolve(route);
    if (!urlParams['mode']) {
      if (isOpenLogStash && includes(canFetchLogList, isDaemonsetNormal.phase)) {
        dispatch(
          logActions.applyFilter({
            clusterId,
            namespace: namespaceSelection
          })
        );
      } else {
        dispatch(
          logActions.fetch({
            noCache: true
          })
        );
      }
    }
  }
});

/** 查询log采集器列表的Action */
const QueryLogDaemonset = generateQueryActionCreator<LogDaemonSetFliter>({
  actionType: ActionType.QueryLogDaemonset,
  bindFetcher: fetchLogDaemonsetActions
});

export const restActions = {
  /**
   * 判断当前集群是否已经开通日志采集的功能
   */
  isOpenLogStash: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      const { logDaemonset } = getState();
      dispatch({
        type: ActionType.IsOpenLogStash,
        payload: logDaemonset.error ? false : true
      });
    };
  },

  /** 判断当前的daemonset是否正常 */
  isDaemonsetNormal: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      const { logDaemonset } = getState();
      let phase = '',
        reason = '';

      if (logDaemonset.error) {
        phase = '404';
        reason = 'not found';
      } else {
        phase = logDaemonset.data.records[0].status.phase;
        reason = logDaemonset.data.records[0].status.reason;
      }

      dispatch({
        type: ActionType.IsDaemonsetNormal,
        payload: {
          phase,
          reason
        }
      });
    };
  }
};

//需要写一个函数获取全部的资源resource
export const logDaemonsetActions = extend({}, fetchLogDaemonsetActions, QueryLogDaemonset, restActions);
