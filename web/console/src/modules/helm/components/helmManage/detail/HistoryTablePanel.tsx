import classNames from 'classnames';
import * as React from 'react';
import { connect } from 'react-redux';

import { Bubble, Button, Card, ContentView, Modal, Text } from '@tea/component';
import { TablePanel, TablePanelColumnProps } from '@tencent/ff-component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { dateFormatter } from '../../../../../../helpers';
import { LinkButton } from '../../../../common/components';
import { DialogBodyLayout } from '../../../../common/layouts';
import { allActions } from '../../../actions';
import { helmStatus } from '../../../constants/Config';
import { HelmHistory } from '../../../models';
import { RootProps } from '../../HelmApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

interface State {
  showRollbackDialog?: boolean;
  rollbackHistory?: HelmHistory;
}

@connect(state => state, mapDispatchToProps)
export class HistoryTablePanel extends React.Component<RootProps, State> {
  state = {
    showRollbackDialog: false,
    rollbackHistory: null
  };
  componentDidMount() {
    let { actions, route } = this.props;
    actions.detail.queryHistory.applyFilter({
      helmName: route.queries['helmName'],
      regionId: +route.queries['rid'],
      clusterId: route.queries['clusterId']
    });
  }
  render() {
    return (
      <ContentView>
        <ContentView.Body>
          {this._renderTablePanel()}
          {this.state.showRollbackDialog && this._renderUpdateConfirmDialog()}
        </ContentView.Body>
      </ContentView>
    );
  }
  private _renderUpdateConfirmDialog() {
    const { actions } = this.props;
    const cancel = () => this.setState({ showRollbackDialog: false, rollbackHistory: null });
    const confirm = () => {
      actions.detail.rollback(this.state.rollbackHistory.name, this.state.rollbackHistory.version);
      cancel();
    };

    return (
      <Modal visible={true} caption={t('回滚Helm应用')} onClose={cancel} size={485} disableEscape={true}>
        <Modal.Body>
          <DialogBodyLayout>
            <p className="til">{t('是否立即回滚Helm应用？')}</p>
          </DialogBodyLayout>
        </Modal.Body>
        <Modal.Footer>
          <Button type="primary" onClick={confirm}>
            {t('确认')}
          </Button>
          <Button onClick={cancel}>{t('取消')}</Button>
        </Modal.Footer>
      </Modal>
    );
  }

  private _renderTablePanel() {
    let {
      actions,
      detailState: { histories, historyQuery },
      route
    } = this.props;

    const columns: TablePanelColumnProps<HelmHistory>[] = [
      {
        key: 'name',
        header: t('应用名'),
        width: '15%'
        // render: x => (
        //   <Text parent="div" overflow>
        //     {x.name}
        //   </Text>
        // )
      },
      {
        key: 'helm@version',
        header: t('部署详情'),
        width: '15%',
        render: x => (
          <Text parent="div" overflow>
            {x.chart.metadata.name}@{x.chart.metadata.version}
          </Text>
        )
      },
      {
        key: 'description',
        header: t('描述'),
        width: '30%',
        render: x => (
          <div className="form-unit">
            <p className="text-overflow" style={{ width: '90%', float: 'left' }}>
              {x.info.Description}
            </p>
            {x.info.Description.length > 30 && (
              <Bubble placement="top" content={x.info.Description || null}>
                <i className="plaint-icon" style={{ marginLeft: '5px' }} />
              </Bubble>
            )}
          </div>
        )
      },
      {
        key: 'status',
        header: t('状态'),
        width: '8%',
        render: x => (
          <div>
            <span
              className={classNames(
                'text-overflow',
                helmStatus[x.info.status.code] && helmStatus[x.info.status.code].classname
              )}
            >
              {helmStatus[x.info.status.code] ? helmStatus[x.info.status.code].text : '-'}
            </span>
          </div>
        )
      },
      {
        key: 'version',
        header: t('版本号'),
        width: '7%'
        // render: x => x.version + ''
      },
      {
        key: 'createTime',
        header: t('部署时间'),
        width: '15%',
        render: x => dateFormatter(new Date(x.info.first_deployed), 'YYYY-MM-DD HH:mm:ss')
      }
      // {
      //   key: 'operation',
      //   header: t('操作'),
      //   width: '10%',
      //   render: x => this._renderOperationCell(x)
      // }
    ];

    return (
      <TablePanel
        columns={columns}
        emptyTips={<div className="text-center">{t('该Helm应用暂无版本历史')}</div>}
        model={{
          list: histories,
          query: historyQuery
        }}
        action={Object.assign({}, actions.detail.queryHistory, actions.detail.fetchHistory)}
        getOperations={x => this._renderOperationButtons(x)}
      />
    );
  }

  /** 渲染操作按钮 */
  private _renderOperationButtons(history: HelmHistory) {
    let buttons = [];
    if (history.info.status.code !== 'DEPLOYED') {
      buttons.push(<LinkButton onClick={e => this.rollback(history)}>{t('回滚')}</LinkButton>);
    }
    return buttons;
    // return (
    //   <div>
    //     {history.info.status.code !== 'DEPLOYED' && (
    //       <LinkButton onClick={e => this.rollback(history)}>{t('回滚')}</LinkButton>
    //     )}
    //   </div>
    // );
  }

  private rollback(history: HelmHistory) {
    this.setState({
      showRollbackDialog: true,
      rollbackHistory: history
    });
  }
}
