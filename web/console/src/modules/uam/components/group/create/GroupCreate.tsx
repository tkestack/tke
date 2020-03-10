import * as React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { HeaderPanel } from './HeaderPanel';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { BaseInfoPanel } from './BaseInfoPanel';
import { RootProps } from '../GroupApp';
import { ContentView, Card, Justify, Icon } from '@tea/component';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class GroupCreate extends React.Component<RootProps, {}> {

  componentWillUnmount() {
    let { actions } = this.props;
    actions.group.create.addGroupWorkflow.reset();
    actions.group.create.clearCreationState();
    actions.group.create.clearValidatorState();
    /** 清理关联状态 */
    actions.commonUser.associate.clearUserAssociation();
  }

  componentDidMount() {
    const { actions } = this.props;
    /** 拉取用户列表 */
    actions.commonUser.associate.userList.performSearch('');
  }

  render() {
    return (
      <ContentView>
        <ContentView.Header>
          <HeaderPanel />
        </ContentView.Header>
        <ContentView.Body>
          <BaseInfoPanel />
        </ContentView.Body>
      </ContentView>
    );
  }
}
