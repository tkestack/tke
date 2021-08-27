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

import { FormItem, LinkButton } from '@src/modules/common';
import { FormPanel } from '@tencent/ff-component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Icon, Text } from '@tencent/tea-component';

import { allActions } from '../../../actions';
import { router } from '../../../router';
import { RootProps } from '../../ClusterApp';

interface WorkloadPodAdvancePanelState {
  /** 是否需要展示高级设置的内容 */
  isOpenAdvanced: boolean;
}
const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class EditResourceImagePullSecretsPanel extends React.Component<RootProps, any> {
  render() {
    let { subRoot, actions, route } = this.props,
      urlParams = router.resolve(route),
      { workloadEdit } = subRoot,
      { configEdit, imagePullSecrets } = workloadEdit,
      { secretList } = configEdit;

    /**
     * 渲染imagePullSecret的列表
     * pre: secret的类型为 kubernetes.io/dockercfg
     */
    let finalSecretList = secretList.data.records.filter(
      item => item.type === 'kubernetes.io/dockercfg' || item.type === 'kubernetes.io/dockerconfigjson'
    );
    let secretListOptions = finalSecretList.map(item => ({
      value: item.metadata.name,
      text: item.metadata.name
    }));

    return (
      <FormItem label="imagePullSecrets">
        {finalSecretList.length ? (
          imagePullSecrets.map((item, index) => {
            let sId = item.id + '';

            return (
              <div key={index}>
                <FormPanel.Select
                  key={index}
                  value={item.secretName}
                  className="tea-mb-1n"
                  size="m"
                  onChange={value => {
                    actions.editWorkload.secret.updateImagePullSecret({ secretName: value }, sId);
                    actions.validate.workload.validateImagePullSecret(value, sId);
                  }}
                  options={secretListOptions}
                  placeholder={t('请选择Secret')}
                />
                <LinkButton onClick={() => actions.editWorkload.secret.deleteImagePullSecret(sId)}>
                  <Icon type="close" />
                </LinkButton>
              </div>
            );
          })
        ) : (
          <Text theme="label" parent="p">
            <Trans>
              当前命名空间下无可用Secret，前往配置项管理进行
              <a
                href="javascript:;"
                onClick={e => {
                  this._handleClickForUpdateSecret(urlParams);
                }}
              >
                新建Secret
              </a>
            </Trans>
          </Text>
        )}
        {finalSecretList.length > 0 && (
          <LinkButton
            tipDirection="right"
            disabled={imagePullSecrets.length >= finalSecretList.length}
            errorTip={t('无更多Secret可添加')}
            onClick={() => {
              actions.editWorkload.secret.addImagePullSecret();
            }}
          >
            {t('添加')}
          </LinkButton>
        )}
      </FormItem>
    );
  }

  /** 处理没有secret的时候，跳转ns进行下发 */
  private _handleClickForUpdateSecret(urlParams) {
    let { actions, route } = this.props;
    // 不需要拉取namespace，用于在 namespace列表上，决定需不需要路由跳转
    actions.resource.initResourceInfoAndFetchData(true, 'secret');
    router.navigate(
      Object.assign({}, urlParams, { mode: 'list', resourceName: 'secret', type: 'config' }),
      route.queries
    );
  }
}
