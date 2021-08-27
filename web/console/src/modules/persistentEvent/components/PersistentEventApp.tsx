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
import { ContentView } from '@tencent/tea-component';

import { resourceConfig } from '../../../../config';
import { ResetStoreAction } from '../../../../helpers';
import { allActions } from '../actions';
import { RootState } from '../models';
import { router } from '../router';
import { configStore } from '../stores/RootStore';
import { ClusterActionPanel } from './ClusterActionPanel';
import { ClusterTablePanel } from './ClusterTablePanel';
import { EditPersistentEventPanel } from './EditPersistentEventPanel';
import { PersistentEventHeadPanel } from './PersistentEventHeadPanel';

const store = configStore();
export class PersistentEventAppContainer extends React.Component<any, any> {
  // 页面离开时，清空store
  componentWillUnmount() {
    store.dispatch({ type: ResetStoreAction });
  }

  render() {
    return (
      <Provider store={store}>
        <PersistentEventApp />
      </Provider>
    );
  }
}

export interface RootProps extends RootState {
  actions?: typeof allActions;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
@((router.serve as any)())
class PersistentEventApp extends React.Component<RootProps, {}> {
  componentDidMount() {
    let { actions } = this.props;

    // 初始化 resourceInfo
    let peResourceInfo = resourceConfig()['pe'];
    actions.pe.initPeResourceInfo(peResourceInfo);
  }

  render() {
    let { route } = this.props,
      urlParams = router.resolve(route);

    if (!urlParams['mode']) {
      return (
        <ContentView>
          <ContentView.Header>
            <PersistentEventHeadPanel />
          </ContentView.Header>
          <ContentView.Body>
            <ClusterActionPanel />
            <ClusterTablePanel />
          </ContentView.Body>
        </ContentView>
      );
    } else if (urlParams['mode'] === 'create' || urlParams['mode'] === 'update') {
      return <EditPersistentEventPanel />;
    }
  }
}
