import * as React from 'react';
import { BaseReactProps } from '@tencent/qcloud-lib';

export class DetailLayout extends React.Component<BaseReactProps, {}> {
  render() {
    let { className = '', children } = this.props;
    return <div className={'tc-panel ' + className}>{children}</div>;
  }
}
