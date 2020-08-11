import * as React from 'react';
import { connect, Provider } from 'react-redux';
import { bindActionCreators } from '@tencent/ff-redux';
import { RootState } from '../../models';
import { allActions } from '../../actions';
import { router } from '../../router';
import { ChartGroupList } from './list/ChartGroupList';
import { ChartGroupCreate } from './create/ChartGroupCreate';
import { ChartGroupDetail } from './detail/ChartGroupDetail';

export interface RootProps extends RootState {
  actions?: typeof allActions;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class ChartGroupApp extends React.Component<RootProps, {}> {

  render() {
    let { route } = this.props,
      urlParam = router.resolve(route);
    if (!urlParam['mode'] || urlParam['mode'] === 'list') {
      return (
        <div className="manage-area">
          <ChartGroupList {...this.props} />
        </div>
      );
    } else if (urlParam['mode'] === 'create') {
      return (
        <div className="manage-area">
          <ChartGroupCreate {...this.props} />
        </div>
      );
    } else if (urlParam['mode'] === 'detail') {
      return (
        <div className="manage-area">
          <ChartGroupDetail {...this.props} />
        </div>
      );
    }
  }
}
