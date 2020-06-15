import { setProjectName } from '@helper';
import { extend, FetchOptions, generateFetcherActionCreator, RecordSet } from '@tencent/ff-redux';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';

import { resourceConfig } from '../../../../config';
import { Cluster } from '../../common/models';
import { uniq } from '../../common/utils';
import * as ActionType from '../constants/ActionType';
import { FFReduxActionName } from '../constants/Config';
import { RootState } from '../models';
import { router } from '../router';
import * as WebAPI from '../WebAPI';
import { namespaceActions } from './namespaceActions.project';

type GetState = () => RootState;
const fetchOptions: FetchOptions = {
  noCache: false
};

/** fetch namespacesetlist */
const fetchProjectNamespaceActions = generateFetcherActionCreator({
  actionType: ActionType.FetchProjectNamespace,
  fetcher: async (getState: GetState, fetchOptions, dispatch: Redux.Dispatch) => {
    let { projectNamespaceQuery } = getState();
    let response = await WebAPI.fetchProjectNamespaceList(projectNamespaceQuery);
    let clusterList = uniq(response.records.map(namespace => namespace.spec.clusterName));
    dispatch(projectNamespaceActions.initClusterList(clusterList));
    return response;
  },
  finish: async (dispatch: Redux.Dispatch, getState: GetState) => {
    dispatch(namespaceActions.fetch());
  }
});

/** query namespace list action */
const queryProjectNamespaceActions = generateQueryActionCreator({
  actionType: ActionType.QueryProjectNamespace,
  bindFetcher: fetchProjectNamespaceActions
});

const restActions = {
  /** 初始化 NamespaceList列表 */
  initProjectList: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { route, projectSelection } = getState();
      let portalResourceInfo = resourceConfig().portal;
      let portal = await WebAPI.fetchUserPortal(portalResourceInfo);
      let userProjectList = Object.keys(portal.projects).map(key => {
        return {
          name: key,
          displayName: portal.projects[key]
        };
      });
      dispatch({
        type: ActionType.InitProjectList,
        payload: userProjectList
      });
      let defaultProjectName = projectSelection
        ? projectSelection
        : route.queries['projectName']
        ? route.queries['projectName']
        : userProjectList.length
        ? userProjectList[0].name
        : '';
      defaultProjectName && dispatch(projectNamespaceActions.selectProject(defaultProjectName));
    };
  },

  /** 选择业务 */
  selectProject: (project: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { route } = getState(),
        urlParams = router.resolve(route);
      dispatch({
        type: ActionType.ProjectSelection,
        payload: project
      });
      router.navigate(
        urlParams,
        Object.assign(route.queries, {
          projectName: project
        })
      );
      setProjectName(project);
      dispatch(projectNamespaceActions.applyFilter({ specificName: project }));
    };
  },

  /** 初始化集群列表 */
  initClusterList: clusterList => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let result: RecordSet<Cluster> = {
        recordCount: clusterList.length,
        records: []
      };
      result.records = clusterList.map(item => {
        return {
          metadata: { name: item },
          spec: { displayName: '-' },
          status: {}
        };
      });
      dispatch({
        type: FFReduxActionName.CLUSTER + '_BaseDone',
        payload: {
          data: result,
          trigger: 'Done'
        }
      });
      // //业务不一样集群不一定一样，导致不能取url上面的做默认值
      // let defaultCluster = result.records[0] ? result.records[0] : null;

      // defaultCluster && dispatch(projectNamespaceActions.selectCluster(defaultCluster));
    };
  },

  /** 集群的选择 */
  selectCluster: cluster => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let { route } = getState(),
        urlParams = router.resolve(route);
      dispatch({
        type: FFReduxActionName.CLUSTER + '_Selection',
        payload: cluster
      });
      router.navigate(
        urlParams,
        Object.assign(route.queries, {
          clusterId: cluster.metadata.name
        })
      );
    };
  }
};

export const projectNamespaceActions = extend(fetchProjectNamespaceActions, queryProjectNamespaceActions, restActions);
