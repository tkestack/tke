import {
  extend,
  generateWorkflowActionCreator,
  OperationResult,
  OperationTrigger,
  createFFListActions
} from '@tencent/ff-redux';
import { router } from '../../router';
import { RootState, Cluster } from '../../models';
import * as ActionTypes from '../../constants/ActionTypes';
import * as WebAPI from '../../WebAPI';
import { namespaceActions } from '../namespace';
type GetState = () => RootState;

/**
 * 列表操作
 */
const fetchClusterActions = createFFListActions<Cluster, void>({
  actionName: ActionTypes.ClusterList,
  fetcher: async (query, getState: GetState, fetchOptions, dispatch: Redux.Dispatch) => {
    let response = await WebAPI.fetchClusterList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().clusterList;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    if (record.data.recordCount > 0) {
      let { route } = getState();
      const cluster = route.queries['cluster'];
      let exist = record.data.records.find(r => {
        return r.metadata.name === cluster;
      });
      if (exist) {
        dispatch(listActions.selectCluster(exist.metadata.name));
      } else {
        dispatch(listActions.selectCluster(record.data.records[0].metadata.name));
      }
    }
  }
});

const restActions = {
  selectCluster: (cluster: string) => {
    return async (dispatch, getState: GetState) => {
      let { route } = getState(),
        urlParams = router.resolve(route);
      router.navigate(urlParams, Object.assign({}, route.queries, { cluster: cluster }));

      dispatch(listActions.selectByValue(cluster));
      dispatch(
        namespaceActions.list.applyFilter({
          cluster: cluster
        })
      );
    };
  }
};

export const listActions = extend({}, fetchClusterActions, restActions);
