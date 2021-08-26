/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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

import { extend, FetchOptions, generateFetcherActionCreator } from '@tencent/ff-redux';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';

import { resourceConfig } from '../../../../config';
import * as ActionType from '../constants/ActionType';
import { PollEventName } from '../constants/Config';
import { PodLogFilter, ResourceFilter, RootState } from '../models';
import * as WebAPI from '../WebAPI';

type GetState = () => RootState;
const fetchOptions: FetchOptions = {
  noCache: false
};

/** ===== start workload的相关操作  ============ */
const fetchWorkloadActions = generateFetcherActionCreator({
  actionType: ActionType.L_FetchWorkloadList,
  fetcher: async (getState: GetState, fetchOptions, dispatch) => {
    let { subRoot, clusterVersion } = getState(),
      { workloadQuery, workloadType } = subRoot.resourceLogOption;

    let workloadResourceInfo = resourceConfig(clusterVersion)[workloadType];
    let isClearData = fetchOptions && fetchOptions.noCache ? true : false;
    let response = await WebAPI.fetchResourceList(workloadQuery, {
      resourceInfo: workloadResourceInfo,
      isClearData
    });
    return response;
  },
  finish: async (dispatch, getState: GetState) => {
    let { workloadList } = getState().subRoot.resourceLogOption;

    // 如果有workload，则默认选择第一个
    workloadList.data.recordCount &&
      dispatch(workloadActions.selectWorkload(workloadList.data.records[0].metadata.name));
  }
});

const queryWorkloadActions = generateQueryActionCreator({
  actionType: ActionType.L_QueryWorkloadList,
  bindFetcher: fetchWorkloadActions
});

const workloadRestActions = {
  /** 选择某个具体的workload */
  selectWorkload: (name: string) => {
    return async (dispatch, getState: GetState) => {
      let { route, subRoot } = getState(),
        { namespaceSelection, workloadType } = subRoot.resourceLogOption;

      dispatch({
        type: ActionType.L_WorkloadSelection,
        payload: name
      });

      // 如果name不为空，则需要拉取该workload下的pod的列表
      name !== '' &&
        namespaceSelection !== '' &&
        workloadType !== '' &&
        dispatch(
          podActions.applyFilter({
            namespace: namespaceSelection,
            regionId: +route.queries['rid'],
            clusterId: route.queries['clusterId'],
            specificName: name
          })
        );
    };
  },

  /** 选择workloadtype */
  selectWorkloadType: (type: string) => {
    return async (dispatch, getState: GetState) => {
      let { route, subRoot } = getState(),
        { namespaceSelection } = subRoot.resourceLogOption;

      dispatch({
        type: ActionType.L_WorkloadType,
        payload: type
      });

      // 切换类型的时候，需要重新进行列表的拉取
      type !== '' &&
        namespaceSelection !== '' &&
        dispatch(
          workloadActions.applyFilter({
            namespace: namespaceSelection,
            clusterId: route.queries['clusterId'],
            regionId: route.queries['rid']
          })
        );
    };
  }
};

const workloadActions = extend(fetchWorkloadActions, queryWorkloadActions, workloadRestActions);
/** ===== end workload的相关操作  ============ */

/** ===== start pod的相关操作  ============ */
const fetchPodActions = generateFetcherActionCreator({
  actionType: ActionType.L_FetchPodList,
  fetcher: async (getState: GetState, fetchOptions, dispatch) => {
    let { subRoot, clusterVersion } = getState(),
      { workloadType, podQuery } = subRoot.resourceLogOption;

    let workloadResourceInfo = resourceConfig(clusterVersion)[workloadType];
    let isClearData = fetchOptions && fetchOptions.noCache ? true : false;
    let response = await WebAPI.fetchExtraResourceList(podQuery, workloadResourceInfo, isClearData, 'pods');
    return response;
  },
  finish: async (dispatch, getState: GetState) => {
    let { podList } = getState().subRoot.resourceLogOption;

    // 如果podlist有，则默认选择第一个，否则选空
    dispatch(podActions.selectPod(podList.data.recordCount ? podList.data.records[0].metadata.name : ''));
  }
});

const queryPodActions = generateQueryActionCreator<ResourceFilter>({
  actionType: ActionType.L_QueryPodList,
  bindFetcher: fetchPodActions
});

const podRestActions = {
  /** 选择pod */
  selectPod: (podName: string) => {
    return async (dispatch, getState: GetState) => {
      let { podList } = getState().subRoot.resourceLogOption;

      dispatch({
        type: ActionType.L_PodSelection,
        payload: podName
      });

      // 选择了pod之后，需要默认选中第一个container
      let finder = podList.data.records.find(p => p.metadata.name === podName);
      finder && dispatch(podActions.selectContainer(finder.spec.containers ? finder.spec.containers[0].name : ''));
    };
  },

  /** 选择container */
  selectContainer: (containerName: string) => {
    return async (dispatch, getState: GetState) => {
      let { subRoot } = getState(),
        { podSelection, tailLines } = subRoot.resourceLogOption;

      dispatch({
        type: ActionType.L_ContainerSelection,
        payload: containerName
      });

      // 进行数据的拉取
      dispatch(resourceLogActions.fetchLogData(podSelection, containerName, tailLines));
    };
  }
};

