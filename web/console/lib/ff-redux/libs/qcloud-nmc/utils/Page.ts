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

import { createRoot, Root } from 'react-dom/client';

export interface PageOptions {
  businessKey: string;
  title: string;
  id: string;
  component: JSX.Element;
  require: nmc.Require;
  requires: any;
}

const browserOutOfDate = `
<p class="tc-15-msg warning" style="margin: 20px;">
    <span class="tip-info">
        <span class="msg-span">
            <i class="n-error-icon"></i>
            <span>您的浏览器版本太低，请使用
                <a href="https://www.google.com/chrome/" target="_blank">Chrome</a>、
                <a href="https://www.mozilla.org/zh-CN/firefox/new/" target="_blank">FireFox</a>
                或者升级
                <a href="http://windows.microsoft.com/zh-cn/internet-explorer/download-ie" target="_blank">Internet Explorer</a>
                访问</span>
        </span>
    </span>
</p>
`;

// cluezhang: 支持引用外部React库（TEA）时的hot-reload
let currentReactRoot;

/**
 * 表示一个页面
 * */
export class Page {
  private businessKey: string;
  private id: string;
  private title: string;
  private component: JSX.Element;
  private container: Element;
  private require: nmc.Require;
  private requires: any;
  private rendered: boolean;
  private destroyed: boolean;
  private root: Root;

  constructor(options: PageOptions) {
    const { businessKey, id, title, component } = options;
    this.businessKey = businessKey;
    this.id = businessKey + '-' + id;
    this.title = process.env.NODE_ENV === 'development' ? '(dev) ' + title : title;
    this.component = component;
    this.require = options.require;
    this.requires = options.requires;
    this.destroyed = false;
  }

  /**
     * React Hot支持。
     * 通常react-hot-loader需要React一起打包，但我们使用的是外部引入，所以使用时需要一点手工处理才能正确支持。
     * 用法示例（在总模板入口使用）：
       import { Page } from "@tencent/qcloud-nmc";

       ...

       if (module.hot) {
            require('react-hot-loader/Injection').RootInstanceProvider.injectProvider({
                getRootInstances: Page.getRootInstances
            });
        }
     */
  static getRootInstances() {
    return [currentReactRoot];
  }

  render() {
    const bee = seajs.require('qccomponent');
    if (bee.utils.ie && bee.utils.ie <= 8) {
      return nmc.render(browserOutOfDate, this.businessKey);
    }

    const { require, requires } = this;

    function loadDependencies(requires: any): Promise<any> {
      if (requires instanceof Array) {
        return new Promise(resolve => require.async(requires, resolve));
      } else if (typeof requires === 'object') {
        const fatherDependencies = Object.keys(requires);
        Promise.all(fatherDependencies.map(key => loadDependencies(requires[key]))).then(() =>
          loadDependencies(fatherDependencies)
        );
      }
      return Promise.resolve();
    }

    loadDependencies(requires).then(() => {
      // 可能加载完依赖之后，用户已经切走了，此时 nmc.render() 不会执行
      nmc.render(`<div id="${this.id}"></div>`, this.businessKey);
      if ((this.container = document.getElementById(this.id))) {
        const root = createRoot(this.container);
        root.render(this.component);
        this.root = root;
      }
    });
  }

  destroy() {
    if (this.root) {
      this.root.unmount();
    }
  }
}
