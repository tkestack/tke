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

import { BaseReactProps } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

export interface SlidePanelProps extends BaseReactProps {
  id: string | number;
  name: string;
  onHide?: () => void;
  headStyle?: any;
}

export class SlidePanel<P extends SlidePanelProps> extends React.Component<P, {}> {
  protected logActSelector: string = '[data-logviewer]';

  private _deferTimer;
  private _handler;

  componentDidMount() {
    // 展示时，监听全局click事件，用于外部点击判断
    this._handler = this._onClick.bind(this);
    // 延时是因为有时click事件触发的面板展示，最后冒泡到body上又触发onHide
    this._deferTimer = setTimeout(() => {
      $(document.body).on('click', this._handler);
      this._deferTimer = null;
    }, 100);
  }

  componentWillUnmount() {
    // 销毁时，清除监听器
    if (this._deferTimer) {
      clearTimeout(this._deferTimer);
    }
    $(document.body).off('click', this._handler);
  }

  render() {
    let { id, name, onHide, headStyle } = this.props;

    if (!headStyle) {
      headStyle = { width: 240 };
    }

    return (
      <div className="sidebar-panel">
        {this.props.onHide && (
          <a className="btn-close" href="javascript:void(0)" onClick={onHide}>
            {t('关闭')}
          </a>
        )}
        <div className="sidebar-panel-container">
          <div className="sidebar-panel-hd">
            <h3 style={headStyle}>
              {name}
              {t('日志')}
            </h3>
            <span className="details-hd-meta">{id}</span>
          </div>
          <div className="sidebar-panel-bd">
            <div className="charts-panel">{this.props.children}</div>
          </div>
        </div>
      </div>
    );
  }

  private _onClick(e: Event) {
    let $target = $(e.target);
    // 如果点击面板外（还需确认节点位于Dom树内，以避免点击时移除的操作导致错误判断为点击外部）
    if (!$.contains(document.body, $target[0] as any)) {
      return;
    }
    if ($target.closest(this.logActSelector).length > 0) {
      return;
    }
    /* eslint-disable */
    if (this.props.onHide && !$.contains(findDOMNode(this) as any, $target[0] as any)) {
      this.props.onHide();
    }
    /* eslint-enable */
  }
}
