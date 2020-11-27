import * as React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { RootProps } from '../ChartGroupApp';
import { HeaderPanel } from './HeaderPanel';
import { BaseInfoPanel } from './BaseInfoPanel';
import { ContentView } from '@tea/component';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ChartGroupDetail extends React.Component<RootProps, {}> {
  componentWillUnmount() {
    let { actions } = this.props;
    actions.chartGroup.detail.updateChartGroupWorkflow.reset();
    actions.chartGroup.detail.clearEditorState();
    actions.chartGroup.detail.clearValidatorState();

    actions.user.associate.clearUserAssociation();
  }

  componentDidMount() {
    const { actions, route } = this.props;
    /** 查询具体仓库，从而Detail可以用到 */
    actions.chartGroup.detail.fetchChartGroup({ name: route.queries['cg'], projectID: route.queries['prj'] });

    /** 获取具备权限的业务列表 */
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
