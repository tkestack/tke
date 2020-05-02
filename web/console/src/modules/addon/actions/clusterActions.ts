import { createFFListActions, extend } from '@tencent/ff-redux';

import { resourceConfig } from '../../../../config';
import { Resource, ResourceFilter, ResourceInfo } from '../../common';
import { CommonAPI } from '../../common/webapi';
import * as ActionType from '../constants/ActionType';
import { AddonStatusEnum, FFReduxActionName } from '../constants/Config';
import { RootState } from '../models';
import { router } from '../router';
import { addonActions } from './addonActions';

type GetState = () => RootState;

/** ========================== start 已存在的add的相关操作 ========================== */
const ListAddonActions = createFFListActions<Resource, ResourceFilter>({
  actionName: FFReduxActionName.OPENADDON,
  fetcher: async (query, getState: GetState) => {
    let clusterInfo: ResourceInfo = resourceConfig()['cluster'];
    let response = await CommonAPI.fetchExtraResourceList({
      query,
      resourceInfo: clusterInfo,
      extraResource: 'addons'
    });

    // 对结果进行排序，保证每次的结果一样，后台是通过promise.all 并行请求的，返回结果顺序不确定
    response.records = response.records.sort((prev, next) => (prev.spec.type < next.spec.type ? 1 : -1));
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().openAddon;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    let { route } = getState(),
      urlParams = router.resolve(route);
    let { resourceIns } = route.queries;

    if (record.data.recordCount) {
      // 只有url里面有resourceIns才需要进行选择
      if (resourceIns) {
        let finder = resourceIns ? record.data.records.find(addon => addon.metadata.name === resourceIns) : null;
        let defaultAddon = finder ? finder : record.data.records[0];
        dispatch(openAddonActions.select(defaultAddon));
      }

      let mode = urlParams['mode'];
      // 只有在列表页面才需要进行轮询
      if (!mode) {
        // 如果addon的列表存在非 running 状态的，则需要继续进行轮询，否则清空轮询
        if (
          record.data.records.filter(
            item =>
              item.status.phase.toLowerCase() !== AddonStatusEnum.Running &&
              item.status.phase.toLowerCase() !== AddonStatusEnum.Failed
          ).length === 0
        ) {
          dispatch(openAddonActions.clearPolling());
        }
      } else {
        dispatch(openAddonActions.clearPolling());
      }
    } else {
      dispatch(openAddonActions.clearPolling());
    }
  }
});

const addonRestActions = {
  poll: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { route } = getState();
      let { rid, clusterId } = route.queries;
      dispatch(
        openAddonActions.polling({
          filter: { regionId: +rid, specificName: clusterId },
          delayTime: 8000
        })
      );
    };
  }
};

const openAddonActions = extend({}, ListAddonActions, addonRestActions);
/** ========================== end 已存在的add的相关操作 ========================== */

/** 集群列表的Actions */
const ListClusterActions = createFFListActions<Resource, ResourceFilter>({
  actionName: FFReduxActionName.CLUSTER,
  fetcher: async (query, getState: GetState) => {
    let clusterInfo: ResourceInfo = resourceConfig()['cluster'];
    let response = await CommonAPI.fetchResourceList({ query, resourceInfo: clusterInfo });
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().cluster;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    let { route } = getState();

    if (record.data.recordCount) {
      let routeClusterId = route.queries['clusterId'];
      let finder = routeClusterId ? record.data.records.find(c => c.metadata.name === routeClusterId) : undefined;
      let defaultCluster = finder ? finder : record.data.records[0];
      dispatch(clusterActions.selectCluster(defaultCluster));
    }
  }
});

const restActions = {
  addon: openAddonActions,

  selectCluster: (cluster: Resource) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      if (cluster) {
        let { route } = getState(),
          urlParams = router.resolve(route);
        dispatch(ListClusterActions.select(cluster));

        let clusterId: string = cluster.metadata.name;

        // 进行路由的跳转
        router.navigate(urlParams, Object.assign({}, route.queries, { clusterId }));

        // 初始化集群的版本
        dispatch(clusterActions.initClusterVersion(cluster.status.version));

        let { rid } = route.queries;

        // 进行所有addon列表的查询
        dispatch(addonActions.applyFilter({ regionId: +rid }));

        // 进行当前集群的addon列表的拉取，需要进行轮询
        dispatch(openAddonActions.poll());
      }
    };
  },

  /** 初始化集群的版本 */
  initClusterVersion: (clusterVersion: string) => {
    return async (dispatch: Redux.Dispatch) => {
      let k8sVersion = clusterVersion.split('.').slice(0, 2).join('.');
      dispatch({
        type: ActionType.ClusterVersion,
        payload: k8sVersion
      });
    };
  }
};

export const clusterActions = extend({}, ListClusterActions, restActions);
