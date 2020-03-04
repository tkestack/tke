import { createFFListActions, extend, FetchOptions, ReduxAction } from '@tencent/ff-redux';

import { resourceConfig } from '../../../../config';
import { Resource, ResourceFilter, ResourceInfo } from '../../common/models';
import { CommonAPI } from '../../common/webapi';
import * as ActionType from '../constants/ActionType';
import { FFReduxActionName } from '../constants/Config';
import { RootState } from '../models/RootState';
import { router } from '../router';
import { peActions } from './peActions';

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
