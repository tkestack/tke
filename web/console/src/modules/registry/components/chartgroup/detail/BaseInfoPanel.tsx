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
import { UserAssociatePanel } from '../associate/UserAssociatePanel';

const tips = seajs.require('tips');

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class BaseInfoPanel extends React.Component<RootProps> {
  render() {
    let {
      actions,
      chartGroupEditor,
      route,
      chartGroupValidator,
      projectList,
      commonUserAssociation,
      userPlainList
    } = this.props;

    const action = actions.chartGroup.detail.updateChartGroupWorkflow;
    const { chartGroupUpdateWorkflow } = this.props;
    const workflow = chartGroupUpdateWorkflow;

    /** 提交 */
    const perform = () => {
      actions.chartGroup.detail.validator.validate(null, async r => {
        if (isValid(r)) {
          let chartGroup: ChartGroup = Object.assign({}, chartGroupEditor);
          if (chartGroup.spec.importedInfo && chartGroup.spec.importedInfo.password) {
            chartGroup.spec.importedInfo.password = btoa(chartGroup.spec.importedInfo.password);
          }
          action.start([chartGroup]);
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
      actions.chartGroup.detail.updateEditorState({ v_editing: false });
    };
    const failed = workflow.operationState === OperationState.Done && !isSuccessWorkflow(workflow);

    const typeMap = {
      SelfBuilt: '自建',
      Imported: '导入',
      System: '平台'
    };
    const visibilityMap = {
      Public: '公共',
      User: '指定用户',
      Project: '指定业务'
    };

    let projects = [];
    if (chartGroupEditor.spec.visibility === 'Project') {
      projectList.list.data.records.forEach(i => {
        if (chartGroupEditor.spec.projects && chartGroupEditor.spec.projects.indexOf(i.metadata.name) > -1) {
          projects.push(i.spec.displayName);
        }
      });
    }
    if (projects.length === 0) {
      projects.push(t('无权限'));
    }

    let users = [];
    if (chartGroupEditor.spec.visibility === 'User') {
      userPlainList.list.data.records.forEach(i => {
        if (chartGroupEditor.spec.users && chartGroupEditor.spec.users.indexOf(i.name) > -1) {
          users.push(i.displayName);
        }
      });
    }
    if (users.length === 0) {
      users.push(t('无'));
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
              {!chartGroupEditor.v_editing ? (
                <FormPanel.Item text label={t('权限范围')}>
                  {visibilityMap[chartGroupEditor.spec.visibility] +
                    (chartGroupEditor.spec.visibility === 'Project'
                      ? '(' + projects.join(',') + ')'
                      : chartGroupEditor.spec.visibility === 'User'
                      ? '(' + users.join(',') + ')'
                      : '') || '-'}
                </FormPanel.Item>
              ) : (
                <FormPanel.Item label={t('权限范围')} vkey="spec.visibility">
                  <FormPanel.Radios
                    value={chartGroupEditor.spec.visibility}
                    options={[
                      { value: 'User', text: '指定用户' },
                      { value: 'Project', text: '指定业务' },
                      { value: 'Public', text: '公共' }
                    ]}
                    onChange={value => {
                      let obj = { visibility: value };
                      switch (value) {
                        case 'User': {
                          /** 已选中的数据 */
                          obj['users'] = commonUserAssociation.users
                            ? commonUserAssociation.users.map(e => e.name)
                            : [];
                          obj['projects'] = [];
                          break;
                        }
                        case 'Project': {
                          obj['users'] = [];
                          obj['projects'] = chartGroupEditor.spec.projects || [];
                          break;
                        }
                        case 'Public': {
                          obj['users'] = [];
                          obj['projects'] = [];
                          break;
                        }
                      }
                      actions.chartGroup.detail.updateEditorState({
                        spec: Object.assign({}, chartGroupEditor.spec, { ...obj })
                      });
                    }}
                  />
                </FormPanel.Item>
              )}
              {chartGroupEditor.v_editing && chartGroupEditor.spec.visibility === 'User' && (
                <FormPanel.Item label={t('绑定用户')} vkey="spec.users">
                  <UserAssociatePanel
                    onChange={selection => {
                      actions.chartGroup.detail.updateEditorState({
                        spec: Object.assign({}, chartGroupEditor.spec, {
                          users: selection.map(e => e.name)
                        })
                      });
                    }}
                  />
                </FormPanel.Item>
              )}
              {chartGroupEditor.v_editing && chartGroupEditor.spec.visibility === 'Project' && (
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
                          projects: value !== '' ? [value] : []
                        })
                      })
                    }
                  />
                </FormPanel.Item>
              )}
              <FormPanel.Item text label={t('仓库类型')}>
                {typeMap[chartGroupEditor.spec.type] || '-'}
              </FormPanel.Item>
              {!chartGroupEditor.v_editing && chartGroupEditor.spec && chartGroupEditor.spec.type === 'Imported' && (
                <React.Fragment>
                  <FormPanel.Item text label={t('第三方仓库地址')}>
                    {chartGroupEditor.spec.importedInfo.addr || '无'}
                  </FormPanel.Item>
                  <FormPanel.Item text label={t('第三方仓库用户名')}>
                    {chartGroupEditor.spec.importedInfo.username || '无'}
                  </FormPanel.Item>
                </React.Fragment>
              )}
              {chartGroupEditor.v_editing && chartGroupEditor.spec && chartGroupEditor.spec.type === 'Imported' && (
                <React.Fragment>
                  <FormPanel.Item
                    label={t('第三方仓库地址')}
                    vkey="spec.importedInfo.addr"
                    input={{
                      placeholder: t('请输入仓库地址'),
                      value: chartGroupEditor.spec.importedInfo.addr,
                      onChange: value => {
                        let info = Object.assign({}, chartGroupEditor.spec.importedInfo);
                        info.addr = value;
                        actions.chartGroup.detail.updateEditorState({
                          spec: Object.assign({}, chartGroupEditor.spec, { importedInfo: info })
                        });
                      }
                    }}
                  />
                  <FormPanel.Item
                    label={t('第三方仓库用户名')}
                    vkey="spec.importedInfo.username"
                    input={{
                      placeholder: t('请输入用户名'),
                      value: chartGroupEditor.spec.importedInfo.username,
                      onChange: value => {
                        let info = Object.assign({}, chartGroupEditor.spec.importedInfo);
                        info.username = value;
                        actions.chartGroup.detail.updateEditorState({
                          spec: Object.assign({}, chartGroupEditor.spec, { importedInfo: info })
                        });
                      }
                    }}
                  />
                  <FormPanel.Item
                    label={t('第三方仓库密码')}
                    vkey="spec.importedInfo.password"
                    input={{
                      type: 'password',
                      placeholder: t('请输入仓库密码'),
                      value: chartGroupEditor.spec.importedInfo.password,
                      onChange: value => {
                        let info = Object.assign({}, chartGroupEditor.spec.importedInfo);
                        info.password = value;
                        actions.chartGroup.detail.updateEditorState({
                          spec: Object.assign({}, chartGroupEditor.spec, { importedInfo: info })
                        });
                      }
                    }}
                  />
                </React.Fragment>
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
