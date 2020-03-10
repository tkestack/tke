import * as React from 'react';
import { connect } from 'react-redux';
import { TablePanel } from '@tencent/ff-component';
import { LinkButton, emptyTips } from '../../../../common/components';
import { Table, TableColumn, Text, Modal, Icon } from '@tea/component';
import { bindActionCreators, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { allActions } from '../../../actions';
import { UserPlain, CommonUserAssociation } from '../../../models';
import { RootProps } from '../RoleApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class UserTablePanel extends React.Component<RootProps, {}> {

  render() {
    let { actions, commonUserAssociation, commonUserAssociatedList } = this.props;

    const columns: TableColumn<UserPlain>[] = [
      {
        key: 'name',
        header: t('用户ID / 名称'),
        render: (user, text, index) => (
          <Text parent="div" overflow>
            {user.name || '-'}{' / '}{user.displayName || '-'}
          </Text>
        )
      },
      { key: 'operation', header: t('操作'), render: user => this._renderOperationCell(user) }
    ];

    return (
      <TablePanel
        columns={columns}
        recordKey={'name'}
        records={commonUserAssociation.originUsers}
        action={actions.commonUser.associate.userAssociatedList}
        model={commonUserAssociatedList}
        emptyTips={emptyTips}
      />
    );
  }

  /** 渲染操作按钮 */
  _renderOperationCell = (user: UserPlain) => {
    let { actions } = this.props;
    return (
      <React.Fragment>
        <LinkButton
          tipDirection="right"
          onClick={(e) => {
            this._removeUser(user);
          }}
        >
          <Trans>解除关联</Trans>
        </LinkButton>
      </React.Fragment>
    );
  }

  _removeUser = async (user: UserPlain) => {
    let { actions, commonUserFilter } = this.props;
    const yes = await Modal.confirm({
      message: t('确认解除当前用户关联') + ` - ${user.displayName}？`,
      okText: t('解除'),
      cancelText: t('取消')
    });
    if (yes) {
      let userAssociation: CommonUserAssociation = { id: uuid(), removeUsers: [user] };
      actions.commonUser.associate.disassociateUserWorkflow.start([userAssociation], commonUserFilter);
      actions.commonUser.associate.disassociateUserWorkflow.perform();
    }
  }

}
