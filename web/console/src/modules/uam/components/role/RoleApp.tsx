import * as React from 'react';
import { connect, Provider } from 'react-redux';
import { bindActionCreators } from '@tencent/ff-redux';
import { RootState } from '../../models';
import { allActions } from '../../actions';
import { router } from '../../router';
import { RoleList } from './list/RoleList';
import { RoleCreate } from './create/RoleCreate';
import { RoleDetail } from './detail/RoleDetail';

export interface RootProps extends RootState {
  actions?: typeof allActions;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class RoleApp extends React.Component<RootProps, {}> {

  render() {
    let { route } = this.props,
      urlParam = router.resolve(route);
    if (!urlParam['sub']) {
      return (
        <div className="manage-area">
          <RoleList {...this.props} />
        </div>
      );
    } else if (urlParam['sub'] === 'create') {
      return (
        <div className="manage-area">
          <RoleCreate {...this.props} />
        </div>
      );
    } else if (urlParam['sub'] === 'detail') {
      return (
        <div className="manage-area">
          <RoleDetail {...this.props} />
        </div>
      );
    }
  }
}
