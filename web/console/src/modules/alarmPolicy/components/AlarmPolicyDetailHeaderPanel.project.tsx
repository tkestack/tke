/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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
import { RootProps } from './AlarmPolicyApp';
import { router } from '../router';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Icon, Justify } from '@tencent/tea-component';
export class AlarmPolicyDetailHeaderPanel extends React.Component<RootProps, {}> {
  goBack() {
    let { route } = this.props;
    // history.back();
    router.navigate(
      {},
      { clusterId: route.queries['clusterId'], projectName: route.queries['projectName'], np: route.queries['np'] }
    );
  }

  componentDidMount() {
    let { regionList, cluster, actions } = this.props;
    if (cluster.list.data.recordCount === 0) {
      actions.projectNamespace.initProjectList();
    } else {
      actions.alarmPolicy.initAlarmPolicyData();
    }
  }

  render() {
    let title = t('告警策略详情');

    return (
      <React.Fragment>
        <Justify
          left={
            <React.Fragment>
              <a href="javascript:;" className="back-link" onClick={this.goBack.bind(this)}>
                <Icon type="btnback" />
                {t('返回')}
              </a>
              <h2>{title}</h2>
            </React.Fragment>
          }
        />
      </React.Fragment>
    );
  }
}
