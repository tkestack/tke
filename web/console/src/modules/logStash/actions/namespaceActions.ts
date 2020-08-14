import { extend, ReduxAction } from '@tencent/ff-redux';
import { generateFetcherActionCreator } from '@tencent/qcloud-redux-fetcher';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';

import { resourceConfig } from '../../../../config';
import { cloneDeep, CommonAPI } from '../../common/';
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
    // 获取当前的资源配置，兼容业务侧和平台侧
    let namesapceInfo = resourceConfig(clusterVersion)[namespaceQuery && namespaceQuery.filter && namespaceQuery.filter.projectName ? 'namespaces' : 'ns'];
    let isClearData = fetchOptions && fetchOptions.noCache;
    let response = await WebAPI.fetchNamespaceList(namespaceQuery, namesapceInfo, isClearData);
    if (namespaceQuery && namespaceQuery.filter && namespaceQuery.filter.projectName) {
      // 如果是在业务侧, 给cluster注入logAgent信息。因为在业务侧操作的是业务和命名空间，只有通过命名空间的信息转换出集群信息来
      let agents = await CommonAPI.fetchLogagents();
      let clusterHasLogAgent = {};
      for (let agent of agents.records) {
        clusterHasLogAgent[agent.spec.clusterName] = { name: agent.metadata.name, status: agent.status.phase };
      }
      for (let ns of response.records) {
        let logagent = clusterHasLogAgent[ns.cluster.metadata.name];
        if (logagent) {
          let { name, status } = logagent;
          ns.cluster.spec.logAgentName = name;
          ns.cluster.spec.logAgentStatus = status;
        }
      }
    }
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
      let { route, namespaceList } = getState();
      let urlParams = router.resolve(route);

      //只有在创建状态下才需要帮用户选择默认选项
      if (urlParams['mode'] === 'create') {
        if (namespaceList.data.recordCount) {
          let namespace = namespaceList.data.records[0].namespace;
          dispatch({ type: ActionType.SelectContainerFileNamespace, payload: namespace });
          let containerLogsArr: ContainerLogs[] = cloneDeep(getState().logStashEdit.containerLogs);
          containerLogsArr[0].namespaceSelection = namespace;
          dispatch({
            type: ActionType.UpdateContainerLogs,
            payload: containerLogsArr
          });
          //进行资源的拉取
          dispatch(
            resourceActions.applyFilter({
              clusterId: route.queries['clusterId'],
              namespace,
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
