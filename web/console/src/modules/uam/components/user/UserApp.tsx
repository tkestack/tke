import * as React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from '@tencent/qcloud-lib';
import { RootState } from '../../models';
import { allActions } from '../../actions';
import { router } from '../../router';
import { UserHeadPanel } from './UserHeadPanel';
import { UserActionPanel } from './UserActionPanel';
import { UserTablePanel } from './UserTablePanel';
import { UserDetailsPanel } from './UserDetailsPanel';
import { ContentView, Card, Justify, Icon } from '@tea/component';

export interface RootProps extends RootState {
  actions?: typeof allActions;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class UserApp extends React.Component<RootProps, {}> {
  componentDidMount() {
    const { actions } = this.props;
    actions.user.applyFilter({});
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
