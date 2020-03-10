import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';

import { emptyTips, LinkButton } from '@src/modules/common';
import { Button, Modal, Card, Input, Form, TableColumn, Tabs, TabPanel } from '@tea/component';
import { TablePanel } from '@tencent/ff-component';
import { bindActionCreators, insertCSS } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { dateFormat } from '../../../../../helpers/dateUtil';
import { allActions } from '../../actions';
import { STRATEGY_TYPE, VALIDATE_EMAIL_RULE, VALIDATE_PHONE_RULE } from '../../constants/Config';
import { Strategy, User } from '../../models';
import { router } from '../../router';

import { RoleActionPanel } from './detail/RoleActionPanel';
import { RoleTablePanel } from './detail/RoleTablePanel';
import { GroupActionPanel } from './detail/GroupActionPanel';
import { GroupTablePanel } from './detail/GroupTablePanel';
const { useState, useEffect, useRef } = React;
const _isEqual = require('lodash/isEqual');

insertCSS(
  'UserDetailsPanel',
  `
    .item-descr-list .is-error {
      color: #e1504a;
      border-color: #e1504a;
    }
`
);

export const UserDetailsPanel = () => {
  const state = useSelector(state => state);
  const dispatch = useDispatch();
  const { actions } = bindActionCreators({ actions: allActions }, dispatch);

  const { route, userList, getUser, updateUser, userStrategyList } = state;
  const getUserData = getUser.data[0];
  const updateUserData = updateUser.data[0];
  const { sub } = router.resolve(route);

  const [basicParamsValue, setBasicParamsValue] = useState({ displayName: '', email: '', phoneNumber: '' });
  const [editValue, setEditValue] = useState({ editBasic: false });
  const [user, setUser] = useState(undefined);

  useEffect(() => {
    // 请求用户详情
    actions.user.getUser.fetch({
      noCache: true,
      data: { name: sub }
    });

    // 进行用户绑定的策略的拉取
    actions.user.strategy.applyFilter({ specificName: sub });
  }, []);

  useEffect(() => {
    // 初始化用户详情
    if (getUserData && getUserData.target.metadata.name === sub) {
      const showUser: User = getUserData.target;
      const { displayName = '', email = '', phoneNumber = '' } = showUser.spec;
      setUser(showUser);
      setBasicParamsValue({ displayName, email, phoneNumber });
    }
  }, [getUserData, sub]);

  useEffect(() => {
    // 更新user后修改state数据: 有个坑 —— 上边初始化用户详情后，下边user会变更，如果updateUserData存储有以往的旧数据，就会在里边setUser旧数据
    if (updateUserData && updateUserData.success && !_isEqual(user, updateUserData.target)) {
      const showUser = updateUserData.target;
      setUser(showUser);
    }
  }, [updateUserData]);

  const { displayName, phoneNumber, email } = basicParamsValue;
  const isNameError = displayName.length <= 0 || displayName.length > 255;
  const { displayName: pDisplayName = '', phoneNumber: pPhoneNumber = '', email: pEmail = '' } = user ? user.spec : {};

  // 都满足，确定才可用
  const enabled =
    (pDisplayName !== displayName || pPhoneNumber !== phoneNumber || pEmail !== email) &&
    !isNameError &&
    (!phoneNumber || VALIDATE_PHONE_RULE.pattern.test(phoneNumber)) &&
    (!email || VALIDATE_EMAIL_RULE.pattern.test(email));

  const columns: TableColumn<Strategy>[] = [
    {
      key: 'name',
      header: t('策略名'),
      width: '20%',
      render: x => x.spec.displayName
    },
    {
      key: 'category',
      header: t('类型'),
      width: '20%',
      render: x => x.spec.category
    },
    {
      key: 'desp',
      header: t('描述'),
      width: '40%',
      render: x => x.spec.description
    }
  ];

  const tabs = [
    { id: 'policies', label: '已关联策略' },
    { id: 'groups', label: '已关联用户组' },
    { id: 'roles', label: '已关联角色' },
  ];

  return (
    <React.Fragment>
      <Card>
        <Card.Body
          title={t('基本信息')}
          subtitle={
            <Button type="link" onClick={_onBasicEdit}>
              编辑
            </Button>
          }
        >
          {user && (
            <ul className="item-descr-list">
              <li>
                <span className="item-descr-tit">用户账号</span>
                <span className="item-descr-txt">{user.spec.username}</span>
              </li>
              <li>
                <span className="item-descr-tit">用户名称</span>
                {editValue.editBasic ? (
                  <React.Fragment>
                    <Input
                      value={displayName}
                      className={isNameError && 'is-error'}
                      onChange={value => {
                        setBasicParamsValue({ ...basicParamsValue, displayName: value });
                      }}
                    />
                    {isNameError ? <p className="is-error">输入不能为空且需要小于256个字符</p> : ''}
                  </React.Fragment>
                ) : (
                  <span className="item-descr-txt">{user.spec.displayName}</span>
                )}
              </li>
              <li>
                <span className="item-descr-tit">手机号</span>
                {editValue.editBasic ? (
                  <React.Fragment>
                    <Input
                      value={phoneNumber}
                      onChange={value => {
                        setBasicParamsValue({ ...basicParamsValue, phoneNumber: value });
                      }}
                    />
                    {VALIDATE_PHONE_RULE.pattern.test(phoneNumber) || !phoneNumber ? (
                      ''
                    ) : (
                      <p className="is-error">{VALIDATE_PHONE_RULE.message}</p>
                    )}
                  </React.Fragment>
                ) : (
                  <span className="item-descr-txt">{user.spec.phoneNumber || '-'}</span>
                )}
              </li>
              <li>
                <span className="item-descr-tit">邮箱</span>
                {editValue.editBasic ? (
                  <React.Fragment>
                    <Input
                      value={email}
                      onChange={value => {
                        setBasicParamsValue({ ...basicParamsValue, email: value });
                      }}
                    />
                    {VALIDATE_EMAIL_RULE.pattern.test(email) || !email ? (
                      ''
                    ) : (
                      <p className="is-error">{VALIDATE_EMAIL_RULE.message}</p>
                    )}
                  </React.Fragment>
                ) : (
                  <span className="item-descr-txt">{user.spec.email || '-'}</span>
                )}
              </li>
              <li>
                <span className="item-descr-tit">创建时间</span>
                <span className="item-descr-txt">
                  {dateFormat(new Date(user.metadata.creationTimestamp), 'yyyy-MM-dd hh:mm:ss')}
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
          <Tabs tabs={tabs}>
            <TabPanel id="policies">
              <TablePanel
                isNeedCard={true}
                columns={columns}
                model={userStrategyList}
                action={actions.user.strategy}
                emptyTips={emptyTips}
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

  function _onBasicEdit() {
    setEditValue({ editBasic: true });
  }

  async function _onSubmitBasic() {
    const { displayName, phoneNumber, email } = basicParamsValue;

    await actions.user.updateUser.fetch({
      noCache: true,
      data: {
        user: {
          metadata: {
            name: user.metadata.name,
            resourceVersion: user.metadata.resourceVersion
          },
          spec: {
            username: user.spec.username,
            displayName,
            phoneNumber,
            email
          }
        }
      }
    });
    setEditValue({ editBasic: false });
  }

  function _onCancelBasicEdit() {
    setEditValue({ editBasic: false });
  }
};
