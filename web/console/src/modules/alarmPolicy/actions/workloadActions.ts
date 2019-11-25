import { extend } from '@tencent/qcloud-lib';
import { generateFetcherActionCreator } from '@tencent/qcloud-redux-fetcher';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';
import { RootState } from '../models';
import * as ActionType from '../constants/ActionType';
import * as WebAPI from '../../cluster/WebAPI';
import { resourceConfig } from '../../../../config';

type GetState = () => RootState;

/** fetch workload list */
const fetchWorkloadActions = generateFetcherActionCreator({
  actionType: ActionType.FetchWorkloadList,
  fetcher: async (getState: GetState, fetchOptions, dispatch: Redux.Dispatch) => {
    let { workloadQuery, clusterVersion } = getState();
    let { filter } = workloadQuery;
    let resourceInfo = resourceConfig(clusterVersion)[filter.workloadType];
    let response = await WebAPI.fetchResourceList(getState().workloadQuery, resourceInfo);
    return response;
  }
});

/** query Pod list action */
const queryWorkloadActions = generateQueryActionCreator({
  actionType: ActionType.QueryWorkloadList,
  bindFetcher: fetchWorkloadActions
});

export const workloadActions = extend(fetchWorkloadActions, queryWorkloadActions);
