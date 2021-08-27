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
import { t } from '@tencent/tea-app/lib/i18n';
import classNames from 'classnames';
import * as React from 'react';
import { CommonBar, FormItem } from '../../../../common/components';
import { FormPanel } from '@tencent/ff-component';
import { helmResourceList } from '../../../constants/Config';
import { RootProps } from '../../HelmApp';
export class BaseInfoPanel extends React.Component<RootProps, {}> {
  componentDidMount() {
    const { actions } = this.props;
    actions.create.fetchRegionList();
    // actions.create.getToken();
    // actions.create.selectTencenthubType(TencentHubType.Public);
  }
  onChangeName(name: string) {
    let {
      actions,
      helmCreation: { isValid }
    } = this.props;
    actions.create.inputName(name);
  }
  onSelectResource(resource: string) {
    this.props.actions.create.selectResource(resource);
  }

  render() {
    let {
      namespaceSelection,
      namespaceList,
      actions,
      listState: { cluster },
      helmCreation: { name, resourceSelection, isValid }
    } = this.props;
    let namespaceOptions = namespaceList.data.records.map((p, index) => ({
      text: p.name,
      value: p.name
    }));

    return (
      <div>
        <ul className="form-list">
          <FormItem label={t('应用名')}>
            <div
              className={classNames('form-unit', {
                'is-error': isValid.name !== ''
              })}
            >
              <input
                type="text"
                className="tc-15-input-text m"
                placeholder={t('请输入应用名称')}
                value={name}
                onChange={e => this.onChangeName((e.target.value + '').trim())}
              />
              <p className="form-input-help">
                {t('最长63个字符，只能包含小写字母、数字及分隔符("-")，且必须以小写字母开头，数字或小写字母结尾')}
              </p>
            </div>
          </FormItem>

          <FormItem label={t('运行集群')}>{cluster.selection ? cluster.selection.metadata.name : '-'}</FormItem>

          <FormItem label={t('命名空间')}>
            <div
              className={classNames('form-unit', {
                'is-error': namespaceSelection === ''
              })}
            >
              <FormPanel.Select
                label={'namespace'}
                options={namespaceOptions}
                value={namespaceSelection}
                onChange={value => actions.namespace.selectNamespace(value)}
              />
            </div>
          </FormItem>
        </ul>
        <hr className="hr-mod" />
        <ul className="form-list" style={{ marginBottom: 16 }}>
          <FormItem label={t('来源')}>
            <div className="form-unit">
              <CommonBar
                list={helmResourceList}
                value={resourceSelection}
                onSelect={item => {
                  this.onSelectResource(item.value + '');
                }}
                isNeedPureText={false}
              />

              {/* {resource === HelmResource.Helm &&
              this.renderHelmPanel()}
            {resource === HelmResource.Other &&
              this.renderOtherPanel()} */}
            </div>
          </FormItem>
        </ul>
      </div>
    );
  }
}
