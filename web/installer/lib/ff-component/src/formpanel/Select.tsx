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

import { FetchState, FFListAction, FFListModel } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import {
    Bubble, Button, Icon, InputAdorment, InputAdornmentProps, Select, SelectOptionWithGroup,
    SelectProps
} from '@tencent/tea-component';

import { FormPanelValidatable, FormPanelValidatableProps } from '../';

interface FormPanelSelectProps extends SelectProps, FormPanelValidatableProps {
  model?: FFListModel;
  displayField?: String | Function;
  valueField?: String | Function;
  groupKeyField?: String | Function;
  action?: FFListAction;
  filter?: any;
  label?: String;
  loading?: boolean;
  disabledLoading?: boolean;
  showRefreshBtn?: boolean;

  inputAdornment?: InputAdornmentProps;

  beforeFormat?: (records: any[]) => any[];

  emptyTip?: React.ReactNode;
}

function getFieldValue(record, field: String | Function) {
  if (typeof field === 'function') {
    return field(record);
  } else {
    return record[field as string];
  }
}

/* eslint-disable */
function FormPanelSelect({
  action,
  model,
  displayField = model ? model.displayField : '',
  valueField = model ? model.valueField : '',
  groupKeyField = model ? model.groupKeyField : '',
  disabledLoading = false,
  showRefreshBtn,
  filter,

  validator,
  formvalidator,
  vkey,
  vactions,
  errorTipsStyle,
  bubblePlacement,

  inputAdornment,

  onChange,

  beforeFormat,

  emptyTip,

  ...props
}: FormPanelSelectProps) {
  let error = false;
  let loading = false;
  let isEmpty = false;
  props.options = props.options || [];

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
    let options: SelectOptionWithGroup[] = [];
    if (model.list.fetchState === FetchState.Ready && model.list.fetched) {
      let records = beforeFormat ? beforeFormat(model.list.data.records) : model.list.data.records;

      options = records.map((record, index) => {
        return {
          text: getFieldValue(record, displayField),
          value: getFieldValue(record, valueField),
          groupKey: getFieldValue(record, groupKeyField)
        } as SelectOptionWithGroup;
      });
    }
    if (model.list.fetchState === FetchState.Fetching && model.list.loading) {
      //loading
      options = [
        {
          text: t('正在加载...'),
          value: ''
        }
      ];
      props.disabled = true;
      loading = true;
    }
    if (model.list.fetchState === FetchState.Failed) {
      //加载失败
      options = [
        {
          text: t('加载失败...'),
          value: ''
        }
      ];
      props.disabled = true;
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
  if (!props.options || props.options.length === 0) {
    //列表为空
    props.options.push({
      text: t('暂无数据'),
      value: ''
    });
    props.disabled = true;
    isEmpty = true;
  }
  if (props.label && !props.placeholder) {
    props.placeholder = t('请选择') + props.label;
  }

  if (props.options.length > 10) {
    props.searchable = props.searchable === undefined ? true : props.searchable;
  }
  props.type = props.type || 'simulate';
  props.appearence = props.appearence || 'button';
  props.boxSizeSync = props.boxSizeSync === undefined ? false : props.boxSizeSync;
  props.size = props.size || 'm';
  //根据size调整弹出层最小宽度
  let sizeWithMap = {
    xs: '60px',
    s: '100px',
    m: '200px',
    l: '420px',
    full: '100%',
    auto: 'auto'
  };
  let boxMinWitdh = sizeWithMap[props.size];
  props.boxStyle = props.boxStyle
    ? Object.assign({}, props.boxStyle, {
        maxWidth: props.boxStyle.maxWidth ? props.boxStyle.maxWidth : '400px',
        minWidth: props.boxStyle.minWidth ? props.boxStyle.minWidth : boxMinWitdh
      })
    : { maxWidth: '400px', minWidth: boxMinWitdh };

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

  if (inputAdornment) {
    //添加一个div.style=inline-block，为了外面包裹bubble时能正常工作
    return (
      <React.Fragment>
        <FormPanelValidatable {...validatableProps}>
          <div style={{ display: 'inline-block' }}>
            {isEmpty && emptyTip ? (
              emptyTip
            ) : (
              <InputAdorment {...inputAdornment}>
                <Select {...props} onChange={onChangeWrap} />
              </InputAdorment>
            )}
          </div>
        </FormPanelValidatable>
        {error && (
          <Bubble placement="right" content={(model && model.list.error && model.list.error.message) || null}>
            <Icon type="error" style={{ marginLeft: '5px' }} />
          </Bubble>
        )}
        {/* {!disabledLoading && loading && <Icon type="loading" />} */}
        {action && (showRefreshBtn || error) && <Button icon="refresh" onClick={() => action.fetch()} />}
      </React.Fragment>
    );
  } else {
    return (
      <React.Fragment>
        <FormPanelValidatable {...validatableProps}>
          {isEmpty && emptyTip ? emptyTip : <Select {...props} onChange={onChangeWrap} />}
        </FormPanelValidatable>
        {error && (
          <Bubble placement="right" content={(model && model.list.error && model.list.error.message) || null}>
            <Icon type="error" style={{ marginLeft: '5px' }} />
          </Bubble>
        )}
        {/* {!disabledLoading && loading && <Icon type="loading" />} */}
        {action && (showRefreshBtn || error) && <Button icon="refresh" onClick={() => action.fetch()} />}
      </React.Fragment>
    );
  }
}

export { FormPanelSelect, FormPanelSelectProps };
