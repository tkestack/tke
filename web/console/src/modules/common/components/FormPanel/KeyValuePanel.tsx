import * as React from 'react';
import { Input, InputAdorment, InputAdornmentProps, SelectOptionWithGroup, Button } from '@tea/component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { FormPanel } from './FormPanel';
import { FormPanelInputProps } from './Input';
import { FormPanelSelectProps } from './Select';
import { FormPanelSegmentProps } from './Segment';

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
          新增属性
        </Button>
      </div>
    </React.Fragment>
  );
}

export { FormPanelKeyValuePanel, FormPanelKeyValuePanelProps };
