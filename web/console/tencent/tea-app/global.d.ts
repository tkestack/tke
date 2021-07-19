/**
 * 声明具有路由关系的业务模块
 * 
 * 模块遵循 [CMD 规范](https://github.com/seajs/seajs/issues/242)
 * 
 * 业务模块的 moduleId 与业务的路由有约定关系：
 * 
 *   - 路由 `/cvm` 的 moduleId 为 `/modules/cvm/index/index`
 *   - 路由 `/cvm/list` 的 moduleId 为 `modules/cvm/list/list`
 * 
 * 业务模块返回一个业务对象，实现 `render()` 和 `destroy()` 两个生命周期方法
 * 
 * @example
 * 
 ```js
  // 业务进入 https://console.cloud.tencent.com/cvm 时加载并使用此模块
  define('/modules/cvm/index/index', function(require) {
    return {
      render: function() {
        nmc.render('<div>This is cvm home page</div>', 'cvm');
      },
      
      destroy: function() {
        // do some clean up here
      }
    }
  })
  ```
  */
declare function define(moduleId: string, factory: any): void;

/**
 * 提供 nmc 全局变量，主要用于业务渲染
 */
declare const nmc: NMCModule;

/**
 * 直出到全局变量中的登录信息
 */
declare const LOGIN_INFO: GlobalLoginInfo;

/**
 * CMD 环境，由 seajs 提供
 */
declare const seajs: {
  /**
   * 获取已加载过的 CMD 模块，控制台内置模块都在业务加载前已经提前加载了
   */
  require: SeajsRequire;

  /**
   *  异步加载新的 CMD 模块
   */
  use: (moduleId: string, cb: Callback) => void;
};

/**
 * 控制台在全局变量 window 上的主要对象是 CMD 环境和 nmc
 */
declare interface Window {
  define: typeof define;
  seajs: typeof seajs;
  nmc: typeof nmc;
  LOGIN_INFO: typeof LOGIN_INFO;
}

/**
 * 跟路由绑定的业务模块
 */
declare interface RoutedBusinessModule {
  render(): void;
  destroy(): void;
}

/**
 * 全局变量 `nmc` 模块，提供 render 方法给业务渲染内容
 */
declare interface NMCModule {
  /**
   * 将 HTML 片段渲染到业务区域（e.g. #appArea）
   *
   * @param html 要渲染的 HTML 代码
   * @param businessKey 业务 key，必传。业务 key 对应控制台一级路由，如果当前路由跟 businessKey 不对应，则 render 方法不会响应
   */
  render(html: string, businessKey: string): void;
}

/**
 * 表示一个回调函数，回调函数获取的数据类型为 <T>
 */
declare type Callback<T = any> = (data: T) => void;

/**
 * 控制台暴露的模块列表
 *
 * CMD 模块通过 require 获取，全局通过 seajs.require() 获取
 */
declare interface SeajsRequire {
  (moduleId: string): any;
}

declare interface GlobalLoginInfo {
  loginUin: number;
  ownerUin: number;
  appId: number;
  identity: null | {
    subjectType: number;
    authType: number;
    authArea: number;
  };
  area: number;
  countryName: string;
  countryCode: string;
}
