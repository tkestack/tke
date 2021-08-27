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
import { t } from '@tencent/tea-app/lib/i18n';
import { Input } from '@tencent/tea-component';

import { allActions } from '../actions';
import { RootProps } from './AddonApp';
import { InputField } from '@src/modules/common';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class EditPersistentEventPanel extends React.Component<RootProps, any> {
  render() {
    let { editAddon, route, actions } = this.props,
      { peEdit } = editAddon,
      { esAddress, indexName, v_esAddress, v_indexName, esUsername, esPassword } = peEdit;

    let { rid } = route.queries;

    return (
      <React.Fragment>
        <FormPanel.Item validator={v_esAddress} label={t('Elasticsearch地址')} errorTipsStyle="Bubble">
          <Input
            value={esAddress}
            onChange={value => {
              actions.editAddon.pe.inputEsAddress(value);
            }}
            placeholder="eg: http://190.0.0.1:9200"
            onBlur={actions.validator.validateEsAddress}
          />
        </FormPanel.Item>

        <FormPanel.Item
          validator={v_indexName}
          label={t('索引')}
          errorTipsStyle="Bubble"
          message={t('最长60个字符，只能包含小写字母、数字及分隔符("-"、"_"、"+")，且必须以小写字母开头')}
        >
          <Input
            value={indexName}
            onChange={value => {
              actions.editAddon.pe.inputIndexName(value);
            }}
            placeholder="eg: fluentd"
            onBlur={actions.validator.validateIndexName}
          />
        </FormPanel.Item>

        <FormPanel.Item label={t('用户名')}>
          <Input
            style={{
              width: '300px'
            }}
            placeholder="仅需要用户验证的 Elasticsearch 需要填入用户名"
            value={esUsername}
            onChange={value => {
              actions.editAddon.pe.inputEsUsername(value);
            }}
          />
        </FormPanel.Item>

        <FormPanel.Item label={t('密码')}>
          <Input
            type="password"
            style={{
              width: '300px'
            }}
            placeholder="仅需要用户验证的 Elasticsearch 需要填入密码"
            value={esPassword}
            onChange={value => {
              actions.editAddon.pe.inputEsPassword(value);
            }}
          />
        </FormPanel.Item>

      </React.Fragment>
    );
  }
}
