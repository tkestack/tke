import * as React from 'react';
import { connect } from 'react-redux';
import { RootProps } from '../RoleApp';
import { TransferTableProps, TransferTable } from '../../../../common/components';
import { UserPlain } from '../../../models';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators } from '@tencent/ff-redux';
import { allActions } from '../../../actions';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

interface Props extends RootProps{
  onChange?: (selection: UserPlain[]) => void;
}
@connect(state => state, mapDispatchToProps)
export class UserAssociatePanel extends React.Component<Props, {}> {

  render() {
    let { commonUserAssociation, actions, userPlainList } = this.props;
    // 表示 ResourceSelector 里要显示和选择的数据类型是 `UserPlain`
    const TransferTableSelector = TransferTable as new () => TransferTable<UserPlain>;

    // 参数配置
    const selectorProps: TransferTableProps<UserPlain> = {
      /** 要供选择的数据 */
      model: userPlainList,

      /** 用于改变model的query值等 */
      action: actions.commonUser.associate.userList,

      /** 已选中的数据 */
      selections: commonUserAssociation.users,

      /** 用户选择发生改变后，应该更新选中的数据状态 */
      onChange: (selection: UserPlain[]) => {
        actions.commonUser.associate.selectUser(selection);
        this.props.onChange && this.props.onChange(selection);
      },

      /** 选择器标题 */
      title: t(`当前角色可关联以下用户`),

      columns: [
        {
          key: 'name',
          header: t('ID/名称'),
          render: (user: UserPlain) => <p>{`${user.displayName}(${user.name})`}</p>
        }
      ],
      recordKey: 'id'
    };
    return <TransferTableSelector {...selectorProps} />;
  }
}
