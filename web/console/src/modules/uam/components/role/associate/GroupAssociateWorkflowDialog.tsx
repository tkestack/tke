import * as React from 'react';
import { connect } from 'react-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators, insertCSS } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { RootProps } from '../RoleApp';
import { WorkflowDialog } from '../../../../common/components';
import { GroupAssociatePanel } from './GroupAssociatePanel';
const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

interface WorkflowDialogProps extends RootProps {
  onPostCancel?: () => void;
}

@connect(state => state, mapDispatchToProps)
export class GroupAssociateWorkflowDialog extends React.Component<WorkflowDialogProps, {}> {

  render() {
    const {
      actions,
      groupAssociation,
      groupFilter,
      associateGroupWorkflow
    } = this.props;
    const { onPostCancel = undefined } = this.props;
    return (
      <WorkflowDialog
        caption={t('关联用户组')}
        workflow={associateGroupWorkflow}
        action={actions.group.associate.associateGroupWorkflow}
        targets={[groupAssociation]}
        params={groupFilter}
        postAction={() => {
          //清空查询条件，重新拉取
          // actions.group.associate.groupList.performSearch('');
          onPostCancel && onPostCancel();
        }}
        width={700}
      >
        <GroupAssociatePanel />
      </WorkflowDialog>
    );
  }
}
