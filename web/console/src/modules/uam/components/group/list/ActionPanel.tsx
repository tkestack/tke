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
import { connect } from 'react-redux';
import { Justify, Icon, Table, Button, SearchBox } from '@tea/component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators } from '@tencent/ff-redux';
import { router } from '../../../router';
import { allActions } from '../../../actions';
import { RootProps } from '../GroupPanel';

const mapDispatchToProps = (dispatch) =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch,
  });

@connect((state) => state, mapDispatchToProps)
export class ActionPanel extends React.Component<RootProps, {}> {
  render() {
    const { actions, route, groupList } = this.props;
    let urlParam = router.resolve(route);
    const { sub } = urlParam;

    return (
      <Table.ActionPanel>
        <Justify
          left={
            <Button
              type="primary"
              onClick={(e) => {
                e.preventDefault();
                router.navigate({ module: 'user', sub: 'group', action: 'create' }, route.queries);
              }}
            >
              {t('新建')}
            </Button>
          }
          right={
            <React.Fragment>
              <SearchBox
                value={groupList.query.keyword || ''}
                onChange={actions.group.list.changeKeyword}
                onSearch={actions.group.list.performSearch}
                onClear={() => {
                  actions.group.list.performSearch('');
                }}
                placeholder={t('请输入用户组名称')}
              />
            </React.Fragment>
          }
        />
      </Table.ActionPanel>
    );
  }
}
