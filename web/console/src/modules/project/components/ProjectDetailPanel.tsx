import * as React from 'react';

import { K8SUNIT, valueLabels1000, valueLabels1024 } from '@helper/k8sUnitUtil';
import { Bubble, Button, Icon, Modal, StatusTip, Table, TableColumn, Text } from '@tea/component';
import { deepClone } from '@tencent/qcloud-lib';
import { isSuccessWorkflow, OperationState } from '@tencent/qcloud-redux-workflow';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { autotip } from '@tencent/tea-component/lib/table/addons';

import { dateFormatter } from '../../../../helpers';
import { getWorkflowError } from '../../common';
import { FormItem, FormPanel, TablePanel, WorkflowDialog } from '../../common/components';
import { DialogBodyLayout } from '../../common/layouts';
import { resourceLimitTypeToText, resourceTypeToUnit } from '../constants/Config';
import { ProjectResourceLimit } from '../models/Project';
import { CreateProjectResourceLimitPanel } from './CreateProjectResourceLimitPanel';
import { EditProjectManagerPanel } from './EditProjectManagerPanel';
import { EditProjectNamePanel } from './EditProjectNamePanel';
import { RootProps } from './ProjectApp';

export class ProjectDetailPanel extends React.Component<RootProps, {}> {
  state = {
    currentClusterIndex: 0,
    isShowDialog: false
  };
  componentDidMount() {
    let { actions, route } = this.props;
    actions.manager.fetch();
  }
  formatManager(managers) {
    if (managers && managers.length) {
      return managers.map((m, index) => {
        return managers.length - 1 === index ? (
          <span key={index} className="text-overflow">
            {m}
          </span>
        ) : (
          <p key={index} className="text-overflow">
            {m}
          </p>
        );
      });
    } else {
      return '-';
    }
  }

  formatResourceLimit(resourceLimit) {
    let resourceLimitKeys = resourceLimit ? Object.keys(resourceLimit) : [];
    let content = resourceLimitKeys.map((item, index) => (
      <Text parent="p" key={index}>{`${resourceLimitTypeToText[item]}:${
        resourceTypeToUnit[item] === 'MiB'
          ? valueLabels1024(resourceLimit[item], K8SUNIT.Mi)
          : valueLabels1000(resourceLimit[item], K8SUNIT.unit)
      }${resourceTypeToUnit[item]}`}</Text>
    ));
    return resourceLimit ? (
      <Bubble content={content}>{content.filter((item, index) => index < 2)}</Bubble>
    ) : (
      <Text parent="p">{t('无限制')}</Text>
    );
  }

  render() {
    let { actions, route, project } = this.props,
      projectItem = project.selections[0] ? project.selections[0] : null;

    return projectItem ? (
      <FormPanel title={t('基本信息')}>
        <FormPanel.Item
          label={t('业务名称')}
          text
          textProps={{
            onEdit: () => {
              actions.project.initEdition(projectItem);
              actions.project.editProjectName.start([]);
            }
          }}
        >
          {projectItem.spec.displayName}
        </FormPanel.Item>
        <FormPanel.Item
          label={t('成员')}
          text
          textProps={{
            onEdit: () => {
              actions.project.initEdition(projectItem);
              actions.project.editProjectManager.start([]);
            }
          }}
        >
          {this.formatManager(projectItem.spec.members)}
        </FormPanel.Item>
        <FormPanel.Item label={t('资源限制')}>{this._renderTablePanel()}</FormPanel.Item>
        <FormPanel.Item text label={t('创建时间')}>
          {dateFormatter(new Date(projectItem.metadata.creationTimestamp), 'YYYY-MM-DD HH:mm:ss')}
        </FormPanel.Item>
        {this._renderEditProjectNameDialog()}
        {this._renderEditProjectManagerDialog()}
        {this._renderEditProjectLimitDialog()}
      </FormPanel>
    ) : (
      <noscript />
    );
  }
  private _renderTablePanel() {
    let { actions, namespace, project } = this.props,
      projectItem = project.selections[0] ? project.selections[0] : null;
    let clusterKeys = projectItem && projectItem.spec.clusters ? Object.keys(projectItem.spec.clusters) : [];
    let finalClusterList = clusterKeys.map(item => {
      return {
        name: item,
        hard: projectItem.spec.clusters[item].hard
      };
    });
    const columns: TableColumn<{ name: string; hard: any }>[] = [
      {
        key: 'name',
        header: t('名称'),
        width: '20%',
        render: x => (
          <div>
            <span className="text-overflow">{x.name}</span>
          </div>
        )
      },
      {
        width: '65%',
        key: 'resourceLimit',
        header: t('集群配额'),
        render: x => <React.Fragment>{this.formatResourceLimit(x.hard)}</React.Fragment>
      },
      {
        width: '15%',
        key: 'operation',
        header: t('操作'),
        render: (x, recordkey, recordIndex) => (
          <Button
            type="link"
            onClick={() => {
              actions.project.initEdition(projectItem);
              this.setState({
                isShowDialog: true,
                currentClusterIndex: recordIndex
              });
            }}
          >
            {t('编辑')}
          </Button>
        )
      }
    ];

    return (
      <div style={{ width: 500 }}>
        <Table
          columns={columns}
          recordKey={'name'}
          records={finalClusterList}
          addons={[
            autotip({
              emptyText: (
                <StatusTip
                  status="empty"
                  emptyText={<div className="text-center">{t('该业务没有集群配额限制')}</div>}
                />
              )
            })
          ]}
        />
      </div>
    );
  }
  private _renderEditProjectNameDialog() {
    const { actions, route, editProjectName, projectEdition } = this.props;
    return (
      <WorkflowDialog
        caption={t('编辑名称')}
        workflow={editProjectName}
        action={actions.project.editProjectName}
        targets={[projectEdition]}
        postAction={() => {
          actions.project.clearEdition();
        }}
        params={{}}
      >
        <DialogBodyLayout>
          <EditProjectNamePanel {...this.props} />
        </DialogBodyLayout>
      </WorkflowDialog>
    );
  }

