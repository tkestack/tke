import * as React from 'react';
import { connect } from 'react-redux';
import { RootProps } from '../ChartApp';
import { FormPanel } from '@tencent/ff-component';
import { TipInfo, getWorkflowError, InputField } from '../../../../../modules/common';
import { Button, Tabs, TabPanel, Card, Bubble, Icon, ContentView, Drawer } from '@tea/component';
import { dateFormat } from '../../../../../../helpers/dateUtil';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators, OperationState, isSuccessWorkflow } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { isValid } from '@tencent/ff-validator';
import { Chart } from '../../../models';
import { DeployPanel } from './DeployPanel';
// @ts-ignore
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
    let { actions, chartEditor, appCreation, route, chartValidator } = this.props;

    let action = actions.chart.detail.updateChartWorkflow;
    const { chartUpdateWorkflow } = this.props;
    const workflow = chartUpdateWorkflow;

    /** 提交 */
    const perform = () => {
      actions.chart.detail.validator.validate(null, async r => {
        if (isValid(r)) {
          let chart: Chart = Object.assign({}, chartEditor);
          action.start([chart], {
            namespace: chartEditor.metadata.namespace,
            name: chartEditor.metadata.name,
            projectID: route.queries['prj']
          });
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
                        let chart = Object.assign({}, appCreation.spec.chart);
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
        </ContentView.Body>
      </ContentView>
    );
  }
}
