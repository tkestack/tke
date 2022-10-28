/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */
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
    const { clusterVersion, namespaceQuery } = getState();
    // 获取当前的资源配置，兼容业务侧和平台侧
    const namesapceInfo =
      resourceConfig(clusterVersion)[
        namespaceQuery && namespaceQuery.filter && namespaceQuery.filter.projectName ? 'namespaces' : 'ns'
      ];
    const isClearData = fetchOptions && fetchOptions.noCache;
    const response = await WebAPI.fetchNamespaceList(namespaceQuery, namesapceInfo, isClearData);
    if (namespaceQuery && namespaceQuery.filter && namespaceQuery.filter.projectName) {
      // 如果是在业务侧, 给cluster注入logAgent信息。因为在业务侧操作的是业务和命名空间，只有通过命名空间的信息转换出集群信息来
      const agents = await CommonAPI.fetchLogagents();
      const clusterHasLogAgent = {};
      for (const agent of agents.records) {
        clusterHasLogAgent[agent.spec.clusterName] = { name: agent.metadata.name, status: agent.status.phase };
      }
      for (const ns of response.records) {
        const logagent = clusterHasLogAgent[ns.cluster.metadata.name];
        if (logagent) {
          const { name, status } = logagent;
          ns.cluster.spec.logAgentName = name;
          ns.cluster.spec.logAgentStatus = status;
        }
      }
    }
    return response;
  },
  finish: (dispatch: Redux.Dispatch, getState: GetState) => {
    const { namespaceList } = getState();

    dispatch(namespaceActions.selectNamespace(namespaceList?.data?.records?.[0]?.namespace ?? ''));
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
      const { route, namespaceList } = getState();
      const urlParams = router.resolve(route);

      //只有在创建状态下才需要帮用户选择默认选项
      if (urlParams['mode'] === 'create') {
        if (namespaceList.data.recordCount) {
          const namespace = namespaceList.data.records[0].namespace;
          dispatch({ type: ActionType.SelectContainerFileNamespace, payload: namespace });
          const containerLogsArr: ContainerLogs[] = cloneDeep(getState().logStashEdit.containerLogs);
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
          const containerLogsArr: ContainerLogs[] = cloneDeep(getState().logStashEdit.containerLogs);
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
