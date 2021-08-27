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
import * as JsYAML from 'js-yaml';
import { connect } from 'react-redux';
import { RootProps } from '../ChartApp';
import { FormPanel } from '@tencent/ff-component';
import { TipInfo, getWorkflowError, InputField } from '../../../../../modules/common';
import { YamlEditorPanel, YamlDialog } from '../../../../common/components';
import { Button, Alert, Drawer, StatusTip } from '@tea/component';
import { dateFormat } from '../../../../../../helpers/dateUtil';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators, OperationState, isSuccessWorkflow } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { isValid } from '@tencent/ff-validator';
import { App, ChartInfoFilter } from '../../../models';
import { NamespacePanel } from './NamespacePanel';
import { router } from '@src/modules/registry/router.project';
let deepEqual = require('deep-equal');

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

interface AppCreateProps extends RootProps {
  showDeploySetting?: boolean;
  onClose?: Function;
  chartInfoFilter: ChartInfoFilter;
}

interface AppCreateState extends RootProps {
  yamlValidator?: {
    result: number;
    message: string;
  };
  chartInfoFilter?: ChartInfoFilter;
  showDryRunManifest?: boolean;
}

@connect(state => state, mapDispatchToProps)
export class DeployPanel extends React.Component<AppCreateProps, AppCreateState> {
  constructor(props, context) {
    super(props, context);
    this.state = {
      yamlValidator: {
        result: 0,
        message: ''
      },
      chartInfoFilter: this.props.chartInfoFilter,
      showDryRunManifest: false
    };
  }

  componentWillReceiveProps(nextProps: AppCreateProps) {
    let { chartInfoFilter } = nextProps;
    if (!deepEqual(chartInfoFilter, this.props.chartInfoFilter)) {
      this.setState({
        chartInfoFilter: chartInfoFilter
      });
    }
  }

  _handleForInputEditor = (value: string) => {
    let { actions, appCreation } = this.props;
    let values = Object.assign({}, appCreation.spec.values);
    values.rawValues = value;
    actions.app.create.updateCreationState({
      spec: Object.assign({}, appCreation.spec, { values: values })
    });
  };

  render() {
    let { actions, chartEditor, chartInfo, appValidator, appCreation, appDryRun, projectList, route } = this.props;
    let action = actions.app.create.addAppWorkflow;
    const { appAddWorkflow } = this.props;
    const workflow = appAddWorkflow;
    const versionOptions = chartEditor
      ? chartEditor.status.versions.map(v => {
          return {
            text: v.version,
            value: v.version
          };
        })
      : [];
    /** 提交 */
    const perform = (dryRun: boolean = false) => {
      if (appCreation.spec.values.rawValuesType === 'yaml') {
        try {
          JsYAML.safeLoad(appCreation.spec.values.rawValues);
          this.setState({
            yamlValidator: { result: 0, message: '' }
          });
        } catch (error) {
          this.setState({
            yamlValidator: { result: 2, message: t('Yaml格式错误') }
          });
          return;
        }
      }

      actions.app.create.validator.validate(null, async r => {
        if (isValid(r)) {
          let app: App = Object.assign({}, appCreation);
          app.spec.dryRun = dryRun;
          this.setState({ showDryRunManifest: dryRun });

          action.start([app], {
            cluster: appCreation.spec.targetCluster,
            namespace: appCreation.metadata.namespace,
            projectId: projectList.selection ? projectList.selection.metadata.name : ''
          });
          action.perform();
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
    };
    const failed = workflow.operationState === OperationState.Done && !isSuccessWorkflow(workflow);
    return (
      <Drawer
        size={'l'}
        placement={'right'}
        outerClickClosable={false}
        visible={this.props.showDeploySetting}
        title={t('创建应用')}
        footer={
          <>
            <Button
              type="primary"
              onClick={e => {
                e.preventDefault();
                perform();
              }}
            >
              {failed ? t('重试') : t('创建')}
            </Button>
            <Button
              type="weak"
              onClick={() => {
                cancel();
                this.props.onClose && this.props.onClose();
              }}
            >
              {t('取消')}
            </Button>
          </>
        }
        onClose={() => {
          cancel();
          this.props.onClose && this.props.onClose();
        }}
      >
        <FormPanel isNeedCard={false} vactions={actions.app.create.validator} formvalidator={appValidator}>
          <FormPanel.Item
            label={t('应用名称')}
            vkey="spec.name"
            message={t('最长60个字符，只能包含小写字母、数字及分隔符("-")，且必须以小写字母开头，数字或小写字母结尾')}
            input={{
              placeholder: t('请输入应用名称'),
              value: appCreation.spec.name,
              onChange: value =>
                actions.app.create.updateCreationState({
                  spec: Object.assign({}, appCreation.spec, { name: value })
                })
            }}
          />
          <NamespacePanel chartInfoFilter={this.state.chartInfoFilter} />
          <FormPanel.Item
            label={t('Chart版本')}
            vkey="spec.chart"
            select={{
              value:
                (appCreation && appCreation.spec && appCreation.spec.chart && appCreation.spec.chart.chartVersion) ||
                '',
              valueField: 'value',
              displayField: 'text',
              options: versionOptions,
              onChange: value => {
                let chart = Object.assign({}, appCreation.spec.chart);
                chart.chartVersion = value;
                actions.app.create.updateCreationState({
                  spec: Object.assign({}, appCreation.spec, { chart: chart })
                });
                //更新版本
                this.setState({
                  chartInfoFilter: {
                    ...this.state.chartInfoFilter,
                    chartVersion: value
                  }
                });
                //加载values.yaml
                actions.app.create.chart.applyFilter({
                  cluster: appCreation.spec.targetCluster,
                  namespace: appCreation.metadata.namespace,
                  metadata: {
                    namespace: chartEditor ? chartEditor.metadata.namespace : '',
                    name: chartEditor ? chartEditor.metadata.name : ''
                  },
                  chartVersion: value,
                  projectID: route.queries['prj']
                });
              }
            }}
          />
          <FormPanel.Item
            label={t('参数')}
            message={
              chartInfo.object && chartInfo.object.loading
                ? t('参数配置正在加载中')
                : chartInfo.object && chartInfo.object.error
                ? t('参数配置加载失败，请稍后重试')
                : null
            }
          >
            <YamlEditorPanel
              config={appCreation.spec.values.rawValues}
              handleInputForEditor={this._handleForInputEditor}
            />
            <TipInfo isShow={failed} type="error" isForm>
              {getWorkflowError(workflow)}
            </TipInfo>
            {this.state.yamlValidator.result === 2 && (
              <Alert type="error" style={{ marginTop: 8 }}>
                {this.state.yamlValidator.message}
              </Alert>
            )}
          </FormPanel.Item>
          <FormPanel.Item label={t('拟运行')} message={t('返回模板渲染清单，不会真正执行安装')}>
            <Button
              style={{ paddingTop: '6px' }}
              type="link"
              onClick={e => {
                e.preventDefault();
                perform(true);
              }}
            >
              {t('点击执行')}
            </Button>
          </FormPanel.Item>
        </FormPanel>
        <YamlDialog
          title={
            <span>
              {t('清单')}
              {workflow.operationState === OperationState.Performing && (
                <StatusTip status="loading" loadingText=""></StatusTip>
              )}
            </span>
          }
          onClose={() => {
            this.setState({
              showDryRunManifest: false
            });

            actions.app.create.clearDryRunState();
          }}
          yamlConfig={appDryRun && appDryRun.status ? appDryRun.status.manifest : ''}
          isShow={this.state.showDryRunManifest}
        />
      </Drawer>
    );
  }
}
