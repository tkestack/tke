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
import * as classnames from 'classnames';
import * as React from 'react';
import { connect } from 'react-redux';

import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Bubble, Icon, Input, Text } from '@tencent/tea-component';

import { FormItem, isEmpty, LinkButton } from '../../../../common';
import { allActions } from '../../../actions';
import { RootProps } from '../../ClusterApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class EditResourceAnnotations extends React.Component<RootProps, {}> {
  render() {
    let { subRoot, actions } = this.props,
      { workloadEdit } = subRoot,
      { workloadAnnotations } = workloadEdit;

    // 判断，只有key 和 value 都为非空时，才能进行添加的操作
    let canAdd = isEmpty(workloadAnnotations.filter(x => !x.labelKey || !x.labelValue));

    return (
      <FormItem label={`${t('注释')}(Annotations)`}>
        <div className="form-unit">
          {this._renderAnnotationsList()}
          <LinkButton
            disabled={!canAdd}
            errorTip={t('请先完成待编辑项')}
            onClick={() => {
              actions.editWorkload.addAnnotations();
            }}
          >
            {t('新增Annotations')}
          </LinkButton>
          <Text theme="label">
            {t('只能包含大小写字母、数字及分隔符"-"、"_"、"."和"/"，且必须以大小写字母、数字开头和结尾')}
          </Text>
        </div>
      </FormItem>
    );
  }

  /** 展示 annotations的选项 */
  private _renderAnnotationsList() {
    let { actions, subRoot } = this.props,
      { workloadEdit } = subRoot,
      { workloadAnnotations } = workloadEdit;

    return workloadAnnotations.map((annotation, index) => {
      return (
        <div className="code-list" key={index} style={{ marginBottom: '5px' }}>
          <div
            className={classnames({
              'is-error': annotation.v_labelKey.status === 2 || annotation.v_labelValue.status === 2
            })}
            style={{ display: 'inline-block' }}
          >
            <Bubble placement="top" content={annotation.v_labelKey.status === 2 ? annotation.v_labelKey.message : null}>
              <Input
                placeholder="Key"
                style={{ width: '150px' }}
                value={annotation.labelKey}
                maxLength={63}
                onChange={value => {
                  actions.editWorkload.updateAnnotations({ labelKey: value }, annotation.id + '');
                }}
                onBlur={e => {
                  actions.validate.workload.validateAllWorkloadAnnotationsKey();
                }}
              />
            </Bubble>
            <span className="inline-help-text tea-mr-1n">=</span>
            <Bubble
              placement="top"
              content={annotation.v_labelValue.status === 2 ? annotation.v_labelValue.message : null}
            >
              <Input
                placeholder="Value"
                style={{ width: '150px' }}
                value={annotation.labelValue}
                maxLength={253}
                onChange={value => {
                  actions.editWorkload.updateAnnotations({ labelValue: value }, annotation.id + '');
                }}
                onBlur={e => {
                  actions.validate.workload.validateAllWorkloadAnnotationsValue();
                }}
              />
            </Bubble>
            <Icon
              style={{ cursor: 'pointer' }}
              className="tea-ml-1n"
              type="close"
              onClick={() => {
                actions.editWorkload.deleteAnnotations(annotation.id + '');
              }}
            />
          </div>
        </div>
      );
    });
  }
}
