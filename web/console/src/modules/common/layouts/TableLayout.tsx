import * as React from 'react';

import { BaseReactProps } from '@tencent/ff-redux';

export class TableLayout extends React.Component<BaseReactProps, {}> {
  render() {
    let { className, children } = this.props;
    return <div className="tc-panel panel-table">{children}</div>;
  }
}
