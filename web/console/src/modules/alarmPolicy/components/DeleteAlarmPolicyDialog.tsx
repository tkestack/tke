import { t } from '@tencent/tea-app/lib/i18n';
import { Table } from '@tencent/tea-component';
import * as React from 'react';
import { WorkflowDialog } from '../../common/components';
import { AlarmPolicyTablePanel } from './AlarmPolicyTablePanel';
export class DeleteAlarmPolicyDialog extends AlarmPolicyTablePanel {
  render() {
    let { route, alarmPolicyDeleteWorkflow, actions, alarmPolicy, regionSelection, cluster } = this.props,
      regionId = route.queries['rid'] || regionSelection.value,
      clusterId = route.queries['clusterId'] || (cluster.selection ? cluster.selection.metadata.name : '');
    return (
      <WorkflowDialog
        width={800}
        caption={t('删除告警设置')}
        workflow={alarmPolicyDeleteWorkflow}
        action={actions.workflow.deleteAlarmPolicy}
        params={{ regionId, clusterId }}
        targets={alarmPolicy.selections}
      >
        <div style={{ fontSize: '14px', lineHeight: '20px' }}>
          <div style={{ fontSize: '14px', lineHeight: '20px' }}>
            <div className="docker-dialog jiqun">
              <p>
                <strong style={{ wordBreak: 'break-all' }}>{t('您确定要删除以下告警设置吗？')}</strong>
              </p>
              <Table bordered columns={this.getColumns()} records={alarmPolicy.selections} />
            </div>
          </div>
        </div>
      </WorkflowDialog>
    );
  }
}
