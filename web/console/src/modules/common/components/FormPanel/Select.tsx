import * as React from 'react';
import { Select, SelectProps, SelectOptionWithGroup, Icon } from '@tea/component';
import { FormPanelText } from './Text';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { ListModel, ListAction } from '@tencent/redux-list';
import { FetchState } from '@tencent/qcloud-redux-fetcher';

interface FormPanelSelectProps extends SelectProps {
  model?: ListModel;
  displayField?: String | Function;
  valueField?: String | Function;
  action?: ListAction;
  label?: String;
}

function getFieldValue(record, field: String | Function) {
  if (typeof field === 'function') {
    return field(record);
  } else {
    return record[field as string];
  }
}

function FormPanelSelect({ model, displayField, valueField, action, ...props }: FormPanelSelectProps) {
  let showRefreshBtn = false;
  if (model && valueField && displayField) {
    let options: SelectOptionWithGroup[] = [];
    if (model.list.fetchState === FetchState.Ready && model.list.fetched) {
      options = model.list.data.records.map((record, index) => {
        return { text: getFieldValue(record, displayField), value: getFieldValue(record, valueField) };
      });
    }
    if (model.list.fetchState === FetchState.Fetching && model.list.loading) {
      //loading
      options = [
        {
          text: t('正在加载...'),
          value: '',
          disabled: true
        }
      ];
      props.disabled = true;
    }
    if (model.list.fetchState === FetchState.Failed) {
      //加载失败
      options = [
        {
          text: t('加载失败...'),
          value: '',
          disabled: true
        }
      ];
      props.disabled = true;
      if (action) {
        //如果有action，这里可以显示一个刷新按钮
        showRefreshBtn = true;
      }
    }
    props.options = options;
    if (action && !props.onChange) {
      props.onChange = value => {
        let selected = model.list.data.records.find(record => getFieldValue(record, valueField) === value);
        action.select(selected);
      };
    }

    if (!('value' in props) && 'selection' in model) {
      props.value = model.selection ? getFieldValue(model.selection, valueField) : '';
    }
  }
  if (!props.options || props.options.length === 0) {
    //列表为空
    props.options.push({
      text: t('没有数据'),
      value: '',
      disabled: true
    });
    props.disabled = true;
  }
  if (props.label && !props.placeholder) {
    props.placeholder = t('请选择') + props.label;
  }

  return (
    <React.Fragment>
      <Select {...props} />
      {showRefreshBtn && action && <Icon type="refresh" onClick={() => action.fetch()} />}
    </React.Fragment>
  );
}

export { FormPanelSelect, FormPanelSelectProps };
