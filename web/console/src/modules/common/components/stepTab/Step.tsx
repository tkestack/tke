import * as React from 'react';
import { BaseReactProps } from '@tencent/qcloud-lib';

export interface StepProps extends BaseReactProps {
  /**步骤序号 */
  no?: number;

  /**显示文本 */
  text?: string;

  /**当前步骤 */
  current?: number;

  /**是否可用 */
  disabled?: boolean;
}

export class Step extends React.Component<StepProps, any> {
  render() {
    const { no, text, current } = this.props;
    return (
      <li className={current === no ? 'current' : current > no ? 'succeed' : 'disabled'}>
        <div className="tc-15-step-name">
          <span className="tc-15-step-num">{no}</span>
          {text}
        </div>
        <div className="tc-15-step-arrow" />
      </li>
    );
  }
}
