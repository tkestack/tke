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
import { GroupList } from './list/GroupList';
import { GroupCreate } from './create/GroupCreate';
import { GroupDetail } from './detail/GroupDetail';

export interface RootProps extends RootState {
  actions?: typeof allActions;
}

const mapDispatchToProps = (dispatch) =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch,
  });

@connect((state) => state, mapDispatchToProps)
export class GroupPanel extends React.Component<RootProps, {}> {
  render() {
    const { route } = this.props;
    const { action } = router.resolve(route);
    if (action === 'create') {
      return (
        <div className="manage-area">
          <GroupCreate {...this.props} />
        </div>
      );
    } else if (action === 'detail') {
      return (
        <div className="manage-area">
          <GroupDetail {...this.props} />
        </div>
      );
    } else {
      return (
        <div className="manage-area">
          <GroupList {...this.props} />
        </div>
      );
    }
  }
}
