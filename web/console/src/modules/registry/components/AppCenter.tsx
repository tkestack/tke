import * as React from 'react';
import { connect, Provider } from 'react-redux';
import { bindActionCreators } from '@tencent/ff-redux';
import { RootState } from '../models';
import { allActions } from '../actions';
import { router } from '../router';
import { ChartList } from './chart/list/ChartList';
import { ChartDetail } from './chart/detail/ChartDetail';
import { ChartGroupList } from './chartgroup/list/ChartGroupList';
import { ChartGroupCreate } from './chartgroup/create/ChartGroupCreate';
import { ChartGroupDetail } from './chartgroup/detail/ChartGroupDetail';
import { ContentView, Tabs, TabPanel, Layout } from '@tea/component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
const { Body, Content } = Layout;

export interface RootProps extends RootState {
  actions?: typeof allActions;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class AppCenter extends React.Component<RootProps, {}> {
  render() {
    let { route } = this.props,
      urlParam = router.resolve(route);
    const tabs = [
      { id: 'chart', label: t('模板') },
      { id: 'chartgroup', label: t('仓库') }
    ];
    let tab = 'chart';
    if (urlParam['sub']) tab = urlParam['sub'];
    if (!urlParam['mode'] || urlParam['mode'] === 'list') {
      return (
        <div className="manage-area">
          <Layout>
            <Body>
              <Content>
                <Content.Header title={t('Helm模板')} />
                <Content.Body>
                  <Tabs
                    ceiling
                    animated={false}
                    tabs={tabs}
                    defaultActiveId={tab}
                    onActive={tab => {
                      router.navigate({ sub: tab.id });
                    }}
                  >
                    <TabPanel id="chart">
                      <ChartList {...this.props} />
                    </TabPanel>
                    <TabPanel id="chartgroup">
                      <ChartGroupList {...this.props} />
                    </TabPanel>
                  </Tabs>
                </Content.Body>
              </Content>
            </Body>
          </Layout>
        </div>
      );
    }

    if (!urlParam['sub'] || urlParam['sub'] === 'chart') {
      if (urlParam['mode'] === 'create') {
        return <div className="manage-area">{/* <ChartCreate {...this.props} /> */}</div>;
      } else if (urlParam['mode'] === 'detail') {
        return (
          <div className="manage-area">
            <ChartDetail {...this.props} />
          </div>
        );
      }
    } else {
      if (urlParam['mode'] === 'create') {
        return (
          <div className="manage-area">
            <ChartGroupCreate {...this.props} />
          </div>
        );
      } else if (urlParam['mode'] === 'detail') {
        return (
          <div className="manage-area">
            <ChartGroupDetail {...this.props} />
          </div>
        );
      }
    }
  }
}
