import * as React from 'react';
import { connect } from 'react-redux';
import { RootProps } from '../ChartGroupApp';
import { FormPanel } from '@tencent/ff-component';
import { TipInfo, getWorkflowError, InputField } from '../../../../../modules/common';
import { Button, Tabs, TabPanel, Card } from '@tea/component';
import { dateFormat } from '../../../../../../helpers/dateUtil';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators, OperationState, isSuccessWorkflow } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { isValid } from '@tencent/ff-validator';
import { ChartGroup } from '../../../models';
// @ts-ignore
const tips = seajs.require('tips');

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class BaseInfoPanel extends React.Component<RootProps> {
  render() {
    let { actions, chartGroupEditor, route, chartGroupValidator, projectList } = this.props;

    let action = actions.chartGroup.detail.updateChartGroupWorkflow;
    const { chartGroupUpdateWorkflow } = this.props;
    const workflow = chartGroupUpdateWorkflow;

    /** 提交 */
    const perform = () => {
      actions.chartGroup.detail.validator.validate(null, async r => {
        if (isValid(r)) {
          let chartGroup: ChartGroup = Object.assign({}, chartGroupEditor);
          action.start([chartGroup]);
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
      actions.chartGroup.detail.updateEditorState({ v_editing: false });
    };
    const failed = workflow.operationState === OperationState.Done && !isSuccessWorkflow(workflow);

    const typeMap = {
      personal: '个人',
      project: '业务',
      system: '系统'
    };
    const visibilityMap = {
      Public: '公有',
      Private: '私有'
    };
    let projects = [];
    if (chartGroupEditor.spec.type === 'project') {
      projectList.list.data.records.forEach(i => {
        if (chartGroupEditor.spec.projects.indexOf(i.metadata.name) > -1) {
          projects.push(i.spec.displayName);
        }
      });
    }
    if (projects.length === 0) {
      projects.push(t('无权限'));
    }
    return (
      <React.Fragment>
        <Card>
          <Card.Body
            title={t('基本信息')}
            subtitle={
              <React.Fragment>
                <Button type="link" onClick={e => actions.chartGroup.detail.updateEditorState({ v_editing: true })}>
                  {t('编辑')}
                </Button>
              </React.Fragment>
            }
          >
            <FormPanel
              isNeedCard={false}
              vactions={actions.chartGroup.detail.validator}
              formvalidator={chartGroupValidator}
            >
              <FormPanel.Item text label={t('仓库ID')}>
                {chartGroupEditor.metadata.name}
              </FormPanel.Item>
              <FormPanel.Item text label={t('仓库名称')}>
                {chartGroupEditor.spec.name}
              </FormPanel.Item>
              {/* {!chartGroupEditor.v_editing ? (
                <FormPanel.Item text label={t('仓库别名')}>
                  {chartGroupEditor.spec.displayName}
                </FormPanel.Item>
              ) : (
                <FormPanel.Item
                  label={t('仓库别名')}
                  vkey="spec.displayName"
                  input={{
                    placeholder: t('请输入仓库别名，不超过60个字符'),
                    value: chartGroupEditor.spec.displayName,
                    onChange: value =>
                      actions.chartGroup.detail.updateEditorState({
                        spec: Object.assign({}, chartGroupEditor.spec, { displayName: value })
                      })
                  }}
                />
              )} */}
              <FormPanel.Item text label={t('仓库类型')}>
                {typeMap[chartGroupEditor.spec.type] +
                  (chartGroupEditor.spec.type === 'project' ? '(' + projects.join(',') + ')' : '') || '-'}
              </FormPanel.Item>
              {chartGroupEditor.v_editing && chartGroupEditor.spec.type === 'project' && (
                <FormPanel.Item label={t('绑定业务')} vkey="spec.projects">
                  <FormPanel.Select
                    showRefreshBtn={true}
                    value={
                      chartGroupEditor.spec.projects && chartGroupEditor.spec.projects.length > 0
                        ? chartGroupEditor.spec.projects[0]
                        : ''
                    }
                    model={projectList}
                    action={actions.project.list}
                    valueField={x => x.metadata.name}
                    displayField={x => `${x.metadata.name}(${x.spec.displayName})`}
                    onChange={value =>
                      actions.chartGroup.detail.updateEditorState({
                        spec: Object.assign({}, chartGroupEditor.spec, {
                          projects: chartGroupEditor.spec.type === 'project' && value !== '' ? [value] : []
                        })
                      })
                    }
                  />
                </FormPanel.Item>
              )}
              {!chartGroupEditor.v_editing ? (
                <FormPanel.Item text label={t('仓库权限')}>
                  {visibilityMap[chartGroupEditor.spec.visibility] || '-'}
                </FormPanel.Item>
              ) : (
                <FormPanel.Item label={t('仓库权限')} vkey="spec.visibility">
                  <FormPanel.Radios
                    value={chartGroupEditor.spec.visibility}
                    options={[
                      { value: 'Public', text: '公有' },
                      { value: 'Private', text: '私有' }
                    ]}
                    onChange={value =>
                      actions.chartGroup.detail.updateEditorState({
                        spec: Object.assign({}, chartGroupEditor.spec, { visibility: value })
                      })
                    }
                  />
                </FormPanel.Item>
              )}
              {!chartGroupEditor.v_editing ? (
                <FormPanel.Item text label={t('仓库描述')}>
                  {chartGroupEditor.spec.description || '无'}
                </FormPanel.Item>
              ) : (
                <FormPanel.Item
                  label={t('仓库描述')}
                  vkey="spec.description"
                  input={{
                    multiline: true,
                    placeholder: t('请输入仓库描述，不超过255个字符'),
                    value: chartGroupEditor.spec.description,
                    onChange: value =>
                      actions.chartGroup.detail.updateEditorState({
                        spec: Object.assign({}, chartGroupEditor.spec, { description: value })
                      })
                  }}
                />
              )}
              <FormPanel.Item text label={t('创建时间')}>
                {dateFormat(new Date(chartGroupEditor.metadata.creationTimestamp), 'yyyy-MM-dd hh:mm:ss')}
              </FormPanel.Item>
              {chartGroupEditor.v_editing && (
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
      </React.Fragment>
    );
  }
}
