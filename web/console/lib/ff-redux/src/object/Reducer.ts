import { combineReducers } from 'redux';

import { createFFObjectActionType } from './ActionType';
import { createBaseReducer, QueryState } from '../base';

export function createFFObjectReducer<T, TFilter>(
  actionName: string,
  id?: string,
  initialData?: {
    object?: T;
    query?: QueryState<TFilter>;
  }
) {
  const ActionType = createFFObjectActionType(actionName, id);
  const { fetchReducer, queryReducer } = createBaseReducer<T, TFilter>({
    actionType: ActionType.Base,
    initData: initialData && initialData.object ? initialData.object : null,
    initQuery: initialData && initialData.query ? initialData.query : null
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
