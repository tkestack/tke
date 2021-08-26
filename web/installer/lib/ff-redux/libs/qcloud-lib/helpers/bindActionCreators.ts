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
 * like redux.bindActionCreators but do it recurisivly
 * */
export function bindActionCreators<T>(actions: T, dispatch: Redux.Dispatch): T {
  if (typeof actions !== 'object' || !actions) {
    throw new RangeError('invalid actions!');
  }

  let result: any = {};

  for (let key in actions) {
    if (!actions.hasOwnProperty(key)) {
      continue;
    }
    const creator = actions[key];
    if (typeof creator === 'object' && creator) {
      result[key] = bindActionCreators(creator, dispatch);
    } else if (typeof creator === 'function') {
      result[key] = (...args: any[]) => dispatch(creator(...args));
    }
  }

  return result as T;
}
