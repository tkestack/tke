import * as React from 'react';
import { findDOMNode } from 'react-dom';
import { BaseReactProps, OnOuterClick } from '@tencent/qcloud-lib';
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
