import { combineReducers } from 'redux';

import { reduceToPayload } from '../../libs/qcloud-lib';
import { createBaseReducer } from '../base';
import { QueryState, RecordSet } from '../base/Model';
import { createFFListActionType } from './ActionType';

export function createFFListReducer<T, TFilter, ExtendParamsT = any, TSFilter = any>({
  actionName,
  id,
  displayField,
  valueField,
  groupField,
  initialData
}: {
  actionName: string;
  id?: string;
  displayField?: String | Function;
  valueField?: String | Function;
  groupField?: String | Function;
  initialData?: {
    list?: RecordSet<T, ExtendParamsT>;
    query?: QueryState<TFilter, TSFilter>;
    selection?: T;
    selections?: T[];
  };
});

export function createFFListReducer<T, TFilter, ExtendParamsT = any, TSFilter = any>(
  actionName: string,
  id?: string,
  displayField?: String | Function,
  valueField?: String | Function,
  initialData?: {
    list?: RecordSet<T, ExtendParamsT>;
    query?: QueryState<TFilter, TSFilter>;
    selection?: T;
    selections?: T[];
  }
);

export function createFFListReducer<T, TFilter, ExtendParamsT = any, TSFilter = any>(arg1, ...args) {
  let actionName,
    id,
    displayField,
    valueField,
    groupField = '',
    initialData;
  if (typeof arg1 === 'string') {
    actionName = arg1;
    id = args[0];
    displayField = args[1];
    valueField = args[2];
    initialData = args[3];
  } else {
    actionName = arg1.actionName;
    id = arg1.id;
    displayField = arg1.displayField;
    valueField = arg1.valueField;
    groupField = arg1.groupField;
    initialData = arg1.initialData;
  }
  const ActionType = createFFListActionType(actionName, id);
  const { fetchReducer, queryReducer } = createBaseReducer<RecordSet<T, ExtendParamsT>, TFilter, TSFilter>({
    actionType: ActionType.Base,
    initData:
      initialData && initialData.list
        ? initialData.list
        : {
            recordCount: 0,
            records: [] as T[]
          },
    initQuery: initialData && initialData.query ? initialData.query : null
  });
  const TempReducer = combineReducers({
    list: fetchReducer,
    query: queryReducer,
    initValue: reduceToPayload<string | number>(ActionType.InitValue, null),
    selection: reduceToPayload<T>(
      ActionType.Selection,
      initialData && initialData.selection ? initialData.selection : null
    ),
    initValues: reduceToPayload<string[] | number[]>(ActionType.InitValues, []),
    selections: reduceToPayload<T[]>(
      ActionType.Selections,
      initialData && initialData.selections ? initialData.selections : []
    ),
    displayField: reduceToPayload<String | Function>(ActionType.DisplayField, displayField || ''),
    valueField: reduceToPayload<String | Function>(ActionType.ValueField, valueField || ''),
    groupKeyField: reduceToPayload<String | Function>(ActionType.GroupKeyField, groupField || '')
  });
  return (state, action) => {
    let newState = state;
    switch (action.type) {
      case ActionType.Clear:
        newState = undefined;
        break;
      case ActionType.ClearSelection:
        newState.selection = null;
        newState.selections = [];
        break;
    }
    return TempReducer(newState, action);
  };
}
