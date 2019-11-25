import * as React from 'react';
import { BaseReactProps } from '@tencent/qcloud-lib';

export class MainTitleLayout extends React.Component<BaseReactProps, {}> {
  render() {
    let { className, children } = this.props;
    return <div className="manage-area-title secondary-title">{children}</div>;
  }
}
