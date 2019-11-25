import * as React from 'react';
import {
  FormProps,
  Form,
  FormItemProps,
  CardProps,
  CardBodyProps,
  Card,
  Icon,
  Bubble,
  SlideTransition
} from '@tea/component';

import { FormPanelText, FormPanelTextProps } from './Text';
import { FormPanelInput, FormPanelInputProps } from './Input';
import { FormPanelSelect, FormPanelSelectProps } from './Select';

import { FormPanelCheckbox, FormPanelCheckboxProps } from './Checkbox';
import { FormPanelSwitch, FormPanelSwitchProps } from './Switch';
import { FormPanelSegment, FormPanelSegmentProps } from './Segment';

import { FormPanelHelpText, FormPanelHelpTextProps } from './HelpText';
import { FormPanelInlineText, FormPanelInlineTextProps } from './InlineText';
import { FormPanelFooter, FormPanelFooterProps } from './Footer';
import { FormPanelKeyValuePanel, FormPanelKeyValuePanelProps } from './KeyValuePanel';

import { insertCSS } from '@tencent/qcloud-lib';
import { Validation } from '../../models';

import { t, Trans } from '@tencent/tea-app/lib/i18n';
import * as classnames from 'classnames';
insertCSS(
  'FormPanelFooter',
  `
.formpanel-footer .tea-btn {
  margin-right:20px !important;
}
`
);
insertCSS(
  'FormPanelHideLabel',
  `
.formpanel-hide-label .tea-form__label{
  margin:0;
  padding:0;
  width:0;
}
`
);

interface FormPanelProps extends FormProps {
  title?: React.ReactNode;
  operation?: React.ReactNode;

  /** 是否需要 Card */
  isNeedCard?: boolean;
  cardProps?: CardProps;
  cardBodyProps?: CardBodyProps;

  labelStyle?: React.CSSProperties;
  fieldStyle?: React.CSSProperties;

  //如果是嵌套的form，需要添加这个属性，不然label宽度会变
  fixed?: boolean;
}

export function FormPanel({
  children,

  title,
  operation,

  isNeedCard = true,
  cardProps = {},
  cardBodyProps = {},

  fixed,

  labelStyle,
  fieldStyle,

  ...formProps
}: FormPanelProps) {
  //如果使用了footer，需要在下方留出足够的空间，避免重叠
  React.Children.map(children, (child: any) => {
    if (child && child.type === FormPanel.Footer) {
      if (isNeedCard) {
        cardProps.style = Object.assign({}, cardProps.style, { marginBottom: '60px' } as React.CSSProperties);
      } else {
        formProps.style = Object.assign({}, formProps.style, { marginBottom: '60px' } as React.CSSProperties);
      }
    }
  });

  let childrenWithProps: React.ReactNode = null;
  if (labelStyle || fieldStyle) {
    childrenWithProps = React.Children.map(children, (child: any) => {
      if (child && child.type === FormPanel.Item) {
        return React.cloneElement(child, {
          labelStyle: Object.assign({}, labelStyle, child.props.labelStyle),
          fieldStyle: Object.assign({}, fieldStyle, child.props.fieldStyle)
        });
      } else {
        return child;
      }
    });
  }

  if (fixed) {
    formProps.className = formProps.className ? formProps.className + ' size-full-width' : 'size-full-width';
  }

  return isNeedCard ? (
    <Card {...cardProps}>
      <Card.Body title={title} operation={operation} {...cardBodyProps}>
        <Form {...formProps}>{childrenWithProps ? childrenWithProps : children}</Form>
      </Card.Body>
    </Card>
  ) : (
    <Form {...formProps}>
      {title && <Form.Title>{title}</Form.Title>}
      {childrenWithProps ? childrenWithProps : children}
    </Form>
  );
}

interface FormPanelItemProps extends FormItemProps {
  isShow?: boolean;
  tips?: React.ReactNode;

  labelStyle?: React.CSSProperties;
  fieldStyle?: React.CSSProperties;

