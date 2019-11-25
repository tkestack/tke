import * as React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { TableColumn, Text, Modal, Form, Input, Button } from '@tea/component';
import { bindActionCreators } from '@tencent/qcloud-lib';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { router } from '../../router';
import { TablePanel, LinkButton } from '../../../common/components';
import { allActions } from '../../actions';
import { User } from '../../models';
import { VALIDATE_PASSWORD_RULE } from '../../constants/Config';

const { useState, useEffect } = React;

export const UserTablePanel = () => {
  const state = useSelector(state => state);
  const dispatch = useDispatch();
  const { actions } = bindActionCreators({ actions: allActions }, dispatch);

  const [pwdModalVisible, setPwdModalVisible] = useState(false);
  const [pwdMessages, setPwdMessages] = useState({ passwordMsg: '', rePasswordMsg: '' });
  const [passwords, setPasswords] = useState({ password: '', rePassword: '' });
  const [btnDisabled, setBtnDisabled] = useState(true);
  const [user, setUser] = useState(undefined);

  const { userList } = state;

  const columns: TableColumn<User>[] = [
    {
      key: 'name',
      header: t('用户ID / 名称'),
      render: (user, text, index) => (
        <Text parent="div" overflow>
          <a
            href="javascript:;"
            onClick={e => {
              // router.navigate({ module: 'user', sub: `${user.uid}` });
              router.navigate({ module: 'user', sub: `${user.name}` });
            }}
          >
            {user.name} / {user.Spec.extra.displayName || '-'}
          </a>
        </Text>
      )
    },
    {
      key: 'phone',
      header: t('关联手机'),
      render: user => <Text>{user.Spec.extra.phoneNumber || '-'}</Text>
    },
    {
      key: 'email',
      header: t('关联邮箱'),
      render: user => <Text>{user.Spec.extra.email || '-'}</Text>
    },
    { key: 'operation', header: t('操作'), render: user => _renderOperationCell(user) }
  ];
  const emptyTips: JSX.Element = (
    <div className="text-center">
      <Trans>暂无内容</Trans>
    </div>
  );

  return (
    <React.Fragment>
      <TablePanel
        columns={columns}
        model={userList}
        action={actions.user}
        emptyTips={emptyTips}
        isNeedPagination={true}
        bodyClassName={'tc-15-table-panel tc-15-table-fixed-body'}
      />
      <Modal
        visible={pwdModalVisible}
        size="s"
        caption={t('修改密码')}
        onClose={() => {
          setPwdModalVisible(false);
        }}
      >
        <Modal.Body>
          <Form>
            <Form.Item
              label={t('用户密码')}
              required
              status={pwdMessages.passwordMsg ? 'error' : passwords.password ? 'success' : undefined}
              message={pwdMessages.passwordMsg ? t(pwdMessages.passwordMsg) : t(VALIDATE_PASSWORD_RULE.message)}
            >
              <Input
                type="password"
                placeholder={t('请输入用户密码')}
                defaultValue={passwords.password}
                onChange={value => {
                  let pwdMsg = '';
                  let disabled = true;
                  if (!value) {
                    pwdMsg = '请输入用户密码';
                  } else if (!VALIDATE_PASSWORD_RULE.pattern.test(value)) {
                    pwdMsg = VALIDATE_PASSWORD_RULE.message;
                  } else if (passwords.rePassword && value !== passwords.rePassword) {
                    pwdMsg = '两次输入密码不一致';
                  }
                  if (value && passwords.rePassword && !pwdMsg && !pwdMessages.rePasswordMsg) {
                    disabled = false;
                  }
                  setPasswords({ ...passwords, password: value });
                  setPwdMessages({ ...pwdMessages, passwordMsg: pwdMsg });
                  setBtnDisabled(disabled);
                }}
              />
            </Form.Item>
            <Form.Item
              label={t('确认密码')}
              required
              status={pwdMessages.rePasswordMsg ? 'error' : passwords.rePassword ? 'success' : undefined}
              message={pwdMessages.rePasswordMsg ? t(pwdMessages.rePasswordMsg) : t(VALIDATE_PASSWORD_RULE.message)}
            >
              <Input
                type="password"
                placeholder={t('请再次输入用户密码')}
                defaultValue={passwords.rePassword}
                onChange={value => {
                  let rePwdMsg = '';
                  let disabled = true;
                  if (!value) {
                    rePwdMsg = '请输入用户密码';
                  } else if (!VALIDATE_PASSWORD_RULE.pattern.test(value)) {
                    rePwdMsg = VALIDATE_PASSWORD_RULE.message;
                  } else if (passwords.password && value !== passwords.password) {
                    rePwdMsg = '两次输入密码不一致';
                  }
                  if (passwords.password && value && !pwdMessages.passwordMsg && !rePwdMsg) {
                    disabled = false;
                  }
                  setPasswords({ ...passwords, rePassword: value });
                  setPwdMessages({ ...pwdMessages, rePasswordMsg: rePwdMsg });
                  setBtnDisabled(disabled);
                }}
              />
            </Form.Item>
          </Form>
        </Modal.Body>
        <Modal.Footer>
          <Form.Action>
            <Button disabled={btnDisabled} type="primary" onClick={_onPwdModalSubmit}>
              <Trans>保存</Trans>
            </Button>
            <Button
              onClick={() => {
                setPwdModalVisible(false);
              }}
            >
              <Trans>取消</Trans>
            </Button>
          </Form.Action>
        </Modal.Footer>
      </Modal>
    </React.Fragment>
  );

  /** 渲染操作按钮 */
  function _renderOperationCell(user: User) {
    if (user.name.toLowerCase() === 'admin') {
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
          onClick={() => {
            setUser(user);
            setPwdModalVisible(true);
          }}
        >
          {t('修改密码')}
        </LinkButton>

        <LinkButton
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
      description: t('删除后，用户{{username}}的所有配置将会被清空，且无法恢复', { username: user.name }),
      okText: t('删除'),
      cancelText: t('取消')
    });
    if (yes) {
      actions.user.removeUser.start([user.name]);
      actions.user.removeUser.perform();
    }
  }

  async function _onPwdModalSubmit() {
    const { password } = passwords;
    await actions.user.updateUser.fetch({
      noCache: true,
      data: {
        user: {
          name: user.name,
          Spec: {
            hashedPassword: btoa(password)
          }
        }
      }
    });
    setPasswords({
      password: '',
      rePassword: ''
    });
    setPwdModalVisible(false);
    setBtnDisabled(true);
  }
};
