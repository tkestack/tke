import { createFFListActions, extend } from '@tencent/ff-redux';
import { FFReduxActionName } from '../../cluster/constants/Config';
import * as WebAPI from '../../cluster/WebAPI';
import { Cluster, ClusterFilter } from '../../common/';
import { RootState } from '../models';

type GetState = () => RootState;

/** 集群列表的Actions */
const FFModelClusterActions = createFFListActions<Cluster, ClusterFilter>({
  actionName: FFReduxActionName.CLUSTER,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchClusterList(query, query.filter.regionId);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().cluster;
  }
});

export const clusterActions = extend({}, FFModelClusterActions);
