import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';

import { Button, Form, Icon, Input, Modal, TableColumn, Text } from '@tea/component';
import { TablePanel } from '@tencent/ff-component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { LinkButton } from '../../../common/components';
import { allActions } from '../../actions';
import { VALIDATE_PASSWORD_RULE } from '../../constants/Config';
import { User } from '../../models';
import { router } from '../../router';

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
              router.navigate({ module: 'user', sub: `${user.metadata.name}` });
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
        recordKey={(record) => {
          return record.metadata.name;
        }}
        columns={columns}
        rowDisabled={record => record.status && record.status.phase === 'Deleting'}
        model={userList}
        action={actions.user}
        emptyTips={emptyTips}
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
            setUser(user);
            setPwdModalVisible(true);
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
      cancelText: t('取消')
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
            name: user.metadata.name,
            resourceVersion: user.metadata.resourceVersion
          },
          spec: Object.assign({}, user.spec, {
            hashedPassword: btoa(password)
          })
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
