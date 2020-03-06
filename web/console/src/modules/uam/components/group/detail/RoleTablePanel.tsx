import * as React from 'react';
import { connect } from 'react-redux';
import { LinkButton, emptyTips } from '../../../../common/components';
import { TablePanel } from '@tencent/ff-component';
import { Table, TableColumn, Text, Modal, Icon } from '@tea/component';
import { bindActionCreators, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { router } from '../../../router';
import { allActions } from '../../../actions';
import { RolePlain, RoleAssociation, GroupAssociation, GroupPlain, GroupFilter } from '../../../models';
import { RootProps } from '../GroupApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class RoleTablePanel extends React.Component<RootProps, {}> {

  render() {
    let { actions, roleAssociation, roleAssociatedList } = this.props;

    const columns: TableColumn<RolePlain>[] = [
      {
        key: 'name',
        header: t('名称'),
        render: (role, text, index) => (
          <Text parent="div" overflow>
            {role.displayName || '-'}
          </Text>
        )
      },
      {
        key: 'description',
        header: t('描述'),
        render: (role, text, index) => (
          <Text parent="div" overflow>
            {role.description || '-'}
          </Text>
        )
      },
      // { key: 'operation', header: t('操作'), render: role => this._renderOperationCell(role) }
    ];

    return (
      <React.Fragment>
        <TablePanel
          columns={columns}
          recordKey={'id'}
          records={roleAssociation.originRoles}
          action={actions.role.associate.roleAssociatedList}
          model={roleAssociatedList}
          emptyTips={emptyTips}
        />
      </React.Fragment>
    );
  }

  /** 渲染操作按钮 */
  _renderOperationCell = (role: RolePlain) => {
    let { actions } = this.props;
    return (
      <React.Fragment>
        <LinkButton
          tipDirection="right"
          onClick={(e) => {
            this._removeRole(role);
          }}
        >
          <Trans>解除关联</Trans>
        </LinkButton>
      </React.Fragment>
    );
  }

  _removeRole = async (role: RolePlain) => {
    let { actions, roleFilter } = this.props;
    const yes = await Modal.confirm({
      message: t('确认解除当前角色关联') + ` - ${role.displayName}？`,
      okText: t('解除'),
      cancelText: t('取消')
    });
    if (yes) {
      /** 目前还没有实现基于用户组解绑角色 */
      let roleAssociation: RoleAssociation = { id: uuid(), removeRoles: [role] };
      actions.role.associate.disassociateRoleWorkflow.start([roleAssociation], roleFilter);
      actions.role.associate.disassociateRoleWorkflow.perform();
    }
  }

}
