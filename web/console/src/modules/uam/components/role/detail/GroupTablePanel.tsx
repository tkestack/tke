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
import { GroupPlain, GroupAssociation } from '../../../models';
import { RootProps } from '../RoleApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class GroupTablePanel extends React.Component<RootProps, {}> {

  render() {
    let { actions, groupAssociation, groupAssociatedList } = this.props;

    const columns: TableColumn<GroupPlain>[] = [
      {
        key: 'name',
        header: t('名称'),
        render: (group, text, index) => (
          <Text parent="div" overflow>
            {group.displayName || '-'}
          </Text>
        )
      },
      {
        key: 'description',
        header: t('描述'),
        render: (group, text, index) => (
          <Text parent="div" overflow>
            {group.description || '-'}
          </Text>
        )
      },
      { key: 'operation', header: t('操作'), render: group => this._renderOperationCell(group) }
    ];

    return (
      <TablePanel
        columns={columns}
        recordKey={'id'}
        records={groupAssociation.originGroups}
        action={actions.group.associate.groupAssociatedList}
        model={groupAssociatedList}
        emptyTips={emptyTips}
      />
    );
  }

  /** 渲染操作按钮 */
  _renderOperationCell = (group: GroupPlain) => {
    let { actions } = this.props;
    return (
      <React.Fragment>
        <LinkButton
          tipDirection="right"
          onClick={(e) => {
            this._removeGroup(group);
          }}
        >
          <Trans>解除关联</Trans>
        </LinkButton>
      </React.Fragment>
    );
  }

  _removeGroup = async (group: GroupPlain) => {
    let { actions, groupFilter } = this.props;
    const yes = await Modal.confirm({
      message: t('确认解除当前用户组关联') + ` - ${group.displayName}？`,
      okText: t('解除'),
      cancelText: t('取消')
    });
    if (yes) {
      let groupAssociation: GroupAssociation = { id: uuid(), removeGroups: [group] };
      actions.group.associate.disassociateGroupWorkflow.start([groupAssociation], groupFilter);
      actions.group.associate.disassociateGroupWorkflow.perform();
    }
  }

}
