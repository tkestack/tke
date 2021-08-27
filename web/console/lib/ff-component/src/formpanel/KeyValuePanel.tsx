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

import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Button, Input, InputAdorment, InputAdornmentProps, SelectOptionWithGroup } from '@tencent/tea-component';

import { FormPanel } from './FormPanel';
import { FormPanelInputProps } from './Input';
import { FormPanelSegmentProps } from './Segment';
import { FormPanelSelectProps } from './Select';

interface KeyValue {
  key?: string;
  value?: string;
}

interface FormPanelKeyValueOptions extends SelectOptionWithGroup {
  input?: FormPanelInputProps;
  select?: FormPanelSelectProps;
  segment?: FormPanelSegmentProps;
}

interface FormPanelKeyValuePanelProps {
  onChange?: (kvs: KeyValue[], option: FormPanelKeyValueOptions) => void;
  value?: KeyValue[];
  options?: FormPanelKeyValueOptions[];
}

function FormPanelKeyValuePanel({ ...props }: FormPanelKeyValuePanelProps) {
  let kvs = props.value || [];
  return (
    <React.Fragment>
      {kvs.map((kv, index) => {
        return (
          <div key={index} style={{ marginBottom: 10 }}>
            <FormPanel.Select
              options={props.options}
              value={kvs[index].key}
              onChange={value => {
                let config = props.options.find(o => o.value === kvs[index].key);
                kvs[index].key = value;
                props.onChange && props.onChange(kvs.slice(0), config);
              }}
            />
            <FormPanel.InlineText style={{ marginLeft: 5, marginRight: 5 }}>=</FormPanel.InlineText>
            {(() => {
              let config = props.options.find(o => o.value === kvs[index].key);
              if (!config) return;
              if (config.select) {
                return (
                  <FormPanel.Select
                    {...config.select}
                    value={kvs[index].value}
                    onChange={value => {
                      let config = props.options.find(o => o.value === kvs[index].key);
                      kvs[index].value = value;
                      props.onChange && props.onChange(kvs.slice(0), config);
                    }}
                  />
                );
              } else if (config.segment) {
                return (
                  <FormPanel.Segment
                    {...config.segment}
                    value={kvs[index].value}
                    onChange={value => {
                      let config = props.options.find(o => o.value === kvs[index].key);
                      kvs[index].value = value;
                      props.onChange && props.onChange(kvs.slice(0), config);
                    }}
                  />
                );
              } else {
                return (
                  <FormPanel.Input
                    {...config.input}
                    value={kvs[index].value}
                    onChange={value => {
                      let config = props.options.find(o => o.value === kvs[index].key);
                      kvs[index].value = value;
                      props.onChange && props.onChange(kvs.slice(0), config);
                    }}
                  />
                );
              }
            })()}

            <Button
              icon="close"
              onClick={() => {
                let config = props.options.find(o => o.value === kvs[index].key);
                kvs.splice(index, 1);
                props.onChange && props.onChange(kvs.slice(0), config);
              }}
            />
          </div>
        );
      })}

      <div>
        <Button
          type="link"
          onClick={() => {
            let config = props.options[0];
            kvs.push({ key: props.options && props.options[0].value, value: '' });
            props.onChange && props.onChange(kvs.slice(0), config);
          }}
        >
          {t('新增属性')}
        </Button>
      </div>
    </React.Fragment>
  );
}

export { FormPanelKeyValuePanel, FormPanelKeyValuePanelProps };
