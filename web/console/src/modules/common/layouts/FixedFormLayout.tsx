/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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

interface FixedFormLayoutProps extends BaseReactProps {
  /** 是否去掉 ul 的 margin-top */
  isRemoveUlMarginTop?: boolean;

  /** style */
  style?: any;
}

export class FixedFormLayout extends React.Component<FixedFormLayoutProps, {}> {
  render() {
    let { isRemoveUlMarginTop, style } = this.props;

    return (
      <div className="run-docker-box" style={style}>
        <div className="edit-param-list">
          <div className="param-box">
            <div className="param-bd">
              <ul className="form-list fixed-layout" style={isRemoveUlMarginTop ? { marginTop: '0' } : {}}>
                {this.props.children}
              </ul>
            </div>
          </div>
        </div>
      </div>
    );
  }
}
