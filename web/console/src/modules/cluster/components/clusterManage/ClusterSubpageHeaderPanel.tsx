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
import { RootProps } from '../ClusterApp';
import { router } from '../../router';
import { Justify, Icon } from '@tencent/tea-component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
export class ClusterSubpageHeaderPanel extends React.Component<RootProps, {}> {
  goBack() {
    history.back();
  }

  render() {
    let title = t('导入集群');

    return (
      <Justify
        left={
          <React.Fragment>
            <a href="javascript:;" className="back-link" onClick={this.goBack.bind(this)}>
              <Icon type="btnback" />
              {t('返回')}
            </a>
            <span className="line-icon">|</span>
            <h2>{title}</h2>
          </React.Fragment>
        }
      />
    );
  }
}
