import * as React from 'react';
import { connect } from 'react-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators, insertCSS, OperationState, isSuccessWorkflow } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { RootProps } from '../RoleApp';
import { Button } from '@tencent/tea-component';
import { router } from '../../../router';
import { FormPanel } from '@tencent/ff-component';
import { InputField, TipInfo, getWorkflowError } from '../../../../../modules/common';
import { UserAssociatePanel } from '../associate/UserAssociatePanel';
import { GroupAssociatePanel } from '../associate/GroupAssociatePanel';
import { PolicyAssociatePanel } from '../associate/PolicyAssociatePanel';
import { UserPlain, PolicyPlain, GroupPlain, Role } from '../../../models';
import { isValid } from '@tencent/ff-validator';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class BaseInfoPanel extends React.Component<RootProps, {}> {

  render() {
    let { actions, route, roleCreation, roleValidator } = this.props;
    let action = actions.role.create.addRoleWorkflow;
    const { roleAddWorkflow } = this.props;
    const workflow = roleAddWorkflow;

    /** 提交 */
    const perform = () => {
      actions.role.create.validator.validate(null, async r => {
        if (isValid(r)) {
          let role: Role = Object.assign({}, roleCreation);
          action.start([role]);
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
      router.navigate({ module: 'role', sub: '' }, route.queries);
    };
    const failed = workflow.operationState === OperationState.Done && !isSuccessWorkflow(workflow);

    return (
      <FormPanel vactions={actions.role.create.validator} formvalidator={roleValidator}>
        <FormPanel.Item
          label={t('角色名称')}
          vkey="spec.displayName"
          input={{
            placeholder: t('请输入角色名称，不超过60个字符'),
            value: roleCreation.spec.displayName,
            onChange: value => actions.role.create.updateCreationState({ spec: Object.assign({}, roleCreation.spec, { displayName: value }) })
          }}
        />
        <FormPanel.Item
          label={t('角色描述')}
          vkey="spec.description"
          input={{
            multiline: true,
            placeholder: t('请输入角色描述，不超过255个字符'),
            value: roleCreation.spec.description,
            onChange: value => actions.role.create.updateCreationState({ spec: Object.assign({}, roleCreation.spec, { description: value }) })
          }}
        />
        <FormPanel.Item label={t('关联策略')}>
          <PolicyAssociatePanel onChange={(selection: PolicyPlain[]) => {
            actions.role.create.updateCreationState({ spec: Object.assign({}, roleCreation.spec, { policies: selection.map((p) => { return p.id }) }) });
          }}
          />
        </FormPanel.Item>
        <FormPanel.Item label={t('关联用户')}>
          <UserAssociatePanel onChange={(selection: UserPlain[]) => {
            actions.role.create.updateCreationState({ status: Object.assign({}, roleCreation.status, { users: selection.map((u) => { return { id: u.id } }) }) });
          }}
          />
        </FormPanel.Item>
        <FormPanel.Item label={t('关联用户组')}>
          <GroupAssociatePanel onChange={(selection: GroupPlain[]) => {
            actions.role.create.updateCreationState({ status: Object.assign({}, roleCreation.status, { groups: selection.map((g) => { return { id: g.id } }) }) });
          }}
          />
        </FormPanel.Item>
        <FormPanel.Footer>
          <React.Fragment>
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
          </React.Fragment>
        </FormPanel.Footer>
      </FormPanel>
    );
  }
}
