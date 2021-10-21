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
import * as Clipboard from 'clipboard';
import * as React from 'react';

import { BaseReactProps } from '@tencent/ff-redux';
import { Button, Icon, Tooltip } from '@tencent/tea-component';

import { TopTips } from '../toptips';

export interface ClipProps extends BaseReactProps {
  /**复制对象 */
  target?: string;

  /**是否显示 */
  isShow?: boolean;

  /**是否显示操作提示 */
  isShowTip?: boolean;

  /**提示方向 */
  tipDirection?: 'top' | 'right' | 'left' | 'bottom';
}

export class Clip extends React.Component<ClipProps> {
  render() {
    const { target, isShow = true, isShowTip, className, style, children } = this.props;
    return isShow ? (
      <Tooltip title="复制">
        <Icon
          style={{ cursor: 'pointer' }}
          type="copy"
          data-clipboard-action="copy"
          data-clipboard-target={target}
          className="copy-trigger hover-icon"
          onClick={e => e.stopPropagation()}
        />
      </Tooltip>
    ) : (
      <noscript />
    );
  }

  componentDidMount() {
    let clipboard = window['oss_clipboard'];

    if (!clipboard) {
      clipboard = new Clipboard('.copy-trigger');
    }

    clipboard.on('success', e => {
      TopTips({ message: '复制成功', theme: 'success', duration: 1000 });
      e.clearSelection();
    });
    clipboard.on('error', e => {
      TopTips({ message: '复制失败', theme: 'error', duration: 1000 });
    });
  }

  componentWillUnmount() {
    const clipboard = window['oss_clipboard'];

    if (clipboard) {
      clipboard.destroy();
    }
  }
}
