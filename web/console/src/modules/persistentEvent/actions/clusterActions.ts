import { FetchOptions, createFFListActions } from '@tencent/ff-redux';
import { extend, ReduxAction } from '@tencent/qcloud-lib';
import * as ActionType from '../constants/ActionType';
import { RootState } from '../models/RootState';
import { ResourceInfo, ResourceFilter, Resource } from '../../common/models';
import { resourceConfig } from '../../../../config';
import { peActions } from './peActions';
import { FFReduxActionName } from '../constants/Config';
import { CommonAPI } from '../../common/webapi';
import { router } from '../router';

type GetState = () => RootState;
const fetchOptions: FetchOptions = {
  noCache: false
};

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
      let { rid, clusterId } = route.queries;
      let finder = clusterId ? record.data.records.find(c => c.metadata.name === clusterId) : undefined;
      let defaultCluster = finder ? finder : record.data.records[0];
      dispatch(clusterActions.selectCluster(defaultCluster));

      // 这里去拉取pe的列表
      dispatch(peActions.poll({ regionId: +rid }));
    }
  }
});

const restActions = {
  /** 选择集群 */
  selectCluster: (cluster: Resource) => {
    return async (dispatch, getState: GetState) => {
      if (cluster) {
        let { route } = getState(),
          urlParams = router.resolve(route);
        dispatch(ListClusterActions.select(cluster));
        let clusterId: string = cluster.metadata.name;
        // 进行路由的跳转
        router.navigate(urlParams, Object.assign({}, route.queries, { clusterId }));
      }
    };
  },

  /** 国际版 */
  toggleIsI18n: (isI18n: boolean): ReduxAction<boolean> => {
    return {
      type: ActionType.IsI18n,
      payload: isI18n
    };
  }
};

export const clusterActions = extend({}, ListClusterActions, restActions);
