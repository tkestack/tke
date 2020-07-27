import { RootState, ClusterOverview, ClusterOverviewFilter } from './../models/RootState';
import { Dispatch } from 'redux';
import * as ActionType from '../constants/ActionType';
import { FetchState, generateFetcherActionCreator, createFFObjectActions } from '@tencent/ff-redux';

import { cloneDeep } from '../../common/utils';
import * as WebAPI from '../WebAPI';

type GetState = () => RootState;

export const overviewActions = {
  clusterOverActions: createFFObjectActions<ClusterOverview, ClusterOverviewFilter>({
    actionName: ActionType.ClusterOverview,
    fetcher: async (query, getState: GetState) => {
      let response = await WebAPI.fetchClusteroverviews(query);
      return response;
    },
    getRecord: (getState: GetState) => {
      return getState().clusterOverview;
    },
    onFinish: (record, dispatch, getState: GetState) => {}
  })
};

export type OverviewActionsType = typeof overviewActions;
