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

import * as classnames from 'classnames';
import * as React from 'react';
import * as TransitionGroup from 'react-addons-css-transition-group';

import { BaseReactProps, fade, OnOuterClick, slide } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

export interface DropdownListItem {
  /**
   * 列表项的标识，在同一个列表中不允许重复
   */
  id: string | number;

  /**
   * 列表显示的标签，可以为字符串或者是特定 React Element
   */
  label?: string | JSX.Element;

  /**
   * 列表是否被禁用，默认为 false
   */
  disabled?: boolean;

  /**
   * 子dropdown（支持分类下拉菜单,暂时只支持到二级下拉菜单）
   */
  sub?: DropdownListItem[];
}

export interface DropdownListProps extends BaseReactProps {
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
  theme?: 'select' | 'dropdown' | 'dropdown-hd';

  /**
   * 是否使用小尺寸，默认为 false
   */
  smallSize?: boolean;

  /**
   * 按钮最大宽度，可以给数字（单位：px）或者 CSS 字符串值，默认为空，不限制最大宽度
   */
  buttonMaxWidth?: number | string;

  /**
   * 菜单最大宽度，可以给数字（单位：px）或者 CSS 字符串值，默认为空，不限制最大宽度
   */
  menuMaxWidth?: number | string;

  /**
   * 菜单最大高度，可以给数字（单位：px）或者 CSS 字符串值，默认为空，不限制最大高度
   */
  menuMaxHeight?: number | string;

  /**
   * 弹出方向：
   *   - `down` (默认) 从下方弹出
   *   - `up` 从上方弹出
   */
  popDirection?: 'down' | 'up';

  /** 下拉按钮的样式设置 */
  buttonClassName?: 'transparent' | string;

  /** 下拉菜单的样式设置 */
  meunClassName?: string;

  /**
   * 下拉组件的标签，simulateSelect 为 false 的时候生效
   * */
  placeholder?: string | JSX.Element;

  /**
   * 下拉显示的业务列表
   * */
  items?: DropdownListItem[];

  /**
   * 当前选中的业务
   */
  selected?: DropdownListItem;

  /**
   * 下拉项被选中的时候
   */
  onSelect?: (item: DropdownListItem) => void;

  /**
   * 是否禁用下拉组件
   */
  disabled?: boolean;
}

interface StyleClassConfig {
  className: string;
  activeClassName?: string;
  buttonClassName: string;
  buttonActiveClassName?: string;
  menuClassName: string;
  menuItemSelectedClassName: string;
  triggerMenuDisplay?: boolean;
  upClassName?: string;
}

type StyleClassConfigMap = {
  [theme: string]: StyleClassConfig;
};

const styleClassConfigMap: StyleClassConfigMap = {
  select: {
    className: 'tc-15-simulate-select-wrap',
    activeClassName: 'show',
    buttonClassName: 'tc-15-simulate-select',
    buttonActiveClassName: 'show',
    menuClassName: 'tc-15-simulate-option',
    menuItemSelectedClassName: 'selected',
    upClassName: '',
    triggerMenuDisplay: true
  },
  dropdown: {
    className: 'tc-15-dropdown',
    buttonClassName: 'tc-15-dropdown-link',
    activeClassName: 'tc-15-menu-active',
    buttonActiveClassName: 'active',
    menuClassName: 'tc-15-dropdown-menu',
    menuItemSelectedClassName: 'selected',
    upClassName: 'tc-15-dropup'
  },
  'dropdown-hd': {
    className: 'tc-15-dropdown tc-15-dropdown-in-hd',
    buttonClassName: 'tc-15-dropdown-link',
    activeClassName: 'tc-15-menu-active',
    buttonActiveClassName: 'active',
    menuClassName: 'tc-15-dropdown-menu',
    menuItemSelectedClassName: 'selected',
    upClassName: 'tc-15-dropup'
  }
};

interface DropdownListState {
  /**
   * 当前是否为打开状态
   */
  isOpened?: boolean;
}

export class DropdownList extends React.Component<DropdownListProps, DropdownListState> {
  _documentClickHandler: (e: MouseEvent) => void;

  state = {
    isOpened: false
  };

  totalSubItems = [];

  _handleButtonClick(e: React.MouseEvent) {
    if (!this.props.disabled) {
      this.setState({
        isOpened: !this.state.isOpened
      });
      e.stopPropagation();
    }
  }

  _handleItemClick(item: DropdownListItem, e: React.MouseEvent) {
    if (!item.disabled) {
      this.select(item);
    }
    this.close();
    e.stopPropagation();
  }

