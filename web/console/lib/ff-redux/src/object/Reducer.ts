/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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
import { combineReducers } from 'redux';

import { createFFObjectActionType } from './ActionType';
import { createBaseReducer, QueryState } from '../base';

export interface CreateFFObjectReducerProps<T, TFilter> {
  actionName: string;
  id?: string;

  initialData?: {
    object?: T;
    query?: QueryState<TFilter>;
  };
}
export function createFFObjectReducer<T, TFilter>(props?: CreateFFObjectReducerProps<T, TFilter>) {
  const { actionName, id, initialData } = props;
  const ActionType = createFFObjectActionType(actionName, id);
  const { fetchReducer, queryReducer } = createBaseReducer<T, TFilter>({
    actionType: ActionType.Base,
    initData: initialData?.object ?? null,
    initQuery: initialData?.query ?? null
  });

  const TempReducer = combineReducers({
    object: fetchReducer,
    query: queryReducer
  });
  return (state, action) => {
    let newState = state;
    switch (action.type) {
      case ActionType.Clear:
        newState = undefined;
        break;
    }
    return TempReducer(newState, action);
  };
}
