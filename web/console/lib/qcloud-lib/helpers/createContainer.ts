import * as React from 'react';
import * as classnames from 'classnames';

export function createContainer(
  tag: string,
  baseClassName?: string,
  NestedContainer?: React.ClassicComponentClass<any>
) {
  return class Instance extends React.Component<any, any> {
    render() {
      let { className, children } = this.props;
      className = classnames(baseClassName, className);

      let passProps: any = {};

      for (let key in this.props) {
        if (this.props.hasOwnProperty(key) && key !== 'children') {
          passProps[key] = this.props[key];
        }
      }

      passProps.className = className;

      if (NestedContainer) {
        children = React.createElement(NestedContainer, {}, React.Children.toArray(children));
      }

      return React.createElement(tag, passProps, children);
    }
  };
}
