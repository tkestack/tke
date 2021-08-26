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

import * as classnames from 'classnames';
import * as React from 'react';
import { connect } from 'react-redux';

import { FormPanel } from '@tencent/ff-component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Justify } from '@tencent/tea-component';

import { LinkButton } from '../../../../common/components';
import { allActions } from '../../../actions';
import { validateWorkloadActions } from '../../../actions/validateWorkloadActions';
import { RootProps } from '../../ClusterApp';

interface ContainerListItemProps extends RootProps {
  /** 容器的id */
  cKey?: string;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ResourceContainerListItem extends React.Component<ContainerListItemProps, {}> {
  render() {
    let { actions, subRoot, cKey } = this.props,
      { workloadEdit } = subRoot,
      { containers, volumes, canAddContainer } = workloadEdit;

    let container = containers.find(c => c.id === cKey),
      editingContainer = containers.find(c => c.status === 'editing'),
      canEdit = canAddContainer,
      canDelete = containers.length > 1;

    if (editingContainer) {
      canEdit = validateWorkloadActions._canAddContainer(editingContainer, volumes);
    } else {
      canEdit = true;
    }

    let validatedStatus = validateWorkloadActions._canAddContainer(container, volumes) || container.isAdvancedError;

    // 容器上展示的内容
    let cText = container.name + '(' + container.registry;
    if (container.tag) {
      cText += ':' + container.tag + ')';
    } else {
      cText += ')';
    }

    return (
      <div className={classnames('run-docker-box', { 'run-docker-error': !validatedStatus })}>
        <Justify
          left={<FormPanel.Text>{cText}</FormPanel.Text>}
          right={
            <React.Fragment>
              <LinkButton
                disabled={!canEdit}
                tip={t('编辑')}
                errorTip={t('请完成待编辑项')}
                onClick={() => this._handleEditButton(cKey)}
              >
                <i className="icon-edit-gray" />
              </LinkButton>
              <LinkButton
                disabled={!canDelete}
                tip={t('删除')}
                errorTip={t('不可删除，至少创建一个容器')}
                onClick={() => actions.editWorkload.deleteContainer(cKey)}
              >
                <i className="icon-cancel-icon" />
              </LinkButton>
            </React.Fragment>
          }
        />
      </div>
    );
  }

  /** 处理编辑按钮 */
  private _handleEditButton(cKey: string) {
    let { actions, subRoot } = this.props,
      { workloadEdit } = subRoot,
      { containers, volumes } = workloadEdit;

    let editingContainer = containers.find(c => c.status === 'editing');

    if (editingContainer) {
      actions.validate.workload.validateContainer(editingContainer);
      if (validateWorkloadActions._validateContainer(editingContainer, volumes, containers)) {
        actions.editWorkload.updateContainer({ status: 'edited' }, editingContainer.id + '');
        actions.editWorkload.updateContainer({ status: 'editing' }, cKey);
      }
    } else {
      actions.editWorkload.updateContainer({ status: 'editing' }, cKey);
    }
  }
}
