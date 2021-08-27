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
import * as classnames from 'classnames';
import * as React from 'react';

import { Button, Modal, TableColumn, Text } from '@tea/component';
import { stylize } from '@tea/component/table/addons/stylize';
import { isSuccessWorkflow, OperationState } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { GridTable, TipInfo } from '../../../../common/components';
import { getWorkflowError } from '../../../../common/utils';
import { Resource } from '../../../models';
import { RootProps } from '../../ClusterApp';

interface BatchDrainComputerDialogState {
  isCollapsed?: boolean;
}

export class BatchDrainComputerDialog extends React.Component<RootProps, BatchDrainComputerDialogState> {
  constructor(props, context) {
    super(props, context);
    this.state = {
      isCollapsed: false
    };
  }

  render() {
    let { isCollapsed } = this.state;
    let classname = classnames({ 'blue-down-icon': isCollapsed }, { 'blue-up-icon': !isCollapsed });

    let { actions, subRoot } = this.props,
      { drainComputer, computerPodList, computerPodQuery } = subRoot.computerState;
    let action = actions.workflow.batchDrainComputer;
    let workflow = drainComputer;

    if (workflow.operationState === OperationState.Pending) {
      return <noscript />;
    }

    const cancel = () => {
      if (workflow.operationState === OperationState.Done) {
        action.reset();
      }
      if (workflow.operationState === OperationState.Started) {
        action.cancel();
      }
    };

    const perform = () => {
      action.start(workflow.targets, {
        clusterId: workflow.params.clusterId
      });
      action.perform();
    };

    const failed = workflow.operationState === OperationState.Done && !isSuccessWorkflow(workflow);

    let colunms: TableColumn<Resource>[] = [
      {
        key: 'name',
        header: t('实例（Pod）名称'),
        width: '55%',
        render: x => (
          <Text parent="div" overflow>
            <span title={x.metadata.name}>{x.metadata.name}</span>
          </Text>
        )
      },
      {
        key: 'namespace',
        header: t('所属集群空间'),
        width: '45%',
        render: x => (
          <Text parent="div" overflow>
            <span title={x.metadata.namespace}>{x.metadata.namespace}</span>
          </Text>
        )
      }
    ];

    let computerPodCount = computerPodList.data.recordCount;

    return (
      <Modal
        visible={true}
        caption={t('您确定对选中节点进行驱逐么？')}
        onClose={cancel}
        size={485}
        disableEscape={true}
      >
        <Modal.Body>
          <div style={{ fontSize: '14px', lineHeight: '20px' }}>
            <div className="docker-dialog jiqun">
              {
                <div className="act-outline">
                  <div className="act-summary">
                    <p>
                      <Trans count={computerPodCount}>
                        <span>
                          节点包含<strong className="text-warning">{{ computerPodCount }}个</strong>实例，
                        </span>
                        <a href="javascript:;" onClick={this.toggleHandle.bind(this)}>
                          查看详情 <i className={classname} />
                        </a>
                      </Trans>
                    </p>
                  </div>
                  {!isCollapsed && (
                    <div className="del-colony-tb">
                      <GridTable
                        columns={colunms}
                        emptyTips={<div className="text-center">{t('节点的实例（Pod）列表为空')}</div>}
                        listModel={{
                          list: computerPodList,
                          query: computerPodQuery
                        }}
                        actionOptions={actions.computerPod}
                        addons={[
                          stylize({
                            className: 'ovm-dialog-tablepanel',
                            bodyStyle: { overflowY: 'auto', height: 160, minHeight: 100 }
                          })
                        ]}
                        isNeedCard={false}
                      />
                    </div>
                  )}
                </div>
              }
              <Trans>
                <div className="text">
                  节点驱逐后，将会把节点内的所有Pod（不包含DaemonSet管理的Pod）从节点中驱逐到集群内其他节点，并将节点设置为封锁状态。
                </div>
                <div className="text text-danger">注意：本地存储的Pod被驱逐后数据将丢失，请谨慎操作</div>
              </Trans>
            </div>
            {failed && <TipInfo type="error">{getWorkflowError(workflow)}</TipInfo>}
          </div>
        </Modal.Body>
        <Modal.Body>
          <Modal.Footer>
            <Button type="primary" disabled={workflow.operationState === OperationState.Performing} onClick={perform}>
              {failed ? t('重试') : t('确定')}
            </Button>
            <Button onClick={cancel}>{t('取消')}</Button>
          </Modal.Footer>
        </Modal.Body>
      </Modal>
    );
  }

  toggleHandle() {
    this.setState({
      isCollapsed: !this.state.isCollapsed
    });
  }
}
