import * as React from 'react';
import { connect, Provider } from 'react-redux';
import { bindActionCreators } from '@tencent/qcloud-lib';
import { RootState } from '../models';
import { allActions } from '../actions';
import { configStore } from '../stores/RootStore';
import { router } from '../router';
import { ResetStoreAction } from '../../../../helpers';
import { resourceConfig } from '../../../../config';
import { PersistentEventHeadPanel } from './PersistentEventHeadPanel';
import { ClusterTablePanel } from './ClusterTablePanel';
import { EditPersistentEventPanel } from './EditPersistentEventPanel';
import { ClusterActionPanel } from './ClusterActionPanel';
import { ContentView } from '@tencent/tea-component';

const store = configStore();
export class PersistentEventAppContainer extends React.Component<any, any> {
  // 页面离开时，清空store
  componentWillUnmount() {
    store.dispatch({ type: ResetStoreAction });
  }

  render() {
    return (
      <Provider store={store}>
        <PersistentEventApp />
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

@connect(
  state => state,
  mapDispatchToProps
)
@((router.serve as any)())
class PersistentEventApp extends React.Component<RootProps, {}> {
  componentDidMount() {
    let { actions } = this.props;

    // 初始化 resourceInfo
    let peResourceInfo = resourceConfig()['pe'];
    actions.pe.initPeResourceInfo(peResourceInfo);
  }

  render() {
    let { route } = this.props,
      urlParams = router.resolve(route);

    if (!urlParams['mode']) {
      return (
        <ContentView>
          <ContentView.Header>
            <PersistentEventHeadPanel />
          </ContentView.Header>
          <ContentView.Body>
            <ClusterActionPanel />
            <ClusterTablePanel />
          </ContentView.Body>
        </ContentView>
      );
    } else if (urlParams['mode'] === 'create' || urlParams['mode'] === 'update') {
      return <EditPersistentEventPanel />;
    }
  }
}
