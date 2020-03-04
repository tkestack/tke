declare namespace nmc {
  interface Router {
    /**
     * 跳转到指定 URL
     */
    navigate(url: string): void;

    /**
     * 动态使用路由规则
     * @method use
     * @param {string} rule 路由匹配字符串
     * @param {function} action 路由操作
     * @author techirdliu
     * */
    use(rule: string, action: Function);

    /**
     * 取消动态路由规则的使用
     * @method unuse
     * @param {string} rule 路由的匹配字符串
     * @author techirdliu
     * */
    unuse(rule: string);

    /**
     * 获得当前 URL 路径
     */
    fragment: string;

    /**
     * 获得当前 URL 路径（含 Debug 信息）
     */
    getFragment(): string;

    /**
     * 返回路由匹配的参数
     * @method matchRoute
     * @param  {String} rule 路由规则
     * @param  {String} url 地址
     * @return {Array}  参数数组
     * @author evanyuan
     */
    matchRoute(rule: string, url: string): string[];

    /** 当前调试配置 */
    debug: string;
  }
}
