import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';

import {
    Button, Card, CodeEditor, Input, LoadingTip, Modal, SearchBox, Table, TableColumn, TabPanel,
    Tabs, Text, Tooltip, Transfer
} from '@tea/component';
import { removeable, selectable } from '@tea/component/table/addons';
import { TablePanel } from '@tencent/ff-component';
import { bindActionCreators, insertCSS } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { dateFormat } from '../../../../../helpers/dateUtil';
import { LinkButton, usePrevious } from '../../../common/components';
import { allActions } from '../../actions';
import { User } from '../../models';
import { router } from '../../router';
import { GroupActionPanel } from './detail/GroupActionPanel';
import { GroupTablePanel } from './detail/GroupTablePanel';
import { RoleActionPanel } from './detail/RoleActionPanel';
import { RoleTablePanel } from './detail/RoleTablePanel';

const { useState, useEffect, useRef } = React;
const _isEqual = require('lodash/isEqual');

insertCSS(
  'StrategyDetailsPanel',
  `
    .item-descr-list .is-error {
      color: #e1504a;
      border-color: #e1504a;
    }
`
);

let editorRef;
const tabs = [
  { id: 'actions', label: '策略语法' },
  { id: 'users', label: '关联用户' },
  { id: 'groups', label: '关联用户组' },
  { id: 'roles', label: '已关联角色' }
];
export const StrategyDetailsPanel = () => {
  const state = useSelector(state => state);
  const dispatch = useDispatch();
  const { actions } = bindActionCreators({ actions: allActions }, dispatch);

  const { route, associatedUsersList, userList, getStrategy, updateStrategy } = state;
  const associatedUsersListRecords = associatedUsersList.list.data.records.map(item => item.metadata.name);
  const userListRecords = userList.list.data.records;
  const getStrategyData = getStrategy.data[0];
  const updateStrategyData = updateStrategy.data[0];

  const { sub } = router.resolve(route);

  const [modalVisible, setModalVisible] = useState(false);
  const [basicParamsValue, setBasicParamsValue] = useState({ name: '', description: '' });
  const [userMsgsValue, setUserMsgsValue] = useState({
    inputValue: '',
    targetKeys: associatedUsersListRecords,
    newTargetKeys: []
  });
  const [editValue, setEditValue] = useState({ editorStatement: {}, editBasic: false });
  const [editorValue, setEditorValue] = useState({ ready: false, readOnly: true });
  const [strategy, setStrategy] = useState(undefined);

  useEffect(() => {
    // 请求策略详情
    actions.strategy.getStrategy.fetch({
      noCache: true,
      data: { id: sub }
    });

    return () => {
      actions.user.performSearch('');
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    // 初始化策略详情
    if (getStrategyData && getStrategyData.target.metadata.name === sub) {
      const showStrategy = getStrategyData.target;
      setStrategy(showStrategy);
      setBasicParamsValue({ name: showStrategy.spec.displayName, description: showStrategy.spec.description });
    }
  }, [getStrategyData, sub]);

  useEffect(() => {
    // 更新strategy
    if (updateStrategyData && updateStrategyData.success && !_isEqual(strategy, updateStrategyData.target)) {
      const showStrategy = updateStrategyData.target;
      setStrategy(showStrategy);
      // setBasicParamsValue({ name: showStrategy.name, description: showStrategy.description });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [updateStrategyData]);

  useEffect(() => {
    // 关联用户
    if (!_isEqual(associatedUsersListRecords, userMsgsValue.targetKeys)) {
      setUserMsgsValue({ ...userMsgsValue, targetKeys: associatedUsersListRecords });
    }
  }, [associatedUsersListRecords, userMsgsValue]);

  const { name, description } = basicParamsValue;
  const { displayName: newName = '' } = strategy ? strategy.spec : {};
  const { description: newDescription = '', email: pEmail = '' } = strategy ? strategy.spec : {};
  const isNameError = name.length <= 0 || name.length > 255;
  const enabled = (name !== newName || description !== newDescription) && !isNameError;
  const modalColumns = [
    {
      key: 'name',
      header: '用户',
      render: user => {
        if (userMsgsValue.targetKeys.indexOf(user.metadata.name) !== -1) {
          return (
            <Tooltip title="用户已被关联">
              {user.spec.name}({user.spec.displayName})
            </Tooltip>
          );
        }
        return (
          <p>
            {user.spec.name}({user.spec.displayName})
          </p>
        );
      }
    }
  ];
  const columns: TableColumn<User>[] = [
    {
      key: 'name',
      header: t('关联用户'),
      render: record => <Text parent="div">{record.spec.displayName}</Text>
    },
    { key: 'operation', header: t('操作'), render: record => _renderOperationCell(record.metadata.name) }
  ];
  return (
    <React.Fragment>
      <Card>
        <Card.Body
          title={t('基本信息')}
          subtitle={
            strategy && strategy.type !== 1 ? (
              <Button type="link" onClick={_onBasicEdit}>
                编辑
              </Button>
            ) : (
              ''
            )
          }
        >
          {strategy && (
            <ul className="item-descr-list">
              <li>
                <span className="item-descr-tit">策略id</span>
                <span className="item-descr-txt">{strategy.metadata.name}</span>
              </li>
              <li>
                <span className="item-descr-tit">策略名称</span>
                {editValue.editBasic ? (
                  <Input
                    value={name}
                    className={isNameError && 'is-error'}
                    onChange={value => {
                      setBasicParamsValue({ ...basicParamsValue, name: value });
                    }}
                  />
                ) : (
                  <span className="item-descr-txt">{strategy.spec.displayName || '-'}</span>
                )}
                {editValue.editBasic && isNameError && <p className="is-error">输入不能为空且需要小于256个字符</p>}
              </li>
              <li>
                <span className="item-descr-tit">描述</span>
                {editValue.editBasic ? (
                  <Input
                    multiline
                    value={description}
                    onChange={value => setBasicParamsValue({ ...basicParamsValue, description: value })}
                  />
                ) : (
                  <span className="item-descr-txt">{strategy.spec.description || '-'}</span>
                )}
              </li>
              <li>
                <span className="item-descr-tit">策略类型</span>
                <span className="item-descr-txt">{strategy.spec.type}</span>
              </li>
              <li>
                <span className="item-descr-tit">创建时间</span>
                <span className="item-descr-txt">
                  {dateFormat(new Date(strategy.metadata.creationTimestamp), 'yyyy-MM-dd hh:mm:ss')}
                </span>
              </li>
            </ul>
          )}
          {editValue.editBasic && (
            <div>
              <Button type="primary" disabled={!enabled} onClick={_onSubmitBasic}>
                保存
              </Button>
              <Button style={{ marginLeft: '10px' }} onClick={_onCancelBasicEdit}>
                取消
              </Button>
            </div>
          )}
        </Card.Body>
      </Card>
      <Card>
        <Card.Body>
          <Tabs
            tabs={tabs}
            onActive={tab => {
              if (tab.id === 'users') {
                actions.associateActions.applyFilter({ search: sub });
              }
            }}
          >
            <TabPanel id="actions">
              {strategy && (
                <React.Fragment>
                  <Button
                    style={{ marginBottom: '10px' }}
                    disabled={!editorValue.ready}
                    onClick={_onStrategyGrammarEdit}
                  >
                    编辑
                  </Button>
                  {editorValue.readOnly && (
                    <CodeEditor
                      style={{ height: 400 }}
                      options={{ language: 'json', readOnly: true }}
                      onReady={_onReady}
                      onEdit={_onEdit}
                      loadingPlaceholder={<LoadingTip style={{ textAlign: 'center', padding: 20 }} />}
                    />
                  )}
                  {!editorValue.readOnly && (
                    <CodeEditor
                      style={{ height: 400 }}
                      options={{ language: 'json', readOnly: false }}
                      onReady={_onReady}
                      onEdit={_onEdit}
                      loadingPlaceholder={<LoadingTip style={{ textAlign: 'center', padding: 20 }} />}
                    />
                  )}
                  {!editorValue.readOnly && (
                    <div>
                      <Button type="primary" onClick={_onSubmitStrategyGrammar}>
                        保存
                      </Button>
                      <Button style={{ marginLeft: '10px' }} onClick={_onCancelStrategyGrammarEdit}>
                        取消
                      </Button>
                    </div>
                  )}
                </React.Fragment>
              )}
            </TabPanel>
            <TabPanel id="users">
              <Button type="primary" onClick={_setModalVisible} style={{ marginBottom: '10px' }}>
                关联用户
              </Button>
              <aside>
                <Modal caption={t('关联用户')} size="l" visible={modalVisible} onClose={_close}>
                  <Modal.Body>
                    <Transfer
                      leftCell={
                        <Transfer.Cell
                          title={t('关联用户')}
                          tip={t('支持按住 shift 键进行多选')}
                          header={
                            <SearchBox
                              value={userMsgsValue.inputValue}
                              onChange={value => {
                                setUserMsgsValue({ ...userMsgsValue, inputValue: value });
                                // 进行用户的搜索
                                actions.user.performSearch(value);
                              }}
                            />
                          }
                        >
                          <Table
                            records={userListRecords.filter(user => {
                              return (
                                (user.spec.name &&
                                  (user.spec.name.toLowerCase().includes(userMsgsValue.inputValue.toLowerCase()) ||
                                    user.spec.name.toLowerCase() !== 'admin')) ||
                                user.spec.displayName.toLowerCase().includes(userMsgsValue.inputValue.toLowerCase())
                              );
                            })}
                            rowDisabled={record => {
                              return userMsgsValue.targetKeys.indexOf(record.metadata.name) !== -1;
                            }}
                            recordKey={record => {
                              return record.metadata.name;
                            }}
                            columns={modalColumns}
                            addons={[
                              selectable({
                                value: userMsgsValue.newTargetKeys.concat(userMsgsValue.targetKeys),
                                onChange: keys => {
                                  const newKeys = [];
                                  keys.forEach(item => {
                                    if (userMsgsValue.targetKeys.indexOf(item) === -1) {
                                      newKeys.push(item);
                                    }
                                  });
                                  setUserMsgsValue({ ...userMsgsValue, newTargetKeys: newKeys });
                                },
                                rowSelect: true
                              })
                            ]}
                          />
                        </Transfer.Cell>
                      }
                      rightCell={
                        <Transfer.Cell title={t(`已选择 (${userMsgsValue.newTargetKeys.length}条)`)}>
                          <Table
                            records={userListRecords.filter(i => userMsgsValue.newTargetKeys.includes(i.metadata.name))}
                            recordKey="name"
                            columns={modalColumns}
                            addons={[
                              removeable({
                                onRemove: key =>
                                  setUserMsgsValue({
                                    ...userMsgsValue,
                                    newTargetKeys: userMsgsValue.newTargetKeys.filter(i => i !== key)
                                  })
                              })
                            ]}
                          />
                        </Transfer.Cell>
                      }
                    />
                  </Modal.Body>
                  <Modal.Footer>
                    <Button type="primary" onClick={_onSubmit}>
                      <Trans>确定</Trans>
                    </Button>
                    <Button type="weak" onClick={_close}>
                      <Trans>取消</Trans>
                    </Button>
                  </Modal.Footer>
                </Modal>
              </aside>
              <TablePanel
                columns={columns}
                model={associatedUsersList}
                action={actions.strategy}
                bodyClassName={'tc-15-table-panel tc-15-table-fixed-body'}
              />
            </TabPanel>
            <TabPanel id="groups">
              <GroupActionPanel />
              <GroupTablePanel />
            </TabPanel>
            <TabPanel id="roles">
              <RoleActionPanel />
              <RoleTablePanel />
            </TabPanel>
          </Tabs>
        </Card.Body>
      </Card>
    </React.Fragment>
  );
  // }

  function _renderOperationCell(name: string) {
    return (
      <LinkButton tipDirection="right" onClick={() => _removeAssociateUser(name)}>
        <Trans>解除关联</Trans>
      </LinkButton>
    );
  }

  function _onBasicEdit() {
    setEditValue({ ...editValue, editBasic: true });
  }

  async function _onSubmitBasic() {
    const { name, description } = basicParamsValue;
    await actions.strategy.updateStrategy.fetch({
      noCache: true,
      data: { ...strategy, name, description }
    });
    setEditValue({ ...editValue, editBasic: false });
  }

  function _onCancelBasicEdit() {
    setEditValue({ ...editValue, editBasic: false });
  }

  function _onStrategyGrammarEdit() {
    setEditorValue({ ...editorValue, readOnly: !editorValue.readOnly });
  }

  async function _onSubmitStrategyGrammar() {
    // strategy.statement = editValue.editorStatement;
    await actions.strategy.updateStrategy.fetch({
      noCache: true,
      data: { ...strategy, statement: editValue.editorStatement }
    });
    setEditValue({ ...editValue, editorStatement: strategy.statement });
    setEditorValue({ ...editorValue, readOnly: true });
  }

  function _onCancelStrategyGrammarEdit() {
    setEditorValue({ ...editorValue, readOnly: true });
  }

  function _onReady(instance) {
    editorRef = instance;
    editorRef.setValue(JSON.stringify(strategy.spec.statement, null, 2));
    editorRef.focus();
    setEditValue({ ...editValue, editorStatement: strategy.spec.statement });
    setEditorValue({ ...editorValue, ready: true });
  }

  function _onEdit(instance) {
    instance.getValue().then(value => {
      setEditValue({ ...editValue, editorStatement: JSON.parse(value) });
    });
  }

  async function _removeAssociateUser(name: string) {
    const yes = await Modal.confirm({
      message: t('确认删除当前所选用户？'),
      okText: t('删除'),
      cancelText: t('取消')
    });
    if (yes) {
      actions.associateActions.removeAssociatedUser.start([{ id: sub, userNames: [name] }]);
      actions.associateActions.removeAssociatedUser.perform();
    }
  }

  function _setModalVisible() {
    actions.user.applyFilter({ ifAll: true, isPolicyUser: true });
    setModalVisible(true);
  }
  function _close() {
    setModalVisible(false);
  }
  function _onSubmit() {
    actions.associateActions.associateUser.start([
      { id: strategy.metadata.name, userNames: userMsgsValue.newTargetKeys }
    ]);
    actions.associateActions.associateUser.perform();
    setModalVisible(false);
    setUserMsgsValue({
      ...userMsgsValue,
      targetKeys: userMsgsValue.targetKeys.concat(userMsgsValue.newTargetKeys),
      newTargetKeys: []
    });
  }
};
