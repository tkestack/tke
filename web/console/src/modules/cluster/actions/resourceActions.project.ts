import { extend, RecordSet, ReduxAction } from '@tencent/qcloud-lib';
import { FetchOptions, generateFetcherActionCreator } from '@tencent/qcloud-redux-fetcher';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';

import { resourceConfig } from '../../../../config';
import { ResourceInfo } from '../../common/models';
import { includes } from '../../common/utils';
import { IsResourceShowLoadingIcon } from '../components/resource/resourceTableOperation/ResourceTablePanel';
import * as ActionType from '../constants/ActionType';
import { PollEventName, ResourceNeedJudgeLoading } from '../constants/Config';
import { Resource, ResourceFilter, RootState } from '../models';
import { router } from '../router';
import * as WebAPI from '../WebAPI';
import { namespaceActions } from './namespaceActions';
import { resourceDetailActions } from './resourceDetailActions';
import { resourceDetailEventActions } from './resourceDetailEventActions';
import { serviceEditActions } from './serviceEditActions';
import { workloadEditActions } from './workloadEditActions';

type GetState = () => RootState;
const fetchOptions: FetchOptions = {
  noCache: false
};

/** fetch resource list */
const fetchResourceActions = generateFetcherActionCreator({
  actionType: ActionType.FetchResourceList,
  fetcher: async (getState: GetState, fetchOptions, dispatch) => {
    let { subRoot, projectNamespaceList, clusterVersion } = getState(),
      { resourceInfo, resourceOption, resourceName } = subRoot,
      { resourceQuery } = resourceOption;

    let isClearData = fetchOptions && fetchOptions.noCache ? true : false;

    let response: any;
    if (resourceName === 'lbcf') {
      response = _reduceGameGateResource(clusterVersion, resourceQuery, resourceInfo, isClearData);
    } else if (resourceName === 'np') {
      let list = [];
      projectNamespaceList.data.records.forEach(item => {
        list.push({
          metadata: { name: item.spec.namespace, creationTimestamp: item.metadata.creationTimestamp },
          spec: {
            clusterId: item.spec.clusterName,
            clusterVersion: item.spec.clusterVersion,
            clusterDisplayName: item.spec.clusterDisplayName,
            hard: item.spec.hard
          },
          status: {
            phase: item.status.phase,
            used: item.status.used || {}
          }
        });
      });
      const result: RecordSet<Resource> = {
        recordCount: list.length,
        records: list
      };
      return result;
    } else {
      response = await WebAPI.fetchResourceList(resourceQuery, resourceInfo, isClearData);
    }

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
        (resourceName === 'deployment' ||
          resourceName === 'statefulset' ||
          resourceName === 'daemonset' ||
          resourceName === 'tapp')
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

async function _reduceGameGateResource(clusterVersion, resourceQuery, resourceInfo, isClearData) {
  let gameBGresourceInfo = resourceConfig(clusterVersion).lbcf_bg;
  let gameBRresourceInfo = resourceConfig(clusterVersion).lbcf_br;
  let gameBGList = await WebAPI.fetchResourceList(resourceQuery, gameBGresourceInfo, isClearData),
    gameLBList = await WebAPI.fetchResourceList(resourceQuery, resourceInfo, isClearData),
    gameBRList = await WebAPI.fetchResourceList(resourceQuery, gameBRresourceInfo, isClearData);
  gameLBList.records.forEach((item, index) => {
    let backGroups = [];
    gameBGList.records.forEach(backgroup => {
      if (backgroup.spec.lbName === item.metadata.name) {
        let backendRecords = gameBRList.records.filter(
          records => records.metadata.labels['lbcf.tkestack.io/backend-group'] === backgroup.metadata.name
        );
        try {
          backGroups.push({
            name: backgroup.metadata.name,
            labels: backgroup.spec.pods.byLabel.selector,
            port: backgroup.spec.pods.port,
            status: backgroup.status,
            backendRecords: backendRecords.map(record => {
              return {
                name: record.metadata.name,
                backendAddr: record.status && record.status.backendAddr ? record.status.backendAddr : '-',
                conditions: record.status && record.status.conditions ? record.status.conditions : []
              };
            })
          });
        } catch (e) {}
      }
    });
    gameLBList.records[index].spec.backGroups = backGroups;
  });
  return gameLBList;
}

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

      //detailResourceInfo初始化
      if (resourceInfo.requestType && resourceInfo.requestType.useDetailInfo) {
        dispatch(resourceActions.initDetailResourceName(name));
      }

      dispatch({
        type: ActionType.InitResourceInfo,
        payload: resourceInfo
      });
    };
  },

  //只有当需要使用detailresourceInfo,每个页面的配置不一致的时候需要触发这个方法选择正确的detailresourceName
  changeDetailTab: (tab: string) => {
    return async (dispatch, getState: GetState) => {
      let { subRoot, clusterVersion, route } = getState(),
        {
          detailResourceOption: { detailResourceName, detailResourceList },
          resourceInfo
        } = subRoot;
      let list = resourceInfo.requestType.detailInfoList[tab];
      if (list) {
        let finder = list.find(item => item.value === detailResourceName);
        if (!finder) {
          dispatch(resourceActions.initDetailResourceName(list[0].value));
        }
      }
    };
  },

  /** 修改当前资源的名称 */
  initDetailResourceName: (resourceName: string) => {
    return async (dispatch, getState: GetState) => {
      let {
        subRoot: { mode }
      } = getState();
      dispatch({
        type: ActionType.InitDetailResourceName,
        payload: resourceName
      });
      // 初始化 detailresourceInfo的信息
      dispatch(resourceActions.initDetailResourceInfo(resourceName));

      mode === 'detail' && dispatch(resourceActions.initDetailResourceList(resourceName));
    };
  },

  //addon里面有些crd是由两个资源组成，所以在detail页面有时需要在不更新当前resourceInfo,切换resourceInfo
  initDetailResourceInfo: (rsName?: string) => {
    return async (dispatch, getState: GetState) => {
      let { subRoot, clusterVersion } = getState();

      let resourceInfo: ResourceInfo = resourceConfig(clusterVersion)[rsName] || {};

      dispatch({
        type: ActionType.InitDetailResourceInfo,
        payload: resourceInfo
      });
    };
  },

  initDetailResourceList: (rsName?: string) => {
    return async (dispatch, getState: GetState) => {
      let {
        route,
        subRoot: {
          resourceName,
          resourceOption: { resourceSelection }
        }
      } = getState();
      let list = [];
      if (rsName === resourceName) {
        let defaultResourceIns =
          route.queries['resourceIns'] || (resourceSelection[0] && resourceSelection[0].metadata.name);
        list.push({ value: defaultResourceIns, text: defaultResourceIns });
      } else if (rsName === 'lbcf_bg') {
        resourceSelection[0] &&
          resourceSelection[0].spec.backGroups &&
          resourceSelection[0].spec.backGroups.forEach(item => {
            list.push({ value: item.name, text: item.name });
          });
      } else if (rsName === 'lbcf_br') {
        resourceSelection[0] &&
          resourceSelection[0].spec.backGroups &&
          resourceSelection[0].spec.backGroups.forEach(item => {
            for (let i = 0; i < item.backendRecords.length; ++i) {
              list.push({ value: item.backendRecords[i].name, text: item.backendRecords[i].name });
            }
          });
      }
      dispatch({
        type: ActionType.InitDetailResourceList,
        payload: list
      });
      dispatch(resourceActions.selectDetailResouceIns(list.length ? list[0].value : ''));
    };
  },

  selectDetailResouceIns: (rsIns: string) => {
    return async (dispatch, getState: GetState) => {
      let { route, subRoot } = getState(),
        { resourceDetailState, resourceInfo } = subRoot,
        { event } = resourceDetailState;
      let { tab } = router.resolve(route);
      dispatch({
        type: ActionType.SelectDetailResourceSelection,
        payload: rsIns
      });
      //如果存在这类资源则重新拉取数据
      if (rsIns) {
        if (tab === 'yaml') {
          dispatch(resourceDetailActions.fetchResourceYaml.fetch());
        } else if (tab === 'event') {
          dispatch(resourceDetailEventActions.poll());
        }
      }
    };
  },

  selectDetailDeleteResouceIns: (rsIns: string) => {
    return async (dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.SelectDetailDeleteResourceSelection,
        payload: rsIns
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
