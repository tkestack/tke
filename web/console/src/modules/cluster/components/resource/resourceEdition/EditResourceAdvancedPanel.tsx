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
import { connect } from 'react-redux';

import { FormPanel } from '@tencent/ff-component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Select } from '@tencent/tea-component';

import { FormItem } from '../../../../common';
import { allActions } from '../../../actions';
import { FloatingIPReleasePolicy, WorkloadNetworkType, WorkloadNetworkTypeEnum } from '../../../constants/Config';
import { RootProps } from '../../ClusterApp';
import { EditResourceAnnotations } from './EditResourceAnnotations';
import { EditResourceImagePullSecretsPanel } from './EditResourceImagePullSecretsPanel';
import { EditResourceNodeAffinityPanel } from './EditResourceNodeAffinityPanel';

interface EditResourceAdvancedPanelProps extends RootProps {
  /** 是否展示高级设置 */
  isOpenAdvanced: boolean;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class EditResourceAdvancedPanel extends React.Component<EditResourceAdvancedPanelProps, {}> {
  render() {
    const { isOpenAdvanced, subRoot, actions } = this.props,
      { workloadEdit } = subRoot,
      { networkType, floatingIPReleasePolicy } = workloadEdit;

    // let isShowPort = networkType !== 'Overlay';

    return isOpenAdvanced ? (
      <React.Fragment>
        <EditResourceImagePullSecretsPanel />
        <EditResourceNodeAffinityPanel />
        <EditResourceAnnotations />
        <FormItem label={t('网络模式')} isShow={false}>
          <Select
            size="m"
            options={WorkloadNetworkType}
            value={networkType}
            onChange={value => {
              actions.editWorkload.selectNetworkType(value);
            }}
          />
        </FormItem>
        <FormItem isShow={networkType === WorkloadNetworkTypeEnum.FloatingIP} label={t('IP回收策略')}>
          <FormPanel.Select
            size="m"
            options={FloatingIPReleasePolicy}
            value={floatingIPReleasePolicy}
            onChange={value => {
              actions.editWorkload.selectFloatingIPReleasePolicy(value);
            }}
          ></FormPanel.Select>
        </FormItem>
        {/* <FormItem label={t('端口')} isShow={isShowPort}>
        </FormItem> */}
      </React.Fragment>
    ) : (
      <noscript />
    );
  }
}
