import { createFFListActions, extend, uuid } from '@tencent/ff-redux';
import { generateFetcherActionCreator } from '@tencent/qcloud-redux-fetcher';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';

import { resourceConfig } from '../../../../config';
import {
  Cluster, ClusterFilter, CreateResource, Resource, ResourceFilter, ResourceInfo
} from '../../common/models';
import { CommonAPI } from '../../common/webapi';
import * as ActionType from '../constants/ActionType';
import { FFReduxActionName } from '../constants/Config';
import { RootState } from '../models';
import { router } from '../router';
import { editLogStashActions } from './editLogStashActions';
import { logActions } from './logActions';
import { logDaemonsetActions } from './logDaemonsetActions';
import { namespaceActions } from './namespaceActions';
import * as WebAPI from '../WebAPI';
import { projectNamespaceActions } from '@src/modules/cluster/actions/projectNamespaceActions.project';

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
    let { route } = getState();

    if (record.data.recordCount) {
      let lcFinder = record.data.records.find(item => item.spec.type === 'LogCollector');
      if (lcFinder) {
        dispatch(ListAddonActions.select(lcFinder));
        // 判断当前集群是否已经开通了集群日志功能，所有页面下都需要去判断当前的集群是否已经开通过logStash的功能,logList列表的获取会在isOpenLogStash后根据该状态选择拉取与否
        dispatch(
          logDaemonsetActions.applyFilter({
            specificName: lcFinder.metadata.name,
            clusterId: route.queries['clusterId']
          })
        );
      }
    }
  }
});
/** ========================== end 已存在的add的相关操作 ========================== */

/** 地域下的集群列表 */
const fetchClusterActions = generateFetcherActionCreator({
  actionType: ActionType.FetchClusterList,
  fetcher: async (getState: GetState, fetchOptions, dispatch: Redux.Dispatch) => {
    // // 根据地域获取集群列表
    let { clusterQuery } = getState();
    let resourceInfo = resourceConfig()['cluster'];
    let isClearData = fetchOptions && fetchOptions.noCache ? true : false;
    let response = await CommonAPI.fetchResourceList({ query: clusterQuery, resourceInfo, isClearData });
    let agents = await CommonAPI.fetchLogagents();
    let clusterHasLogAgent = {};
    for (let agent of agents.records) {
      clusterHasLogAgent[agent.spec.clusterName] = agent.metadata.name;
    }
    for (let cluster of response.records) {
      cluster.spec.logAgentName = clusterHasLogAgent[cluster.metadata.name];
    }

    return response;
  },
  finish: (dispatch: Redux.Dispatch, getState: GetState) => {
    let { clusterList, route } = getState();
    let clusterId: string = '';

    if (clusterList.data.recordCount) {
      let oldClusterId = route.queries['clusterId'];
      clusterId = clusterList.data.records.find(item => item.metadata.name === oldClusterId)
        ? oldClusterId
        : clusterList.data.records[0].metadata.name;
    }
    // 如果集群列表为空，则选择为空
    dispatch(clusterActions.selectCluster(clusterId, true));
  }
});

/** 地域列表的查询 */
const queryClusterActions = generateQueryActionCreator<ClusterFilter>({
  actionType: ActionType.QueryClusterList,
  bindFetcher: fetchClusterActions
});

