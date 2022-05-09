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

import { Checkbox, Radio } from '@tea/component';
import { FetchState } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { FormItem, SelectList } from '../../common/components';
import { AlarmObjectsType, workloadTypeList } from '../constants/Config';
import { RootProps } from './AlarmPolicyApp';

export class EditAlarmPolicyObject extends React.Component<RootProps, {}> {
  renderPodList() {
    const Tip = content => {
      return (
        <div className="colony" style={{ fontSize: '12px' }}>
          <span>{content}</span>
        </div>
      );
    };
    const { workloadList, alarmPolicyEdition, actions } = this.props;
    if (workloadList.fetchState === FetchState.Fetching) {
      return Tip(t('加载中'));
    } else if (workloadList.fetchState === FetchState.Failed) {
      return Tip(t('加载失败'));
    } else if (workloadList.data.recordCount === 0) {
      return Tip(t('该命名空间下无workload'));
    } else {
      // 根据 PodList 初始化 checkedList
      const workloadOptions = workloadList?.data?.records?.map(workload => workload?.metadata?.name) ?? [];

      const checkedList = alarmPolicyEdition?.alarmObjects?.filter(item => workloadOptions.includes(item));

      return (
        <Checkbox.Group
          onChange={items => {
            actions.alarmPolicy.inputAlarmPolicyObjects(items);
          }}
          value={checkedList}
          layout="column"
        >
          {workloadOptions.map(name => <Checkbox name={name}>{name}</Checkbox>) ?? null}
        </Checkbox.Group>
      );
    }
  }

  renderRadioList(type) {
    const { alarmPolicyEdition, actions, namespaceList, addons } = this.props;
    if (type === 'cluster' || type === '') {
      return <noscript />;
    }
    const finalWorkloadTypeList = workloadTypeList.slice();
    const finalWorkloadNsList = namespaceList.data.records.map(ns => ({
      value: ns.name,
      label: ns.displayName
    }));
    if (addons['TappController']) {
      finalWorkloadTypeList.push({
        value: 'TApp',
        label: 'TApp'
      });
    }

    if (type === 'pod' && alarmPolicyEdition.alarmObjectsType === 'all') {
      finalWorkloadTypeList.unshift({
        value: 'ALL',
        label: 'ALL'
      });
      finalWorkloadNsList.unshift({
        value: 'ALL',
        label: 'ALL'
      });
    }
    const radioList: JSX.Element[] = [];
    AlarmObjectsType[type].forEach((item, index) => {
      radioList.push(
        <div className="form-unit unit-group new-strategy-alarm-object">
          <div className="alarm-select">
            <Radio key={index} name={item.value} disabled={item.value === 'k8sLabel'}>
              {item.text}
              <span className="text-label">{item.tip}</span>
            </Radio>
          </div>
          <div className="alarm-write">
            {alarmPolicyEdition.alarmPolicyType === 'pod' &&
              ((item.value === 'part' && alarmPolicyEdition.alarmObjectsType === 'part') ||
                (item.value === 'all' && alarmPolicyEdition.alarmObjectsType === 'all')) && (
                <ul className="form-list fixed-layout">
                  <FormItem label="Namespace">
                    <SelectList
                      value={alarmPolicyEdition.alarmObjectNamespace + ''}
                      recordList={finalWorkloadNsList}
                      valueField="value"
                      textFields={['label']}
                      textFormat={`\${label}`}
                      className="tc-15-select m"
                      style={{ marginRight: '5px' }}
                      onSelect={value => {
                        actions.namespace.selectNamespace(value);
                      }}
                      isUnshiftDefaultItem={false}
                    />
                  </FormItem>
                  <FormItem label="WorkloadType" isNeedFormInput={false}>
                    <SelectList
                      value={alarmPolicyEdition.alarmObjectWorkloadType + ''}
                      recordList={finalWorkloadTypeList}
                      valueField="value"
                      textFields={['label']}
                      textFormat={`\${label}`}
                      className="tc-15-select m"
                      style={{ marginRight: '5px' }}
                      onSelect={value => actions.alarmPolicy.inputAlarmObjectWorkloadType(value)}
                      isUnshiftDefaultItem={false}
                    />
                    {alarmPolicyEdition.alarmObjectsType === 'part' && (
                      <>
                        <div
                          className="param-box"
                          style={{
                            backgroundColor: '#fff',
                            padding: '5px 10px',
                            marginTop: '5px',
                            marginBottom: '10px',
                            width: '250px'
                          }}
                        >
                          <div className="param-bd">
                            <ul>
                              <div
                                className="tc-g-u-1-3"
                                style={{ width: '100%', maxHeight: '180px', overflowY: 'auto' }}
                              >
                                <div className="colony-list">{this.renderPodList()}</div>
                              </div>
                            </ul>
                          </div>
                        </div>
                        <div className="is-error">
                          <p className="form-input-help" style={{ fontSize: '12px' }}>
                            {alarmPolicyEdition.v_alarmObjects.status === 2 &&
                              alarmPolicyEdition.v_alarmObjects.message}
                          </p>
                        </div>
                      </>
                    )}
                  </FormItem>
                </ul>
              )}
          </div>
        </div>
      );
    });
    return radioList;
  }
  render() {
    const { actions, alarmPolicyEdition } = this.props;

    const isShow = alarmPolicyEdition.alarmPolicyType !== 'cluster' && alarmPolicyEdition.alarmPolicyType !== '';
    return (
      <FormItem isShow={isShow} label={t('告警对象')} isNeedFormInput={false}>
        <Radio.Group
          value={alarmPolicyEdition.alarmObjectsType}
          onChange={value => actions.alarmPolicy.inputAlarmPolicyObjectsType(value)}
          // className="form-unit"
        >
          {this.renderRadioList(alarmPolicyEdition.alarmPolicyType)}
        </Radio.Group>
      </FormItem>
    );
  }
}
