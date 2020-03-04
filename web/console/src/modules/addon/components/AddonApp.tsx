import * as React from 'react';
import { connect, Provider } from 'react-redux';

import { bindActionCreators } from '@tencent/ff-redux';
import { ContentView } from '@tencent/tea-component';

import { ResetStoreAction } from '../../../../helpers';
import { allActions } from '../actions';
import { RootState } from '../models';
import { router } from '../router';
import { configStore } from '../stores/RootStore';
import { AddonActionPanel } from './AddonActionPanel';
import { AddonDeleteDialog } from './AddonDeleteDialog';
import { AddonDetail } from './AddonDetail';
import { AddonHeadPanel } from './AddonHeadPanel';
import { AddonSubpageHeadPanel } from './AddonSubpageHeadPanel';
import { AddonTablePanel } from './AddonTablePanel';
import { EditAddonPanel } from './EditAddonPanel';

export const store = configStore();

export class AddonAppContainer extends React.Component<any, any> {
  // 页面离开时，清空store
  componentWillUnmount() {
    store.dispatch({ type: ResetStoreAction });
  }
  render() {
    return (
      <Provider store={store}>
        <AddonApp />
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
class AddonApp extends React.Component<RootProps, any> {
  render() {
    let { route } = this.props,
      urlParams = router.resolve(route);

    let mode = urlParams['mode'];
    if (!mode) {
      return (
        <React.Fragment>
          <ContentView>
            <ContentView.Header>
              <AddonHeadPanel />
            </ContentView.Header>
            <ContentView.Body>
              <AddonActionPanel />
              <AddonTablePanel />
            </ContentView.Body>
          </ContentView>
          <AddonDeleteDialog />
        </React.Fragment>
      );
    } else if (mode === 'detail') {
      return <AddonDetail />;
    } else if (mode === 'create') {
      return (
        <ContentView>
          <ContentView.Header>
            <AddonSubpageHeadPanel route={route} />
          </ContentView.Header>
          <ContentView.Body>
            <EditAddonPanel />
          </ContentView.Body>
        </ContentView>
      );
    }
  }
}
