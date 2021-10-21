/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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
import * as classnames from 'classnames';
import * as React from 'react';

import { Bubble, Button, Text } from '@tea/component';
import { BaseReactProps } from '@tencent/ff-redux';

export interface LinkButtonProps extends BaseReactProps {
  /**是否禁用 */
  disabled?: boolean;

  /**点击操作 */
  onClick?: (e) => void;

  tipDirection?: 'top' | 'right' | 'left' | 'bottom';

  /**提示 */
  tip?: string | JSX.Element;

  /**禁用操作提示 只有在禁用时显示*/
  errorTip?: string | JSX.Element;

  title?: string;
  /**是否显示 */
  isShow?: boolean;

  overflow?: boolean;
}

export class LinkButton extends React.Component<LinkButtonProps, {}> {
  render() {
    const {
        disabled,
        onClick,
        tip,
        errorTip,
        children,
        className,
        isShow = true,
        title,
        tipDirection,
        overflow
      } = this.props,
      defaultStyle = { fontSize: '12px', marginRight: '10px', verticalAlign: 'middle' },
      disableStyle = {
        color: 'gray',
        cursor: 'not-allowed',
        fontSize: '12px',
        textDecoration: 'none',
        marginRight: '10px',
        verticalAlign: 'middle'
      };

    let bubbleContent: string | React.ReactNode = null;
    if (disabled) {
      bubbleContent = errorTip ? errorTip : '';
    } else {
      bubbleContent = tip ? tip : '';
    }
    return isShow ? (
      <Button
        title={title}
        type="link"
        className={className + (overflow ? ' tea-text-overflow' : '')}
        disabled={disabled}
        onClick={e => onClick(e)}
      >
        <Bubble placement={tipDirection || 'bottom'} content={bubbleContent || null}>
          <p style={{ marginBottom: 0 }}>{children}</p>
        </Bubble>
      </Button>
    ) : (
      <noscript />
    );
  }
}
