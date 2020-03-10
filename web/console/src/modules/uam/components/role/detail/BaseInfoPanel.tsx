import * as React from 'react';
import { connect } from 'react-redux';
import { RootProps } from '../RoleApp';
import { FormPanel } from '@tencent/ff-component';
import { TipInfo, getWorkflowError, InputField } from '../../../../../modules/common';
import { Button, Tabs, TabPanel, Card } from '@tea/component';
import { dateFormat } from '../../../../../../helpers/dateUtil';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators, OperationState, isSuccessWorkflow  } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { PolicyActionPanel } from './PolicyActionPanel';
import { PolicyTablePanel } from './PolicyTablePanel';
import { UserActionPanel } from './UserActionPanel';
import { UserTablePanel } from './UserTablePanel';
import { GroupActionPanel } from './GroupActionPanel';
import { GroupTablePanel } from './GroupTablePanel';
import { isValid } from '@tencent/ff-validator';
import { Role } from '../../../models';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class BaseInfoPanel extends React.Component<RootProps> {

  render() {
    let { actions, roleEditor, route, roleValidator } = this.props;

    let action = actions.role.detail.updateRoleWorkflow;
    const { roleUpdateWorkflow } = this.props;
    const workflow = roleUpdateWorkflow;

    /** 提交 */
    const perform = () => {
      actions.role.detail.validator.validate(null, async r => {
        if (isValid(r)) {
          let role: Role = Object.assign({}, roleEditor);
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
      actions.role.detail.updateEditorState({ v_editing: false });
    };
    const failed = workflow.operationState === OperationState.Done && !isSuccessWorkflow(workflow);

    const tabs = [
      { id: 'policies', label: '关联策略' },
      { id: 'users', label: '关联用户' },
      { id: 'groups', label: '关联用户组' }
    ];

    return (
      <React.Fragment>
        <Card>
          <Card.Body
            title={t('基本信息')}
            subtitle={
              <React.Fragment>
                <Button type="link" onClick={e =>
                  actions.role.detail.updateEditorState({ v_editing: true })
                }>
                  {t('编辑')}
                </Button>
              </React.Fragment>
              }
            >
            <FormPanel isNeedCard={false} vactions={actions.role.detail.validator} formvalidator={roleValidator}>
              <FormPanel.Item text label={t('角色ID')}>
                {roleEditor.metadata.name}
              </FormPanel.Item>
              {!roleEditor.v_editing ?
                <FormPanel.Item text label={t('角色名称')}>
                  {roleEditor.spec.displayName}
                </FormPanel.Item> :
                <FormPanel.Item
                  label={t('角色名称')}
                  vkey="spec.displayName"
                  input={{
                    placeholder: t('请输入角色名称，不超过60个字符'),
                    value: roleEditor.spec.displayName,
                    onChange: value => actions.role.detail.updateEditorState({ spec: Object.assign({}, roleEditor.spec, { displayName: value }) })
                  }}
                />
              }
              {!roleEditor.v_editing ?
                <FormPanel.Item text label={t('角色描述')}>
                  {roleEditor.spec.description}
                </FormPanel.Item> :
                <FormPanel.Item
                  label={t('角色描述')}
                  vkey="spec.description"
                  input={{
                    multiline: true,
                    placeholder: t('请输入角色描述，不超过255个字符'),
                    value: roleEditor.spec.description,
                    onChange: value => actions.role.detail.updateEditorState({ spec: Object.assign({}, roleEditor.spec, { description: value }) })
                  }}
                />
              }
              <FormPanel.Item text label={t('创建时间')}>
                {dateFormat(new Date(roleEditor.metadata.creationTimestamp), 'yyyy-MM-dd hh:mm:ss')}
              </FormPanel.Item>
              {roleEditor.v_editing && (
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
            )}
            </FormPanel>
          </Card.Body>
        </Card>
        <Card>
          <Card.Body>
            <Tabs
              tabs={tabs}
              defaultActiveId={'policies'}
              onActive={tab => {
              }}
            >
              <TabPanel id="policies">
                <PolicyActionPanel />
                <PolicyTablePanel />
              </TabPanel>
              <TabPanel id="users">
                <UserActionPanel />
                <UserTablePanel />
              </TabPanel>
              <TabPanel id="groups">
                <GroupActionPanel />
                <GroupTablePanel />
              </TabPanel>
            </Tabs>
          </Card.Body>
        </Card>
      </React.Fragment>
    );
  }
}
