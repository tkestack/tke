import * as React from 'react';
import { connect, Provider } from 'react-redux';

import { bindActionCreators } from '@tencent/ff-redux';
import { ContentView, Row, Col } from '@tencent/tea-component';

import { ResetStoreAction } from '../../../../helpers';
import { overviewActions } from '../actions/overviewActions';
import { router } from '../router';
import { configStore } from '../stores/RootStore';
import { OverviewHeadPanel } from './OverviewHeadPanel';
import { RootState } from '../models/RootState';
import { ClusterOverviewPanel } from './ClusterOverview';
import { QuickHelpPanel } from './QuickHelpPanel';
import { TipsPanel } from './TipsPanel';
import { ClusterDetailPanel } from './ClusterDetailPanel';
const store = configStore();

export class OverviewAppContainer extends React.Component<any, any> {
  // 页面离开时，清空store
  componentWillUnmount() {
    store.dispatch({ type: ResetStoreAction });
  }

  render() {
    return (
      <Provider store={store}>
        <OverviewApp />
      </Provider>
    );
  }
}

export interface RootProps extends RootState {
  actions?: typeof overviewActions;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: overviewActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
@((router.serve as any)())
class OverviewApp extends React.Component<RootProps, {}> {
  render() {
    let { clusterOverview } = this.props;
    return (
      <ContentView>
        <ContentView.Header>
          <OverviewHeadPanel {...this.props} />
        </ContentView.Header>
        <ContentView.Body>
          <Row>
            <Col span={18}>
              <ClusterOverviewPanel
                clusterData={clusterOverview.object && clusterOverview.object.data ? clusterOverview.object.data : null}
              />
              <ClusterDetailPanel
                clusterData={clusterOverview.object && clusterOverview.object.data ? clusterOverview.object.data : null}
              />
            </Col>
            <Col span={6}>
              <QuickHelpPanel />
              <TipsPanel />
            </Col>
          </Row>
        </ContentView.Body>
      </ContentView>
    );
  }
}
