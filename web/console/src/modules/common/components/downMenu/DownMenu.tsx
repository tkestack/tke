import * as React from 'react';

import { BaseReactProps, insertCSS } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { DropdownList, DropdownListItem } from '../dropdown';

/**插入自定义样式覆盖bubble样式 */
insertCSS(
  'ButtonBarCss',
  `
.tc-15-rich-radio .tc-15-bubble-icon{
    font-size: 12px;
}
`
);
//regionbar样式调整
insertCSS(
  'DownMenu',
  `
.dropdown-hd-min-width{
    min-width: 60px;
}
`
);
insertCSS(
  'DownMenuCss',
  `
.tc-15-dropdown-allow-hover:hover .tc-15-dropdown-menu, .tc-15-menu-active .tc-15-dropdown-menu-min-width{
    min-width: 60px;
}
`
);

insertCSS(
  'DownMenuACss',
  `
.tc-15-dropdown-in-hd .tc-15-dropdown-menu li a{
    min-width: 60px;
}
`
);
export interface DownMenuItem {
  /* 文案 */
  name?: string | JSX.Element;

  /* 值 */
  value?: number | string;

  /* 是否可用 */
  disabled?: boolean;

  /* 提示 */
  tip?: string | JSX.Element;

  /**地域 */
  area?: string;
}

export interface DownMenuStateProps extends BaseReactProps {
  /* 列表数据 */
  list: DownMenuItem[];

  /* 选中的button */
  selected?: DownMenuItem;

  /* 选择后的回调 */
  onSelect?: (value: any) => void;
}

interface DownMenuState {
  /**
   * 当前是否为打开状态
   */
  isOpened?: boolean;
}
export class DownMenu extends React.Component<DownMenuStateProps, DownMenuState> {
  state = {
    isOpened: false
  };
  public select(item: DownMenuItem) {
    const { onSelect } = this.props;

    typeof onSelect === 'function' && onSelect(item);
  }

  renderButton(list: DownMenuItem[]) {
    const { selected } = this.props;
    let selectItem: DropdownListItem;
    let options = [];
    list.forEach(item => {
      let option: DropdownListItem = {
        id: item.value,
        label: item.area ? `${item.area}(${item.name})` : `${item.name}`
      };
      if (selected && selected.value === item.value) {
        selectItem = option;
      }
      options.push(option);
    });
    if (list.length === 0) {
      let option: DropdownListItem = {
        id: 'empty',
        label: t('无')
      };
      options.push(option);
      selectItem = option;
    }
    return { selectItem, options };
  }

  // @OnOuterClick
  // public close() {
  //     this.setState({ isOpened: false });
  // }

  // private _handleButtonClick(e: React.MouseEvent) {
  //     this.setState({
  //         isOpened: !this.state.isOpened
  //     });
  //     e.stopPropagation();
  // }
  // private _handleRegionClick(item) {
  //     this.select(item)
  //     this.setState({
  //         isOpened: !this.state.isOpened
  //     });
  // }
  render() {
    const { list, onSelect, style } = this.props;
    let rendList = this.renderButton(list);
    return (
      <div style={{ display: 'inline-block', border: '1px solid #ddd', height: '30px' }} className="form-unit">
        <DropdownList
          items={rendList.options}
          simulateSelect
          selected={rendList.selectItem}
          onSelect={item => {
            onSelect(item);
          }}
          buttonMaxWidth="200px"
          menuMaxWidth="200px"
          menuMaxHeight="400px"
          theme="dropdown-hd"
          className="dropdown-hd-min-width"
          meunClassName="tc-15-dropdown-menu-min-width"
          buttonClassName="dropdown-head transparent"
        />
      </div>
    );
  }
}
