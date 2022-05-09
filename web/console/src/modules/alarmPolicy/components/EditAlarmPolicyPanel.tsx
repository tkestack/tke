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
import classNames from 'classnames';
import * as React from 'react';

import { Bubble, Button, Col, ExternalLink, Form, Row, Select } from '@tea/component';
import { isSuccessWorkflow, OperationState, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { FormItem, InputField, SelectList, TipInfo } from '../../common/components';
import { FormLayout, MainBodyLayout } from '../../common/layouts';
import { getWorkflowError } from '../../common/utils';
import { router as notifyRouter } from '../../notify/router';
import { validatorActions } from '../actions/validatorActions';
import {
  AlarmPolicyMetricsContinuePeriod,
  AlarmPolicyMetricsEvaluatorType,
  AlarmPolicyMetricsEvaluatorValue,
  AlarmPolicyMetricsStatisticsPeriod,
  AlarmPolicyType,
  MetricNameMap
} from '../constants/Config';
// import { EditAlarmPolicyReceiverTunnel } from './EditAlarmPolicyReceiverTunnel';
import { MetricsObjectEdition } from '../models/AlarmPolicy';
import { router } from '../router';
import { RootProps } from './AlarmPolicyApp';
import { EditAlarmPolicyObject } from './EditAlarmPolicyObject';
import { EditAlarmPolicyReceiverGroup } from './EditAlarmPolicyReceiverGroup';
import { EditAlarmPolicyVM } from './EditAlarmPolicyVM';

export class EditAlarmPolicyPanel extends React.Component<RootProps, {}> {
  componentWillUnmount() {
    const { actions } = this.props;
    actions.alarmPolicy.clearAlarmPolicyEdit();
    actions.workflow.editAlarmPolicy.reset();
  }
  _renderAlarmMetrics(alarmMetrics: MetricsObjectEdition[]) {
    const { actions, alarmPolicyEdition, channel, template } = this.props;
    const content: JSX.Element[] = [];
    alarmMetrics.forEach((item, index) => {
      content.push(
        <div key={index}>
          <label className="form-ctrl-label" style={{ display: 'inline-block', width: 195 }}>
            <input
              type="checkbox"
              className="tc-15-checkbox"
              checked={item.enable}
              onClick={() => actions.alarmPolicy.inputAlarmMetrics(item.id + '', { enable: !item.enable })}
            />
            {MetricNameMap[item.metricName] || item.metricName}
            {item.tip && (
              <Bubble placement="right" content={item.tip || null}>
                <i className="plaint-icon" />
              </Bubble>
            )}
          </label>

          {item.type === 'boolean' ? (
            <div className="form-unit" style={{ display: 'inline-block', fontSize: 12 }}>
              <input
                type="text"
                className="tc-15-input-text s "
                value="="
                readOnly={true}
                style={{ width: 70, marginTop: 6, marginRight: 10 }}
              />
            </div>
          ) : (
            <SelectList
              key={uuid()}
              value={item.evaluatorType + ''}
              recordList={AlarmPolicyMetricsEvaluatorType}
              valueField="value"
              textFields={['text']}
              textFormat={`\${text}`}
              className="tc-15-select"
              style={{ display: 'inline-block', maxWidth: 80, minWidth: 70, marginRight: 5 }}
              onSelect={value => actions.alarmPolicy.inputAlarmMetrics(item.id + '', { evaluatorType: value })}
              isUnshiftDefaultItem={false}
            />
          )}
          {item.type === 'boolean' ? (
            //事件型告警，True=> 1, lt
            //False=>0, gt
            <SelectList
              key={uuid()}
              value={item.evaluatorValue + ''}
              recordList={AlarmPolicyMetricsEvaluatorValue}
              valueField="value"
              textFields={['text']}
              textFormat={`\${text}`}
              className="tc-15-select"
              outerStyle={{ display: 'inline-block' }}
              style={{ display: 'inline-block', maxWidth: 80, minWidth: 80 }}
              onSelect={value => {
                actions.alarmPolicy.inputAlarmMetrics(item.id + '', {
                  evaluatorValue: value
                  // evaluatorType: value === '1' ? 'lt' : 'gt'
                });
              }}
              isUnshiftDefaultItem={false}
            />
          ) : (
            <InputField
              type="text"
              className="tc-15-input-text s "
              popDirection="right"
              placeholder={t('请输入数值')}
              value={item.evaluatorValue}
              style={{ width: 80, marginTop: 6 }}
              validator={item.v_evaluatorValue}
              tipMode="popup"
              onChange={value => actions.alarmPolicy.inputAlarmMetrics(item.id + '', { evaluatorValue: value })}
              onBlur={() => item.enable && actions.validator.validateEvaluatorValue(item.id + '')}
            />
          )}
          <div className="form-unit" style={{ display: 'inline-block', fontSize: 12 }}>
            <input
              type="text"
              className="tc-15-input-text s "
              value={item.unit}
              readOnly={true}
              style={{ width: 55, marginTop: 6, marginRight: 10 }}
            />
          </div>
          <SelectList
            key={uuid()}
            value={item.continuePeriod + ''}
            recordList={AlarmPolicyMetricsContinuePeriod}
            valueField="value"
            textFields={['value']}
            // textFormat={t('持续') + `\${value}` + t('个周期')}
            textFormat={value => {
              return t('持续{{count}}个周期', { count: value });
            }}
            className="tc-15-select m"
            style={{ display: 'inline-block' }}
            onSelect={value => actions.alarmPolicy.inputAlarmMetrics(item.id + '', { continuePeriod: +value })}
            isUnshiftDefaultItem={false}
          />
        </div>
      );
    });
    return <div className="form-unit unit-group">{content}</div>;
  }

  render() {
    const {
      actions,
      alarmPolicyEdition,
      cluster,
      route,
      alarmPolicyCreateWorkflow,
      alarmPolicyUpdateWorkflow,
      channel,
      template,
      namespaceList
    } = this.props;
    const urlParams = router.resolve(route);
    const workflow = urlParams['sub'] === 'update' ? alarmPolicyUpdateWorkflow : alarmPolicyCreateWorkflow;
    const failed = workflow.operationState === OperationState.Done && !isSuccessWorkflow(workflow);
    return (
      <MainBodyLayout className="secondary-main">
        <FormLayout>
          <div className="param-box add" style={{ paddingBottom: '50px' }}>
            <div className="param-bd">
              <ul className="form-list" style={{ paddingBottom: '0' }}>
                {/* <FormItem label={t('地域')}>
                  <p className="form-input-help">{regionSelection.name && regionSelection.name}</p>
                </FormItem> */}
                <FormItem label={t('集群')}>
                  <p className="form-input-help">
                    {cluster.selection &&
                      cluster.selection.metadata &&
                      cluster.selection.metadata.name &&
                      `${cluster.selection.metadata.name}(${cluster.selection.metadata.name})`}
                  </p>
                </FormItem>
                <FormItem label={t('告警策略名称')}>
                  {urlParams['sub'] === 'update' ? (
                    <p className="form-input-help">{alarmPolicyEdition.alarmPolicyName}</p>
                  ) : (
                    <InputField
                      type="text"
                      popDirection="right"
                      className="tc-15-input-text m"
                      placeholder={t('请输入告警策略名称')}
                      value={alarmPolicyEdition.alarmPolicyName}
                      tip={t('最长60个字符')}
                      validator={alarmPolicyEdition.v_alarmPolicyName}
                      onChange={value => actions.alarmPolicy.inputAlarmPolicyName(value)}
                      onBlur={actions.validator.validateAlarmPolicyName}
                      tipMode="popup"
                    />
                  )}
                </FormItem>
                {/* <FormItem label={t('备注')}>
                  <InputField
                    type="textarea"
                    popDirection="right"
                    className="tc-15-input-text m"
                    placeholder={t('请输入策略备注')}
                    tip={t('最长100个字符')}
                    validator={alarmPolicyEdition.v_alarmPolicyDescription}
                    value={alarmPolicyEdition.alarmPolicyDescription}
                    onChange={value => actions.alarmPolicy.inputAlarmPolicyDescription(value)}
                    onBlur={actions.validator.validateDescription}
                    tipMode="popup"
                  />
                </FormItem> */}
                {
                  /// #if tke
                  <FormItem label={t('策略类型')}>
                    <SelectList
                      value={alarmPolicyEdition.alarmPolicyType}
                      recordList={AlarmPolicyType}
                      valueField="value"
                      textField="text"
                      textFields={['text']}
                      textFormat={`\${text}`}
                      className="tc-15-select s"
                      isUnshiftDefaultItem={false}
                      style={{ marginRight: '5px' }}
                      validator={alarmPolicyEdition.v_alarmPolicyType}
                      onSelect={value => {
                        actions.alarmPolicy.inputAlarmPolicyType(value);
                        actions.validator.validateAlarmPolicyType();
                      }}
                    />
                  </FormItem>
                  /// #endif
                }

                {alarmPolicyEdition.alarmPolicyType === 'virtualMachine' ? (
                  <FormItem label={t('告警对象')} isNeedFormInput={false}>
                    <EditAlarmPolicyVM
                      clusterId={cluster?.selection?.metadata?.name}
                      namespaceList={namespaceList?.data?.records ?? []}
                      type={alarmPolicyEdition?.alarmObjectsType}
                      setType={actions?.alarmPolicy?.inputAlarmPolicyObjectsType}
                      namespaceSelection={alarmPolicyEdition?.alarmObjectNamespace}
                      setNamespaceSelection={actions?.alarmPolicy?.selectsWorkLoadNamespace}
                      vmSelections={alarmPolicyEdition?.alarmObjects}
                      setVmSelections={actions?.alarmPolicy?.inputAlarmPolicyObjects}
                    />
                  </FormItem>
                ) : (
                  <EditAlarmPolicyObject {...this.props} />
                )}

                <FormItem label={t('统计周期')}>
                  <SelectList
                    value={alarmPolicyEdition.statisticsPeriod + ''}
                    recordList={AlarmPolicyMetricsStatisticsPeriod}
                    valueField="value"
                    textFields={['value']}
                    textFormat={value => {
                      return t('{{count}}分钟', { count: value });
                    }}
                    className="tc-15-select s"
                    style={{ display: 'inline-block', marginRight: 5 }}
                    onSelect={value => actions.alarmPolicy.inputAlarmPolicyStatisticsPeriod(+value)}
                    isUnshiftDefaultItem={false}
                  />
                </FormItem>
                <FormItem label={t('指标')} isShow={alarmPolicyEdition.alarmMetrics.length !== 0}>
                  {this._renderAlarmMetrics(alarmPolicyEdition.alarmMetrics)}
                  <div
                    className={classNames('', {
                      'is-error': alarmPolicyEdition.v_alarmMetrics.status === 2
                    })}
                  >
                    <p className="form-input-help">
                      {alarmPolicyEdition.v_alarmMetrics.status === 2 && alarmPolicyEdition.v_alarmMetrics.message}
                    </p>
                  </div>
                </FormItem>
                <EditAlarmPolicyReceiverGroup {...this.props} />

                <FormItem
                  label={t('通知方式')}
                  tips={
                    <ExternalLink
                      href={
                        window.location.pathname.indexOf('tkestack-project') !== -1
                          ? '/tkestack-project/notify/create/channel'
                          : '/tkestack/notify/create/channel'
                      }
                    >
                      {t('新建通知渠道')}
                    </ExternalLink>
                  }
                >
                  {alarmPolicyEdition.notifyWays.map((notifyWay, index) => (
                    <Row key={notifyWay.id}>
                      <Col span={6}>
                        <Select
                          value={notifyWay.channel}
                          options={channel.list.data.records.map(c => ({
                            value: c.metadata.name,
                            text: `${c.metadata.name}(${c.spec.displayName})`
                          }))}
                          placeholder={t('请选择通知渠道')}
                          onChange={value => {
                            actions.alarmPolicy.inputAlarmNotifyWay(notifyWay.id, {
                              channel: value,
                              template: undefined
                            });
                          }}
                        />
                      </Col>
                      <Col span={6}>
                        <Select
                          placeholder={t('请选择消息模版')}
                          value={notifyWay.template}
                          options={template.list.data.records
                            .filter(t => t.metadata.namespace === notifyWay.channel)
                            .map(c => ({
                              value: c.metadata.name,
                              text: `${c.metadata.name}(${c.spec.displayName})`
                            }))}
                          onChange={value => {
                            actions.alarmPolicy.inputAlarmNotifyWay(notifyWay.id, {
                              template: value
                            });
                          }}
                        />
                      </Col>
                      <Col>
                        <Button
                          type="icon"
                          disabled={alarmPolicyEdition.notifyWays.length === 1}
                          icon="close"
                          onClick={() => actions.alarmPolicy.deleteAlarmNotifyWay(notifyWay.id)}
                        />
                      </Col>
                    </Row>
                  ))}

                  <div className="is-error">
                    <p className="form-input-help" style={{ fontSize: '12px', marginTop: '5px' }}>
                      {alarmPolicyEdition.v_notifyWay.message}
                    </p>
                  </div>

                  <Row>
                    <Col>
                      <Button type="link" onClick={() => actions.alarmPolicy.addAlarmNotifyWay()}>
                        {t('添加通知方式')}
                      </Button>
                    </Col>
                  </Row>
                </FormItem>

                {/* <EditAlarmPolicyReceiverTunnel {...this.props} /> */}
                <li className="pure-text-row" style={{ position: 'absolute' }}>
                  <Button
                    type="primary"
                    className="mr10"
                    disabled={workflow.operationState === OperationState.Performing}
                    onClick={this._handleSubmit.bind(this)}
                  >
                    {failed ? t('重试') : t('提交')}
                  </Button>
                  <Button
                    onClick={e =>
                      router.navigate({}, { rid: route.queries['rid'], clusterId: route.queries['clusterId'] })
                    }
                  >
                    {t('取消')}
                  </Button>
                  {failed && (
                    <TipInfo style={{ display: 'inline-block', marginBottom: 10 }} type="error" className="error">
                      {getWorkflowError(workflow)}
                    </TipInfo>
                  )}
                </li>
              </ul>
            </div>
          </div>
        </FormLayout>
      </MainBodyLayout>
    );
  }
  /** 处理提交请求 */
  private _handleSubmit() {
    const { actions, alarmPolicyEdition, route, regionSelection, cluster, receiverGroup } = this.props;

    actions.validator.validateAlarmPolicyEdition();

    if (validatorActions._validateAlarmPolicyEdition(alarmPolicyEdition, receiverGroup)) {
      actions.workflow.editAlarmPolicy.start([alarmPolicyEdition], {
        regionId: +regionSelection.value,
        clusterId: cluster.selection ? cluster.selection.metadata.name : ''
      });
      actions.workflow.editAlarmPolicy.perform();
    }
  }
}
