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

import { bindActionCreators } from '@tencent/ff-redux';
import { Col, ContentView, Row } from '@tencent/tea-component';

import { ResetStoreAction } from '../../../../helpers';
import { overviewActions } from '../actions/overviewActions';
import { RootState } from '../models/RootState';
import { router } from '../router';
import { configStore } from '../stores/RootStore';
import { ClusterDetailPanel } from './ClusterDetailPanel';
import { ClusterOverviewPanel } from './ClusterOverview';
import { OverviewHeadPanel } from './OverviewHeadPanel';
import { QuickHelpPanel } from './QuickHelpPanel';
import { TipsPanel } from './TipsPanel';
import { PermissionProvider } from '@common';

const store = configStore();

export class OverviewAppContainer extends React.Component<any, any> {
  // 页面离开时，清空store
  componentWillUnmount() {
    store.dispatch({ type: ResetStoreAction });
  }

  render() {
    return (
      <Provider store={store}>
        <OverviewApp />
      </Provider>
    );
  }
}

export interface RootProps extends RootState {
  actions?: typeof overviewActions;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: overviewActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
@((router.serve as any)())
class OverviewApp extends React.Component<RootProps, {}> {
  render() {
    const { clusterOverview } = this.props;
    return (
      <ContentView>
        <ContentView.Header>
          <OverviewHeadPanel {...this.props} />
        </ContentView.Header>
        <ContentView.Body>
          <Row>
            <Col span={18}>
              <ClusterOverviewPanel
                clusterData={clusterOverview.object && clusterOverview.object.data ? clusterOverview.object.data : null}
              />
              <ClusterDetailPanel
                clusterData={clusterOverview.object && clusterOverview.object.data ? clusterOverview.object.data : null}
              />
            </Col>
            <PermissionProvider value="platform.overview.help">
              <Col span={6}>
                <QuickHelpPanel />
                <TipsPanel />
              </Col>
            </PermissionProvider>
          </Row>
        </ContentView.Body>
      </ContentView>
    );
  }
}
