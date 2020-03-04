import { FetcherState, QueryState, RecordSet } from '../base/Model';

export type FFListModel<T = any, TFilter = any, ExtendParamsT = any, TSFilter = any> = {
  list?: FetcherState<RecordSet<T, ExtendParamsT>>;
  query?: QueryState<TFilter, TSFilter>;
  initValue?: string | number;
  selection?: T;
  initValues?: string[] | number[];
  selections?: T[];
  displayField?: String | Function;
  valueField?: String | Function;
  groupKeyField?: String | Function;
};
