import * as React from 'react';

import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Card, CardBodyProps, CardProps, Form, FormProps } from '@tencent/tea-component';

import { FormPanelCheckbox, FormPanelCheckboxProps } from './Checkbox';
import { FormPanelCheckboxs, FormPanelCheckboxsProps } from './Checkboxs';
import { FormPanelFooter, FormPanelFooterProps } from './Footer';
import { FormPanelHelpText, FormPanelHelpTextProps } from './HelpText';
import { FormPanelInlineText, FormPanelInlineTextProps } from './InlineText';
import { FormPanelInput, FormPanelInputProps } from './Input';
import { FormPanelInputNumber, FormPanelInputNumberProps } from './InputNumber';
import { FormPanelItem, FormPanelItemProps } from './Item';
import { FormPanelKeyValuePanel, FormPanelKeyValuePanelProps } from './KeyValuePanel';
import { FormPanelRadios, FormPanelRadiosProps } from './Radios';
import { FormPanelSegment, FormPanelSegmentProps } from './Segment';
import { FormPanelSelect, FormPanelSelectProps } from './Select';
import { FormPanelSwitch, FormPanelSwitchProps } from './Switch';
import { FormPanelText, FormPanelTextProps } from './Text';
import { FormPanelValidatable, FormPanelValidatableProps } from './Validatable';

// import { Validation } from '../../models';

export interface Validation {
  /**验证状态 0: 初始状态；1：校验通过；2：校验不通过；*/
  status?: number;

  /**结果描述 */
  message?: string | React.ReactNode;

  /**
   * 返回的校验列表
   * 目前仅 CIDR 有使用
   */
  list?: any[];
}

interface FormPanelProps extends FormProps, FormPanelValidatableProps {
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

  validator,
  formvalidator,
  vkey,
  vactions,
  errorTipsStyle,

  ...formProps
}: FormPanelProps) {
  let cardEl = React.createRef() as any;

  //如果使用了footer，需要在下方留出足够的空间，避免重叠
  React.Children.map(children, (child: any, index) => {
    if (child && child.type === FormPanel.Footer) {
      if (isNeedCard) {
        cardProps.style = Object.assign({}, cardProps.style, { marginBottom: '60px' } as React.CSSProperties);
      } else {
        formProps.style = Object.assign({}, formProps.style, { marginBottom: '60px' } as React.CSSProperties);
      }
    }
  });

  let childrenWithProps: React.ReactNode = React.Children.map(children, (child: any) => {
    if (
      child &&
      child.type === FormPanel.Item &&
      (labelStyle || fieldStyle || validator || formvalidator || vkey || vactions || errorTipsStyle)
    ) {
      //向下传递样式 & ValidatableProps
      return React.cloneElement(child, {
        labelStyle: Object.assign({}, labelStyle, child.props.labelStyle),
        fieldStyle: Object.assign({}, fieldStyle, child.props.fieldStyle),
        validator: child.props.validator || validator,
        formvalidator: child.props.formvalidator || formvalidator,
        vkey: child.props.vkey || vkey,
        vactions: child.props.vactions || vactions,
        errorTipsStyle: child.props.errorTipsStyle || errorTipsStyle
      });
    } else if (child && child.type === FormPanel.Footer && isNeedCard) {
      return React.cloneElement(child, { cardRef: cardEl });
    } else {
      return child;
    }
  });

  if (fixed) {
    formProps.className = formProps.className ? formProps.className + ' size-full-width' : 'size-full-width';
  }

  return isNeedCard ? (
    <Card {...cardProps} ref={cardEl}>
      <Card.Body title={title} operation={operation} {...cardBodyProps}>
        <Form {...formProps}>{childrenWithProps}</Form>
      </Card.Body>
    </Card>
  ) : (
    <Form {...formProps}>
      {title && <Form.Title>{title}</Form.Title>}
      {childrenWithProps}
    </Form>
  );
}

FormPanel.Checkbox = FormPanelCheckbox;
FormPanel.Checkboxs = FormPanelCheckboxs;
FormPanel.Footer = FormPanelFooter;
FormPanel.HelpText = FormPanelHelpText;
FormPanel.InlineText = FormPanelInlineText;
FormPanel.Input = FormPanelInput;
FormPanel.InputNumber = FormPanelInputNumber;
FormPanel.Item = FormPanelItem;
FormPanel.KeyValuePanel = FormPanelKeyValuePanel;
FormPanel.Radios = FormPanelRadios;
FormPanel.Segment = FormPanelSegment;
FormPanel.Select = FormPanelSelect;
FormPanel.Switch = FormPanelSwitch;
FormPanel.Text = FormPanelText;

FormPanel.Title = Form.Title;
FormPanel.Action = Form.Action;
FormPanel.Control = Form.Control;

FormPanel.Validatable = FormPanelValidatable;
