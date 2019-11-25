export interface PagingQuery {
  /** 请求的页码，从 1 开始索引 */
  pageIndex?: number;

  /** 请求的每页记录数 */
  pageSize?: number;
}

// 集合分页
export function collectionPaging<T>(collection: T[], paging: PagingQuery) {
  const start = (paging.pageIndex - 1) * paging.pageSize;
  const end = start + paging.pageSize;
  return collection.slice(start, end);
}
