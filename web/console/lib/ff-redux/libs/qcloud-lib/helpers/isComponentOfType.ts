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

export type ComponentType =
  | string
  | (new () => React.Component<any, any>)
  | React.ClassicComponentClass<any>
  | Function;

export function isComponentOfType(component: any, type: ComponentType): component is ComponentType {
  if (!component) return false;
  if (typeof component === 'string') {
    return type === 'text';
  }
  if (typeof type === 'string') {
    return component.type === type;
  }
  if (typeof type === 'function') {
    // Stateless/Class
    if (component.type === type) {
      return true;
    }
    // Inherit
    if (component.type && component.type.prototype) {
      return component.type.prototype instanceof <any>type;
    }
  }
  return false;
}
