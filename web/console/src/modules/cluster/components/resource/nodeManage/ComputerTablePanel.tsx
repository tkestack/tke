import * as React from 'react';
import { connect } from 'react-redux';

import { Button, Icon, Modal, Text } from '@tea/component';
import { selectable } from '@tea/component/table/addons/selectable';
import { TablePanel, TablePanelColumnProps } from '@tencent/ff-component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Bubble } from '@tencent/tea-component';
import { sortable, SortBy } from '@tencent/tea-component/lib/table/addons/sortable';

import { dateFormatter } from '../../../../../../helpers';
import { Clip, LinkButton } from '../../../../common/components';
import { includes } from '../../../../common/utils';
import { allActions } from '../../../actions';
import { Computer, DialogNameEnum } from '../../../models';
import { router } from '../../../router';
import { RootProps } from '../../ClusterApp';
import { ReduceRequest } from '../resourceDetail/ResourcePodPanel';

export const ComputerStatus = {
  Running: 'success',
  Initializing: 'label',
  Failed: 'danger'
};

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

interface State {
  showOsTips?: boolean;
  selectCluster?: any;
  sorts?: SortBy[];
}

@connect(state => state, mapDispatchToProps)
export class ComputerTablePanel extends React.Component<RootProps, State> {
  state = {
    showOsTips: false,
    selectCluster: null,
    sorts: []
  };

  handleExpand() {
    let { route, cluster } = this.props;

    if ((cluster as any).os) {
      router.navigate({ sub: 'expand' }, { rid: route.queries['rid'], clusterId: route.queries['clusterId'] });
    } else {
      this.setState({ showOsTips: true, selectCluster: cluster });
    }
  }

  renderOsTipsDialog = () => {
    let { showOsTips, selectCluster } = this.state;
    let { route } = this.props,
      urlParams = router.resolve(route);

    if (!showOsTips) return <noscript />;
    const gotoDetail = () => {
      router.navigate(
        Object.assign({}, urlParams, { mode: 'list', type: 'basic', resourceName: 'info' }),
        Object.assign({}, route.queries, {
          rid: route.queries['rid'],
          clusterId: selectCluster.clusterId
        })
      );
      hide();
    };
    const hide = () => {
      this.setState({ showOsTips: false, selectCluster: null });
    };
    return (
      <Modal visible={true} caption={t('设置集群操作系统')} onClose={hide} size={485} disableEscape={true}>
        <Modal.Body>{t('当前集群未设置操作系统，请前往集群基本信息页设置集群节点操作系统。')}</Modal.Body>
        <Modal.Footer>
          <Button type="primary" onClick={gotoDetail}>
            {t('前往设置')}
          </Button>
          <Button onClick={hide}>{t('取消')}</Button>
        </Modal.Footer>
      </Modal>
    );
  };

  render() {
    return (
      <React.Fragment>
        {this._renderTablePanel()}
        {this.renderOsTipsDialog()}
      </React.Fragment>
    );
  }

