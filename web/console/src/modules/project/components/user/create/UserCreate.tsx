import * as React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators, insertCSS, uuid } from '@tencent/ff-redux';
import { useForm, useField } from 'react-final-form-hooks';
import { allActions } from '../../../actions';
import { Button, Text, Form, Affix, Card, Radio, Transfer, Table, SearchBox } from '@tencent/tea-component';
import { router } from '../../../router';
import { User, UserPlain } from '../../../models';
import { getStatus } from '../../../../common/validate';

const { useState, useEffect, useRef } = React;
const { scrollable, selectable, removeable } = Table.addons;

export const UserCreate = props => {
  const state = useSelector(state => state);
  const dispatch = useDispatch();
  const { actions } = bindActionCreators({ actions: allActions }, dispatch);
  const { route, manager, policyPlainList } = state;
  const userList = manager.list.data.records || [];
  let strategyList = policyPlainList.list.data.records || [];
  strategyList = strategyList.filter(
    item => ['业务管理员', '业务成员', '业务只读'].includes(item.displayName) === false
  );
  const tenantID = strategyList.filter(item => item.displayName === '业务管理员').tenantID;

  const [inputValue, setInputValue] = useState('');
  const [targetKeys, setTargetKeys] = useState([]);
  const [userInputValue, setUserInputValue] = useState('');
  const [userTargetKeys, setUserTargetKeys] = useState([]);

  useEffect(() => {
    actions.manager.fetch();
  }, []);

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
    // console.log('submit .....', values, targetKeys, userTargetKeys);
    const { role } = values;
    let userInfo = {
      id: uuid(),
      projectId: route.queries.projectId,
      users: userTargetKeys.map(id => ({
        id
      })),
      policies: role === 'custom' ? targetKeys : [role]
    };
    // console.log('submit userInfo: ', userInfo);
    actions.user.addUser.start([userInfo]);
    actions.user.addUser.perform();
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
    validate: ({ displayName, description, role }) => {
      const errors = {
        role: undefined
      };

      if (!role) {
        errors.role = t('请选择平台角色');
      }

      return errors;
    }
  });

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
  // onChange={(value) => setUserInputValue(value)}
  return (
    <form onSubmit={handleSubmit}>
      <Card>
        <Card.Body>
          <Form>
            <Form.Item label={t('编辑成员')}>
              <Transfer
                leftCell={
                  <Transfer.Cell
                    scrollable={false}
                    title="当前账户可分配以下责任人"
                    tip="支持按住 shift 键进行多选"
                    header={
                      <SearchBox
                        value={userInputValue}
                        onChange={keyword => {
                          setUserInputValue(keyword);
                          // actions.manager.changeKeyword((keyword || '').trim());
                          actions.manager.performSearch((keyword || '').trim());
                        }}
                        onSearch={keyword => {
                          actions.manager.performSearch((keyword || '').trim());
                        }}
                        onClear={() => {
                          actions.manager.changeKeyword('');
                          actions.manager.performSearch('');
                        }}
                      />
                    }
                  >
                    <UserAssociateSourceTable
                      dataSource={userList}
                      // dataSource={userList.filter((i) => i.displayName.includes(userInputValue))}
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
              label={t('业务角色')}
              required
              status={getStatus(role.meta, validating)}
              message={getStatus(role.meta, validating) === 'error' ? role.meta.error : ''}
            >
              <Radio.Group {...role.input} layout="column">
                <Radio name={tenantID ? `pol-${tenantID}-project-owner` : 'pol-default-project-owner'}>
                  <Text>业务管理员</Text>
                  <Text parent="div">预设业务角色，允许管理业务自身和业务下的所有功能和资源</Text>
                </Radio>
                <Radio name={tenantID ? `pol-${tenantID}-project-member` : 'pol-default-project-member'}>
                  <Text>业务成员</Text>
                  <Text parent="div">预设业务角色，允许访问和管理所在业务下的所有功能和资源</Text>
                </Radio>
                <Radio name={tenantID ? `pol-${tenantID}-project-viewer` : 'pol-default-project-viewer'}>
                  <Text>只读成员</Text>
                  <Text parent="div">预设业务角色，仅能够查看业务下资源</Text>
                </Radio>
                <Radio name="custom">
                  <Text>自定义</Text>
                  <Transfer
                    leftCell={
                      <Transfer.Cell
                        scrollable={false}
                        title="为这个用户选择单个角色"
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
                  history.back();
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
