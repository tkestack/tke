import * as React from 'react';
import * as ReactDom from 'react-dom';
import { appendFunction } from '../helpers/appendFunction';

/* eslint-disable */
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

  function outerClickHandler(e: MouseEvent) {
    const _this = this as Instance;
    let isClickedOutside = !ReactDom.findDOMNode(_this).contains(e.target as HTMLElement);
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
  target.componentWillUnmount = onUnmount ? appendFunction(onUnmount, unbind) : unbind;

  return descriptor;
}
