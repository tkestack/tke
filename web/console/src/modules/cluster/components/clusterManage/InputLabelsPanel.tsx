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
import { FormPanel } from '@tencent/ff-component';
import { t } from '@tencent/tea-app/lib/i18n';
import { Bubble, Button, Input } from '@tencent/tea-component';
import * as React from 'react';
import { initValidator, Validation } from '../../../common';
import { RuleTypeEnum, validateValue } from '../../../common/validate';

const rules = {
  key: {
    label: 'Key',
    rules: [
      RuleTypeEnum.isRequire,
      { type: RuleTypeEnum.maxLength, limit: 256 },
      { type: RuleTypeEnum.regExp, limit: /^([A-Za-z0-9][-A-Za-z0-9_\.]*)?[A-Za-z0-9]$/ }
    ]
  },
  value: {
    label: 'Value',
    rules: [
      RuleTypeEnum.isRequire,
      { type: RuleTypeEnum.maxLength, limit: 256 },
      { type: RuleTypeEnum.regExp, limit: /^([A-Za-z0-9][-A-Za-z0-9_\.]*)?[A-Za-z0-9]$/ }
    ]
  }
};

interface LabelKeyValue {
  key?: string;
  value?: string;
  v_key?: Validation;
  v_value?: Validation;
}

export function InputLabelsPanel({
  value,
  onChange
}: {
  value?: LabelKeyValue[];
  onChange: (kvs: LabelKeyValue[], isValid: boolean) => void;
}) {
  //如果使用了footer，需要在下方留出足够的空间，避免重叠

  let [kvList, setKVList] = React.useState<LabelKeyValue[]>([]);

  let canAdd = true;

  kvList.forEach(item => {
    if (item.v_key.status > 1 || item.v_value.status > 1) {
      canAdd = false;
    }
  });

  function setValue(v: LabelKeyValue[]) {
    setKVList(v);
    let isValid = true;
    v.forEach(item => {
      if (item.v_key.status > 1 || item.v_value.status > 1) isValid = false;
    });
    onChange && onChange(v, isValid);
  }

  return (
    <React.Fragment>
      {kvList.map((item, index) => {
        return (
          <div key={index} style={{ marginBottom: 5 }}>
            <span className={item.v_key.status > 1 ? 'is-error' : ''}>
              <Bubble content={item.v_key.message || null}>
                <Input
                  style={{ width: 100 }}
                  placeholder={t('请输入Key')}
                  value={item.key}
                  onChange={value => {
                    item.key = value;
                    setValue(kvList);
                  }}
                  onBlur={() => {
                    item.v_key = validateValue(item.key, rules.key);
                    setValue(kvList);
                  }}
                />
              </Bubble>
            </span>
            <FormPanel.InlineText style={{ margin: '0 5px' }}>=</FormPanel.InlineText>
            <span className={item.v_value.status > 1 ? 'is-error' : ''}>
              <Bubble content={item.v_value.message || null}>
                <Input
                  style={{ width: 150 }}
                  placeholder={t('请输入Value')}
                  value={item.value}
                  onChange={value => {
                    item.value = value;
                    setValue(kvList);
                  }}
                  onBlur={() => {
                    item.v_value = validateValue(item.value, rules.value);
                    setValue(kvList);
                  }}
                />
              </Bubble>
            </span>
            <Button
              icon="close"
              onClick={() => {
                kvList.splice(index, 1);
                setValue(kvList);
              }}
            />
          </div>
        );
      })}

      <Bubble content={canAdd ? null : t('请先完成待编辑项')}>
        <Button
          type="link"
          disabled={!canAdd}
          onClick={() => {
            kvList.push({
              key: '',
              value: '',
              v_key: initValidator,
              v_value: initValidator
            });
            setValue(kvList);
          }}
        >
          添加
        </Button>
      </Bubble>
    </React.Fragment>
  );
}
