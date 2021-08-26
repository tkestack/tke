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

import * as classnames from 'classnames';
import * as React from 'react';

import { Bubble } from '@tea/component';
import { BaseReactProps } from '@tencent/ff-redux';

export interface FormItemProps extends BaseReactProps {
  /**显示的文本 */
  label?: string | JSX.Element;

  /**是否纯文本显示 */
  isPureText?: boolean;

  /**提示 */
  tips?: string | JSX.Element;

  /**样式 */
  className?: any;

  /**是否显示 */
  isShow?: boolean;

  minWidth?: number;

  isNeedFormInput?: boolean;
}

export class FormItem extends React.Component<FormItemProps, {}> {
  render() {
    const {
      isShow = true,
      label = '',
      isPureText = false,
      tips = '',
      children,
      className,
      style,
      isNeedFormInput = true
    } = this.props;
    return isShow ? (
      <li className={classnames(className, { 'pure-text-row': isPureText })} style={style}>
        <div className="form-label" style={{ minWidth: this.props.minWidth || '80px', verticalAlign: 'top' }}>
          <label>
            {label}
            {tips ? (
              <Bubble placement="top" content={tips || null}>
                <i className="plaint-icon" style={{ marginLeft: '5px' }} />
              </Bubble>
            ) : (
              <noscript />
            )}
          </label>
        </div>
        {isNeedFormInput ? (
          <div className="form-input">{children}</div>
        ) : (
          <div style={{ paddingBottom: 16 }}>{children}</div>
        )}
      </li>
    ) : (
      <noscript />
    );
  }
}