const podActions = extend(fetchPodActions, queryPodActions, podRestActions);
/** ===== end pod的相关操作  ============ */

/** ===== start log的相关操作  ============ */
const fetchLogActions = generateFetcherActionCreator({
  actionType: ActionType.L_FetchLogList,
  fetcher: async (getState: GetState, fetchOptions, dispatch) => {
    let { subRoot, clusterVersion } = getState(),
      { logQuery, workloadType, tailLines, containerSelection } = subRoot.resourceLogOption;

    let isClearData = fetchOptions && fetchOptions.noCache ? true : false;
    let podResourceInfo = resourceConfig(clusterVersion)['pods'];

    let k8sQueryObj = {
      container: containerSelection ? containerSelection : undefined,
      timestamps: true,
      tailLines: tailLines === 'all' ? undefined : tailLines
    };

    k8sQueryObj = JSON.parse(JSON.stringify(k8sQueryObj));
    let response = await WebAPI.fetchResourceLogList(logQuery, podResourceInfo, isClearData, k8sQueryObj);
    return response;
  }
});

const queryLogActions = generateQueryActionCreator<PodLogFilter>({
  actionType: ActionType.L_QueryLogList,
  bindFetcher: fetchLogActions
});

const logRestActions = {
  poll: (queryObj: any) => {
    return async (dispatch, getState: GetState) => {
      // 触发日志的轮询
      dispatch(logActions.applyFilter(queryObj));
      // 每次轮询之前先清空之前的轮询
      dispatch(logActions.clearPollLog());

      window[PollEventName['resourceLog']] = setInterval(() => {
        dispatch(logActions.poll(queryObj));
      }, 8000);
    };
  },

  clearPollLog: () => {
    return async (dispatch, getState: GetState) => {
      clearInterval(window[PollEventName['resourceLog']]);
    };
  },

  /** 切换自动刷新 */
  toggleAutoRenew: () => {
    return async (dispatch, getState: GetState) => {
      let { isAutoRenew } = getState().subRoot.resourceLogOption;
      dispatch({
        type: ActionType.L_IsAutoRenew,
        payload: !isAutoRenew
      });
    };
  },

  /** 选择展示数据的条数 */
  selectTailLine: (tailLine: string) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.L_TailLines,
        payload: tailLine
      });
    };
  }
};

const logActions = extend(fetchLogActions, queryLogActions, logRestActions);
/** ===== end log的相关操作  ============ */

export const resourceLogActions = {
  /** workload的相关操作 */
  workload: workloadActions,

  /** pod的相关操作 */
  pod: podActions,

  /** log的相关操作 */
  log: logActions,

  /** 选择命名空间 */
  selectNamespace: (namespace: string) => {
    return async (dispatch, getState: GetState) => {
      let { route, subRoot } = getState(),
        { workloadType } = subRoot.resourceLogOption;

      dispatch({
        type: ActionType.L_NamespaceSelection,
        payload: namespace
      });

      // 切换命名空间时候，需要进行workload列表的拉取
      namespace !== '' &&
        workloadType !== '' &&
        dispatch(workloadActions.applyFilter({ namespace, clusterId: route.queries['clusterId'] }));
    };
  },

  /** 获取日志的数据 */
  fetchLogData: (podName: string, containerName: string, tailLines: string) => {
    return async (dispatch, getState: GetState) => {
      let { route, subRoot } = getState(),
        { isAutoRenew, namespaceSelection } = subRoot.resourceLogOption;

      // 进行log的拉取
      let filterObj: PodLogFilter = {
        regionId: +route.queries['rid'],
        namespace: namespaceSelection,
        clusterId: route.queries['clusterId'],
        specificName: podName,
        container: containerName,
        tailLines
      };

      // 只有实例名称和选择对应的容器才会进行日志的拉取
      if (containerName !== '' && podName !== '') {
        !isAutoRenew && dispatch(logActions.toggleAutoRenew());
        dispatch(logActions.poll(filterObj));
      } else {
        dispatch(resourceLogActions.closeAutoRenewAndClearLog());
      }
    };
  },

  /** 关闭自动刷新，并且清空原来的日志数据 */
  closeAutoRenewAndClearLog: () => {
    return async (dispatch, getState: GetState) => {
      let { isAutoRenew } = getState().subRoot.resourceLogOption;

      isAutoRenew && dispatch(logActions.toggleAutoRenew());
      dispatch(logActions.fetch({ noCache: true }));
      dispatch(logActions.clearPollLog());
    };
  }
};
