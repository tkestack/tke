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
import { connect } from 'react-redux';

import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { ContentView, TabPanel, Tabs } from '@tencent/tea-component';

import { allActions } from '../actions';
import { router } from '../router';
import { RootProps } from './AddonApp';
import { AddonDetailHeadPanel } from './AddonDetailHeadPanel';
import { AddonDetailPanel } from './AddonDetailPanel';

/** 详情页的tab列表 */
const tabList: any[] = [
  {
    id: 'info',
    label: t('详情')
  }
];

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class AddonDetail extends React.Component<RootProps, {}> {
  componentDidMount() {
    let { actions, region } = this.props;
    // 进行地域的拉取
    region.list.fetched !== true && actions.region.applyFilter({});
  }

  render() {
    let { route } = this.props,
      urlParams = router.resolve(route);

    let { tab } = urlParams;
    // 默认选择第一个
    let selected = tabList[0];
    if (tab) {
      let finder = tabList.find(x => x.id === tab);
      if (finder) {
        selected = finder;
      }
    }

    return (
      <ContentView>
        <ContentView.Header>
          <AddonDetailHeadPanel />
        </ContentView.Header>
        <ContentView.Body>
          <Tabs
            ceiling
            tabs={tabList}
            activeId={selected ? selected.id : ''}
            onActive={tab => {
              router.navigate(Object.assign({}, urlParams, { tab: tab.id }), route.queries);
            }}
          >
            <TabPanel id="info">
              <AddonDetailPanel />
            </TabPanel>
          </Tabs>
        </ContentView.Body>
      </ContentView>
    );
  }
}
