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

import { BaseReactProps } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { SelectList, SelectListProps } from '../select';

export interface VpcNetworkProps extends BaseReactProps {
  /**VPC列表 */
  vpc: SelectListProps;

  /**isShowCIDR */
  isShowCIDR?: boolean;

  /**mode当前的模式 */
  mode?: string;
}

export class VpcNetwork extends React.Component<VpcNetworkProps, {}> {
  render() {
    let { vpc, isShowCIDR, mode } = this.props,
      totalIPNum = 0,
      availableIPNum = 0,
      cidr = '';

    if (vpc.value) {
      if (isShowCIDR) {
        vpc.recordData.data.records.forEach(v => {
          if (v.unVpcId === vpc.value) {
            cidr = v.cidrBlock;
          }
        });
      }
    }

    return (
      <div>
        <SelectList {...vpc} name={t('集群网络')} className="tc-15-select m" style={{ display: 'inline-block' }} />
        {cidr && mode === 'create' && (
          <span className="inline-help-text text-weak" style={{ marginLeft: '5px' }}>
            CIDR: {cidr}
          </span>
        )}
      </div>
    );
  }
}
