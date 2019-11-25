import { extend } from '@tencent/qcloud-lib';
import { RootState } from '../models';
import { generateFetcherActionCreator } from '@tencent/qcloud-redux-fetcher';
import * as ActionType from '../constants/ActionType';
import { CommonAPI } from '../../common/webapi';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';
import { ClusterFilter, Cluster, initValidator } from '../../common/models';
import { router } from '../router';
import { logActions } from './logActions';
import { namespaceActions } from './namespaceActions';
import { editLogStashActions } from './editLogStashActions';
import { resourceConfig } from '../../../../config';
import { logDaemonsetActions } from './logDaemonsetActions';
import { ResourceListMapForContainerLog } from '../constants/Config';

type GetState = () => RootState;

/** 地域下的集群列表 */
const fetchClusterActions = generateFetcherActionCreator({
  actionType: ActionType.FetchClusterList,
  fetcher: async (getState: GetState, fetchOptions, dispatch: Redux.Dispatch) => {
    // // 根据地域获取集群列表
    let { clusterQuery } = getState();
    let resourceInfo = resourceConfig()['cluster'];
    let isClearData = fetchOptions && fetchOptions.noCache ? true : false;
    let response = await CommonAPI.fetchResourceList({ query: clusterQuery, resourceInfo, isClearData });
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
  /** 集群的选择
   * @params clusterId: string  集群Id
   * @params isNeedInitClusterVersion: boolean | number 是否需要初始化集群的版本
   */
  selectCluster: (clusterId: string, isNeedInitClusterVersion: boolean | number = false) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { clusterList, route, regionList } = getState(),
        urlParams = router.resolve(route),
        { rid } = route.queries,
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

        dispatch(namespaceActions.applyFilter({ clusterId, regionId: +rid }));

        // 判断当前集群是否已经开通了集群日志功能，所有页面下都需要去判断当前的集群是否已经开通过logStash的功能,logList列表的获取会在isOpenLogStash后根据该状态选择拉取与否
        dispatch(
          logDaemonsetActions.applyFilter({
            specificName: clusterId,
            clusterId
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
  }
};

export const clusterActions = extend({}, queryClusterActions, fetchClusterActions, restActions);
