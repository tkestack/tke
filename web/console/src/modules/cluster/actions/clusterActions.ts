import {
    createFFListActions, extend, FetchOptions, generateFetcherActionCreator, uuid
} from '@tencent/ff-redux';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';

import { resourceConfig } from '../../../../config';
import { Cluster, ClusterFilter, ResourceInfo } from '../../common/';
import {
    TellIsNeedFetchNS, TellIsNotNeedFetchResource
} from '../components/resource/ResourceSidebarPanel';
import * as ActionType from '../constants/ActionType';
import { FFReduxActionName } from '../constants/Config';
import { initAllcationRatioEdition } from '../constants/initState';
import { DialogNameEnum, RootState } from '../models';
import { AddonStatus } from '../models/Addon';
import { router } from '../router';
import * as WebAPI from '../WebAPI';
import { resourceActions } from './resourceActions';

type GetState = () => RootState;
const fetchOptions: FetchOptions = {
  noCache: false
};

/** ===================== 通过k8s接口拉取详情，获取k8s的版本详情 ==============================
 * 这里主要有如下的情况
 * 1. 如果在集群详情页内，需要判断clusterVersion
 * 2. 拉取集群列表是分页的，如果刚开始进来，clusterId不在当前列表页，并且clusterVersion 刚好不对，则会出问题
 *
 * pre:
 * 1. 拉取地域时，如果不是在集群详情页，不会调用此方法去获取集群的详情
 * 2. 在详情页直接刷新，会判断当前的clusterInfoList是否存在 和地域是否为空，没有则进行拉取并且初始化 k8sVersion
 */
const fetchClusterInfoActions = generateFetcherActionCreator({
  actionType: ActionType.FetchClusterInfo,
  fetcher: async (getState: GetState, fetchOptions, dispatch) => {
    let { clusterInfoQuery } = getState();

    let resourceInfo = resourceConfig()['cluster'];
    let isClearData = fetchOptions && fetchOptions.noCache ? true : false;
    let response = await WebAPI.fetchSpecificResourceList(clusterInfoQuery, resourceInfo, isClearData, true);
    return response;
  },
  finish: async (dispatch, getState: GetState) => {
    let { clusterInfoList, route } = getState(),
      urlParams = router.resolve(route);

    /**
     * 如果存在该集群，并且当前是在集群详情页里面，则初始化集群版本
     * 如果是从集群列表页进行 clusterInfo的拉取的话，则不需要进行下面的步骤，因为点击集群id，已经进行了resource数据的拉取，并且更快
     */
    if (clusterInfoList.data.recordCount) {
      let version = clusterInfoList.data.records[0] && clusterInfoList.data.records[0].status.version;
      version && dispatch(clusterActions.initClusterVersion(version));

      // 初始化当前资源的名称，是deployment | cronjob 还是其他
      let { resourceName: resource } = urlParams;

      // 如果不需要进行resource的拉取，则不拉取
      if (!TellIsNotNeedFetchResource(resource)) {
        let isNeedFetchNamespace = TellIsNeedFetchNS(resource);
        // 进行相关资源的拉取
        dispatch(resourceActions.initResourceInfoAndFetchData(isNeedFetchNamespace, resource, false));
      }
    }
  }
});

const queryClusterInfoActions = generateQueryActionCreator({
  actionType: ActionType.QueryClusterInfo,
  bindFetcher: fetchClusterInfoActions
});

const clusterInfoActions = extend(fetchClusterInfoActions, queryClusterInfoActions);
/** =================================================================================== */

/** 集群列表的Actions */
const FFModelClusterActions = createFFListActions<Cluster, ClusterFilter>({
  actionName: FFReduxActionName.CLUSTER,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchClusterList(query, query.filter.regionId);
    let ps = await WebAPI.fetchPrometheuses();
    let clusterHasPs = {};
    for (let p of ps.records) {
      clusterHasPs[p.spec.clusterName] = true;
    }
    for (let record of response.records) {
      record.spec.hasPrometheus = clusterHasPs[record.metadata.name];
    }
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().cluster;
  },
  onFinish: (record, dispatch, getState: GetState) => {
    let { route, dialogState, cluster } = getState(),
      urlParams = router.resolve(route);

    let { sub } = urlParams;

    if (record.data.recordCount) {
      let routeClusterId = route.queries['clusterId'];
      let finder = routeClusterId ? record.data.records.find(c => c.metadata.name === routeClusterId) : undefined;
      let defaultCluster = finder ? finder : record.data.records[0];
      // 只有在资源列表页才选择具体的集群
      sub === 'sub' && dispatch(clusterActions.selectCluster([defaultCluster]));

      /** =========================== 只有状态的都正常了，才停止轮询 =========================== */
      if (!sub) {
        // 如果当前是打开了创建详情的话，需要选择具体的cluster
        if (dialogState[DialogNameEnum.clusterStatusDialog]) {
          let clusterInfo = record.data.records.find(c => c.metadata.name === cluster.selection.metadata.name);
          clusterInfo && dispatch(clusterActions.selectCluster([clusterInfo]));
        }

        if (record.data.records.filter(item => item.status.phase !== 'Running').length === 0) {
          dispatch(clusterActions.clearPolling());
        }
      } else {
        dispatch(clusterActions.clearPolling());
      }
      /** =========================== 只有状态的都正常了，才停止轮询 =========================== */
    } else {
      dispatch(clusterActions.clearPolling());
    }
  }
});

