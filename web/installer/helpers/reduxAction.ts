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

import { Dispatch } from 'redux';

import { FetcherPayload, FetcherTrigger, FetchOptions, ReduxAction } from '@tencent/ff-redux';

interface Fetcher<T> {
  (getState: () => any, options: FetchOptions, dispatch: Redux.Dispatch): Promise<T>;
}

interface FetcherCreator<T> {
  actionType: number | string;
  fetcher: Fetcher<T>;
  loadingTolerance?: number;
}

export interface FetcherActions {
  fetch(options?: FetchOptions, meta?: any): any;
}

export function createFetcherActions<T>(creator: FetcherCreator<T>): FetcherActions {
  type ActionType = ReduxAction<FetcherPayload<T>>;

  let { actionType, fetcher, loadingTolerance } = creator;

  let syncId = 0;
  let lastLoadingTimeout = 0;

  function fetch(options?: FetchOptions, meta?: any) {
    return (dispatch: Dispatch, getState: () => any) => {
      const fetchAction: ActionType = {
        type: actionType + (FetcherTrigger.Start as any),
        payload: {
          trigger: FetcherTrigger.Start
        },
        meta
      };
      dispatch(fetchAction);

      // keep the action is always dispatch with the latest promise result by `fetch()`
      const currentSyncId = ++syncId;
      const dispatchOnSync = (action: any) => {
        if (syncId === currentSyncId) {
          dispatch(action);
        }
        clearTimeout(lastLoadingTimeout);
      };

      lastLoadingTimeout = window.setTimeout(() => {
        dispatch(loading(meta));
      }, loadingTolerance || 0);

      fetcher(getState, options, dispatch).then(
        data => dispatchOnSync(done(data, meta)),
        error => dispatchOnSync(fail(error, meta))
      );
    };
  }

  function loading(meta): ActionType {
    return {
      type: actionType + (FetcherTrigger.Loading as any),
      payload: {
        trigger: FetcherTrigger.Loading
      },
      meta
    };
  }

  function done(data: T, meta?): ActionType {
    return {
      type: actionType + (FetcherTrigger.Done as any),
      payload: {
        trigger: FetcherTrigger.Done,
        data: data
      },
      meta
    };
  }

  function fail(error: Error, meta?): ActionType {
    return {
      type: actionType + (FetcherTrigger.Fail as any),
      payload: {
        trigger: FetcherTrigger.Fail,
        error: error
      },
      meta
    };
  }

  return { fetch };
}
