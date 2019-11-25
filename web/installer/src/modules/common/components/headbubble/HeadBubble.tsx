import * as React from 'react';
import { BaseReactProps } from '@tencent/qcloud-lib';
import { Bubble, Icon } from '@tencent/tea-component';

export interface HeadBubbleProps extends BaseReactProps {
  /**显示标题 */
  title?: string | JSX.Element;

  /**显示的文本 */
  text?: string | JSX.Element;

  /**气泡显示方式 */
  position?: 'top' | 'bottom' | 'left' | 'right';

  /**对齐方式 */
  align?: 'start' | 'end';

  /** 用于title隐藏 */
  autoflow?: boolean;
}

export class HeadBubble extends React.Component<HeadBubbleProps> {
  render() {
    const { title = '', text = '', position, align, autoflow } = this.props;
    return (
      <div>
        {autoflow ? <span className="text-overflow">{title}</span> : <span>{title}</span>}
        <Bubble placement={position} content={text}>
          <Icon type="info" />
        </Bubble>
      </div>
    );
  }
}
