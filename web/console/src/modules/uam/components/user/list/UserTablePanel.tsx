import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';

import { Button, Form, Icon, Input, Modal, TableColumn, Text } from '@tea/component';
import { TablePanel } from '@tencent/ff-component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { useModal } from '@src/modules/common/utils/tHooks';
import { LinkButton } from '../../../../common/components';
import { allActions } from '../../../actions';
import { VALIDATE_PASSWORD_RULE } from '../../../constants/Config';
import { User } from '../../../models';
import { router } from '../../../router';
import { RoleModifyDialog } from '@src/modules/uam/components/user/list/RoleModifyDialog';
import { PasswordModifyDialog } from '@src/modules/uam/components/user/list/PasswordModifyDialog';

const { useState, useEffect } = React;

export const UserTablePanel = () => {
  const state = useSelector((state) => state);
  const dispatch = useDispatch();
  const { actions } = bindActionCreators({ actions: allActions }, dispatch);

  const [pwdModalVisible, setPwdModalVisible] = useState(false);
  const [pwdMessages, setPwdMessages] = useState({ passwordMsg: '', rePasswordMsg: '' });
  const [passwords, setPasswords] = useState({ password: '', rePassword: '' });
  const [btnDisabled, setBtnDisabled] = useState(true);
  // const [user, setUser] = useState(undefined);
  const { isShowing, toggle } = useModal(false);
  const { isShowing: pwdIsShowing, toggle: pwdToggle } = useModal(false);
  const [editUser, setEditUser] = useState();
  const { userList } = state;

  useEffect(() => {
    // 这种调用数据可以缓存，不会和下边的performSearch等互相影响
    actions.policy.associate.policyList.applyFilter({ resource: 'platform', resourceID: '' });
    // actions.policy.associate.policyList.performSearch('platform');
    // actions.strategy.getPlatformCategories.fetch();
    // actions.policy.associate.policyList.fetch({
    //   noCache: true,
    //   data: { flag: 'platform' },
    // });
    // actions.policy.associate.policyList.applyFilter({ queryFilter: 'platform' });
  }, []);

  const columns: TableColumn<User>[] = [
    {
      key: 'name',
      header: t('用户ID / 名称'),
      render: (user, text, index) => (
        <Text parent="div" overflow>
          <a
            href="javascript:;"
            onClick={(e) => {
              router.navigate(
                { module: 'user', sub: 'normal', action: 'detail' },
                { username: user.spec.username, name: user.metadata.name }
              );
            }}
          >
            {user.spec.username} / {user.spec.displayName || '-'}
          </a>
          {user.status.phase === 'Deleting' && (
            <React.Fragment>
              <Icon type="loading" />
            </React.Fragment>
          )}
        </Text>
      ),
    },
    {
      key: 'phone',
      header: t('关联手机'),
      render: (user) => <Text>{user.spec.phoneNumber || '-'}</Text>,
    },
    {
      key: 'email',
      header: t('关联邮箱'),
      render: (user) => <Text>{user.spec.email || '-'}</Text>,
    },
    {
      key: 'policies',
      header: t('角色'),
      render: (user) => {
        const content = Object.values(JSON.parse(user.spec.extra.policies)).join(',');
        return (
          <Text>
            {content || '-'}
            <Icon
              onClick={() => {
                // actions.cluster.selectCluster([x]);
                // actions.workflow.modifyClusterName.start([x], 1);
                toggle();
                setEditUser(user);
              }}
              style={{ cursor: 'pointer' }}
              type="pencil"
            />
          </Text>
        );
      },
    },
    { key: 'operation', header: t('操作'), render: (user) => _renderOperationCell(user) },
  ];
  const emptyTips: JSX.Element = (
    <div className="text-center">
      <Trans>暂无内容</Trans>
    </div>
  );

  return (
    <>
      <TablePanel
        recordKey={(record) => {
          return record.metadata.name;
        }}
        columns={columns}
        rowDisabled={(record) => record.status && record.status.phase === 'Deleting'}
        model={userList}
        action={actions.user}
        emptyTips={emptyTips}
        bodyClassName={'tc-15-table-panel tc-15-table-fixed-body'}
      />
      <RoleModifyDialog isShowing={isShowing} toggle={toggle} user={editUser} />
      <PasswordModifyDialog isShowing={pwdIsShowing} toggle={pwdToggle} user={editUser} />
      {/*<Modal*/}
      {/*  visible={pwdModalVisible}*/}
      {/*  size="s"*/}
      {/*  caption={t('修改密码')}*/}
      {/*  onClose={() => {*/}
      {/*    setPwdModalVisible(false);*/}
      {/*  }}*/}
      {/*>*/}
      {/*  <Modal.Body>*/}
      {/*    <Form>*/}
      {/*      <Form.Item*/}
      {/*        label={t('用户密码')}*/}
      {/*        required*/}
      {/*        status={pwdMessages.passwordMsg ? 'error' : passwords.password ? 'success' : undefined}*/}
      {/*        message={pwdMessages.passwordMsg ? t(pwdMessages.passwordMsg) : t(VALIDATE_PASSWORD_RULE.message)}*/}
      {/*      >*/}
      {/*        <Input*/}
      {/*          type="password"*/}
      {/*          placeholder={t('请输入用户密码')}*/}
      {/*          defaultValue={passwords.password}*/}
      {/*          onChange={(value) => {*/}
      {/*            let pwdMsg = '';*/}
      {/*            let disabled = true;*/}
      {/*            if (!value) {*/}
      {/*              pwdMsg = '请输入用户密码';*/}
      {/*            } else if (!VALIDATE_PASSWORD_RULE.pattern.test(value)) {*/}
      {/*              pwdMsg = VALIDATE_PASSWORD_RULE.message;*/}
      {/*            } else if (passwords.rePassword && value !== passwords.rePassword) {*/}
      {/*              pwdMsg = '两次输入密码不一致';*/}
      {/*            }*/}
      {/*            if (value && passwords.rePassword && !pwdMsg && !pwdMessages.rePasswordMsg) {*/}
      {/*              disabled = false;*/}
      {/*            }*/}
      {/*            setPasswords({ ...passwords, password: value });*/}
      {/*            setPwdMessages({ ...pwdMessages, passwordMsg: pwdMsg });*/}
      {/*            setBtnDisabled(disabled);*/}
      {/*          }}*/}
      {/*        />*/}
      {/*      </Form.Item>*/}
      {/*      <Form.Item*/}
      {/*        label={t('确认密码')}*/}
      {/*        required*/}
      {/*        status={pwdMessages.rePasswordMsg ? 'error' : passwords.rePassword ? 'success' : undefined}*/}
      {/*        message={pwdMessages.rePasswordMsg ? t(pwdMessages.rePasswordMsg) : t(VALIDATE_PASSWORD_RULE.message)}*/}
      {/*      >*/}
      {/*        <Input*/}
      {/*          type="password"*/}
      {/*          placeholder={t('请再次输入用户密码')}*/}
      {/*          defaultValue={passwords.rePassword}*/}
      {/*          onChange={(value) => {*/}
      {/*            let rePwdMsg = '';*/}
      {/*            let disabled = true;*/}
      {/*            if (!value) {*/}
      {/*              rePwdMsg = '请输入用户密码';*/}
      {/*            } else if (!VALIDATE_PASSWORD_RULE.pattern.test(value)) {*/}
      {/*              rePwdMsg = VALIDATE_PASSWORD_RULE.message;*/}
      {/*            } else if (passwords.password && value !== passwords.password) {*/}
      {/*              rePwdMsg = '两次输入密码不一致';*/}
      {/*            }*/}
      {/*            if (passwords.password && value && !pwdMessages.passwordMsg && !rePwdMsg) {*/}
      {/*              disabled = false;*/}
      {/*            }*/}
      {/*            setPasswords({ ...passwords, rePassword: value });*/}
      {/*            setPwdMessages({ ...pwdMessages, rePasswordMsg: rePwdMsg });*/}
      {/*            setBtnDisabled(disabled);*/}
      {/*          }}*/}
      {/*        />*/}
      {/*      </Form.Item>*/}
      {/*    </Form>*/}
      {/*  </Modal.Body>*/}
      {/*  <Modal.Footer>*/}
      {/*    <Form.Action>*/}
      {/*      <Button disabled={btnDisabled} type="primary" onClick={_onPwdModalSubmit}>*/}
      {/*        <Trans>保存</Trans>*/}
      {/*      </Button>*/}
      {/*      <Button*/}
      {/*        onClick={() => {*/}
      {/*          setPwdModalVisible(false);*/}
      {/*        }}*/}
      {/*      >*/}
      {/*        <Trans>取消</Trans>*/}
      {/*      </Button>*/}
      {/*    </Form.Action>*/}
      {/*  </Modal.Footer>*/}
      {/*</Modal>*/}
    </>
  );

  /** 渲染操作按钮 */
  function _renderOperationCell(user: User) {
    const isDisable = user.status.phase === 'Deleting';

    if (user.spec.username.toLowerCase() === 'admin') {
      return (
        <LinkButton tipDirection="right" errorTip="管理员不能被删除" disabled>
          <Trans>删除</Trans>
        </LinkButton>
      );
    }
    return (
      <React.Fragment>
        <LinkButton
          tipDirection="left"
          disabled={isDisable}
          onClick={() => {
            setEditUser(user);
            pwdToggle();
            // setPwdModalVisible(true);
          }}
        >
          {t('修改密码')}
        </LinkButton>

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
    );
  }

  async function _removeUser(user: User) {
    const yes = await Modal.confirm({
      message: t('确认删除当前所选用户？'),
      description: t('删除后，用户{{username}}的所有配置将会被清空，且无法恢复', { username: user.spec.username }),
      okText: t('删除'),
      cancelText: t('取消'),
    });
    if (yes) {
      actions.user.removeUser.start([user.metadata.name]);
      actions.user.removeUser.perform();
    }
  }

  async function _onPwdModalSubmit() {
    const { password } = passwords;
    await actions.user.updateUser.fetch({
      noCache: true,
      data: {
        user: {
          metadata: {
            name: editUser.metadata.name,
            resourceVersion: editUser.metadata.resourceVersion,
          },
          spec: Object.assign({}, editUser.spec, {
            hashedPassword: btoa(password),
          }),
        },
      },
    });
    setPasswords({
      password: '',
      rePassword: '',
    });
    setPwdModalVisible(false);
    setBtnDisabled(true);
  }
};