  private _renderTablePanel() {
    const { subRoot, actions, route, cluster } = this.props,
      urlParams = router.resolve(route),
      { computer } = subRoot.computerState;

    const columns: TablePanelColumnProps<Computer>[] = [
      {
        key: 'instanceId',
        header: t('节点名'),
        width: '15%',
        render: x => {
          let instanceId = x.metadata.name;

          return (
            <Text overflow>
              <a
                id={x.id + ''}
                href="javascript:;"
                onClick={() => {
                  this._handleClickForNavigate(instanceId);
                }}
              >
                {instanceId}
              </a>
              <Clip target={`#${x.id}`} />
            </Text>
          );
        }
      },
      {
        key: 'status',
        header: t('状态'),
        width: '8%',
        render: x => (
          <React.Fragment>
            <Text theme={ComputerStatus[x.status.phase]} verticalAlign="middle" parent={'p'}>
              {x.status.phase || '-'}
            </Text>
            <div className="sl-editor-name">
              {x.spec.unschedulable && (
                <span className="text-overflow m-width text-danger" title={t('已封锁')}>
                  {t('已封锁')}
                </span>
              )}
            </div>
            {x.status.phase === 'Initializing' && (
              // <Bubble content={t('点击查看详情')}>

              <Button
                type="link"
                // icon={x.status.phase === 'Failed' ? 'error' : 'loading'}
                onClick={() => {
                  actions.computer.select(x);
                  actions.dialog.updateDialogState(DialogNameEnum.computerStatusDialog);
                }}
              >
                {t('查看创建详情')}
              </Button>
              // </Bubble>
            )}
          </React.Fragment>
        )
      },
      {
        key: 'role',
        header: t('角色'),
        width: '10%',
        render: (x: Computer) => {
          return <Text verticalAlign="middle">{x.metadata.role}</Text>;
        }
      },
      {
        key: 'capacity',
        header: t('配置'),
        width: '12%',
        render: x => {
          let capacity = x.status.capacity;
          let capacityInfo = {
            cpu: capacity.cpu,
            memory: capacity.memory
          };
          let finalCpu = ReduceRequest('cpu', capacityInfo),
            finalmem = (ReduceRequest('memory', capacity) / 1024).toFixed(2);

          return (
            <React.Fragment>
              <Text verticalAlign="middle">
                {t('{{count}} 核, ', {
                  count: finalCpu
                })}
              </Text>
              <Text verticalAlign="middle">{`${finalmem} GB`}</Text>
            </React.Fragment>
          );
        }
      },
      {
        key: 'address',
        header: t('IP地址'),
        width: '15%',
        render: x => {
          let finalIPInfo = x.status.addresses.filter(item => item.type !== 'Hostname');

          return (
            <React.Fragment>
              {finalIPInfo.map((item, index) => (
                <p key={index}>
                  <Text id={item.type} verticalAlign="middle">
                    {item.address}
                  </Text>
                  <Clip target={`#${item.type}`} />
                </p>
              ))}
            </React.Fragment>
          );
        }
      },
      {
        key: 'podCIDR',
        header: t('PodCIDR'),
        width: '12%',
        render: x => <Text>{x.spec.podCIDR}</Text>
      },
      {
        key: 'createTime',
        header: t('创建时间'),
        width: '15%',
        render: x => {
          let time = dateFormatter(new Date(x.metadata.creationTimestamp), 'YYYY-MM-DD HH:mm:ss');

          let [year, currentTime] = time.split(' ');
          return (
            <React.Fragment>
              <Text parent="p" key={'year'}>
                {year}
              </Text>
              <Text parent="p" key={'currentTime'}>
                {currentTime}
              </Text>
            </React.Fragment>
          );
        }
      }
    ];

    let emptyTips: JSX.Element = (
      <div className="text-center">
        <Trans>
          您选择的集群的节点列表为空，您可以
          <a style={{ verticalAlign: '0' }} href="javascript:;" onClick={this.handleExpand.bind(this)}>
            新建节点
          </a>
          ，或切换到其他集群
        </Trans>
      </div>
    );
    return (
      <TablePanel
        columns={columns}
        getOperations={x => this._renderOperationCell(x)}
        action={actions.computer}
        model={computer}
        emptyTips={emptyTips}
        addons={[
          selectable({
            value: computer.selections.map(item => item.id as string),
            onChange: keys => {
              actions.computer.selects(
                computer.list.data.records.filter(record => {
                  return keys.indexOf(record.id as string) !== -1;
                })
              );
            }
          })
        ]}
      />
    );
  }

  /** 链接的跳转
   * @param resourceIns:string  node节点的instanceId
   */
  private _handleClickForNavigate(resourceIns: string) {
    let { actions, route, subRoot } = this.props,
      { resourceList } = subRoot.resourceOption,
      urlParams = router.resolve(route);

    // 选择当前选择的具体的resource
    let resourceSelection = resourceList.data.records.find(item => item.metadata.name === resourceIns);
    actions.resource.selectResource([resourceSelection]);
    // 进行路由的跳转
    router.navigate(
      Object.assign({}, urlParams, { mode: 'detail' }),
      Object.assign({}, route.queries, { resourceIns })
    );
  }

  /** 打开 unSchedule 的弹框 */
  private _handleBatchUnScheduleComputer(computer: Computer, clusterId: string) {
    let { actions } = this.props;
    actions.workflow.batchUnScheduleComputer.start();
    actions.computer.selects([computer]);
  }

  /** 打开 turn on Schedule 的弹框 */
  private _handeBatchTurnOnScheduleComputer(computer: Computer, clusterId: string) {
    let { actions } = this.props;
    actions.workflow.batchTurnOnScheduleComputer.start();
    actions.computer.selects([computer]);
  }

  /** 打开驱逐操作的弹框 */
  private _handleBatchDrainComputer(computer: Computer, clusterId: string) {
    let { actions, route } = this.props;
    actions.computerPod.applyFilter({ clusterId, specificName: computer.metadata.name });
    actions.workflow.batchDrainComputer.start([computer], { clusterId });
    actions.computer.selects([computer]);
  }

  /** 打开编辑标签操作的弹框 */
  private _handleUpdatelabel(computer: Computer, clusterId: string) {
    let { actions, route } = this.props,
      { rid } = route.queries;
    actions.computer.initLabelEdition(computer.metadata.labels, computer.metadata.name);
    actions.workflow.updateNodeLabel.start();
    actions.computer.selects([computer]);
  }

  /** 打开编辑标签操作的弹框 */
  private _handleUpdateTaint(computer: Computer, clusterId: string) {
    let { actions, route } = this.props,
      { rid } = route.queries;
    actions.computer.initTaintEdition(computer.spec.taints, computer.metadata.name);
    actions.workflow.updateNodeTaint.start();
    actions.computer.selects([computer]);
  }

