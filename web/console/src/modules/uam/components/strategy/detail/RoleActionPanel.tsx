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
import { connect } from 'react-redux';
import { Justify, Icon, Table, Button, SearchBox } from '@tea/component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators } from '@tencent/ff-redux';
import { router } from '../../../router';
import { allActions } from '../../../actions';
import { RootProps } from '../StrategyApp';
// import { RoleAssociateWorkflowDialog } from '../associate/RoleAssociateWorkflowDialog';
import { RoleFilter } from '../../../models';

const mapDispatchToProps = (dispatch) =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect((state) => state, mapDispatchToProps)
export class RoleActionPanel extends React.Component<RootProps, {}> {
  componentWillUnmount() {
    let { actions } = this.props;
    actions.role.associate.clearRoleAssociation();
  }

  componentDidMount() {
    const { actions, route } = this.props;
    /** 设置角色关联场景 */
    let filter: RoleFilter = {
      resource: 'policy',
      // resourceID: route.queries['groupName']
      // resourceID: router.resolve(route).sub,
      resourceID: route.queries['id'],
      /** 关联/解关联回调函数 */
      callback: () => {
        /** 重新加载策略 */
      },
    };
    actions.role.associate.setupRoleFilter(filter);
    /** 拉取关联角色列表，拉取后自动更新roleAssociation */
    actions.role.associate.roleAssociatedList.applyFilter(filter);
    /** 拉取角色列表 */
    // actions.role.associate.roleList.performSearch('');
  }

  render() {
    const { actions, route } = this.props;

    return (
      <noscript />
      // <React.Fragment>
      //   <Table.ActionPanel>
      //     <Justify
      //       left={
      //         <Button type="primary" onClick={e => {
      //           /** 开始关联角色工作流 */
      //           actions.role.associate.associateRoleWorkflow.start();
      //         }}>
      //           {t('关联角色')}
      //         </Button>
      //       }
      //     />
      //   </Table.ActionPanel>
      //   <RoleAssociateWorkflowDialog />
      // </React.Fragment>
    );
  }
}
