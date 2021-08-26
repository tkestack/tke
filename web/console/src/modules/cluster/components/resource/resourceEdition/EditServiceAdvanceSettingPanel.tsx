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
import { RadioGroup, Radio } from '@tea/component';
import { FormItem, InputField } from '../../../../common/components';
import { SessionAffinity, ExternalTrafficPolicy } from '../../../constants/Config';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Validation } from '../../../../common/models';
interface EditServiceAdvanceSettingProps {
  isShow?: boolean;

  validatesessionAffinityTimeout?: () => void;

  chooseChoosesessionAffinityMode?: (mode: string) => void;

  inputsessionAffinityTimeout?: (seconds: number) => void;

  chooseExternalTrafficPolicyMode?: (mode: string) => void;

  /**访问方式 */
  communicationType: string;

  /** externalTrafficPolicy */
  externalTrafficPolicy?: string;

  /**会话保持 */
  sessionAffinity?: string;

  sessionAffinityTimeout?: number;

  v_sessionAffinityTimeout?: Validation;
}

const ExternalTrafficPolicyItems = [
  {
    value: 'Cluster',
    label: 'Cluster'
  },
  {
    value: 'Local',
    label: 'Local'
  }
];

const sessionAffinityItems = [
  {
    value: 'ClientIP',
    label: 'ClientIP'
  },
  {
    value: 'None',
    label: 'None'
  }
];
export class EditServiceAdvanceSettingPanel extends React.Component<EditServiceAdvanceSettingProps, {}> {
  /** 构造数据盘选择列表 */
  renderRadioList(items: any) {
    let list: JSX.Element[] = [];
    items.forEach(item => {
      list.push(
        <Radio key={item.value} name={item.value}>
          {item.label}
        </Radio>
      );
    });
    return list;
  }
  render() {
    let {
      isShow,
      validatesessionAffinityTimeout,
      chooseChoosesessionAffinityMode,
      inputsessionAffinityTimeout,
      chooseExternalTrafficPolicyMode,
      externalTrafficPolicy,
      sessionAffinity,
      sessionAffinityTimeout,
      v_sessionAffinityTimeout,
      communicationType
    } = this.props;
    return isShow ? (
      <ul className="form-list fixed-layout jiqun">
        {communicationType !== 'ClusterIP' && (
          <FormItem label="ExternalTrafficPolicy">
            <div className="form-unit">
              <Radio.Group
                value={externalTrafficPolicy}
                onChange={value => chooseExternalTrafficPolicyMode(value)}
                // style={{ fontSize: '12px', display: 'inline-block' }}
              >
                {this.renderRadioList(ExternalTrafficPolicyItems)}
              </Radio.Group>
              <p className="form-input-help text-weak">
                {externalTrafficPolicy === ExternalTrafficPolicy.Cluster
                  ? t('默认均衡转发到工作负载的所有Pod')
                  : t(
                      '能够保留来源IP，并可以保证公网、VPC内网访问（LoadBalancer）和主机端口访问（NodePort）模式下流量仅在本节点转发。Local转发使部分没有业务Pod存在的节点健康检查失败，可能存在流量不均衡的转发的风险。'
                    )}
              </p>
            </div>
          </FormItem>
        )}
        <FormItem label="Session Affinity">
          <div className="form-unit">
            <Radio.Group
              value={sessionAffinity}
              onChange={value => chooseChoosesessionAffinityMode(value)}
              // style={{ fontSize: '12px', display: 'inline-block' }}
            >
              {this.renderRadioList(sessionAffinityItems)}
            </Radio.Group>
            <p className="form-input-help text-weak">
              {sessionAffinity === SessionAffinity.ClientIP && t('基于来源IP做会话保持。')}
            </p>
          </div>
        </FormItem>
        <FormItem isShow={sessionAffinity === SessionAffinity.ClientIP} label={t('最大会话保持时间')}>
          <div className="form-unit">
            <InputField
              type="text"
              value={sessionAffinityTimeout}
              validator={v_sessionAffinityTimeout}
              onChange={value => inputsessionAffinityTimeout(value)}
              onBlur={() => validatesessionAffinityTimeout()}
              tipMode="popup"
              tip={
                communicationType !== 'ClusterIP' && communicationType !== 'NodePort'
                  ? t(
                      '会话保持时间范围为30-3600，若您的访问方式是公网或VPC内网访问（LoadBalancer）模式，设置成CLB监听器的会话保持时间一致。'
                    )
                  : t('会话保持时间范围为0-86400')
              }
            />
          </div>
        </FormItem>
      </ul>
    ) : (
      <noscript />
    );
  }
}
