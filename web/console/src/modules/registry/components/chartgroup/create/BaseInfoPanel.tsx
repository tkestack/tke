import * as React from 'react';
import { connect } from 'react-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators, insertCSS, OperationState, isSuccessWorkflow } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { RootProps } from '../ChartGroupApp';
import { Button } from '@tencent/tea-component';
import { router } from '../../../router';
import { FormPanel } from '@tencent/ff-component';
import { InputField, TipInfo, getWorkflowError } from '../../../../../modules/common';
import { ChartGroup } from '../../../models';
import { isValid } from '@tencent/ff-validator';
// @ts-ignore
const tips = seajs.require('tips');

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class BaseInfoPanel extends React.Component<RootProps, {}> {
  render() {
    let { actions, route, chartGroupCreation, chartGroupValidator, projectList, userInfo } = this.props;
    let action = actions.chartGroup.create.addChartGroupWorkflow;
    const { chartGroupAddWorkflow } = this.props;
    const workflow = chartGroupAddWorkflow;

    /** 提交 */
    const perform = () => {
      actions.chartGroup.create.validator.validate(null, async r => {
        if (isValid(r)) {
          let chartGroup: ChartGroup = Object.assign({}, chartGroupCreation);
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
      router.navigate({ mode: '', sub: 'chartgroup' }, route.queries);
    };
    const failed = workflow.operationState === OperationState.Done && !isSuccessWorkflow(workflow);

    return (
      <FormPanel vactions={actions.chartGroup.create.validator} formvalidator={chartGroupValidator}>
        <FormPanel.Item
          label={t('仓库名称')}
          vkey="spec.name"
          input={{
            placeholder: t('请输入仓库名称，不超过60个字符'),
            value: chartGroupCreation.spec.name,
            onChange: value =>
              actions.chartGroup.create.updateCreationState({
                spec: Object.assign({}, chartGroupCreation.spec, { name: value, displayName: value })
              })
          }}
        />
        {/* <FormPanel.Item
          label={t('仓库别名')}
          vkey="spec.displayName"
          input={{
            placeholder: t('请输入仓库名称，不超过60个字符'),
            value: chartGroupCreation.spec.displayName,
            onChange: value =>
              actions.chartGroup.create.updateCreationState({
                spec: Object.assign({}, chartGroupCreation.spec, { displayName: value })
              })
          }}
        /> */}
        <FormPanel.Item label={t('仓库类型')} vkey="spec.type">
          <FormPanel.Radios
            value={chartGroupCreation.spec.type}
            options={[
              { value: 'personal', text: '个人' },
              { value: 'project', text: '业务' }
            ]}
            onChange={value => {
              let obj = { type: value };
              if (value !== 'project') {
                obj['projects'] = [];
              }
              if (value === 'personal') {
                obj['name'] = userInfo.name;
                obj['displayName'] = userInfo.name;
              }
              actions.chartGroup.create.updateCreationState({
                spec: Object.assign({}, chartGroupCreation.spec, obj)
              });
            }}
          />
        </FormPanel.Item>
        {chartGroupCreation.spec && chartGroupCreation.spec.type === 'project' && (
          <FormPanel.Item label={t('绑定业务')} vkey="spec.projects">
            <FormPanel.Select
              showRefreshBtn={true}
              value={
                chartGroupCreation.spec.projects && chartGroupCreation.spec.projects.length > 0
                  ? chartGroupCreation.spec.projects[0]
                  : ''
              }
              model={projectList}
              action={actions.project.list}
              valueField={x => x.metadata.name}
              displayField={x => `${x.metadata.name}(${x.spec.displayName})`}
              onChange={value =>
                actions.chartGroup.create.updateCreationState({
                  spec: Object.assign({}, chartGroupCreation.spec, {
                    projects: chartGroupCreation.spec.type === 'project' && value !== '' ? [value] : []
                  })
                })
              }
            />
          </FormPanel.Item>
        )}
        <FormPanel.Item label={t('仓库权限')} vkey="spec.visibility">
          <FormPanel.Radios
            value={chartGroupCreation.spec.visibility}
            options={[
              { value: 'Public', text: '公有' },
              { value: 'Private', text: '私有' }
            ]}
            onChange={value =>
              actions.chartGroup.create.updateCreationState({
                spec: Object.assign({}, chartGroupCreation.spec, { visibility: value })
              })
            }
          />
        </FormPanel.Item>
        <FormPanel.Item
          label={t('仓库描述')}
          vkey="spec.description"
          input={{
            multiline: true,
            placeholder: t('请输入仓库描述，不超过255个字符'),
            value: chartGroupCreation.spec.description,
            onChange: value =>
              actions.chartGroup.create.updateCreationState({
                spec: Object.assign({}, chartGroupCreation.spec, { description: value })
              })
          }}
        />
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
