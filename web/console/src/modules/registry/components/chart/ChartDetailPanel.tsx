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

import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { ContentView, Icon, Justify, TabPanel, Tabs, Text } from '@tencent/tea-component';

import { router } from '../../router';
import { RootProps } from '../RegistryApp';
import { ChartTablePanel } from './ChartTablePanel';

export class ChartDetailPanel extends React.Component<RootProps, any> {
  render() {
    return (
      <ContentView>
        <ContentView.Header>
          <Justify
            left={
              <React.Fragment>
                <a
                  href="javascript:;"
                  className="back-link"
                  onClick={() => {
                    this.goBack();
                  }}
                >
                  <Icon type="btnback" />
                  {t('返回')}
                </a>
                <h2>
                  {t('Chart包列表')} <span>({this.props.route.queries['cgName'] || '-'})</span>
                </h2>
              </React.Fragment>
            }
          />
          ;
        </ContentView.Header>
        <ContentView.Body>
          <ChartTablePanel {...this.props} />
        </ContentView.Body>
      </ContentView>
    );
  }

  private goBack() {
    let urlParams = router.resolve(this.props.route);
    router.navigate(Object.assign({}, urlParams, { sub: 'chart', mode: 'list' }), this.props.route.queries);
  }
}
