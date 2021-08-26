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
import { connect } from 'react-redux';
import { ContentView, TabPanel, Tabs } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { allActions } from '../actions';
import { router } from '../router';
import { NamespaceActionPanel } from './NamespaceActionPanel';
import { NamespaceTablePanel } from './NamespaceTablePanel';
import { RootProps } from './ProjectApp';
import { ProjectDetailPanel } from './ProjectDetailPanel';
import { SubpageHeadPanel } from './SubpageHeadPanel';
import { ProjectHeadPanel } from '@src/modules/project/components/ProjectHeadPanel';
import { CreateNamespacePanel } from '@src/modules/project/components/CreateNamespacePanel';
import { UserPanel } from './user/UserPanel';
import { DetailSubProjectPanel } from './DetailSubProjectPanel';
import { DetailSubProjectActionPanel } from './DetailSubProjectActionPanel';

interface ProjectDetailState {
  /** tabKey */
  tabId?: string;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ProjectDetail extends React.Component<RootProps, ProjectDetailState> {
  constructor(props, context) {
    super(props, context);
    let { route } = props;
    let urlParams = router.resolve(route);
    this.state = {
      tabId: urlParams['tab'] || 'info'
    };
  }

  componentDidMount() {
    let { actions, route } = this.props;
    actions.project.fetchDetail(route.queries['projectId']);
    actions.project.fetch();
  }

  render() {
    let tabs = [
      {
        id: 'info',
        label: t('业务信息')
      },
      { id: 'member', label: t('成员列表') },
      {
        id: 'subProject',
        label: t('子业务')
      },
      {
        id: 'namespace',
        label: t('Namespace列表')
      }
    ];

    /** 默认选中第一个tab */
    let selected = tabs[0];
    let finder = tabs.find(x => x.id === this.state.tabId);
    if (finder) {
      selected = finder;
    }
    let { route } = this.props;
    const urlParams = router.resolve(route);
    const { action } = urlParams;
    let header;
    if (action === 'createNS') {
      header = <ProjectHeadPanel isNeedBack={true} title={t('新建Namespace')} />;
    } else if (action === 'create') {
      header = <ProjectHeadPanel isNeedBack={true} title={t('添加成员')} />;
    } else {
      header = <SubpageHeadPanel />;
    }
    return (
      <ContentView>
        <ContentView.Header>{header}</ContentView.Header>
        <ContentView.Body>
          <Tabs
            ceiling
            tabs={tabs}
            activeId={selected.id}
            onActive={tab => {
              router.navigate(Object.assign({}, urlParams, { tab: tab.id, action: '' }), route.queries);
              this.setState({ tabId: tab.id });
            }}
          >
            <TabPanel id="namespace">
              {action === 'createNS' ? (
                <CreateNamespacePanel />
              ) : (
                <>
                  <NamespaceActionPanel {...this.props} />
                  <NamespaceTablePanel {...this.props} />
                </>
              )}
            </TabPanel>
            <TabPanel id="member">
              <UserPanel />
            </TabPanel>
            <TabPanel id="subProject">
              <DetailSubProjectActionPanel />
              <DetailSubProjectPanel />
            </TabPanel>
            <TabPanel id="info">
              <ProjectDetailPanel {...this.props} />
            </TabPanel>
          </Tabs>
        </ContentView.Body>
      </ContentView>
    );
  }
}
