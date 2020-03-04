import * as React from 'react';
import { connect } from 'react-redux';

import { ContentView, TabPanel, Tabs } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { FormLayout, MainBodyLayout } from '../../common/layouts';
import { allActions } from '../actions';
import { router } from '../router';
import { NamespaceActionPanel } from './NamespaceActionPanel';
import { NamespaceTablePanel } from './NamespaceTablePanel';
import { RootProps } from './ProjectApp';
import { ProjectDetailPanel } from './ProjectDetailPanel';
import { SubpageHeadPanel } from './SubpageHeadPanel';

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
  }
  componentWillUnmount() {
    let { actions } = this.props;

    // actions.project.selectProject([]);
  }

  render() {
    let tabs = [
      {
        id: 'info',
        label: t('业务信息')
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
    let urlParams = router.resolve(route);
    return (
      <ContentView>
        <ContentView.Header>
          <SubpageHeadPanel />
        </ContentView.Header>
        <ContentView.Body>
          <Tabs
            ceiling
            tabs={tabs}
            activeId={selected.id}
            onActive={tab => {
              router.navigate(Object.assign({}, urlParams, { tab: tab.id }), route.queries);
              this.setState({ tabId: tab.id });
            }}
          >
            <TabPanel id="namespace">
              <NamespaceActionPanel {...this.props} />
              <NamespaceTablePanel {...this.props} />
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
