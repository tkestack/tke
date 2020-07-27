import * as React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { RootProps } from '../ChartApp';
import { HeaderPanel } from './HeaderPanel';
import { MainBodyLayout } from '../../../../common/layouts';
import { Button, Tabs, TabPanel, Card } from '@tea/component';
import { router } from '../../../router';
import { BasicInfoPanel } from './BasicInfoPanel';
import { VersionTablePanel } from './VersionTablePanel';
import { FileTreePanel } from './FileTreePanel';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ChartDetail extends React.Component<RootProps, {}> {
  componentWillUnmount() {
    let { actions } = this.props;
    actions.chart.detail.updateChartWorkflow.reset();
    actions.chart.detail.clearEditorState();
    actions.chart.detail.clearValidatorState();

    /** 取消轮询 */
    actions.chart.detail.clearPolling();
  }

  componentDidMount() {
    const { actions, route } = this.props;
    let urlParam = router.resolve(route);
    let tab = 'detail';
    if (urlParam['tab']) tab = urlParam['tab'];
    this.changeTab(tab);
  }

  changeTab(tab: string) {
    const { actions, route, chartEditor } = this.props;
    switch (tab) {
      case 'detail':
      case 'file': {
        /** 取消轮询 */
        actions.chart.detail.clearPolling();
        /** 查询具体chart，从而Detail可以用到 */
        actions.chart.detail.applyFilter({
          namespace: route.queries['cg'],
          name: route.queries['chart'],
          projectID: route.queries['prj']
        });

        break;
      }
      case 'version': {
        /** 使用轮询是因为需要更新已删除版本的状态 */
        /** 查询具体chart，从而Detail可以用到 */
        actions.chart.detail.poll({
          namespace: route.queries['cg'],
          name: route.queries['chart'],
          projectID: route.queries['prj']
        });

        break;
      }
    }
  }

  render() {
    const tabs = [
      { id: 'detail', label: t('基本信息') },
      { id: 'version', label: t('版本管理') },
      { id: 'file', label: t('目录树') }
    ];

    let { actions, route } = this.props,
      urlParam = router.resolve(route);
    let tab = 'detail';
    if (urlParam['tab']) tab = urlParam['tab'];
    return (
      <React.Fragment>
        <HeaderPanel />
        <MainBodyLayout className="secondary-main">
          <Tabs
            tabs={tabs}
            defaultActiveId={tab}
            onActive={tab => {
              router.navigate({ sub: 'chart', mode: 'detail', tab: tab.id }, route.queries);
              this.changeTab(tab.id);
            }}
            className="tea-tabs--ceiling"
          >
            <TabPanel id="detail">
              <BasicInfoPanel {...this.props} />
            </TabPanel>
            <TabPanel id="version">
              <VersionTablePanel {...this.props} />
            </TabPanel>
            <TabPanel id="file">
              <FileTreePanel {...this.props} />
            </TabPanel>
          </Tabs>
        </MainBodyLayout>
      </React.Fragment>
    );
  }
}
