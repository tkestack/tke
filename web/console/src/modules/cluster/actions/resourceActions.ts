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
import { createFFListActions, extend, FetchOptions, ReduxAction } from '@tencent/ff-redux';
import { resourceConfig } from '../../../../config';
import { ResourceInfo } from '../../common/models';
import { includes } from '../../common/utils';
import { IsResourceShowLoadingIcon } from '../components/resource/resourceTableOperation/ResourceTablePanel';
import * as ActionType from '../constants/ActionType';
import { FFReduxActionName, ResourceNeedJudgeLoading } from '../constants/Config';
import { Resource, ResourceFilter, RootState } from '../models';
import { router } from '../router';
import * as WebAPI from '../WebAPI';
import { lbcfEditActions } from './lbcfEditActions';
import { namespaceActions } from './namespaceActions';
import { resourceDetailActions } from './resourceDetailActions';
import { resourceDetailEventActions } from './resourceDetailEventActions';
import { serviceEditActions } from './serviceEditActions';
import { workloadEditActions } from './workloadEditActions';

type GetState = () => RootState;
const fetchOptions: FetchOptions = {
  noCache: false
};

const listResourceActions = createFFListActions<Resource, ResourceFilter>({
  actionName: FFReduxActionName.Resource_Workload,
  fetcher: async (query, getState: GetState, fetchOptions) => {
    const { subRoot, route, clusterVersion } = getState(),
      { resourceInfo, resourceOption, resourceName } = subRoot,
      { ffResourceList } = resourceOption;

    const isClearData = fetchOptions && fetchOptions.noCache ? true : false;

    let response: any;
    // response = await WebAPI.fetchResourceList(resourceQuery, resourceInfo, isClearData);
    if (resourceName === 'lbcf') {
      response = _reduceGameGateResource(clusterVersion, ffResourceList.query, resourceInfo, isClearData);
    } else {
      response = await WebAPI.fetchResourceList(ffResourceList.query, {
        resourceInfo,
        isClearData,
        isContinue: true
      });
    }

    console.log('ffResourceList.query:  ', ffResourceList.query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().subRoot.resourceOption.ffResourceList;
  },
  onFinish: (record, dispatch, getState: GetState) => {
    let { subRoot, route } = getState(),
      { tab } = router.resolve(route),
      { resourceOption, mode, resourceName } = subRoot,
      { ffResourceList } = resourceOption;

    if (ffResourceList.list.data.recordCount) {
      const defaultResourceIns = route.queries['resourceIns'];
      const finder = ffResourceList.list.data.records.find(
        item => item.metadata && item.metadata.name === defaultResourceIns
      );
      dispatch(resourceActions.select(finder ? finder : ffResourceList.list.data.records[0]));

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
      } else if (mode === 'update' && tab === 'updateBG') {
        dispatch(lbcfEditActions.initGameBGEdition(finder.spec.backGroups));
      }
      /** ============== end 更新的时候，进行一些页面的初始化 =============  */

      /** ============== start 列表页，需要进行资源的轮询 ================= */
      if (mode === 'list' && includes(ResourceNeedJudgeLoading, resourceName)) {
        if (record.data.records.filter(item => IsResourceShowLoadingIcon(resourceName, item)).length === 0) {
          dispatch(resourceActions.clearPolling());
        } else {
          dispatch(resourceActions.startPolling({ delayTime: 8000 }));
        }
      } else {
        dispatch(resourceActions.clearPolling());
      }
      /** ============== end 列表页，需要进行资源的轮询 ================= */
    } else {
      dispatch(resourceActions.clearPolling());
    }
  }
});

async function _reduceGameGateResource(clusterVersion, resourceQuery, resourceInfo, isClearData) {
  const gameBGresourceInfo = resourceConfig(clusterVersion).lbcf_bg;
  const gameBRresourceInfo = resourceConfig(clusterVersion).lbcf_br;
  const gameBGList = await WebAPI.fetchResourceList(resourceQuery, {
      resourceInfo: gameBGresourceInfo,
      isClearData
    }),
    gameLBList = await WebAPI.fetchResourceList(resourceQuery, {
      resourceInfo,
      isClearData,
      isContinue: true
    }),
    gameBRList = await WebAPI.fetchResourceList(resourceQuery, {
      resourceInfo: gameBRresourceInfo,
      isClearData
    });
  gameLBList.records.forEach((item, index) => {
    const backGroups = [];
    gameBGList.records.forEach(backgroup => {
      if (backgroup.spec.lbName === item.metadata.name) {
        const backendRecords = gameBRList.records.filter(
          records => records.metadata.labels['lbcf.tkestack.io/backend-group'] === backgroup.metadata.name
        );
        try {
          const backGroup = {
            name: backgroup.metadata.name,
            status: backgroup.status,
            backendRecords: backendRecords.map(record => {
              return {
                name: record.metadata.name,
                backendAddr: record.status && record.status.backendAddr ? record.status.backendAddr : '-',
                conditions: record.status && record.status.conditions ? record.status.conditions : []
              };
            })
          };
          if (backgroup.spec.pods) {
            backGroup['pods'] = {
              labels: backgroup.spec.pods.byLabel.selector,
              port: backgroup.spec.pods.port,
              byName: backgroup.spec.pods.byName
            };
          } else if (backgroup.spec.service) {
            backGroup['service'] = {
              name: backgroup.spec.service.name,
              port: backgroup.spec.service.port,
              nodeSelector: backgroup.spec.service.nodeSelector
            };
          } else {
            backGroup['static'] = backgroup.spec.static;
          }
          backGroups.push(backGroup);
        } catch (e) {}
      }
    });
    gameLBList.records[index].spec.backGroups = backGroups;
  });
  return gameLBList;
}

