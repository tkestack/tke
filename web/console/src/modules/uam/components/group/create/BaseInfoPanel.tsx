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
import { UserPlain } from '../../../models';
import { router } from '../../../router';

const { useState, useEffect, useRef } = React;
const { scrollable, selectable, removeable } = Table.addons;

export const BaseInfoPanel = props => {
  const state = useSelector(state => state);
  const dispatch = useDispatch();
  const { actions } = bindActionCreators({ actions: allActions }, dispatch);

  const { userPlainList, policyPlainList } = state;
  const userList = userPlainList.list.data.records || [];
  let strategyList = policyPlainList.list.data.records || [];
  strategyList = strategyList.filter(item => ['平台管理员', '平台用户', '租户'].includes(item.displayName) === false);
  const tenantID = strategyList.filter(item => item.displayName === '平台管理员').tenantID;

  const [inputValue, setInputValue] = useState('');
  const [targetKeys, setTargetKeys] = useState([]);
  const [userInputValue, setUserInputValue] = useState('');
  const [userTargetKeys, setUserTargetKeys] = useState([]);

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
    const { displayName, description, role } = values;
    actions.group.create.addGroupWorkflow.start([
      {
        id: uuid(),
        spec: {
          displayName,
          description,
          extra: {
            policies: role === 'custom' ? targetKeys.join(',') : role
          }
        },
        status: {
          users: userTargetKeys.map(id => ({
            id
          }))
        }
      }
    ]);
    actions.group.create.addGroupWorkflow.perform();
  }

  const { form, handleSubmit, validating, submitting } = useForm({
    onSubmit,
    /**
     * 默认为 shallowEqual
     * 如果初始值有多层，会导致重渲染，也可以使用 `useEffect` 设置初始值：
     * useEffect(() => form.initialize({ }), []);
     */
    initialValuesEqual: () => true,
    initialValues: { displayName: '', description: '', role: '' },
    validate: ({ displayName, description, role }) => {
      const errors = {
        displayName: undefined,
        description: undefined,
        role: undefined
      };
      if (!displayName) {
        errors.displayName = t('请输入用户账号');
      } else if (displayName.length > 60) {
        errors.displayName = t('请输入用户组名称，不超过60个字符');
      }

      if (description.length > 255) {
        errors.description = t('请输入用户组描述，不超过255个字符');
      }

      if (!role) {
        errors.role = t('请选择平台角色');
      }

      return errors;
    }
  });

  const displayName = useField('displayName', form);
  const description = useField('description', form);
  const role = useField('role', form);

  const roleValue = role.input.value;
  useEffect(() => {
    if (targetKeys.length > 0 && !roleValue) {
      form.change('role', 'custom');
    }
    if (roleValue && roleValue !== 'custom') {
      // 选择的时候是替换数组，所以引用不同，这里会被触发；这里清空的时候，让引用不变，所以这个useEffect不会被再次触发
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
              label={t('用户组名称')}
              required
              status={getStatus(displayName.meta, validating)}
              message={getStatus(displayName.meta, validating) === 'error' ? displayName.meta.error : ''}
            >
              <Input
                {...displayName.input}
                size="l"
                autoComplete="off"
                placeholder={t('请输入用户组名称，不超过60个字符')}
              />
            </Form.Item>
            <Form.Item
              label={t('用户组描述')}
              status={getStatus(description.meta, validating)}
              message={getStatus(description.meta, validating) === 'error' ? description.meta.error : ''}
            >
              <Input
                {...description.input}
                multiline
                size="l"
                autoComplete="off"
                placeholder={t('请输入用户组描述，不超过255个字符')}
              />
            </Form.Item>
            <Form.Item label={t('关联用户')}>
              {/*下边SearchBox的改造*/}
              {/*<SearchBox*/}
              {/*    value={userInputValue}*/}
              {/*    onChange={(keyword) => {*/}
              {/*      action.changeKeyword((keyword || '').trim());*/}
              {/*      setUserInputValue(keyword);*/}
              {/*    }}*/}
              {/*    onSearch={(keyword) => {*/}
              {/*      action.performSearch((keyword || '').trim());*/}
              {/*    }}*/}
              {/*    onClear={() => {*/}
              {/*      action.changeKeyword('');*/}
              {/*      action.performSearch('');*/}
              {/*      setUserInputValue('');*/}
              {/*    }}*/}
              {/*/>*/}
              <Transfer
                leftCell={
                  <Transfer.Cell
                    scrollable={false}
                    title="当前用户组可关联以下用户"
                    tip="支持按住 shift 键进行多选"
                    header={<SearchBox value={userInputValue} onChange={value => setUserInputValue(value)} />}
                  >
                    <UserAssociateSourceTable
                      dataSource={userList.filter(i => i.displayName.includes(userInputValue))}
                      targetKeys={userTargetKeys}
                      onChange={keys => setUserTargetKeys(keys)}
                    />
                  </Transfer.Cell>
                }
                rightCell={
                  <Transfer.Cell title={`已选择 (${userTargetKeys.length})`}>
                    <UserAssociateTargetTable
                      dataSource={userList.filter(i => userTargetKeys.includes(i.id))}
                      onRemove={key => setUserTargetKeys(userTargetKeys.filter(i => i !== key))}
                    />
                  </Transfer.Cell>
                }
              />
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
                  {roleValue === 'custom' && (
                    <Transfer
                      leftCell={
                        <Transfer.Cell
                          scrollable={false}
                          title="为这个用户自定义独立的权限"
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
                  router.navigate({ module: 'user', sub: 'group' });
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

const userAssociateColumns = [
  {
    key: 'name',
    header: t('ID/名称'),
    render: (user: UserPlain) => <p>{`${user.displayName}(${user.name})`}</p>
  }
];

function UserAssociateSourceTable({ dataSource, targetKeys, onChange }) {
  return (
    <Table
      records={dataSource}
      recordKey="id"
      columns={userAssociateColumns}
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

function UserAssociateTargetTable({ dataSource, onRemove }) {
  return (
    <Table records={dataSource} recordKey="id" columns={userAssociateColumns} addons={[removeable({ onRemove })]} />
  );
}

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
