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

import { RootProps } from './InstallerApp';

export class InstallerHeadPanel extends React.Component<RootProps, void> {
  render() {
    return (
      <div className="qc-header-nav" id="topnav" data-reactid=".1.0">
        <div className="qc-header-inner" data-reactid=".1.0.0">
          <div
            className="qc-header-unit qc-header-logo"
            data-reactid=".1.0.0.0"
          >
            <div className="qc-nav-logo" data-reactid=".1.0.0.0.0">
              <a
                className="qc-logo-inner"
                href="javascript:;"
                title="容器服务"
                data-reactid=".1.0.0.0.0.0"
              >
                <span className="qc-logo-icon" data-reactid=".1.0.0.0.0.0.0">
                  腾讯云
                </span>
              </a>
            </div>
          </div>
        </div>
      </div>
    );
  }
}