const restActions = {
  clusterInfo: clusterInfoActions,

  poll: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { cluster, route } = getState();
      dispatch(
        clusterActions.polling({
          filter: Object.assign({}, cluster.query.filter, { regionId: +route.queries['rid'] }),
          delayTime: 8000
        })
      );
    };
  },

  selectCluster(cluster: Cluster[], isNeedInitClusterVersion: boolean = false) {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      cluster[0] && dispatch(FFModelClusterActions.select(cluster[0]));

      // 初始化集群的版本
      isNeedInitClusterVersion && cluster[0] && dispatch(clusterActions.initClusterVersion(cluster[0].status.version));
    };
  },

  /** 初始化集群的版本 */
  initClusterVersion: (clusterVersion: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let k8sVersion = clusterVersion
        .split('.')
        .slice(0, 2)
        .join('.');
      dispatch({
        type: ActionType.ClusterVersion,
        payload: k8sVersion
      });
    };
  },

  /** 获取当前集群开启的Addon */
  fetchClusterAddon: (clusterId, regionId) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { clusterVersion } = getState();
      let clusterInfo: ResourceInfo = resourceConfig(clusterVersion).cluster;
      let response = await WebAPI.fetchExtraResourceList(
        {
          filter: {
            regionId: regionId,
            specificName: clusterId
          }
        },
        clusterInfo,
        false,
        'addons'
      );
      let addons: AddonStatus = {};
      response.records.forEach(item => {
        addons[item.spec.type] = {
          status: item.status.phase,
          name: item.metadata.name
        };
      });
      dispatch({
        type: ActionType.FetchClusterAddons,
        payload: addons
      });
    };
  },
  fetchClustercredential: (clusterId: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let clustercredentialResourceInfo = resourceConfig().clustercredential;
      // 过滤条件
      let k8sQueryObj = {
        fieldSelector: {
          clusterName: clusterId
        }
      };
      k8sQueryObj = JSON.parse(JSON.stringify(k8sQueryObj));
      let response = await WebAPI.fetchResourceList(
        { filter: {}, search: '' },
        clustercredentialResourceInfo,
        false,
        k8sQueryObj
      );
      if (response.records.length) {
        dispatch({
          type: ActionType.FetchClustercredential,
          payload: {
            name: response.records[0].metadata.name,
            clusterName: response.records[0].clusterName,
            caCert: response.records[0].caCert,
            token: response.records[0].token
          }
        });
      }
    };
  },
  clearClustercredential: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.FetchClustercredential,
        payload: {
          name: '',
          clusterName: '',
          caCert: '',
          token: ''
        }
      });
    };
  },
  initClusterAllocationRatio: (cluster: Cluster) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let clusterAllocationRatioEdition = Object.assign({}, initAllcationRatioEdition, { id: uuid() });
      if (cluster.spec.properties && cluster.spec.properties.oversoldRatio) {
        let oversoldRatio = cluster.spec.properties.oversoldRatio;
        clusterAllocationRatioEdition.isUseCpu = oversoldRatio.cpu ? true : false;
        clusterAllocationRatioEdition.isUseMemory = oversoldRatio.memory ? true : false;
        clusterAllocationRatioEdition.cpuRatio = clusterAllocationRatioEdition.isUseCpu ? oversoldRatio.cpu : '';
        clusterAllocationRatioEdition.memoryRatio = clusterAllocationRatioEdition.isUseMemory
          ? oversoldRatio.memory
          : '';
      }
      dispatch({
        type: ActionType.UpdateClusterAllocationRatioEdition,
        payload: clusterAllocationRatioEdition
      });
    };
  },
  updateClusterAllocationRatio: (object: any) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let {
        subRoot: { clusterAllocationRatioEdition }
      } = getState();
      dispatch({
        type: ActionType.UpdateClusterAllocationRatioEdition,
        payload: Object.assign({}, clusterAllocationRatioEdition, object)
      });
    };
  }
};

export const clusterActions = extend({}, FFModelClusterActions, restActions);
