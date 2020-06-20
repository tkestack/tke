import { FetcherState, QueryState } from '../base/Model';

export type FFObjectModel<T = any, TFilter = any> = {
  object?: FetcherState<T>;
  query?: QueryState<TFilter>;
};
