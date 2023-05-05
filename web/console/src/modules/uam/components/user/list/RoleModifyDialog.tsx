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
import { Button, Form, Modal, Radio, SearchBox, Table, Text, Transfer } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';
import { Trans, t } from '@tencent/tea-app/lib/i18n';
import * as React from 'react';
import { useField, useForm } from 'react-final-form-hooks';
import { useDispatch, useSelector } from 'react-redux';
import { getStatus } from '../../../../common/validate';
import { allActions } from '../../../actions';

const { useState, useEffect } = React;
const { scrollable, selectable, removeable } = Table.addons;

export function RoleModifyDialog(props) {
  const state = useSelector(state => state);
  const dispatch = useDispatch();
  const { actions } = bindActionCreators({ actions: allActions }, dispatch);

  const { policyPlainList } = state;
  let strategyList = policyPlainList.list.data.records || [];
  strategyList = strategyList.filter(item => ['平台管理员', '平台用户', '租户'].includes(item.displayName) === false);

  const { isShowing, toggle, user } = props;

  const [targetKeys, setTargetKeys] = useState([]);
  const [inputValue, setInputValue] = useState('');
  const [tenantID, setTenantID] = useState('default');

  function onSubmit(values, form) {
    console.log('RoleModifyDialog submit values:', values, targetKeys);
    const { role } = values;
    const extraObj =
      role === 'custom'
        ? { ...user.spec.extra, policies: targetKeys.join(',') }
        : { ...user.spec.extra, policies: role };
    actions.user.updateUser.fetch({
      noCache: true,
      data: {
        user: {
          metadata: {
            name: user.metadata.name,
            resourceVersion: user.metadata.resourceVersion
          },
          spec: { ...user.spec, extra: extraObj }
        }
      }
    });
    setTimeout(form.reset);
    toggle();
  }
  const { form, handleSubmit, validating, submitting } = useForm({
    onSubmit,
    /**
     * 默认为 shallowEqual
     * 如果初始值有多层，会导致重渲染，也可以使用 `useEffect` 设置初始值：
     * useEffect(() => form.initialize({ }), []);
     */
    initialValuesEqual: () => true,
    initialValues: { role: '' },
    validate: ({ role }) => ({
      role: !role ? t('请选择平台角色') : undefined
    })
  });
  const role = useField('role', form);

  useEffect(() => {
    if (user) {
      const {
        tenantID,
        extra: { policies }
      } = user.spec;
      setTenantID(tenantID);
      const policiesParse = JSON.parse(policies);
      const keys = Object.keys(policiesParse);
      const roleArray = [`pol-${tenantID}-administrator`, `pol-${tenantID}-platform`, `pol-${tenantID}-viewer`];
      if (keys.length === 1 && roleArray.includes(keys[0])) {
        form.change('role', keys[0]);
        setTargetKeys([]);
      } else if (keys.length >= 1) {
        form.change('role', 'custom');
        setTargetKeys(keys);
      }
    }
  }, [user]);

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
    <Modal
      visible={isShowing}
      size="l"
      caption={t('选择平台角色')}
      onClose={() => {
        toggle();
        setTimeout(form.reset);
      }}
    >
      <Modal.Body>
        <form onSubmit={handleSubmit}>
          <Form>
            <Form.Item
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
                            dataSource={strategyList}
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
          <Form.Action style={{ textAlign: 'center' }}>
            <Button type="primary" htmlType="submit" loading={submitting} disabled={validating}>
              <Trans>确定</Trans>
            </Button>
            <Button
              type="weak"
              htmlType="reset"
              onClick={() => {
                toggle();
              }}
            >
              <Trans>取消</Trans>
            </Button>
          </Form.Action>
        </form>
      </Modal.Body>
    </Modal>
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
    width: 150,
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
