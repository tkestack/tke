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
import { isComponentOfType } from '../helpers/isComponentOfType';

type FindType = string | (new () => React.Component<any, any>) | React.ClassicComponentClass<any> | Function;

/**
 * 查找指定标签的组件
 */
export function findChildren(children, type: string): JSX.Element[];

/**
 * 查找文本碎片
 */
export function findChildren(children, type: 'text'): JSX.Element[];

/**
 * 查找指定类型的组件
 */
export function findChildren<P, S>(children, type: new () => React.Component<P, S>): JSX.Element[];

/**
 * 查找指定类型的组件
 */
export function findChildren<P>(children, type: React.ClassicComponentClass<P>): JSX.Element[];

/**
 * 查找指定的 Stateless 组件
 * */
export function findChildren<P>(children, type: Function): JSX.Element[];

export function findChildren(children, type: FindType): JSX.Element[] {
  let found = [];
  React.Children.forEach(children, (child: JSX.Element) => {
    if (isComponentOfType(child, type)) {
      found.push(child);
    }
  });
  return found;
}
