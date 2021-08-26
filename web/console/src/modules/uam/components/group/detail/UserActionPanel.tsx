/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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
import { connect } from 'react-redux';
import { Justify, Icon, Table, Button, SearchBox } from '@tea/component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators } from '@tencent/ff-redux';
import { router } from '../../../router';
import { allActions } from '../../../actions';
import { RootProps } from '../GroupPanel';
import { UserAssociateWorkflowDialog } from '../associate/UserAssociateWorkflowDialog';
import { CommonUserFilter } from '../../../models';

const mapDispatchToProps = (dispatch) =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect((state) => state, mapDispatchToProps)
export class UserActionPanel extends React.Component<RootProps, {}> {
  componentWillUnmount() {
    let { actions } = this.props;
    actions.commonUser.associate.clearUserAssociation();
  }

  componentDidMount() {
    const { actions, route } = this.props;
    /** 设置用户关联场景 */
    let filter: CommonUserFilter = {
      resource: 'localgroup',
      resourceID: route.queries['groupName'],
      /** 关联/解关联回调函数 */
      callback: () => {
        actions.group.detail.fetchGroup({ name: route.queries['groupName'] });
      },
    };
    actions.commonUser.associate.setupUserFilter(filter);
    /** 拉取关联用户列表，拉取后自动更新commonUserAssociation */
    actions.commonUser.associate.userAssociatedList.applyFilter(filter);
    /** 拉取用户列表 */
    actions.commonUser.associate.userList.performSearch('');
  }

  render() {
    const { actions, route } = this.props;
    let urlParam = router.resolve(route);

    return (
      <React.Fragment>
        <Table.ActionPanel>
          <Justify
            left={
              <Button
                type="primary"
                onClick={(e) => {
                  /** 开始关联用户工作流 */
                  actions.commonUser.associate.associateUserWorkflow.start();
                }}
              >
                {t('关联用户')}
              </Button>
            }
          />
        </Table.ActionPanel>
        <UserAssociateWorkflowDialog
          onPostCancel={() => {
            /** 不清理commonUserAssociation */
            // actions.commonUser.associate.clearUserAssociation();
          }}
        />
      </React.Fragment>
    );
  }
}
