import * as React from 'react';
import { connect, Provider } from 'react-redux';

import { bindActionCreators } from '@tencent/ff-redux';
import { Card, ContentView } from '@tencent/tea-component';

import { ResetStoreAction } from '../../../../helpers';
import { allActions } from '../actions';
import { RootState } from '../models';
import { router } from '../router';
import { configStore } from '../stores/RootStore';
import { ClusterActionPanel } from './clusterManage/ClusterActionPanel';
import { ClusterDeleteDialog } from './clusterManage/ClusterDeleteDialog';
import { ClusterHeadPanel } from './clusterManage/ClusterHeadPanel';
import { ClusterStatusDialog } from './clusterManage/ClusterStatusDialog';
import { ClusterTablePanel } from './clusterManage/ClusterTablePanel';
import { CreateClusterPanel } from './clusterManage/CreateClusterPanel';
import { CreateICPanel } from './clusterManage/CreateICPanel';
import { ModifyClusterNameDialog } from './clusterManage/ModifyClusterNameDialog';
import { TcrRegistyDeployDialog } from './clusterManage/TcrRegistyDeployDialog';
import { ResourceContainerPanel } from './resource/ResourceContainerPanel';
import { ConfigPromethus } from './clusterManage/ConfigPromethus';
import { RecoilRoot } from 'recoil';

export const store = configStore();

export class ClusterAppContainer extends React.Component<any, any> {
  // 页面离开时，清空store
  componentWillUnmount() {
    store.dispatch({ type: ResetStoreAction });
  }
  render() {
    return (
      <Provider store={store}>
        <RecoilRoot>
          <ClusterApp />
        </RecoilRoot>
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
    } else if (urlParam['sub'] === 'config-promethus') {
      return <ConfigPromethus {...this.props} />;
    }
  }
}
