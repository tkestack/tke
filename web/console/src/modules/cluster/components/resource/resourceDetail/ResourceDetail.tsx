import * as React from 'react';
import { connect } from 'react-redux';

import { ContentView, TabPanel, Tabs } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';

import { FormLayout, MainBodyLayout } from '../../../../common/layouts';
import { isEmpty } from '../../../../common/utils';
import { allActions } from '../../../actions';
import { router } from '../../../router';
import { RootProps } from '../../ClusterApp';
import { ResourceHeaderPanel } from '../ResourceHeaderPanel';
import { ResourceDetailEventActionPanel } from './ResourceDetailEventActionPanel';
import { ResourceDetailEventTablePanel } from './ResourceDetailEventTablePanel';
import { ResourceDetailPanel } from './ResourceDetailPanel';
import { ResourceGrayUpgradeDialog } from './ResourceGrayUpgradeDialog';
import { ResourceModifyHistoryPanel } from './ResourceModifyHistoryPanel';
import { ResourceNamespaceDetailPanel } from './ResourceNamespaceDetailPanel';
import { ResourceNodeDetailPanel } from './ResourceNodeDetailPanel';
import { ResourcePodActionPanel } from './ResourcePodActionPanel';
import { ResourcePodDeleteDialog } from './ResourcePodDeleteDialog';
import { ResourcePodLogPanel } from './ResourcePodLogPanel';
import { ResourcePodPanel } from './ResourcePodPanel';
import { ResourcePodRemoteLoginDialog } from './ResourcePodRemoteLoginDialog';
import { ResourceRollbackDialog } from './ResourceRollbackDialog';
import { ResourceTappPodDeleteDialog } from './ResourceTappPodDeleteDialog';
import { ResourceYamlActionPanel } from './ResourceYamlActionPanel';
import { ResourceYamlPanel } from './ResourceYamlPanel';

/** 判断当前详情页面是否在节点里面 */
export const IsInNodeManageDetail = (type: string) => {
  return type === 'nodeManage';
};

/** 判断yaml是否需要编辑yaml */
export const IsNeedModifyYaml = (resourceName: string) => {
  return resourceName !== 'sc';
};

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ResourceDetail extends React.Component<RootProps, {}> {
  componentWillUnmount() {
    let { actions } = this.props;
    // 清除detail里面的所有内容
    actions.resourceDetail.clearDetail();
  }

  componentWillMount() {
    let {
      actions,
      subRoot: {
        resourceInfo,
        detailResourceOption: { detailResourceName },
        resourceOption: { resourceList }
      },
      route
    } = this.props;
    let urlParams = router.resolve(route);
    let tab = urlParams['tab'];
    //当直接跳转到其他tab页的时候，这时需要对集群内资源进行初始化。
    if (resourceInfo.requestType && resourceInfo.requestType.useDetailInfo && resourceList.data.recordCount > 0) {
      let list = tab ? resourceInfo.requestType.detailInfoList[tab] : resourceInfo.requestType.detailInfoList['info'];
      if (list) {
        actions.resource.initDetailResourceName(list[0].value);
      }
    }
  }

  render() {
    let { subRoot, actions, route } = this.props,
      urlParams = router.resolve(route),
      { resourceInfo } = subRoot;

    let tabs = !isEmpty(resourceInfo) ? resourceInfo.detailField.tabList : [];

    // 默认选中第一个
    let selected = tabs[0];
    if (urlParams['tab']) {
      let finder = tabs.find(x => x.id === urlParams['tab']);
      if (finder) {
        selected = finder;
      }
    }

    return (
      <div className="manage-area">
        <ContentView>
          <ContentView.Header>
            <ResourceHeaderPanel />
          </ContentView.Header>
          <ContentView.Body>
            <Tabs
              ceiling
              tabs={tabs}
              activeId={selected ? selected.id : ''}
              onActive={tab => {
                router.navigate(Object.assign({}, urlParams, { tab: tab.id }), route.queries);
                if (resourceInfo.requestType.useDetailInfo) {
                  actions.resource.changeDetailTab(tab.id);
                }
              }}
            >
              <TabPanel id="pod">
                <ResourcePodActionPanel />
                <ResourcePodPanel />
              </TabPanel>
              <TabPanel id="history">
                <ResourceModifyHistoryPanel />
              </TabPanel>
              <TabPanel id="event">
                <ResourceDetailEventActionPanel />
                <ResourceDetailEventTablePanel />
              </TabPanel>
              <TabPanel id="log">
                <ResourcePodLogPanel />
              </TabPanel>
              <TabPanel id="info">
                <ResourceDetailPanel />
              </TabPanel>
              <TabPanel id="yaml">
                {!IsInNodeManageDetail(urlParams['type']) && IsNeedModifyYaml(urlParams['resourceName']) && (
                  <ResourceYamlActionPanel />
                )}
                <ResourceYamlPanel />
              </TabPanel>
              <TabPanel id="nsInfo">
                <ResourceNamespaceDetailPanel />
              </TabPanel>
              <TabPanel id="nodeInfo">
                <ResourceNodeDetailPanel />
              </TabPanel>
            </Tabs>

            <ResourceRollbackDialog />
            <ResourcePodDeleteDialog />
            <ResourcePodRemoteLoginDialog />
            <ResourceTappPodDeleteDialog />
            <ResourceGrayUpgradeDialog />
          </ContentView.Body>
        </ContentView>
        <MainBodyLayout />
      </div>
    );
  }
}
