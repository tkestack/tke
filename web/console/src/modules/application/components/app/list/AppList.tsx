import * as React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { HeaderPanel } from './HeaderPanel';
import { ActionPanel } from './ActionPanel';
import { TablePanel } from './TablePanel';
import { RootProps } from '../AppContainer';
import { ContentView, Card, Justify, Icon } from '@tea/component';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class AppList extends React.Component<RootProps, {}> {
  componentWillUnmount() {
    const { actions } = this.props;
    /** 取消轮询 */
    actions.app.list.clearPolling();
  }

  componentDidMount() {
    // const { actions } = this.props;
    /** 拉取应用列表 */
    // actions.app.list.poll();
    //不要保存filter旧数据
    // actions.cluster.list.reset();
    // actions.cluster.list.applyFilter({
    //   callback: {
    //     call: (cluster: string, namespace: string): void => {
    //       actions.app.list.poll({
    //         cluster: cluster,
    //         namespace: namespace
    //       });
    //     }
    //   }
    // });
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
