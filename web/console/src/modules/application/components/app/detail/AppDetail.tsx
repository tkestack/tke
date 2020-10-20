import * as React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { RootProps } from '../AppContainer';
import { HeaderPanel } from './HeaderPanel';
import { MainBodyLayout } from '../../../../common/layouts';
import { Button, Tabs, TabPanel, Card } from '@tea/component';
import { router } from '../../../router';
import { BasicInfoPanel } from './BasicInfoPanel';
import { ResourceTablePanel } from './ResourceTablePanel';
import { HistoryTablePanel } from './HistoryTablePanel';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class AppDetail extends React.Component<RootProps, {}> {
  componentWillUnmount() {
    let { actions } = this.props;
    actions.app.detail.updateAppWorkflow.reset();
    actions.app.detail.clearEditorState();
    actions.app.detail.clearValidatorState();

    actions.app.resource.clearPolling();
  }

  componentDidMount() {
    const { actions, route, appEditor } = this.props;
    let urlParam = router.resolve(route);
    let tab = 'resource';
    if (urlParam['tab']) tab = urlParam['tab'];
    this.changeTab(tab);
  }

  changeTab(tab: string) {
    const { actions, route, appEditor } = this.props;
    switch (tab) {
      case 'detail': {
        actions.app.resource.clearPolling();

        /** 拉取仓库列表 */
        actions.chartGroup.list.applyFilter({});
        /** 查询具体应用，从而Detail可以用到 */
        actions.app.detail.fetchApp({
          cluster: route.queries['cluster'],
          namespace: route.queries['namespace'],
          name: route.queries['app']
        });
        break;
      }
      case 'resource': {
        actions.app.resource.poll({
          cluster: route.queries['cluster'],
          namespace: route.queries['namespace'],
          name: route.queries['app']
        });
        break;
      }
      case 'history': {
        actions.app.resource.clearPolling();

        actions.app.history.applyFilter({
          cluster: route.queries['cluster'],
          namespace: route.queries['namespace'],
          name: route.queries['app']
        });
        break;
      }
    }
  }

  render() {
    const tabs = [
      { id: 'resource', label: t('资源列表') },
      { id: 'detail', label: t('应用详情') },
      { id: 'history', label: t('版本历史') }
    ];

    let { actions, route } = this.props,
      urlParam = router.resolve(route);
    let tab = 'resource';
    if (urlParam['tab']) tab = urlParam['tab'];
    return (
      <React.Fragment>
        <HeaderPanel />
        <MainBodyLayout className="secondary-main">
          <Tabs
            tabs={tabs}
            defaultActiveId={tab}
            onActive={tab => {
              router.navigate({ sub: 'app', mode: 'detail', tab: tab.id }, route.queries);
              this.changeTab(tab.id);
            }}
            className="tea-tabs--ceiling"
          >
            <TabPanel id="detail">
              <BasicInfoPanel {...this.props} />
            </TabPanel>
            <TabPanel id="resource">
              <ResourceTablePanel {...this.props} />
            </TabPanel>
            <TabPanel id="history">
              <HistoryTablePanel {...this.props} />
            </TabPanel>
          </Tabs>
        </MainBodyLayout>
      </React.Fragment>
    );
  }
}
