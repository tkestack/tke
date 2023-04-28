import { t } from '@/tencent/tea-app/lib/i18n';
import { getParamByUrl, getReactHookFormStatusWithMessage } from '@helper';
import { zodResolver } from '@hookform/resolvers/zod';
import { nodeApi } from '@src/webApi';
import { useRequest } from 'ahooks';
import React, { useState } from 'react';
import { Controller, FormProvider, useFieldArray, useForm, useFormContext, useFormState } from 'react-hook-form';
import {
  Alert,
  Button,
  Form,
  Input,
  InputNumber,
  Justify,
  Modal,
  Radio,
  RadioGroup,
  Select,
  Table,
  TableColumn,
  Transfer
} from 'tea-component';
import {
  AppendAffinityRuleFormType,
  NodeAffinityOperatorEnum,
  ScheduleTypeEnum,
  affinityRuleOperatorList,
  appendAffinityRuleSchema,
  defaultAppendAffinityRuleFormData
} from './constants';

export function AddRuleButton({ append }) {
  const [visible, setVisible] = useState(false);

  const [type, setType] = useState(ScheduleTypeEnum.ScheduleByLabel);

  const [nodeKeys, setNodeKeys] = useState([]);

  const formProps = useForm<AppendAffinityRuleFormType>({
    mode: 'onBlur',
    defaultValues: defaultAppendAffinityRuleFormData,
    resolver: zodResolver(appendAffinityRuleSchema)
  });

  const { handleSubmit, reset } = formProps;

  function onCancel() {
    setVisible(false);
    reset();
  }

  function handleOk() {
    if (type === ScheduleTypeEnum.ScheduleByNode && nodeKeys.length >= 1) {
      append({
        weight: 1,
        subRules: nodeKeys.map(key => ({
          key: 'kubernetes.io/hostname',
          operator: NodeAffinityOperatorEnum.In,
          value: key
        }))
      });

      onCancel();
    }

    if (type === ScheduleTypeEnum.ScheduleByLabel) {
      handleSubmit(({ rules }) => {
        append(rules);

        onCancel();
      })();
    }
  }

  return (
    <>
      <Button type="link" onClick={() => setVisible(true)}>
        添加条件
      </Button>

      <Modal caption="编辑强制满足条件" size="xl" visible={visible} onClose={() => setVisible(false)}>
        <Modal.Body>
          <FormProvider {...formProps}>
            <Form>
              <Form.Item label="条件选择">
                <RadioGroup value={type} onChange={(type: ScheduleTypeEnum) => setType(type)}>
                  <Radio name={ScheduleTypeEnum.ScheduleByNode}>指定节点调度</Radio>
                  <Radio name={ScheduleTypeEnum.ScheduleByLabel}>自定义Label规则</Radio>
                </RadioGroup>
              </Form.Item>

              {type === ScheduleTypeEnum.ScheduleByLabel && (
                <Form.Item>
                  <AffinityRulePanel />
                </Form.Item>
              )}

              {type === ScheduleTypeEnum.ScheduleByNode && (
                <Form.Item>
                  <NodeTransferPanel nodeKeys={nodeKeys} setNodeKeys={setNodeKeys} />
                </Form.Item>
              )}
            </Form>
          </FormProvider>
        </Modal.Body>

        <Modal.Footer>
          <Button type="primary" onClick={handleOk}>
            确定
          </Button>

          <Button onClick={onCancel}>取消</Button>
        </Modal.Footer>
      </Modal>
    </>
  );
}

const { scrollable, selectable, removeable } = Table.addons;

