import * as React from 'react';
import { connect, Provider } from 'react-redux';

import { bindActionCreators } from '@tencent/ff-redux';
import { ContentView } from '@tencent/tea-component';

import { ResetStoreAction } from '../../../../helpers';
import { allActions } from '../actions/';
import * as ActionType from '../constants/ActionType';
import { RootState } from '../models';
import { router } from '../router';
import { configStore } from '../stores/RootStore';
import { AlarmPolicyDetailHeaderPanel } from './AlarmPolicyDetailHeaderPanel';
import { AlarmPolicyDetailPanel } from './AlarmPolicyDetailPanel';
import { AlarmPolicyHeadPanel } from './AlarmPolicyHeadPanel';
import { AlarmPolicySubpageHeaderPanel } from './AlarmPolicySubpageHeaderPanel';
import { AlarmPolicyTablePanel } from './AlarmPolicyTablePanel';
import { DeleteAlarmPolicyDialog } from './DeleteAlarmPolicyDialog';
import { EditAlarmPolicyPanel } from './EditAlarmPolicyPanel';

const store = configStore();

export class AlarmPolicyAppContainer extends React.Component<any, any> {
  // 页面离开时，清空store
  componentWillUnmount() {
    store.dispatch({ type: ResetStoreAction });
  }
  render() {
    return (
      <Provider store={store}>
        <AlarmPolicyApp />
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
class AlarmPolicyApp extends React.Component<RootProps, any> {
  componentDidMount() {
    store.dispatch({ type: ActionType.isI18n, payload: false });
  }

  render() {
    let { actions, route } = this.props;

    const urlParams = router.resolve(route);
    if (!urlParams['sub']) {
      return (
        <ContentView>
          <ContentView.Header>
            <AlarmPolicyHeadPanel {...this.props} />
          </ContentView.Header>
          <ContentView.Body>
            <AlarmPolicyTablePanel {...this.props} />
            <DeleteAlarmPolicyDialog {...this.props} />
          </ContentView.Body>
        </ContentView>
      );
    } else if (urlParams['sub'] === 'create' || urlParams['sub'] === 'copy' || urlParams['sub'] === 'update') {
      return (
        <div className="manage-area server-add-box">
          <AlarmPolicySubpageHeaderPanel {...this.props} />
          <EditAlarmPolicyPanel {...this.props} />
        </div>
      );
    } else if (urlParams['sub'] === 'detail') {
      return (
        <ContentView>
          <ContentView.Header>
            <AlarmPolicyDetailHeaderPanel {...this.props} />
          </ContentView.Header>
          <ContentView.Body>
            <AlarmPolicyDetailPanel {...this.props} />
          </ContentView.Body>
        </ContentView>
      );
    }
  }
}
