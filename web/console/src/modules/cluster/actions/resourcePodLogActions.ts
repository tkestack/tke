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
import { extend, FetchOptions, generateFetcherActionCreator, ReduxAction } from '@tencent/ff-redux';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';

import { LogAgent } from 'src/modules/common/models';
import { resourceConfig } from '../../../../config/resourceConfig';
import * as ActionType from '../constants/ActionType';
import { PollEventName } from '../constants/Config';
import { PodLogFilter, RootState, LogHierarchyQuery, LogContentQuery, DownloadLogQuery } from '../models';
import { router } from '../router';
import * as WebAPI from '../WebAPI';

type GetState = () => RootState;
const fetchOptions: FetchOptions = {
  noCache: false
};

/** 获取pod的日志 action */
const fetchLogActions = generateFetcherActionCreator({
  actionType: ActionType.FetchLogList,
  fetcher: async (getState: GetState, fetchOptions, dispatch) => {
    let { subRoot, clusterVersion } = getState(),
      { resourceDetailState } = subRoot,
      { logQuery } = resourceDetailState,
      { container, tailLines } = logQuery.filter;

    let isClearData = fetchOptions && fetchOptions.noCache ? true : false;
    let podResourceInfo = resourceConfig(clusterVersion)['pods'];

    let k8sQueryObj = {
      container: container ? container : undefined,
      timestamps: true,
      tailLines: tailLines === 'all' ? undefined : tailLines
    };

    k8sQueryObj = JSON.parse(JSON.stringify(k8sQueryObj));
    let response = await WebAPI.fetchResourceLogList(logQuery, podResourceInfo, isClearData, k8sQueryObj);
    return response;
  }
});

/** 查询pod的日志列表action */
const queryLogActions = generateQueryActionCreator<PodLogFilter>({
  actionType: ActionType.QueryLogList,
  bindFetcher: fetchLogActions
});

