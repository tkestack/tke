import * as React from 'react';
import { BaseReactProps } from '@tencent/qcloud-lib';

export class MainBodyLayout extends React.Component<BaseReactProps, {}> {
  render() {
    let { className = '', style = {}, children } = this.props;
    let finalStyle = Object.assign({}, style);
    return (
      <div className={'manage-area-main secondary-main ' + className} style={finalStyle}>
        <div className="wrap-mod-box">{children}</div>
      </div>
    );
  }
}
