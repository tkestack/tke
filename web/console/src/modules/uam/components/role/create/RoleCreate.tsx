import * as React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { HeaderPanel } from './HeaderPanel';
import { BaseInfoPanel } from './BaseInfoPanel';
import { RootProps } from '../RoleApp';
import { ContentView, Card, Justify, Icon } from '@tea/component';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class RoleCreate extends React.Component<RootProps, {}> {

  componentWillUnmount() {
    let { actions } = this.props;
    actions.role.create.addRoleWorkflow.reset();
    actions.role.create.clearCreationState();
    actions.role.create.clearValidatorState();
    /** 清理关联状态 */
    actions.commonUser.associate.clearUserAssociation();
    actions.policy.associate.clearPolicyAssociation();
    actions.group.associate.clearGroupAssociation();
  }

  componentDidMount() {
    const { actions } = this.props;
    /** 拉取用户列表 */
    actions.commonUser.associate.userList.performSearch('');
    /** 拉取用户组列表 */
    actions.group.associate.groupList.performSearch('');
    /** 拉取策略列表 */
    actions.policy.associate.policyList.performSearch('');
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
