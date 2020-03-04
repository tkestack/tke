import * as React from 'react';

import { BaseReactProps } from '@tencent/ff-redux';
import { Alert, AlertProps } from '@tencent/tea-component';

export interface TipInfoProps extends AlertProps {
  /**是否显示组件 */
  isShow?: boolean;

  /** 是否在表单当中去展示 */
  isForm?: boolean;
}

export class TipInfo extends React.Component<TipInfoProps, {}> {
  render() {
    let { style = {}, isShow = true, isForm = false, ...restProps } = this.props,
      renderStyle = style;

    // 用于在创建表单当中 展示错误信息
    if (isForm) {
      renderStyle = Object.assign({}, renderStyle, {
        display: 'inline-block',
        marginLeft: '20px',
        marginBottom: '0px',
        maxWidth: '750px',
        maxHeight: '120px',
        overflow: 'auto'
      });
    }

    return isShow ? (
      <Alert style={renderStyle} {...restProps}>
        {this.props.children}
      </Alert>
    ) : (
      <noscript />
    );
  }
}
