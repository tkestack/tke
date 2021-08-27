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

import { ResetStoreAction } from '../../../../helpers';
import { allActions } from '../actions';
import { RootState } from '../models';
import { router } from '../router';
import { configStore } from '../stores/RootStore';
import { HelmCreate } from './helmManage/create/HelmCreate';
import { HelmDetail } from './helmManage/detail/HelmDetail';
import { HelmActionPanel } from './helmManage/list/HelmActionPanel';
import { HelmHeadPanel } from './helmManage/list/HelmHeadPanel';
import { HelmTablePanel } from './helmManage/list/HelmTablePanel';

const store = configStore();

export class HelmAppContainer extends React.Component<any, any> {
  // 页面离开时，清空store
  componentWillUnmount() {
    store.dispatch({ type: ResetStoreAction });
  }

  render() {
    return (
      <Provider store={store}>
        <HelmApp />
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
class HelmApp extends React.Component<RootProps, {}> {
  componentDidMount() {
    let { actions } = this.props;
  }

  render() {
    let { route, isShowTips, actions } = this.props,
      urlParam = router.resolve(route);
    if (!urlParam['sub']) {
      return (
        <ContentView>
          <ContentView.Header>
            <HelmHeadPanel {...this.props} />
          </ContentView.Header>
          <ContentView.Body>
            <HelmActionPanel {...this.props} />
            <HelmTablePanel {...this.props} />
          </ContentView.Body>
        </ContentView>
      );
    } else if (urlParam['sub'] === 'create') {
      return (
        <div className="manage-area">
          <HelmCreate {...this.props} />
        </div>
      );
    } else if (urlParam['sub'] === 'detail') {
      return (
        <div className="manage-area">
          <HelmDetail {...this.props} />
        </div>
      );
    }
  }
}
