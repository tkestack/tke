import * as React from 'react';
import { connect } from 'react-redux';
import { TablePanel as CTablePanel } from '@tencent/ff-component';
import { LinkButton, emptyTips } from '../../../../common/components';
import { TableColumn, Text, Modal, Icon } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { router } from '../../../router';
import { allActions } from '../../../actions';
import { Group, CommonUserFilter } from '../../../models';
import { RootProps } from '../GroupApp';
import { UserAssociateWorkflowDialog } from '../associate/UserAssociateWorkflowDialog';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class TablePanel extends React.Component<RootProps, {}> {

  render() {
    let { actions, groupList, route } = this.props;

    const columns: TableColumn<Group>[] = [
      {
        key: 'name',
        header: t('用户组名'),
        render: (item, text, index) => (
          <Text parent="div" overflow>
            <a
              href="javascript:;"
              onClick={e => {
                router.navigate({ module: 'group', sub: 'detail' }, { groupName: item.metadata.name });
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
      { key: 'operation', header: t('操作'), render: group => this._renderOperationCell(group) }
    ];

    return (
      <React.Fragment>
        <CTablePanel
          recordKey={(record) => {
            return record.metadata.name;
          }}
          columns={columns}
          model={groupList}
          action={actions.group.list}
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
      </React.Fragment>
    );
  }

  /** 渲染操作按钮 */
  _renderOperationCell = (group: Group) => {
    let { actions } = this.props;
    return (
      <React.Fragment>
        <LinkButton
          tipDirection="right"
          disabled={group.status['phase'] === 'Terminating'}
          onClick={(e) => {
            /** 设置用户关联场景 */
            let filter: CommonUserFilter = {
              resource: 'localgroup',
              resourceID: group.metadata.name,
              /** 关联/解关联回调函数 */
              callback: () => {
                actions.group.list.fetch();
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
        <LinkButton onClick={() => this._removeGroup(group)}><Trans>删除</Trans></LinkButton>
      </React.Fragment>
    );
  }

  _removeGroup = async (group: Group) => {
    let { actions } = this.props;
    const yes = await Modal.confirm({
      message: t('确认删除当前所选用户组') + ` - ${group.spec.displayName}？`,
      okText: t('删除'),
      cancelText: t('取消')
    });
    if (yes) {
      actions.group.list.removeGroupWorkflow.start([group]);
      actions.group.list.removeGroupWorkflow.perform();
    }
  }
}
