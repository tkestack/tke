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
