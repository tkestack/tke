import * as React from "react";
import * as ReactDOM from "react-dom";

declare var seajs;
declare var nmc;
let router = seajs.require("router");

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
let currentReactRoot;
let lastComponent;
let container;
/**
 * 表示一个页面
 * */
export default class Page {
  static getRootInstances() {
    return [currentReactRoot];
  }

  [key: string]: any;

  constructor(options) {
    const { businessKey, id, title, Component, getComponent, render } = options;
    this.options = options;
    this.businessKey = businessKey;
    this.id = businessKey + "-" + id;
    this.title = title;
    this.Component = Component;
    this.getComponent = getComponent;
    this.renderer = render;
    this.require = options.require;
    this.requires = options.requires;
    this.destroyed = false;
  }

  async render(...params: any[]) {
    const bee = seajs.require("qccomponent");
    if (bee.utils.ie && bee.utils.ie <= 11) {
      return nmc.render(browserOutOfDate, this.businessKey);
    }

    let Component = this.Component;
    if (!Component) {
      if (typeof params[0] === "object") {
        params.unshift("");
      }
      if (this.getComponent) {
        Component = await this.getComponent(params);
      }
    }
    if (!document.getElementById(this.id)) {
      // 可能加载完依赖之后，用户已经切走了，此时 nmc.render() 不会执行
      nmc.render(`<div id="${this.id}" style="min-height:100%"></div>`, this.businessKey);
    }

    if ((container = document.getElementById(this.id))) {
      currentReactRoot = ReactDOM.render(<Component params={params} />, container);
    }
    // 记录当前页面对象
    lastComponent = Component;
  }

  destroy() {
    if (container && router.getFragment().split("/")[1] !== "argus") {
      // ReactDOM.unmountComponentAtNode &&
      //     ReactDOM.unmountComponentAtNode(container);
      container = null;
    }
  }
}
