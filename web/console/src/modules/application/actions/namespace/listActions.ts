import {
  extend,
  generateWorkflowActionCreator,
  OperationResult,
  OperationTrigger,
  createFFListActions
} from '@tencent/ff-redux';
import { router } from '../../router';
import { RootState, Namespace, NamespaceFilter } from '../../models';
import * as ActionTypes from '../../constants/ActionTypes';
import * as WebAPI from '../../WebAPI';
import { appActions } from '../app';
type GetState = () => RootState;

/**
 * 列表操作
 */
const fetchNamespaceActions = createFFListActions<Namespace, NamespaceFilter>({
  actionName: ActionTypes.NamespaceList,
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchNamespaceList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().namespaceList;
  },
  onFinish: (record, dispatch: Redux.Dispatch, getState: GetState) => {
    if (record.data.recordCount > 0) {
      let { route } = getState();
      const namespace = route.queries['namespace'];
      let exist = record.data.records.find(r => {
        return r.metadata.name === namespace;
      });
      if (exist) {
        dispatch(listActions.selectNamespace(exist.metadata.name));
      } else {
        dispatch(listActions.selectNamespace(record.data.records[0].metadata.name));
      }
    }
  }
});

const restActions = {
  selectNamespace: (namespace: string) => {
    return async (dispatch, getState: GetState) => {
      let { route, clusterList, appCreation } = getState(),
        urlParams = router.resolve(route);
      router.navigate(urlParams, Object.assign({}, route.queries, { namespace: namespace }));

      dispatch(listActions.selectByValue(namespace));

      const cluster = clusterList.selection ? clusterList.selection.metadata.name : '';
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

export const listActions = extend({}, fetchNamespaceActions, restActions);
