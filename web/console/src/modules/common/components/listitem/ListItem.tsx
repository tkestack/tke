import * as classnames from 'classnames';
import * as React from 'react';

import { Bubble } from '@tea/component';
import { BaseReactProps } from '@tencent/ff-redux';

export interface ListItemProps extends BaseReactProps {
  /**显示的标题文本 */
  label?: string | JSX.Element;

  /**是否显示 */
  isShow?: boolean;

  /**提示 */
  tips?: string | JSX.Element;
}

export class ListItem extends React.Component<ListItemProps, {}> {
  render() {
    const { label, tips, isShow = true, children } = this.props;
    return isShow ? (
      <li style={{ fontSize: '12px' }}>
        <span className="item-descr-tit">
          <span style={{ verticalAlign: 'middle' }}>{label}</span>
          {tips && (
            <Bubble placement="left" content={<p style={{ whiteSpace: 'normal' }}>{tips}</p>}>
              <i className="plaint-icon" style={{ verticalAlign: 'middle' }} />
            </Bubble>
          )}
        </span>
        <span className="item-descr-txt">{children}</span>
      </li>
    ) : (
      <noscript />
    );
  }
}
