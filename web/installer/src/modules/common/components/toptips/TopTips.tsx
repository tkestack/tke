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
