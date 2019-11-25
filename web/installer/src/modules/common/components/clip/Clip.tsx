import * as React from 'react';
import { BaseReactProps } from '@tencent/qcloud-lib';
import * as Clipboard from 'clipboard';
import { Tooltip, Button, Icon } from '@tencent/tea-component';
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
