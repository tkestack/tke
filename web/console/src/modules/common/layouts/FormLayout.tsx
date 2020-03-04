import * as React from 'react';

import { BaseReactProps } from '@tencent/ff-redux';

export class FormLayout extends React.Component<BaseReactProps, {}> {
  render() {
    let { className, children, style } = this.props;

    return (
      <div className="tc-panel" style={style}>
        {children}
      </div>
    );
  }
}
