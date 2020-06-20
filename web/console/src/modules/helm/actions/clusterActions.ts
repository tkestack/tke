import { createFFListActions, extend } from '@tencent/ff-redux';

import { resourceConfig } from '../../../../config';
import { CommonAPI, Resource, ResourceFilter, ResourceInfo } from '../../common';
import * as ActionType from '../constants/ActionType';
import { FFReduxActionName } from '../constants/Config';
import { RootState } from '../models';
import { router } from '../router';
import { helmActions } from './helmActions';
import { namespaceActions } from './namespaceActions';

type GetState = () => RootState;

/** 集群列表的Actions */
const ListClusterActions = createFFListActions<Resource, ResourceFilter>({
  actionName: FFReduxActionName.CLUSTER,
  fetcher: async query => {
    let clusterInfo: ResourceInfo = resourceConfig()['cluster'];
    let response = await CommonAPI.fetchResourceList({ query, resourceInfo: clusterInfo });
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().listState.cluster;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    let { route } = getState();
    if (record.data.recordCount) {
      let defaultClusterId = route.queries['clusterId'];
      let defaultCluster = record.data.records.find(item => item.metadata.name === defaultClusterId);
      dispatch(clusterActions.selectCluster(defaultCluster || record.data.records[0]));
    }
  }
});

const restActions = {
  selectCluster: (cluster: Resource) => {
    return async (dispatch, getState: GetState) => {
      let { route } = getState(),
        urlParams = router.resolve(route);
      router.navigate(urlParams, Object.assign({}, route.queries, { clusterId: cluster.metadata.name }));
      dispatch(ListClusterActions.select(cluster));

      dispatch({
        type: ActionType.FetchHelmList + 'Done',
        payload: {
          data: {
            recordCount: 0,
            records: []
          },
          trigger: 'Done'
        }
      });
      dispatch({
        type: ActionType.FetchInstallingHelmList + 'Done',
        payload: {
          data: {
            recordCount: 0,
            records: []
          },
          trigger: 'Done'
        }
      });
      dispatch(helmActions.checkClusterHelmStatus());
      /// #if tke
      dispatch(namespaceActions.applyFilter({ clusterId: cluster.metadata.name }));

      /// #endif
    };
  }
};

export const clusterActions = extend({}, ListClusterActions, restActions);
