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
import { RootProps } from '../RoleApp';
import { PolicyAssociateWorkflowDialog } from '../associate/PolicyAssociateWorkflowDialog';
import { PolicyFilter } from '../../../models';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class PolicyActionPanel extends React.Component<RootProps, {}> {

  componentWillUnmount() {
    let { actions } = this.props;
    actions.policy.associate.clearPolicyAssociation();
  }

  componentDidMount() {
    const { actions, route } = this.props;
    /** 设置策略关联场景 */
    let filter: PolicyFilter = {
      resource: 'role',
      resourceID: route.queries['roleName'],
      /** 关联/解关联回调函数 */
      callback: () => {
        actions.role.detail.fetchRole({ name: route.queries['roleName'] });
      }
    };
    actions.policy.associate.setupPolicyFilter(filter);
    /** 拉取关联策略列表，拉取后自动更新policyAssociation */
    actions.policy.associate.policyAssociatedList.applyFilter(filter);
    /** 拉取策略列表 */
    actions.policy.associate.policyList.performSearch('');
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
              <Button type="primary" onClick={e => {
                /** 开始关联策略工作流 */
                actions.policy.associate.associatePolicyWorkflow.start();
              }}>
                {t('关联策略')}
              </Button>
            }
          />
        </Table.ActionPanel>
        <PolicyAssociateWorkflowDialog />
      </React.Fragment>
    );
  }

}

