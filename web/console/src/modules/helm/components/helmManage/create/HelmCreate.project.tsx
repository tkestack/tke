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

import { Button } from '@tea/component/button';
import { t } from '@tencent/tea-app/lib/i18n';
import * as React from 'react';
import { FormLayout, MainBodyLayout } from '../../../../common/layouts';
import { HelmResource, OtherType } from '../../../constants/Config';
import { router } from '../../../router';
import { RootProps } from '../../HelmApp';
import { BaseInfoPanel } from './BaseInfoPanel';
import { KeyValuesPanel } from './KeyValuesPanel';
import { OtherChartPanel } from './OtherChartPanel';
import { TencentHubChartPanel } from './TencentHubChartPanel';

export class HelmCreate extends React.Component<RootProps, {}> {
  componentDidMount() {
    // this.props.actions.create.clear();
    this.props.actions.create.inputKeyValue([]);
  }
  componentWillUnmount() {
    this.props.actions.create.clear();
    this.props.actions.create.inputKeyValue([]);
    // 去除错误信息
  }
  goBack() {
    let { actions, route } = this.props,
      urlParams = router.resolve(route);
    router.navigate({}, route.queries);
  }
  onOk() {
    const {
      actions,
      helmCreation: { name }
    } = this.props;
    actions.create.validAll();
    let canSave = this.isCanSave();
    if (canSave) {
      this.props.actions.create.createHelm();
    }
  }
  onCancel() {
    this.goBack();
  }

  isCanSave() {
    let {
      namespaceSelection,
      helmCreation: { isValid, resourceSelection, otherTypeSelection }
    } = this.props;
    let canSave = true;
    if (isValid.name !== '' || namespaceSelection === '') {
      canSave = false;
    } else {
      if (resourceSelection === HelmResource.Other) {
        if (isValid.otherChartUrl !== '') {
          canSave = false;
        } else {
          if (otherTypeSelection === OtherType.Private) {
            if (isValid.otherUserName !== '' || isValid.otherPassword !== '') {
              canSave = false;
            }
          }
        }
      }
    }
    return canSave;
  }
  render() {
    let {
      helmCreation: { resourceSelection }
    } = this.props;

    let canSave = this.isCanSave();
    return (
      <div>
        <div className="manage-area-title secondary-title">
          <a
            href="javascript:void(0)"
            className="back-link"
            onClick={() => {
              this.goBack();
            }}
          >
            <i className="btn-back-icon" />
            <span>{t('返回')}</span>
          </a>
          <span className="line-icon"> |</span>
          <h2>{t('新建 Helm 应用')}</h2>
        </div>
        <MainBodyLayout>
          <div className="manage-area-main-inner">
            {/* <TipInfo>
              <span style={{ verticalAlign: 'middle' }}>
                {t(
                  '创建Helm应用，若应用中包含了公网CLB类型的Services或Ingress，将按照腾讯云CLB对应价格收费。若应用中包含PV/PVC/StorageClass，其创建的存储资源将按对应的产品价格收费。'
                )}
              </span>
            </TipInfo> */}
            <FormLayout>
              <div className="param-box">
                <div className="param-bd">
                  <BaseInfoPanel {...this.props} />
                  {resourceSelection === HelmResource.TencentHub && <TencentHubChartPanel {...this.props} />}
                  {resourceSelection === HelmResource.Other && <OtherChartPanel {...this.props} />}
                  <KeyValuesPanel
                    onChangeKeyValue={this.props.actions.create.inputKeyValue}
                    kvs={this.props.helmCreation.kvs}
                  />
                </div>
                <div className="param-ft">
                  <Button
                    className="mr10"
                    title={t('完成')}
                    disabled={!canSave}
                    onClick={() => {
                      this.onOk();
                    }}
                    type="primary"
                  >
                    {t('完成')}
                  </Button>
                  <Button
                    title={t('取消')}
                    onClick={() => {
                      this.onCancel();
                    }}
                  >
                    {t('取消')}
                  </Button>
                </div>
              </div>
            </FormLayout>
          </div>
        </MainBodyLayout>
      </div>
    );
  }
}
