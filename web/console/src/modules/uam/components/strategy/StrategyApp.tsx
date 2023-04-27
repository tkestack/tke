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
import { useDispatch, useSelector } from 'react-redux';

import { ContentView, Layout, Tabs, TabPanel } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { allActions } from '../../actions';
import { RootState } from '../../models';
import { router } from '../../router';
import { StrategyActionPanel } from './StrategyActionPanel';
import { StrategyDetailsPanel } from './StrategyDetailsPanel';
import { StrategyHeadPanel } from './StrategyHeadPanel';
import { StrategyTablePanel } from './StrategyTablePanel';

import { UserPanel } from '@src/modules/uam/components/user/UserPanel';
import { GroupPanel } from '@src/modules/uam/components/group/GroupPanel';
import { checkCustomVisible } from '@src/modules/common/components/permission-provider';
const { Body, Content } = Layout;
const { useState, useEffect } = React;

export interface RootProps extends RootState {
  actions?: typeof allActions;
}
//
// const mapDispatchToProps = dispatch =>
//   Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

// @connect(state => state, mapDispatchToProps)
export const StrategyApp = props => {
  const state = useSelector(state => state);
  const dispatch = useDispatch();
  const { actions } = bindActionCreators({ actions: allActions }, dispatch);
  const { route } = state;
  // useEffect(() => {
  //   actions.strategy.poll();
  // }, []);

  const urlParam = router.resolve(route);
  const { module, sub, action } = urlParam;

  const tabs = [
    { id: 'platform', label: '平台策略' },
    ...(checkCustomVisible('platform.uam.business') ? [{ id: 'business', label: '业务策略' }] : [])
  ];

  let header;
  if (action === 'detail') {
    // 策略详情
    header = <Content.Header showBackButton onBackButtonClick={() => history.back()} title={route.queries['id']} />;
  } else {
    // 策略管理
    header = <Content.Header showBackButton onBackButtonClick={() => history.back()} title={t('策略管理')} />;
  }

  return (
    <Layout>
      <Body>
        <Content>
          {header}
          <Content.Body>
            <Tabs
              ceiling
              animated={false}
              tabs={tabs}
              activeId={sub || 'platform'}
              onActive={value => {
                router.navigate({ module: 'strategy', sub: value.id });
              }}
            >
              <TabPanel id="platform">
                {action === 'detail' ? (
                  <StrategyDetailsPanel />
                ) : (
                  <>
                    <StrategyActionPanel type="platform" />
                    <StrategyTablePanel type="platform" />
                  </>
                )}
              </TabPanel>
              <TabPanel id="business">
                {action === 'detail' ? (
                  <StrategyDetailsPanel />
                ) : (
                  <>
                    <StrategyActionPanel type="business" />
                    <StrategyTablePanel type="business" />
                  </>
                )}
              </TabPanel>
            </Tabs>
          </Content.Body>
        </Content>
      </Body>
    </Layout>
  );
  //
  // return (
  //   <React.Fragment>
  //     {sub ? (
  //       <ContentView>
  //         <ContentView.Header>
  //           <StrategyHeadPanel />
  //         </ContentView.Header>
  //         <ContentView.Body>
  //           <StrategyDetailsPanel />
  //         </ContentView.Body>
  //       </ContentView>
  //     ) : (
  //       <ContentView>
  //         <ContentView.Header>
  //           <StrategyHeadPanel />
  //         </ContentView.Header>
  //         <ContentView.Body>
  //           <StrategyActionPanel />
  //           <StrategyTablePanel />
  //         </ContentView.Body>
  //       </ContentView>
  //     )}
  //   </React.Fragment>
  // );
};
