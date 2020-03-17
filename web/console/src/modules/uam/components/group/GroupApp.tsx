import * as React from 'react';
import { connect, Provider } from 'react-redux';
import { bindActionCreators } from '@tencent/ff-redux';
import { RootState } from '../../models';
import { allActions } from '../../actions';
import { router } from '../../router';
import { GroupList } from './list/GroupList';
import { GroupCreate } from './create/GroupCreate';
import { GroupDetail } from './detail/GroupDetail';

export interface RootProps extends RootState {
  actions?: typeof allActions;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class GroupApp extends React.Component<RootProps, {}> {

  render() {
    let { route } = this.props,
      urlParam = router.resolve(route);
    if (!urlParam['sub']) {
      return (
        <div className="manage-area">
          <GroupList {...this.props} />
        </div>
      );
    } else if (urlParam['sub'] === 'create') {
      return (
        <div className="manage-area">
          <GroupCreate {...this.props} />
        </div>
      );
    } else if (urlParam['sub'] === 'detail') {
      return (
        <div className="manage-area">
          <GroupDetail {...this.props} />
        </div>
      );
    }
  }
}
