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
import { RootProps } from '../StrategyApp';
import { GroupAssociateWorkflowDialog } from '../associate/GroupAssociateWorkflowDialog';
import { GroupFilter } from '../../../models';

const mapDispatchToProps = (dispatch) =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect((state) => state, mapDispatchToProps)
export class GroupActionPanel extends React.Component<RootProps, {}> {
  componentWillUnmount() {
    let { actions } = this.props;
    actions.group.associate.clearGroupAssociation();
  }

  componentDidMount() {
    const { actions, route } = this.props;
    /** 设置用户组关联场景 */
    let filter: GroupFilter = {
      resource: 'policy',
      // resourceID: route.queries['roleName']
      // resourceID: router.resolve(route).sub,
      resourceID: route.queries['id'],
      /** 关联/解关联回调函数 */
      callback: () => {
        /** 重新加载策略 */
      },
    };
    actions.group.associate.setupGroupFilter(filter);
    /** 拉取关联用户组列表，拉取后自动更新groupAssociation */
    actions.group.associate.groupAssociatedList.applyFilter(filter);
    /** 拉取用户组列表 */
    actions.group.associate.groupList.performSearch('');
  }

  render() {
    const { actions, route } = this.props;
    let urlParam = router.resolve(route);
    const { sub } = urlParam;

    return (
      <React.Fragment>
        <Table.ActionPanel>
          <Justify
            left={
              <Button
                type="primary"
                onClick={(e) => {
                  /** 开始关联用户组工作流 */
                  actions.group.associate.associateGroupWorkflow.start();
                }}
              >
                {t('关联用户组')}
              </Button>
            }
          />
        </Table.ActionPanel>
        <GroupAssociateWorkflowDialog />
      </React.Fragment>
    );
  }
}
