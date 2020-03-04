import * as React from 'react';
import { connect, Provider } from 'react-redux';

import { bindActionCreators } from '@tencent/ff-redux';
import { ContentView, Icon } from '@tencent/tea-component';

import { notifySubRouter } from '../../../../config/routerConfig';
import { ResetStoreAction } from '../../../../helpers';
import { allActions } from '../actions';
import * as ActionType from '../constants/ActionType';
import { RootState } from '../models';
import { router } from '../router';
import { configStore } from '../stores/RootStore';
import { DeleteResourceDialog } from './DeleteResourceDialog';
import { NotifyHead } from './NotifyHead';
import { ResourceDetail } from './resourceDetail/ResourceDetail';
import { ResourceDetailChannel } from './resourceDetail/ResourceDetailChannel';
import { ResourceDetailReceiver } from './resourceDetail/ResourceDetailReceiver';
import { ResourceDetailReceiverGroup } from './resourceDetail/ResourceDetailReceiverGroup';
import { ResourceDetailTempalte } from './resourceDetail/ResourceDetailTemplate';
import { ResourceHeader } from './resourceDetail/ResourceHeader';
import { EditResourceChannel } from './resourceEdition/EditResourceChannel';
import { EditResourceReceiver } from './resourceEdition/EditResourceReceiver';
import { EditResourceReceiverGroup } from './resourceEdition/EditResourceReceiverGroup';
import { EditResourceTemplate } from './resourceEdition/EditResourceTemplate';
import { ResourceTable } from './resourceList/ResourceTable';
import { ResourceTableChannel } from './resourceList/ResourceTableChannel';
import { ResourceTableReceiver } from './resourceList/ResourceTableReceiver';
import { ResourceTableTemplate } from './resourceList/ResourceTableTemplate';
import { ResourceSidebar } from './ResourceSidebar';

const store = configStore();

export class NotifyAppContainer extends React.Component<any, any> {
  // 页面离开时，清空store
  componentWillUnmount() {
    store.dispatch({ type: ResetStoreAction });
  }
  render() {
    return (
      <Provider store={store}>
        <NotifyApp />
      </Provider>
    );
  }
}

export interface RootProps extends RootState {
  actions?: typeof allActions;
}

const mapDispatchToProps = dispatch =>
  Object.assign(
    {},
    bindActionCreators(
      {
        actions: allActions
      },
      dispatch
    ),
    { dispatch }
  );

@connect(state => state, mapDispatchToProps)
@((router.serve as any)())
class NotifyApp extends React.Component<RootProps, any> {
  componentDidMount() {
    let resourceName = router.resolve(this.props.route).resourceName || 'channel';
    this.props.actions.resource[resourceName].fetch();
  }

  componentWillReceiveProps(nextProps: RootProps) {
    let resourceName = router.resolve(this.props.route).resourceName || 'channel';
    let nextResourceName = router.resolve(nextProps.route).resourceName || 'channel';
    if (resourceName !== nextResourceName) {
      this.props.actions.resource[resourceName].clear();
      nextProps.actions.resource[nextResourceName].clear();
      nextProps.actions.resource[nextResourceName].polling({
        filter: {},
        delayTime: 1e4
      });
    }
  }

  componentWillUnmount() {
    let { actions } = this.props;
    for (let resourceName in actions.resource) {
      actions.resource[resourceName].clearPolling && actions.resource[resourceName].clearPolling();
    }
  }

  render() {
    let { route } = this.props;
    const urlParams = router.resolve(route);
    let sideBar = <ResourceSidebar {...this.props} subRouterList={notifySubRouter} />;
    let content;
    let header;
    switch (urlParams.mode) {
      case 'detail':
        header = <ResourceHeader {...this.props} />;

        switch (urlParams.resourceName) {
          case 'channel':
            content = <ResourceDetailChannel {...this.props} />;
            break;
          case 'template':
            content = <ResourceDetailTempalte {...this.props} />;
            break;
          case 'receiver':
            content = <ResourceDetailReceiver {...this.props} />;
            break;
          case 'receiverGroup':
            content = <ResourceDetailReceiverGroup {...this.props} />;
            break;
          default:
            content = <ResourceDetail {...this.props} />;
        }
        break;

      case 'create':
      case 'copy':
      case 'update': {
        let resourceName = urlParams.resourceName || 'channel';
        header = <ResourceHeader {...this.props} />;

        let ins = this.props[resourceName].list.data.records.find(
          ins => ins.metadata.name === route.queries.resourceIns
        );
        if (!ins) {
          content = <Icon type="loading" />;
        }
        switch (resourceName) {
          case 'channel':
            content = <EditResourceChannel {...this.props} instance={ins} />;
            break;
          case 'template':
            content = <EditResourceTemplate {...this.props} instance={ins} />;
            break;
          case 'receiver':
            content = <EditResourceReceiver {...this.props} instance={ins} />;
            break;
          case 'receiverGroup':
            content = <EditResourceReceiverGroup {...this.props} instance={ins} />;
            break;
        }
        break;
      }
      default:
        header = <NotifyHead {...this.props} />;
        switch (urlParams.resourceName) {
          case 'channel':
            content = <ResourceTableChannel {...this.props} />;
            break;
          case 'template':
            content = <ResourceTableTemplate {...this.props} />;
            break;
          case 'receiver':
            content = <ResourceTableReceiver {...this.props} />;
            break;
          default:
            content = <ResourceTableChannel {...this.props} />;
        }
    }
    return (
      <ContentView>
        <ContentView.Header>{header}</ContentView.Header>
        <ContentView.Body sidebar={sideBar}>
          <ContentView>
            <ContentView.Body>
              {content}
              <DeleteResourceDialog {...this.props} />
            </ContentView.Body>
          </ContentView>
        </ContentView.Body>
      </ContentView>
    );
  }
}
