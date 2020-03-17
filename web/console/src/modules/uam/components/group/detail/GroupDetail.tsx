import * as React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { RootProps } from '../GroupApp';
import { HeaderPanel } from './HeaderPanel';
import { BaseInfoPanel } from './BaseInfoPanel';
import { ContentView } from '@tea/component';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class GroupDetail extends React.Component<RootProps, {}> {

  componentWillUnmount() {
    let { actions } = this.props;
    actions.group.detail.updateGroupWorkflow.reset();
    actions.group.detail.clearEditorState();
    actions.group.detail.clearValidatorState();
  }

  componentDidMount() {
    const { actions, route } = this.props;
    /** 查询具体用户组，从而Detail可以用到 */
    actions.group.detail.fetchGroup({ name: route.queries['groupName'] });
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
