import * as React from 'react';
import { ICComponter } from '../../models';
import { FormPanel, LinkButton, TipInfo } from '../../../common/components';
import { Justify, Button, Text, Radio, Segment, Input, InputAdorment, Bubble } from '@tencent/tea-component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { authTypeList, computerRoleList } from '../../constants/Config';
import { Validation, initValidator } from '../../../common';
import { validateValue, Rule, RuleTypeEnum } from '../../../common/validate';

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
