import * as React from 'react';
import { connect } from 'react-redux';

import { downloadKubeconfig } from '@helper/downloadText';
import { K8SUNIT, valueLabels1000, valueLabels1024 } from '@helper/k8sUnitUtil';
import { Alert, Bubble, Button, ExternalLink, Icon, Modal, TableColumn, Text } from '@tea/component';
import { FormPanel, TablePanel } from '@tencent/ff-component';
import { bindActionCreators, isSuccessWorkflow, OperationState, WorkflowState } from '@tencent/ff-redux';
import { initValidator } from '@tencent/ff-validator';
import { t } from '@tencent/tea-app/lib/i18n';

import { dateFormatter, downloadCrt } from '../../../../helpers';
import { getWorkflowError } from '../../common';
import { Clip, GridTable, LinkButton, WorkflowDialog } from '../../common/components';
import { DialogBodyLayout } from '../../common/layouts';
import { allActions } from '../actions';
import { namespaceActions } from '../actions/namespaceActions';
import { NamespaceStatus, resourceLimitTypeToText, resourceTypeToUnit, PlatformTypeEnum } from '../constants/Config';
import { Namespace, NamespaceOperator, Project } from '../models';
import { router } from '../router';
import { CreateProjectResourceLimitPanel } from './CreateProjectResourceLimitPanel';
import { RootProps } from './ProjectApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class NamespaceTablePanel extends React.Component<RootProps, {}> {
  state = {
    isShowDialog: false,
    isShowKuctlDialog: false
  };
  render() {
    return (
      <React.Fragment>
        {this._renderTablePanel()}
        {this._renderDeleteNamespaceDialog()}
        {this._renderEditProjectLimitDialog()}
        {this._renderKubectlDialog()}
        <MigarteamespaceDialog {...this.props} />
      </React.Fragment>
    );
  }
  formatResourceLimit(resourceLimit) {
    let resourceLimitKeys = Object.keys(resourceLimit);
    let content = resourceLimitKeys.map((item, index) => (
      <Text parent="p" key={index}>{`${resourceLimitTypeToText[item]}:${
        resourceTypeToUnit[item] === 'MiB'
          ? valueLabels1024(resourceLimit[item], K8SUNIT.Mi)
          : valueLabels1000(resourceLimit[item], K8SUNIT.unit)
      }${resourceTypeToUnit[item]}`}</Text>
    ));
    return (
      <Bubble content={content}>
        <p style={{ display: 'inline-block' }}>{content.filter((item, index) => index < 3)}</p>
      </Bubble>
    );
  }
  private _renderTablePanel() {
    let { actions, namespace, namespaceEdition } = this.props;
    const columns: TableColumn<Namespace>[] = [
      {
        key: 'name',
        header: t('名称'),
        render: x => {
          let disabledOp = x.status.phase === 'Terminating';
          let url = `/tkestack/cluster/sub/list/resource/deployment?rid=1&clusterId=${x.spec.clusterName}&np=${x.spec.namespace}`;
          /// #if project
          url = `/tkestack-project/application/list/resource/deployment?rid=1&clusterId=${x.spec.clusterName}&np=${x.spec.namespace}`;
          /// #endif
          return <Text overflow>{!disabledOp ? <a href={url}>{x.spec.namespace}</a> : x.spec.namespace}</Text>;
        }
      },
      {
        key: 'clusterName',
        header: t('所属集群'),
        render: x => (
          <div>
            <span className="text-overflow">{x.spec.clusterName}</span>
          </div>
        )
      },
      {
        key: 'status',
        header: t('状态'),
        render: x => (
          <React.Fragment>
            <Text theme={NamespaceStatus[x.status.phase]} verticalAlign="middle">
              {x.status.phase || '-'}
            </Text>
            {(x.status.phase === 'Terminating' || x.status.phase === 'Pending') && <Icon type="loding" />}
            {x.status.phase === 'Failed' && (
              <Bubble content={x.status.reason || '-'}>
                <Icon type="error" />
              </Bubble>
            )}
          </React.Fragment>
        )
      },
      {
        key: 'resourceLimit',
        header: t('资源配额'),
        render: x => {
          let disabledOp = x.status.phase === 'Terminating';
          return (
            <React.Fragment>
              {x.spec.hard ? this.formatResourceLimit(x.spec.hard) : '无限制'}
              {!disabledOp && (
                <Icon
                  onClick={() => {
                    this.setState({ isShowDialog: true });
                    actions.namespace.initNamespaceEdition(x);
                    actions.namespace.editNamespaceResourceLimit.start([namespaceEdition]);
                  }}
                  style={{ cursor: 'pointer' }}
                  type="pencil"
                />
              )}
            </React.Fragment>
          );
        }
      },
      {
        key: 'createdTime',
        header: t('创建时间'),
        render: x => (
          <Text parent="div" overflow>
            <span className="text">{dateFormatter(new Date(x.metadata.creationTimestamp), 'YYYY-MM-DD HH:mm:ss')}</span>
          </Text>
        )
      }
    ];

    return (
      <TablePanel
        columns={columns}
        emptyTips={<div className="text-center">{t('您选择的该业务的命名空间为空')}</div>}
        model={namespace}
        action={actions.namespace}
        getOperations={x => this._renderOperationCell(x)}
        operationsWidth={300}
        // isNeedPagination={true}
      />
    );
  }

  private _renderOperationCell(namespace: Namespace) {
    const {
      actions,
      route,
      deleteNamespace,
      namespaceEdition,
      platformType,
      projectDetail,
      userManagedProjects
    } = this.props;
    const urlParams = router.resolve(route);
    let enableOp =
      platformType === PlatformTypeEnum.Manager ||
      (platformType === PlatformTypeEnum.Business &&
        userManagedProjects.list.data.records.find(
          item => item.name === (projectDetail ? projectDetail.metadata.name : null)
        ));
    const matchPerformingWorkflow = (workflow: WorkflowState<Namespace, NamespaceOperator>) => {
      return (
        workflow.operationState === OperationState.Performing &&
        workflow.targets &&
        workflow.targets[0] &&
        workflow.targets[0].id === namespace.id
      );
    };
    const isDeleting = matchPerformingWorkflow(deleteNamespace),
      errTip = <p>{t('当前状态下不可进行该操作')}</p>;

    let disabledOp = namespace.status.phase === 'Terminating';
    let disabledMigartion = namespace.status.phase !== 'Available';
    let disabledCert = namespace.spec.clusterType === 'Imported';

    const renderDeleteButton = () => {
      return (
        <LinkButton
          key={'delete'}
          disabled={isDeleting || disabledOp}
          errorTip={errTip}
          tipDirection="right"
          onClick={() => actions.namespace.deleteNamespace.start([namespace])}
        >
          {t('删除')}
        </LinkButton>
      );
    };

    const renderKubctlConfigButton = () => {
      return (
        <LinkButton
          key={'kubectl'}
          disabled={isDeleting || disabledOp || disabledCert}
          tipDirection="right"
          onClick={() => {
            this.setState({ isShowKuctlDialog: true });
            actions.namespace.select(namespace);
            actions.namespace.namespaceKubectlConfig.applyFilter({
              projectId: projectDetail ? projectDetail.metadata.name : route.queries['projectId'],
              np: namespace.metadata.name
            });
          }}
        >
          {t('查看访问凭证')}
        </LinkButton>
      );
    };

    const renderMigartionButton = () => {
      return (
        <LinkButton
          key={'nigartion'}
          disabled={isDeleting || disabledMigartion}
          tipDirection="right"
          onClick={() => {
            actions.namespace.selects([namespace]);
            actions.namespace.migrateNamesapce.start([]);
          }}
        >
          {t('迁移')}
        </LinkButton>
      );
    };

    let buttons = [];
    buttons.push([renderKubctlConfigButton()]);

    if (enableOp) {
      buttons.push([renderDeleteButton()]);
    }
    if (platformType === PlatformTypeEnum.Manager) {
      buttons.push([renderMigartionButton()]);
    }

    return buttons;
  }
  private _renderEditProjectLimitDialog() {
    const { actions, project, editNamespaceResourceLimit, namespaceEdition, projectDetail } = this.props;
    let isShowDialog = this.state.isShowDialog;

    let parentResourceLimits =
      projectDetail && namespaceEdition.clusterName
        ? projectDetail.spec.clusters[namespaceEdition.clusterName].hard
        : {};

    let failed =
      editNamespaceResourceLimit.operationState === OperationState.Done &&
      !isSuccessWorkflow(editNamespaceResourceLimit);

    const cancel = () => {
      this.setState({ isShowDialog: false });
      actions.namespace.clearEdition();
      if (editNamespaceResourceLimit.operationState === OperationState.Done) {
        actions.namespace.editNamespaceResourceLimit.reset();
      }
      if (editNamespaceResourceLimit.operationState === OperationState.Started) {
        actions.namespace.editNamespaceResourceLimit.cancel();
      }
    };
    return (
      <Modal
        visible={isShowDialog || editNamespaceResourceLimit.operationState !== OperationState.Pending}
        caption={t('编辑资源限制')}
        onClose={() => cancel()}
      >
        <CreateProjectResourceLimitPanel
          parentResourceLimits={parentResourceLimits}
          onCancel={() => cancel()}
          failMessage={failed ? getWorkflowError(editNamespaceResourceLimit) : null}
          resourceLimits={namespaceEdition.resourceLimits}
          onSubmit={resourceLimits => {
            namespaceEdition.resourceLimits = resourceLimits;
            actions.namespace.editNamespaceResourceLimit.start([namespaceEdition], {
              projectId: projectDetail.metadata.name
            });
            actions.namespace.editNamespaceResourceLimit.perform();
            cancel();
          }}
        />
      </Modal>
    );
  }

  private _renderDeleteNamespaceDialog() {
    const { actions, route, deleteNamespace } = this.props;
    return (
      <WorkflowDialog
        caption={t('删除Namespace')}
        workflow={deleteNamespace}
        action={actions.namespace.deleteNamespace}
        targets={deleteNamespace.targets}
        postAction={() => {
          actions.namespace.clearEdition();
        }}
        params={{ projectId: route.queries['projectId'] }}
      >
        <DialogBodyLayout>
          <p className="til">
            <strong className="tip-top">
              {t('确定要删除Namespace {{name}}么？', {
                name: deleteNamespace.targets ? deleteNamespace.targets[0].metadata.name : ''
              })}
            </strong>
          </p>
          <p className="text-danger">{t('删除Namespace将删除该Namespace下所有资源，该操作不可逆，请谨慎操作。')}</p>
        </DialogBodyLayout>
      </WorkflowDialog>
    );
  }
  private _renderKubectlDialog() {
    let {
      namespaceKubectlConfig,
      namespace: { selection },
      userInfo
    } = this.props;
    const cancel = () => {
      this.setState({ isShowKuctlDialog: false });
    };
    let certInfo = namespaceKubectlConfig.object && namespaceKubectlConfig.object.data;
    let clusterId = selection && selection.spec.clusterName;
    let np = selection && selection.spec.namespace;
    let userName = userInfo.object.data ? userInfo.object.data.name : '';
    let kubectlConfig = certInfo ? namespaceActions.getKubectlConfig(certInfo, clusterId, np, userName) : '';
    return (
      <Modal visible={this.state.isShowKuctlDialog} caption={t('访问凭证')} onClose={() => cancel()} size={700}>
        <Modal.Body>
          <FormPanel isNeedCard={false}>
            <FormPanel.Item text label={'Kubeconfig'}>
              <div className="form-unit">
                <div className="rich-textarea hide-number" style={{ width: '100%' }}>
                  <Clip target={'#kubeconfig'} className="copy-btn">
                    {t('复制')}
                  </Clip>
                  <a
                    href="javascript:void(0)"
                    onClick={e => downloadKubeconfig(kubectlConfig, `config`)}
                    className="copy-btn"
                    style={{ right: '50px' }}
                  >
                    {t('下载')}
                  </a>
                  <div className="rich-content">
                    <pre
                      className="rich-text"
                      id="kubeconfig"
                      style={{
                        whiteSpace: 'pre-wrap',
                        overflow: 'auto',
                        height: '300px'
                      }}
                    >
                      {kubectlConfig}
                    </pre>
                  </div>
                </div>
              </div>
            </FormPanel.Item>
          </FormPanel>
          <div
            style={{
              textAlign: 'left',
              borderTop: '1px solid #D1D2D3',
              paddingTop: '10px',
              marginTop: '10px',
              color: '#444'
            }}
          >
            <h3 style={{ marginBottom: '1em' }}>通过Kubectl连接Kubernetes集群操作说明:</h3>
            <p style={{ marginBottom: '5px' }}>
              1. 安装 Kubectl 客户端：从
              <ExternalLink href="https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG.md">
                Kubernetes 版本页面
              </ExternalLink>
              下载最新的 kubectl 客户端，并安装和设置 kubectl 客户端，具体可参考
              <ExternalLink href="https://kubernetes.io/docs/tasks/tools/install-kubectl/">
                安装和设置 kubectl
              </ExternalLink>
              。
            </p>
            <p style={{ marginBottom: '5px' }}>2. 配置 Kubeconfig：</p>
            <ul>
              <li style={{ listStyle: 'disc', marginLeft: '15px' }}>
                <p style={{ marginBottom: '5px' }}>
                  若当前访问客户端尚未配置任何集群的访问凭证，即 ~/.kube/config 内容为空，可直接复制上方 kubeconfig
                  访问凭证内容并粘贴入 ~/.kube/config 中。
                </p>
              </li>
            </ul>
            <p style={{ marginBottom: '5px' }}>
              3. 可执行 kubectl get pod -n {np}
              测试是否可正常访问您的命名空间下的资源。如果无法连接请查看是否已经开启公网访问或内网访问入口，并确保访问客户端在指定的网络环境内。
              如果返回 (Forbidden) 错误，请确保用户具有所在业务相应的权限。
            </p>
          </div>
        </Modal.Body>
        <Modal.Footer>
          <Button type="primary" onClick={cancel}>
            {t('关闭')}
          </Button>
        </Modal.Footer>
      </Modal>
    );
  }
}

