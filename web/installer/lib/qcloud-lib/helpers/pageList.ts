import { PagingQuery } from '../../..';

export function pageList<T>(list: T[], paging: PagingQuery): T[] {
  const start = (paging.pageIndex - 1) * paging.pageSize;
  const end = start + paging.pageSize;

  return list.slice(start, end);
}
