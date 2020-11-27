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
    const { actions } = this.props;
    // 清除detail里面的所有内容
    actions.resourceDetail.clearDetail();
  }

  componentWillMount() {
    const {
      actions,
      clusterVersion,
      subRoot: {
        resourceDetailState: { resourceDetailInfo }
      },
      route
    } = this.props;
    const { rid, np, clusterId, resourceIns } = route.queries;

    // 进入即进行资源详情的拉取，前提是当前resourceDetailInfo没有拉取过，且clusterVersion为''
    !resourceDetailInfo.selection &&
      clusterVersion !== '' &&
      actions.resourceDetail.resourceInfo.applyFilter({
        regionId: +rid,
        namespace: np,
        clusterId,
        specificName: resourceIns
      });
  }

  render() {
    const { subRoot, actions, route } = this.props,
      urlParams = router.resolve(route),
      { resourceInfo } = subRoot;

    const tabs = !isEmpty(resourceInfo) ? resourceInfo.detailField.tabList : [];

    // 默认选中第一个
    let selected = tabs[0];
    if (urlParams['tab']) {
      const finder = tabs.find(x => x.id === urlParams['tab']);
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
      </div>
    );
  }
}
