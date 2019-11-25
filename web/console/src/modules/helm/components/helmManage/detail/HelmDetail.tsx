import * as React from 'react';
import { Tabs, TabPanel } from '@tea/component/tabs';

import { RootProps } from '../../HelmApp';
import { HelmDetailBasicInfoPanel } from './HelmDetailBasicInfoPanel';
import { HistoryTablePanel } from './HistoryTablePanel';
// import { ValueYamlPanel } from './ValueYamlPanel';
import { router } from '../../../router';
import { MainBodyLayout } from '../../../../common/layouts';
import { insertCSS } from '@tencent/qcloud-lib';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
export class HelmDetail extends React.Component<RootProps, {}> {
  componentWillUnmount() {
    this.props.actions.detail.clear();
  }
  componentDidMount() {
    insertCSS(
      'codemirror',
      `.height400 .react-codemirror2 {
        height:400px !important;
      }
      .height200 .react-codemirror2 {
        height: 200px !important;
      }
      .CodeMirror {
      /* Set height, width, borders, and global font properties here */
      font-family: monospace;
      height: 100%;
    }`
    );
    let {
      actions,
      route
      // listState: { regionList }
    } = this.props;
    actions.detail.fetchHelm(route.queries['helmName']);

    // let isNeedFetchRegion = regionList.data.recordCount ? false : true;
    // isNeedFetchRegion && actions.region.applyFilter({});
  }

  goBack() {
    let { actions, route } = this.props,
      urlParams = router.resolve(route);
    router.navigate(
      {},
      {
        rid: route.queries['rid'],
        clusterId: route.queries['clusterId']
      }
    );
  }
  render() {
    let {
      route,
      detailState: { helm }
    } = this.props;
    const urlParams = router.resolve(route);

    let tabs = [
      {
        id: 'info',
        label: t('应用详情')
      },
      {
        id: 'history',
        label: t('版本历史')
      }
      // {
      //   id: 'value',
      //   label: 'Value.yaml'
      // }
    ];
    //默认选中第一个
    let selected = tabs[0];
    if (urlParams['tab']) {
      let finder = tabs.find(x => x.id === urlParams['tab']);
      if (finder) {
        selected = finder;
      }
    }
    if (!helm) {
      return <noscript />;
    }
    return (
      <div className="manage-area">
        <div className="manage-area-title secondary-title">
          <a href="javascript:;" onClick={() => this.goBack()} className="back-link">
            <i className="btn-back-icon" />
            <span>{t('返回')}</span>
          </a>
          <span className="line-icon">|</span>
          <h2>
            {helm.name} {t('详情')}
          </h2>
        </div>
        <MainBodyLayout className="secondary-main">
          <Tabs
            tabs={tabs}
            activeId={selected.id}
            onActive={tab => {
              router.navigate(Object.assign({}, urlParams, { tab: tab.id }), route.queries);
            }}
            className="tea-tabs--ceiling"
          >
            <TabPanel id="info">
              <HelmDetailBasicInfoPanel {...this.props} />
            </TabPanel>
            <TabPanel id="history">
              <HistoryTablePanel {...this.props} />
            </TabPanel>
            {/* <TabPanel id="value">
              <ValueYamlPanel {...this.props} />
            </TabPanel> */}
          </Tabs>
        </MainBodyLayout>
      </div>
    );
  }
}
