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

import { insertCSS } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Bubble, Form, FormItemProps, Icon } from '@tencent/tea-component';

import {
  FormPanel,
  FormPanelCheckboxProps,
  FormPanelCheckboxsProps,
  FormPanelInputNumberProps,
  FormPanelInputProps,
  FormPanelKeyValuePanelProps,
  FormPanelRadiosProps,
  FormPanelSegmentProps,
  FormPanelSelectProps,
  FormPanelSwitchProps,
  FormPanelText,
  FormPanelTextProps,
  FormPanelValidatable,
  FormPanelValidatableProps
} from '../';
import { classNames } from '../lib/classname';

insertCSS(
  '@tencent/ff-component/hidelabel',
  `
.ff-formpanel-hide-label .app-tke-fe-form__label{
  margin:0;
  padding:0;
  width:0;
}
`
);

interface FormPanelItemProps extends FormItemProps, FormPanelValidatableProps {
  isShow?: boolean;
  tips?: React.ReactNode;

  labelStyle?: React.CSSProperties;
  fieldStyle?: React.CSSProperties;

  //vkey?: string; //如果FormPanel定义了ValidatorInstance,这个key用来匹配对应的ValidatorConfig
  errorTips?: React.ReactNode;

  loading?: boolean;
  loadingElement?: React.ReactNode;

  after?: React.ReactNode;

  //支持的组件类型
  text?: boolean;
  checkbox?: FormPanelCheckboxProps;
  checkboxs?: FormPanelCheckboxsProps;
  input?: FormPanelInputProps;
  inputNumber?: FormPanelInputNumberProps;
  keyvalue?: FormPanelKeyValuePanelProps;
  radios?: FormPanelRadiosProps;
  segment?: FormPanelSegmentProps;
  select?: FormPanelSelectProps;
  Switch?: FormPanelSwitchProps;
  textProps?: FormPanelTextProps;
}
function FormPanelItem({
  isShow = true,
  tips,

  labelStyle,
  fieldStyle,

  // label,
  message,
  status,

  after,

  // label,

  ...formItemProps
}: FormPanelItemProps) {
  let align =
    formItemProps.align ||
    (formItemProps.text ||
    formItemProps.textProps ||
    formItemProps.checkbox ||
    formItemProps.checkboxs ||
    formItemProps.radios ||
    formItemProps.Switch
      ? 'top'
      : 'middle');

  let llabelStyle = Object.assign({}, { minWidth: 100 }, labelStyle);

  return isShow ? (
    <Form.Item
      {...formItemProps}
      className={classNames({
        'ff-formpanel-hide-label': !formItemProps.label //label为空时去掉空白
      })}
      label={
        formItemProps.label ? (
          <div style={llabelStyle}>
            {formItemProps.label}
            {tips && (
              <Bubble placement="right" content={tips}>
                <Icon type="info" />
              </Bubble>
            )}
          </div>
        ) : null
      }
      align={align}
    >
      <div style={fieldStyle}>
        {renderField(formItemProps)}
        {after && <FormPanel.InlineText>{after}</FormPanel.InlineText>}
        {message && <FormPanel.HelpText parent="div">{message}</FormPanel.HelpText>}
      </div>
    </Form.Item>
  ) : (
    <noscript />
  );
}

function renderField({
  children,
  text,

  checkbox,
  checkboxs,
  input,
  inputNumber,
  keyvalue,
  radios,
  segment,
  select,
  Switch,
  textProps,

  loading,
  loadingElement,

  validator,
  formvalidator,
  vkey,
  vactions,
  errorTipsStyle,
  bubblePlacement,

  ...props
}: FormPanelItemProps) {
  if (props.errorTips) {
    if (typeof props.errorTips === 'string') {
      return <FormPanelText>{props.errorTips}</FormPanelText>;
    } else {
      return <React.Fragment>{props.errorTips}</React.Fragment>;
    }
  }

  if (loading) {
    return loadingElement ? (
      loadingElement
    ) : (
      <FormPanelText>
        <Icon type="loading" />
        {t('加载中...')}
      </FormPanelText>
    );
  }

  let validatableProps = {
    validator,
    formvalidator,
    vkey,
    vactions,
    errorTipsStyle,
    bubblePlacement
  };

  if (checkbox) return <FormPanel.Checkbox {...checkbox}>{children}</FormPanel.Checkbox>;
  if (checkboxs) return <FormPanel.Checkboxs {...checkboxs}>{children}</FormPanel.Checkboxs>;

  if (input) {
    if (typeof props.label === 'string') {
      input.label = input.label || props.label;
    }
    return <FormPanel.Input {...input} {...validatableProps} />;
  }
  if (inputNumber) {
    return <FormPanel.InputNumber {...inputNumber} {...validatableProps} />;
  }
  if (keyvalue) return <FormPanel.KeyValuePanel {...keyvalue} />;
  if (radios) return <FormPanel.Radios {...radios} />;
  if (segment) return <FormPanel.Segment {...segment} {...validatableProps} />;
  if (select) {
    if (typeof props.label === 'string') {
      select.label = select.label || props.label;
    }
    return <FormPanel.Select {...select} {...validatableProps} />;
  }
  if (Switch) return <FormPanel.Switch {...Switch} />;
  if (text || textProps) {
    return <FormPanel.Text {...textProps}>{children}</FormPanel.Text>;
  }

  if (validator) {
    return <FormPanelValidatable {...validatableProps}>{children}</FormPanelValidatable>;
  } else {
    return children;
  }
}

export { FormPanelItem, FormPanelItemProps };
