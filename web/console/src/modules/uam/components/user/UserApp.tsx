import * as React from 'react';
import { connect } from 'react-redux';

import { Card, ContentView, Icon, Justify } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';

import { allActions } from '../../actions';
import { RootState } from '../../models';
import { router } from '../../router';
import { UserActionPanel } from './UserActionPanel';
import { UserDetailsPanel } from './UserDetailsPanel';
import { UserHeadPanel } from './UserHeadPanel';
import { UserTablePanel } from './UserTablePanel';

export interface RootProps extends RootState {
  actions?: typeof allActions;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class UserApp extends React.Component<RootProps, {}> {
  componentDidMount() {
    const { actions } = this.props;
    actions.user.poll();
  }
  render() {
    let { route } = this.props;
    let urlParam = router.resolve(route);
    const { sub } = urlParam;
    return (
      <React.Fragment>
        {sub ? (
          <ContentView>
            <ContentView.Header>
              <UserHeadPanel />
            </ContentView.Header>
            <ContentView.Body>
              <UserDetailsPanel />
            </ContentView.Body>
          </ContentView>
        ) : (
          <ContentView>
            <ContentView.Header>
              <UserHeadPanel />
            </ContentView.Header>
            <ContentView.Body>
              <UserActionPanel />
              <UserTablePanel />
            </ContentView.Body>
          </ContentView>
        )}
      </React.Fragment>
    );
  }
}
