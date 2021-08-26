/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

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
