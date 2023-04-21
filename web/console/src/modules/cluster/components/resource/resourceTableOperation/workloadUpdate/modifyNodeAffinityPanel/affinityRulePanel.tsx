import { t } from '@/tencent/tea-app/lib/i18n';
import { getReactHookFormStatusWithMessage } from '@helper';
import React from 'react';
import { Controller, useFieldArray, useFormContext, useFormState } from 'react-hook-form';
import { Button, Form, Input, InputNumber, Justify, Select } from 'tea-component';
import { AddRuleButton } from './appendAffinityRuleButton';
import {
  AffinityTypeEnum,
  NodeAffinityFormType,
  NodeAffinityOperatorEnum,
  affinityRuleOperatorList
} from './constants';

export function AffinityRulePanel({ subName }: { subName: AffinityTypeEnum }) {
  const showWeight = subName === AffinityTypeEnum.Attempt;

  const { control } = useFormContext<NodeAffinityFormType>();

  const {
    fields: rules,
    append,
    remove
  } = useFieldArray({
    control,
    name: `affinityRules.${subName}`
  });

  return (
    <>
      {rules.map(({ id }, index) => (
        <Form key={id} style={{ marginBottom: 10 }}>
          <Form.Item>
            <Justify right={<Button type="icon" icon="close" onClick={() => remove(index)} />} />
          </Form.Item>
          {showWeight && (
            <Controller
              control={control}
              name={`affinityRules.force.${index}.weight`}
              render={({ field, ...another }) => (
                <Form.Item label="权重" showStatusIcon={false} {...getReactHookFormStatusWithMessage(another)}>
                  <InputNumber hideButton size="l" min={1} max={100} step={1} {...field} />
                </Form.Item>
              )}
            />
          )}

          <SubRulePanel ruleIndex={index} subName={subName} />
        </Form>
      ))}

      {subName === AffinityTypeEnum.Force && <AddRuleButton append={append} />}

      {subName === AffinityTypeEnum.Attempt && (
        <Button
          type="link"
          onClick={() =>
            append({ weight: 1, subRules: [{ key: '', operator: NodeAffinityOperatorEnum.In, value: '' }] })
          }
        >
          添加条件
        </Button>
      )}
    </>
  );
}

const SubRulePanel = ({ ruleIndex, subName }: { ruleIndex: number; subName: AffinityTypeEnum }) => {
  const { control, trigger, watch } = useFormContext<NodeAffinityFormType>();

  const prePath:
    | `affinityRules.force.${number}.subRules`
    | `affinityRules.attempt.${number}.subRules` = `affinityRules.${subName}.${ruleIndex}.subRules`;

  const {
    fields: subRules,
    append: appendSubRule,
    remove: removeSubRule,
    update
  } = useFieldArray({
    control,
    name: prePath
  });

  function needDisabledValue(operator) {
    return operator === NodeAffinityOperatorEnum.Exists || operator === NodeAffinityOperatorEnum.DoesNotExist;
  }

  const { errors } = useFormState({ control });

  const disabled = Boolean(errors?.affinityRules?.[subName]?.[ruleIndex]?.subRules);

  return (
    <>
      {subRules.map(({ id }, index) => (
        <Form.Item label="条件" key={id}>
          <Controller
            control={control}
            name={`${prePath}.${index}.key`}
            render={({ field, ...another }) => (
              <Form.Control
                showStatusIcon={false}
                style={{ display: 'inline-block', width: 'auto' }}
                {...getReactHookFormStatusWithMessage(another)}
              >
                <Input {...field} />
              </Form.Control>
            )}
          />

          <Controller
            control={control}
            name={`${prePath}.${index}.operator`}
            render={({ field }) => (
              <Form.Control showStatusIcon={false} style={{ display: 'inline-block', width: 'auto' }}>
                <Select
                  appearance="button"
                  size="s"
                  matchButtonWidth
                  options={affinityRuleOperatorList}
                  {...field}
                  onChange={operator => {
                    if (needDisabledValue(operator)) {
                      update(index, { value: '', key: watch(`${prePath}.${index}.key`) });
                    }

                    field.onChange(operator);

                    trigger(`${prePath}.${index}.value`);
                  }}
                />
              </Form.Control>
            )}
          />

          <Controller
            control={control}
            name={`${prePath}.${index}.value`}
            render={({ field, ...another }) => (
              <Form.Control
                showStatusIcon={false}
                style={{ display: 'inline-block', width: 'auto' }}
                {...getReactHookFormStatusWithMessage(another)}
              >
                <Input
                  {...field}
                  placeholder={
                    needDisabledValue(watch(`${prePath}.${index}.operator`))
                      ? 'DoesNotExist,Exists操作符不需要填写value'
                      : `多个Label Value请以 ';' 分隔符隔开`
                  }
                  disabled={needDisabledValue(watch(`${prePath}.${index}.operator`))}
                />
              </Form.Control>
            )}
          />

          <Button
            type="icon"
            icon="close"
            disabled={subRules.length <= 1}
            tooltip={subRules.length <= 1 ? '至少保留一条规则' : ''}
            onClick={() => removeSubRule(index)}
          />
        </Form.Item>
      ))}

      <Form.Item>
        <Button
          type="link"
          disabled={disabled}
          tooltip={disabled ? t('请先完成编辑项') : ''}
          onClick={() => appendSubRule({ key: '', operator: NodeAffinityOperatorEnum.In, value: '' })}
        >
          添加规则
        </Button>
      </Form.Item>
    </>
  );
};
