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
