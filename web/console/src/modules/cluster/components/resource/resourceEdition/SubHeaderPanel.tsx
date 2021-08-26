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

import { allActions } from '../../../actions';
import { router } from '../../../router';
import { RootProps } from '../../ClusterApp';

interface SubHeaderPanelProps extends RootProps {
  /** 标题名称 */
  headTitle?: string;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class SubHeaderPanel extends React.Component<SubHeaderPanelProps, {}> {
  render() {
    let { headTitle } = this.props;

    return (
      <div className="manage-area-title secondary-title">
        <a href="javascript:;" className="back-link" onClick={this.goBack.bind(this)}>
          <i className="btn-back-icon" />
          {t('返回')}
        </a>
        <h2>{headTitle}</h2>
        {/* <div className="manage-area-title-right"></div> */}
      </div>
    );
  }

  private goBack() {
    let { route } = this.props,
      urlParam = router.resolve(route);
    // 回到列表处
    let routeQueries = JSON.parse(JSON.stringify(Object.assign({}, route.queries, { resourceIns: undefined })));
    let newUrlParmas = JSON.parse(JSON.stringify(Object.assign({}, urlParam, { mode: 'list', tab: undefined })));
    router.navigate(newUrlParmas, routeQueries);
  }
}
