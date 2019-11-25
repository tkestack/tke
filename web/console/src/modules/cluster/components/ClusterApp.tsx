import * as React from 'react';
import { connect, Provider } from 'react-redux';
import { bindActionCreators } from '@tencent/qcloud-lib';
import { RootState } from '../models';
import { allActions } from '../actions';
import { configStore } from '../stores/RootStore';
import { router } from '../router';
import { ResetStoreAction } from '../../../../helpers';
import { ClusterHeadPanel } from './clusterManage/ClusterHeadPanel';
import { ClusterActionPanel } from './clusterManage/ClusterActionPanel';
import { ClusterTablePanel } from './clusterManage/ClusterTablePanel';
import { ResourceContainerPanel } from './resource/ResourceContainerPanel';
import { ContentView, Card } from '@tencent/tea-component';
import { ClusterDeleteDialog } from './clusterManage/ClusterDeleteDialog';
import { ClusterStatusDialog } from './clusterManage/ClusterStatusDialog';
import { CreateClusterPanel } from './clusterManage/CreateClusterPanel';
import { CreateICPanel } from './clusterManage/CreateICPanel';
import { ModifyClusterNameDialog } from './clusterManage/ModifyClusterNameDialog';
import { TcrRegistyDeployDialog } from './clusterManage/TcrRegistyDeployDialog';

export const store = configStore();

export class ClusterAppContainer extends React.Component<any, any> {
  // 页面离开时，清空store
  componentWillUnmount() {
    store.dispatch({ type: ResetStoreAction });
  }
  render() {
    return (
      <Provider store={store}>
        <ClusterApp />
      </Provider>
    );
  }
}

export interface RootProps extends RootState {
  actions?: typeof allActions;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
@((router.serve as any)())
class ClusterApp extends React.Component<RootProps, {}> {
  render() {
    let { route } = this.props;
    let urlParam = router.resolve(route);
    if (!urlParam['sub']) {
      return (
        <ContentView>
          <ContentView.Header>
            <ClusterHeadPanel />
          </ContentView.Header>
          <ContentView.Body>
            <ClusterActionPanel />
            <ClusterTablePanel />
            <ClusterDeleteDialog />
            <ClusterStatusDialog />
            <ModifyClusterNameDialog />
            <TcrRegistyDeployDialog />
          </ContentView.Body>
        </ContentView>
      );
    } else if (urlParam['sub'] === 'sub') {
      return <ResourceContainerPanel />;
    } else if (urlParam['sub'] === 'create') {
      return <CreateClusterPanel />;
    } else if (urlParam['sub'] === 'createIC') {
      return <CreateICPanel />;
    }
  }
}
