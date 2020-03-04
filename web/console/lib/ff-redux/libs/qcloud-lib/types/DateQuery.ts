export interface DateQuery {
  /**
   * 查询的起始日期
   */
  from?: string;

  /**
   * 查询的结束日期
   */
  to?: string;

  /**
   * 跨越的天数
   */
  length?: number;

  /**
   * 区间的时间粒度
   */
  stride?: number;
}
