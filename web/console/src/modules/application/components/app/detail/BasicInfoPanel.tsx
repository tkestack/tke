import * as React from 'react';
import { connect } from 'react-redux';
import { RootProps } from '../AppContainer';
import { FormPanel } from '@tencent/ff-component';
import { TipInfo, getWorkflowError, InputField } from '../../../../../modules/common';
import { Button, Tabs, TabPanel, Card, Bubble, Icon, ContentView, StatusTip } from '@tea/component';
import { dateFormat } from '../../../../../../helpers/dateUtil';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators, OperationState, isSuccessWorkflow } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { isValid } from '@tencent/ff-validator';
import { App, Chart } from '../../../models';
import { ChartTablePanel } from '../ChartTablePanel';
import { ChartActionPanel } from '../ChartActionPanel';
import { ChartValueYamlDialog } from '../ChartValueYamlDialog';
import { YamlDialog } from '../../../../common/components';

const tips = seajs.require('tips');

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

interface AppCreateState {
  showValueSetting?: boolean;
  projectID?: string;
  showDryRunManifest?: boolean;
}

@connect(state => state, mapDispatchToProps)
export class BasicInfoPanel extends React.Component<RootProps, AppCreateState> {
  constructor(props, context) {
    super(props, context);
    this.state = {
      showValueSetting: false,
      projectID: '',
      showDryRunManifest: false
    };
  }

