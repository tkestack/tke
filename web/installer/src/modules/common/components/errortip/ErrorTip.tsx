import * as React from 'react';
import { BaseReactProps } from '@tencent/qcloud-lib';
import { WorkflowState } from '@tencent/qcloud-redux-workflow';
import { Link } from '../../models';
import { TipInfo, LinkHref } from '..';
import { getWorkflowError, getWorkflowErrorCode } from '../../utils';

export interface ErrorGuide {
  /**链接 */
  link: Link;

  /**错误码 如果有错误码，则在指定错误码下显示指定链接；如果未指定，则在所有错误返回下显示指定链接*/
  code?: number;
}

export interface ErrorTipProps extends BaseReactProps {
  /**是否显示组件 */
  isShow?: boolean;

  /**工作流 */
  workflow?: WorkflowState<any, any>;

  /**错误指引 */
  guide?: ErrorGuide;
}

export class ErrorTip extends React.Component<ErrorTipProps> {
  render() {
    let { workflow, isShow = true, guide } = this.props,
      isShowGuide = guide && guide.code === getWorkflowErrorCode(workflow);

    return (
      <TipInfo
        isShow={isShow}
        className="error"
        style={{ display: 'inline-block', marginLeft: '20px', marginBottom: '0px' }}
      >
        {getWorkflowError(workflow)}
        {isShowGuide && (
          <LinkHref href={guide.link.href} target={guide.link.target} title={guide.link.text}>
            {guide.link.text}
          </LinkHref>
        )}
      </TipInfo>
    );
  }
}
