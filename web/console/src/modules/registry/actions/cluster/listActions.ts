import {
  extend,
  generateWorkflowActionCreator,
  OperationResult,
  OperationTrigger,
  createFFListActions
} from '@tencent/ff-redux';
import { router } from '../../router';
import { RootState, Cluster, ChartInfoFilter, ClusterFilter } from '../../models';
import * as ActionType from '../../constants/ActionType';
import * as WebAPI from '../../WebAPI';
import { namespaceActions } from '../namespace';
type GetState = () => RootState;

/**
 * 列表操作
 */
const fetchClusterActions = createFFListActions<Cluster, ClusterFilter, ChartInfoFilter>({
  actionName: ActionType.ClusterList,
  fetcher: async (query, getState: GetState, fetchOptions, dispatch: Redux.Dispatch) => {
    let response = await WebAPI.fetchClusterList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().clusterList;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    if (record.data.recordCount > 0) {
      dispatch(listActions.selectCluster(record.data.records[0].metadata.name, record.data.data));
    }
  }
});

const restActions = {
  selectCluster: (cluster: string, chartInfoFilter?: ChartInfoFilter) => {
    return async (dispatch, getState: GetState) => {
      dispatch(listActions.selectByValue(cluster));
      dispatch(
        namespaceActions.list.applyFilter({
          cluster: cluster,
          chartInfoFilter: chartInfoFilter
        })
      );
    };
  }
};

export const listActions = extend({}, fetchClusterActions, restActions);
