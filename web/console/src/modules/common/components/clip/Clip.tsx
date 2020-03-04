import * as Clipboard from 'clipboard';
import * as React from 'react';

import { Bubble } from '@tea/component';
import { BaseReactProps } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

const tips = seajs.require('tips');
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

export class Clip extends React.Component<ClipProps, {}> {
  render() {
    const { target, isShow = true, isShowTip, className, tipDirection, style, children } = this.props;
    let renderClass = children ? 'copy-trigger ' + className : 'copy-trigger copy-icon ' + className;
    return isShow ? (
      <Bubble
        // className={className}
        placement={tipDirection || 'bottom'}
        content={isShowTip ? t('复制') : null}
      >
        <a
          href="javascript:;"
          data-clipboard-action="copy"
          data-clipboard-target={target}
          className={renderClass}
          style={style}
        >
          {children}
        </a>
      </Bubble>
    ) : (
      <noscript />
    );
  }

  componentDidMount() {
    if (!window['paas_is_init_clipboard']) {
      window['paas_is_init_clipboard'] = new Clipboard('.copy-trigger');
      window['paas_is_init_clipboard'].on('success', e => {
        tips.success(t('复制成功'), 1000);
        e.clearSelection();
      });
      window['paas_is_init_clipboard'].on('error', e => {
        tips.error(t('复制失败'), 1000);
      });
    }
  }
}
