import { extend, ReduxAction, RecordSet } from '@tencent/qcloud-lib';
import { generateFetcherActionCreator, FetchOptions } from '@tencent/qcloud-redux-fetcher';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';
import { RootState, Resource, ResourceFilter } from '../models';
import * as ActionType from '../constants/ActionType';
import * as WebAPI from '../WebAPI';
import { ResourceInfo } from '../../common/models';
import { resourceConfig } from '../../../../config';
import { namespaceActions } from './namespaceActions';
import { serviceEditActions } from './serviceEditActions';
import { workloadEditActions } from './workloadEditActions';
import { PollEventName, ResourceNeedJudgeLoading } from '../constants/Config';
import { includes } from '../../common/utils';
import { IsResourceShowLoadingIcon } from '../components/resource/resourceTableOperation/ResourceTablePanel';
import { router } from '../router';

type GetState = () => RootState;
const fetchOptions: FetchOptions = {
  noCache: false
};

/** fetch resource list */
const fetchResourceActions = generateFetcherActionCreator({
  actionType: ActionType.FetchResourceList,
  fetcher: async (getState: GetState, fetchOptions, dispatch) => {
    let { subRoot, projectNamespaceList } = getState(),
      { resourceInfo, resourceOption, resourceName } = subRoot,
      { resourceQuery } = resourceOption;

    let isClearData = fetchOptions && fetchOptions.noCache ? true : false;
    if (resourceName === 'np') {
      let list = [];
      projectNamespaceList.data.records.forEach(item => {
        list.push({
          metadata: { name: item.spec.namespace, creationTimestamp: item.metadata.creationTimestamp },
          spec: {},
          status: { phase: item.status.phase }
        });
      });

      const result: RecordSet<Resource> = {
        recordCount: list.length,
        records: list
      };
      return result;
    }

    let response = await WebAPI.fetchResourceList(resourceQuery, resourceInfo, isClearData);
    return response;
  },
  finish: async (dispatch, getState: GetState) => {
    let { subRoot, route } = getState(),
      { tab } = router.resolve(route),
      { resourceOption, mode, resourceName } = subRoot,
      { resourceList } = resourceOption;

    if (resourceList.data.recordCount) {
      let defaultResourceIns = route.queries['resourceIns'];
      let finder = resourceList.data.records.find(item => item.metadata.name === defaultResourceIns);
      dispatch(resourceActions.selectResource([finder ? finder : resourceList.data.records[0]]));

      /** ============== start 更新的时候，进行一些页面的初始化 =============  */
      if (mode === 'update' && resourceName === 'svc') {
        dispatch(serviceEditActions.initServiceEditForUpdate(finder));
      } else if (
        mode === 'update' &&
        tab === 'modifyRegistry' &&
        (resourceName === 'deployment' || resourceName === 'statefulset' || resourceName === 'daemonset')
      ) {
        dispatch(workloadEditActions.initWorkloadEditForUpdateRegistry(finder));
      } else if (mode === 'update' && tab === 'modifyPod') {
        dispatch(workloadEditActions.updateContainerNum(finder.spec.replicas));
      }
      /** ============== end 更新的时候，进行一些页面的初始化 =============  */

      /** ============== start 列表页，需要进行资源的轮询 ================= */
      if (mode === 'list' && includes(ResourceNeedJudgeLoading, resourceName)) {
        if (resourceList.data.records.filter(item => IsResourceShowLoadingIcon(resourceName, item)).length === 0) {
          dispatch(resourceActions.clearPollEvent());
        }
      } else {
        dispatch(resourceActions.clearPollEvent());
      }
      /** ============== end 列表页，需要进行资源的轮询 ================= */
    } else {
      dispatch(resourceActions.clearPollEvent());
    }
  }
});

/** query resource list action */
const queryResourceActions = generateQueryActionCreator({
  actionType: ActionType.QueryResourceList,
  bindFetcher: fetchResourceActions
});

const restActions = {
  /** 在列表上选择具体的资源，如在deploymentList当中选择某个 deployment */
  selectResource: (resource: Resource[]) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.SelectResource,
        payload: resource
      });
    };
  },

  /** 在列表上选择多个具体的资源，如在deploymentList当中选择某几个具体的deployment */
  selectMultipleResource: (resource: Resource[]) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.SelectMultipleResource,
        payload: resource
      });
    };
  },

  /** 选择删除的资源 */
  selectDeleteResource: (resource: Resource[]) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.SelectDeleteResource,
        payload: resource
      });
    };
  },

  /** 修改当前资源的名称 */
  initResourceName: (resourceName: string) => {
    return async dispatch => {
      dispatch({
        type: ActionType.InitResourceName,
        payload: resourceName
      });

      // 初始化 resourceInfo的信息
      dispatch(resourceActions.initResourceInfo(resourceName));
    };
  },

  /** 初始化 resource */
  initResourceInfo: (rsName?: string) => {
    return async (dispatch, getState: GetState) => {
      let { subRoot, clusterVersion } = getState(),
        { resourceName } = subRoot;

      let resourceInfo: ResourceInfo,
        name = rsName ? rsName : resourceName;
      resourceInfo = resourceConfig(clusterVersion)[name] || {};

      dispatch({
        type: ActionType.InitResourceInfo,
        payload: resourceInfo
      });
    };
  },

  /** 变更当前的模式 */
  selectMode: (mode: string) => {
    return async dispatch => {
      dispatch({
        type: ActionType.SelectMode,
        payload: mode
      });
    };
  },

  /** 判断当前是否需要拉取资源的namespace列表 */
  toggleIsNeedFetchNamespace: (isNeedFetch: boolean): ReduxAction<boolean> => {
    return {
      type: ActionType.IsNeedFetchNamespace,
      payload: isNeedFetch
    };
  },

  /** 路由变化，不同的资源切换的时候，需要进行数据的初始化 */
  initResourceInfoAndFetchData: (
    isNeedFetchNamespace: boolean = true,
    resourceName: string,
    isNeedClear: boolean = true
  ) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { clusterId, rid } = getState().route.queries;
      // 判断是否需要展示ns
      dispatch(resourceActions.toggleIsNeedFetchNamespace(isNeedFetchNamespace));
      // 初始化当前的资源的名称
      dispatch(resourceActions.initResourceName(resourceName));
      // 进行ns的拉取
      dispatch(namespaceActions.applyFilter({ clusterId, regionId: +rid }));
      // 是否需要清空resourceList
      isNeedClear && dispatch(resourceActions.fetch({ noCache: true }));
    };
  },

  /** 轮询拉取条件 */
  poll: (queryObj: ResourceFilter) => {
    return async (dispatch, getState: GetState) => {
      // 每次轮询之前先清空之前的轮询
      dispatch(resourceActions.clearPollEvent());
      // 触发列表的查询
      dispatch(resourceActions.applyFilter(queryObj));

      window[PollEventName['resourceList']] = setInterval(() => {
        dispatch(resourceActions.poll(queryObj));
      }, 8000);
    };
  },

  /** 清空轮询条件 */
  clearPollEvent: () => {
    return async (dispatch, getState: GetState) => {
      clearInterval(window[PollEventName['resourceList']]);
    };
  },

  /** 清除subRoot的信息 */
  clearSubRoot: (): ReduxAction<any> => {
    return {
      type: ActionType.ClearSubRoot
    };
  }
};

export const resourceActions = extend(fetchResourceActions, queryResourceActions, restActions);