  private _handleDeleteComputer(computer: Computer, clusterId: string) {
    let { actions } = this.props;

    actions.workflow.deleteComputer.start([computer]);
    actions.computerPod.applyFilter({ clusterId, specificName: computer.metadata.name });
    actions.computer.selects([computer]);
  }

  /**
   * 判断节点是否可以进行unschedule的操作
   * @param computer:Computer 当前操作的节点
   * @param mode:unSchedule|turn on scheduling    判断是否可进行操作
   */
  private getCanUnSchedule(computer: Computer, mode: 'unSchedule' | 'turnOnScheduling' = 'turnOnScheduling') {
    let computerStatus = computer.status.phase === 'Running';

    if (mode === 'unSchedule') {
      return computerStatus && (computer.spec.unschedulable == null || computer.spec.unschedulable === false);
    } else {
      return computerStatus && computer.spec.unschedulable === true;
    }
  }

  private _renderOperationCell(com: Computer) {
    const { cluster, subRoot, route, actions } = this.props,
      { computer } = subRoot.computerState;

    let clusterId = route.queries['clusterId'];

    let deleteDisable =
      com.metadata.role === 'Master&Etcd' || (cluster.selection && cluster.selection.spec.type === 'Imported');

    const renderDeleteButton = () => {
      return (
        <LinkButton
          key={'delete'}
          disabled={deleteDisable}
          errorTip={com.metadata.role === 'Master&Etcd' ? t('Master&Etcd节点不允许移除') : t('导入集群不允许移除节点')}
          onClick={() => this._handleDeleteComputer(com, clusterId)}
        >
          {t('移出')}
        </LinkButton>
      );
    };

    const renderUnScheduleButton = () => {
      const canUnSchedule = this.getCanUnSchedule(com, 'unSchedule');

      return (
        <LinkButton
          key={'unSchedule'}
          disabled={!canUnSchedule}
          errorTip={t('当前节点状态不能进行封锁的操作')}
          tip={t('封锁节点后，将不接受新的Pod调度到该节点')}
          tipDirection={'right'}
          onClick={() => this._handleBatchUnScheduleComputer(com, clusterId)}
        >
          {t('封锁')}
        </LinkButton>
      );
    };

    const renderTurnOnScheduleButton = () => {
      const canTurnOnSchedule = this.getCanUnSchedule(com, 'turnOnScheduling');
      const disabled = !canTurnOnSchedule;

      return (
        <LinkButton
          key={'schedule'}
          disabled={disabled}
          errorTip={t('当前节点状态不能进行取消封锁的操作')}
          tip={t('节点取消封锁后将允许新的Pod调度到该节点')}
          tipDirection={'right'}
          onClick={() => disabled || this._handeBatchTurnOnScheduleComputer(com, clusterId)}
        >
          {t('取消封锁')}
        </LinkButton>
      );
    };
    const renderDrainButton = () => {
      const canDrain = true;
      const disabled = !canDrain;

      return (
        <LinkButton
          key={'drain'}
          disabled={disabled}
          errorTip={t('当前节点状态不能进行驱逐的操作')}
          tip={t('将节点内的Pod从节点中驱逐到集群内其他节点')}
          tipDirection={'right'}
          onClick={() => disabled || this._handleBatchDrainComputer(com, clusterId)}
        >
          {t('驱逐')}
        </LinkButton>
      );
    };

    const renderLabelUpdateButton = () => {
      const disabled = false;

      return (
        <LinkButton
          key={'label'}
          disabled={disabled}
          errorTip={t('当前节点状态不能进行编辑标签的操作')}
          tip={t('编辑当前节点标签')}
          tipDirection={'right'}
          onClick={() => disabled || this._handleUpdatelabel(com, clusterId)}
        >
          {t('编辑标签')}
        </LinkButton>
      );
    };

    const renderTaintUpdateButton = () => {
      const disabled = false;

      return (
        <LinkButton
          key={'taint'}
          disabled={disabled}
          errorTip={t('当前节点状态不能进行编辑Taint的操作')}
          tip={t('编辑当前节点Taint')}
          tipDirection={'right'}
          onClick={() => disabled || this._handleUpdateTaint(com, clusterId)}
        >
          {t('编辑Taint')}
        </LinkButton>
      );
    };
    let btns = [renderDeleteButton(), renderDrainButton(), renderLabelUpdateButton(), renderTaintUpdateButton()];
    if (this.getCanUnSchedule(com, 'turnOnScheduling')) {
      btns.push(renderTurnOnScheduleButton());
    } else if (this.getCanUnSchedule(com, 'unSchedule')) {
      btns.push(renderUnScheduleButton());
    }
    return btns;
  }
}
