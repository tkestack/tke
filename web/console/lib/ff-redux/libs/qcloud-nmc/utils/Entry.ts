import { Page } from './Page';

/* eslint-disable */
export namespace Entry {
  export interface ModuleRequiresTree {
    [parentDependency: string]: string[] | ModuleRequiresTree;
  }

  export interface ModuleConfig {
    /**
     * 模块标题，将体现在浏览器的标题栏上
     * */
    title: string;

    /**
     * 要渲染的容器组件，容器组件应该已经链接到 Redux 中
     * */
    container: JSX.Element;

    /**
     * 需要依赖的 NMC 模块
     */
    requires?: string[] | ModuleRequiresTree;
  }

  export interface Registration {
    /**
     * 业务 Key，对应唯一业务，业务入口的 URL 通过该 Key 推导
     * */
    businessKey: string;

    /**
     * 路由配置，index 必配。其他路由名称对应业务二级路由，如：
     *
     * {
     *     index: {...},  // 对应 `https://console.qcloud.com/{businessKey}`
     *     foo: {...},    // 对应 `https://console.qcloud.com/{businessKey}/foo`
     *     bar: {...}     // 对应 `https://console.qcloud.com/{businessKey}/bar`
     * }
     *
     *
     * */
    routes: {
      /**
       * 主页路由配置，对应 URL 为 `https://console.qcloud.com/{businessKey}`
       * */
      index?: ModuleConfig;

      /**
       * 其他业务模块路由定义
       * */
      [moduleKey: string]: ModuleConfig;
    };
  }

  /**
   * 注册控制台模块
   * */
  export function register({ businessKey, routes }: Registration) {
    for (const moduleKey in routes) {
      if (routes.hasOwnProperty(moduleKey)) {
        const modulePath = `modules/${businessKey}/${moduleKey}/${moduleKey}`;
        const moduleConfig = routes[moduleKey];
        window.define(modulePath, require => {
          return new Page({
            businessKey,
            id: moduleKey,
            title: moduleConfig.title,
            component: moduleConfig.container,
            require,
            requires: moduleConfig.requires
          });
        });
      }
    }
  }
}
