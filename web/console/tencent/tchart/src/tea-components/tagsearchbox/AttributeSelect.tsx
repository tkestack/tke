import * as React from 'react';
import * as classNames from 'classnames';

export interface AttributeValue {

  /**
   * 为资源属性需求值的类型
   */
  type: string;

  /**
   * 属性的唯一标识，会在结果中返回
   */
  key: string;

  /**
   * 资源属性值名称
   */
  name: string;

  /**
   * 属性是否可重复选择
   */
  values?: Array<any> | Function;

  /**
   * 该属性是否可重复选择
   */
  reusable?: boolean;

  /**
   * 该属性是否可移除
   */
  removeable?: boolean;
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
  "8": 'backspace',
  "9": 'tab',
  "13": 'enter',
  "37": 'left',
  "38": 'up',
  "39": 'right',
  "40": 'down'
};

export class AttributeSelect extends React.Component<AttributeSelectProps, any> {

  state: AttributeSelectState = {
    select: -1,
  }

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
  }

  getUseableList(): Array<AttributeValue> {
    const { attributes, inputValue } = this.props;

    // 获取冒号前字符串模糊查询
    const fuzzyValue = /(.*?)(:|：).*/.test(inputValue) ? RegExp.$1 : inputValue;

    return attributes.filter(item => item.name.includes(inputValue) || item.name.includes(fuzzyValue));
  }

  getAttribute(select: number): AttributeValue {
    const list = this.getUseableList();
    if (select < list.length) {
      return list[select];
    }
  }

  move = (step: number): void => {
    const select = this.state.select;
    const list = this.getUseableList();
    if (list.length <= 0) return;
    this.setState({ select: (select + step + list.length) % list.length });
  }


  handleClick = (e, index: number): void => {
    e.stopPropagation();
    if (this.props.onSelect) {
      this.props.onSelect(this.getAttribute(index));
    }
  }

  render() {
    const select = this.state.select;
    const list = this.getUseableList().map((item, index) => {
      if (select === index) {
        return (
          <li role="presentation" key={index} className="autocomplete-cur" onClick={(e) => this.handleClick(e, index)}>
            <a className="text-truncate" role="menuitem" href="javascript:;">{item.name}</a>
          </li>
        )
      }
      return (
        <li role="presentation" key={index} onClick={(e) => this.handleClick(e, index)}>
          <a className="text-truncate" role="menuitem" href="javascript:;">{item.name}</a>
        </li>
      )
    });

    if (list.length === 0) return null;

    let maxHeight = document.body.clientHeight ? document.body.clientHeight - 450 : 400;
    maxHeight = maxHeight > 240 ? maxHeight : 240;

    return (
      <div className="tc-15-autocomplete" style={{ minWidth: 180, width: 'auto' }}>
        <ul className="tc-15-autocomplete-menu" role="menu" style={{ maxHeight: `${maxHeight}px` }}>
          <li role="presentation">
            <a className="autocomplete-empty" role="menuitem" href="javascript:;">选择维度</a>
          </li>
          {list}
        </ul>
      </div>
    )
  }
}