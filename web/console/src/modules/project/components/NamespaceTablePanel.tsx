import * as React from 'react';
import { connect } from 'react-redux';

import { K8SUNIT, valueLabels1000, valueLabels1024 } from '@helper/k8sUnitUtil';
import { Bubble, Icon, Modal, TableColumn, Text, Button, Alert, ExternalLink } from '@tea/component';
import { TablePanel, FormPanel } from '@tencent/ff-component';
import { bindActionCreators, isSuccessWorkflow, OperationState, WorkflowState } from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';

import { dateFormatter, downloadCrt } from '../../../../helpers';
import { getWorkflowError } from '../../common';
import { GridTable, LinkButton, WorkflowDialog, Clip } from '../../common/components';
import { DialogBodyLayout } from '../../common/layouts';
import { allActions } from '../actions';
import { NamespaceStatus, resourceLimitTypeToText, resourceTypeToUnit } from '../constants/Config';
import { Namespace, NamespaceOperator, Project } from '../models';
import { router } from '../router';
import { CreateProjectResourceLimitPanel } from './CreateProjectResourceLimitPanel';
import { RootProps } from './ProjectApp';
import { initValidator } from '@tencent/ff-validator';
import { downloadKubeconfig } from '@helper/downloadText';
import { namespaceActions } from '../actions/namespaceActions';

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
    return <Bubble content={content}>{content.filter((item, index) => index < 3)}</Bubble>;
  }
  private _renderTablePanel() {
    let { actions, namespace } = this.props;

    const columns: TableColumn<Namespace>[] = [
      {
        key: 'name',
        header: t('名称'),
        render: x => (
          <div>
            <span className="text-overflow">
              {x.metadata.name.includes('cls')
                ? x.metadata.name.split('-').splice(2).join('-')
                : x.metadata.name.split('-').splice(1).join('-')}
            </span>
          </div>
        )
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
        render: x => <React.Fragment>{x.spec.hard ? this.formatResourceLimit(x.spec.hard) : '无限制'}</React.Fragment>
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
    const { actions, route, deleteNamespace, namespaceEdition, projectDetail } = this.props;
    const urlParams = router.resolve(route);

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

    const renderEditResourceLimitButton = () => {
      return (
        <LinkButton
          key={'edit'}
          disabled={isDeleting || disabledOp}
          tipDirection="right"
          onClick={() => {
            this.setState({ isShowDialog: true });
            actions.namespace.initNamespaceEdition(namespace);
            actions.namespace.editNamespaceResourceLimit.start([namespaceEdition]);
          }}
        >
          {t('编辑资源限制')}
        </LinkButton>
      );
    };

    const renderKubctlConfigButton = () => {
      return (
        <LinkButton
          key={'kubectl'}
          disabled={isDeleting || disabledOp}
          tipDirection="right"
          onClick={() => {
            this.setState({ isShowKuctlDialog: true });
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
    buttons.push([
      renderDeleteButton(),
      renderEditResourceLimitButton(),
      renderKubctlConfigButton(),
      renderMigartionButton()
    ]);
    return buttons;
  }
  private _renderEditProjectLimitDialog() {
    const { actions, project, editNamespaceResourceLimit, namespaceEdition } = this.props;
    let isShowDialog = this.state.isShowDialog;
    let projectSelection = project.selections[0] ? project.selections[0] : null;

    let parentResourceLimits =
      projectSelection && namespaceEdition.clusterName
        ? projectSelection.spec.clusters[namespaceEdition.clusterName].hard
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
              projectId: projectSelection.metadata.name
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
      namespace: { selection }
    } = this.props;
    const cancel = () => {
      this.setState({ isShowKuctlDialog: false });
    };
    let certInfo = namespaceKubectlConfig.object && namespaceKubectlConfig.object.data;
    let clusterId = selection && selection.spec.clusterName;
    let kubectlConfig = certInfo ? namespaceActions.getKubectlConfig(certInfo, clusterId) : '';
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
                    onClick={e => downloadKubeconfig(kubectlConfig, `${clusterId}&{}-config`)}
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
              <li style={{ listStyle: 'disc', marginLeft: '15px' }}>
                <p style={{ marginBottom: '5px' }}>
                  若当前访问客户端已配置了其他集群的访问凭证，你可下载上方 kubeconfig
                  至指定位置，并执行以下指令以合并多个集群的 config。
                </p>
                <div className="rich-textarea hide-number" style={{ width: '100%' }}>
                  <div className="rich-content">
                    <Clip target={'#kubeconfig-merge'} className="copy-btn">
                      复制
                    </Clip>
                    <pre
                      className="rich-text"
                      id="kubeconfig-merge"
                      style={{
                        whiteSpace: 'pre-wrap',
                        overflow: 'auto'
                      }}
                    >
                      KUBECONFIG=~/.kube/config:~/Downloads/{clusterId}-config kubectl config view --merge --flatten
                      &gt; ~/.kube/config
                      <br />
                      export KUBECONFIG=~/.kube/config
                    </pre>
                  </div>
                </div>
                <p style={{ marginBottom: '5px' }}>
                  其中，~/Downloads/{clusterId}-config 为本集群的 kubeconfig
                  的文件路径，请替换为下载至本地后的实际路径。
                </p>
              </li>
            </ul>
            <p style={{ marginBottom: '5px' }}>3. 访问 Kubernetes 集群：</p>
            <ul>
              <li style={{ marginLeft: '15px' }}>
                <p style={{ marginBottom: '5px' }}>
                  完成 kubeconfig 配置后，执行以下指令查看并切换 context 以访问本集群：
                </p>
                <div className="rich-textarea hide-number" style={{ width: '100%' }}>
                  <div className="rich-content">
                    <Clip target={'#kubeconfig-visit'} className="copy-btn">
                      复制
                    </Clip>
                    <pre
                      className="rich-text"
                      id="kubeconfig-visit"
                      style={{
                        whiteSpace: 'pre-wrap',
                        overflow: 'auto'
                      }}
                    >
                      kubectl config get-contexts
                      <br />
                      kubectl config use-context {clusterId}-context-default
                    </pre>
                  </div>
                </div>
                <p style={{ marginBottom: '5px' }}>
                  而后可执行 kubectl get node
                  测试是否可正常访问集群。如果无法连接请查看是否已经开启公网访问或内网访问入口，并确保访问客户端在指定的网络环境内。
                </p>
              </li>
            </ul>
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
