import * as classnames from 'classnames';
import * as React from 'react';

import { Bubble } from '@tea/component';
import { BaseReactProps } from '@tencent/ff-redux';

export interface FormItemProps extends BaseReactProps {
  /**显示的文本 */
  label?: string | JSX.Element;

  /**是否纯文本显示 */
  isPureText?: boolean;

  /**提示 */
  tips?: string | JSX.Element;

  /**样式 */
  className?: any;

  /**是否显示 */
  isShow?: boolean;

  minWidth?: number;

  isNeedFormInput?: boolean;
}

export class FormItem extends React.Component<FormItemProps, {}> {
  render() {
    const {
      isShow = true,
      label = '',
      isPureText = false,
      tips = '',
      children,
      className,
      style,
      isNeedFormInput = true
    } = this.props;
    return isShow ? (
      <li className={classnames(className, { 'pure-text-row': isPureText })} style={style}>
        <div className="form-label" style={{ minWidth: this.props.minWidth || '80px', verticalAlign: 'top' }}>
          <label>
            {label}
            {tips ? (
              <Bubble placement="top" content={tips || null}>
                <i className="plaint-icon" style={{ marginLeft: '5px' }} />
              </Bubble>
            ) : (
              <noscript />
            )}
          </label>
        </div>
        {isNeedFormInput ? (
          <div className="form-input">{children}</div>
        ) : (
          <div style={{ paddingBottom: 16 }}>{children}</div>
        )}
      </li>
    ) : (
      <noscript />
    );
  }
}