  public getTotalSubItems(items) {
    items.forEach(item => {
      if (item.sub) {
        this.getTotalSubItems(item.sub);
      } else {
        this.totalSubItems.push(item);
      }
    });
  }

  public select(item: DropdownListItem) {
    if (item.sub) {
      //分类菜单的标题不支持点击选中
      return false;
    }

    this.getTotalSubItems(this.props.items);
    let selectedIndex = this.totalSubItems.indexOf(item);

    if (this.props.onSelect) {
      this.props.onSelect(selectedIndex > -1 ? item : null);
    }
  }

  public open() {
    if (!this.props.disabled) {
      this.setState({ isOpened: true });
    }
  }

  @OnOuterClick
  public close() {
    this.setState({ isOpened: false });
  }

  public render() {
    const { selected, buttonMaxWidth, menuMaxWidth, menuMaxHeight, disabled } = this.props;
    const { isOpened } = this.state;
    const {
      items = [],
      simulateSelect = false,
      placeholder = t('请选择'),
      theme,
      popDirection,
      smallSize
    } = this.props;

    const classConfig = styleClassConfigMap[theme] || styleClassConfigMap['select'];

    const className = classnames(classConfig.className, this.props.className, {
      [classConfig.activeClassName]: isOpened,
      [classConfig.upClassName]: popDirection === 'up',
      m: smallSize,
      disabled
    });
    const buttonClassName = classnames(classConfig.buttonClassName, this.props.buttonClassName, {
      [classConfig.buttonActiveClassName]: isOpened,
      m: smallSize
    });
    const menuClassName = classnames(classConfig.menuClassName, this.props.meunClassName);

    const buttonContent = (simulateSelect && selected && selected.label) || placeholder;

    const getItemsBody = items => {
      let itemsBody = items.map(item => {
        return (
          <li
            role="presentation"
            className={classnames({
              disabled: item.disabled,
              [classConfig.menuItemSelectedClassName]: item === selected,
              'tc-15-optgroup': !!item.sub
            })}
            key={item.key}
            onClick={e => this._handleItemClick(item, e)}
            title={typeof item.label === 'string' && !item.sub ? item.label : null}
          >
            {//分类dropdown的标题
            !!item.sub && (
              <h6 role="menuitem" className="tc-15-optgroup-label" style={{ fontSize: 12 }}>
                {item.label}
              </h6>
            )}
            {//递归生成分类dropdown的内容
            !!item.sub && (
              <ul key="menu" className={menuClassName} style={menuStyle}>
                {getItemsBody(item.sub)}
              </ul>
            )}
            {!item.sub && (
              <a
                role="menuitem"
                style={{ cursor: item.disabled ? 'not-allowed' : 'pointer' }}
                className={classnames({ [classConfig.menuItemSelectedClassName]: item === selected })}
              >
                {item.label}
              </a>
            )}
          </li>
        );
      });

      return itemsBody;
    };

    let wrapperStyle: React.CSSProperties = Object.assign(
      {
        fontSize: 12
      },
      this.props.style || {}
    );

    let buttonStyle: React.CSSProperties = {
      display: 'inline-block',
      cursor: disabled ? 'not-allowed' : 'pointer',
      color: disabled ? '#999' : null,
      backgroundColor: buttonClassName === 'transparent' ? 'transparent' : 'white'
    };
    if (buttonMaxWidth) {
      buttonStyle = Object.assign(buttonStyle, {
        maxWidth: buttonMaxWidth,
        overflow: 'hidden',
        textOverflow: 'ellipsis',
        whiteSpace: 'nowrap'
      });
    }

    let menuStyle: React.CSSProperties = { maxWidth: menuMaxWidth, maxHeight: menuMaxHeight, overflow: 'auto' };
    if (classConfig.triggerMenuDisplay) {
      menuStyle.display = isOpened ? 'block' : 'none';
    }

    return (
      <div className={className} style={wrapperStyle}>
        <a
          style={buttonStyle}
          className={buttonClassName}
          onClick={this._handleButtonClick.bind(this)}
          title={typeof buttonContent === 'string' ? buttonContent : null}
        >
          {buttonContent}
          <i className="caret" />
        </a>
        {Array.isArray(items) && items.length > 0 && (
          <ul key="menu" className={menuClassName} style={menuStyle}>
            {getItemsBody(items)}
          </ul>
        )}
      </div>
    );
  }
}
