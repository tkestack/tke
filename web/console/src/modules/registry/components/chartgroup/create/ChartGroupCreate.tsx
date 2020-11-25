import * as React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { HeaderPanel } from './HeaderPanel';
import { BaseInfoPanel } from './BaseInfoPanel';
import { RootProps } from '../ChartGroupApp';
import { ContentView, Card, Justify, Icon } from '@tea/component';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ChartGroupCreate extends React.Component<RootProps, {}> {
  componentWillUnmount() {
    let { actions } = this.props;
    actions.chartGroup.create.addChartGroupWorkflow.reset();
    actions.chartGroup.create.clearCreationState();
    actions.chartGroup.create.clearValidatorState();

    actions.user.associate.clearUserAssociation();
  }

  componentDidMount() {
    const { actions } = this.props;
    /** 拉取业务列表 */
    actions.project.list.fetch();
    /** 拉取用户信息 */
    actions.user.detail.fetchUserInfo();
    /** 拉取用户列表 */
    actions.user.associate.userList.performSearch('');
  }

  render() {
    return (
      <React.Fragment>
        <ContentView>
          <ContentView.Header>
            <HeaderPanel />
          </ContentView.Header>
          <ContentView.Body>
            <BaseInfoPanel />
          </ContentView.Body>
        </ContentView>
      </React.Fragment>
    );
  }
}
