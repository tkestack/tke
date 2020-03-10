import * as React from 'react';
import { connect, Provider } from 'react-redux';

import { bindActionCreators } from '@tencent/ff-redux';

import { ResetStoreAction } from '../../../../helpers';
import { allActions } from '../actions';
import { RootState } from '../models';
import { router } from '../router';
import { configStore } from '../stores/RootStore';
import { StrategyApp } from './strategy/StrategyApp';
import { UserApp } from './user/UserApp';
import { RoleApp } from './role/RoleApp';
import { GroupApp } from './group/GroupApp';

const store = configStore();

export class UamAppContainer extends React.Component<any, any> {
  // 页面离开时，清空store
  componentWillUnmount() {
    store.dispatch({ type: ResetStoreAction });
  }
  render() {
    return (
      <Provider store={store}>
        <UamApp />
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
class UamApp extends React.Component<RootProps, {}> {
  render() {
    let { route } = this.props;
    let urlParam = router.resolve(route);
    const { module } = urlParam;

    let content: React.ReactNode;
    if (module === 'user') {
      content = <UserApp />;
    } else if (module === 'strategy') {
      content = <StrategyApp />;
    } else if (module === 'role') {
      content = <RoleApp />;
    } else if (module === 'group') {
      content = <GroupApp />;
    } else {
      content = <UserApp />;
    }

    return content;
  }
}
