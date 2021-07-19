import { AppEntry } from "./entry";
/**
 * 注册应用路由
 *
 * @param routemap 路由映射，每个路由对应一个路由定义
 *
 * @example
```js
app.routes({
  'cvm': {
    render: () => <CvmIndex />,
  },
  'cvm/overview': {
    render: () => <CvmOverview />,
  }
});
```
 *
 * @todo 支持 React-Hot
 */
export declare function routes(routemap: AppRouteMap, ...args: any[]): void;
/**
 * 路由映射
 */
export interface AppRouteMap {
    [key: string]: AppEntry;
}
