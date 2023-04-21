import { getReactHookFormStatusWithMessage } from '@helper';
import { ValidateProvider } from '@src/modules/common/components';
import React from 'react';
import { Controller, useFieldArray, useFormContext } from 'react-hook-form';
import { Button, Input, InputNumber, Select, Table, TableColumn } from 'tea-component';
import {
  NodeAffinityFormType,
  TolerationEffectEnum,
  TolerationOperatorEnum,
  tolerationEffectOptions,
  tolerationOperatorOptions
} from './constants';

export const TolerationRulePanel = () => {
  const { control, watch } = useFormContext<NodeAffinityFormType>();

  const { fields, remove, append } = useFieldArray({
    control,
    name: 'tolerationRules'
  });

  const rulesWatch = watch('tolerationRules');

  const columns: TableColumn[] = [
    {
      key: 'key',
      header: '标签名',
      render(_, __, index) {
        return (
          <Controller
            control={control}
            name={`tolerationRules.${index}.key`}
            render={({ field, ...another }) => (
              <ValidateProvider {...getReactHookFormStatusWithMessage(another)}>
                <Input size="full" {...field} />
              </ValidateProvider>
            )}
          />
        );
      }
    },

    {
      key: 'operator',
      header: '操作符',
      render(_, __, index) {
        return (
          <Controller
            control={control}
            name={`tolerationRules.${index}.operator`}
            render={({ field, ...another }) => (
              <ValidateProvider {...getReactHookFormStatusWithMessage(another)}>
                <Select
                  appearance="button"
                  size="full"
                  matchButtonWidth
                  options={tolerationOperatorOptions}
                  {...field}
                />
              </ValidateProvider>
            )}
          />
        );
      }
    },

    {
      key: 'value',
      header: '标签值',
      render(_, __, index) {
        return (
          <Controller
            control={control}
            name={`tolerationRules.${index}.value`}
            render={({ field, ...another }) => (
              <ValidateProvider {...getReactHookFormStatusWithMessage(another)}>
                <Input
                  size="full"
                  disabled={rulesWatch?.[index]?.operator === TolerationOperatorEnum.Exists}
                  {...field}
                />
              </ValidateProvider>
            )}
          />
        );
      }
    },

    {
      key: 'effect',
      header: '效果',
      render(_, __, index) {
        return (
          <Controller
            control={control}
            name={`tolerationRules.${index}.effect`}
            render={({ field, ...another }) => (
              <ValidateProvider {...getReactHookFormStatusWithMessage(another)}>
                <Select appearance="button" size="full" matchButtonWidth options={tolerationEffectOptions} {...field} />
              </ValidateProvider>
            )}
          />
        );
      }
    },

    {
      key: 'time',
      header: '时间（秒）',
      render(_, __, index) {
        return (
          <Controller
            control={control}
            name={`tolerationRules.${index}.time`}
            render={({ field, ...another }) => (
              <ValidateProvider {...getReactHookFormStatusWithMessage(another)}>
                <InputNumber
                  size="l"
                  hideButton
                  min={0}
                  step={1}
                  disabled={rulesWatch?.[index]?.effect !== TolerationEffectEnum.NoExecute}
                  {...field}
                />
              </ValidateProvider>
            )}
          />
        );
      }
    },

    {
      key: 'action',
      header: '',
      render(_, __, index) {
        return (
          <Button
            type="icon"
            icon="close"
            onClick={() => {
              console.log('remove index', index);
              remove(index);
            }}
          />
        );
      }
    }
  ];

  return (
    <>
      <Table.ActionPanel>
        <Button
          type="link"
          onClick={() =>
            append({
              key: '',
              operator: TolerationOperatorEnum.Exists,
              value: '',
              effect: TolerationEffectEnum.All,
              time: 0
            })
          }
        >
          添加
        </Button>
      </Table.ActionPanel>
      <Table columns={columns} records={fields} recordKey="id" />
    </>
  );
};