const restActions = {
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
      const { subRoot, clusterVersion } = getState(),
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
      const list = resourceInfo.requestType.detailInfoList[tab];
      if (list) {
        const finder = list.find(item => item.value === detailResourceName);
        if (!finder) {
          dispatch(resourceActions.initDetailResourceName(list[0].value));
        }
      }
    };
  },

  /** 修改当前资源的名称 */
  initDetailResourceName: (resourceName: string, defaultIns?: string) => {
    return async (dispatch, getState: GetState) => {
      const {
        subRoot: { mode }
      } = getState();
      dispatch({
        type: ActionType.InitDetailResourceName,
        payload: resourceName
      });
      // 初始化 detailresourceInfo的信息
      dispatch(resourceActions.initDetailResourceInfo(resourceName));

      mode === 'detail' && dispatch(resourceActions.initDetailResourceList(resourceName, defaultIns));
    };
  },

  //addon里面有些crd是由两个资源组成，所以在detail页面有时需要在不更新当前resourceInfo,切换resourceInfo
  initDetailResourceInfo: (rsName?: string) => {
    return async (dispatch, getState: GetState) => {
      const { subRoot, clusterVersion } = getState();

      const resourceInfo: ResourceInfo = resourceConfig(clusterVersion)[rsName] || {};

      dispatch({
        type: ActionType.InitDetailResourceInfo,
        payload: resourceInfo
      });
    };
  },

  initDetailResourceList: (rsName?: string, defaultIns?: string) => {
    return async (dispatch, getState: GetState) => {
      const {
        route,
        subRoot: {
          resourceName,
          resourceOption: { ffResourceList }
        }
      } = getState();
      const list = [];
      if (rsName === resourceName) {
        const defaultResourceIns =
          route.queries['resourceIns'] || (ffResourceList.selection && ffResourceList.selection.metadata.name);
        list.push({ value: defaultResourceIns, text: defaultResourceIns });
      } else if (rsName === 'lbcf_bg') {
        ffResourceList.selection &&
          ffResourceList.selection.spec.backGroups &&
          ffResourceList.selection.spec.backGroups.forEach(item => {
            list.push({ value: item.name, text: item.name });
          });
      } else if (rsName === 'lbcf_br') {
        ffResourceList.selection &&
          ffResourceList.selection.spec.backGroups &&
          ffResourceList.selection.spec.backGroups.forEach(item => {
            for (let i = 0; i < item.backendRecords.length; ++i) {
              list.push({ value: item.backendRecords[i].name, text: item.backendRecords[i].name });
            }
          });
      }
      dispatch({
        type: ActionType.InitDetailResourceList,
        payload: list
      });
      dispatch(resourceActions.selectDetailResouceIns(defaultIns ? defaultIns : list.length ? list[0].value : ''));
    };
  },

  selectDetailResouceIns: (rsIns: string) => {
    return async (dispatch, getState: GetState) => {
      let { route, subRoot } = getState(),
        { resourceDetailState, resourceInfo } = subRoot,
        { event } = resourceDetailState;
      const { tab } = router.resolve(route);
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
  initResourceInfoAndFetchData: (isNeedFetchNamespace = true, resourceName: string, isNeedClear = true) => {
    return async (dispatch, getState: GetState) => {
      const { clusterId, rid } = getState().route.queries;
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
  poll: () => {
    return async (dispatch, getState: GetState) => {
      const { route } = getState();
      const { np, clusterId, rid, meshId } = route.queries;

      const filterObj: ResourceFilter = {
        namespace: np,
        clusterId,
        regionId: +rid,
        meshId
      };

      dispatch(
        resourceActions.polling({
          filter: filterObj,
          delayTime: 8000
        })
      );
    };
  },

  /** 清除subRoot的信息 */
  clearSubRoot: (): ReduxAction<any> => {
    return {
      type: ActionType.ClearSubRoot
    };
  }
};

export const resourceActions = extend(listResourceActions, restActions);
