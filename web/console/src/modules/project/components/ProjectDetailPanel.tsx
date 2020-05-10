import * as React from 'react';

import { K8SUNIT, valueLabels1000, valueLabels1024 } from '@helper/k8sUnitUtil';
import { Bubble, Button, Modal, StatusTip, Table, TableColumn, Text } from '@tea/component';
import { FormPanel } from '@tencent/ff-component';
import { t } from '@tencent/tea-app/lib/i18n';
import { autotip } from '@tencent/tea-component/lib/table/addons';

import { dateFormatter } from '../../../../helpers';
import { getWorkflowError } from '../../common';
import { WorkflowDialog } from '../../common/components';
import { DialogBodyLayout } from '../../common/layouts';
import { resourceLimitTypeToText, resourceTypeToUnit } from '../constants/Config';
import { EditProjectManagerPanel } from './EditProjectManagerPanel';
import { EditProjectNamePanel } from './EditProjectNamePanel';
import { RootProps } from './ProjectApp';
import { ProjectDetailResourcePanel } from './ProjectDetailResourcePanel';

export class ProjectDetailPanel extends React.Component<RootProps, {}> {
  componentDidMount() {
    let { actions, cluster } = this.props;
    if (cluster.list.data.recordCount === 0) {
      actions.cluster.applyFilter({});
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
        <ProjectDetailResourcePanel {...this.props} />
        <FormPanel.Item text label={t('创建时间')}>
          {dateFormatter(new Date(projectItem.metadata.creationTimestamp), 'YYYY-MM-DD HH:mm:ss')}
        </FormPanel.Item>
        {this._renderEditProjectNameDialog()}
        {this._renderEditProjectManagerDialog()}
      </FormPanel>
    ) : (
      <noscript />
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
}
