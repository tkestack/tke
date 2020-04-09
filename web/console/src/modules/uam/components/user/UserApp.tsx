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
  let urlParam = router.resolve(route);
  const { sub, action } = urlParam;

  const tabs = [
    { id: 'normal', label: '用户' },
    { id: 'group', label: '用户组' },
  ];

  let header;
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
                console.log('tab value:', value);
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
