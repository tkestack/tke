import * as React from 'react';

import { FetchState, FFListAction, FFListModel } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Bubble, Button, Checkbox, CheckboxGroup, CheckboxGroupProps, Icon, Text } from '@tencent/tea-component';
import { SegmentOption } from '@tencent/tea-component/lib/segment/SegmentOption';

import { FormPanel, FormPanelText, FormPanelValidatable, FormPanelValidatableProps } from '../';

interface FormPanelCheckboxsProps extends CheckboxGroupProps, FormPanelValidatableProps {
  model?: FFListModel;
  displayField?: String | Function;
  valueField?: String | Function;
  action?: FFListAction;
  filter?: any;
  label?: String;
  loading?: boolean;
  disabledLoading?: boolean;

  // showRefreshBtn?: boolean;
  /**
   * Segment 中选项
   */
  options?: SegmentOption[];
  /**
   * 是否为无边框样式
   * @default false
   */
  rimless?: boolean;
}

function getFieldValue(record, field: String | Function) {
  if (typeof field === 'function') {
    return field(record);
  } else {
    return record[field as string];
  }
}

function FormPanelCheckboxs({
  action,
  model,
  displayField = model ? model.displayField : '',
  valueField = model ? model.valueField : '',
  disabledLoading = false,
  // showRefreshBtn,
  filter,

  validator,
  formvalidator,
  vkey,
  vactions,
  errorTipsStyle,
  bubblePlacement,

  onChange,

  ...props
}: FormPanelCheckboxsProps) {
  let error = false;
  let loading = false;
  let empty = false;

  let rOnChange = onChange;

  if (filter) {
    const values = Object.keys(filter).map(key => filter[key]);
    React.useEffect(() => {
      if (filter && action) {
        action.applyFilter(filter);
      }
    }, values);
  }

  if (model && valueField && displayField) {
    let options: SegmentOption[] = [];
    if (model.list.fetchState === FetchState.Ready && model.list.fetched) {
      options = model.list.data.records.map((record, index) => {
        return { text: getFieldValue(record, displayField), value: getFieldValue(record, valueField) };
      });
      empty = model.list.data.records.length === 0;
    }
    if (model.list.fetchState === FetchState.Fetching && model.list.loading) {
      loading = true;
    }
    if (model.list.fetchState === FetchState.Failed) {
      error = true;
    }
    props.options = options;
    if (action && !rOnChange) {
      rOnChange = values => {
        let selected = values.map(value =>
          model.list.data.records.find(record => getFieldValue(record, valueField) === value)
        );
        action.selects(selected);
      };
    }
    if (!('value' in props) && model.selections) {
      props.value = model.selections.map(r => getFieldValue(r, valueField));
    }
  }

  if (!disabledLoading && loading) {
    return <Icon type="loading" />;
  }

  if (error) {
    return (
      <FormPanelText>
        <Bubble placement="right" content={(model && model.list.error && model.list.error.message) || null}>
          <Icon
            type="error"
            style={{
              lineHeight: '24px',
              verticalAlign: 'middle'
            }}
          />
        </Bubble>
        <Text theme="danger" style={{ lineHeight: '24px', verticalAlign: 'middle' }}>
          {t('加载失败')}
        </Text>
        {action && (
          <Button
            icon="refresh"
            style={
              {
                // lineHeight: '20px',
                // verticalAlign: 'middle'
              }
            }
            onClick={() => action.fetch()}
          />
        )}
      </FormPanelText>
    );
  }

  if (empty) {
    return <FormPanel.Text>{t('暂无数据')}</FormPanel.Text>;
  }

  let validatableProps = {
    validator,
    formvalidator,
    vkey,
    vactions,
    errorTipsStyle,
    bubblePlacement
  };

  let onChangeWrap =
    vactions && vkey
      ? (values: string[], context) => {
          rOnChange && rOnChange(values, context);
          vactions && vkey && vactions.validate(vkey);
        }
      : rOnChange;
  return (
    <React.Fragment>
      <FormPanelValidatable {...validatableProps}>
        <CheckboxGroup value={props.value} onChange={onChangeWrap}>
          {props.options.map((option, index) => (
            <Checkbox key={index} name={option.value}>
              {option.text}
            </Checkbox>
          ))}
        </CheckboxGroup>
      </FormPanelValidatable>
    </React.Fragment>
  );
}

export { FormPanelCheckboxs, FormPanelCheckboxsProps };
