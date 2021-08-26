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
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators, insertCSS } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { RootProps } from '../RoleApp';
import { WorkflowDialog } from '../../../../common/components';
import { GroupAssociatePanel } from './GroupAssociatePanel';
const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

interface WorkflowDialogProps extends RootProps {
  onPostCancel?: () => void;
}

@connect(state => state, mapDispatchToProps)
export class GroupAssociateWorkflowDialog extends React.Component<WorkflowDialogProps, {}> {

  render() {
    const {
      actions,
      groupAssociation,
      groupFilter,
      associateGroupWorkflow
    } = this.props;
    const { onPostCancel = undefined } = this.props;
    return (
      <WorkflowDialog
        caption={t('关联用户组')}
        workflow={associateGroupWorkflow}
        action={actions.group.associate.associateGroupWorkflow}
        targets={[groupAssociation]}
        params={groupFilter}
        postAction={() => {
          //清空查询条件，重新拉取
          // actions.group.associate.groupList.performSearch('');
          onPostCancel && onPostCancel();
        }}
        width={700}
      >
        <GroupAssociatePanel />
      </WorkflowDialog>
    );
  }
}
