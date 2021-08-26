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

/* eslint-disable */
declare module Redux {
  interface ActionCreator extends Function {
    (...args: any[]): any;
  }

  interface Reducer extends Function {
    (state: any, action: any): any;
  }

  interface Dispatch extends Function {
    (action: any): any;
  }

  interface StoreMethods {
    dispatch: Dispatch;
    getState(): any;
  }

  interface MiddlewareArg {
    dispatch: Dispatch;
    getState: Function;
  }

  interface Middleware extends Function {
    (obj: MiddlewareArg): Function;
  }

  class Store {
    getReducer(): Reducer;
    replaceReducer(nextReducer: Reducer): void;
    dispatch(action: any): any;
    getState(): any;
    subscribe(listener: Function): Function;
  }

  function createStore(reducer: Reducer, initialState?: any, enhancer?: () => any): Store;
  function bindActionCreators<T>(actionCreators: T, dispatch: Dispatch): T;
  function combineReducers(reducers: any): Reducer;
  function applyMiddleware(...middlewares: Middleware[]): Function;
  function compose<T extends Function>(...functions: Function[]): T;
}

declare module 'redux' {
  export = Redux;
}
