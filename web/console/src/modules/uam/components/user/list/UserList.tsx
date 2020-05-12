import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Layout } from '@tencent/tea-component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { UserActionPanel } from './UserActionPanel';
import { UserTablePanel } from './UserTablePanel';
import { bindActionCreators } from '@tencent/ff-redux/libs/qcloud-lib';
import { allActions } from '@src/modules/uam/actions';
import { router } from '@src/modules/uam/router';
const { useState, useEffect } = React;

export const UserList = (props) => {
  // const state = useSelector((state) => state);
  const dispatch = useDispatch();
  const { actions } = bindActionCreators({ actions: allActions }, dispatch);

  useEffect(() => {
    actions.user.poll();
  }, []);

  return (
    <>
      <UserActionPanel />
      <UserTablePanel />
    </>
  );
};