  //vkey?: string; //如果FormPanel定义了ValidatorInstance,这个key用来匹配对应的ValidatorConfig
  validator?: Validation;
  errorTips?: React.ReactNode;
  errorTipsStyle?: 'Icon' | 'Bubble';

  loading?: boolean;
  loadingElement?: React.ReactNode;

  //支持的组件类型
  text?: boolean;
  textProps?: FormPanelTextProps;
  input?: FormPanelInputProps;
  select?: FormPanelSelectProps;
  checkbox?: FormPanelCheckboxProps;
  Switch?: FormPanelSwitchProps;
  segment?: FormPanelSegmentProps;
  keyvalue?: FormPanelKeyValuePanelProps;
}
FormPanel.Item = function({
  isShow = true,
  tips,
  validator,
  errorTipsStyle = 'Icon',

  labelStyle,
  fieldStyle,

  // label,
  message,
  status,

  // label,

  ...formItemProps
}: FormPanelItemProps) {
  let tStatus = status;
  if (validator && validator.status > 1) {
    // tMessage = validator.message;
    tStatus = 'error';
  }

  let align =
    formItemProps.align || (formItemProps.text || formItemProps.textProps || formItemProps.checkbox ? 'top' : 'middle');

  let llabelStyle = Object.assign({}, { minWidth: 100 }, labelStyle);

  return isShow ? (
    <Form.Item
      {...formItemProps}
      className={classnames({
        ['is-' + tStatus]: tStatus,
        'formpanel-hide-label': !formItemProps.label //label为空时去掉空白
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
        {errorTipsStyle === 'Icon' && (
          <React.Fragment>
            {renderField(formItemProps)}
            {tStatus === 'error' && (
              <SlideTransition in={validator && validator.status > 1} from={[-10, 10]}>
                <Bubble placement="right" content={(validator && validator.message) || null}>
                  <Icon type="error" style={{ marginLeft: 5 }} />
                </Bubble>
              </SlideTransition>
            )}
            {message && <FormPanel.HelpText parent="div">{message}</FormPanel.HelpText>}
          </React.Fragment>
        )}
        {errorTipsStyle === 'Bubble' && (
          <React.Fragment>
            {tStatus === 'error' ? (
              <Bubble placement="right" content={(validator && validator.message) || null}>
                {renderField(formItemProps)}
              </Bubble>
            ) : (
              renderField(formItemProps)
            )}
            {message && <FormPanel.HelpText parent="div">{message}</FormPanel.HelpText>}
          </React.Fragment>
        )}
      </div>
    </Form.Item>
  ) : (
    <noscript />
  );
};

function renderField({
  children,
  text,
  textProps,
  loading,
  input,
  select,
  checkbox,
  Switch,
  segment,
  keyvalue,
  loadingElement,
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

  if (input) {
    if (typeof props.label === 'string') {
      input.label = input.label || props.label;
    }
    return <FormPanel.Input {...input} />;
  }
  if (select) {
    if (typeof props.label === 'string') {
      select.label = select.label || props.label;
    }
    return <FormPanel.Select {...select} />;
  }

  if (checkbox) return <FormPanel.Checkbox {...checkbox}>{children}</FormPanel.Checkbox>;
  if (Switch) return <FormPanel.Switch {...Switch} />;
  if (segment) return <FormPanel.Segment {...segment} />;

  if (keyvalue) return <FormPanel.KeyValuePanel {...keyvalue} />;

  if (text || textProps) {
    return <FormPanel.Text {...textProps}>{children}</FormPanel.Text>;
  }
  return children;
}

FormPanel.Title = Form.Title;
FormPanel.Action = Form.Action;
FormPanel.Control = Form.Control;
FormPanel.Text = FormPanelText;
FormPanel.Input = FormPanelInput;
FormPanel.Select = FormPanelSelect;
FormPanel.Checkbox = FormPanelCheckbox;
FormPanel.Switch = FormPanelSwitch;
FormPanel.Segment = FormPanelSegment;

FormPanel.HelpText = FormPanelHelpText;
FormPanel.InlineText = FormPanelInlineText;
FormPanel.Footer = FormPanelFooter;

FormPanel.KeyValuePanel = FormPanelKeyValuePanel;
