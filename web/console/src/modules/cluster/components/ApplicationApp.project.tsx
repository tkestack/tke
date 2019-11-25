import * as React from 'react';
import { connect, Provider } from 'react-redux';
import { bindActionCreators } from '@tencent/qcloud-lib';
import { RootState } from '../models';
import { allActions } from '../actions';
import { MainBodyLayout } from '../../common/layouts';
import { configStore } from '../stores/RootStore';
import { router } from '../router';
import { ResetStoreAction } from '../../../../helpers';
import { ApplicationHeadPanel } from './ApplicationHeadPanel.project';
import { ResourceListPanel } from './resource/ResourceListPanel';
import { ResourceDetail } from './resource/resourceDetail/ResourceDetail';
import { EditResourcePanel } from './resource/resourceEdition/EditResourcePanel';
import { UpdateResourcePanel } from './resource/resourceEdition/UpdateResourcePanel';

const store = configStore();

export class ApplicationAppContainer extends React.Component<any, any> {
  // 页面离开时，清空store
  componentWillUnmount() {
    store.dispatch({ type: ResetStoreAction });
  }

  render() {
    return (
      <Provider store={store}>
        <ApplicationApp />
      </Provider>
    );
  }
}

export interface RootProps extends RootState {
  actions?: typeof allActions;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
@((router.serve as any)())
class ApplicationApp extends React.Component<RootProps, {}> {
  render() {
    return <ApplicationList {...this.props} />;
  }
}

class ApplicationList extends React.Component<RootProps, {}> {
  componentDidMount() {
    let { actions, route } = this.props,
      urlParams = router.resolve(route);
    actions.region.fetch();
    // 这里需要去判断一下当前的resource是否需要进行namespace 路由的更新，参考resourceTabelPanel
    let { resourceName: resource, type: resourceType } = urlParams;
    resource ? actions.resource.initResourceName(resource) : actions.resource.initResourceName('np');
    // 判断当前是否需要去更新np的路由
    let isNeedFetchNamespace =
      resourceType === 'resource' || resourceType === 'service' || resourceType === 'config' || resource === 'pvc';
    actions.resource.toggleIsNeedFetchNamespace(isNeedFetchNamespace ? true : false);
    actions.projectNamespace.initProjectList();
  }

  componentWillReceiveProps(nextProps: RootProps) {
    let { route, actions, namespaceSelection, subRoot } = nextProps,
      newUrlParam = router.resolve(route),
      { mode } = subRoot;
    let newMode = newUrlParam['mode'];
    let oldMode = this.props.subRoot.mode;

    if (newMode !== '' && oldMode !== newMode && newMode !== mode) {
      actions.resource.selectMode(newMode);
      // 这里是判断回退动作，取消动作等的时候
      newUrlParam['mode'] !== 'list' &&
        actions.resource.applyFilter({ namespace: namespaceSelection, clusterId: route.queries['clusterId'] });
    }
  }

  render() {
    let { route, subRoot } = this.props,
      urlParams = router.resolve(route);
    let urlMode = urlParams['mode'];
    if (!urlMode || urlMode === 'list') {
      return (
        <div className="manage-area manage-area-secondary">
          <ApplicationHeadPanel />
          <ResourceListPanel subRouterList={subRoot.subRouterList.data.records} />
        </div>
      );
    } else if (urlMode === 'detail') {
      return <ResourceDetail />;
    } else if (urlMode === 'create' || urlMode === 'modify') {
      return <EditResourcePanel {...this.props} />;
    } else if (urlMode === 'update') {
      return <UpdateResourcePanel />;
    }
  }
}
