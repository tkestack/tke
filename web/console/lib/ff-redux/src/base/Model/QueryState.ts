import { PagingQuery } from './PagingQuery';
import { SortQuery } from './SortQuery';

export interface QueryState<TFilter, TSFilter = any> {
  paging?: PagingQuery;
  search?: string;
  keyword?: string;
  filter?: TFilter;
  searchFilter?: TSFilter;
  sort?: SortQuery;
  continueToken?: string;
}
