import * as React from 'react';
import { bindActionCreators, uuid } from '@tencent/qcloud-lib';
import { allActions } from '../../../actions';
import { connect } from 'react-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Modal, TableColumn, Text, Bubble, Icon, Button, Table } from '@tencent/tea-component';
import { ClusterCondition } from '../../../../common';
import { dateFormatter } from '../../../../../../helpers';
import { scrollable } from '@tencent/tea-component/lib/table/addons';
import { ListModel } from '@tencent/redux-list';
import { Computer, ComputerFilter, DialogNameEnum, DialogState } from '../../../../../modules/cluster/models';

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

let nodeConditionStaus = ['OutOfDisk', 'MemoryPressure', 'DiskPressure', 'PIDPressure'];
interface ComputerStatusDialogProps {
  computer: ListModel<Computer, ComputerFilter>;
  dialogState: DialogState;
  actions?: typeof allActions;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(
  state => state,
  mapDispatchToProps
)
export class ComputerStatusDialog extends React.Component<ComputerStatusDialogProps, {}> {
  render() {
    let { computer, dialogState } = this.props;

    if (!computer.selection) return <noscript />;
    let isShowDialog = dialogState[DialogNameEnum.computerStatusDialog];
    let columns: TableColumn<ClusterCondition>[] = [
      {
        key: 'type',
        header: t('类型'),
        render: x => <Text>{x.type}</Text>
      },
      {
        key: 'status',
        header: t('状态'),
        render: x => {
          let status;
          if (nodeConditionStaus.indexOf(x.type) !== -1) {
            status = x.status === 'False' ? 'success' : 'danger';
          } else {
            status = conditionStatusType[x.status];
          }
          return <Text theme={status}>{conditionStatusMap[x.status] ? conditionStatusMap[x.status] : '-'}</Text>;
        }
      },
      {
        key: 'probeTime',
        header: t('最后探测时间'),
        render: x => (
          <Text>{dateFormatter(new Date(x.lastProbeTime || x.lastHeartbeatTime), 'YYYY-MM-DD HH:mm:ss') || '-'}</Text>
        )
      },
      {
        key: 'reason',
        header: t('原因'),
        render: x => {
          let isFailed;
          if (nodeConditionStaus.indexOf(x.type) !== -1) {
            isFailed = x.status !== 'False';
          } else {
            isFailed = x.status === 'False';
          }
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
        caption={`集群 ${computer.selection.metadata.name} 的状态`}
        onClose={this._handleClose.bind(this)}
      >
        <Modal.Body>
          <Table
            records={
              computer.selection
                ? computer.selection.status.conditions
                  ? computer.selection.status.conditions
                  : []
                : []
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
    let { actions } = this.props;
    actions.dialog.updateDialogState(DialogNameEnum.computerStatusDialog);
  }
}
