import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { bindActionCreators } from '@tencent/ff-redux';
import { allActions } from '../../actions';
import { router } from '../../router';
import { UserCreate } from './create/UserCreate';
import { UserList } from './list/UserList';
import { UserDetail } from './detail/UserDetail';

export const UserPanel = (props) => {
  const state = useSelector((state) => state);
  // const dispatch = useDispatch();
  // const { actions } = bindActionCreators({ actions: allActions }, dispatch);
  const { route } = state;
  const { sub, action } = router.resolve(route);

  let content;
  if (!action && (!sub || sub === 'normal')) {
    content = <UserList />;
  } else if (action === 'create') {
    content = <UserCreate />;
  } else if (action === 'detail') {
    content = <UserDetail />;
  }
  return content;
};
