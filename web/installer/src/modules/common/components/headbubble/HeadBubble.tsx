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
import { Bubble, Icon } from '@tencent/tea-component';

export interface HeadBubbleProps extends BaseReactProps {
  /**显示标题 */
  title?: string | JSX.Element;

  /**显示的文本 */
  text?: string | JSX.Element;

  /**气泡显示方式 */
  position?: 'top' | 'bottom' | 'left' | 'right';

  /**对齐方式 */
  align?: 'start' | 'end';

  /** 用于title隐藏 */
  autoflow?: boolean;
}

export class HeadBubble extends React.Component<HeadBubbleProps> {
  render() {
    const { title = '', text = '', position, align, autoflow } = this.props;
    return (
      <div>
        {autoflow ? <span className="text-overflow">{title}</span> : <span>{title}</span>}
        <Bubble placement={position} content={text}>
          <Icon type="info" />
        </Bubble>
      </div>
    );
  }
}
