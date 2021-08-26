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

import { extend, merge } from '../../libs/qcloud-lib';
import {
    ActionTypesEnum, DataType, FetcherAction, FetcherState, FetcherTrigger, FetchState, PagingQuery,
    QueryState, ReduxAction
} from './Model';

export interface BaseReducerParams<TData, TFilter, TSFilter = any> {
  actionType: string;
  initData?: TData;
  initQuery?: QueryState<TFilter, TSFilter>;
  dataType?: DataType;
}
export interface BaseReducer<TData, TFilter, TSFilter = any> {
  fetchReducer: any;
  queryReducer: any;
}

const defaultInitialState: QueryState<any, any> = {
  paging: { pageIndex: 1, pageSize: 20 },
  keyword: '',
  search: '',
  filter: {},
  searchFilter: {},
  sort: {}
};

/* eslint-disable */
export function createBaseReducer<TData, TFilter, TSFilter = any>({
  actionType,
  initData,
  initQuery,
  dataType = DataType.List
}: BaseReducerParams<TData, TFilter>): BaseReducer<TData, TFilter, TSFilter> {
  const actionTypes = Object.keys(FetcherTrigger).map(trigger => actionType + trigger.toString());
  const fetchReducer = function(
    state: FetcherState<TData> = {
      fetchState: FetchState.Ready,
      data: initData,
      fetched: false,
      loading: false,
      error: null,
      pages: []
    },
    action: FetcherAction<TData>
  ): FetcherState<TData> {
    let { fetchState, data, loading, error, fetched, pages } = state;
    if (actionTypes.indexOf(action.type.toString()) === -1) {
      return state;
    }

    switch (action.payload.trigger) {
      case FetcherTrigger.Start:
        fetchState = FetchState.Fetching;
        error = null;
        break;
      case FetcherTrigger.Loading:
        fetchState = FetchState.Fetching;
        loading = true;
        error = null;
        break;
      case FetcherTrigger.Done:
        if (typeof action.payload.pageIndex === 'number' && dataType === DataType.List) {
          fetchState = FetchState.Ready;
          if (action.payload.clear) {
            pages = [];
          }
          //这里-1是因为query.pageing.pageIndex是从1开始的
          pages[action.payload.pageIndex - 1] = {
            fetchState: FetchState.Ready,
            fetched: true,
            data: action.payload.data,
            loading: false,
            error: null
          };

          fetched = true;
          if (action.payload.append) {
            data = {} as any;
            data['recordCount'] = action.payload.data['recordCount'];
            data['continue'] = action.payload.data['continue'];
            data['continueToken'] = action.payload.data['continueToken'];
            let records = [];
            pages.forEach(page => {
              records = records.concat(page.data['records']);
            });
            data['records'] = records;
          } else {
            data = action.payload.data;
          }

          loading = false;
          error = null;
        } else {
          fetchState = FetchState.Ready;
          fetched = true;
          data = action.payload.data;
          loading = false;
          error = null;
          pages = [];
        }
        break;
      case FetcherTrigger.Fail:
        if (typeof action.payload.pageIndex === 'number' && dataType === DataType.List) {
          fetchState = FetchState.Failed;
          if (action.payload.clear) {
            pages = [];
          }
          //这里-1是因为query.pageing.pageIndex是从1开始的
          pages[action.payload.pageIndex - 1] = {
            fetchState: FetchState.Failed,
            fetched: true,
            data: initData,
            loading: false,
            error: action.payload.error
          };

          fetched = true;
          if (action.payload.append) {
            data = {} as any;
            data['recordCount'] = action.payload.data['recordCount'];
            data['continue'] = action.payload.data['continue'];
            data['continueToken'] = action.payload.data['continueToken'];
            let records = [];
            pages.forEach(page => {
              records = records.concat(page.data['records']);
            });
            data['records'] = records;
          } else {
            data = initData;
          }

          loading = false;
          error = action.payload.error;
        } else {
          fetchState = FetchState.Failed;
          data = initData;
          fetched = true;
          loading = false;
          error = action.payload.error;
          pages = [];
        }
        break;
      case FetcherTrigger.Update:
        fetchState = FetchState.Ready;
        fetched = true;
        data = action.payload.data;
        loading = false;
        error = null;
        break;
      case FetcherTrigger.Clear:
        fetchState = FetchState.Ready;
        fetched = false;
        data = initData;
        loading = false;
        error = null;
        pages = [];
        break;
    }
    return { fetchState, data, loading, error, fetched, pages };
  };

  initQuery = merge({}, defaultInitialState, initQuery);

  function resetPaging(paging: PagingQuery): PagingQuery {
    return {
      pageIndex: 1,
      pageSize: paging.pageSize,
      append: false,
      clear: true
    };
  }

  const queryReducer = function(
    state: QueryState<TFilter, TSFilter> = initQuery,
    action: ReduxAction<any>
  ): QueryState<TFilter, TSFilter> {
    let { paging, keyword, search, filter, sort, searchFilter, continueToken } = state;

    continueToken = '';

    switch (action.type) {
      //设置搜索关键词，不会执行查询
      case actionType + ActionTypesEnum.ChangeKeyword:
        keyword = action.payload;
        break;
      //设置关键词并执行查询
      case actionType + ActionTypesEnum.PerformSearch:
        search = action.payload;
        keyword = search;
        paging = resetPaging(paging);
        break;

      //设置查询条件，不会执行查询
      case actionType + ActionTypesEnum.ChangeFilter:
        filter = extend({}, filter, action.payload);
        break;
      //设置查询条件并执行查询
      case actionType + ActionTypesEnum.ApplyFilter:
        filter = extend({}, filter, action.payload);
        paging = resetPaging(paging);
        break;
      case actionType + ActionTypesEnum.ApplyPolling:
        filter = extend({}, filter, action.payload);
        break;

      //tagSearchFilter
      //设置搜索关键词，不会执行查询
      case actionType + ActionTypesEnum.ChangeSearchFilter:
        searchFilter = extend({}, searchFilter, action.payload);
        break;
      //设置关键词并执行查询
      case actionType + ActionTypesEnum.ApplySearchFilter:
        searchFilter = extend({}, searchFilter, action.payload);
        paging = resetPaging(paging);
        break;

      // 设置轮询的token
      case actionType + ActionTypesEnum.ChangeContinueToken:
        continueToken = action.payload;
        break;

      //修改页码并重新加载数据
      case actionType + ActionTypesEnum.ChangePaging:
        paging = action.payload;
        break;

      // 修改页码
      case actionType + ActionTypesEnum.ChangePagingIndex:
        paging = {
          pageIndex: action.payload,
          pageSize: paging.pageSize,
          append: false,
          clear: false
        };
        break;

      //重置页码并重新加载数据
      case actionType + ActionTypesEnum.RestPaging:
        paging = resetPaging(paging);
        break;
      case actionType + ActionTypesEnum.Next:
        paging = {
          pageIndex: paging.pageIndex + 1,
          pageSize: paging.pageSize,
          append: true,
          clear: false
        };
        break;

      case actionType + ActionTypesEnum.SortBy:
        sort = action.payload;
        break;

      case actionType + ActionTypesEnum.Reset:
        return initQuery;

      default:
        return state;
    }

    return {
      paging,
      keyword,
      search,
      filter,
      searchFilter,
      sort,
      continueToken
    };
  };

  return {
    fetchReducer,
    queryReducer
  };
}
