import * as classNames from 'classnames';
import * as React from 'react';

import { Bubble } from '@tea/component';
import { BaseReactProps, insertCSS } from '@tencent/ff-redux';

/**插入自定义样式覆盖bubble样式 */
insertCSS(
  'ButtonBarCss',
  `
.tc-15-rich-radio .tc-15-bubble-icon{
    font-size: 12px;
}
`
);

export interface ButtonBarItem {
  /* 文案 */
  name?: string | JSX.Element;

  /* 值 */
  value?: number | string;

  /* 是否可用 */
  disabled?: boolean;

  /* 提示 */
  tip?: string | JSX.Element;

  /**显示的小图标 */
  icon?: string | JSX.Element;

  /**icon 的样式 */
  iconClassName?: string;
}

export interface ButtonBarProps extends BaseReactProps {
  /* 列表数据 */
  list: ButtonBarItem[];

  /**
   * 尺寸
   *   - "m" 小尺寸
   */
  size?: string;

  /* 选中的button */
  selected?: ButtonBarItem;

  /* 选择后的回调 */
  onSelect?: (value: ButtonBarItem) => void;

  bubbleDirection?: 'top' | 'right' | 'left' | 'bottom';

  /* 是否为国际版用户 */
  isI18n?: boolean;

  /** 单个list的时候是否显示纯文本 */
  isNeedPureText?: boolean;

  style?: object;
  buttonStyle?: object;
}

export class ButtonBar extends React.Component<ButtonBarProps, {}> {
  public select(item: ButtonBarItem) {
    const { onSelect } = this.props;

    typeof onSelect === 'function' && onSelect(item);
  }

  renderButton(list: ButtonBarItem[]) {
    const { selected, size, isI18n, bubbleDirection = 'top', buttonStyle } = this.props;
    const length = list.length;

    return list.map((item, index) => {
      // 判断是否为国际版，隐藏包年包月的相关内容
      if (isI18n && (item.value === 'monthly' || item.value === 'PayByHour')) {
        /* eslint-disable */
        return;
        /* eslint-enable */
      }

      let classname = classNames(
        'tc-15-btn',
        { checked: selected && item.value === selected.value },
        { first: index === 0 },
        { last: index === length - 1 },
        { disabled: item.disabled },
        size
      );

      let bubbleContent: string | React.ReactNode = '';

      if (item.tip) {
        bubbleContent = <p dangerouslySetInnerHTML={{ __html: item.tip as string }} />;
      }
      if (item.tip && typeof item.tip !== 'string') {
        bubbleContent = item.tip;
      }

      return (
        <Bubble key={index} placement={bubbleDirection} content={bubbleContent || null}>
          {item.iconClassName ? (
            <a
              href="javascript:;"
              style={Object.assign({ paddingRight: '0', marginBottom: '0' }, buttonStyle)}
              className={classname}
              onClick={item.disabled ? null : () => this.select(item)}
              key={index}
            >
              {item.name}
              {item.icon && (
                <i
                  style={item.iconClassName ? { position: 'relative', right: '-1px', marginLeft: '8px' } : {}}
                  className={item.iconClassName ? item.iconClassName : 'shop-ui-block-icon'}
                >
                  <em>{item.icon}</em>
                </i>
              )}
            </a>
          ) : (
            <a
              href="javascript:;"
              style={Object.assign({ marginBottom: '0' }, buttonStyle)}
              className={classname}
              onClick={item.disabled ? null : () => this.select(item)}
              key={index}
            >
              {item.name}
              {item.icon && (
                <i className="shop-ui-block-icon">
                  <em>{item.icon}</em>
                </i>
              )}
            </a>
          )}
        </Bubble>
      );
    });
  }

  render() {
    const { list, className, isNeedPureText, style } = this.props;

    return isNeedPureText === false ? (
      <div className={classNames('tc-15-rich-radio', className)} style={Object.assign({ overflow: 'visible' }, style)}>
        {this.renderButton(list)}
      </div>
    ) : list.length === 1 ? (
      <span style={{ fontSize: '12px', lineHeight: '2' }}>{list[0].name}</span>
    ) : (
      <div className={classNames('tc-15-rich-radio', className)} style={Object.assign({ overflow: 'visible' }, style)}>
        {this.renderButton(list)}
      </div>
    );
  }
}
