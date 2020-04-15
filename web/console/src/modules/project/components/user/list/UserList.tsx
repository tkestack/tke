import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Layout } from '@tencent/tea-component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators } from '@tencent/ff-redux/libs/qcloud-lib';
import { UserActionPanel } from './UserActionPanel';
import { UserTablePanel } from './UserTablePanel';
import { allActions } from '../../../actions';
import { router } from '../../../router';
const { useState, useEffect } = React;

export const UserList = (props) => {
  const state = useSelector((state) => state);
  const dispatch = useDispatch();
  const { actions } = bindActionCreators({ actions: allActions }, dispatch);
  const { route } = state;
  console.log('member route:', route);
  useEffect(() => {
    actions.user.poll(route.queries);
  }, []);

  return (
    <>
      <UserActionPanel />
      <UserTablePanel />
    </>
  );
};