/** 剩余的log的操作 */
const restActions = {
  /** 选择pod */
  selectPod: (podName: string) => {
    return async (dispatch, getState: GetState) => {
      let { podList } = getState().subRoot.resourceDetailState;
      dispatch({
        type: ActionType.PodName,
        payload: podName
      });

      let finder = podList.data.records.find(item => item.metadata.name === podName),
        containerList = finder ? finder.spec.containers : [],
        containerName = '';
      if (finder && containerList.length > 0) {
        containerName = containerList[0].name;
      }
      dispatch(resourcePodLogActions.selectContainer(containerName));
    };
  },

  /** 选择container */
  selectContainer: (containerName: string) => {
    return async (dispatch, getState: GetState) => {
      let { route, subRoot } = getState(),
        urlParams = router.resolve(route),
        { logAgent, logOption } = subRoot.resourceDetailState,
        { podName, logFile, tailLines } = logOption;

      dispatch({
        type: ActionType.ContainerName,
        payload: containerName
      });

      // 进行数据的拉取，如果tab 不在 log 页面的话，不需要进行拉取
      if (urlParams['tab'] === 'log') {
        dispatch(resourcePodLogActions.handleFetchData(podName, containerName, tailLines));
        // 拉取日志目录结构
        if (logAgent && logAgent['metadata'] && logAgent['metadata']['name']) {
          let agentName = logAgent['metadata']['name'];
          const query: LogHierarchyQuery = {
            agentName,
            namespace: route.queries['np'],
            clusterId: route.queries['clusterId'],
            pod: podName,
            container: containerName,
          };
          dispatch(resourcePodLogActions.getLogHierarchy(query));
        }
      }
    };
  },

  /** 选择日志文件 */
  selectLogFile: (logFile: string) => {
    return async (dispatch, getState: GetState) => {
      let { route, subRoot } = getState(),
        urlParams = router.resolve(route),
        { logAgent } = subRoot.resourceDetailState,
        { podName, containerName, tailLines } = subRoot.resourceDetailState.logOption;

      dispatch({
        type: ActionType.LogFile,
        payload: logFile
      });

      if (logFile === 'stdout') {
        urlParams['tab'] === 'log' && dispatch(resourcePodLogActions.handleFetchData(podName, containerName, tailLines));
      } else {
        let agentName = '';
        if (logAgent && logAgent['metadata'] && logAgent['metadata']['name']) {
          agentName = logAgent['metadata']['name'];
        }
        const query: LogContentQuery = {
          agentName,
          namespace: route.queries['np'],
          clusterId: route.queries['clusterId'],
          pod: podName,
          container: containerName,
          filepath: logFile,
          start: 0,
          length: Number(tailLines)
        };
        urlParams['tab'] === 'log' && dispatch(resourcePodLogActions.getLogContent(query));
      }
    };
  },

  /** 选择展示的日志的条数 */
  selectTailLine: (tailLines: string) => {
    return async (dispatch, getState: GetState) => {
      let { podName, containerName } = getState().subRoot.resourceDetailState.logOption;

      dispatch({
        type: ActionType.TailLines,
        payload: tailLines
      });

      // 进行数据的拉取
      dispatch(resourcePodLogActions.handleFetchData(podName, containerName, tailLines));
    };
  },

  /** 进行日志的拉取 */
  handleFetchData: (podName: string, containerName: string, tailLines: string) => {
    return async (dispatch, getState: GetState) => {
      let { route, subRoot } = getState(),
        { isAutoRenew } = subRoot.resourceDetailState.logOption;

      // 进行log的拉取
      let filterObj: PodLogFilter = {
        regionId: +route.queries['rid'],
        namespace: route.queries['np'],
        clusterId: route.queries['clusterId'],
        specificName: podName,
        container: containerName,
        tailLines
      };

      // 只有有实例名称 和 选择对应的容器才会进行日志的拉取
      if (containerName !== '' && podName !== '') {
        !isAutoRenew && dispatch(resourcePodLogActions.toggleAutoRenew());
        dispatch(resourcePodLogActions.poll(filterObj));
      } else {
        dispatch(resourcePodLogActions.fetch({ noCache: true }));
        dispatch(resourcePodLogActions.clearPollLog());
        // 如果不符合拉取日志的条件，则关闭自动刷新
        isAutoRenew && dispatch(resourcePodLogActions.toggleAutoRenew());
      }
    };
  },

  setLogAgent: (logAgent: LogAgent) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.PodLogAgent,
        payload: logAgent
      });
    };
  },

  /** 拉取采集日志结构 **/
  getLogHierarchy: (query: LogHierarchyQuery) => {
    return async (dispatch, getState: GetState) => {
      let response = await WebAPI.fetchResourceLogHierarchy(query);
      dispatch({
        type: ActionType.PodLogHierarchy,
        payload: response
      });
    };
  },

  /** 获取日志文件内容 **/
  getLogContent: (query: LogContentQuery) => {
    return async (dispatch, getState: GetState) => {
      let response = await WebAPI.fetchResourceLogContent(query);
      dispatch({
        type: ActionType.PodLogContent,
        payload: response
      });
    };
  },

  /** 下载日志文件 **/
  downloadLogFile: (query: DownloadLogQuery) => {
    return async (dispatch, getState: GetState) => {
      let response = await WebAPI.downloadLogFile(query);
    };
  },

  /** 切换自动刷新 */
  toggleAutoRenew: () => {
    return async (dispatch, getState: GetState) => {
      let { isAutoRenew } = getState().subRoot.resourceDetailState.logOption;
      dispatch({
        type: ActionType.IsAutoRenewPodLog,
        payload: !isAutoRenew
      });
    };
  },

  /** 轮训日志的拉取 */
  poll: (queryObj: any) => {
    return async (dispatch, getState: GetState) => {
      let { route, subRoot } = getState(),
        { logAgent, logOption } = subRoot.resourceDetailState,
        { podName, containerName, logFile, tailLines, isAutoRenew } = logOption,
      urlParams = router.resolve(route);

      // 触发日志的轮询
      if (logFile === 'stdout') {
        dispatch(resourcePodLogActions.applyFilter(queryObj));
      } else {
        let agentName = '';
        if (logAgent && logAgent['metadata'] && logAgent['metadata']['name']) {
          agentName = logAgent['metadata']['name'];
        }
        let { namespace, clusterId, specificName: pod, container, tailLines } = queryObj;
        let logContentQuery: LogContentQuery = {
          agentName,
          namespace,
          clusterId,
          pod,
          container,
          start: 0,
          length: Number(tailLines),
          filepath: logFile,
        };
        dispatch(resourcePodLogActions.getLogContent(logContentQuery));
      }
      // 每次轮询之前先清空之前的轮询
      dispatch(resourcePodLogActions.clearPollLog());

      window[PollEventName['resourcePodLog']] = setInterval(() => {
        urlParams['tab'] === 'log' && dispatch(resourcePodLogActions.poll(queryObj));
      }, 8000);
    };
  },

  clearPollLog: () => {
    return async (dispatch, getState: GetState) => {
      clearInterval(window[PollEventName['resourcePodLog']]);
    };
  }
};

export const resourcePodLogActions = extend(fetchLogActions, queryLogActions, restActions);
