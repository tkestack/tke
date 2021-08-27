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
import { LinkButton, emptyTips } from '../../../../common/components';
import { TablePanel } from '@tencent/ff-component';
import { Table, TableColumn, Text, Modal, Icon } from '@tea/component';
import { bindActionCreators, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { router } from '../../../router';
import { allActions } from '../../../actions';
import { PolicyPlain, PolicyAssociation, GroupAssociation, GroupPlain, GroupFilter } from '../../../models';
import { RootProps } from '../GroupPanel';

const mapDispatchToProps = (dispatch) =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch,
  });

@connect((state) => state, mapDispatchToProps)
export class PolicyTablePanel extends React.Component<RootProps, {}> {
  render() {
    let { actions, policyAssociation, policyAssociatedList } = this.props;

    const columns: TableColumn<PolicyPlain>[] = [
      {
        key: 'name',
        header: t('名称'),
        render: (policy, text, index) => (
          <Text parent="div" overflow>
            {policy.displayName || '-'}
          </Text>
        ),
      },
      {
        key: 'category',
        header: t('类型'),
        render: (policy, text, index) => (
          <Text parent="div" overflow>
            {policy.category || '-'}
          </Text>
        ),
      },
      {
        key: 'description',
        header: t('描述'),
        render: (policy, text, index) => (
          <Text parent="div" overflow>
            {policy.description || '-'}
          </Text>
        ),
      },
      // { key: 'operation', header: t('操作'), render: policy => this._renderOperationCell(policy) }
    ];

    return (
      <TablePanel
        columns={columns}
        recordKey={'id'}
        records={policyAssociation.originPolicies}
        action={actions.policy.associate.policyAssociatedList}
        model={policyAssociatedList}
        emptyTips={emptyTips}
      />
    );
  }

  /** 渲染操作按钮 */
  _renderOperationCell = (policy: PolicyPlain) => {
    let { actions } = this.props;
    return (
      <React.Fragment>
        <LinkButton
          tipDirection="right"
          onClick={(e) => {
            this._removePolicy(policy);
          }}
        >
          <Trans>解除关联</Trans>
        </LinkButton>
      </React.Fragment>
    );
  };

  _removePolicy = async (policy: PolicyPlain) => {
    let { actions, policyFilter } = this.props;
    const yes = await Modal.confirm({
      message: t('确认解除当前策略关联') + ` - ${policy.displayName}？`,
      okText: t('解除'),
      cancelText: t('取消'),
    });
    if (yes) {
      /** 目前还没有实现基于用户组解绑策略 */
      let policyAssociation: PolicyAssociation = { id: uuid(), removePolicies: [policy] };
      actions.policy.associate.disassociatePolicyWorkflow.start([policyAssociation], policyFilter);
      actions.policy.associate.disassociatePolicyWorkflow.perform();
    }
  };
}
