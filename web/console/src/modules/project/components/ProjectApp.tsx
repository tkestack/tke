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
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { ContentView } from '@tencent/tea-component';

import { ResetStoreAction } from '../../../../helpers';
import { allActions } from '../actions';
import { RootState } from '../models';
import { router } from '../router';
import { configStore } from '../stores/RootStore';
import { CreateProjectPanel } from './CreateProjectPanel';
import { ProjectActionPanel } from './ProjectActionPanel';
import { ProjectDetail } from './ProjectDetail';
import { ProjectHeadPanel } from './ProjectHeadPanel';
import { ProjectTablePanel } from './ProjectTablePanel';
import { PlatformTypeEnum } from '../constants/Config';

const store = configStore();

export class ProjectAppContainer extends React.Component<any, any> {
  // 页面离开时，清空store
  componentWillUnmount() {
    store.dispatch({ type: ResetStoreAction });
  }

  render() {
    return (
      <Provider store={store}>
        <ProjectApp />
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
class ProjectApp extends React.Component<RootProps, {}> {
  constructor(props, context) {
    super(props, context);
    /// #if project
    props.actions.bussiness.initPlatformType(PlatformTypeEnum.Business);
    props.actions.bussiness.userInfo.fetch();
    /// #endif

    /// #if tke
    props.actions.bussiness.initPlatformType(PlatformTypeEnum.Manager);
    /// #endif
  }
  render() {
    const { route } = this.props;
    const { sub } = router.resolve(route);
    if (!sub) {
      return (
        <ContentView>
          <ContentView.Header>
            <ProjectHeadPanel title={t('业务管理')} />
          </ContentView.Header>
          <ContentView.Body>
            <ProjectActionPanel />
            <ProjectTablePanel />
          </ContentView.Body>
        </ContentView>
      );
    } else if (sub === 'detail') {
      return <ProjectDetail {...this.props} />;
    } else if (sub === 'create') {
      return (
        <ContentView>
          <ContentView.Header>
            <ProjectHeadPanel isNeedBack={true} title={t('新建业务')} />
          </ContentView.Header>
          <ContentView.Body>
            <CreateProjectPanel />
          </ContentView.Body>
        </ContentView>
      );
    }
  }
}
