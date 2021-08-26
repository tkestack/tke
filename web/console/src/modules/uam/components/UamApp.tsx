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

import { ResetStoreAction } from '../../../../helpers';
import { allActions } from '../actions';
import { RootState } from '../models';
import { router } from '../router';
import { configStore } from '../stores/RootStore';
import { StrategyApp } from './strategy/StrategyApp';
import { UserApp } from './user/UserApp';
import { RoleApp } from './role/RoleApp';

const store = configStore();

export class UamAppContainer extends React.Component<any, any> {
  // 页面离开时，清空store
  componentWillUnmount() {
    store.dispatch({ type: ResetStoreAction });
  }
  render() {
    return (
      <Provider store={store}>
        <UamApp />
      </Provider>
    );
  }
}

export interface RootProps extends RootState {
  actions?: typeof allActions;
}
const mapDispatchToProps = (dispatch) =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect((state) => state, mapDispatchToProps)
@((router.serve as any)())
class UamApp extends React.Component<RootProps, {}> {
  render() {
    let { route } = this.props;
    let urlParam = router.resolve(route);
    const { module } = urlParam;

    let content: React.ReactNode;
    if (module === 'strategy') {
      content = <StrategyApp />;
    } else if (module === 'role') {
      content = <RoleApp />;
    } else {
      content = <UserApp />;
    }

    return content;
  }
}
