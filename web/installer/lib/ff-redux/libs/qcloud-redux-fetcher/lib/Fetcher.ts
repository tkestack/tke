import { Dispatch } from 'redux';

import { ReduxAction } from '../../../';
import {
    FetcherAction, FetcherPayload, FetcherState, FetcherTrigger, FetchState
} from '../../../src/base';

// export enum FetchState {
//   /** indicates the data is up to date and ready to use */
//   Ready = "Ready" as any,

//   /** indicates the data is out of date, and the new data is fetching */
//   Fetching = "Fetching" as any,

//   /**
//    * indicates the data is out of date, and the new data fetches failed
//    */
//   Failed = "Failed" as any
// }

// export enum FetcherTrigger {
//   /**
//    * trigger a load operation
//    */
//   Start = "Start" as any,

//   /**
//    * trigger when load for the tolerance duration
//    * */
//   Loading = "Loading" as any,

//   /** trigger a receive operation */
//   Done = "Done" as any,

//   /** trigger a failed result */
//   Fail = "Fail" as any,

//   /** trigger a manual update */
//   Update = "Update" as any,

//   Clear = "Clear" as any
// }

/** state for data fetcher */
// export interface FetcherState<TData> {
//   /**
//    * current fetch state
//    * */
//   fetchState: FetchState;

//   /**
//    * 请求是否已完成
//    */
//   fetched?: boolean;

//   /**
//    * data fetched from the last time
//    * */
//   data?: TData;

//   /**
//    * error object when in fail state
//    */
//   error?: any;

//   /**
//    * If the fetch started for a while, the loading will be true.
//    * You can specific the duration by passing `loadingTolerance` when generating action creator.
//    * If the duration is not specific, loading will be true as well as the fetchState gets to `Fetching`
//    * */
//   loading?: boolean;
// }

/** action payload for trigger */
// export interface FetcherPayload<TData> {
//   trigger: FetcherTrigger;
//   data?: TData;
//   error?: Error;
// }

// export type FetcherAction<TData> = ReduxAction<FetcherPayload<TData>>;

/** generate reducer for fetcher */
export function generateFetcherReducer<TData>({
  actionType,
  initialData,
  resetOnStart
}: {
  actionType: number | string;
  initialData: TData;
  resetOnStart?: boolean;
}) {
  const actionTypes = [
    FetcherTrigger.Start,
    FetcherTrigger.Loading,
    FetcherTrigger.Done,
    FetcherTrigger.Fail,
    FetcherTrigger.Update,
    FetcherTrigger.Clear
  ].map(trigger => actionType + trigger.toString());

  return function FetcherReducer(
    state: FetcherState<TData> = {
      fetchState: FetchState.Ready,
      data: initialData,
      fetched: false
    },
    action: FetcherAction<TData>
  ): FetcherState<TData> {
    if (actionTypes.indexOf(action.type.toString()) === -1) {
      return state;
    }

    const trigger = action.payload.trigger;

    switch (trigger) {
      case FetcherTrigger.Start:
        return {
          fetchState: FetchState.Fetching,
          data: resetOnStart ? initialData : state.data,
          loading: state.loading,
          error: null
        };
      case FetcherTrigger.Loading:
        return {
          fetchState: FetchState.Fetching,
          data: state.data,
          loading: true,
          error: null
        };
      case FetcherTrigger.Done:
        return {
          fetchState: FetchState.Ready,
          fetched: true,
          data: action.payload.data,
          loading: false,
          error: null
        };
      case FetcherTrigger.Fail:
        return {
          fetchState: FetchState.Failed,
          data: initialData,
          fetched: true,
          loading: false,
          error: action.payload.error
        };
      case FetcherTrigger.Update:
        return {
          fetchState: state.fetchState,
          data: action.payload.data,
          loading: state.loading,
          error: null
        };
      case FetcherTrigger.Clear:
        return {
          fetchState: FetchState.Ready,
          fetched: false,
          data: initialData,
          loading: false,
          error: null
        };
    }
    return state;
  };
}

interface FetchOptions {
  /**
   * 是否要求强制无缓存拉取
   */
  noCache?: boolean;

  /**
   * 需要传递的数据
   */
  data?: any;

  /**
   * 是否一次性拉取全部数据
   */
  fetchAll?: boolean;

  maxFetchTimes?: boolean;

  orginData?: any;
}

export interface FetcherActionCreator {
  clearFetch?: () => void;
  fetch(options?: FetchOptions): any;
  update(data: any): any;
}

export function generateFetcherActionCreator<TData>({
  actionType,
  fetcher,
  loadingTolerance,
  finish
}: {
  actionType: number | string;
  fetcher: (getState: () => any, options: FetchOptions, dispatch: Redux.Dispatch) => Promise<TData>;
  loadingTolerance?: number;
  finish?: (dispatch: Redux.Dispatch, getState: () => any) => any;
}): FetcherActionCreator {
  type ActionType = ReduxAction<FetcherPayload<TData>>;

  let syncId = 0;
  let lastLoadingTimeout = 0;

  function start(options?: FetchOptions) {
    return (dispatch: Dispatch, getState: () => any) => {
      const fetchAction: ActionType = {
        type: actionType + (FetcherTrigger.Start as any),
        payload: {
          trigger: FetcherTrigger.Start
        }
      };
      dispatch(fetchAction);

      // keep the action is always dispatch with the latest promise result by `start()`
      const currentSyncId = ++syncId;
      const dispatchOnSync = (action: any) => {
        if (syncId === currentSyncId) {
          dispatch(action);
        }
        clearTimeout(lastLoadingTimeout);
      };

      lastLoadingTimeout = window.setTimeout(() => {
        dispatch(loading());
      }, loadingTolerance || 0);

      const fetched = fetcher(getState, options, dispatch).then(
        data => dispatchOnSync(done(data)),
        error => dispatchOnSync(fail(error))
      );

      if (typeof finish === 'function') {
        fetched.then(() => finish(dispatch, getState));
      }
    };
  }

  function loading(): ActionType {
    return {
      type: actionType + (FetcherTrigger.Loading as any),
      payload: {
        trigger: FetcherTrigger.Loading
      }
    };
  }

  function done(data: TData): ActionType {
    return {
      type: actionType + (FetcherTrigger.Done as any),
      payload: {
        trigger: FetcherTrigger.Done,
        data
      }
    };
  }

  function fail(error: Error): ActionType {
    return {
      type: actionType + (FetcherTrigger.Fail as any),
      payload: {
        trigger: FetcherTrigger.Fail,
        error
      }
    };
  }

  function update(data: TData): ActionType {
    return {
      type: actionType + (FetcherTrigger.Update as any),
      payload: {
        trigger: FetcherTrigger.Update,
        data
      }
    };
  }

  function clear(): ActionType {
    syncId = -1;
    return {
      type: actionType + (FetcherTrigger.Clear as any),
      payload: {
        trigger: FetcherTrigger.Clear
      }
    };
  }

  return { fetch: start, update, clearFetch: clear };
}
