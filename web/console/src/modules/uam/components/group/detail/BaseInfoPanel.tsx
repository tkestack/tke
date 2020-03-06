import * as React from 'react';
import { connect } from 'react-redux';
import { RootProps } from '../GroupApp';
import { FormPanel } from '@tencent/ff-component';
import { TipInfo, getWorkflowError, InputField } from '../../../../../modules/common';
import { Button, Tabs, TabPanel, Card } from '@tea/component';
import { dateFormat } from '../../../../../../helpers/dateUtil';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators, OperationState, isSuccessWorkflow } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { UserActionPanel } from './UserActionPanel';
import { UserTablePanel } from './UserTablePanel';
import { RoleActionPanel } from './RoleActionPanel';
import { RoleTablePanel } from './RoleTablePanel';
import { Group } from '../../../models/Group';
import { PolicyActionPanel } from './PolicyActionPanel';
import { PolicyTablePanel } from './PolicyTablePanel';
import { isValid } from '@tencent/ff-validator';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class BaseInfoPanel extends React.Component<RootProps> {

  render() {
    let { actions, groupEditor, route, groupValidator } = this.props;

    let action = actions.group.detail.updateGroupWorkflow;
    const { groupUpdateWorkflow } = this.props;
    const workflow = groupUpdateWorkflow;

    /** 提交 */
    const perform = () => {
      actions.group.detail.validator.validate(null, async r => {
        if (isValid(r)) {
          let group: Group = Object.assign({}, groupEditor);
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
      actions.group.detail.updateEditorState({ v_editing: false });
    };
    const failed = workflow.operationState === OperationState.Done && !isSuccessWorkflow(workflow);

    const tabs = [
      { id: 'users', label: '关联用户' },
      { id: 'roles', label: '已关联角色' },
      { id: 'policies', label: '已关联策略' },
    ];

    return (
      <React.Fragment>
        <Card>
          <Card.Body
            title={t('基本信息')}
            subtitle={
              <React.Fragment>
                <Button type="link" onClick={e =>
                  actions.group.detail.updateEditorState({ v_editing: true })
                }>
                  {t('编辑')}
                </Button>
              </React.Fragment>
              }
            >
            <FormPanel isNeedCard={false} vactions={actions.group.detail.validator} formvalidator={groupValidator}>
              <FormPanel.Item text label={t('用户组ID')}>
                {groupEditor.metadata.name}
              </FormPanel.Item>
              {!groupEditor.v_editing ?
              (<FormPanel.Item text label={t('用户组名称')}>
                {groupEditor.spec.displayName}
              </FormPanel.Item>) :
              (<FormPanel.Item
                label={t('用户组名称')}
                vkey="spec.displayName"
                input={{
                  placeholder: t('请输入用户组名称，不超过60个字符'),
                  value: groupEditor.spec.displayName,
                  onChange: value => actions.group.detail.updateEditorState({ spec: Object.assign({}, groupEditor.spec, { displayName: value }) })
                }}
              />)
              }
              {!groupEditor.v_editing ?
              (<FormPanel.Item text label={t('用户组描述')}>
                {groupEditor.spec.description}
              </FormPanel.Item>) :
              (<FormPanel.Item
                label={t('用户组描述')}
                vkey="spec.description"
                input={{
                  placeholder: t('请输入用户组描述，不超过255个字符'),
                  value: groupEditor.spec.description,
                  onChange: value => actions.group.detail.updateEditorState({ spec: Object.assign({}, groupEditor.spec, { description: value }) })
                }}
              />)
              }
              <FormPanel.Item text label={t('创建时间')}>
                {dateFormat(new Date(groupEditor.metadata.creationTimestamp), 'yyyy-MM-dd hh:mm:ss')}
              </FormPanel.Item>
              {groupEditor.v_editing && (
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
            <Tabs tabs={tabs}>
              <TabPanel id="users">
                <UserActionPanel />
                <UserTablePanel />
              </TabPanel>
              <TabPanel id="roles">
                <RoleActionPanel />
                <RoleTablePanel />
              </TabPanel>
              <TabPanel id="policies">
                <PolicyActionPanel />
                <PolicyTablePanel />
              </TabPanel>
            </Tabs>
          </Card.Body>
        </Card>
      </React.Fragment>
    );
  }
}
