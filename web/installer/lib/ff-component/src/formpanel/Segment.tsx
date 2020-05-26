import * as React from 'react';

import { FetchState, FFListAction, FFListModel, insertCSS } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import {
    Bubble, Button, ControlledProps, FormText, Icon, Segment, SegmentMultiple, SegmentMultipleProps,
    SegmentProps, Text
} from '@tencent/tea-component';
import { Combine, StyledProps } from '@tencent/tea-component/lib/_type';
import { SegmentOption } from '@tencent/tea-component/lib/segment/SegmentOption';

import { FormPanelValidatable, FormPanelValidatableProps } from '../';
import { classNames } from '../lib/classname';
import { FormPanel } from './FormPanel';
import { FormPanelText } from './Text';

insertCSS(
  '@tencent/ff-component/formpanel/segment',
  `
  .is-error .ff-formpanel-segment .app-tke-fe-btn {
    border-color: #e1504a;
    color: #e1504a;
  }
);`
);

interface FormPanelSegmentProps extends Combine<StyledProps, ControlledProps<string>>, FormPanelValidatableProps {
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

function FormPanelSegment({
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
}: FormPanelSegmentProps) {
  props.className = classNames(props.className, 'ff-formpanel-segment');

  let error = false;
  let loading = false;
  let empty = false;

  let rOnChange = onChange;

  if (filter && model && action) {
    const values = Object.keys(filter)
      .map(key => filter[key])
      .concat([model]);
    React.useEffect(() => {
      let same = true;
      Object.keys(filter).forEach(key => {
        if (filter[key] !== model.query.filter[key]) {
          same = false;
        }
      });
      if (!same) {
        //如果参数不一样，重新拉取
        action.applyFilter(filter);
      } else {
        if (model.list.fetched === false && model.list.fetchState === FetchState.Ready) {
          //如果列表还没有加载过，就加载
          action.applyFilter(filter);
        }
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
      rOnChange = value => {
        let selected = model.list.data.records.find(record => getFieldValue(record, valueField) === value);
        action.select(selected);
      };
    }
    if (!('value' in props) && model.selection) {
      props.value = model.selection ? getFieldValue(model.selection, valueField) : '';
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
              lineHeight: '30px',
              verticalAlign: 'middle'
            }}
          />
        </Bubble>
        <Text theme="danger" style={{ lineHeight: '30px', verticalAlign: 'middle' }}>
          {t('加载失败')}
        </Text>
        {action && (
          <Button
            icon="refresh"
            style={{
              lineHeight: '20px',
              verticalAlign: 'middle'
            }}
            onClick={() => action.fetch()}
          />
        )}
      </FormPanelText>
    );
  }

  if (empty) {
    return <FormPanel.Text style={{ lineHeight: '30px', verticalAlign: 'middle' }}>{t('暂无数据')}</FormPanel.Text>;
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
      ? (value, context) => {
          rOnChange && rOnChange(value, context);
          vactions && vkey && vactions.validate(vkey);
        }
      : rOnChange;
  return (
    <React.Fragment>
      {error && (
        <FormPanelText>
          <Bubble placement="right" content={(model && model.list.error && model.list.error.message) || null}>
            <Icon
              type="error"
              style={{
                lineHeight: '30px',
                verticalAlign: 'middle'
              }}
            />
          </Bubble>
          <Text theme="danger" style={{ lineHeight: '30px', verticalAlign: 'middle' }}>
            {t('加载失败')}
          </Text>
          {action && (
            <Button
              icon="refresh"
              style={{
                lineHeight: '20px',
                verticalAlign: 'middle'
              }}
              onClick={() => action.fetch()}
            />
          )}
        </FormPanelText>
      )}
      {!error && (
        <FormPanelValidatable {...validatableProps}>
          <Segment {...props} options={props.options || []} onChange={onChangeWrap} />
        </FormPanelValidatable>
      )}
    </React.Fragment>
  );
}

export { FormPanelSegment, FormPanelSegmentProps };