function MigarteamespaceDialog(props: RootProps) {
  let [projectSelection, setProjectSelection] = React.useState('');
  let [v_projectSelection, setVProjectSelection] = React.useState(initValidator);
  const { actions, migrateNamesapce, route, project, projectDetail, namespace } = props;
  let failed = migrateNamesapce.operationState === OperationState.Done && !isSuccessWorkflow(migrateNamesapce);

  const cancel = () => {
    actions.namespace.clearSelection();
    if (migrateNamesapce.operationState === OperationState.Done) {
      actions.namespace.migrateNamesapce.reset();
    }
    if (migrateNamesapce.operationState === OperationState.Started) {
      actions.namespace.migrateNamesapce.cancel();
    }
  };
  return (
    <Modal
      visible={migrateNamesapce.operationState !== OperationState.Pending}
      caption={t('迁移命名空间')}
      onClose={() => cancel()}
    >
      <Modal.Body>
        <FormPanel isNeedCard={false}>
          <FormPanel.Item label={t('当前业务')}>
            <FormPanel.Text>
              {projectDetail ? `${projectDetail.metadata.name}(${projectDetail.spec.displayName})` : null}
            </FormPanel.Text>
          </FormPanel.Item>
          <FormPanel.Item
            label={'目标业务'}
            validator={v_projectSelection}
            select={{
              value: projectSelection,
              model: project,
              onChange: value => {
                setProjectSelection(value);
              },
              displayField: (r: Project) => `${r.metadata.name}(${r.spec.displayName})`,
              valueField: (r: Project) => r.metadata.name
            }}
          ></FormPanel.Item>
        </FormPanel>
      </Modal.Body>
      <Modal.Footer>
        <Button
          type="primary"
          style={{ margin: '0px 5px 0px 40px' }}
          onClick={() => {
            if (projectSelection === '') {
              setVProjectSelection({
                status: 2,
                message: t('目标业务不能为空')
              });
            } else if (projectDetail && projectSelection === projectDetail.metadata.name) {
              setVProjectSelection({
                status: 2,
                message: t('目标业务不能和当前业务一致')
              });
            } else {
              setVProjectSelection({
                status: 1,
                message: t('')
              });
              actions.namespace.migrateNamesapce.start(namespace.selections, {
                projectId: projectDetail ? projectDetail.metadata.name : route.queries['rid'],
                desProjectId: projectSelection
              });
              actions.namespace.migrateNamesapce.perform();
            }
          }}
        >
          {failed ? t('重试') : t('完成')}
        </Button>
        <Button
          type="weak"
          onClick={() => {
            cancel();
          }}
        >
          {t('取消')}
        </Button>
        {failed ? (
          <Alert
            type="error"
            style={{ display: 'inline-block', marginLeft: '20px', marginBottom: '0px', maxWidth: '750px' }}
          >
            {getWorkflowError(migrateNamesapce)}
          </Alert>
        ) : (
          <noscript />
        )}
      </Modal.Footer>
    </Modal>
  );
}
