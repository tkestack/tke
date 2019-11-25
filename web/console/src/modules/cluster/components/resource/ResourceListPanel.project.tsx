import * as React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from '@tencent/qcloud-lib';
import { allActions } from '../../actions';
import { RootProps } from '../ClusterApp';
import { router } from '../../router';
import { ResourceActionPanel } from './resourceTableOperation/ResourceActionPanel';
import { ResourceTablePanel } from './resourceTableOperation/ResourceTablePanel';
import { ResourceDeleteDialog } from './resourceTableOperation/ResourceDeleteDialog';
import { ResourceSidebarPanel } from './ResourceSidebarPanel';
import { ComputerActionPanel } from './nodeManage/ComputerActionPanel';
import { ComputerTablePanel } from './nodeManage/ComputerTablePanel';
import { BatchDrainComputerDialog } from './nodeManage/BatchDrainComputerDialog';
import { ClusterDetailPanel } from './clusterInfomation/ClusterDetail';
import { ResourceLogPanel } from './resourceTableOperation/ResourceLogPanel';
import { ResourceEventPanel } from './resourceTableOperation/ResourceEventPanel';
import { isEmpty } from '../../../common/utils';
import { UpdateNodeLabelDialog } from './nodeManage/UpdateNodeLabelDialog';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { ContentView, Justify } from '@tencent/tea-component';
import { SubRouter } from '../../models';

const loadingElement: JSX.Element = (
  <div>
    <i className="n-loading-icon" />
    &nbsp; <span className="text">{t('加载中...')}</span>
  </div>
);

export interface ResourceListPanelProps extends RootProps {
  /** subRouterList */
  subRouterList: SubRouter[];
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(
  state => state,
  mapDispatchToProps
)
export class ResourceListPanel extends React.Component<ResourceListPanelProps, {}> {
  componentDidMount() {
    let { actions, subRoot } = this.props,
      { subRouterList } = subRoot;

    // 这里去拉取侧边栏的配置，侧边路由
    !subRouterList.fetched && actions.subRouter.applyFilter({});
  }
  render() {
    let { subRoot, route, namespaceList, subRouterList } = this.props,
      urlParams = router.resolve(route),
      { resourceInfo } = subRoot;

    let content: JSX.Element;
    let headTitle: string = '';
    let resource = urlParams['resourceName'];
    // 判断应该展示什么组件
    switch (resource) {
      case 'info':
        content = <ClusterDetailPanel {...this.props} />;
        headTitle = t('基础信息');
        break;

      case 'log':
        content = <ResourceLogPanel />;
        headTitle = t('日志');
        break;

      case 'event':
        content = <ResourceEventPanel />;
        headTitle = t('事件');
        break;
      default:
        content = isEmpty(resourceInfo) ? (
          loadingElement
        ) : (
          <React.Fragment>
            <ResourceActionPanel />
            <ResourceTablePanel />
          </React.Fragment>
        );
        headTitle = resourceInfo.headTitle;
        break;
    }

    return (
      <React.Fragment>
        <ContentView>
          <ContentView.Body sidebar={<ResourceSidebarPanel subRouterList={subRouterList} />}>
            {namespaceList.fetched ? (
              <ContentView>
                <ContentView.Header>
                  <Justify left={<h2 className="tea-h2">{headTitle || ''}</h2>} />
                </ContentView.Header>
                <ContentView.Body>{content}</ContentView.Body>
              </ContentView>
            ) : (
              <ContentView>
                <ContentView.Body>
                  <div style={{ marginTop: '20px' }}>
                    <i className="n-loading-icon" />
                    &nbsp; <span className="text">加载中...</span>
                  </div>
                </ContentView.Body>
              </ContentView>
            )}
          </ContentView.Body>
        </ContentView>
        <ResourceDeleteDialog />
      </React.Fragment>
    );
  }
}
