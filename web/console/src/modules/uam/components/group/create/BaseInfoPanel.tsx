import * as React from 'react';
import { connect } from 'react-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators, insertCSS, OperationState, isSuccessWorkflow } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { RootProps } from '../GroupApp';
import { Button, Tabs, TabPanel, Card } from '@tea/component';
import { UserAssociatePanel } from '../associate/UserAssociatePanel';
import { Group } from '../../../models/Group';
import { router } from '../../../router';
import { UserPlain } from '../../../models';
import { FormPanel } from '@tencent/ff-component';
import { InputField, TipInfo, getWorkflowError } from '../../../../../modules/common';
import { isValid } from '@tencent/ff-validator';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class BaseInfoPanel extends React.Component<RootProps, {}> {

  render() {
    let { actions, route, groupCreation, groupValidator } = this.props;
    let action = actions.group.create.addGroupWorkflow;
    const { groupAddWorkflow } = this.props;
    const workflow = groupAddWorkflow;

    /** 提交 */
    const perform = () => {
      actions.group.create.validator.validate(null, async r => {
        if (isValid(r)) {
          let group: Group = Object.assign({}, groupCreation);
          action.start([group]);
          action.perform();
        }
      });
    };
    /** 取消 */
    const cancel = () => {
      if (workflow.operationState === OperationState.Done) {
        action.reset();
      }
      if (workflow.operationState === OperationState.Started) {
        action.cancel();
      }
      router.navigate({ module: 'group', sub: '' }, route.queries);
    };
    const failed = workflow.operationState === OperationState.Done && !isSuccessWorkflow(workflow);

    return (
      <FormPanel vactions={actions.group.create.validator} formvalidator={groupValidator}>
        <FormPanel.Item
          label={t('用户组名称')}
          vkey="spec.displayName"
          input={{
            placeholder: t('请输入用户组名称，不超过60个字符'),
            value: groupCreation.spec.displayName,
            onChange: value => actions.group.create.updateCreationState({ spec: Object.assign({}, groupCreation.spec, { displayName: value }) })
          }}
        />
        <FormPanel.Item
          label={t('用户组描述')}
          vkey="spec.description"
          input={{
            multiline: true,
            placeholder: t('请输入用户组描述，不超过255个字符'),
            value: groupCreation.spec.description,
            onChange: value => actions.group.create.updateCreationState({ spec: Object.assign({}, groupCreation.spec, { description: value }) })
          }}
        />
        <FormPanel.Item label={t('关联用户')}>
          <UserAssociatePanel onChange={(selection: UserPlain[]) => {
            actions.group.create.updateCreationState({ status: Object.assign({}, groupCreation.status, { users: selection.map((u) => { return { id: u.id } }) }) });
          }}
          />
        </FormPanel.Item>
        <FormPanel.Footer>
          <Button
            className="m"
            type="primary"
            disabled={workflow.operationState === OperationState.Performing}
            onClick={e => { e.preventDefault(); perform() }}>
            {failed ? t('重试') : t('提交')}
          </Button>
          <Button type="weak" onClick={e => { e.preventDefault(); cancel() }}>
            {t('取消')}
          </Button>
          <TipInfo type="error" isForm isShow={failed}>
            {getWorkflowError(workflow)}
          </TipInfo>
        </FormPanel.Footer>
      </FormPanel>
    );
  }
}
