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
import { UserAssociatePanel } from '../associate/UserAssociatePanel';
import { checkCustomVisible } from '@src/modules/common/components/permission-provider';

const tips = seajs.require('tips');

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class BaseInfoPanel extends React.Component<RootProps, {}> {
  render() {
    const { actions, route, chartGroupCreation, chartGroupValidator, projectList, commonUserAssociation } = this.props;
    const action = actions.chartGroup.create.addChartGroupWorkflow;
    const { chartGroupAddWorkflow } = this.props;
    const workflow = chartGroupAddWorkflow;

    const isBase64 = str => {
      if (str === '' || str.trim() === '') {
        return false;
      }
      try {
        return btoa(atob(str)) === str;
      } catch (error) {
        return false;
      }
    };
    /** 提交 */
    const perform = () => {
      actions.chartGroup.create.validator.validate(null, async r => {
        if (isValid(r)) {
          const chartGroup: ChartGroup = Object.assign({}, chartGroupCreation);
          if (chartGroup.spec.importedInfo && chartGroup.spec.importedInfo.password) {
            if (!isBase64(chartGroup.spec.importedInfo.password)) {
              chartGroup.spec.importedInfo.password = btoa(chartGroup.spec.importedInfo.password);
            }
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
      router.navigate({ mode: '', sub: 'chartgroup' }, route.queries);
    };
    const failed = workflow.operationState === OperationState.Done && !isSuccessWorkflow(workflow);

    return (
      <FormPanel vactions={actions.chartGroup.create.validator} formvalidator={chartGroupValidator}>
        <FormPanel.Item
          required
          label={t('仓库名称')}
          vkey="spec.name"
          input={{
            style: { width: 320 },
            placeholder: t('请输入仓库名称，不超过60个字符'),
            value: chartGroupCreation.spec.name,
            onChange: value =>
              actions.chartGroup.create.updateCreationState({
                spec: Object.assign({}, chartGroupCreation.spec, { name: value, displayName: value })
              })
          }}
        />
        <FormPanel.Item label={t('权限范围')} vkey="spec.visibility">
          <FormPanel.Radios
            value={chartGroupCreation.spec.visibility}
            options={[
              { value: 'User', text: '指定用户' },
              ...(checkCustomVisible('platform.registry.project') ? [{ value: 'Project', text: '指定业务' }] : []),
              { value: 'Public', text: '公共' }
            ]}
            onChange={value => {
              const obj = { visibility: value };
              switch (value) {
                case 'User': {
                  /** 已选中的数据 */
                  obj['users'] = commonUserAssociation.users ? commonUserAssociation.users.map(e => e.name) : [];
                  obj['projects'] = [];
                  break;
                }
                case 'Project': {
                  obj['users'] = [];
                  obj['projects'] = chartGroupCreation.spec.projects || [];
                  break;
                }
                case 'Public': {
                  obj['users'] = [];
                  obj['projects'] = [];
                  break;
                }
              }
              actions.chartGroup.create.updateCreationState({
                spec: Object.assign({}, chartGroupCreation.spec, { ...obj })
              });
            }}
          />
        </FormPanel.Item>
        {chartGroupCreation.spec && chartGroupCreation.spec.visibility === 'User' && (
          <FormPanel.Item label={t('绑定用户')} vkey="spec.users">
            <UserAssociatePanel
              onChange={selection => {
                actions.chartGroup.create.updateCreationState({
                  spec: Object.assign({}, chartGroupCreation.spec, {
                    users: selection.map(e => e.name)
                  })
                });
              }}
            />
          </FormPanel.Item>
        )}
        {chartGroupCreation.spec && chartGroupCreation.spec.visibility === 'Project' && (
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
                    projects: value !== '' ? [value] : []
                  })
                })
              }
            />
          </FormPanel.Item>
        )}
        <FormPanel.Item
          label={t('导入第三方仓库')}
          vkey="spec.imported"
          checkbox={{
            value: chartGroupCreation.spec.type === 'Imported',
            onChange: value => {
              actions.chartGroup.create.updateCreationState({
                spec: Object.assign({}, chartGroupCreation.spec, { type: value ? 'Imported' : 'SelfBuilt' })
              });
            }
          }}
        />

        <FormPanel.Item
          required
          isShow={chartGroupCreation?.spec?.type === 'Imported'}
          label={t('第三方仓库地址')}
          vkey="spec.importedInfo.addr"
          input={{
            style: { width: 320 },
            placeholder: t('请输入仓库地址, 例如:  http://mirror.azure.cn/kubernetes/charts'),
            value: chartGroupCreation.spec.importedInfo.addr,
            onChange: value => {
              const info = Object.assign({}, chartGroupCreation.spec.importedInfo);
              info.addr = value;
              actions.chartGroup.create.updateCreationState({
                spec: Object.assign({}, chartGroupCreation.spec, { importedInfo: info })
              });
            }
          }}
        />
        <FormPanel.Item
          isShow={chartGroupCreation?.spec?.type === 'Imported'}
          label={t('第三方仓库用户名')}
          vkey="spec.importedInfo.username"
          input={{
            style: { width: 320 },
            placeholder: t('请输入用户名'),
            value: chartGroupCreation.spec.importedInfo.username,
            onChange: value => {
              const info = Object.assign({}, chartGroupCreation.spec.importedInfo);
              info.username = value;
              actions.chartGroup.create.updateCreationState({
                spec: Object.assign({}, chartGroupCreation.spec, { importedInfo: info })
              });
            }
          }}
        />
        <FormPanel.Item
          isShow={chartGroupCreation?.spec?.type === 'Imported'}
          label={t('第三方仓库密码')}
          vkey="spec.importedInfo.password"
          input={{
            style: { width: 320 },
            type: 'password',
            placeholder: t('请输入仓库密码'),
            value: chartGroupCreation.spec.importedInfo.password,
            onChange: value => {
              const info = Object.assign({}, chartGroupCreation.spec.importedInfo);
              info.password = value;
              actions.chartGroup.create.updateCreationState({
                spec: Object.assign({}, chartGroupCreation.spec, { importedInfo: info })
              });
            }
          }}
        />

        <FormPanel.Item
          label={t('仓库描述')}
          vkey="spec.description"
          input={{
            style: { width: 320 },
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
