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
/**
 * router for nmc
 *
 * 1. a state model
 * 2. a reducer
 * 3. an action creator
 * 4. an decorator
 */
import { Dispatch } from 'redux';

import { appendFunction, ReduxConnectedProps } from '../../qcloud-lib';
import { ReduxAction } from '../../..';

const nmcRouter = seajs.require('router');
const pageManager = seajs.require('nmc/main/pagemanager');

function getFragment() {
  const debug = nmcRouter.debug || '';
  const debugReg = new RegExp('^' + debug.replace('/', '\\/'));
  return nmcRouter.getFragment().replace(debugReg, '');
}

function getInitialState(rule: string) {
  const fragment = getFragment();
  const params = nmcRouter.matchRoute(rule, fragment) || [];
  return { fragment, params };
}

/**
 * the route state
 * */
export interface RouteState {
  /**
   * current route
   **/
  fragment?: string;

  /**
   * parameters
   **/
  params?: string[];
}

const RouterNavigateAction = 'RouterNavigate';

type NavigateAction = ReduxAction<RouteState>;

function generateRouterReducer(rule: string) {
  function routerReducer(state: RouteState = getInitialState(rule), action: NavigateAction) {
    if (action.type === RouterNavigateAction) {
      return action.payload;
    }
    return state;
  }
  return routerReducer;
}

function navigateAction(fragment: string, params: string[]): NavigateAction {
  return {
    type: RouterNavigateAction,
    payload: { fragment, params }
  };
}

function startAction(rule: string) {
  return (dispatch: Dispatch) => {
    const { fragment, params } = getInitialState(rule);
    dispatch(navigateAction(fragment, params));

    nmcRouter.use(rule, (...args: string[]) => {
      const params = args.slice();
      const fragment = getFragment();
      dispatch(navigateAction(fragment, params));

      // 更新导航状态
      const parts = fragment.split('/');
      parts.shift();
      const navArgs = [parts.shift(), parts.shift(), parts];
      pageManager.fragment = fragment;
      pageManager.changeNavStatus.apply(pageManager, navArgs);
    });
  };
}

type ConnectedComponent = React.Component<ReduxConnectedProps, any>;

function generateRouterDecorator(rule: string) {
  function decorator(target: new () => ConnectedComponent) {
    function onMount() {
      const _this = this as ConnectedComponent;
      _this.props.dispatch(startAction(rule));
    }

    function onUnmount() {
      const _this = this as ConnectedComponent;
      nmcRouter.unuse(rule);
    }

    const proto = target.prototype as React.ComponentLifecycle<ReduxConnectedProps, any>;

    proto.componentWillMount = proto.componentDidMount ? appendFunction(proto.componentDidMount, onMount) : onMount;
    proto.componentWillUnmount = proto.componentWillUnmount
      ? appendFunction(proto.componentWillUnmount, onUnmount)
      : onUnmount;
  }
  return decorator;
}

const NAMED_PARAM_REGEX = /\(\/:(\w+)\)/g;

/**
 * NMC 路由，使用方法
 *
 *  // route.ts 配置局部路由
    import { Router } from "@tencent/qcloud-nmc";
    export const router = new Router("/business(/:tabId)");

    // RootState.ts 使用路由状态
    import { RouteState } from "@tencent/qcloud-nmc";
    export interface RootState {
        route?: RouteState;
        ...
    }

    // RootReducer.ts 配置 Reducer
    import { router } from "../router";
    export const RootReducer = combineReducers({
        route: router.getReducer()
    });

    // ExampleContainer.ts 绑定 Router
    import { router } from "../router";
    @router.serve()
    export class ExampleContainer extends React.Component<any, any> {
        ...
    }
 */
export class Router {
  private paramNames: string[];

  constructor(private rule: string, private defaults: { [key: string]: string }) {
    this.paramNames = [];

    this.rule.replace(NAMED_PARAM_REGEX, (match, name) => {
      this.paramNames.push(name);
      return match;
    });
  }

  public getReducer() {
    return generateRouterReducer(this.rule);
  }

  public serve() {
    return generateRouterDecorator(this.rule);
  }

  public buildFragment(params: { [key: string]: string } = {}) {
    // TODO 优化默认路由生成
    return this.rule.replace(NAMED_PARAM_REGEX, (matched, name) => {
      if (typeof params[name] === 'undefined' || params[name] === this.defaults[name]) {
        return '/' + this.defaults[name];
      }
      return '/' + params[name];
    });
  }

  public resolve(state: RouteState) {
    let resolved: any = {};

    const paramNames = this.paramNames.slice();
    const paramValues = state.params.slice();

    paramNames.forEach(x => {
      resolved[x] = paramValues.shift() || this.defaults[x];
    });

    return resolved;
  }

  public navigate(params: { [key: string]: string } = {}) {
    nmcRouter.navigate(this.buildFragment(params));
  }
}
