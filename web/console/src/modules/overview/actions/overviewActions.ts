/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

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
