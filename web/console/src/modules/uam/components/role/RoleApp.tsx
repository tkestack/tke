/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */
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
