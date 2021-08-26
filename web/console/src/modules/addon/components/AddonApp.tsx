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
import { connect, Provider } from 'react-redux';

import { bindActionCreators } from '@tencent/ff-redux';
import { ContentView } from '@tencent/tea-component';

import { ResetStoreAction } from '../../../../helpers';
import { allActions } from '../actions';
import { RootState } from '../models';
import { router } from '../router';
import { configStore } from '../stores/RootStore';
import { AddonActionPanel } from './AddonActionPanel';
import { AddonDeleteDialog } from './AddonDeleteDialog';
import { AddonDetail } from './AddonDetail';
import { AddonHeadPanel } from './AddonHeadPanel';
import { AddonSubpageHeadPanel } from './AddonSubpageHeadPanel';
import { AddonTablePanel } from './AddonTablePanel';
import { EditAddonPanel } from './EditAddonPanel';

export const store = configStore();

export class AddonAppContainer extends React.Component<any, any> {
  // 页面离开时，清空store
  componentWillUnmount() {
    store.dispatch({ type: ResetStoreAction });
  }
  render() {
    return (
      <Provider store={store}>
        <AddonApp />
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
class AddonApp extends React.Component<RootProps, any> {
  render() {
    let { route } = this.props,
      urlParams = router.resolve(route);

    let mode = urlParams['mode'];
    if (!mode) {
      return (
        <React.Fragment>
          <ContentView>
            <ContentView.Header>
              <AddonHeadPanel />
            </ContentView.Header>
            <ContentView.Body>
              <AddonActionPanel />
              <AddonTablePanel />
            </ContentView.Body>
          </ContentView>
          <AddonDeleteDialog />
        </React.Fragment>
      );
    } else if (mode === 'detail') {
      return <AddonDetail />;
    } else if (mode === 'create') {
      return (
        <ContentView>
          <ContentView.Header>
            <AddonSubpageHeadPanel route={route} />
          </ContentView.Header>
          <ContentView.Body>
            <EditAddonPanel />
          </ContentView.Body>
        </ContentView>
      );
    }
  }
}
