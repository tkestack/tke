/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */
import { PermissionProvider } from '@common';
import { bindActionCreators, uuid } from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';
import { Affix, Button, Card, Form, Input, Radio, SearchBox, Table, Text, Transfer } from '@tencent/tea-component';
import * as React from 'react';
import { useField, useForm } from 'react-final-form-hooks';
import { useDispatch, useSelector } from 'react-redux';
import { getStatus } from '../../../../common/validate';
import { allActions } from '../../../actions';
import {
  VALIDATE_EMAIL_RULE,
  VALIDATE_NAME_RULE,
  VALIDATE_PASSWORD_RULE,
  VALIDATE_PHONE_RULE
} from '../../../constants/Config';
import { User } from '../../../models';
import { router } from '../../../router';

const { useState, useEffect, useRef } = React;
const { scrollable, selectable, removeable } = Table.addons;

export const UserCreate = props => {
  const state = useSelector(state => state);
  const dispatch = useDispatch();
  const { actions } = bindActionCreators({ actions: allActions }, dispatch);

  const { filterUsers, policyPlainList } = state;
  let strategyList = policyPlainList.list.data.records || [];
  strategyList = strategyList.filter(item => ['平台管理员', '平台用户', '租户'].includes(item.displayName) === false);
  const tenantID = strategyList.filter(item => item.displayName === '平台管理员').tenantID;

  const [targetKeys, setTargetKeys] = useState([]);
  const [inputValue, setInputValue] = useState('');

  // 处理外层滚动
  const bottomAffixRef = useRef(null);
  useEffect(() => {
    const body = document.querySelector('.tea-web-body');
    if (!body) {
      return () => null;
    }
    const handleScroll = () => {
      bottomAffixRef.current.update();
    };
    body.addEventListener('scroll', handleScroll);
    return () => body.removeEventListener('scroll', handleScroll);
  }, []);

  function onSubmit(values, form) {
    console.log('submit .....', values, targetKeys);
    const { name, displayName, password, phone, email, role } = values;
    const userInfo: User = {
      id: uuid(),
      spec: {
        username: name,
        hashedPassword: btoa(password),
        displayName,
        email,
        phoneNumber: phone,
        extra: {
          policies: role === 'custom' ? targetKeys.join(',') : role
        }
      }
    };
    console.log('submit userInfo: ', userInfo);
    actions.user.addUser.start([userInfo]);
    actions.user.addUser.perform();
    // router.navigate({ module: 'user' });
    // setTimeout(form.reset);
  }

  const { form, handleSubmit, validating, submitting } = useForm({
    onSubmit,
    /**
     * 默认为 shallowEqual
     * 如果初始值有多层，会导致重渲染，也可以使用 `useEffect` 设置初始值：
     * useEffect(() => form.initialize({ }), []);
     */
    initialValuesEqual: () => true,
    initialValues: { name: '', displayName: '', password: '', rePassword: '', phone: '', email: '', role: '' },
    validate: ({ name, displayName, password, rePassword, phone, email, role }) => {
      const errors = {
        name: undefined,
        displayName: undefined,
        password: undefined,
        rePassword: undefined,
        phone: undefined,
        email: undefined,
        role: undefined
      };
      if (!name) {
        errors.name = t('请输入用户账号');
      } else if (!VALIDATE_NAME_RULE.pattern.test(name)) {
        errors.name = VALIDATE_NAME_RULE.message;
      } else if (filterUsers && filterUsers.length) {
        errors.name = t('账号名被占用，请更换新的尝试');
      }

      if (!displayName) {
        errors.displayName = t('请输入用户名称');
      } else if (displayName.length >= 256) {
        errors.displayName = t('长度需要小于256个字符');
      }

      if (!password) {
        errors.password = t('请输入密码');
      } else if (!VALIDATE_PASSWORD_RULE.pattern.test(password)) {
        errors.password = VALIDATE_PASSWORD_RULE.message;
      }

      if (!rePassword) {
        errors.rePassword = t('请再次输入密码');
      } else if (!VALIDATE_PASSWORD_RULE.pattern.test(rePassword)) {
        errors.rePassword = VALIDATE_PASSWORD_RULE.message;
      } else if (password !== rePassword) {
        errors.rePassword = t('两次输入密码需一致');
      }

      if (!phone) {
        errors.phone = undefined;
      } else if (!VALIDATE_PHONE_RULE.pattern.test(phone)) {
        errors.phone = VALIDATE_PHONE_RULE.message;
      }

      if (!email) {
        errors.email = undefined;
      } else if (!VALIDATE_EMAIL_RULE.pattern.test(email)) {
        errors.email = VALIDATE_EMAIL_RULE.message;
      }

      if (!role) {
        errors.role = t('请选择平台角色');
      }

      return errors;
    }
  });

  const name = useField('name', form);
  const displayName = useField('displayName', form);
  const password = useField('password', form);
  const rePassword = useField('rePassword', form);
  const phone = useField('phone', form);
  const email = useField('email', form);
  const role = useField('role', form);

  const roleValue = role.input.value;
  useEffect(() => {
    if (targetKeys.length > 0 && !roleValue) {
      form.change('role', 'custom');
    }
    if (roleValue && roleValue !== 'custom') {
      // 对于targetKeys选择的时候是替换数组，所以引用不同，这里会被触发；这里清空的时候，让引用不变，所以这个useEffect不会被再次触发
      const newTargetKeys = targetKeys;
      newTargetKeys.length = 0;
      setTargetKeys(newTargetKeys);
    }
  }, [roleValue, targetKeys]);

  return (
    <form onSubmit={handleSubmit}>
      <Card>
        <Card.Body>
          <Form>
            <Form.Item
              label={t('用户账号')}
              required
              status={getStatus(name.meta, validating)}
              message={getStatus(name.meta, validating) === 'error' ? name.meta.error : VALIDATE_NAME_RULE.message}
            >
              <Input
                {...name.input}
                onChange={value => {
                  name.input.onChange(value);
                  actions.user.getUsersByName(value);
                }}
                size="l"
                autoComplete="off"
                placeholder={t('请输入用户账号')}
              />
            </Form.Item>
            <Form.Item
              label={t('用户名称')}
              required
              status={getStatus(displayName.meta, validating)}
              message={
                getStatus(displayName.meta, validating) === 'error'
                  ? displayName.meta.error
                  : t('长度需要小于256个字符')
              }
            >
              <Input {...displayName.input} size="l" autoComplete="off" placeholder={t('请输入用户名称')} />
            </Form.Item>
            <Form.Item
              label={t('用户密码')}
              required
              status={getStatus(password.meta, validating)}
              message={
                getStatus(password.meta, validating) === 'error' ? password.meta.error : VALIDATE_PASSWORD_RULE.message
              }
            >
              <Input {...password.input} type="password" size="l" autoComplete="off" placeholder={t('请输入密码')} />
            </Form.Item>
            <Form.Item
              label={t('确认密码')}
              required
              status={getStatus(rePassword.meta, validating)}
              message={
                getStatus(rePassword.meta, validating) === 'error'
                  ? rePassword.meta.error
                  : VALIDATE_PASSWORD_RULE.message
              }
            >
              <Input
                {...rePassword.input}
                type="password"
                size="l"
                autoComplete="off"
                placeholder={t('请再次输入密码')}
              />
            </Form.Item>
            <Form.Item
              label={t('手机号')}
              status={getStatus(phone.meta, validating)}
              message={getStatus(phone.meta, validating) === 'error' ? phone.meta.error : ''}
            >
              <Input {...phone.input} size="l" autoComplete="off" placeholder={t('请输入用户手机号')} />
            </Form.Item>
            <Form.Item
              label={t('邮箱')}
              status={getStatus(email.meta, validating)}
              message={getStatus(email.meta, validating) === 'error' ? email.meta.error : ''}
            >
              <Input {...email.input} size="l" autoComplete="off" placeholder={t('请输入用户邮箱')} />
            </Form.Item>
            <Form.Item
              label={t('平台角色')}
              required
              status={getStatus(role.meta, validating)}
              message={getStatus(role.meta, validating) === 'error' ? role.meta.error : ''}
            >
              <Radio.Group {...role.input} layout="column">
                <Radio name={tenantID ? `pol-${tenantID}-administrator` : 'pol-default-administrator'}>
                  <Text>管理员</Text>
                  <Text parent="div">平台预设角色，允许访问和管理所有平台和业务的功能和资源</Text>
                </Radio>
                <Radio name={tenantID ? `pol-${tenantID}-platform` : 'pol-default-platform'}>
                  <Text>平台用户</Text>
                  <Text parent="div">平台预设角色，允许访问和管理大部分平台功能，可以新建集群及业务</Text>
                </Radio>
                <PermissionProvider value="platform.uam.tenant">
                  <Radio name={tenantID ? `pol-${tenantID}-viewer` : 'pol-default-viewer'}>
                    <Text>租户</Text>
                    <Text parent="div">平台预设角色，不绑定任何平台权限，仅能登录</Text>
                  </Radio>
                </PermissionProvider>
                <Radio name="custom">
                  <Text>自定义</Text>
                  <Text parent="div">为这个用户自定义独立的权限</Text>
                  {roleValue === 'custom' && (
                    <Transfer
                      leftCell={
                        <Transfer.Cell
                          scrollable={false}
                          tip="支持按住 shift 键进行多选"
                          header={<SearchBox value={inputValue} onChange={value => setInputValue(value)} />}
                        >
                          <SourceTable
                            dataSource={strategyList.filter(i => i.displayName.includes(inputValue))}
                            targetKeys={targetKeys}
                            onChange={keys => setTargetKeys(keys)}
                          />
                        </Transfer.Cell>
                      }
                      rightCell={
                        <Transfer.Cell title={`已选择 (${targetKeys.length})`}>
                          <TargetTable
                            dataSource={strategyList.filter(i => targetKeys.includes(i.id))}
                            onRemove={key => setTargetKeys(targetKeys.filter(i => i !== key))}
                          />
                        </Transfer.Cell>
                      }
                    />
                  )}
                </Radio>
              </Radio.Group>
            </Form.Item>
          </Form>
        </Card.Body>
      </Card>
      <Affix ref={bottomAffixRef} offsetBottom={0} style={{ zIndex: 5 }}>
        <Card>
          <Card.Body style={{ borderTop: '1px solid #ddd' }}>
            <Form.Action style={{ borderTop: 0, marginTop: 0, paddingTop: 0 }}>
              <Button type="primary">保存</Button>
              <Button
                onClick={e => {
                  e.preventDefault();
                  router.navigate({ module: 'user' });
                }}
              >
                取消
              </Button>
            </Form.Action>
          </Card.Body>
        </Card>
      </Affix>
    </form>
  );
};

const columns = [
  {
    key: 'displayName',
    header: '策略名称',
    render: strategy => <p>{strategy.displayName}</p>
  },
  {
    key: 'description',
    header: '描述',
    width: 300,
    render: strategy => <p>{strategy.description || '-'}</p>
  }
];

function SourceTable({ dataSource, targetKeys, onChange }) {
  return (
    <Table
      records={dataSource}
      recordKey="id"
      columns={columns}
      addons={[
        scrollable({
          maxHeight: 310,
          onScrollBottom: () => console.log('到达底部')
        }),
        selectable({
          value: targetKeys,
          onChange,
          rowSelect: true
        })
      ]}
    />
  );
}

function TargetTable({ dataSource, onRemove }) {
  return <Table records={dataSource} recordKey="id" columns={columns} addons={[removeable({ onRemove })]} />;
}
