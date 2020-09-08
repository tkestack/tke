import {
  extend,
  generateWorkflowActionCreator,
  OperationResult,
  OperationTrigger,
  createFFListActions
} from '@tencent/ff-redux';
import { router } from '../../router';
import { RootState, Namespace, NamespaceFilter, ChartInfoFilter } from '../../models';
import * as ActionType from '../../constants/ActionType';
import * as WebAPI from '../../WebAPI';
import { appActions } from '../app';
type GetState = () => RootState;

/**
 * 列表操作
 */
const fetchNamespaceActions = createFFListActions<Namespace, NamespaceFilter>({
  actionName: ActionType.NamespaceList,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchNamespaceList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().namespaceList;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    if (record.data.recordCount > 0) {
      dispatch(listActions.selectNamespace(record.data.records[0].metadata.name, record.data.data));
    }
  }
});

const restActions = {
  selectNamespace: (namespace: string, chartInfoFilter?: ChartInfoFilter) => {
    return async (dispatch, getState: GetState) => {
      let { route, clusterList, appCreation, chartEditor } = getState(),
        urlParams = router.resolve(route);

      dispatch(listActions.selectByValue(namespace));

      const cluster = clusterList.selection ? clusterList.selection.metadata.name : '';
      if (!urlParams['sub'] || urlParams['sub'] === 'chart') {
        if (!urlParams['mode'] || urlParams['mode'] === 'detail') {
          dispatch(
            appActions.create.updateCreationState({
              metadata: Object.assign({}, appCreation.metadata, {
                namespace: namespace
              }),
              spec: Object.assign({}, appCreation.spec, {
                targetCluster: cluster
              })
            })
          );

          //加载values.yaml
          chartInfoFilter &&
            dispatch(
              appActions.create.chart.applyFilter({
                cluster: cluster,
                namespace: namespace,
                metadata: { ...chartInfoFilter.metadata },
                chartVersion: chartInfoFilter.chartVersion,
                projectID: route.queries['prj']
              })
            );
        }
      }
    };
  }
};

export const listActions = extend({}, fetchNamespaceActions, restActions);
