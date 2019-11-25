import { extend } from '@tencent/qcloud-lib';
import { generateFetcherActionCreator } from '@tencent/qcloud-redux-fetcher';
import { generateQueryActionCreator } from '@tencent/qcloud-redux-query';
import { RootState } from '../models';
import * as ActionType from '../constants/ActionType';
import * as WebAPI from '../../cluster/WebAPI';
import { ResourceInfo } from '../../common/models';
import { resourceConfig } from '../../../../config';
import { alarmPolicyActions } from './alarmPolicyActions';

type GetState = () => RootState;

/** fetch namespace list */
const fetchNamespaceActions = generateFetcherActionCreator({
  actionType: ActionType.FetchNamespaceList,
  fetcher: async (getState: GetState, fetchOptions, dispatch: Redux.Dispatch) => {
    let { clusterVersion } = getState();
    let namespaceInfo: ResourceInfo = resourceConfig(clusterVersion)['np'];
    let response = await WebAPI.fetchNamespaceList(getState().namespaceQuery, namespaceInfo);
    return response;
  },
  finish: (dispatch: Redux.Dispatch, getState: GetState) => {
    let { namespaceList, route, namespaceQuery } = getState();

    if (namespaceQuery.filter.default) {
      let defauleNamespace =
        (namespaceList.data.recordCount && namespaceList.data.records.find(n => n.name === 'default').name) ||
        'default';

      dispatch(alarmPolicyActions.selectsWorkLoadNamespace(defauleNamespace));
    }
  }
});

/** query namespace list action */
const queryNamespaceActions = generateQueryActionCreator({
  actionType: ActionType.QueryNamespaceList,
  bindFetcher: fetchNamespaceActions
});

export const namespaceActions = extend(fetchNamespaceActions, queryNamespaceActions);
