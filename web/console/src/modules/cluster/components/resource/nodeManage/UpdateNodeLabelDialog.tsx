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

import { bindActionCreators, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { FormItem, InputField, LinkButton, WorkflowDialog } from '../../../../common/components';
import { isEmpty } from '../../../../common/utils';
import { allActions } from '../../../actions';
import { validatorActions } from '../../../actions/validatorActions';
import { RootProps } from '../../ClusterApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class UpdateNodeLabelDialog extends React.Component<RootProps, {}> {
  render() {
    let { actions, route, subRoot, cluster } = this.props,
      { computerState } = subRoot,
      { updateNodeLabel, labelEdition } = computerState,
      { labels } = labelEdition;
    let canAdd = isEmpty(labels.filter(x => !x.key));

    return (
      <WorkflowDialog
        caption={t('编辑Label')}
        workflow={updateNodeLabel}
        width={500}
        action={actions.workflow.updateNodeLabel}
        preAction={actions.validate.validateAllComputerLabel}
        validateAction={() => validatorActions._validateAllComputerLabel(labelEdition.labels)}
        params={{ regionId: route.queries['rid'], clusterId: route.queries['clusterId'] }}
        targets={[labelEdition]}
      >
        <FormItem label="Label" tips={t('设置节点的label')}>
          <div className="form-unit is-success">
            {this._renderLabelList()}
            <LinkButton
              disabled={!canAdd}
              errorTip={t('请先完成待编辑项')}
              onClick={() => {
                actions.computer.addLabel();
              }}
            >
              {t('新增Label')}
            </LinkButton>
            <p className="form-input-help text-weak">
              {t(
                '长度不超过63个字符，只能包含字母、数字及"-./"，必须以字母或者数字开头结尾，且不能包含"kubernetes"保留字'
              )}
            </p>
          </div>
        </FormItem>
      </WorkflowDialog>
    );
  }
  /** 展示Label的选项 */
  private _renderLabelList() {
    let { actions, subRoot } = this.props,
      { labels } = subRoot.computerState.labelEdition;
    return labels.map((label, index) => {
      return (
        <div className="code-list" key={index}>
          <div style={{ display: 'inline-block' }}>
            <InputField
              type="text"
              placeholder={t('Label名')}
              className="tc-15-input-text m"
              style={{ width: '180px' }}
              value={label.key}
              validator={label.v_key}
              disabled={label.disabled}
              disabeldTip={t('默认标签不可以编辑')}
              tipMode="popup"
              onChange={value => actions.computer.updateLabel({ key: value }, label.id + '')}
              onBlur={value => actions.validate.validateComputerLabelKey(label.id)}
            />
            <span className="inline-help-text">=</span>
            <InputField
              type="text"
              placeholder={t('Label值')}
              className="tc-15-input-text m"
              style={{ marginLeft: '5px', width: '128px' }}
              value={label.value}
              validator={label.v_value}
              disabled={label.disabled}
              tipMode="popup"
              onChange={value => actions.computer.updateLabel({ value: value }, label.id + '')}
              onBlur={value => actions.validate.validateComputerLabelValue(label.id)}
            />
            <span className="inline-help-text">
              <LinkButton disabled={label.disabled} onClick={() => actions.computer.deleteLabel(label.id + '')}>
                <i className="icon-cancel-icon" />
              </LinkButton>
            </span>
          </div>
        </div>
      );
    });
  }
}
