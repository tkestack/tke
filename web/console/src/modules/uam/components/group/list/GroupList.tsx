import * as React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { HeaderPanel } from './HeaderPanel';
import { ActionPanel } from './ActionPanel';
import { TablePanel } from './TablePanel';
import { RootProps } from '../GroupApp';
import { ContentView, Card, Justify, Icon } from '@tea/component';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class GroupList extends React.Component<RootProps, {}> {

  componentWillUnmount() {
    const { actions } = this.props;
    /** 取消轮询 */
    actions.group.list.clearPolling();
  }

  componentDidMount() {
    const { actions } = this.props;
    /** 拉取用户组列表 */
    actions.group.list.poll();
  }

  render() {
    return (
      <React.Fragment>
        <ContentView>
          <ContentView.Header>
            <HeaderPanel />
          </ContentView.Header>
          <ContentView.Body>
            <ActionPanel />
            <TablePanel />
          </ContentView.Body>
        </ContentView>
      </React.Fragment>
    );
  }
}
