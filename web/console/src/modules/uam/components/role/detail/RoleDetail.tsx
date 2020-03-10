import * as React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { RootProps } from '../RoleApp';
import { HeaderPanel } from './HeaderPanel';
import { BaseInfoPanel } from './BaseInfoPanel';
import { ContentView } from '@tea/component';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class RoleDetail extends React.Component<RootProps, {}> {

  componentWillUnmount() {
    let { actions } = this.props;
    actions.role.detail.updateRoleWorkflow.reset();
    actions.role.detail.clearEditorState();
    actions.role.detail.clearValidatorState();
  }

  componentDidMount() {
    const { actions, route } = this.props;
    /** 查询具体角色，从而Detail可以用到 */
    actions.role.detail.fetchRole({ name: route.queries['roleName'] });
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
