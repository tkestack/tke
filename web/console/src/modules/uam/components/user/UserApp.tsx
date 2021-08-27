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
import { Layout, Tabs, TabPanel } from '@tencent/tea-component';
import { GroupPanel } from '../group/GroupPanel';
import { UserPanel } from './UserPanel';
import { router } from '@src/modules/uam/router';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators } from '@tencent/ff-redux/libs/qcloud-lib';
import { allActions } from '@src/modules/uam/actions';
import { RootState } from '@src/modules/uam/models';
const { Body, Content } = Layout;

export interface RootProps extends RootState {
  actions?: typeof allActions;
}

export const UserApp = (props) => {
  const state = useSelector((state) => state);
  // const dispatch = useDispatch();
  // const { actions } = bindActionCreators({ actions: allActions }, dispatch);
  const { route } = state;
  const { sub, action } = router.resolve(route);

  const tabs = [
    { id: 'normal', label: '用户' },
    { id: 'group', label: '用户组' },
  ];

  let header: React.ReactNode;
  if (!action && (!sub || sub === 'normal' || sub === 'group')) {
    // 用户管理页头
    header = <Content.Header showBackButton onBackButtonClick={() => history.back()} title={t('用户管理')} />;
  } else if (sub === 'normal' && action === 'create') {
    // 创建用户头
    header = <Content.Header showBackButton onBackButtonClick={() => history.back()} title={t('新建用户')} />;
  } else if (sub === 'normal' && action === 'detail') {
    // 用户详情头
    header = <Content.Header showBackButton onBackButtonClick={() => history.back()} title={route.queries['name']} />;
  } else if (sub === 'group' && action === 'create') {
    // 新建用户组头
    header = <Content.Header showBackButton onBackButtonClick={() => history.back()} title={t('新建用户组')} />;
  } else if (sub === 'group' && action === 'detail') {
    // 用户组详情头
    header = (
      <Content.Header showBackButton onBackButtonClick={() => history.back()} title={route.queries['groupName']} />
    );
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
              activeId={sub || 'normal'}
              onActive={(value) => {
                router.navigate({ module: 'user', sub: value.id });
              }}
            >
              <TabPanel id="normal">
                <UserPanel />
              </TabPanel>
              <TabPanel id="group">
                <GroupPanel />
              </TabPanel>
            </Tabs>
          </Content.Body>
        </Content>
      </Body>
    </Layout>
  );
};
