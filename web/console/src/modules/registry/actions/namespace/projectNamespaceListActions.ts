import {
  extend,
  generateWorkflowActionCreator,
  OperationResult,
  OperationTrigger,
  createFFListActions
} from '@tencent/ff-redux';
import { router } from '../../router';
import { RootState, ProjectNamespace, ProjectNamespaceFilter, ChartInfoFilter } from '../../models';
import * as ActionType from '../../constants/ActionType';
import * as WebAPI from '../../WebAPI';
import { appActions } from '../app';
type GetState = () => RootState;

/**
 * 列表操作
 */
const fetchProjectNamespaceActions = createFFListActions<ProjectNamespace, ProjectNamespaceFilter>({
  actionName: ActionType.ProjectNamespaceList,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchProjectNamespaceList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().projectNamespaceList;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    if (record.data.recordCount > 0) {
      dispatch(
        projectNamespaceListActions.selectProjectNamespace(
          record.data.records[0].metadata.namespace,
          record.data.records[0].spec.clusterName,
          record.data.records[0].spec.namespace,
          record.data.data
        )
      );
    }
  }
});

const restActions = {
  selectProjectNamespace: (
    projectID: string,
    cluster: string,
    namespace: string,
    chartInfoFilter?: ChartInfoFilter
  ) => {
    return async (dispatch, getState: GetState) => {
      let { route, appCreation, chartEditor } = getState(),
        urlParams = router.resolve(route);

      dispatch(projectNamespaceListActions.selectByValue(cluster + '/' + namespace));

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
          dispatch(
            appActions.create.chart.applyFilter({
              cluster: cluster,
              namespace: namespace,
              metadata: { ...chartInfoFilter.metadata },
              chartVersion: chartInfoFilter.chartVersion,
              projectID: projectID
            })
          );
        }
      }
    };
  }
};

export const projectNamespaceListActions = extend({}, fetchProjectNamespaceActions, restActions);
