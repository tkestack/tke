/** 表示一次分页请求 */
export interface PagingQuery {
  /** 请求的页码，从 1 开始索引 */
  pageIndex?: number;

  /** 请求的每页记录数 */
  pageSize?: number;

  append?: boolean;

  /** 翻页清空pages数组 */
  clear?: boolean;
}
