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
import { findDOMNode } from 'react-dom';

import { BaseReactProps, OnOuterClick } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

export interface SidePanelProps extends BaseReactProps {
  /**侧边面板标题 */
  title?: string | JSX.Element;

  /**关闭操作 */
  onClose?: () => void;

  /**左侧宽度 默认795px */
  width?: string | number;
}

export class SidePanel extends React.Component<SidePanelProps, {}> {
  render() {
    let { title, width, children } = this.props;

    return (
      <div className="sidebar-panel" style={{ width: width || '795px' }}>
        <a className="btn-close" href="javascript:void(0)" onClick={this.onHide.bind(this)}>
          {t('关闭')}
        </a>
        <div className="sidebar-panel-container">
          <div className="sidebar-panel-hd">
            <h3 style={{ width: '240px' }}>{title}</h3>
          </div>
          <div className="sidebar-panel-bd">{this.props.children}</div>
        </div>
      </div>
    );
  }

  private onHide() {
    this.props.onClose();
  }
}
