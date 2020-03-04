import { extend } from '@tencent/qcloud-lib';
import { FetchOptions, generateFetcherActionCreator } from '@tencent/ff-redux';
import { RootState, ResourceFilter, PodLogFilter, EventFilter } from '../models';
import * as ActionType from '../constants/ActionType';
import { resourceConfig } from '../../../../config';
import * as WebAPI from '../WebAPI';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';
import { PollEventName } from '../constants/Config';
import { ResourceInfo } from '../../common/models';

type GetState = () => RootState;
const fetchOptions: FetchOptions = {
  noCache: false
};

/** ======= start workload的相关操作 ============== */
const fetchWorkloadActions = generateFetcherActionCreator({
  actionType: ActionType.E_FetchWorkloadList,
  fetcher: async (getState: GetState, fetchOptions, dispatch) => {
    let { clusterVersion, subRoot } = getState(),
      { workloadQuery, workloadType } = subRoot.resourceEventOption;

    let workloadResourceInfo = resourceConfig(clusterVersion)[workloadType];

    let isClearData = fetchOptions && fetchOptions.noCache ? true : false;
    let response = await WebAPI.fetchResourceList(workloadQuery, workloadResourceInfo, isClearData);
    return response;
  }
});

const queryWorkloadActions = generateQueryActionCreator({
  actionType: ActionType.E_QueryWorkloadList,
  bindFetcher: fetchWorkloadActions
});

const workloadRestActions = {
  /** 选择某个具体的workload */
  selectWorkload: (name: string) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.E_WorkloadSelection,
        payload: name
      });
    };
  },

  /** 选择workloadtype */
  selectWorkloadType: (type: string) => {
    return async (dispatch, getState: GetState) => {
      let { route, subRoot } = getState(),
        { namespaceSelection } = subRoot.resourceEventOption;

      dispatch({
        type: ActionType.E_WorkloadType,
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
/** ======= end workload的相关操作 ============== */

const fetchEventActions = generateFetcherActionCreator({
  actionType: ActionType.E_FetchEventList,
  fetcher: async (getState: GetState, fetchOptions, dispatch) => {
    let { subRoot, clusterVersion } = getState(),
      { eventQuery } = subRoot.resourceEventOption,
      { kind, name } = eventQuery.filter;

    let eventResourceInfo = resourceConfig(clusterVersion)['event'];

    let isClearData = fetchOptions && fetchOptions.noCache ? true : false;

    let k8sQueryObj = {
      fieldSelector: {
        'involvedObject.kind': kind ? kind : undefined,
        'involvedObject.name': name ? name : undefined
      },
      limit: 20
    };

    k8sQueryObj = JSON.parse(JSON.stringify(k8sQueryObj));
    let response = await WebAPI.fetchResourceList(eventQuery, eventResourceInfo, isClearData, k8sQueryObj);
    return response;
  }
});

const queryEventActions = generateQueryActionCreator({
  actionType: ActionType.E_QueryEventList,
  bindFetcher: fetchEventActions
});

const restActons = {
  /** workload的相关操作 */
  workload: workloadActions,

  /** 选择命名空间 */
  selectNamespace: (namespace: string) => {
    return async (dispatch, getState: GetState) => {
      let { route, subRoot } = getState(),
        { workloadType } = subRoot.resourceEventOption;

      dispatch({
        type: ActionType.E_NamespaceSelection,
        payload: namespace
      });

      // 切换命名空间时候，需要进行workload列表的拉取
      namespace !== '' &&
        workloadType !== '' &&
        dispatch(workloadActions.applyFilter({ namespace, clusterId: route.queries['clusterId'] }));
    };
  },

  /** 轮询拉取事件 */
  poll: (queryObj: any) => {
    return async (dispatch, getState: GetState) => {
      // 每次轮询之前先清空之前的轮询
      dispatch(resourceEventActions.clearPollEvent());
      // 触发事件的轮询
      dispatch(resourceEventActions.applyFilter(queryObj));

      window[PollEventName['resourceEvent']] = setInterval(() => {
        dispatch(resourceEventActions.poll(queryObj));
      }, 10000);
    };
  },

  /** 清空轮询条件 */
  clearPollEvent: () => {
    return async (dispatch, getState: GetState) => {
      clearInterval(window[PollEventName['resourceEvent']]);
    };
  },

  /** 切换自动刷新 */
  toggleAutoRenew: () => {
    return async (dispatch, getState: GetState) => {
      let { isAutoRenew } = getState().subRoot.resourceEventOption;
      dispatch({
        type: ActionType.E_IsAutoRenew,
        payload: !isAutoRenew
      });
    };
  },

  /** 获取事件的数据 */
  fetchEventData: (kind: string, namespace: string, name: string, isFromAutoSwicth: boolean = false) => {
    return async (dispatch, getState: GetState) => {
      let { route, subRoot, clusterVersion } = getState(),
        { isAutoRenew } = subRoot.resourceEventOption;

      let resourceInfo: ResourceInfo = resourceConfig(clusterVersion)[kind];

      // 进行log的拉取
      let filterObj: EventFilter = {
        regionId: +route.queries['rid'],
        namespace,
        clusterId: route.queries['clusterId'],
        kind: resourceInfo ? resourceInfo.headTitle : undefined,
        name
      };

      // 如果是开启了自动刷新，则用Poll进行拉取
      if (isAutoRenew || isFromAutoSwicth) {
        dispatch(resourceEventActions.poll(filterObj));
      } else {
        dispatch(resourceEventActions.applyFilter(filterObj));
      }
    };
  }
};

export const resourceEventActions = extend(fetchEventActions, queryEventActions, restActons);
