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
import * as ReactDom from 'react-dom';

import { FadeTransition } from '@tencent/tea-component';

export interface TopTipsProps {
  /** 提示消息 */
  message?: string;

  /** 类型 success/error */
  theme?: string;

  /** 间隔时间 */
  duration?: number;
}

let instant;
let timeout;

export function TopTips(options: TopTipsProps) {
  clearTimeout(timeout);
  let { message, theme, duration = 1500 } = options;
  let classname = classnames({
    'top-alert-icon-done': theme === 'success',
    'top-alert-icon-waring': theme === 'error'
  });

  if (!instant) {
    instant = document.createElement('div');
    document.body.appendChild(instant);
  }

  ReactDom.render(
    <FadeTransition in={true}>
      <div className="top-alert" style={{ marginLeft: '-200px', zIndex: 1100 }}>
        <span className={classname}>{message}</span>
      </div>
    </FadeTransition>,
    instant
  );

  timeout = setTimeout(function() {
    ReactDom.unmountComponentAtNode(instant);
  }, duration);
}
