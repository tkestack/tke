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

import { BaseReactProps } from '@tencent/ff-redux';

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