  render() {
    const { actions, appEditor, appDryRun, appValidator, chartInfo, chartList } = this.props;

    const action = actions.app.detail.updateAppWorkflow;
    const { appUpdateWorkflow } = this.props;
    const workflow = appUpdateWorkflow;

    const valueDisable = chartInfo && chartInfo.object && (chartInfo.object.loading || chartInfo.object.error);

    const initSelectedChart = chartList.list.data.records.find(c => {
      return (
        c.spec.chartGroupName === appEditor.spec.chart.chartGroupName &&
        c.spec.name === appEditor.spec.chart.chartName &&
        c.spec.tenantID === appEditor.spec.chart.tenantID
      );
    });
    //考虑两种情况：1、手动选中；2、未手工选择，但匹配了appEditor
    const targetChart = chartList.selection ? chartList.selection : initSelectedChart;
    const versionOptions = targetChart
      ? targetChart.status.versions.map(v => {
          return {
            text: v.version,
            value: v.version
          };
        })
      : [];
    /** 提交 */
    const perform = (dryRun = false) => {
      actions.app.detail.validator.validate(null, async r => {
        if (isValid(r)) {
          const app: App = Object.assign({}, appEditor);
          app.spec.dryRun = dryRun;
          this.setState({ showDryRunManifest: dryRun });

          action.start([app]);
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
      actions.app.detail.updateEditorState({ v_editing: false });
    };
    const failed = workflow.operationState === OperationState.Done && !isSuccessWorkflow(workflow);
    return (
      <ContentView>
        <ContentView.Body>
          <Card>
            <Card.Body
              title={t('基本信息')}
              subtitle={
                <React.Fragment>
                  <Button type="link" onClick={e => actions.app.detail.updateEditorState({ v_editing: true })}>
                    {t('编辑')}
                  </Button>
                </React.Fragment>
              }
            >
              <FormPanel isNeedCard={false} vactions={actions.app.detail.validator} formvalidator={appValidator}>
                <FormPanel.Item text label={t('应用ID')}>
                  {appEditor.metadata.name}
                </FormPanel.Item>
                <FormPanel.Item text label={t('应用名称')}>
                  {appEditor.spec.name}
                </FormPanel.Item>
                <FormPanel.Item text label={t('运行集群')}>
                  {appEditor.spec.targetCluster}
                </FormPanel.Item>
                <FormPanel.Item text label={t('命名空间')}>
                  {appEditor.metadata.namespace}
                </FormPanel.Item>
                <FormPanel.Item text label={t('类型')}>
                  {appEditor.spec.type}
                </FormPanel.Item>
                {!appEditor.v_editing ? (
                  <FormPanel.Item text label={t('Chart')}>
                    {appEditor.spec.chart
                      ? appEditor.spec.chart.chartGroupName +
                        '/' +
                        appEditor.spec.chart.chartName +
                        ':' +
                        appEditor.spec.chart.chartVersion
                      : ''}
                  </FormPanel.Item>
                ) : (
                  <React.Fragment>
                    <FormPanel.Item label={t('Chart')} vkey="spec.chart">
                      <ChartActionPanel />
                      <ChartTablePanel
                        SelectedChart={targetChart ? targetChart.metadata.name : ''}
                        onSelectChart={(chart: Chart, projectID: string) => {
                          const specChart = Object.assign({}, appEditor.spec.chart);
                          specChart.chartGroupName = chart.spec.chartGroupName;
                          specChart.chartName = chart.spec.name;
                          specChart.chartVersion = '';
                          specChart.tenantID = chart.spec.tenantID;
                          actions.app.detail.updateEditorState({
                            spec: Object.assign({}, appEditor.spec, { chart: specChart })
                          });

                          this.setState({ projectID: projectID });
                        }}
                      />
                      {chartList.selection &&
                        (appEditor.spec.chart.tenantID !== chartList.selection.spec.tenantID ||
                          appEditor.spec.chart.chartName !== chartList.selection.spec.name ||
                          appEditor.spec.chart.chartGroupName !== chartList.selection.spec.chartGroupName) && (
                          <FormPanel.InlineText theme={'warning'}>
                            {t('当前选择Chart名称与旧版本不一致')}
                          </FormPanel.InlineText>
                        )}
                    </FormPanel.Item>
                    <FormPanel.Item
                      label={t('Chart版本')}
                      vkey="spec.chart"
                      select={{
                        value: appEditor.spec && appEditor.spec.chart ? appEditor.spec.chart.chartVersion : '',
                        valueField: 'value',
                        displayField: 'text',
                        options: versionOptions,
                        onChange: value => {
                          const chart = Object.assign({}, appEditor.spec.chart);
                          chart.chartVersion = value;
                          actions.app.detail.updateEditorState({
                            spec: Object.assign({}, appEditor.spec, { chart: chart })
                          });
                          //加载values.yaml
                          actions.app.detail.chart.applyFilter({
                            cluster: appEditor.spec.targetCluster,
                            namespace: appEditor.metadata.namespace,
                            metadata: {
                              namespace: targetChart ? targetChart.metadata.namespace : '',
                              name: targetChart ? targetChart.metadata.name : ''
                            },
                            chartVersion: value,
                            projectID: this.state.projectID
                          });
                        }
                      }}
                    ></FormPanel.Item>
                  </React.Fragment>
                )}
                <FormPanel.Item
                  vkey={'spec.values.rawValues'}
                  errorTipsStyle={'Message'}
                  label={t('参数')}
                  message={t('更新时如果选择不同版本的Helm Chart,参数设置将被覆盖')}
                  text={true}
                  textProps={{
                    onEdit: () => {
                      if (!valueDisable) {
                        this.setState({ showValueSetting: true });
                      }
                    }
                  }}
                >
                  values.yaml
                  {valueDisable && (
                    <Bubble
                      content={
                        chartInfo.object && chartInfo.object.loading
                          ? t('参数配置正在加载中')
                          : chartInfo.object && chartInfo.object.error
                          ? t('参数配置加载失败，请稍后重试')
                          : null
                      }
                    >
                      <Icon type={chartInfo.object && chartInfo.object.loading ? 'warning' : 'error'} />
                    </Bubble>
                  )}
                  <ChartValueYamlDialog
                    onChange={value => {
                      const values = Object.assign({}, appEditor.spec.values);
                      values.rawValues = value;
                      actions.app.detail.updateEditorState({
                        spec: Object.assign({}, appEditor.spec, { values: values })
                      });
                    }}
                    onClose={() => {
                      this.setState({
                        showValueSetting: false
                      });
                    }}
                    yamlConfig={appEditor.spec.values.rawValues || ''}
                    isShow={this.state.showValueSetting}
                  />
                </FormPanel.Item>
                {appEditor.v_editing && (
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
                )}
                <FormPanel.Item text label={t('创建时间')}>
                  {dateFormat(new Date(appEditor.metadata.creationTimestamp), 'yyyy-MM-dd hh:mm:ss')}
                </FormPanel.Item>
                {appEditor.v_editing && (
                  <FormPanel.Item>
                    <React.Fragment>
                      <Button
                        type="primary"
                        disabled={workflow.operationState === OperationState.Performing}
                        onClick={e => {
                          e.preventDefault();
                          perform();
                        }}
                      >
                        {failed ? t('重试') : t('提交')}
                      </Button>
                      <Button
                        type="weak"
                        onClick={e => {
                          e.preventDefault();
                          cancel();
                        }}
                      >
                        {t('取消')}
                      </Button>
                      <TipInfo type="error" isForm isShow={failed}>
                        {getWorkflowError(workflow)}
                      </TipInfo>
                    </React.Fragment>
                  </FormPanel.Item>
                )}
              </FormPanel>
            </Card.Body>
          </Card>
        </ContentView.Body>
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
            actions.app.detail.clearDryRunState();
          }}
          yamlConfig={appDryRun && appDryRun.status ? appDryRun.status.manifest : ''}
          isShow={this.state.showDryRunManifest}
        />
      </ContentView>
    );
  }
}
