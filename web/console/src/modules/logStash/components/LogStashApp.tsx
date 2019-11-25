import * as React from 'react';
import { connect, Provider } from 'react-redux';
import { bindActionCreators } from '@tencent/qcloud-lib';
import { RootState } from '../models';
import { allActions } from '../actions';
import { configStore } from '../stores/RootStore';
import { router } from '../router';
import { ResetStoreAction } from '../../../../helpers';
import { LogStashHeadPanel } from './LogStashHeadPanel';
import { LogStashActionPanel } from './LogStashActionPanel';
import { LogStashTablePanel } from './LogStashTablePanel';
import { OpenLogStashDialog } from './OpenLogStashDialog';
import { LogStashSubHeadPanel } from './LogStashSubHeadPanel';
import { EditLogStashPanel } from './EditLogStashPanel';
import { ContentView } from '@tencent/tea-component';
import { LogStashDetailPanel } from './LogDetailPanel';
import { DeleteLogDialog } from './DeleteLogDialog';

const store = configStore();

export class LogStasgAppContainer extends React.Component<any, any> {
  // 页面离开时，清空store
  componentWillUnmount() {
    store.dispatch({ type: ResetStoreAction });
  }

  render() {
    return (
      <Provider store={store}>
        <LogStashApp />
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

@connect(
  state => state,
  mapDispatchToProps
)
@((router.serve as any)())
class LogStashApp extends React.Component<RootProps, any> {
  render() {
    let { actions, route } = this.props;
    const urlParams = router.resolve(route);

    let content: JSX.Element;
    let { mode } = urlParams;

    if (!mode) {
      content = (
        <React.Fragment>
          <ContentView>
            <ContentView.Header>
              <LogStashHeadPanel />
            </ContentView.Header>
            <ContentView.Body>
              <LogStashActionPanel />
              <LogStashTablePanel />
              <DeleteLogDialog />
            </ContentView.Body>
          </ContentView>
          <OpenLogStashDialog />
        </React.Fragment>
      );
    } else if (mode === 'create') {
      content = (
        <React.Fragment>
          <ContentView>
            <ContentView.Header>
              <LogStashSubHeadPanel />
            </ContentView.Header>
            <ContentView.Body>
              <EditLogStashPanel />
            </ContentView.Body>
          </ContentView>
          <OpenLogStashDialog />
        </React.Fragment>
      );
    } else if (mode === 'update') {
      content = (
        <React.Fragment>
          <ContentView.Header>
            <LogStashSubHeadPanel />
          </ContentView.Header>
          <ContentView.Body>
            <EditLogStashPanel />
          </ContentView.Body>
        </React.Fragment>
      );
    } else if (mode === 'detail') {
      content = (
        <ContentView>
          <ContentView.Header>
            <LogStashSubHeadPanel />
          </ContentView.Header>
          <ContentView.Body>
            <LogStashDetailPanel />
          </ContentView.Body>
        </ContentView>
      );
    }

    return content;
  }
}
