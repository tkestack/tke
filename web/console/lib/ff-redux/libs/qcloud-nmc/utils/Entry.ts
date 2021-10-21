/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */
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
