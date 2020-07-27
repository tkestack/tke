import {
  extend,
  generateWorkflowActionCreator,
  OperationResult,
  OperationTrigger,
  createFFListActions
} from '@tencent/ff-redux';
import { router } from '../../router';
import { RootState, ProjectNamespace, ProjectNamespaceFilter } from '../../models';
import * as ActionTypes from '../../constants/ActionTypes';
import * as WebAPI from '../../WebAPI';
import { appActions } from '../app';
type GetState = () => RootState;

/**
 * 列表操作
 */
const fetchProjectNamespaceActions = createFFListActions<ProjectNamespace, ProjectNamespaceFilter>({
  actionName: ActionTypes.ProjectNamespaceList,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchProjectNamespaceList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().projectNamespaceList;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    if (record.data.recordCount > 0) {
      let { route } = getState();
      const namespace = route.queries['namespace'];
      const cluster = route.queries['cluster'];
      let exist = record.data.records.find(r => {
        return r.spec.namespace === namespace && r.spec.clusterName === cluster;
      });
      if (exist) {
        dispatch(
          projectNamespaceListActions.selectProjectNamespace(
            exist.metadata.namespace,
            exist.spec.clusterName,
            exist.spec.namespace
          )
        );
      } else {
        dispatch(
          projectNamespaceListActions.selectProjectNamespace(
            record.data.records[0].metadata.namespace,
            record.data.records[0].spec.clusterName,
            record.data.records[0].spec.namespace
          )
        );
      }
    }
  }
});

const restActions = {
  selectProjectNamespace: (projectId: string, cluster: string, namespace: string) => {
    return async (dispatch, getState: GetState) => {
      let { route, appCreation } = getState(),
        urlParams = router.resolve(route);
      router.navigate(
        urlParams,
        Object.assign({}, route.queries, { cluster: cluster, namespace: namespace, projectId: projectId })
      );

      dispatch(projectNamespaceListActions.selectByValue(cluster + '/' + namespace));

      if (!urlParams['sub'] || urlParams['sub'] === 'app') {
        if (!urlParams['mode'] || urlParams['mode'] === 'list') {
          dispatch(
            appActions.list.poll({
              cluster: cluster,
              namespace: namespace
            })
          );
        } else if (urlParams['mode'] === 'create') {
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
        }
      }
    };
  }
};

export const projectNamespaceListActions = extend({}, fetchProjectNamespaceActions, restActions);
