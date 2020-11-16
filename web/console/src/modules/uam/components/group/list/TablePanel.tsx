import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { TablePanel as CTablePanel } from '@tencent/ff-component';
import { LinkButton, emptyTips } from '../../../../common/components';
import { TableColumn, Text, Modal, Icon } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { router } from '../../../router';
import { allActions } from '../../../actions';
import { Group, CommonUserFilter } from '../../../models';
import { RootProps } from '../GroupPanel';
import { UserAssociateWorkflowDialog } from '../associate/UserAssociateWorkflowDialog';
import { useModal } from '@src/modules/common/utils/tHooks';
import { RoleModifyDialog } from './RoleModifyDialog';
const { useState, useEffect } = React;
// const mapDispatchToProps = (dispatch) =>
//   Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
//     dispatch,
//   });
//
// @connect((state) => state, mapDispatchToProps)
export const TablePanel = props => {
  const state = useSelector(state => state);
  const dispatch = useDispatch();
  const { actions } = bindActionCreators({ actions: allActions }, dispatch);
  const { groupList, route } = state;

  const { isShowing, toggle } = useModal(false);
  const [editUserGroup, setEditUserGroup] = useState<Group | undefined>();

  useEffect(() => {
    actions.policy.associate.policyList.applyFilter({ resource: 'platform', resourceID: '' });
  }, []);

  const columns: TableColumn<Group>[] = [
    {
      key: 'name',
      header: t('用户组名'),
      render: (item, text, index) => (
        <Text parent="div" overflow>
          <a
            href="javascript:;"
            onClick={e => {
              router.navigate({ module: 'user', sub: 'group', action: 'detail' }, { groupName: item.metadata.name });
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
    {
      key: 'policies',
      header: t('角色'),
      render: item => {
        const content = Object.values(JSON.parse(item.spec.extra.policies)).join(',');
        return (
          <Text>
            {content || '-'}
            <Icon
              onClick={() => {
                toggle();
                setEditUserGroup({ ...item });
              }}
              style={{ cursor: 'pointer' }}
              type="pencil"
            />
          </Text>
        );
      }
    },
    { key: 'operation', header: t('操作'), render: group => _renderOperationCell(group) }
  ];

  return (
    <React.Fragment>
      <CTablePanel
        recordKey={record => {
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
      <RoleModifyDialog isShowing={isShowing} toggle={toggle} user={editUserGroup} />
      <UserAssociateWorkflowDialog
        onPostCancel={() => {
          //取消按钮时，清理编辑状态
          actions.commonUser.associate.clearUserAssociation();
        }}
      />
    </React.Fragment>
  );

  /** 渲染操作按钮 */
  function _renderOperationCell(group: Group) {
    // let { actions } = this.props;
    return (
      <React.Fragment>
        <LinkButton
          tipDirection="right"
          disabled={group.status['phase'] === 'Terminating'}
          onClick={e => {
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
        <LinkButton onClick={() => _removeGroup(group)}>
          <Trans>删除</Trans>
        </LinkButton>
      </React.Fragment>
    );
  }

  async function _removeGroup(group: Group) {
    // let { actions } = this.props;
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
};
