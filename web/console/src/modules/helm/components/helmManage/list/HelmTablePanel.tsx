import * as React from 'react';
import { Modal, Button, Bubble, Card, Text } from '@tea/component';
import { bindActionCreators } from '@tencent/qcloud-lib';
import { connect } from 'react-redux';
import { allActions } from '../../../actions';
import { Helm } from '../../../models';
import { RootProps } from '../../HelmApp';
import classNames from 'classnames';
import { helmStatus, ClusterHelmStatus, HelmResource } from '../../../constants/Config';
import { TipInfo } from '../../../../common/components/';
import { dateFormatter } from '../../../../../../helpers';
import { InstallingHelmContent } from './InstallingDialog';
import { UpdateHelmDialog } from './UpdateHelmDialog';
import { UpdateHelmDialogOther } from './UpdateHelmDialogOther';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { router } from '../../../router';
import { TablePanel } from '@tencent/ff-component';

interface State {
  showSetupHelmDialog?: boolean;
  showDeleteHelmDialog?: boolean;
  deleteHelm?: Helm;
  showUpdateHelmDialog?: boolean;

  showInstallingHelmDialog?: boolean;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class HelmTablePanel extends React.Component<RootProps, State> {
  state = {
    showSetupHelmDialog: false,
    showDeleteHelmDialog: false,
    deleteHelm: null,
    showUpdateHelmDialog: false,

    showInstallingHelmDialog: false
  };
  render() {
    const {
      listState: { clusterHelmStatus, installingHelmList, cluster }
    } = this.props;

    const reason = clusterHelmStatus.reason;

    let version = cluster.selection ? cluster.selection.status.version.split('.') : [];
    let enableHelm = version.length > 2 && +version[1] >= 8;
    return (
      <div>
        {clusterHelmStatus.code === ClusterHelmStatus.NONE &&
          // <div className="manage-area-main-inner">
          (enableHelm ? (
            <TipInfo className="warning">
              <span style={{ verticalAlign: 'middle' }}>
                <Trans>
                  该集群暂未开通Helm应用，开通需在集群内安装Helm tiller组件，需要占用一定资源，如需使用请
                  <a
                    href="javascript:void(0);"
                    onClick={event => {
                      this.setupHelm();
                    }}
                  >
                    申请开通
                  </a>
                </Trans>
              </span>
            </TipInfo>
          ) : (
            <TipInfo className="warning">
              <span style={{ verticalAlign: 'middle' }}>
                <Trans>Helm应用管理仅支持kubernetes 1.8以上版本的集群。</Trans>
              </span>
            </TipInfo>
          ))

        // </div>
        }
        {clusterHelmStatus.code === ClusterHelmStatus.CHECKING && (
          // <div className="manage-area-main-inner">
          <TipInfo className="warning">
            <span style={{ verticalAlign: 'middle' }}>
              {t('正在校验Helm应用管理功能')}
              <i className="n-loading-icon" />
            </span>
          </TipInfo>
          // </div>
        )}
        {clusterHelmStatus.code === ClusterHelmStatus.INIT && (
          // <div className="manage-area-main-inner">
          <TipInfo className="warning">
            <span style={{ verticalAlign: 'middle' }}>
              {t('正在开通Helm应用管理功能')}
              <i className="n-loading-icon" />
            </span>
          </TipInfo>
          // </div>
        )}
        {clusterHelmStatus.code === ClusterHelmStatus.REINIT && (
          // <div className="manage-area-main-inner">
          <TipInfo className="warning">
            <span style={{ verticalAlign: 'middle' }}>
              <Trans>
                开通失败
                <Bubble placement="top" content={reason || null}>
                  <i className="plaint-icon" style={{ marginLeft: '5px' }} />
                </Bubble>
                ，正在重新开通Helm应用管理功能
                <i className="n-loading-icon" />
              </Trans>
            </span>
          </TipInfo>
          // </div>
        )}
        {clusterHelmStatus.code === ClusterHelmStatus.ERROR && (
          // <div className="manage-area-main-inner">
          <TipInfo className="warning">
            <span style={{ verticalAlign: 'middle' }}>
              <Trans>
                开通失败
                <Bubble placement="top" content={reason || null}>
                  <i className="plaint-icon" style={{ marginLeft: '5px' }} />
                </Bubble>
                ，请确认集群保留足够的空闲资源，并
                <a
                  href="javascript:void(0);"
                  onClick={event => {
                    this.setupHelm();
                  }}
                >
                  重新开通
                </a>
              </Trans>
            </span>
          </TipInfo>
          // </div>
        )}
        {/* {clusterHelmStatus.code === ClusterHelmStatus.RUNNING && ( */}
        <Card>
          <Card.Body>
            {installingHelmList.fetched && installingHelmList.data.recordCount > 0 && this._renderInstallingPanel()}
            {this._renderTablePanel()}
          </Card.Body>
        </Card>
        {/* )} */}
        {this._renderSetupHelmDialog()}
        {this._renderUpdateHelmDialog()}
        {this._renderDeleteConfirmDialog()}
        {this._renderInstallingDialog()}
      </div>
    );
  }

  private _renderInstallingPanel() {
    return (
      <TipInfo>
        <span>
          {t('正在创建新的Helm应用...')}
          <i className="n-loading-icon" />
          &nbsp;
          <a href="javascript:void(0);" onClick={() => this.showInstalling()}>
            {t('查看详情')}
          </a>
        </span>
      </TipInfo>
    );
  }

  private _renderUpdateHelmDialog() {
    let {
      listState: { helmSelection }
    } = this.props;
    if (this.state.showUpdateHelmDialog) {
      if (helmSelection.chart_metadata.repo === HelmResource.TencentHub) {
        return (
          <UpdateHelmDialog
            onCancel={e => {
              this.setState({
                showUpdateHelmDialog: false
              });
            }}
            {...this.props}
          />
        );
      } else {
        return (
          <UpdateHelmDialogOther
            onCancel={e => {
              this.setState({
                showUpdateHelmDialog: false
              });
            }}
            {...this.props}
          />
        );
      }
    }
  }

  private _renderSetupHelmDialog() {
    const { actions } = this.props;
    const cancel = () => this.setState({ showSetupHelmDialog: false });
    const confirm = () => {
      actions.helm.setupHelm();
      cancel();
    };

    return (
      <Modal
        visible={this.state.showSetupHelmDialog}
        caption={t('集群Helm 应用管理功能')}
        onClose={cancel}
        size={485}
        disableEscape={true}
      >
        <Modal.Body>
          <Trans>
            <p className="til">新建Helm应用需要先开通Helm应用，当前所选集群暂未开通。</p>
            <p className="til">开通Helm应用功能：</p>
            <ul style={{ marginLeft: 30 }}>
              <li>
                <p className="til">1.将在集群内安装Helm tiller组件</p>
              </li>
              <li>
                <p className="til">
                  2.将占用集群 <em className="text-warning">0.28 核CPU 180Mi </em>的资源
                </p>
              </li>
            </ul>
          </Trans>
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

  private _renderDeleteConfirmDialog() {
    const { actions } = this.props;
    const cancel = () => this.setState({ showDeleteHelmDialog: false, deleteHelm: null });
    const confirm = () => {
      actions.helm.delete(this.state.deleteHelm);
      cancel();
    };

    if (!this.state.showDeleteHelmDialog) {
      return <React.Fragment />;
    }

    return (
      <Modal
        visible={this.state.showDeleteHelmDialog}
        caption={t('您确定要删除【{{name}}】吗？', {
          name: this.state.deleteHelm.name
        })}
        onClose={cancel}
        size={485}
        disableEscape={true}
      >
        <Modal.Body>
          {t('删除应用将删除该应用创建的所有K8S资源，删除后所有数据将被清除且不可恢复,请提前备份数据。')}
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

  private _renderInstallingDialog() {
    const { actions } = this.props;
    const cancel = () => this.setState({ showInstallingHelmDialog: false });
    return (
      <Modal
        visible={this.state.showInstallingHelmDialog}
        caption={t('Helm应用安装日志')}
        onClose={cancel}
        size={960}
        disableEscape={true}
      >
        <Modal.Body>
          <InstallingHelmContent {...this.props} />
        </Modal.Body>
      </Modal>
    );
  }

  private _renderMoreBtn() {
    let {
      actions,
      listState: { helmList, helmQuery, helmSelection },
      route
    } = this.props;
    if (helmList.data.recordCount >= helmQuery.paging.pageSize) {
      return (
        <span>
          <div
            style={{
              overflow: 'visible',
              fontSize: '14px',
              lineHeight: '54px',
              textAlign: 'center'
            }}
          >
            <a
              href="javascript:void(0);"
              onClick={() => {
                actions.helm.changePaging({
                  pageIndex: 1,
                  pageSize: helmQuery.paging.pageSize + 10
                });
              }}
              className="text-center"
            >
              {t('加载更多')}
            </a>
          </div>
        </span>
      );
    }
  }

  private _renderTablePanel() {
    let {
      actions,
      listState: { helmList, helmQuery, helmSelection },
      route
    } = this.props;

    const columns = [
      {
        key: 'name',
        header: t('应用名'),
        width: '20%',
        render: x => (
          <Text parent="div" overflow>
            <a
              href="javascript:void(0);"
              onClick={() => {
                this.click(x);
              }}
            >
              {x.name}
            </a>
          </Text>
        )
      },
      {
        key: 'status',
        header: t('状态'),
        width: '10%',
        render: x => (
          <Text>
            <span
              className={classNames(
                'text-overflow',
                helmStatus[x.info.status.code] && helmStatus[x.info.status.code].classname
              )}
            >
              {helmStatus[x.info.status.code] ? helmStatus[x.info.status.code].text : '-'}
            </span>
            {x.info.status.code === 'FAILED' && (
              <Bubble placement="top" content={x.info.Description || null}>
                <i className="plaint-icon" style={{ marginLeft: '5px' }} />
              </Bubble>
            )}
          </Text>
        )
      },
      {
        key: 'version',
        header: t('版本号'),
        width: '10%',
        render: x => <Text overflow> {x.version + ''}</Text>
      },
      {
        key: 'createTime',
        header: t('创建时间'),
        width: '15%',
        render: x => <Text overflow> {dateFormatter(new Date(x.info.first_deployed), 'YYYY-MM-DD HH:mm:ss')}</Text>
      },
      {
        key: 'resource',
        header: t('Chart仓库'),
        width: '10%',
        render: x => <Text overflow> {x.chart_metadata.repo}</Text>
      },
      {
        key: 'namespace',
        header: t('Chart命名空间'),
        width: '10%',
        render: x => <Text overflow> {x.chart_metadata.chart_ns || '-'}</Text>
      },
      {
        key: 'chartversion',
        header: t('Chart版本'),
        width: '10%',
        render: x => <Text overflow> {x.chart_metadata.version + ''}</Text>
      },
      {
        key: 'operation',
        header: t('操作'),
        width: '15%',
        render: x => this._renderOperationCell(x)
      }
    ];

    return (
      <React.Fragment>
        <TablePanel
          columns={columns}
          model={{
            list: helmList,
            query: helmQuery
          }}
          action={actions.helm}
          isNeedCard={false}
        />
        {this._renderMoreBtn()}
      </React.Fragment>
    );
  }

  /** 渲染操作按钮 */
  private _renderOperationCell(helm: Helm) {
    return (
      <div>
        <a href="javascript:void(0)" onClick={e => this.update(helm)}>
          {t('更新应用')}
        </a>
        &nbsp;&nbsp;
        <a href="javascript:void(0)" onClick={e => this.delete(helm)}>
          {t('删除')}
        </a>
      </div>
    );
  }

  private setupHelm() {
    this.setState({
      showSetupHelmDialog: true
    });
  }

  private showInstalling() {
    this.setState({
      showInstallingHelmDialog: true
    });
  }
  private click(helm: Helm) {
    const {
      route,
      listState: { cluster, region }
    } = this.props;
    router.navigate(
      { sub: 'detail' },
      Object.assign({}, route.queries, {
        rid: region.selection.value + '',
        clusterId: cluster.selection.metadata.name,
        helmName: helm.name
      })
    );
  }
  private update(helm: Helm) {
    this.setState({
      showUpdateHelmDialog: true
    });
    this.props.actions.helm.select(helm);
  }
  private delete(helm: Helm) {
    this.setState({
      showDeleteHelmDialog: true,
      deleteHelm: helm
    });
  }
}