function NodeTransferPanel({ nodeKeys, setNodeKeys }) {
  const clusterId = getParamByUrl('clusterId')!;

  const { data: nodeList = [] } = useRequest(
    async () => {
      const rsp = await nodeApi.fetchNodeList({ clusterId });

      console.log('NodeTransferPanel--->', rsp);

      return rsp?.items ?? [];
    },
    { ready: Boolean(clusterId) }
  );

  const columns: TableColumn[] = [
    {
      key: 'metadata.name',
      header: 'ID/节点名'
    },

    {
      key: 'ip',
      header: 'IP地址',
      render: ({ metadata }) => metadata?.name
    }
  ];

  return (
    <>
      {nodeKeys.length < 1 && <Alert type="error">节点不能为空，请选择节点</Alert>}
      <Transfer
        leftCell={
          <Transfer.Cell title={`当前集群下有以下可用节点 (共3项 已加载 3 项) 已选择 0 项`}>
            <Table
              columns={columns}
              records={nodeList}
              recordKey="metadata.name"
              addons={[
                scrollable({
                  maxHeight: 310
                }),

                selectable({
                  value: nodeKeys,
                  onChange: keys => setNodeKeys(keys),
                  rowSelect: true
                })
              ]}
            />
          </Transfer.Cell>
        }
        rightCell={
          <Transfer.Cell title="">
            <Table
              columns={columns}
              records={nodeList.filter(item => nodeKeys.includes(item?.metadata?.name))}
              recordKey="metadata.name"
              addons={[
                removeable({
                  onRemove: key => {
                    console.log('remove--->', key);

                    setNodeKeys(items => items.filter(item => item !== key));
                  }
                })
              ]}
            />
          </Transfer.Cell>
        }
      />
    </>
  );
}

export function AffinityRulePanel() {
  const showWeight = false;

  const { control } = useFormContext<AppendAffinityRuleFormType>();

  const {
    fields: rules,
    append,
    remove
  } = useFieldArray({
    control,
    name: `rules`
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
              name={`rules.${index}.weight`}
              render={({ field, ...another }) => (
                <Form.Item label="权重" showStatusIcon={false} {...getReactHookFormStatusWithMessage(another)}>
                  <InputNumber hideButton size="l" min={1} max={100} step={1} {...field} />
                </Form.Item>
              )}
            />
          )}

          <SubRulePanel ruleIndex={index} />
        </Form>
      ))}

      <Button
        type="link"
        onClick={() => append({ weight: 1, subRules: [{ key: '', operator: NodeAffinityOperatorEnum.In, value: '' }] })}
      >
        添加条件
      </Button>
    </>
  );
}

const SubRulePanel = ({ ruleIndex }: { ruleIndex: number }) => {
  const { control, trigger, watch } = useFormContext<AppendAffinityRuleFormType>();

  const {
    fields: subRules,
    append: appendSubRule,
    remove: removeSubRule,
    update
  } = useFieldArray({
    control,
    name: `rules.${ruleIndex}.subRules`
  });

  function needDisabledValue(operator) {
    return operator === NodeAffinityOperatorEnum.Exists || operator === NodeAffinityOperatorEnum.DoesNotExist;
  }

  const { errors } = useFormState({ control });

  const disabled = Boolean(errors?.rules?.[ruleIndex]?.subRules);

  return (
    <>
      {subRules.map(({ id }, index) => (
        <Form.Item label="条件" key={id}>
          <Controller
            control={control}
            name={`rules.${ruleIndex}.subRules.${index}.key`}
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
            name={`rules.${ruleIndex}.subRules.${index}.operator`}
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
                      update(index, { value: '', key: watch(`rules.${ruleIndex}.subRules.${index}.key`) });
                    }

                    field.onChange(operator);

                    trigger(`rules.${ruleIndex}.subRules.${index}.value`);
                  }}
                />
              </Form.Control>
            )}
          />

          <Controller
            control={control}
            name={`rules.${ruleIndex}.subRules.${index}.value`}
            render={({ field, ...another }) => (
              <Form.Control
                showStatusIcon={false}
                style={{ display: 'inline-block', width: 'auto' }}
                {...getReactHookFormStatusWithMessage(another)}
              >
                <Input
                  {...field}
                  placeholder={
                    needDisabledValue(watch(`rules.${ruleIndex}.subRules.${index}.operator`))
                      ? 'DoesNotExist,Exists操作符不需要填写value'
                      : `多个Label Value请以 ';' 分隔符隔开`
                  }
                  disabled={needDisabledValue(watch(`rules.${ruleIndex}.subRules.${index}.operator`))}
                />
              </Form.Control>
            )}
          />

          <Button type="icon" icon="close" onClick={() => removeSubRule(index)} />
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

// TODO 泛型组件
