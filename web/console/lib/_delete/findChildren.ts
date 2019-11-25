import * as React from 'react';
import { isComponentOfType } from './isComponentOfType';

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
