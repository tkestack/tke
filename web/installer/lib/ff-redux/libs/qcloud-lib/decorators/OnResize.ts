/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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

import { appendFunction } from '../helpers/appendFunction';

/**
 * 在窗口 resize 的时候，会调用被装饰的方法
 */
export function OnResize(
  target: React.ComponentLifecycle<any, any>,
  propertyKey: string,
  descriptor: TypedPropertyDescriptor<() => any>
) {
  let onMount = target.componentDidMount;
  let onUnmount = target.componentWillUnmount;
  let callbackMethod = descriptor.value;

  function resizeHandler(e: MouseEvent) {
    callbackMethod.call(this);
  }

  let bind = function() {
    this.$$resizeHandler = resizeHandler.bind(this);
    window.addEventListener('resize', this.$$resizeHandler, false);
  };

  let unbind = function() {
    window.removeEventListener('resize', this.$$resizeHandler);
  };

  target.componentDidMount = onMount ? appendFunction(onMount, bind) : bind;
  target.componentWillUnmount = onUnmount ? appendFunction(onUnmount, unbind) : unbind;

  return descriptor;
}
