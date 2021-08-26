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

/**
 * router for nmc
 *
 * 1. a state model
 * 2. a reducer
 * 3. an action creator
 * 4. an decorator
 */
import * as React from 'react';
import { appendFunction, ReduxConnectedProps, ReduxAction } from '@tencent/ff-redux';
import { buildQueryString, parseQueryString } from './urlUtil';

/**
 * 判断是否和路由跳转之前的一样
 * @param prevpath: string  之前的路由
 * @param currentpath: string 当前的路由
 */
const isInSameModule = (prevpath: string, currentpath: string) => {
  let [prevBusiness, prevModule, ...prevRest] = prevpath.split('/').filter(item => item !== ''),
    [currentBusiness, currentModule, ...currentRest] = currentpath.split('/').filter(item => item !== '');
  if (prevModule !== currentModule && prevModule !== undefined) {
    return false;
  }
  return true;
};

const nmcRouter = seajs.require('router');
const pageManager = seajs.require('nmc/main/pagemanager');

function getFragment() {
  const debug = nmcRouter.debug || '';
  const debugReg = new RegExp('^' + debug.replace('/', '\\/'));
  const fragment = nmcRouter.getFragment().split('?');
  return fragment[0].replace(debugReg, '');
}

function getQueryString() {
  const str = nmcRouter.getFragment().split('?');
  return str[1] ? '?' + str[1] : '';
}

function getInitialState(rule: string) {
  const fragment = getFragment();
  const params = nmcRouter.matchRoute(rule, fragment) || [];
  const queryString = getQueryString();
  const queries = queryString ? parseQueryString(queryString) : {};
  return { fragment, params, queryString, queries };
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

  /**queryString */
  queryString?: string;

  /**queries */
  queries?: { [key: string]: string };
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

let curFragment = '',
  curQueryString = '',
  curQueries = {};

function navigateAction(
  fragment: string,
  params: string[],
  queryString: string,
  queries: { [key: string]: string }
): NavigateAction {
  return {
    type: RouterNavigateAction,
    payload: { fragment, params, queryString, queries }
  };
}

function getCurrentRouterStatus(
  fragment: string,
  params: string[],
  qString: string,
  queries: { [key: string]: string }
) {
  let queryString = qString;
  let flag = false; /**记录是否变化， 默认无变化 */
  if (fragment === curFragment && curQueryString && !queryString) {
    /**path无变化，参数变化 */
    queryString = curQueryString;
  } else {
    /**path变化，重置当前缓存 */
    curFragment = fragment;
    curQueryString = queryString;
    curQueries = queries;
    flag = true;
  }

  return {
    flag,
    queries: curQueries
  };
}

function startAction(rule: string) {
  return dispatch => {
    const { fragment, params, queryString, queries } = getInitialState(rule);
    dispatch(navigateAction(fragment, params, queryString, queries));

    nmcRouter.use(rule, (...args: string[]) => {
      const params = args.slice();
      const fragment = getFragment();

      /* eslint-disable */
      typeof params[params.length - 1] === 'object' ? params.pop() : {}; // nmcRouter parse error
      /* eslint-enable */

      const queryString = getQueryString();
      const queries = parseQueryString(queryString);

      dispatch(navigateAction(fragment, params, queryString, queries));

      const curStatus = getCurrentRouterStatus(fragment, params, queryString, queries);
      if (curStatus.flag) {
        // 更新导航状态
        const parts = fragment.split('/');
        parts.shift();
        const navArgs = [parts.shift(), parts.shift(), parts];
        pageManager.fragment = fragment;
        pageManager.changeNavStatus.apply(pageManager, navArgs);
      } else {
        nmcRouter.navigate(fragment + buildQueryString(curStatus.queries));
      }
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

    proto.componentWillMount = proto.componentWillMount ? appendFunction(proto.componentWillMount, onMount) : onMount;
    proto.componentWillUnmount = proto.componentWillUnmount
      ? appendFunction(proto.componentWillUnmount, onUnmount)
      : onUnmount;
  }
  return decorator;
}

const NAMED_PARAM_REGEX = /\(\/:(\w+)\)/g;

export class Router {
  private paramNames: string[];

  /* eslint-disable */
  constructor(private rule: string, private defaults: { [key: string]: string }) {
    this.paramNames = [];

    this.rule.replace(NAMED_PARAM_REGEX, (match, name) => {
      this.paramNames.push(name);
      return match;
    });
  }
  /* eslint-enable */

  public getReducer() {
    return generateRouterReducer(this.rule);
  }

  public serve() {
    return generateRouterDecorator(this.rule);
  }

  public buildFragment(params: { [key: string]: string } = {}) {
    // TODO 优化默认路由生成
    let temp = this.rule.replace(NAMED_PARAM_REGEX, (matched, name) => {
      if (typeof params[name] === 'undefined' || params[name] === this.defaults[name]) {
        return '/' + this.defaults[name];
      }
      return '/' + params[name];
    });

    temp = temp.replace(/([^:])[\/\\\\]{2,}/, '$1/');
    temp = temp.endsWith('/') ? temp.substring(0, temp.length - 1) : temp;
    return temp;
  }

  public resolve(state: RouteState) {
    const resolved: any = {};

    const paramNames = this.paramNames.slice();
    const paramValues = state.params.slice();

    paramNames.forEach(x => {
      resolved[x] = paramValues.shift() || this.defaults[x];
    });

    return resolved;
  }

  public buildUrl(params: { [key: string]: string } = {}, queries?: { [key: string]: string }) {
    return this.buildFragment(params) + buildQueryString(queries);
  }

  public navigate(params: { [key: string]: string } = {}, queries?: { [key: string]: string }, url?: string) {
    if (url) {
      nmcRouter.navigate(url);
    } else {
      const nextLocationPath = this.buildFragment(params);
      const prevLocationPath = location.pathname;
      if (isInSameModule(prevLocationPath, nextLocationPath)) {
        nmcRouter.navigate(nextLocationPath + buildQueryString(queries));
      }
    }
  }

  public newNavigate({
    params = {},
    queries,
    url
  }: {
    params?: { [key: string]: string };
    queries?: { [key: string]: string };
    url?: string;
  }) {
    if (url) {
      nmcRouter.navigate(url);
    } else {
      const nextLocationPath = this.buildFragment(params);
      const prevLocationPath = location.pathname;
      if (isInSameModule(prevLocationPath, nextLocationPath)) {
        nmcRouter.navigate(nextLocationPath + buildQueryString(queries));
      }
    }
  }
}