  private _renderEditProjectManagerDialog() {
    const { actions, route, editProjectManager, projectEdition } = this.props;
    return (
      <WorkflowDialog
        caption={t('编辑成员')}
        workflow={editProjectManager}
        action={actions.project.editProjectManager}
        targets={[projectEdition]}
        postAction={() => {
          actions.project.clearEdition();
        }}
        params={{}}
        width={600}
      >
        <DialogBodyLayout>
          <EditProjectManagerPanel {...this.props} />
        </DialogBodyLayout>
      </WorkflowDialog>
    );
  }

  private _renderEditProjectLimitDialog() {
    const { actions, project, projectEdition, editProjecResourceLimit } = this.props;
    let { currentClusterIndex, isShowDialog } = this.state;
    let parentProjectSelection = projectEdition.parentProject
      ? project.list.data.records.find(item => item.metadata.name === projectEdition.parentProject)
      : null;
    let clusterName = projectEdition.clusters.length ? projectEdition.clusters[currentClusterIndex].name : '-';

    let parentResourceLimits =
      parentProjectSelection && clusterName ? parentProjectSelection.spec.clusters[clusterName].hard : {};

    let failed =
      editProjecResourceLimit.operationState === OperationState.Done && !isSuccessWorkflow(editProjecResourceLimit);

    const cancel = () => {
      this.setState({ isShowDialog: false, currentClusterIndex: 0 });
      actions.project.clearEdition();

      if (editProjecResourceLimit.operationState === OperationState.Done) {
        actions.project.editProjecResourceLimit.reset();
      }
      if (editProjecResourceLimit.operationState === OperationState.Started) {
        actions.project.editProjecResourceLimit.cancel();
      }
    };
    return (
      <Modal
        visible={isShowDialog || editProjecResourceLimit.operationState !== OperationState.Pending}
        caption={t('编辑资源限制')}
        onClose={() => cancel()}
      >
        <CreateProjectResourceLimitPanel
          parentResourceLimits={parentResourceLimits}
          onCancel={() => cancel()}
          failMessage={failed ? getWorkflowError(editProjecResourceLimit) : null}
          resourceLimits={projectEdition.clusters[currentClusterIndex].resourceLimits}
          onSubmit={resourceLimits => {
            if (projectEdition.clusters[currentClusterIndex]) {
              projectEdition.clusters[currentClusterIndex] = Object.assign(
                {},
                projectEdition.clusters[currentClusterIndex],
                {
                  resourceLimits
                }
              );
            }
            actions.project.editProjecResourceLimit.start([projectEdition]);
            actions.project.editProjecResourceLimit.perform();
          }}
        />
      </Modal>
    );
  }
}
