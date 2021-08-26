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
import { LinkButton, emptyTips } from '../../../../common/components';
import { Table, TableColumn, Text, Modal, Card, Bubble, Icon, ContentView } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { router } from '../../../router';
import { allActions } from '../../../actions';
import { RootProps } from '../ChartApp';
import { dateFormat } from '../../../../../../helpers/dateUtil';
import { valueLabels1024, K8SUNIT } from '../../../../../../helpers/k8sUnitUtil';
import { DeployPanel } from './DeployPanel';
import { ChartVersion } from '../../../models';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

interface AppCreateState {
  showDeploySetting?: boolean;
  selectedVersion?: string;
}

@connect(state => state, mapDispatchToProps)
export class VersionTablePanel extends React.Component<RootProps, AppCreateState> {
  constructor(props, context) {
    super(props, context);
    this.state = {
      showDeploySetting: false,
      selectedVersion: ''
    };
  }

  render() {
    let { actions, chartEditor, removedChartVersions, route } = this.props;
    const columns: TableColumn<ChartVersion>[] = [
      {
        key: 'version',
        header: t('版本'),
        render: (x: ChartVersion) => (
          <Text parent="div" overflow>
            {x.version || '-'}
          </Text>
        )
      },
      {
        key: 'size',
        header: t('大小'),
        render: (x: ChartVersion) => {
          let size = x.chartSize ? valueLabels1024(x.chartSize, K8SUNIT.Ki) : '';
          let index = size.indexOf('.');
          return <Text parent="div">{index > -1 ? size.slice(0, index + 2) + 'K' : '-'}</Text>;
        }
      },
      {
        key: 'description',
        header: t('描述'),
        render: (x: ChartVersion) => {
          return (
            <Text
              title={x.description ? x.description : ''}
              style={{
                display: '-webkit-box',
                WebkitBoxOrient: 'vertical',
                WebkitLineClamp: 2,
                overflow: 'hidden'
              }}
            >
              {x.description || '-'}
            </Text>
          );
        }
      },
      {
        key: 'timeCreated',
        header: t('创建时间'),
        render: (x: ChartVersion) => (
          <Text parent="div">{x.timeCreated ? dateFormat(new Date(x.timeCreated), 'yyyy-MM-dd hh:mm:ss') : '-'}</Text>
        )
      },
      { key: 'operation', header: t('操作'), render: chart => this._renderOperationCell(chart) }
    ];

    return (
      <ContentView>
        <ContentView.Body>
          <Card>
            <Card.Body>
              <Table
                recordKey={(record, index) => {
                  return index.toString();
                }}
                rowDisabled={record =>
                  removedChartVersions.versions.find(v => {
                    return (
                      v.version === record.version &&
                      v.namespace === chartEditor.metadata.namespace &&
                      v.name === chartEditor.metadata.name
                    );
                  }) !== undefined
                }
                records={chartEditor.sortedVersions}
                columns={columns}
              />
            </Card.Body>
          </Card>
          <DeployPanel
            showDeploySetting={this.state.showDeploySetting}
            onClose={() => {
              this.setState({ showDeploySetting: false });
            }}
            chartInfoFilter={{
              cluster: '',
              namespace: '',
              metadata: {
                namespace: chartEditor.metadata.namespace,
                name: chartEditor.metadata.name
              },
              chartVersion: this.state.selectedVersion,
              projectID: route.queries['prj']
            }}
          />
        </ContentView.Body>
      </ContentView>
    );
  }

  /** 渲染操作按钮 */
  _renderOperationCell = (chart: ChartVersion) => {
    return (
      <React.Fragment>
        <LinkButton onClick={() => this._deployChart(chart)}>{t('部署')}</LinkButton>
        <LinkButton onClick={() => this._deleteChart(chart)}>{t('删除')}</LinkButton>
      </React.Fragment>
    );
  };

  _deployChart = (chart: ChartVersion) => {
    let { route, appCreation, chartEditor, actions } = this.props;
    //设置chart版本，加载完namespace列表后根据version自动加载values.yaml
    this.setState({ selectedVersion: chart.version });

    //设置选中的版本
    let specChart = Object.assign({}, appCreation.spec.chart);
    specChart.chartGroupName = chartEditor.spec.chartGroupName;
    specChart.chartName = chartEditor.spec.name;
    specChart.chartVersion = chart.version;
    actions.app.create.updateCreationState({
      spec: Object.assign({}, appCreation.spec, { chart: specChart })
    });

    this.setState({ showDeploySetting: true });
  };

  _deleteChart = async (chart: ChartVersion) => {
    let { actions, route, chartEditor } = this.props;
    const yes = await Modal.confirm({
      message: t('确定删除版本：') + `${chart.version}？`,
      description: <p className="text-danger">{t('删除该Chart后，相关数据将永久删除，请谨慎操作。')}</p>,
      okText: t('删除'),
      cancelText: t('取消')
    });
    if (yes) {
      actions.chart.detail.addRemovedChartVersion({
        namespace: chartEditor.metadata.namespace,
        name: chartEditor.metadata.name,
        version: chart.version
      });

      actions.chart.detail.removeChartVersionWorkflow.start([chart], {
        chartGroupName: chartEditor.spec.chartGroupName,
        chartName: chartEditor.spec.name,
        chartVersion: chart.version,
        chartDetailFilter: {
          namespace: chartEditor.metadata.namespace,
          name: chartEditor.metadata.name,
          projectID: route.queries['prj']
        }
      });
      actions.chart.detail.removeChartVersionWorkflow.perform();
    }
  };
}
