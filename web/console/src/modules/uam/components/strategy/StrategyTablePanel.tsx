import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';

import {
    Button, Icon, Modal, SearchBox, Table, TableColumn, Text, Tooltip, Transfer
} from '@tea/component';
import { removeable, selectable } from '@tea/component/table/addons';
import { TablePanel } from '@tencent/ff-component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { emptyTips, LinkButton } from '../../../common/components';
import { allActions } from '../../actions';
import { router } from '../../router';
import { Strategy, GroupFilter } from '../../models';
import { GroupAssociateWorkflowDialog } from './associate/GroupAssociateWorkflowDialog';
const { useState, useEffect } = React;
const _isEqual = require('lodash/isEqual');

export const StrategyTablePanel = () => {
  const state = useSelector(state => state);
  const dispatch = useDispatch();
  const { actions } = bindActionCreators({ actions: allActions }, dispatch);

  const { strategyList, userList, associatedUsersList } = state;
  const associatedUsersListRecords = associatedUsersList.list.data.records.map(item => item.metadata.name);
  const userListRecords = userList.list.data.records;

  const [modalVisible, setModalVisible] = useState(false);
  const [userMsgsValue, setUserMsgsValue] = useState({
    inputValue: '',
    targetKeys: associatedUsersListRecords,
    newTargetKeys: []
  });
  const [currentStrategy, setCurrentStrategy] = useState({ id: undefined });

  useEffect(() => {
    // 关联用户
    if (!_isEqual(associatedUsersListRecords, userMsgsValue.targetKeys)) {
      setUserMsgsValue({ ...userMsgsValue, targetKeys: associatedUsersListRecords });
    }
  }, [associatedUsersListRecords, userMsgsValue]);

  useEffect(() => {
    return () => {
      actions.user.performSearch('');
    };
  }, []);

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
  const columns: TableColumn<Strategy>[] = [
    {
      key: 'name',
      header: t('策略名'),
      render: (item, text, index) => (
        <Text parent="div" overflow>
          <a
            href="javascript:;"
            onClick={e => {
              router.navigate({ module: 'strategy', sub: `${item.metadata.name}` });
            }}
          >
            {item.spec.displayName}
          </a>
          {item.status['phase'] === 'Terminating' && <Icon type="loading" />}
        </Text>
      )
    },
    {
      key: 'description',
      header: t('描述'),
      render: item => <Text parent="div">{item.spec.description || '-'}</Text>
    },
    {
      key: 'service',
      header: t('服务类型'),
      render: item => <Text parent="div">{item.spec.category || '-'}</Text>
    },
    { key: 'operation', header: t('操作'), render: user => _renderOperationCell(user) }
  ];

  return (
    <React.Fragment>
      <TablePanel
        recordKey={(record) => {
          return record.metadata.name;
        }}
        columns={columns}
        model={strategyList}
        action={actions.strategy}
        rowDisabled={record => record.status['phase'] === 'Terminating'}
        emptyTips={emptyTips}
        isNeedPagination={true}
        bodyClassName={'tc-15-table-panel tc-15-table-fixed-body'}
      />
      <GroupAssociateWorkflowDialog onPostCancel={() => {
        //取消按钮时，清理编辑状态
        actions.group.associate.clearGroupAssociation();
      }}
      />
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
                        actions.user.performSearch(value);
                      }}
                    />
                  }
                >
                  <Table
                    records={
                      userListRecords &&
                      userListRecords.filter(
                        user =>
                          (user.spec.name &&
                            (user.spec.name.toLowerCase().includes(userMsgsValue.inputValue.toLowerCase()) ||
                              user.spec.name.toLowerCase() !== 'admin')) ||
                          user.spec.displayName.toLowerCase().includes(userMsgsValue.inputValue.toLowerCase())
                      )
                    }
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
                    records={
                      userListRecords && userListRecords.filter(i => userMsgsValue.newTargetKeys.includes(i.name))
                    }
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
    </React.Fragment>
  );

  /** 渲染操作按钮 */
  function _renderOperationCell(strategy: Strategy) {
    return (
      <React.Fragment>
        <LinkButton
          tipDirection="right"
          onClick={() => _setModalVisible(strategy)}
          disabled={strategy.status['phase'] === 'Terminating'}
        >
          <Trans>关联用户</Trans>
        </LinkButton>
        <LinkButton
          tipDirection="right"
          disabled={strategy.status['phase'] === 'Terminating'}
          onClick={(e) => {
            /** 设置用户组关联场景 */
            let filter: GroupFilter = {
              resource: 'policy',
              resourceID: strategy.metadata.name,
              /** 关联/解关联回调函数 */
              callback: () => {
                actions.strategy.fetch();
              }
            };
            actions.group.associate.setupGroupFilter(filter);
            /** 拉取关联用户组列表，拉取后自动更新groupAssociation */
            actions.group.associate.groupAssociatedList.applyFilter(filter);
            /** 拉取用户组列表 */
            actions.group.associate.groupList.performSearch('');
            /** 开始关联用户组工作流 */
            actions.group.associate.associateGroupWorkflow.start();
          }}
        >
          <Trans>关联用户组</Trans>
        </LinkButton>
        {strategy.type !== 1 && <LinkButton onClick={() => _removeCategory(strategy)}>删除</LinkButton>}
      </React.Fragment>
    );
  }
  function _setModalVisible(strategy: Strategy) {
    actions.user.applyFilter({ ifAll: true, isPolicyUser: true });
    actions.associateActions.applyFilter({ search: strategy.metadata.name + '' });
    setModalVisible(true);
    setCurrentStrategy({
      id: strategy.metadata.name
    });
  }
  function _close() {
    setModalVisible(false);
    setUserMsgsValue({
      ...userMsgsValue,
      newTargetKeys: []
    });
  }
  function _onSubmit() {
    actions.associateActions.associateUser.start([{ id: currentStrategy.id, userNames: userMsgsValue.newTargetKeys }]);
    actions.associateActions.associateUser.perform();
    setModalVisible(false);
    setUserMsgsValue({
      ...userMsgsValue,
      targetKeys: userMsgsValue.targetKeys.concat(userMsgsValue.newTargetKeys),
      newTargetKeys: []
    });
  }
  async function _removeCategory(strategy: Strategy) {
    const yes = await Modal.confirm({
      message: t('确认删除当前所选策略？'),
      okText: t('删除'),
      cancelText: t('取消')
    });
    if (yes) {
      actions.strategy.removeStrategy.start([strategy.metadata.name]);
      actions.strategy.removeStrategy.perform();
    }
  }
};
