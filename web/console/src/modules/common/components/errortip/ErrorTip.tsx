import * as React from 'react';

import { BaseReactProps, WorkflowState } from '@tencent/ff-redux';
import { ExternalLink } from '@tencent/tea-component';

import { TipInfo } from '../';
import { Link } from '../../models';
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

export class ErrorTip extends React.Component<ErrorTipProps, {}> {
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
        {isShowGuide && <ExternalLink href={guide.link.href}>{guide.link.text}</ExternalLink>}
      </TipInfo>
    );
  }
}
