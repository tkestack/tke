import * as React from 'react';
import { ButtonBar, ButtonBarItem, ButtonBarProps } from '../buttonbar';

/**
 * 和普通Buttonbar组件相比，区别在于选中的为值，而不是item
 *
 * @export
 * @interface CommonBarItem
 */
export interface CommonBarItem extends ButtonBarItem {}

export interface CommonBarProps extends ButtonBarProps {
  /* 选中的button */
  value?: string | number;

  /* 列表数据 */
  list: CommonBarItem[];

  /* 选择后的回调 */
  onSelect?: (value: ButtonBarItem) => void;

  /* 判断是否为国际版 */
  isI18n?: boolean;

  isNeedPureText?: boolean;

  style?: object;
  buttonStyle?: object;
}

export class CommonBar extends React.Component<CommonBarProps, {}> {
  render() {
    let { list, value, onSelect, isI18n, isNeedPureText, style, buttonStyle } = this.props,
      barList: ButtonBarItem[] = [],
      selected: ButtonBarItem;
    list.forEach(item => {
      let buttonItem: ButtonBarItem = {
        name: item.name,
        value: item.value,
        tip: item.tip
      };

      if (value && value === item.value) {
        selected = buttonItem;
      }

      barList.push(buttonItem);
    });

    return list.length === 1 && isNeedPureText ? (
      <span>{list[0].name}</span>
    ) : (
      <ButtonBar
        style={style}
        buttonStyle={buttonStyle}
        list={list}
        selected={selected}
        size="m"
        onSelect={onSelect}
        isI18n={isI18n}
        isNeedPureText={isNeedPureText}
      />
    );
  }
}
