import { PagingQuery } from '../../../';

// 集合分页
export function collectionPaging<T>(collection: T[], paging: PagingQuery) {
  const start = (paging.pageIndex - 1) * paging.pageSize;
  const end = start + paging.pageSize;
  return collection.slice(start, end);
}
