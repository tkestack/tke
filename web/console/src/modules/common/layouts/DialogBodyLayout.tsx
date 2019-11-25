import * as React from 'react';
import { BaseReactProps } from '@tencent/qcloud-lib';

export class DialogBodyLayout extends React.Component<BaseReactProps, {}> {
  render() {
    return (
      <div className="docker-dialog jiqun">
        <div className="act-outline">{this.props.children}</div>
      </div>
    );
  }
}
