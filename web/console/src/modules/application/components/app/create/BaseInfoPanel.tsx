import * as React from 'react';
import { connect } from 'react-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators, insertCSS, OperationState, isSuccessWorkflow } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { RootProps } from '../AppContainer';
import { Button, Bubble, Icon } from '@tencent/tea-component';
import { router } from '../../../router';
import { FormPanel } from '@tencent/ff-component';
import { InputField, TipInfo, getWorkflowError } from '../../../../../modules/common';
import { App, Chart } from '../../../models';
import { isValid } from '@tencent/ff-validator';
import { ChartTablePanel } from '../ChartTablePanel';
import { ChartActionPanel } from '../ChartActionPanel';
import { ChartValueYamlDialog } from '../ChartValueYamlDialog';
import { NamespacePanel } from './NamespacePanel';
// @ts-ignore
const tips = seajs.require('tips');

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

interface AppCreateState {
  showValueSetting?: boolean;
  projectID?: string;
}

@connect(state => state, mapDispatchToProps)
export class BaseInfoPanel extends React.Component<RootProps, AppCreateState> {
  constructor(props, context) {
    super(props, context);
    this.state = {
      showValueSetting: false,
      projectID: ''
    };
  }

  render() {
    let { actions, route, appCreation, appValidator, clusterList, namespaceList, chartList, chartInfo } = this.props;
    let action = actions.app.create.addAppWorkflow;
    const { appAddWorkflow } = this.props;
    const workflow = appAddWorkflow;

    let valueDisable = chartInfo && chartInfo.object && (chartInfo.object.loading || chartInfo.object.error);

    const versionOptions = chartList.selection
      ? chartList.selection.status.versions.map(v => {
          return {
            text: v.version,
            value: v.version
          };
        })
      : [];

    /** 提交 */
    const perform = () => {
      actions.app.create.validator.validate(null, async r => {
        if (isValid(r)) {
          let app: App = Object.assign({}, appCreation);
          action.start([app]);
          action.perform();
        } else {
          let invalid = Object.keys(r).filter(v => {
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
      router.navigate({ mode: '', sub: 'app' }, route.queries);
    };
    const failed = workflow.operationState === OperationState.Done && !isSuccessWorkflow(workflow);

    return (
      <FormPanel vactions={actions.app.create.validator} formvalidator={appValidator}>
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
        <NamespacePanel />
        <FormPanel.Item
          label={t('类型')}
          vkey="spec.type"
          segment={{
            value: appCreation.spec.type,
            options: [{ value: 'HelmV3', text: 'HelmV3' }],
            onChange: value => {
              actions.app.create.updateCreationState({
                spec: Object.assign({}, appCreation.spec, { type: value })
              });
            }
          }}
        ></FormPanel.Item>
        <FormPanel.Item label={t('Chart')} vkey="spec.chart">
          <ChartActionPanel />
          <ChartTablePanel
            onSelectChart={(chart: Chart, projectID: string) => {
              let specChart = Object.assign({}, appCreation.spec.chart);
              specChart.chartGroupName = chart.spec.chartGroupName;
              specChart.chartName = chart.spec.name;
              specChart.chartVersion = '';
              specChart.tenantID = chart.spec.tenantID;
              actions.app.create.updateCreationState({
                spec: Object.assign({}, appCreation.spec, { chart: specChart })
              });

              this.setState({ projectID: projectID });
            }}
          />
          {/* <FormPanel.InlineText theme={'warning'}>{t('当前选择Chart名称与旧版本不一致')}</FormPanel.InlineText> */}
        </FormPanel.Item>
        <FormPanel.Item
          label={t('Chart版本')}
          vkey="spec.chart"
          select={{
            value: appCreation.spec && appCreation.spec.chart ? appCreation.spec.chart.chartVersion : '',
            valueField: 'value',
            displayField: 'text',
            options: versionOptions,
            onChange: value => {
              let chart = Object.assign({}, appCreation.spec.chart);
              chart.chartVersion = value;
              actions.app.create.updateCreationState({
                spec: Object.assign({}, appCreation.spec, { chart: chart })
              });
              //加载values.yaml
              actions.app.create.chart.applyFilter({
                cluster: appCreation.spec.targetCluster,
                namespace: appCreation.metadata.namespace,
                metadata: {
                  namespace: chartList.selection ? chartList.selection.metadata.namespace : '',
                  name: chartList.selection ? chartList.selection.metadata.name : ''
                },
                chartVersion: value,
                projectID: this.state.projectID
              });
            }
          }}
        ></FormPanel.Item>
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
              let values = Object.assign({}, appCreation.spec.values);
              values.rawValues = value;
              actions.app.create.updateCreationState({
                spec: Object.assign({}, appCreation.spec, { values: values })
              });
            }}
            onClose={() => {
              this.setState({
                showValueSetting: false
              });
            }}
            yamlConfig={appCreation.spec.values.rawValues || ''}
            isShow={this.state.showValueSetting}
          />
        </FormPanel.Item>
        <FormPanel.Footer>
          <React.Fragment>
            <Button
              className="m"
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
        </FormPanel.Footer>
      </FormPanel>
    );
  }
}
