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
import * as React from 'react';
import * as ReactDom from 'react-dom';

import { appendFunction } from '../helpers/appendFunction';

type Handler = (...args: any[]) => any;

interface Instance extends React.Component<any, any> {
  $$outsideClickHandler: Handler;
  $$outsideClickRegisterMethods: Handler[];
}

/**
 * 在组件外部被点击的时候，会调用被装饰的方法
 */
export function OnOuterClick(
  target: React.ComponentLifecycle<any, any>,
  propertyKey: string,
  descriptor: TypedPropertyDescriptor<() => any>
) {
  let onMount = target.componentDidMount;
  let onUnmount = target.componentWillUnmount;
  let callbackMethod = descriptor.value;
  let bindStack = null;
  let errorWarned = false;

  function outerClickHandler(e: MouseEvent) {
    const _this = this as Instance;
    let clickTarget: HTMLElement = null;
    try {
      /* eslint-disable */
      clickTarget = ReactDom.findDOMNode(_this) as HTMLElement;
      /** eslint-enable */
    } catch (err) {}
    if (!clickTarget) {
      if (!errorWarned) {
        const style = 'background-color: #fdf0f0; color: red;';
        const message =
          'OnOuterClick 无法通过绑定的组件找到对应的 DOM，这就意味着你的组件并没有正确地 Unmount。请确保使用 `ReactDOM.unmountComponentAtNode()` 方法来销毁组件，而不是直接清除其使用的 DOM。';
        if (console.groupCollapsed) {
          (console.groupCollapsed as any)('%c%s', style, message);
        } else {
          console.log('%c%s', style, message);
        }
        console.log('绑定组件：%o', target.constructor);
        console.log('绑定方法：%o', callbackMethod);
        if (bindStack) {
          console.info('绑定堆栈: %s', bindStack);
        }
        if (console.groupEnd) {
          console.groupEnd();
        }
      }
      errorWarned = true;
      return;
    }
    let isClickedOutside = !clickTarget.contains(e.target as HTMLElement);
    if (isClickedOutside) {
      _this.$$outsideClickRegisterMethods.forEach(invoke => invoke());
    }
  }

  let bind = function() {
    const _this = this as Instance;
    if (!_this.$$outsideClickHandler) {
      _this.$$outsideClickHandler = outerClickHandler.bind(_this);
      _this.$$outsideClickRegisterMethods = [];
      document.addEventListener('mousedown', _this.$$outsideClickHandler, false);
    }
    _this.$$outsideClickRegisterMethods.push(callbackMethod.bind(_this));
    try {
      bindStack = (new Error('Bind OnOuterClick') as any).stack;
    } catch (err) {}
  };

  let unbind = function() {
    const _this = this as Instance;
    if (_this.$$outsideClickHandler) {
      document.removeEventListener('mousedown', _this.$$outsideClickHandler);
      _this.$$outsideClickHandler = undefined;
      _this.$$outsideClickRegisterMethods = undefined;
    }
  };

  target.componentDidMount = onMount ? appendFunction(onMount, bind) : bind;
  target.componentWillUnmount = onUnmount ? appendFunction(unbind, onUnmount) : unbind;

  return descriptor;
}
