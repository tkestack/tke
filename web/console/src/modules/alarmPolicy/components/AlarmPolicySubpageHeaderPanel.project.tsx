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

import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Icon, Justify } from '@tencent/tea-component';

import { router } from '../router';
import { RootProps } from './AlarmPolicyApp';
export class AlarmPolicySubpageHeaderPanel extends React.Component<RootProps, {}> {
  goBack() {
    let { route } = this.props;
    // history.back();
    router.navigate(
      {},
      { clusterId: route.queries['clusterId'], projectName: route.queries['projectName'], np: route.queries['np'] }
    );
  }

  componentDidMount() {
    // 根据 queries 来判断在update或者 detail当中显示信息
    let { cluster, actions, route } = this.props;
    if (cluster.list.data.recordCount === 0) {
      actions.projectNamespace.initProjectList();
    } else {
      actions.alarmPolicy.initAlarmPolicyData();
    }
  }

  render() {
    let { route } = this.props;
    let urlParams = router.resolve(route);
    let title = '';
    switch (urlParams['sub']) {
      case 'create':
        title = t('新建策略');
        break;
      case 'update':
        title = t('更新策略');
        break;
      case 'copy':
        title = t('复制策略');
        break;
      default:
        title = '';
    }

    return (
      <div className="manage-area-title secondary-title">
        <a href="javascript:;" className="back-link" onClick={this.goBack.bind(this)}>
          <i className="btn-back-icon" />
          {t('返回')}
        </a>
        <span className="line-icon">|</span>
        <h2>{title}</h2>
      </div>
    );
  }
}
