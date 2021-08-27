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
import * as ReactCSSTransitionGroup from 'react-addons-css-transition-group';

import { Button } from '@tea/component';
import {
    BaseReactProps, isComponentOfType, isSuccessWorkflow, OperationState, slide, WorkflowState
} from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';

import { ErrorGuide, ErrorTip, TipInfo } from '../';
import { FormLayout } from '../../layouts';
import { Step } from './Step';

export interface StepItem {
  /**步骤序号 */
  no: number;

  /**步骤显示的文本 */
  text: string;

  /**下一步执行的函数 */
  stepNext?: () => any;
}

export interface StepTabProps extends BaseReactProps {
  /**步骤列表 */
  list: StepItem[];

  /**操作流 */
  workflow?: WorkflowState<any, any>;

  /**错误指引 */
  guide?: ErrorGuide;

  errorTips?: string;
}

interface StepTabState {
  /**当前步骤 */
  current: number;
}

export class StepTab extends React.Component<StepTabProps, StepTabState> {
  state = {
    current: 1
  };

  handStepPre() {
    this.setState({ current: this.state.current - 1 });
  }

  handStepNext() {
    let { list } = this.props,
      { current } = this.state;
    let finder = list.find(item => item.no === current);
    if (finder.stepNext()) {
      current < list.length && this.setState({ current: this.state.current + 1 });
    }
  }

  render() {
    const { list, workflow, guide, errorTips } = this.props,
      { current } = this.state;
    let failed = workflow && workflow.operationState === OperationState.Done && !isSuccessWorkflow(workflow);
    return (
      <div className="server-add-box">
        <FormLayout>
          <div className={'tc-15-step col' + list.length}>
            <ol>
              {list.map((item, index) => (
                <Step key={index} {...item} current={current} />
              ))}
            </ol>
          </div>
          <div>
            <ReactCSSTransitionGroup transitionLeave={false} {...slide({ x: 0, y: -10 })}>
              {this._getBody()}
            </ReactCSSTransitionGroup>
          </div>
          <ul className="form-list jiqun">
            <li className="pure-text-row fixed">
              {current !== 1 && (
                <Button className="mr10" onClick={this.handStepPre.bind(this)}>
                  {t('上一步')}
                </Button>
              )}
              <Button type="primary" disabled={!!errorTips} onClick={this.handStepNext.bind(this)}>
                {current === list.length ? (failed ? t('重试') : t('完成')) : t('下一步')}
              </Button>
              <ErrorTip isShow={failed} workflow={workflow} guide={guide} />
              {errorTips && (
                <TipInfo className="error" style={{ display: 'inline-block', marginLeft: '20px', marginBottom: '0px' }}>
                  {errorTips}
                </TipInfo>
              )}
            </li>
          </ul>
        </FormLayout>
      </div>
    );
  }

  private _getBody() {
    let body: JSX.Element[] = [];
    let { children } = this.props,
      { current } = this.state;

    React.Children.forEach(this.props.children, (child: JSX.Element, index: number) => {
      if (isComponentOfType(child, StepTabBody)) {
        if (+child.key === +current) {
          body.push(child);
        }
      } else {
        body.push(child);
      }
    });
    return body;
  }
}

export interface StepTabBodyProps extends BaseReactProps {
  /**键值 */
  key: number;
}

export class StepTabBody extends React.Component<StepTabBodyProps, any> {
  render() {
    return <div className="tab-panel">{this.props.children}</div>;
  }
}
