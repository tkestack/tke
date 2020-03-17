import * as React from 'react';
import { connect } from 'react-redux';
import { RootProps } from '../RoleApp';
import { TransferTableProps, TransferTable } from '../../../../common/components';
import { PolicyPlain } from '../../../models';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators } from '@tencent/ff-redux';
import { allActions } from '../../../actions';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

interface Props extends RootProps{
  onChange?: (selection: PolicyPlain[]) => void;
}
@connect(state => state, mapDispatchToProps)
export class PolicyAssociatePanel extends React.Component<Props, {}> {

  render() {
    let { policyAssociation, actions, policyPlainList } = this.props;
    // 表示 ResourceSelector 里要显示和选择的数据类型是 `PolicyPlain`
    const TransferTableSelector = TransferTable as new () => TransferTable<PolicyPlain>;

    // 参数配置
    const selectorProps: TransferTableProps<PolicyPlain> = {
      /** 要供选择的数据 */
      model: policyPlainList,

      /** 用于改变model的query值等 */
      action: actions.policy.associate.policyList,

      /** 已选中的数据 */
      selections: policyAssociation.policies,

      /** 用户选择发生改变后，应该更新选中的数据状态 */
      onChange: (selection: PolicyPlain[]) => {
        actions.policy.associate.selectPolicy(selection);
        this.props.onChange && this.props.onChange(selection);
      },

      /** 选择器标题 */
      title: t(`当前角色可关联以下策略`),

      columns: [
        {
          key: 'name',
          header: t('名称'),
          render: (policy: PolicyPlain) => <p>{`${policy.displayName}`}</p>
        },
        {
          key: 'category',
          header: t('类别'),
          render: (policy: PolicyPlain) => <p>{`${policy.category}`}</p>
        }
      ],
      recordKey: 'id'
    };
    return <TransferTableSelector {...selectorProps} />;
  }
}
