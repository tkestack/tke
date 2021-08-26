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

import { Card, Text } from '@tea/component';
import { FormPanel } from '@tencent/ff-component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';

import { dateFormatter } from '../../../../helpers';
import { allActions } from '../actions';
import { MetadataItem } from '../models';
import { RootProps } from './LogStashApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(
  state => state,
  mapDispatchToProps
)
export class LogStashDetailPanel extends React.Component<RootProps, {}> {
  componentWillUnmount() {
    let { actions } = this.props;
    actions.editLogStash.clearLogStashEdit();
  }

  formatLabels(labels: MetadataItem[]) {
    let labs: JSX.Element[] = [];
    for (let label of labels) {
      labs.push(
        <a className="tag" style={{ marginRight: '5px', cursor: 'default' }} href="javascript:;">
          {label.metadataKey + ' : ' + label.metadataValue}
        </a>
      );
    }
    return <div className="tag-cont-sub">{labs}</div>;
  }

  _renderKafkaInfo() {
    let { logStashEdit } = this.props,
      { addressIP, addressPort, topic } = logStashEdit;
    return (
      <React.Fragment>
        <FormPanel.Item text label={t('类型')}>
          Kafka
        </FormPanel.Item>
        <FormPanel.Item text label={t('访问地址IP')}>
          {addressIP}
        </FormPanel.Item>
        <FormPanel.Item text label={t('访问地址端口')}>
          {addressPort}
        </FormPanel.Item>
        <FormPanel.Item text label={t('主题（Topic）')}>
          {topic}
        </FormPanel.Item>
      </React.Fragment>
    );
  }

  _renderContainerServices() {
    let { logStashEdit } = this.props,
      { isSelectedAllNamespace } = logStashEdit;

    return (
      <React.Fragment>
        <FormPanel.Item text label={t('日志类型')}>
          {t('指定容器日志')}
        </FormPanel.Item>
        {isSelectedAllNamespace === 'selectAll' ? (
          <FormPanel.Item text label={t('日志源')}>
            {t('所有容器')}
          </FormPanel.Item>
        ) : (
          this._renderSpecificContainerLogList()
        )}
      </React.Fragment>
    );
  }

  _renderSpecificContainerLogList() {
    let { logSelection } = this.props;

    let { containerLogs } = this.props.logStashEdit;
    return (
      <FormPanel.Item text label={t('日志源')}>
        {containerLogs.map((contaienrlog, index) => {
          if (contaienrlog.collectorWay === 'workload') {
            return Object.keys(contaienrlog.workloadSelection).map((workload, index1) => {
              return contaienrlog.workloadSelection[workload].map((service, index2) => {
                return (
                  <Text key={index1 + ' ' + index2} parent="div">
                    {contaienrlog.namespaceSelection} / {workload} / {service}
                  </Text>
                );
              });
            });
          } else {
            return (
              <Text key={index} parent="div">
                {contaienrlog.namespaceSelection} / {'全部容器'}
              </Text>
            );
          }
        })}
      </FormPanel.Item>
    );
  }
  _renderContainerFileServices() {
    const {
      containerFileNamespace,
      containerFileWorkloadType,
      containerFileWorkload,
      containerFilePaths
    } = this.props.logStashEdit;
    return (
      <React.Fragment>
        <FormPanel.Item text label={t('日志类型')}>
          {t('指定容器文件')}
        </FormPanel.Item>
        <FormPanel.Item label={t('工作负载')} text>
          <span>
            {containerFileNamespace} / {containerFileWorkloadType} / {containerFileWorkload}
            <br />
          </span>
        </FormPanel.Item>
        <FormPanel.Item label={t('采集路径')} text>
          {containerFilePaths.map((item, index) => (
            <span key={index}>
              容器名称：{item.containerName} <span style={{ marginLeft: '5px' }}>路径：{item.containerFilePath}</span>
              <br />
            </span>
          ))}
        </FormPanel.Item>
      </React.Fragment>
    );
  }

  _renderEsInfo() {
    const { esAddress, indexName, esUsername, esPassword } = this.props.logStashEdit;
    return (
      <React.Fragment>
        <FormPanel.Item text label={t('类型')}>
          elasticsearch
        </FormPanel.Item>
        <FormPanel.Item text label={t('Elasticsearch地址')}>
          {esAddress}
        </FormPanel.Item>
        <FormPanel.Item text label={t('索引')}>
          {indexName}
        </FormPanel.Item>
        <FormPanel.Item text label={t('用户名')}>
          {esUsername}
        </FormPanel.Item>
        <FormPanel.Item text label={t('密码')}>
          {esPassword}
        </FormPanel.Item>
      </React.Fragment>
    );
  }
  render() {
    let { logSelection, route, logStashEdit, clusterSelection, regionSelection, projectSelection } = this.props,
      { logStashName, logMode, metadatas, consumerMode, isSelectedAllNamespace, nodeLogPath } = logStashEdit;

    return (
      <Card>
        <Card.Body>
          <FormPanel isNeedCard={false} title={t('基本信息')}>
            <FormPanel.Item text label={t('日志规则名称')}>
              {route.queries['stashName']}
            </FormPanel.Item>
            {window.location.href.includes('tkestack-project') ? (
              <FormPanel.Item text label={t('所属业务')} >
                {projectSelection}
              </FormPanel.Item>
            ) : (
              <FormPanel.Item text label={t('所属集群')}>
                <a
                  href={
                    clusterSelection[0] && regionSelection
                      ? '/tkestack/cluster/sub/list/basic/info?rid=' +
                        regionSelection.value +
                        '&clusterId=' +
                        clusterSelection[0].metadata.name +
                        '&np=default'
                      : 'javascript:;'
                  }
                >
                  {clusterSelection[0] && clusterSelection[0].id}
                </a>
                <span className="text-weak">
                  {clusterSelection[0] && '（' + clusterSelection[0].metadata.name + '）'}
                </span>
              </FormPanel.Item>
            )}
            <FormPanel.Item text label={t('创建时间')}>
              {logSelection[0] &&
                dateFormatter(new Date(logSelection[0].metadata.creationTimestamp), 'YYYY-MM-DD HH:mm:ss')}
            </FormPanel.Item>
          </FormPanel>
          <hr />
          <FormPanel isNeedCard={false} title={t('日志信息')}>
            {logMode === 'container' && this._renderContainerServices()}
            {logMode === 'containerFile' && this._renderContainerFileServices()}
            {logMode === 'node' && (
              <React.Fragment>
                <FormPanel.Item text label={t('日志类型')}>
                  {t('指定主机文件')}
                </FormPanel.Item>
                <FormPanel.Item text label={t('收集路径')}>
                  {nodeLogPath}
                </FormPanel.Item>
                <FormPanel.Item text label={t('标签')}>
                  {metadatas.length === 0 ? t('无') : this.formatLabels(metadatas)}
                </FormPanel.Item>
              </React.Fragment>
            )}
          </FormPanel>
          <hr />
          <FormPanel isNeedCard={false} title={t('消费端')}>
            {consumerMode === 'kafka' && this._renderKafkaInfo()}
            {consumerMode === 'es' && this._renderEsInfo()}
          </FormPanel>
        </Card.Body>
      </Card>
    );
  }
}
