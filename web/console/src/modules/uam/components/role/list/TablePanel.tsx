import * as React from 'react';
import { connect } from 'react-redux';
import { TablePanel as CTablePanel } from '@tencent/ff-component';
import { LinkButton, emptyTips } from '../../../../common/components';
import { TableColumn, Text, Modal, Icon } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { router } from '../../../router';
import { allActions } from '../../../actions';
import { Role, CommonUserFilter, PolicyFilter, GroupFilter } from '../../../models';
import { RootProps } from '../RoleApp';
import { UserAssociateWorkflowDialog } from '../associate/UserAssociateWorkflowDialog';
import { GroupAssociateWorkflowDialog } from '../associate/GroupAssociateWorkflowDialog';
import { PolicyAssociateWorkflowDialog } from '../associate/PolicyAssociateWorkflowDialog';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class TablePanel extends React.Component<RootProps, {}> {

  render() {
    let { actions, roleList, route } = this.props;

    const columns: TableColumn<Role>[] = [
      {
        key: 'name',
        header: t('角色名'),
        render: (item, text, index) => (
          <Text parent="div" overflow>
            <a
              href="javascript:;"
              onClick={e => {
                router.navigate({ module: 'role', sub: 'detail' }, { roleName: item.metadata.name });
              }}
            >
              {item.spec.displayName || '-'}
            </a>
            {item.status['phase'] === 'Terminating' && <Icon type="loading" />}
          </Text>
        )
      },
      {
        key: 'description',
        header: t('描述'),
        render: item => <Text parent="div">{item.spec.description || '-'}</Text>
      },
      { key: 'operation', header: t('操作'), render: role => this._renderOperationCell(role) }
    ];

    return (
      <React.Fragment>
        <CTablePanel
          recordKey={(record) => {
            return record.metadata.name;
          }}
          columns={columns}
          model={roleList}
          action={actions.role.list}
          rowDisabled={record => record.status['phase'] === 'Terminating'}
          emptyTips={emptyTips}
          isNeedPagination={true}
          bodyClassName={'tc-15-table-panel tc-15-table-fixed-body'}
        />
        <UserAssociateWorkflowDialog onPostCancel={() => {
          //取消按钮时，清理编辑状态
          actions.commonUser.associate.clearUserAssociation();
        }}
        />
        <GroupAssociateWorkflowDialog onPostCancel={() => {
          //取消按钮时，清理编辑状态
          actions.group.associate.clearGroupAssociation();
        }}
        />
        <PolicyAssociateWorkflowDialog onPostCancel={() => {
          //取消按钮时，清理编辑状态
          actions.policy.associate.clearPolicyAssociation();
        }}
        />
      </React.Fragment>
    );
  }

  /** 渲染操作按钮 */
  _renderOperationCell = (role: Role) => {
    let { actions } = this.props;
    return (
      <React.Fragment>
        <LinkButton
          tipDirection="right"
          disabled={role.status['phase'] === 'Terminating'}
          onClick={(e) => {
            /** 设置策略关联场景 */
            let filter: PolicyFilter = {
              resource: 'role',
              resourceID: role.metadata.name,
              /** 关联/解关联回调函数 */
              callback: () => {
                actions.role.list.fetch();
              }
            };
            actions.policy.associate.setupPolicyFilter(filter);
            /** 拉取关联策略列表，拉取后自动更新policyAssociation */
            actions.policy.associate.policyAssociatedList.applyFilter(filter);
            /** 拉取策略列表 */
            actions.policy.associate.policyList.performSearch('');
            /** 开始关联策略工作流 */
            actions.policy.associate.associatePolicyWorkflow.start();
          }}
        >
          <Trans>关联策略</Trans>
        </LinkButton>
        <LinkButton
          tipDirection="right"
          disabled={role.status['phase'] === 'Terminating'}
          onClick={(e) => {
            /** 设置用户关联场景 */
            let filter: CommonUserFilter = {
              resource: 'role',
              resourceID: role.metadata.name,
              /** 关联/解关联回调函数 */
              callback: () => {
                actions.role.list.fetch();
              }
            };
            actions.commonUser.associate.setupUserFilter(filter);
            /** 拉取关联用户列表，拉取后自动更新commonUserAssociation */
            actions.commonUser.associate.userAssociatedList.applyFilter(filter);
            /** 拉取用户列表 */
            actions.commonUser.associate.userList.performSearch('');
            /** 开始关联用户工作流 */
            actions.commonUser.associate.associateUserWorkflow.start();
          }}
        >
          <Trans>关联用户</Trans>
        </LinkButton>
        <LinkButton
          tipDirection="right"
          disabled={role.status['phase'] === 'Terminating'}
          onClick={(e) => {
            /** 设置用户组关联场景 */
            let filter: GroupFilter = {
              resource: 'role',
              resourceID: role.metadata.name,
              /** 关联/解关联回调函数 */
              callback: () => {
                actions.role.list.fetch();
              }
            };
            actions.group.associate.setupGroupFilter(filter);
            /** 拉取关联用户组列表，拉取后自动更新groupAssociation */
            actions.group.associate.groupAssociatedList.applyFilter(filter);
            /** 拉取用户组列表 */
            actions.group.associate.groupList.performSearch('');
            /** 开始关联用户组工作流 */
            actions.group.associate.associateGroupWorkflow.start();
          }}
        >
          <Trans>关联用户组</Trans>
        </LinkButton>
        <LinkButton onClick={() => this._removeRole(role)}><Trans>删除</Trans></LinkButton>
      </React.Fragment>
    );
  }

  _removeRole = async (role: Role) => {
    let { actions } = this.props;
    const yes = await Modal.confirm({
      message: t('确认删除当前所选角色') + ` - ${role.spec.displayName}？`,
      okText: t('删除'),
      cancelText: t('取消')
    });
    if (yes) {
      actions.role.list.removeRoleWorkflow.start([role]);
      actions.role.list.removeRoleWorkflow.perform();
    }
  }
}
