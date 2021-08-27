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

import { Bubble, Drawer, TableColumn, Text } from '@tea/component';
import { bindActionCreators, OperationState, WorkflowState } from '@tencent/ff-redux';
import ChartPanel from '@tencent/tchart';
import { t } from '@tencent/tea-app/lib/i18n';

import { dateFormatter } from '../../../../helpers';
import { projectFields } from '../../cluster/models/MonitorPanel';
import { GridTable, LinkButton, WorkflowDialog } from '../../common/components';
import { DialogBodyLayout } from '../../common/layouts';
import { allActions } from '../actions';
import { projectActions } from '../actions/projectActions';
import { projectStatus } from '../constants/Config';
import { Project } from '../models';
import { router } from '../router';
import { EditProjectManagerPanel } from './EditProjectManagerPanel';
import { EditProjectNamePanel } from './EditProjectNamePanel';
import { RootProps } from './ProjectApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class ProjectTablePanel extends React.Component<RootProps, any> {
  state = {
    monitorPanelProps: undefined,
    isShowMonitor: false
  };

  render() {
    return (
      <React.Fragment>
        {this._renderTablePanel()}
        {this._renderDeleteProjectDialog()}
        {this._renderEditProjectNameDialog()}
        {this.renderMonitor()}
      </React.Fragment>
    );
  }

  handleDeleteProject(project: Project) {
    let { actions } = this.props;
    actions.project.deleteProject.start([project]);
  }

  formatManager(managers) {
    if (managers) {
      return managers.map((m, index) => {
        return (
          <p key={index} className="text-overflow">
            {m.username}
          </p>
        );
      });
    }
  }

  private _renderTablePanel() {
    let { actions, project, projectUserInfo } = this.props;
    const columns: TableColumn<Project>[] = [
      {
        key: 'name',
        header: t('ID/名称'),
        render: x => (
          <div>
            {x.status.phase === 'Terminating' ? (
              <Text parent="div" overflow>
                {x.metadata.name}
              </Text>
            ) : (
              <React.Fragment>
                <Text parent="div" overflow>
                  <a
                    href="javascript:;"
                    onClick={e => {
                      router.navigate({ sub: 'detail', tab: 'info' }, { projectId: x.metadata.name });
                    }}
                  >
                    {x.metadata.name}
                  </a>
                </Text>
                <div className="sl-editor-name">
                  <span className="text-overflow m-width" title={x.spec.displayName}>
                    {x.spec.displayName || t('未命名')}
                  </span>
                  <span className="hover-icon">
                    <a
                      href="javascript:;"
                      className="pencil-icon hover-icon"
                      onClick={() => {
                        actions.project.initEdition(x);
                        actions.project.editProjectName.start([]);
                      }}
                    />
                  </span>
                </div>
              </React.Fragment>
            )}
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
        key: 'parentProject',
        header: t('上级业务'),
        render: x => (
          <Text parent="div" overflow>
            {x.spec.parentProjectName ? x.spec.parentProjectName : '无'}
          </Text>
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
        key: 'managers',
        header: t('成员'),
        render: x => {
          let { projectUserInfo } = this.props;
          let members =
            projectUserInfo.object.data && projectUserInfo.object.data[x.metadata.name]
              ? projectUserInfo.object.data[x.metadata.name]
              : [];
          return (
            <div>
              <Bubble placement="left" content={this.formatManager(members) || null}>
                <span className="text">
                  {this.formatManager(members.length ? members.slice(0, 1) : [])}
                  <Text parent="div" overflow>
                    {members && members.length > 1 ? '...' : ''}
                  </Text>
                </span>
              </Bubble>
              {x.status.phase === 'Terminating' ? (
                <noscript />
              ) : (
                <span>
                  <a
                    href="javascript:;"
                    className="pencil-icon"
                    onClick={e => {
                      router.navigate({ sub: 'detail', tab: 'member' }, { projectId: x.metadata.name });
                    }}
                  />
                </span>
              )}
            </div>
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
      <GridTable
        columns={columns}
        emptyTips={<div className="text-center">{t('业务列表为空')}</div>}
        listModel={{
          list: project.list,
          query: project.query
        }}
        actionOptions={actions.project}
      />
    );
  }

  private _renderOperationCell(project: Project) {
    const { deleteProject } = this.props;

    const matchPerformingWorkflow = (workflow: WorkflowState<Project, string>) => {
      return (
        workflow.operationState === OperationState.Performing &&
        workflow.targets &&
        workflow.targets[0] &&
        workflow.targets[0].id === project.id
      );
    };

    const isDeleting = matchPerformingWorkflow(deleteProject);

    const errTip = <p>{t('当前状态下不可进行该操作')}</p>;

    const renderDeleteButton = () => {
      return (
        <LinkButton
          disabled={isDeleting || project.status.phase === 'Terminating'}
          errorTip={errTip}
          tipDirection="right"
          onClick={() => this.handleDeleteProject(project)}
        >
          {t('删除')}
        </LinkButton>
      );
    };

    return <div>{renderDeleteButton()}</div>;
  }

  private _renderEditProjectNameDialog() {
    const { actions, editProjectName, projectEdition } = this.props;
    return (
      <WorkflowDialog
        caption={t('编辑名称')}
        workflow={editProjectName}
        action={actions.project.editProjectName}
        targets={[projectEdition]}
        validateAction={() => {
          return projectActions._validateDisplayName(projectEdition.displayName).status === 1;
        }}
        preAction={() => {
          actions.project.validateDisplayName(projectEdition.displayName);
        }}
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

  private _renderDeleteProjectDialog() {
    const { actions, deleteProject } = this.props;
    return (
      <WorkflowDialog
        caption={t('删除业务')}
        workflow={deleteProject}
        action={actions.project.deleteProject}
        targets={deleteProject.targets}
        postAction={() => {}}
        params={{}}
        confirmMode={{
          label: t('业务Id'),
          value: deleteProject.targets ? deleteProject.targets[0].metadata.name : ''
        }}
      >
        <DialogBodyLayout>
          <p className="til">
            <strong className="tip-top">
              {t('确定要删除业务{{displayName}}({{name}})么？', {
                displayName: deleteProject.targets ? deleteProject.targets[0].spec.displayName : '',
                name: deleteProject.targets ? deleteProject.targets[0].id : ''
              })}
            </strong>
          </p>
          <p className="text-danger">{t('删除业务将删除该业务下所有资源，该操作不可逆，请谨慎操作。')}</p>
        </DialogBodyLayout>
      </WorkflowDialog>
    );
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
}
