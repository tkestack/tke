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

import { bindActionCreators, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Bubble, Button, Icon, Modal, Table, TableColumn, Text } from '@tencent/tea-component';
import { scrollable } from 'tea-component/es/table/addons';

import { dateFormatter } from '../../../../../helpers';
import { ClusterCondition } from '../../../common';
import { allActions } from '../../actions';
import { DialogNameEnum } from '../../models';
import { RootProps } from '../ClusterApp';

const conditionStatusMap = {
  True: '成功',
  Unknown: '待处理',
  False: '失败'
};

const conditionStatusType = {
  True: 'success',
  False: 'danger',
  Unknown: 'text'
};

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ClusterStatusDialog extends React.Component<RootProps, {}> {
  render() {
    const { cluster, dialogState } = this.props;

    if (!cluster.selection) return <noscript />;

    const isShowDialog = dialogState[DialogNameEnum.clusterStatusDialog];

    const columns: TableColumn<ClusterCondition>[] = [
      {
        key: 'type',
        header: t('类型'),
        render: x => <Text>{x.type}</Text>
      },
      {
        key: 'status',
        header: t('状态'),
        render: x => (
          <Text theme={conditionStatusType[x.status]}>
            {conditionStatusMap[x.status] ? conditionStatusMap[x.status] : '-'}
          </Text>
        )
      },
      {
        key: 'probeTime',
        header: t('最后探测时间'),
        render: x => <Text>{dateFormatter(new Date(x.lastProbeTime), 'YYYY-MM-DD HH:mm:ss') || '-'}</Text>
      },
      {
        key: 'reason',
        header: t('原因'),
        render: x => {
          const isFailed = x.status === 'False';

          return (
            <React.Fragment>
              <Text verticalAlign="middle">{isFailed && x.reason ? x.reason : '-'}</Text>
              {isFailed && (
                <Bubble content={isFailed ? (x.message ? x.message : '未知错误') : null}>
                  <Icon type="error" />
                </Bubble>
              )}
            </React.Fragment>
          );
        }
      }
    ];

    return (
      <Modal
        size={700}
        visible={isShowDialog}
        caption={`集群 ${cluster.selection.metadata.name} 的状态`}
        onClose={this._handleClose.bind(this)}
      >
        <Modal.Body>
          <Table
            records={
              cluster.selection ? (cluster.selection.status.conditions ? cluster.selection.status.conditions : []) : []
            }
            columns={columns}
            addons={[
              scrollable({
                maxHeight: 600
              })
            ]}
          />
        </Modal.Body>
        <Modal.Footer>
          <Button type="primary" onClick={this._handleClose.bind(this)}>
            关闭
          </Button>
        </Modal.Footer>
      </Modal>
    );
  }

  /** 关闭按钮 */
  private _handleClose() {
    const { actions } = this.props;
    actions.dialog.updateDialogState(DialogNameEnum.clusterStatusDialog);
  }
}
