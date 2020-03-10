import * as classnames from 'classnames';
import * as React from 'react';

import { BaseReactProps, fade, OnOuterClick, slide } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

export interface DropdownMenuItem {
  id: string | number;
  item: JSX.Element;
}

export interface DropdownMenuProps extends BaseReactProps {
  /**
   * 是否模拟为 select 组件，默认为 false
   * */
  simulateSelect?: boolean;

  /**
   * 外观主题：
   *   - `select` （默认）表示渲染为类似 select 组件的下拉框
   *   - `dropdown` 表示渲染为弹出菜单
   *   - `dropdown-hd` 表示渲染为标题使用的弹出菜单
   */
  theme?: 'dropdown-hd';

  /**
   * 是否使用小尺寸，默认为 false
   */
  smallSize?: boolean;

  /**
   * 弹出方向：
   *   - `down` (默认) 从下方弹出
   *   - `up` 从上方弹出
   */
  popDirection?: 'down' | 'up';

  /** 下拉按钮的样式设置 */
  buttonClassName?: string;

  /** 下拉菜单的样式设置 */
  meunClassName?: string;

  /**
   * 下拉组件的标签，simulateSelect 为 false 的时候生效
   * */
  placeholder?: string | JSX.Element;

  /**
   * 下拉显示的业务列表
   * */
  items?: DropdownMenuItem[];
}

interface DropdownMenuState {
  /**
   * 当前是否为打开状态
   */
  isOpened?: boolean;
}

export class DropdownMenu extends React.Component<DropdownMenuProps, DropdownMenuState> {
  state = {
    isOpened: false
  };

  public open() {
    this.setState({ isOpened: true });
  }

  @(OnOuterClick as any)
  public close() {
    this.setState({ isOpened: false });
  }

  public render() {
    const { isOpened } = this.state;
    const {
      items = [],
      simulateSelect = false,
      placeholder = t('更多'),
      theme,
      popDirection,
      smallSize,
      style
    } = this.props;

    const itemsBody = items.map(item => (
      <li role="presentation" key={item.id}>
        {item.item}
      </li>
    ));

    return (
      <div
        className={classnames('tc-15-dropdown m', {
          'tc-15-menu-active': isOpened,
          'tc-15-dropup': popDirection === 'up'
        })}
        style={style}
      >
        <a href="javascript:;" className="tc-15-dropdown-link" onClick={this._handleButtonClick.bind(this)}>
          <span style={{ verticalAlign: 'baseline', lineHeight: '16px', height: '16px' }}>{t('更多')}</span>
          <i className="caret" />
        </a>
        <ul key="menu" className="tc-15-dropdown-menu" role="menu">
          {itemsBody}
        </ul>
      </div>
    );
  }

  private _handleButtonClick(e: React.MouseEvent) {
    this.setState({
      isOpened: !this.state.isOpened
    });
    e.stopPropagation();
  }
}
