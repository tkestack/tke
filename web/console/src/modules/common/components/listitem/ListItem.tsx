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

export interface ListItemProps extends BaseReactProps {
  /**显示的标题文本 */
  label?: string | JSX.Element;

  /**是否显示 */
  isShow?: boolean;

  /**提示 */
  tips?: string | JSX.Element;
}

export class ListItem extends React.Component<ListItemProps, {}> {
  render() {
    const { label, tips, isShow = true, children } = this.props;
    return isShow ? (
      <li style={{ fontSize: '12px' }}>
        <span className="item-descr-tit">
          <span style={{ verticalAlign: 'middle' }}>{label}</span>
          {tips && (
            <Bubble placement="left" content={<p style={{ whiteSpace: 'normal' }}>{tips}</p>}>
              <i className="plaint-icon" style={{ verticalAlign: 'middle' }} />
            </Bubble>
          )}
        </span>
        <span className="item-descr-txt">{children}</span>
      </li>
    ) : (
      <noscript />
    );
  }
}
