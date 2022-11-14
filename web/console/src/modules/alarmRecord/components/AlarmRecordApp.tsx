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
import { connect, Provider } from 'react-redux';
import { t } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators } from '@tencent/ff-redux';
import { Layout, Card } from '@tea/component';
import { ResetStoreAction } from '../../../../helpers';
import { allActions } from '../actions';
import { RootState } from '../models';
import { router } from '../router';
import { configStore } from '../stores/RootStore';
import { AlarmRecordHeadPanel } from './AlarmHeaderPanel';
import { AlarmTablePanel } from './AlarmTablePanel';

const { useState, useEffect } = React;
const { Body, Content } = Layout;
const store = configStore();

export class AlarmRecordContainer extends React.Component<any, any> {
  // 页面离开时，清空store
  componentWillUnmount() {
    store.dispatch({ type: ResetStoreAction });
  }
  render() {
    return (
      <Provider store={store}>
        <AlarmRecordApp />
      </Provider>
    );
  }
}

export interface RootProps extends RootState {
  actions?: typeof allActions;
}
const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
@((router.serve as any)())
class AlarmRecordApp extends React.Component<RootProps, {}> {
  render() {
    return (
      <Layout>
        <Body>
          <Content>
            <Content.Header title={t('历史告警记录')}>
              <AlarmRecordHeadPanel />
            </Content.Header>
            <Content.Body full>
              <AlarmTablePanel clusterId={this?.props?.cluster?.selection?.metadata?.name} />
            </Content.Body>
          </Content>
        </Body>
      </Layout>
    );
  }
}
