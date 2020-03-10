import * as React from 'react';
import { connect } from 'react-redux';
import { TablePanel } from '@tencent/ff-component';
import { LinkButton, emptyTips } from '../../../../common/components';
import { Table, TableColumn, Text, Modal, Icon } from '@tea/component';
import { bindActionCreators, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { router } from '../../../router';
import { allActions } from '../../../actions';
import { GroupAssociation, GroupPlain, GroupFilter } from '../../../models';
import { RootProps } from '../UserApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class GroupTablePanel extends React.Component<RootProps, {}> {

  render() {
    let { actions, groupAssociation, groupAssociatedList } = this.props;

    const columns: TableColumn<GroupPlain>[] = [
      {
        key: 'name',
        header: t('名称'),
        render: (group, text, index) => (
          <Text parent="div" overflow>
            {group.displayName || '-'}
          </Text>
        )
      },
      {
        key: 'description',
        header: t('描述'),
        render: (group, text, index) => (
          <Text parent="div" overflow>
            {group.description || '-'}
          </Text>
        )
      },
      // { key: 'operation', header: t('操作'), render: group => this._renderOperationCell(group) }
    ];

    return (
      <React.Fragment>
        <TablePanel
          columns={columns}
          recordKey={'id'}
          records={groupAssociation.originGroups}
          action={actions.group.associate.groupAssociatedList}
          model={groupAssociatedList}
          emptyTips={emptyTips}
        />
      </React.Fragment>
    );
  }

  /** 渲染操作按钮 */
  _renderOperationCell = (group: GroupPlain) => {
    let { actions } = this.props;
    return (
      <React.Fragment>
        <LinkButton
          tipDirection="right"
          onClick={(e) => {
            this._removeGroup(group);
          }}
        >
          <Trans>解除关联</Trans>
        </LinkButton>
      </React.Fragment>
    );
  }

  _removeGroup = async (group: GroupPlain) => {
    let { actions, groupFilter } = this.props;
    const yes = await Modal.confirm({
      message: t('确认解除当前用户组关联') + ` - ${group.displayName}？`,
      okText: t('解除'),
      cancelText: t('取消')
    });
    if (yes) {
      /** 目前还没有实现基于用户解绑用户组 */
      let groupAssociation: GroupAssociation = { id: uuid(), removeGroups: [group] };
      actions.group.associate.disassociateGroupWorkflow.start([groupAssociation], groupFilter);
      actions.group.associate.disassociateGroupWorkflow.perform();
    }
  }

}
