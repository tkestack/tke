import * as React from 'react';
import { RootProps } from './NotifyApp';
import { WorkflowDialog } from '../../common/components';
import { t } from '@tencent/tea-app/lib/i18n';
import { router } from '../router';
import { resourceConfig } from '../../../../config';
const rc = resourceConfig();
export class DeleteResourceDialog extends React.Component<RootProps, {}> {
  render() {
    let { route, resourceDeleteWorkflow, actions } = this.props;
    let urlParams = router.resolve(route);
    let resource = this.props[urlParams.resourceName] || this.props.channel;
    let nameStr = resource.selections.map(item => `${item.metadata.name}(${item.spec.displayName})`).join(', ');
    return (
      <WorkflowDialog
        caption={t('删除{{headTitle}}', rc[urlParams.resourceName])}
        workflow={resourceDeleteWorkflow}
        action={actions.workflow.deleteResource}
        params={{}}
        targets={resource.selections}
      >
        <div style={{ fontSize: '14px', lineHeight: '20px' }}>
          <div style={{ fontSize: '14px', lineHeight: '20px' }}>
            <div className="docker-dialog jiqun">
              <p>
                <strong>{t('您确定要删除 {{nameStr}} 吗？', { nameStr })}</strong>
              </p>
            </div>
          </div>
        </div>
      </WorkflowDialog>
    );
  }
}
