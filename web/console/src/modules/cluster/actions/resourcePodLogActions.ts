import { ReduxAction, extend } from '@tencent/qcloud-lib';
import { generateFetcherActionCreator, FetchOptions } from '@tencent/qcloud-redux-fetcher';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';
import { RootState, PodLogFilter } from '../models';
import * as ActionType from '../constants/ActionType';
import * as WebAPI from '../WebAPI';
import { PollEventName } from '../constants/Config';
import { router } from '../router';
import { resourceConfig } from '../../../../config/resourceConfig';

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
        { podName, tailLines } = subRoot.resourceDetailState.logOption;

      dispatch({
        type: ActionType.ContainerName,
        payload: containerName
      });

      // 进行数据的拉取，如果tab 不在 log 页面的话，不需要进行拉取
      urlParams['tab'] === 'log' && dispatch(resourcePodLogActions.handleFetchData(podName, containerName, tailLines));
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
      let { route } = getState(),
        urlParams = router.resolve(route);

      // 触发日志的轮询
      dispatch(resourcePodLogActions.applyFilter(queryObj));
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
