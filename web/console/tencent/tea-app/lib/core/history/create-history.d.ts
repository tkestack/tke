import { BrowserHistoryBuildOptions, History } from "history";
export interface TeaHistoryBuildOptions extends BrowserHistoryBuildOptions {
    controller?: string;
    action?: string;
}
/**
 * 利用控制台的 Router 对象创建兼容 React Router 的 history
 *
 * 参考 createBrowserHistory 实现
 * @see https://github.com/ReactTraining/history/blob/master/modules/createBrowserHistory.js
 */
export declare function createHistory<S>(props: TeaHistoryBuildOptions): History<S>;
