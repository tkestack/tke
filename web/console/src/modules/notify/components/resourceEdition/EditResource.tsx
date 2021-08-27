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

import { Button } from '@tea/component/button';
import { Card } from '@tea/component/card';
import { Form } from '@tea/component/form';
import { Input } from '@tea/component/input';
import { OperationState, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { InputNumber, Switch, Bubble, Icon, Table, ExternalLink } from '@tencent/tea-component';

import { resourceConfig } from '../../../../../config';
import { CreateResource, getWorkflowError, TipInfo } from '../../../common';
import { router } from '../../router';
import { channelSchema } from '../../schema/channelSchema';
import { getState, onChange, schemaObjToJSON } from '../../schema/schemaUtil';
import { RootProps } from '../NotifyApp';

const rc = resourceConfig();

interface Props extends RootProps {
  instance?;
}
interface State {
  resource?;
  error?;
  searchValue?;
}

export class EditResource extends React.Component<Props, State> {
  constructor(props) {
    super(props);
    const schema = this.getSchema();
    const state = this.getState();
    state.resource = props.instance
      ? getState(schema, this, this.props.instance, this.props.instance)
      : getState(schema, this);

    // 转化webhook中headers的值为字符串
    const theSpec = state.resource.properties.spec;
    if (theSpec.pick === 'webhook') {
      const headers = theSpec.properties.webhook.properties.headers;
      if (headers && headers.value) {
        const headerArr = [];
        Object.keys(headers.value).forEach(key => {
          headerArr.push(key + ':' + headers.value[key]);
        });
        headers.value = headerArr.join(';');
      }
    }

    this.state = state;
  }

  componentDidMount() {
    this.props.actions.workflow.modifyResource.reset();
  }

  getState() {
    return {
      resource: undefined,
      error: undefined
    };
  }

  getSchema(): any {
    return channelSchema;
  }

  render() {
    const { modifyResourceFlow } = this.props;
    return (
      <Card>
        <Card.Body>
          {this.renderForm()}
          <div className="tea-form-operate">
            <Button type="primary" onClick={this.submit}>
              {t('保存')}
            </Button>
            <Button type="weak" onClick={this.back}>
              {t('取消')}
            </Button>
            <TipInfo
              className="tea-mt-1n"
              style={{ marginBottom: '0' }}
              isShow={modifyResourceFlow.operationState === OperationState.Done && modifyResourceFlow.results[0].error}
              type="error"
            >
              {getWorkflowError(modifyResourceFlow)}
            </TipInfo>
          </div>
        </Card.Body>
      </Card>
    );
  }

  renderForm() {
    return <React.Fragment />;
  }

  renderFields(obj) {
    return Object.keys(obj.properties).map(key => (
      <Form.Item
        label={key}
        key={key}
        required={obj.properties[key].required}
        {...(obj.properties[key].validator ? obj.properties[key].validator(obj.properties[key].value) : {})}
      >
        {obj.properties[key].type === 'boolean' ? (
          <Switch value={obj.properties[key].value} onChange={onChange(obj.properties[key])} />
        ) : obj.properties[key].type === 'number' ? (
          <InputNumber value={obj.properties[key].value} onChange={onChange(obj.properties[key])} />
        ) : (
          <Input
            type={key === 'password' ? 'password' : undefined}
            size="l"
            multiline={key === 'body' || key === 'headers'}
            placeholder={obj.properties[key].placeholder}
            value={obj.properties[key].value}
            onChange={onChange(obj.properties[key])}
            onBlur={e => onChange(obj.properties[key])(e.target.value)}
          />
        )}
        {obj.properties[key].bodyTip && (
          <Bubble content={<ExternalLink href="/tkestack/notify/bodyintro">body模板说明</ExternalLink>}>
            <Icon style={{ marginLeft: '5px' }} type="info" />
          </Bubble>
        )}
        {obj.properties[key].smsSignTip && (
          <Bubble
            content={
              <ExternalLink href="https://console.cloud.tencent.com/smsv2/csms-sign">
                查看腾讯云短信服务签名ID信息
              </ExternalLink>
            }
          >
            <Icon style={{ marginLeft: '5px' }} type="info" />
          </Bubble>
        )}
        {obj.properties[key].smsTemplateIDTip && (
          <Bubble
            content={
              <aside>
                <p>
                  通过
                  <ExternalLink href="https://console.cloud.tencent.com/smsv2">腾讯云短信服务</ExternalLink>
                  创建消息模板，腾讯云短信服务中自定义内容编号和TKE消息模板的自定义内容一一对应即可，例：
                </p>
                <p>
                  sms短信模板：注意{1},发生{2}
                </p>
                <p>
                  TKE Body模板：
                  {'注意{{.clusterID}},发生{{.alertName}}'}
                </p>
              </aside>
            }
          >
            <Icon style={{ marginLeft: '5px' }} type="info" />
          </Bubble>
        )}
      </Form.Item>
    ));
  }

  submit = () => {
    const { actions, route } = this.props;
    const urlParams = router.resolve(route);

    const resourceInfo = rc[urlParams.resourceName] || rc.channel;
    const mode = this.props.instance ? 'modify' : 'create';
    const json = schemaObjToJSON(this.state.resource);
    if (this.props.instance) {
      json.metadata.resourceVersion = this.props.instance.metadata.resourceVersion;
    }

    // 将headers字符串转换为对象
    if (json.spec && json.spec.webhook) {
      if (json.spec.webhook.headers) {
        const headersObj = {};
        json.spec.webhook.headers.split(';').forEach(headerStr => {
          if (headerStr) {
            const headerArr = headerStr.split(':');
            headersObj[headerArr[0]] = headerArr[1];
          }
        });
        json.spec.webhook.headers = headersObj;
      }
      if (urlParams.resourceName === 'template') {
        json.spec.text = Object.assign({}, json.spec.webhook);
        json.spec.webhook = undefined;
      }
    }

    const jsonData = JSON.stringify(json);

    const resource: CreateResource = {
      id: uuid(),
      resourceInfo,
      mode,
      namespace: (json.metadata && json.metadata.namespace) || route.queries['np'] || 'default',
      isSpecialNamespace: true,
      jsonData,
      clusterId: route.queries['clusterId'],
      resourceIns: mode === 'modify' ? this.props.instance.metadata.name : ''
    };

    actions.workflow.modifyResource.start([resource]);
    actions.workflow.modifyResource.perform();
  };
  back = () => {
    const { route } = this.props;
    const urlParams = router.resolve(route);
    router.navigate(
      Object.assign({}, urlParams, {
        mode: urlParams.mode === 'update' ? 'detail' : 'list'
      }),
      Object.assign({}, route.queries)
    );
  };
}
