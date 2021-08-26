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

import { FormPanel } from '@tencent/ff-component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Button, Icon, Text } from '@tencent/tea-component';

import { isEmpty, LinkButton, FormItem } from '../../../../common';
import { allActions } from '../../../actions';
import { ContainerEnv, Resource } from '@src/modules/cluster/models';
import { RootProps } from '../../ClusterApp';

interface ContainerEnvItemProps extends RootProps {
  /** 容器的id */
  cKey?: string;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class EditResourceContainerEnvItem extends React.Component<ContainerEnvItemProps, any> {
  render() {
    let { actions, cKey, subRoot } = this.props,
      { workloadEdit } = subRoot,
      { containers, configEdit } = workloadEdit,
      { configList, secretList } = configEdit;

    let container = containers.find(c => c.id === cKey);
    let envItems: ContainerEnv.ItemWithId[] = container && container.envItems ? container.envItems : [];
    const configOptions = this._reduceResourceListOptions(configList.data.records);
    const secretOptions = this._reduceResourceListOptions(secretList.data.records);
    let canNotAdd =
      envItems.filter(envItem => {
        let disabled = false;
        if (envItem.name === '') {
          disabled = true;
        }
        if (!disabled) {
          if (envItem.type === ContainerEnv.EnvTypeEnum.ConfigMapRef) {
            disabled = envItem.configMapName === '' || envItem.configMapDataKey === '';
          } else if (envItem.type === ContainerEnv.EnvTypeEnum.SecretKeyRef) {
            disabled = envItem.secretName === '' || envItem.secretDataKey === '';
          }
        }
        return disabled;
      }).length > 0;

    return (
      <FormItem label={t('环境变量')} tips={t('设置容器中的变量')}>
        {envItems.map(envItem => {
          let envItemId = envItem.id + '';
          let envType = envItem.type;
          let isField = envType === ContainerEnv.EnvTypeEnum.FieldRef;
          let isResourceField = envType === ContainerEnv.EnvTypeEnum.ResourceFieldRef;
          let isSecret = envType === ContainerEnv.EnvTypeEnum.SecretKeyRef;
          let isConfig = envType === ContainerEnv.EnvTypeEnum.ConfigMapRef;

          let configMapKeyOptions = [];
          let secretKeyOptions = [];
          if (isSecret) {
            secretKeyOptions = this._reduceConfigKeyOptions(secretList.data.records, envItem.secretName);
          }
          if (isConfig) {
            configMapKeyOptions = this._reduceConfigKeyOptions(configList.data.records, envItem.configMapName);
          }

          return (
            <div key={envItemId}>
              <FormPanel.Select
                size="s"
                className="tea-mr-1n"
                options={ContainerEnv.EnvTypeOptions}
                value={envType}
                onChange={value => {
                  actions.editWorkload.updateEnvItem({ type: value as ContainerEnv.EnvTypeEnum }, cKey, envItemId);
                }}
              />

              <FormPanel.Input
                placeholder={t('变量名称')}
                size="s"
                validator={envItem.v_name}
                errorTipsStyle="Bubble"
                className="tea-mr-1n"
                value={envItem.name}
                onChange={value => actions.editWorkload.updateEnvItem({ name: value }, cKey, envItemId)}
                onBlur={e => actions.validate.workload.validateAllEnvName(container)}
                maxLength={63}
              />

              {envType === ContainerEnv.EnvTypeEnum.UserDefined && (
                <textarea
                  placeholder={t('变量值')}
                  className="tc-15-input-text m"
                  style={{
                    maxWidth: '304px',
                    minHeight: '30px',
                    minWidth: '304px',
                    overflowY: 'visible'
                  }}
                  value={envItem.value}
                  onChange={e => actions.editWorkload.updateEnvItem({ value: e.target.value }, cKey, envItemId)}
                />
              )}

              {(isField || isResourceField) && (
                <>
                  <FormPanel.Input
                    size="s"
                    disabled={true}
                    value={isField ? 'fieldPath' : 'resource'}
                    className="tea-mr-1n"
                  />
                  <FormPanel.Select
                    options={isField ? ContainerEnv.FieldRefOptions : ContainerEnv.ResourceFieldRefOptions}
                    value={isField ? envItem.fieldName : envItem.resourceFieldName}
                    onChange={value =>
                      actions.editWorkload.updateEnvItem(
                        isField
                          ? { fieldName: value as ContainerEnv.FieldKeyNameEnum }
                          : { resourceFieldName: value as ContainerEnv.ResourceFieldKeyNameEnum },
                        cKey,
                        envItemId
                      )
                    }
                  />
                </>
              )}

              {(isSecret || isConfig) && (
                <>
                  <FormPanel.Select
                    className="tea-mr-1n"
                    size="s"
                    options={isSecret ? secretOptions : configOptions}
                    value={isSecret ? envItem.secretName : envItem.configMapName}
                    onChange={value => {
                      actions.editWorkload.updateEnvItem(
                        isSecret ? { secretName: value } : { configMapName: value },
                        cKey,
                        envItemId
                      );
                    }}
                    placeholder={t('请选择资源')}
                    validator={isSecret ? envItem.v_secretName : envItem.v_configMapName}
                    errorTipsStyle="Bubble"
                  />
                  <FormPanel.Select
                    options={isSecret ? secretKeyOptions : configMapKeyOptions}
                    value={isSecret ? envItem.secretDataKey : envItem.configMapDataKey}
                    onChange={value => {
                      actions.editWorkload.updateEnvItem(
                        isSecret ? { secretDataKey: value } : { configMapDataKey: value },
                        cKey,
                        envItemId
                      );
                    }}
                    placeholder={t('请选择Key')}
                    validator={isSecret ? envItem.v_secretDataKey : envItem.v_configMapDataKey}
                    errorTipsStyle="Bubble"
                  />
                </>
              )}

              <LinkButton
                style={{ position: 'absolute', right: '25px', top: '5px' }}
                onClick={() => actions.editWorkload.deleteEnvItem(cKey, envItemId)}
                className="tea-ml-1n"
              >
                <Icon type="close" />
              </LinkButton>
            </div>
          );
        })}

        <div className={envItems.length ? 'tea-mt-1n' : ''}>
          <LinkButton
            className="tea-mr-1n"
            disabled={canNotAdd}
            tipDirection="right"
            errorTip={t('请先完成待编辑项')}
            onClick={() => {
              actions.editWorkload.addEnvItem(cKey);
            }}
          >
            {t('新增变量')}
          </LinkButton>
        </div>

        <p className="text-label">{t('变量名只能包含大小写字母、数字及下划线，并且不能以数字开头')}</p>
      </FormItem>
    );
  }

  /** 渲染resourceList的options选项 */
  private _reduceResourceListOptions(list: Resource[]) {
    let configListOptions = list.map(item => ({
      value: item.metadata.name,
      text: item.metadata.name
    }));

    return configListOptions;
  }

  /** 渲染configmap/secret的options选项 */
  private _reduceConfigKeyOptions(list: Resource[], resourceName: string) {
    let options = [];
    if (resourceName) {
      let finder = list.find(item => item.metadata.name === resourceName);
      let dataKeys = [];
      if (finder && finder.data) {
        dataKeys = Object.keys(finder.data);
      }
      options = dataKeys.map(item => ({
        value: item,
        text: item
      }));
    }
    return options;
  }
}
