import { Dispatch } from 'redux';
import { extend, merge } from '../../qcloud-lib';
import { QueryState, PagingQuery, ReduxAction, ActionTypesEnum } from '../../../src/base';
import { FetcherActionCreator } from '../../qcloud-redux-fetcher';

// export interface QueryState<TFilter> {
//   date?: DateQuery;
//   paging?: PagingQuery;
//   search?: string;
//   keyword?: string;
//   filter?: TFilter;
//   sort?: SortQuery;
// }

// enum ActionTypesEnum {
//   ChangeDate = "ChangeDate",
//   ChangePaging = "ChangePaging",
//   RestPaging = "RestPaging",
//   ChangeKeyword = "ChangeKeyword",
//   ApplyFilter = "ApplyFilter",
//   ApplyPolling = "ApplyPolling",
//   SortBy = "SortBy",
//   PerformSearch = "PerformSearch",
//   Reset = "Reset"
// }

/* eslint-disable */
const defaultInitialState: QueryState<any> = {
  paging: { pageIndex: 1, pageSize: 20 },
  keyword: '',
  search: '',
  filter: {},
  sort: {},
  continueToken: '',
  pageMapContinueToken: {
    1: ''
  },
  recordCount: Infinity
};

export function generateQueryReducer<TFilter>({
  actionType,
  initialState
}: {
  actionType: string | number;
  initialState?: QueryState<TFilter>;
}) {
  const prefix = actionType as string;
  initialState = merge({}, defaultInitialState, initialState);

  function reducer(state: QueryState<TFilter> = initialState, action: ReduxAction<any>) {
    let { paging, keyword, search, filter, sort, continueToken, pageMapContinueToken, recordCount } = state;

    function resetPaging(newPaging = paging) {
      paging = { pageIndex: 1, pageSize: newPaging.pageSize };
      resetContinue();
      recordCount = Infinity;
    }

    function resetContinue() {
      continueToken = '';
      pageMapContinueToken = {
        1: ''
      };
    }

    switch (action.type) {
      case prefix + ActionTypesEnum.ChangePaging:
        if (action.payload.pageSize !== paging.pageSize) {
          resetPaging(action.payload);
        } else {
          paging = action.payload;
          continueToken = pageMapContinueToken[paging.pageIndex];
        }

        break;

      case prefix + ActionTypesEnum.RestPaging:
        resetPaging();
        break;

      case prefix + ActionTypesEnum.ChangeKeyword:
        keyword = action.payload;
        break;

      case prefix + ActionTypesEnum.PerformSearch:
        search = action.payload;
        keyword = search;
        resetPaging();
        break;

      case prefix + ActionTypesEnum.ApplyFilter:
        if (JSON.stringify(filter) !== JSON.stringify(action.payload)) {
          filter = extend({}, filter, action.payload);
          resetPaging();
        }

        break;

      case prefix + ActionTypesEnum.ApplyPolling:
        filter = extend({}, filter, action.payload);
        break;

      case prefix + ActionTypesEnum.SortBy:
        sort = action.payload;
        break;

      case prefix + ActionTypesEnum.ChangeContinueToken:
        if (!action.payload) {
          recordCount = paging.pageIndex * paging.pageSize;
        } else {
          pageMapContinueToken = {
            ...pageMapContinueToken,
            [paging.pageIndex + 1]: action.payload
          };
        }

        break;

      case prefix + ActionTypesEnum.Reset:
        return initialState;

      default:
        return state;
    }

    return { paging, keyword, search, filter, sort, continueToken, pageMapContinueToken, recordCount };
  }

  return reducer;
}

export interface QueryActionCreator<TFilter> {
  changePaging(nextPaging: PagingQuery): any;
  resetPaging(): any;
  changeKeyword(nextKeyword: string): any;
  performSearch(keyword: string): any;
  applyFilter(nextFilter: TFilter): any;
  applyPolling(nextFilter: TFilter): any;
  sortBy(by: string, desc?: boolean): any;
  changeContinueToken(nextContinueToken: string): any;
  reset(): any;
}

export function generateQueryActionCreator<TFilter>({
  actionType,
  bindFetcher,
  bindFetchers
}: {
  actionType: number | string;
  bindFetcher?: FetcherActionCreator;
  bindFetchers?: Array<{
    actions: ActionTypesEnum[] | 'all';
    fetcher: FetcherActionCreator;
  }>;
}): QueryActionCreator<TFilter> {
  const prefix = actionType as string;

  function changePaging(nextPaging: PagingQuery): ReduxAction<PagingQuery> {
    return {
      type: prefix + ActionTypesEnum.ChangePaging,
      payload: nextPaging
    };
  }

  function resetPaging(): ReduxAction<void> {
    return {
      type: prefix + ActionTypesEnum.RestPaging
    };
  }

  function changeKeyword(nextKeyword: string): ReduxAction<string> {
    return {
      type: prefix + ActionTypesEnum.ChangeKeyword,
      payload: nextKeyword
    };
  }

  function changeContinueToken(nextContinueToken: string): ReduxAction<string> {
    return {
      type: prefix + ActionTypesEnum.ChangeContinueToken,
      payload: nextContinueToken
    };
  }

  function performSearch(nextSearch: string): ReduxAction<string> {
    return {
      type: prefix + ActionTypesEnum.PerformSearch,
      payload: nextSearch
    };
  }

  function applyFilter(nextFilter: TFilter): ReduxAction<TFilter> {
    return {
      type: prefix + ActionTypesEnum.ApplyFilter,
      payload: nextFilter
    };
  }
  function applyPolling(nextFilter: TFilter): ReduxAction<TFilter> {
    return {
      type: prefix + ActionTypesEnum.ApplyPolling,
      payload: nextFilter
    };
  }

  function sortBy(by: string, desc?: boolean) {
    return {
      type: prefix + ActionTypesEnum.SortBy,
      payload: { by, desc }
    };
  }

  function reset() {
    return {
      type: prefix + ActionTypesEnum.Reset,
      payload: ''
    };
  }

  const creator: QueryActionCreator<TFilter> = {
    changePaging,
    applyFilter,
    applyPolling,
    resetPaging,
    changeKeyword,
    performSearch,
    sortBy,
    changeContinueToken,
    reset
  };

  if (bindFetcher || bindFetchers) {
    const wrap = (creator: Function, noFetch = false) => (...args: any[]) => (
      dispatch: Dispatch,
      getState: () => any
    ) => {
      const queryAction: ReduxAction<any> = creator.apply(null, args);
      dispatch(queryAction);

      if (noFetch) return;

      if (bindFetcher) {
        dispatch(bindFetcher.fetch());
      }
      if (bindFetchers) {
        bindFetchers.forEach(bind => {
          if (bind.actions === 'all') {
            dispatch(bind.fetcher.fetch());
          } else if (bind.actions instanceof Array) {
            const bindActions = bind.actions as ActionTypesEnum[];
            if (bindActions.map(x => (actionType as string) + x).indexOf(queryAction.type as string) > -1) {
              dispatch(bind.fetcher.fetch());
            }
          }
        });
      }
    };
    return {
      changeKeyword,
      changePaging: wrap(creator.changePaging),
      resetPaging: wrap(creator.resetPaging),
      applyFilter: wrap(creator.applyFilter),
      applyPolling: wrap(creator.applyPolling),
      performSearch: wrap(creator.performSearch),
      changeContinueToken: wrap(creator.changeContinueToken, true),
      sortBy: wrap(creator.sortBy),
      reset
    };
  }
  return creator;
}
