import * as React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { HeaderPanel } from './HeaderPanel';
import { BaseInfoPanel } from './BaseInfoPanel';
import { RootProps } from '../AppContainer';
import { ContentView, Card, Justify, Icon } from '@tea/component';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class AppCreate extends React.Component<RootProps, {}> {
  componentWillUnmount() {
    let { actions } = this.props;
    actions.app.create.addAppWorkflow.reset();
    actions.app.create.clearCreationState();
    actions.app.create.clearValidatorState();
  }

  componentDidMount() {
    const { actions, appCreation } = this.props;
    /** 拉取仓库列表 */
    actions.chartGroup.list.applyFilter({});
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