const restActions = {
  initProjectList: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { route, projectSelection } = getState();
      let portalResourceInfo = resourceConfig()['portal'];
      let portal = await WebAPI.fetchUserPortal(portalResourceInfo);
      let userProjectList = Object.keys(portal.projects).map(key => {
        return {
          name: key,
          displayName: portal.projects[key]
        };
      });
      dispatch({
        type: ActionType.InitProjectList,
        payload: userProjectList
      });
      let defaultProjectName = projectSelection
        ? projectSelection
        : route.queries['projectName']
          ? route.queries['projectName']
          : userProjectList.length
            ? userProjectList[0].name
            : '';
      defaultProjectName && dispatch(clusterActions.selectProject(defaultProjectName));
    };
  },

  /** 选择业务 */
  selectProject: (project: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { route, projectNamespaceQuery } = getState(),
        urlParams = router.resolve(route);
      dispatch({
        type: ActionType.ProjectSelection,
        payload: project
      });
      router.navigate(urlParams, Object.assign({}, route.queries, { projectName: project }));
      dispatch(namespaceActions.applyFilter({ projectName: project }));
    };
  },
  // 业务侧下，在列表页表格操作区切换命名空间的时候设置一下当前选中的集群（兼容平台侧下的处理逻辑）
  selectClusterFromNamespace: (cluster: Cluster) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { clusterList, route, regionList } = getState(),
        urlParams = router.resolve(route),
        { rid, projectName } = route.queries,
        { mode } = urlParams;
      dispatch({
        type: ActionType.SelectCluster,
        payload: cluster ? [cluster] : []
      });
      let clusterId = cluster.metadata.name;
      router.navigate(urlParams, Object.assign({}, route.queries, { clusterId }));
      // 拉取已经开通的addon的列表
      dispatch(
        ListAddonActions.applyFilter({
          specificName: clusterId
        })
      );
    };
  },
  // 业务侧下，在编辑页切换日志源下的的命名空间时设置选中的集群（同样是兼容平台侧下的处理逻辑）
  selectClusterFromEditNamespace: (cluster: Cluster) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.SelectCluster,
        payload: cluster ? [cluster] : []
      });
    };
  },

  /** 集群的选择-平台侧
   * @params clusterId: string  集群Id
   * @params isNeedInitClusterVersion: boolean | number 是否需要初始化集群的版本
   */
  selectCluster: (clusterId: string, isNeedInitClusterVersion: boolean | number = false) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { clusterList, route, regionList } = getState(),
        urlParams = router.resolve(route),
        { rid, projectName } = route.queries,
        { mode } = urlParams,
        isCreate = mode === 'create',
        isUpdate = mode === 'update',
        isDetail = mode === 'datail';

      //只有在mode为create或者''或者刷新的时候才能进行切换集群的选择
      //如果在创建日志采集规则页面下切换集群,需要清空编辑信息
      if (isCreate) {
        dispatch(editLogStashActions.reInitEdit());
      }

      // 判断当前的cluster是否合法、存在
      let clusterSelection: Cluster;
      if (clusterList.data.recordCount) {
        clusterSelection = clusterList.data.records.find(item => item.metadata.name === clusterId);
      }
      dispatch({
        type: ActionType.SelectCluster,
        payload: clusterSelection ? [clusterSelection] : []
      });

      // 进行路由的更改
      router.navigate(urlParams, Object.assign({}, route.queries, { clusterId }));

      //切换集群，默认搜索的命名空间选择为' '
      dispatch(namespaceActions.selectNamespace(''));

      if (clusterSelection) {
        // 初始化集群的版本
        if (isNeedInitClusterVersion) {
          dispatch(clusterActions.initClusterVersion(clusterSelection.status.version, clusterId));
        }

        dispatch(namespaceActions.applyFilter({ projectName, clusterId, regionId: +rid }));

        // 拉取已经开通的addon的列表
        dispatch(
          ListAddonActions.applyFilter({
            specificName: clusterId
          })
        );

        //在编辑状态或者详情页面下获取具体的log信息
        if (isUpdate || isDetail) {
          dispatch(
            logActions.fetchSpecificLog(
              route.queries['stashName'],
              route.queries['clusterId'],
              route.queries['namespace'],
              mode
            )
          );
        }
      } else {
        dispatch(namespaceActions.fetch({ noCache: true }));
        dispatch(logActions.fetch({ noCache: true }));
      }
    };
  },

  /** 初始化集群的版本 */
  initClusterVersion: (clusterVersion: string, clusterId: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      // 初始化集群的版本
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

  enableLogAgent: (cluster: Cluster) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let resourceInfo = resourceConfig()['logagent'];
      let resource: CreateResource = {
        id: uuid(),
        resourceInfo,
        clusterId: cluster.metadata.name
      };
      let response = await CommonAPI.createLogAgent(resource);
    };
  },

  disableLogAgent: (cluster: Cluster) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let resourceInfo = resourceConfig()['logagent'];
      let resource: CreateResource = {
        id: uuid(),
        resourceInfo,
        clusterId: cluster.metadata.name
      };
      let response = await CommonAPI.deleteLogAgent(resource, cluster.spec.logAgentName);
    };
  },
};

export const clusterActions = extend({}, queryClusterActions, fetchClusterActions, restActions);
