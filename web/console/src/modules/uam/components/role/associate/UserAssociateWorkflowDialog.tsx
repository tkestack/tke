import * as React from 'react';
import { connect } from 'react-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators, insertCSS } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { RootProps } from '../RoleApp';
import { WorkflowDialog } from '../../../../common/components';
import { UserAssociatePanel } from './UserAssociatePanel';
const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

interface WorkflowDialogProps extends RootProps {
  onPostCancel?: () => void;
}

@connect(state => state, mapDispatchToProps)
export class UserAssociateWorkflowDialog extends React.Component<WorkflowDialogProps, {}> {

  render() {
    const {
      actions,
      commonUserAssociation,
      commonUserFilter,
      commonAssociateUserWorkflow
    } = this.props;
    const { onPostCancel = undefined } = this.props;
    return (
      <WorkflowDialog
        caption={t('关联用户')}
        workflow={commonAssociateUserWorkflow}
        action={actions.commonUser.associate.associateUserWorkflow}
        targets={[commonUserAssociation]}
        params={commonUserFilter}
        postAction={() => {
          //清空查询条件，重新拉取
          // actions.commonUser.associate.userList.performSearch('');
          onPostCancel && onPostCancel();
        }}
        width={700}
      >
        <UserAssociatePanel />
      </WorkflowDialog>
    );
  }
}
