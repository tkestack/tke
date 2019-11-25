import * as React from 'react';
import * as classNames from 'classnames';

import { t, Trans } from '@tencent/tea-app/lib/i18n';

export interface AttributeValue {
  type: string;
  key: string;
  name: string;
  values?: Array<any> | Function;
  reusable?: boolean;
}

export interface AttributeSelectProps {
  attributes: Array<AttributeValue>;
  inputValue: string;
  onSelect?: (attribute: AttributeValue) => void;
}

export interface AttributeSelectState {
  select: number;
}

const keys = {
  '8': 'backspace',
  '9': 'tab',
  '13': 'enter',
  '37': 'left',
  '38': 'up',
  '39': 'right',
  '40': 'down'
};

export class AttributeSelect extends React.Component<AttributeSelectProps, any> {
  state: AttributeSelectState = {
    select: -1
  };

  componentWillReceiveProps(nextProps: AttributeSelectProps) {
    if (this.props.inputValue !== nextProps.inputValue) {
      this.setState({ select: -1 });
    }
  }

  // 父组件调用
  handleKeyDown = (keyCode: number): boolean => {
    if (!keys[keyCode]) return;
    const { onSelect } = this.props;
    const select = this.state.select;

    switch (keys[keyCode]) {
      case 'enter':
      case 'tab':
        if (select < 0) break;
        if (onSelect) {
          onSelect(this.getAttribute(select));
        }
        return false;

      case 'up':
        this.move(-1);
        break;

      case 'down':
        this.move(1);
        break;
    }
  };

  getAttribute(select: number): AttributeValue {
    const { attributes, inputValue } = this.props;
    const list = attributes.filter(item => item.name.indexOf(inputValue) >= 0);
    if (select < list.length) {
      return list[select];
    }
  }

  move = (step: number): void => {
    const select = this.state.select;
    const { attributes, inputValue } = this.props;
    const list = attributes.filter(item => item.name.indexOf(inputValue) >= 0);
    if (list.length <= 0) return;
    this.setState({ select: (select + step + list.length) % list.length });
  };

  handleClick = (e, index: number): void => {
    e.stopPropagation();
    if (this.props.onSelect) {
      this.props.onSelect(this.getAttribute(index));
    }
  };

  render() {
    const select = this.state.select;
    const { inputValue, attributes } = this.props;

    const list = attributes
      .filter(item => item.name.indexOf(inputValue) >= 0)
      .map((item, index) => {
        if (select === index) {
          return (
            <li role="presentation" key={index} className="autocomplete-cur" onClick={e => this.handleClick(e, index)}>
              <a className="text-truncate" role="menuitem" href="javascript:;">
                {item.name}
              </a>
            </li>
          );
        }
        return (
          <li role="presentation" key={index} onClick={e => this.handleClick(e, index)}>
            <a className="text-truncate" role="menuitem" href="javascript:;">
              {item.name}
            </a>
          </li>
        );
      });

    if (list.length === 0) return null;

    return (
      <div className="tc-15-autocomplete">
        <ul className="tc-15-autocomplete-menu" role="menu">
          <li role="presentation">
            <a className="autocomplete-empty" role="menuitem" href="javascript:;">
              {t('选择资源属性进行过滤')}
            </a>
          </li>
          {list}
        </ul>
      </div>
    );
  }
}
