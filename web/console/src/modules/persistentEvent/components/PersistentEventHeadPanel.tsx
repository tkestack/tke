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

import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Justify } from '@tencent/tea-component';

import { allActions } from '../actions';
import { RootProps } from './PersistentEventApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class PersistentEventHeadPanel extends React.Component<RootProps, {}> {
  componentDidMount() {
    const { actions, region, route } = this.props;

    // 这里对从创建界面返回之后，判断当前的状态
    let isNeedFetchRegion = region.list.data.recordCount ? false : true;
    isNeedFetchRegion && actions.region.applyFilter({});
    !isNeedFetchRegion && actions.cluster.applyFilter({ regionId: +route.queries['rid'] });
  }

  render() {
    return (
      <React.Fragment>
        <Justify
          left={
            <React.Fragment>
              <h2 style={{ float: 'left' }}>{t('事件持久化')}</h2>
            </React.Fragment>
          }
        />
      </React.Fragment>
    );
  }
}
