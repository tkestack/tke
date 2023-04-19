import { t } from '@/tencent/tea-app/lib/i18n';
import { getParamByUrl, getReactHookFormStatusWithMessage } from '@helper';
import { zodResolver } from '@hookform/resolvers/zod';
import { nodeApi } from '@src/webApi';
import { useRequest } from 'ahooks';
import React, { useEffect, useState } from 'react';
import { Control, Controller, UseFormTrigger, useFieldArray, useForm, useFormState, useWatch } from 'react-hook-form';
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
  NodeAffinityOperatorEnum,
  RuleType,
  affinityRuleOperatorList,
  generateDefaultRules,
  ruleSchema
} from './constants';

const SubRulePanel = ({
  control,
  ruleIndex,
  trigger
}: {
  control: Control<RuleType>;
  ruleIndex: number;
  trigger: UseFormTrigger<RuleType>;
}) => {
  const {
    fields: subRules,
    append: appendSubRule,
    remove: removeSubRule,
    update
  } = useFieldArray({
    control,
    name: `rules.${ruleIndex}.subRules`
  });

  const watchSubRules = useWatch({ control, name: `rules.${ruleIndex}.subRules` });

  function needDisabledValue(operator) {
    return operator === NodeAffinityOperatorEnum.Exists || operator === NodeAffinityOperatorEnum.DoesNotExist;
  }

  const { errors } = useFormState({ control });
  console.log('SubRulePanel errors', errors);
  const disabled = Boolean(errors?.rules?.[ruleIndex]);

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
                      update(index, { value: '', key: watchSubRules?.[index].key });
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
                    needDisabledValue(watchSubRules?.[index]?.operator)
                      ? 'DoesNotExist,Exists操作符不需要填写value'
                      : `多个Label Value请以 ';' 分隔符隔开`
                  }
                  disabled={needDisabledValue(watchSubRules?.[index]?.operator)}
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

export const AffinityRulePanel = ({ showWeight = true, submiting, onSubmit, defaultRules = [] }) => {
  const hasModal = !showWeight;

  const {
    control,
    formState: { errors },
    handleSubmit,
    trigger
  } = useForm<RuleType>({
    mode: 'onBlur',

    defaultValues: {
      rules: defaultRules
    },
    resolver: zodResolver(ruleSchema)
  });

  const {
    fields: rules,
    append,
    remove
  } = useFieldArray({
    control,
    name: 'rules'
  });

  console.log('errors', errors);

  useEffect(() => {
    if (!submiting) return;

    handleSubmit(
      data => {
        console.log('submiting', data);

        onSubmit(data?.rules ?? []);
      },
      () => onSubmit()
    )();
  }, [submiting, handleSubmit, onSubmit]);

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

          <SubRulePanel control={control} trigger={trigger} ruleIndex={index} />
        </Form>
      ))}

      <AddRuleButton append={append} hasModal={hasModal} />
    </>
  );
};

function AddRuleButton({ append, hasModal }) {
  const [visible, setVisible] = useState(false);
  const [submiting, setSubmiting] = useState(false);
  const [type, setType] = useState('scheduleByNode');

  function handleClick() {
    if (!hasModal) {
      append({ weight: 1, subRules: [{ key: '', operator: NodeAffinityOperatorEnum.In, value: '' }] });
    } else {
      setVisible(true);
    }
  }

  function handleOk() {
    setSubmiting(true);
  }

  function onSubmit(data) {
    if (data) {
      console.log('append data', data);
      append(data);

      setVisible(false);
    }

    setSubmiting(false);
  }

  return (
    <>
      <Button type="link" onClick={handleClick}>
        添加条件
      </Button>

      {hasModal && (
        <Modal caption="编辑强制满足条件" size="xl" visible={visible} onClose={() => setVisible(false)}>
          <Modal.Body>
            <Form>
              <Form.Item label="条件选择">
                <RadioGroup value={type} onChange={value => setType(value)}>
                  <Radio name="scheduleByNode">指定节点调度</Radio>
                  <Radio name="customLabel">自定义Label规则</Radio>
                </RadioGroup>
              </Form.Item>

              {type === 'customLabel' && (
                <Form.Item>
                  <AffinityRulePanel
                    showWeight={false}
                    submiting={submiting}
                    onSubmit={onSubmit}
                    defaultRules={generateDefaultRules()}
                  />
                </Form.Item>
              )}

              {type === 'scheduleByNode' && (
                <Form.Item>
                  <NodeTransferPanel submiting={submiting} onSubmit={onSubmit} />
                </Form.Item>
              )}
            </Form>
          </Modal.Body>

          <Modal.Footer>
            <Button type="primary" onClick={handleOk}>
              确定
            </Button>

            <Button onClick={() => setVisible(false)}>取消</Button>
          </Modal.Footer>
        </Modal>
      )}
    </>
  );
}

const { scrollable, selectable, removeable } = Table.addons;

function NodeTransferPanel({ submiting, onSubmit }) {
  const clusterId = getParamByUrl('clusterId')!;

  const [selectedKeys, setSelectedKeys] = useState([]);
  const [errorMessage, setErrorMessage] = useState('');

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

  useEffect(() => {
    if (selectedKeys.length > 0) {
      setErrorMessage('');
    }

    if (!submiting) return;

    if (selectedKeys.length < 1) {
      setErrorMessage('节点不能为空，请选择节点');
      onSubmit();

      return;
    }

    onSubmit([
      {
        weight: 1,
        subRules: [
          {
            key: 'kubernetes.io/hostname',
            operator: NodeAffinityOperatorEnum.In,
            value: selectedKeys.join(';')
          }
        ]
      }
    ]);
  }, [submiting, onSubmit, selectedKeys]);

  return (
    <>
      {errorMessage && <Alert type="error">{errorMessage}</Alert>}
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
                  value: selectedKeys,
                  onChange: keys => setSelectedKeys(keys),
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
              records={nodeList.filter(item => selectedKeys.includes(item?.metadata?.name))}
              recordKey="metadata.name"
              addons={[
                removeable({
                  onRemove: key => {
                    console.log('remove--->', key);

                    setSelectedKeys(keys => keys.filter(item => item !== key));
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
