import { createFFListActions, extend } from '@tencent/ff-redux';

import { Cluster, ClusterFilter, RootState } from '../models';
import * as WebAPI from '../WebAPI';

type GetState = () => RootState;

/** 集群列表的Actions */
const FFModelClusterActions = createFFListActions<Cluster, ClusterFilter>({
  actionName: 'cluster',
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchClusterList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().cluster;
  },
  onFinish: (record, dispatch, getState: GetState) => {}
});

const restActions = {
  selectCluster: (cluster: Cluster[]) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch(FFModelClusterActions.select(cluster[0]));
    };
  }
};

export const clusterActions = extend({}, FFModelClusterActions, restActions);
