import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';

import { Button, Form, Icon, Input, Modal, TableColumn, Text } from '@tea/component';
import { TablePanel } from '@tencent/ff-component';
import { bindActionCreators, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { useModal } from '@src/modules/common/utils/tHooks';
import { LinkButton } from '../../../../common/components';
import { allActions } from '../../../actions';
import { User } from '../../../models';
import { router } from '../../../router';
import { RoleModifyDialog } from './RoleModifyDialog';
import { PlatformTypeEnum } from '@src/modules/project/constants/Config';
const { useState, useEffect } = React;

export const UserTablePanel = () => {
  const state = useSelector(state => state);
  const dispatch = useDispatch();
  const { actions } = bindActionCreators({ actions: allActions }, dispatch);
  const { isShowing, toggle } = useModal(false);
  const [editUser, setEditUser] = useState<User | undefined>();
  const { userList, route, platformType, userManagedProjects, projectDetail } = state;

  let enableOp =
    platformType === PlatformTypeEnum.Manager ||
    (platformType === PlatformTypeEnum.Business &&
      userManagedProjects.list.data.records.find(
        item => item.name === (projectDetail ? projectDetail.metadata.name : null)
      ));

  useEffect(() => {
    actions.policy.associate.policyList.applyFilter({ resource: 'project', resourceID: '' });
  }, []);

  const columns: TableColumn<User>[] = [
    {
      key: 'name',
      header: t('用户ID / 名称'),
      render: (user, text, index) => (
        <Text parent="div" overflow>
          {user.spec.name} / {user.spec.displayName || '-'}
          {!!user.status && user.status.phase === 'Deleting' && (
            <React.Fragment>
              <Icon type="loading" />
            </React.Fragment>
          )}
        </Text>
      )
    },
    {
      key: 'phone',
      header: t('关联手机'),
      render: user => <Text>{user.spec.phoneNumber || '-'}</Text>
    },
    {
      key: 'email',
      header: t('关联邮箱'),
      render: user => <Text>{user.spec.email || '-'}</Text>
    },
    {
      key: 'policies',
      header: t('角色'),
      render: user => {
        const content = Object.values(JSON.parse(user.spec.extra.policies)).join(',');
        return (
          <Text>
            {content || '-'}
            {enableOp && (
              <Icon
                onClick={() => {
                  toggle();
                  setEditUser({ ...user });
                }}
                style={{ cursor: 'pointer' }}
                type="pencil"
              />
            )}
          </Text>
        );
      }
    },
    { key: 'operation', header: t('操作'), render: user => _renderOperationCell(user) }
  ];
  const emptyTips: JSX.Element = (
    <div className="text-center">
      <Trans>暂无内容</Trans>
    </div>
  );

  /** 渲染操作按钮 */
  function _renderOperationCell(user: User) {
    const isDisable = !!user.status && user.status.phase === 'Deleting';

    if (user.spec.name.toLowerCase() === 'admin') {
      return (
        <LinkButton tipDirection="right" errorTip="管理员不能被删除" disabled>
          <Trans>删除</Trans>
        </LinkButton>
      );
    }
    return enableOp ? (
      <React.Fragment>
        <LinkButton
          disabled={isDisable}
          tipDirection="right"
          onClick={() => {
            _removeUser(user);
          }}
        >
          <Trans>删除</Trans>
        </LinkButton>
      </React.Fragment>
    ) : null;
  }

  async function _removeUser(user: User) {
    const yes = await Modal.confirm({
      message: t('确认删除当前所选用户？'),
      description: t('删除后，用户{{username}}的所有配置将会被清空，且无法恢复', { username: user.spec.name }),
      okText: t('删除'),
      cancelText: t('取消')
    });
    if (yes) {
      let userInfo = {
        id: uuid(),
        projectId: route.queries.projectId,
        users: [{ id: user.metadata.name }],
        policies: []
      };
      actions.user.addUser.start([userInfo]);
      actions.user.addUser.perform();
    }
  }

  return (
    <>
      <TablePanel
        recordKey={record => {
          return record.metadata.name;
        }}
        columns={columns}
        rowDisabled={record => record.status && record.status.phase === 'Deleting'}
        model={userList}
        action={actions.user}
        emptyTips={emptyTips}
        bodyClassName={'tc-15-table-panel tc-15-table-fixed-body'}
      />
      <RoleModifyDialog isShowing={isShowing} toggle={toggle} user={editUser} />
    </>
  );
};
