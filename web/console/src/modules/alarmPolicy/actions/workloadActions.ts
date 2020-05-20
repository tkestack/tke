import { extend } from '@tencent/ff-redux';
import { generateFetcherActionCreator } from '@tencent/qcloud-redux-fetcher';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';

import { resourceConfig } from '../../../../config';
import * as WebAPI from '../../cluster/WebAPI';
import * as ActionType from '../constants/ActionType';
import { RootState } from '../models';

type GetState = () => RootState;

/** fetch workload list */
const fetchWorkloadActions = generateFetcherActionCreator({
  actionType: ActionType.FetchWorkloadList,
  fetcher: async (getState: GetState, fetchOptions, dispatch: Redux.Dispatch) => {
    let { workloadQuery, clusterVersion } = getState();
    let { filter } = workloadQuery;
    let workloadTypeMap = {
      Deployment: 'deployment',
      StatefulSet: 'statefulset',
      DaemonSet: 'daemonset',
      TApp: 'tapp'
    };
    let resourceInfo = resourceConfig(clusterVersion)[workloadTypeMap[filter.workloadType]];
    let response = await WebAPI.fetchResourceList(getState().workloadQuery, { resourceInfo });
    return response;
  }
});

/** query Pod list action */
const queryWorkloadActions = generateQueryActionCreator({
  actionType: ActionType.QueryWorkloadList,
  bindFetcher: fetchWorkloadActions
});

export const workloadActions = extend(fetchWorkloadActions, queryWorkloadActions);
