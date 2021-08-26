/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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

import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Justify, Icon, Table, Button, SearchBox } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';

import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { router } from '../../../router';
import { allActions } from '../../../actions';
import { PlatformTypeEnum } from '@src/modules/project/constants/Config';

export const UserActionPanel = props => {
  const state = useSelector(state => state);
  const dispatch = useDispatch();
  const { actions } = bindActionCreators({ actions: allActions }, dispatch);
  const { route, userList, projectDetail, platformType, userManagedProjects } = state;
  let enableOp =
    platformType === PlatformTypeEnum.Manager ||
    (platformType === PlatformTypeEnum.Business &&
      userManagedProjects.list.data.records.find(
        item => item.name === (projectDetail ? projectDetail.metadata.name : null)
      ));
  return (
    <Table.ActionPanel>
      <Justify
        left={
          enableOp ? (
            <Button
              type="primary"
              onClick={e => {
                e.preventDefault();
                router.navigate({ sub: 'detail', tab: 'member', action: 'create' }, route.queries);
              }}
            >
              {t('新建')}
            </Button>
          ) : null
        }
        right={
          <React.Fragment>
            <SearchBox
              value={userList.query.keyword || ''}
              onChange={actions.user.changeKeyword}
              onSearch={actions.user.performSearch}
              onClear={() => {
                actions.user.performSearch('');
              }}
              placeholder={t('请输入用户名称')}
            />
          </React.Fragment>
        }
      />
    </Table.ActionPanel>
  );
};
