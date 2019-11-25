import * as React from 'react';
import { RootProps } from './AddonApp';
import { router } from '../router';
import { bindActionCreators } from '@tencent/qcloud-lib';
import { connect } from 'react-redux';
import { allActions } from '../actions';
import { ContentView, Tabs, TabPanel } from '@tencent/tea-component';
import { AddonDetailHeadPanel } from './AddonDetailHeadPanel';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { AddonDetailPanel } from './AddonDetailPanel';

/** 详情页的tab列表 */
const tabList: any[] = [
  {
    id: 'info',
    label: t('详情')
  }
];

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(
  state => state,
  mapDispatchToProps
)
export class AddonDetail extends React.Component<RootProps, {}> {
  componentDidMount() {
    let { actions, region } = this.props;
    // 进行地域的拉取
    region.list.fetched !== true && actions.region.applyFilter({});
  }

  render() {
    let { route } = this.props,
      urlParams = router.resolve(route);

    let { tab } = urlParams;
    // 默认选择第一个
    let selected = tabList[0];
    if (tab) {
      let finder = tabList.find(x => x.id === tab);
      if (finder) {
        selected = finder;
      }
    }

    return (
      <ContentView>
        <ContentView.Header>
          <AddonDetailHeadPanel />
        </ContentView.Header>
        <ContentView.Body>
          <Tabs
            ceiling
            tabs={tabList}
            activeId={selected ? selected.id : ''}
            onActive={tab => {
              router.navigate(Object.assign({}, urlParams, { tab: tab.id }), route.queries);
            }}
          >
            <TabPanel id="info">
              <AddonDetailPanel />
            </TabPanel>
          </Tabs>
        </ContentView.Body>
      </ContentView>
    );
  }
}
