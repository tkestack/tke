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
import { DeleteLogDialog } from './DeleteLogDialog';
import { EditLogStashPanel } from './EditLogStashPanel';
import { LogStashDetailPanel } from './LogDetailPanel';
import { LogStashActionPanel } from './LogStashActionPanel';
import { LogStashHeadPanel } from './LogStashHeadPanel';
import { LogStashSubHeadPanel } from './LogStashSubHeadPanel';
import { LogStashTablePanel } from './LogStashTablePanel';
import { OpenLogStashDialog } from './OpenLogStashDialog';
import { LogSettingTablePanel } from './LogSettingTablePanel';

const store = configStore();

export class LogStashAppContainer extends React.Component<any, any> {
  // 页面离开时，清空store
  componentWillUnmount() {
    store.dispatch({ type: ResetStoreAction });
  }

  render() {
    return (
      <Provider store={store}>
        <LogStashApp />
      </Provider>
    );
  }
}

export interface RootProps extends RootState {
  actions?: typeof allActions;
}

const mapDispatchToProps = dispatch =>
  Object.assign(
    {},
    bindActionCreators(
      {
        actions: allActions
      },
      dispatch
    ),
    { dispatch }
  );

@connect(state => state, mapDispatchToProps)
@((router.serve as any)())
class LogStashApp extends React.Component<RootProps, any> {
  render() {
    let { actions, route } = this.props;
    const urlParams = router.resolve(route);

    let content: JSX.Element;
    let { mode } = urlParams;

    if (!mode) {
      content = (
        <React.Fragment>
          <ContentView>
            <ContentView.Header>
              <LogStashHeadPanel />
            </ContentView.Header>
            <ContentView.Body>
              <LogStashActionPanel />
              <LogStashTablePanel />
              <DeleteLogDialog />
            </ContentView.Body>
          </ContentView>
          <OpenLogStashDialog />
        </React.Fragment>
      );
    } else if (mode === 'create') {
      content = (
        <React.Fragment>
          <ContentView>
            <ContentView.Header>
              <LogStashSubHeadPanel />
            </ContentView.Header>
            <ContentView.Body>
              <EditLogStashPanel />
            </ContentView.Body>
          </ContentView>
          <OpenLogStashDialog />
        </React.Fragment>
      );
    } else if (mode === 'update') {
      content = (
        <ContentView>
          <ContentView.Header>
            <LogStashSubHeadPanel />
          </ContentView.Header>
          <ContentView.Body>
            <EditLogStashPanel />
          </ContentView.Body>
        </ContentView>
      );
    } else if (mode === 'detail') {
      content = (
        <ContentView>
          <ContentView.Header>
            <LogStashSubHeadPanel />
          </ContentView.Header>
          <ContentView.Body>
            <LogStashDetailPanel />
          </ContentView.Body>
        </ContentView>
      );
    } else if (mode === 'setting') {
      content = (
        <LogSettingTablePanel />
      );
    }

    return content;
  }
}
