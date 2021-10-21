/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */
import * as React from 'react';

import { BaseReactProps, WorkflowState } from '@tencent/ff-redux';

import { LinkHref, TipInfo } from '..';
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
