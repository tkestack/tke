import * as React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { bindActionCreators, uuid, insertCSS } from '@tencent/qcloud-lib';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Justify } from '@tea/component/justify';
import { allActions } from '../../actions';
import { User, UserFilter } from '../../models';
import { Modal, Button, Form, Input, SearchBox, Switch, Table, Icon } from '@tea/component';
import {
  VALIDATE_PASSWORD_RULE,
  VALIDATE_PHONE_RULE,
  VALIDATE_EMAIL_RULE,
  VALIDATE_NAME_RULE
} from '../../constants/Config';

const { useState, useEffect } = React;
const _isEqual = require('lodash/isEqual');

insertCSS(
  'UserActionPanel',
  `
    .add-user-form input {
      width: 350px;
    }
    .add-user-form .tea-form__help-text {
      height: 18px
    }
`
);

export const UserActionPanel = () => {
  const state = useSelector(state => state);
  const dispatch = useDispatch();
  const { actions } = bindActionCreators({ actions: allActions }, dispatch);

  const { userList, filterUsers } = state;

  const [modalVisible, setModalVisible] = useState(false);
  const [modalBtnDisabled, setModalBtnDisabled] = useState(true);
  const [formParamsValue, setFormParamsValue] = useState({
    name: '',
    displayName: '',
    password: '',
    rePassword: '',
    phone: '',
    email: ''
  });
  const [messages, setMessages] = useState({
    displayName: '',
    name: '',
    password: '',
    rePassword: '',
    phone: '',
    email: ''
  });

  useEffect(() => {
    // 用来判断新增用户时name是否重复
    if (filterUsers && filterUsers.length) {
      setMessages({ ...messages, name: '账号名被占用，请更换新的尝试' });
    }
  }, [filterUsers, messages]);

  useEffect(() => {
    const { displayName, name, password, rePassword } = formParamsValue;
    let disabled = false;
    Object.keys(messages).forEach(item => {
      if (messages[item]) {
        disabled = true;
      }
    });
    if (!displayName || !name || !password || !rePassword) {
      disabled = true;
    }
    setModalBtnDisabled(disabled);
  }, [formParamsValue, messages]);

  return (
    <React.Fragment>
      <Table.ActionPanel>
        <Justify
          left={
            <Button type="primary" onClick={_open}>
              {t('新建')}
            </Button>
          }
          right={
            <React.Fragment>
              <SearchBox
                value={userList.query.keyword || ''}
                onChange={actions.user.changeKeyword}
                onSearch={actions.user.performSearch}
                onClear={() => {
                  actions.user.performSearch('');
                }}
                placeholder={t('请输入用户名称')}
              />
            </React.Fragment>
          }
        />
      </Table.ActionPanel>
      <Modal visible={modalVisible} caption={t('添加用户')} onClose={_close}>
        <Modal.Body>
          <Form className="add-user-form">
            <Form.Item
              label={t('用户账号')}
              required
              status={messages.name ? 'error' : name ? 'success' : undefined}
              message={messages.name ? t(messages.name) : t(VALIDATE_NAME_RULE.message)}
            >
              <Input
                placeholder={t('请输入用户名称')}
                defaultValue={formParamsValue.name}
                onChange={async value => {
                  let msg = '';
                  if (!value) {
                    msg = '请输入用户账号';
                  } else if (!VALIDATE_NAME_RULE.pattern.test(value)) {
                    msg = VALIDATE_NAME_RULE.message;
                  }
                  setFormParamsValue({ ...formParamsValue, name: value });
                  setMessages({ ...messages, name: msg });
                }}
                onBlur={e => {
                  const value = e.target.value;
                  actions.user.getUsersByName(value);
                  let msg = '';
                  if (!value) {
                    msg = '请输入用户账号';
                  } else if (!VALIDATE_NAME_RULE.pattern.test(value)) {
                    msg = VALIDATE_NAME_RULE.message;
                  }
                  setFormParamsValue({ ...formParamsValue, name: value });
                  setMessages({ ...messages, name: msg });
                }}
              />
            </Form.Item>
            <Form.Item
              label={t('用户名称')}
              required
              status={messages.displayName ? 'error' : formParamsValue.displayName ? 'success' : undefined}
              message={messages.displayName ? t(messages.displayName) : t('长度需要小于256个字符')}
            >
              <Input
                placeholder={t('请输入用户名称')}
                defaultValue={formParamsValue.displayName}
                onChange={value => {
                  let msg = '';
                  if (!value) {
                    msg = '请输入用户名称';
                  } else if (value.length > 255) {
                    msg = '长度需要小于256个字符';
                  }
                  setFormParamsValue({ ...formParamsValue, displayName: value });
                  setMessages({ ...messages, displayName: msg });
                }}
              />
            </Form.Item>
            <Form.Item
              label={t('用户密码')}
              required
              status={messages.password ? 'error' : formParamsValue.password ? 'success' : undefined}
              message={messages.password ? t(messages.password) : t(VALIDATE_PASSWORD_RULE.message)}
            >
              <Input
                type="password"
                placeholder={t('请输入用户密码')}
                defaultValue={formParamsValue.password}
                onChange={value => {
                  let passwordMsg = '';
                  if (!value) {
                    passwordMsg = '请输入用户密码';
                  } else if (!VALIDATE_PASSWORD_RULE.pattern.test(value)) {
                    passwordMsg = VALIDATE_PASSWORD_RULE.message;
                  } else if (formParamsValue.rePassword && value !== formParamsValue.rePassword) {
                    passwordMsg = '两次输入密码不一致';
                  }
                  setFormParamsValue({ ...formParamsValue, password: value });
                  setMessages({ ...messages, password: passwordMsg });
                }}
              />
            </Form.Item>
            <Form.Item
              label={t('确认密码')}
              required
              status={messages.rePassword ? 'error' : formParamsValue.rePassword ? 'success' : undefined}
              message={messages.rePassword ? t(messages.rePassword) : t(VALIDATE_PASSWORD_RULE.message)}
            >
              <Input
                type="password"
                placeholder={t('请再次输入用户密码')}
                defaultValue={formParamsValue.rePassword}
                onChange={value => {
                  let passwordMsg = '';
                  if (!value) {
                    passwordMsg = '请输入用户密码';
                  } else if (!VALIDATE_PASSWORD_RULE.pattern.test(value)) {
                    passwordMsg = VALIDATE_PASSWORD_RULE.message;
                  } else if (formParamsValue.password && value !== formParamsValue.password) {
                    passwordMsg = '两次输入密码不一致';
                  }
                  setFormParamsValue({ ...formParamsValue, rePassword: value });
                  setMessages({ ...messages, rePassword: passwordMsg });
                }}
              />
            </Form.Item>
            <Form.Item
              label={t('手机号')}
              status={messages.phone ? 'error' : formParamsValue.phone ? 'success' : undefined}
              message={messages.phone ? t(messages.phone) : ' '}
            >
              <Input
                placeholder={t('请输入用户手机号')}
                defaultValue={formParamsValue.phone}
                onChange={value => {
                  setFormParamsValue({ ...formParamsValue, phone: value });
                  setMessages({
                    ...messages,
                    phone: !value || VALIDATE_PHONE_RULE.pattern.test(value) ? '' : '请输入正确的电话号码'
                  });
                }}
              />
            </Form.Item>
            <Form.Item
              label={t('邮箱')}
              status={messages.email ? 'error' : formParamsValue.email ? 'success' : undefined}
              message={messages.email ? t(messages.email) : ' '}
            >
              <Input
                placeholder={t('请输入用户邮箱')}
                defaultValue={formParamsValue.email}
                onChange={value => {
                  setFormParamsValue({ ...formParamsValue, email: value });
                  setMessages({
                    ...messages,
                    email: !value || VALIDATE_EMAIL_RULE.pattern.test(value) ? '' : '请输入正确的邮箱'
                  });
                }}
              />
            </Form.Item>
          </Form>
        </Modal.Body>
        <Modal.Footer>
          <Form.Action>
            <Button disabled={modalBtnDisabled} type="primary" onClick={_onSubmit}>
              <Trans>保存</Trans>
            </Button>
            <Button onClick={_close}>
              <Trans>取消</Trans>
            </Button>
          </Form.Action>
        </Modal.Footer>
      </Modal>
    </React.Fragment>
  );

  function _open() {
    setModalVisible(true);
  }
  function _close() {
    setModalVisible(false);
  }
  function _onSubmit() {
    const { displayName, name, password, phone, email } = formParamsValue;

    let userInfo: User = {
      id: uuid(),
      spec: {
        username: name,
        hashedPassword: btoa(password),
        displayName,
        email,
        phoneNumber: phone
      }
    };
    actions.user.addUser.start([userInfo]);
    actions.user.addUser.perform();
    setModalVisible(false);
    setModalBtnDisabled(true);
    setFormParamsValue({ name: '', displayName: '', password: '', rePassword: '', phone: '', email: '' });
  }
};
