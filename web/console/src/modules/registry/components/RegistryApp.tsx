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
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { ResetStoreAction } from '../../../../helpers';
import { allActions } from '../actions';
import { RootState } from '../models';
import { router } from '../router';
import { configStore } from '../stores/RootStore';
import { ApiKeyContainer } from './apikey/ApiKeyContainer';
import { RepoContainer } from './repo/RepoContainer';
import { ChartApp } from './chart/ChartApp';
import { ChartGroupApp } from './chartgroup/ChartGroupApp';
import { AppCenter } from './AppCenter';

const store = configStore();

export class RegistryAppContainer extends React.Component<any, any> {
  componentWillUnmount() {
    store.dispatch({ type: ResetStoreAction });
  }

  render() {
    return (
      <Provider store={store}>
        <RegistryApp />
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
class RegistryApp extends React.Component<RootProps, {}> {
  componentDidMount() {
    this.props.actions.image.fetchDockerRegUrl.fetch();
  }

  render() {
    let { route } = this.props,
      urlParam = router.resolve(route);
    if (!urlParam['sub'] || urlParam['sub'] === 'chart' || urlParam['sub'] === 'chartgroup') {
      return <AppCenter {...this.props} />;
    } else if (urlParam['sub'] === 'apikey') {
      return <ApiKeyContainer {...this.props} />;
    } else if (urlParam['sub'] === 'repo') {
      return <RepoContainer {...this.props} />;
    } else {
      return <RepoContainer {...this.props} />;
    }
  }
}
