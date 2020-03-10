import * as React from 'react';
import { connect } from 'react-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators, insertCSS } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { RootProps } from '../RoleApp';
import { WorkflowDialog } from '../../../../common/components';
import { PolicyAssociatePanel } from './PolicyAssociatePanel';
const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

interface WorkflowDialogProps extends RootProps {
  onPostCancel?: () => void;
}

@connect(state => state, mapDispatchToProps)
export class PolicyAssociateWorkflowDialog extends React.Component<WorkflowDialogProps, {}> {

  render() {
    const {
      actions,
      policyAssociation,
      policyFilter,
      associatePolicyWorkflow
    } = this.props;
    const { onPostCancel = undefined } = this.props;
    return (
      <WorkflowDialog
        caption={t('关联策略')}
        workflow={associatePolicyWorkflow}
        action={actions.policy.associate.associatePolicyWorkflow}
        targets={[policyAssociation]}
        params={policyFilter}
        postAction={() => {
          //清空查询条件，重新拉取
          // actions.policy.associate.policyList.performSearch('');
          onPostCancel && onPostCancel();
        }}
        width={700}
      >
        <PolicyAssociatePanel />
      </WorkflowDialog>
    );
  }
}
