import * as React from 'react';
import { connect } from 'react-redux';
import { getWorkflowError } from '../../common';

import { Bubble, Button, Icon, Modal, Pagination, TableColumn, Text } from '@tea/component';
import { bindActionCreators } from '@tencent/qcloud-lib';
import { isSuccessWorkflow, OperationState, WorkflowState } from '@tencent/qcloud-redux-workflow';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { dateFormatter } from '../../../../helpers';
import { GridTable, LinkButton, ResourceList, TipDialog, WorkflowDialog } from '../../common/components';
import { DialogBodyLayout } from '../../common/layouts';
import { allActions } from '../actions';
import { NamespaceStatus, resourceLimitTypeToText, resourceTypeToUnit } from '../constants/Config';
import { K8SUNIT, valueLabels1024, valueLabels1000 } from '@helper/k8sUnitUtil';
import { Namespace, NamespaceOperator } from '../models';
import { router } from '../router';
import { CreateProjectResourceLimitPanel } from './CreateProjectResourceLimitPanel';
import { RootProps } from './ProjectApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class NamespaceTablePanel extends React.Component<RootProps, {}> {
  state = {
    isShowDialog: false
  };
  render() {
    return (
      <React.Fragment>
        {this._renderTablePanel()}
        {this._renderDeleteNamespaceDialog()}
        {this._renderEditProjectLimitDialog()}
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
                ? x.metadata.name
                    .split('-')
                    .splice(2)
                    .join('-')
                : x.metadata.name
                    .split('-')
                    .splice(1)
                    .join('-')}
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
      },
      { key: 'operation', header: t('操作'), render: x => this._renderOperationCell(x) }
    ];

    return (
      <GridTable
        columns={columns}
        emptyTips={<div className="text-center">{t('您选择的该业务的命名空间为空')}</div>}
        listModel={{
          list: namespace.list,
          query: namespace.query
        }}
        actionOptions={actions.namespace}
        isNeedPagination={true}
      />
    );
  }

  private _renderOperationCell(namespace: Namespace) {
    const { actions, route, deleteNamespace, namespaceEdition } = this.props;
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

    return <div>{[renderDeleteButton(), renderEditResourceLimitButton()]}</div>;
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
}
