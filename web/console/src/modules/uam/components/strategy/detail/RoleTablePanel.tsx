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
import { TablePanel } from '@tencent/ff-component';
import { LinkButton, emptyTips } from '../../../../common/components';
import { Table, TableColumn, Text, Modal, Icon } from '@tea/component';
import { bindActionCreators, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { router } from '../../../router';
import { allActions } from '../../../actions';
import { RolePlain, RoleAssociation, GroupAssociation, GroupPlain, GroupFilter } from '../../../models';
import { RootProps } from '../StrategyApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class RoleTablePanel extends React.Component<RootProps, {}> {

  render() {
    let { actions, roleAssociation, roleAssociatedList } = this.props;

    const columns: TableColumn<RolePlain>[] = [
      {
        key: 'name',
        header: t('名称'),
        render: (role, text, index) => (
          <Text parent="div" overflow>
            {role.displayName || '-'}
          </Text>
        )
      },
      {
        key: 'description',
        header: t('描述'),
        render: (role, text, index) => (
          <Text parent="div" overflow>
            {role.description || '-'}
          </Text>
        )
      },
      // { key: 'operation', header: t('操作'), render: role => this._renderOperationCell(role) }
    ];

    return (
      <React.Fragment>
        <TablePanel
          columns={columns}
          recordKey={'id'}
          records={roleAssociation.originRoles}
          action={actions.role.associate.roleAssociatedList}
          model={roleAssociatedList}
          emptyTips={emptyTips}
        />
      </React.Fragment>
    );
  }

  /** 渲染操作按钮 */
  _renderOperationCell = (role: RolePlain) => {
    let { actions } = this.props;
    return (
      <React.Fragment>
        <LinkButton
          tipDirection="right"
          onClick={(e) => {
            this._removeRole(role);
          }}
        >
          <Trans>解除关联</Trans>
        </LinkButton>
      </React.Fragment>
    );
  }

  _removeRole = async (role: RolePlain) => {
    let { actions, roleFilter } = this.props;
    const yes = await Modal.confirm({
      message: t('确认解除当前角色关联') + ` - ${role.displayName}？`,
      okText: t('解除'),
      cancelText: t('取消')
    });
    if (yes) {
      /** 目前还没有实现基于用户组解绑角色 */
      let roleAssociation: RoleAssociation = { id: uuid(), removeRoles: [role] };
      actions.role.associate.disassociateRoleWorkflow.start([roleAssociation], roleFilter);
      actions.role.associate.disassociateRoleWorkflow.perform();
    }
  }

}
