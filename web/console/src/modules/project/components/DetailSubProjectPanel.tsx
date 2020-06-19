import * as React from 'react';
import { connect } from 'react-redux';

import { K8SUNIT, valueLabels1000, valueLabels1024 } from '@helper/k8sUnitUtil';
import { Alert, Bubble, Button, Drawer, Modal, StatusTip, Table, TableColumn, Text } from '@tea/component';
import { FormPanel, TablePanel } from '@tencent/ff-component';
import { bindActionCreators, isSuccessWorkflow, OperationState, WorkflowState } from '@tencent/ff-redux';
import ChartPanel from '@tencent/tchart';
import { t } from '@tencent/tea-app/lib/i18n';

import { dateFormatter } from '../../../../helpers';
import { projectFields } from '../../cluster/models/MonitorPanel';
import { getWorkflowError } from '../../common';
import { allActions } from '../actions';
import { projectStatus, resourceLimitTypeToText, resourceTypeToUnit, PlatformTypeEnum } from '../constants/Config';
import { Project } from '../models';
import { RootProps } from './ProjectApp';
import { SelectExistProjectDialog } from './SelectExistProjectDialog';
import { router } from '../router';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class DetailSubProjectPanel extends React.Component<RootProps, any> {
  state = {
    monitorPanelProps: undefined,
    isShowMonitor: false
  };
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
    return (
      <React.Fragment>
        {this._renderTablePanel()}
        {this.renderMonitor()}
        {this._renderDeleteProjectLimitDialog()}
        <SelectExistProjectDialog />
      </React.Fragment>
    );
  }

  private _renderTablePanel() {
    let { actions, detailProject } = this.props;
    let hostPrefix = 'tkestack';
    /// #if project
    hostPrefix = 'tkestack-project';
    /// #endif

    const columns: TableColumn<Project>[] = [
      {
        key: 'name',
        header: t('ID/名称'),
        render: x => (
          <div>
            <Text parent="div" overflow>
              <a href={`/${hostPrefix}/project/detail/info?projectId=${x.metadata.name}`}>
                {`${x.metadata.name}(${x.spec.displayName ? x.spec.displayName : '未命名'})`}
              </a>
            </Text>
          </div>
        )
      },
      {
        key: 'monitor',
        header: t('监控'),
        width: '10%',
        render: x => (
          <div>
            <p className="text-overflow m-width">
              <i
                className="dosage-icon"
                style={{ cursor: 'pointer' }}
                data-monitor
                data-title={t('查看监控')}
                onClick={() => {
                  this._handleMonitor(x);
                }}
              />
            </p>
          </div>
        )
      },

      {
        key: 'phase',
        header: t('状态'),
        render: x => (
          <React.Fragment>
            <Text parent="div" overflow theme={projectStatus[x.status.phase]}>
              {x.status.phase}
            </Text>
          </React.Fragment>
        )
      },
      {
        width: '25%',
        key: 'resourceLimit',
        header: t('集群配额'),
        render: x => {
          let clusterKeys = x && x.spec.clusters ? Object.keys(x.spec.clusters) : [];
          let content = clusterKeys.map((key, index) => (
            <p key={index}>{this.formatResourceLimit(x.spec.clusters[key].hard)}</p>
          ));
          return (
            <Bubble placement={'top-start'} content={content}>
              <>{content.slice().splice(0, 1)}</>
            </Bubble>
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
      },
      {
        key: 'operation',
        header: t('操作'),
        width: '18%',
        render: x => this._renderOperationCell(x)
      }
    ];

    return (
      <TablePanel
        columns={columns}
        emptyTips={<div className="text-center">{t('子业务列表为空')}</div>}
        model={detailProject}
        action={actions.detail.project}
      />
    );
  }

  private _renderOperationCell(project: Project) {
    const { actions, platformType } = this.props;
    let enableOp = platformType === PlatformTypeEnum.Manager;
    const renderDeleteButton = () => {
      return (
        <Button
          type={'link'}
          onClick={() => {
            actions.detail.project.select(project);
            actions.project.deleteParentProject.start([]);
          }}
        >
          {t('解除')}
        </Button>
      );
    };

    return <div>{enableOp ? renderDeleteButton() : null}</div>;
  }

  private _handleMonitor(project: Project) {
    this.setState({
      isShowMonitor: true,
      monitorPanelProps: {
        title: project.metadata.name,
        subTitle: project.spec.displayName,
        tables: [
          {
            table: 'k8s_project',
            fields: projectFields,
            conditions: [['project_name', '=', project.metadata.name]]
          }
        ],
        groupBy: []
      }
    });
  }

  private renderMonitor() {
    return (
      <Drawer
        visible={this.state.isShowMonitor}
        title={(this.state.monitorPanelProps && this.state.monitorPanelProps.title) || ''}
        subTitle={(this.state.monitorPanelProps && this.state.monitorPanelProps.subTitle) || ''}
        onClose={() => this.setState({ isShowMonitor: false })}
        outerClickClosable={true}
        placement={'right'}
        size={'l'}
        style={{ zIndex: 4 }}
        // style={{ width: 600 }}
      >
        {this.state.monitorPanelProps && (
          <ChartPanel
            tables={this.state.monitorPanelProps.tables}
            groupBy={this.state.monitorPanelProps.groupBy}
            height={250}
          />
        )}
      </Drawer>
    );
  }
  private _renderDeleteProjectLimitDialog() {
    const { actions, detailProject, projectEdition, deleteParentProject } = this.props;

    let failed = deleteParentProject.operationState === OperationState.Done && !isSuccessWorkflow(deleteParentProject);

    const cancel = () => {
      actions.detail.project.clearSelection();

      if (deleteParentProject.operationState === OperationState.Done) {
        actions.project.deleteParentProject.reset();
      }
      if (deleteParentProject.operationState === OperationState.Started) {
        actions.project.deleteParentProject.cancel();
      }
    };
    return (
      <Modal
        visible={deleteParentProject.operationState !== OperationState.Pending}
        caption={t('解除业务关联')}
        onClose={() => cancel()}
      >
        <Modal.Body>
          <FormPanel.Text>
            {t('确定要删除业务{{projectId}}与父业务的关联么？', {
              projectId: detailProject.selection ? detailProject.selection.metadata.name : ''
            })}
          </FormPanel.Text>
        </Modal.Body>
        <Modal.Footer>
          <Button
            type="primary"
            style={{ margin: '0px 5px 0px 40px' }}
            onClick={() => {
              actions.project.deleteParentProject.start([detailProject.selection]);
              actions.project.deleteParentProject.perform();
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
              {getWorkflowError(deleteParentProject)}
            </Alert>
          ) : (
            <noscript />
          )}
        </Modal.Footer>
      </Modal>
    );
  }
}
