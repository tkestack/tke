export interface FetchOptions {
  /**
   * 是否要求强制无缓存拉取
   */
  noCache?: boolean;

  /**
   * 需要传递的数据
   */
  data?: any;

  /**
   * 是否一次性拉取全部数据
   */
  fetchAll?: boolean;

  maxFetchTimes?: boolean;

  orginData?: any;
}
