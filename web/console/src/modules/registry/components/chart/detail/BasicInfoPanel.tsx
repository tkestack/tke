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
import { RootProps } from '../ChartApp';
import { FormPanel } from '@tencent/ff-component';
import { TipInfo, getWorkflowError, InputField, Markdown } from '../../../../../modules/common';
import { Button, Tabs, TabPanel, Card, Bubble, Icon, ContentView, Drawer } from '@tea/component';
import { dateFormat } from '../../../../../../helpers/dateUtil';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators, OperationState, isSuccessWorkflow } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { isValid } from '@tencent/ff-validator';
import { Chart } from '../../../models';
import { DeployPanel } from './DeployPanel';

const tips = seajs.require('tips');

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

interface AppCreateState {
  showDeploySetting?: boolean;
}

@connect(state => state, mapDispatchToProps)
export class BasicInfoPanel extends React.Component<RootProps, AppCreateState> {
  constructor(props, context) {
    super(props, context);
    this.state = {
      showDeploySetting: false
    };
  }

  render() {
    const { actions, chartEditor, chartInfo, appCreation, route, chartValidator } = this.props;
    const action = actions.chart.detail.updateChartWorkflow;
    const { chartUpdateWorkflow } = this.props;
    const workflow = chartUpdateWorkflow;

    /** 提交 */
    const perform = () => {
      actions.chart.detail.validator.validate(null, async r => {
        if (isValid(r)) {
          const chart: Chart = Object.assign({}, chartEditor);
          action.start([chart], {
            namespace: chartEditor.metadata.namespace,
            name: chartEditor.metadata.name,
            projectID: route.queries['prj']
          });
          action.perform();
        } else {
          const invalid = Object.keys(r).filter(v => {
            return r[v].status === 2;
          });
          invalid.length > 0 && tips.error(r[invalid[0]].message.toString(), 2000);
        }
      });
    };
    /** 取消 */
    const cancel = () => {
      if (workflow.operationState === OperationState.Done) {
        action.reset();
      }
      if (workflow.operationState === OperationState.Started) {
        action.cancel();
      }
      actions.chart.detail.updateEditorState({ v_editing: false });
    };
    const failed = workflow.operationState === OperationState.Done && !isSuccessWorkflow(workflow);
    return (
      <ContentView>
        <ContentView.Body>
          <Card>
            <Card.Body title={t('基本信息')}>
              <FormPanel isNeedCard={false} vactions={actions.chart.detail.validator} formvalidator={chartValidator}>
                <FormPanel.Item text label={t('仓库名称')}>
                  {chartEditor.spec.chartGroupName}
                </FormPanel.Item>
                <FormPanel.Item text label={t('Chart名称')}>
                  {chartEditor.spec.name}
                </FormPanel.Item>
                <FormPanel.Item text label={t('Chart版本(最新修改)')}>
                  {(chartEditor.selectedVersion && chartEditor.selectedVersion.version) || '无'}
                </FormPanel.Item>
                <FormPanel.Item text label={t('模板描述')}>
                  {(chartEditor.selectedVersion && chartEditor.selectedVersion.description) || '无'}
                </FormPanel.Item>
                <FormPanel.Item>
                  <React.Fragment>
                    <Button
                      type="primary"
                      disabled={workflow.operationState === OperationState.Performing}
                      onClick={e => {
                        e.preventDefault();
                        //设置选中的版本
                        const chart = Object.assign({}, appCreation.spec.chart);
                        chart.chartGroupName = chartEditor.spec.chartGroupName;
                        chart.chartName = chartEditor.spec.name;
                        chart.chartVersion = (chartEditor.selectedVersion && chartEditor.selectedVersion.version) || '';
                        chart.tenantID = chartEditor.spec.tenantID;
                        actions.app.create.updateCreationState({
                          spec: Object.assign({}, appCreation.spec, { chart: chart })
                        });

                        this.setState({ showDeploySetting: true });
                      }}
                    >
                      {t('部署')}
                    </Button>
                    <TipInfo type="error" isForm isShow={failed}>
                      {getWorkflowError(workflow)}
                    </TipInfo>
                  </React.Fragment>
                </FormPanel.Item>
              </FormPanel>
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
              chartVersion: (chartEditor.selectedVersion && chartEditor.selectedVersion.version) || '',
              projectID: route.queries['prj']
            }}
          />
          {chartInfo &&
            chartInfo.object &&
            chartInfo.object.data &&
            chartInfo.object.data.spec &&
            chartInfo.object.data.spec.readme &&
            JSON.stringify(chartInfo.object.data.spec.readme) !== '{}' && (
              <Card>
                <Card.Body>
                  <Markdown
                    style={{ maxHeight: 700, overflow: 'auto' }}
                    text={Object.values(chartInfo.object.data.spec.readme)[0] || t('空')}
                  />
                </Card.Body>
              </Card>
            )}
        </ContentView.Body>
      </ContentView>
    );
  }
}
