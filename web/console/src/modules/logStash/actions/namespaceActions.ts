import { extend, ReduxAction } from '@tencent/ff-redux';
import { generateFetcherActionCreator } from '@tencent/qcloud-redux-fetcher';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';

import { resourceConfig } from '../../../../config';
import { cloneDeep } from '../../common/';
import { NamespaceFilter } from '../../common/models';
import * as ActionType from '../constants/ActionType';
import { ContainerLogs, RootState } from '../models';
import { router } from '../router';
import * as WebAPI from '../WebAPI';
import { resourceActions } from './resourceActions';

type GetState = () => RootState;

/** 拉取namesapce列表 */
const fetchNamespaceListActions = generateFetcherActionCreator({
  actionType: ActionType.FetchNamespaceList,
  fetcher: async (getState: GetState, fetchOptions, dispatch: Redux.Dispatch) => {
    let { clusterVersion, namespaceQuery } = getState();
    // 获取当前的资源配置
    let namesapceInfo = resourceConfig(clusterVersion)['ns'];
    let isClearData = fetchOptions && fetchOptions.noCache;
    let response = await WebAPI.fetchNamespaceList(namespaceQuery, namesapceInfo, isClearData);
    return response;
  },
  finish: (dispatch: Redux.Dispatch, getState: GetState) => {
    dispatch(namespaceActions.autoSelectNamespaceForCreate());
  }
});

/** namespace列表的查询 */
const queryNamesapceListActions = generateQueryActionCreator<NamespaceFilter>({
  actionType: ActionType.QueryNamespaceList,
  bindFetcher: fetchNamespaceListActions
});

const restActions = {
  selectNamespace: (namespace: string): ReduxAction<string> => {
    return {
      type: ActionType.NamespaceSelection,
      payload: namespace
    };
  },
  /**
   * 创建日志采集页面的时候
   *    帮用户选择namespace
   */
  autoSelectNamespaceForCreate: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { route } = getState();
      let urlParams = router.resolve(route);

      //只有在创建状态下才需要帮用户选择默认选项
      if (urlParams['mode'] === 'create') {
        if (getState().namespaceList.data.recordCount) {
          dispatch({ type: ActionType.SelectContainerFileNamespace, payload: 'default' });
          let containerLogsArr: ContainerLogs[] = cloneDeep(getState().logStashEdit.containerLogs);
          containerLogsArr[0].namespaceSelection = 'default';
          dispatch({
            type: ActionType.UpdateContainerLogs,
            payload: containerLogsArr
          });
          //进行资源的拉取
          dispatch(
            resourceActions.applyFilter({
              clusterId: route.queries['clusterId'],
              namespace: 'default',
              workloadType: 'deployment',
              regionId: +route.queries['rid']
            })
          );
        } else {
          dispatch({ type: ActionType.SelectContainerFileNamespace, payload: '' });
          let containerLogsArr: ContainerLogs[] = cloneDeep(getState().logStashEdit.containerLogs);
          containerLogsArr[0].namespaceSelection = '';
          dispatch({
            type: ActionType.UpdateContainerLogs,
            payload: containerLogsArr
          });
          dispatch(
            resourceActions.fetch({
              noCache: true
            })
          );
        }
      }
    };
  }
};

export const namespaceActions = extend(fetchNamespaceListActions, queryNamesapceListActions, restActions);
