import React from 'react';
import { enablePromethus, EnablePromethusParams } from '@/src/webApi/promethus';
import { RootProps } from '../ClusterApp';
import { AntdLayout } from '@src/modules/common/layouts';
import { Button, Form, InputNumber, Checkbox, Input } from 'tea-component';
import { useForm, Controller } from 'react-hook-form';
import validatorjs from 'validator';
import { getReactHookFormStatusWithMessage } from '@helper';

type LocalConfigType = Omit<EnablePromethusParams, 'clusterName'>;

export function ConfigPromethus({ route, actions }: RootProps) {
  const initialConfig = (): LocalConfigType => ({
    resources: {
      limits: {
        cpu: 4,
        memory: 8096
      },
      requests: {
        cpu: 0.1,
        memory: 128
      }
    },
    runOnMaster: false,
    notifyWebhook: '',
    alertRepeatInterval: 20
  });

  const {
    control,
    handleSubmit,

    getValues
  } = useForm<LocalConfigType>({
    defaultValues: initialConfig(),
    mode: 'onBlur'
  });

  const limitValidate = (type: 'cpu' | 'memory') => () => {
    const {
      resources: { limits, requests }
    } = getValues();

    if (requests[type] > limits[type]) return `${type}预留不能超过限制`;
  };

  async function onSubmit(values) {
    await enablePromethus({ clusterName: route.queries.clusterId, ...values });
    actions.cluster.applyFilter({});
    cancelBack();
  }

  function cancelBack() {
    history.back();
  }

  return (
    <AntdLayout
      title="配置告警"
      footer={
        <>
          <Button type="primary" style={{ marginRight: 10 }} onClick={handleSubmit(onSubmit)}>
            提交
          </Button>
          <Button onClick={cancelBack}>取消</Button>
        </>
      }
    >
      <Form>
        <Controller
          control={control}
          name="resources.limits.cpu"
          rules={{
            validate: limitValidate('cpu')
          }}
          render={({ field, ...others }) => (
            <Form.Item label="Promethus CPU限制" {...getReactHookFormStatusWithMessage(others)}>
              <InputNumber {...field} unit="核" min={0} step={0.01} precision={2} size="l" />
            </Form.Item>
          )}
        />

        <Controller
          control={control}
          name="resources.requests.cpu"
          rules={{
            validate: limitValidate('cpu')
          }}
          render={({ field, ...others }) => (
            <Form.Item label="Promethus CPU预留" {...getReactHookFormStatusWithMessage(others)}>
              <InputNumber {...field} unit="核" min={0} step={0.01} precision={2} size="l" />
            </Form.Item>
          )}
        />

        <Controller
          control={control}
          name="resources.limits.memory"
          rules={{
            validate: limitValidate('memory')
          }}
          render={({ field, ...others }) => (
            <Form.Item label="Promethus 内存限制" {...getReactHookFormStatusWithMessage(others)}>
              <InputNumber {...field} unit="Mi" min={4} precision={0} size="l" />
            </Form.Item>
          )}
        />

        <Controller
          control={control}
          name="resources.requests.memory"
          rules={{
            validate: limitValidate('memory')
          }}
          render={({ field, ...others }) => (
            <Form.Item label="Promethus 内存预留" {...getReactHookFormStatusWithMessage(others)}>
              <InputNumber {...field} unit="Mi" min={4} precision={0} size="l" />
            </Form.Item>
          )}
        />

        <Controller
          control={control}
          name="runOnMaster"
          render={({ field }) => (
            <Form.Item label="Master节点上运行">
              <Checkbox {...field}>runOnMaster</Checkbox>
            </Form.Item>
          )}
        />

        <Controller
          name="notifyWebhook"
          control={control}
          rules={{
            validate(value) {
              if (value && !validatorjs.isURL(value)) {
                return 'webhook 格式不正确!';
              }
            }
          }}
          render={({ field, ...others }) => (
            <Form.Item label="指定告警webhook地址" {...getReactHookFormStatusWithMessage(others)}>
              <Input {...field} />
            </Form.Item>
          )}
        />

        <Controller
          control={control}
          name="alertRepeatInterval"
          render={({ field }) => (
            <Form.Item label="重复告警的间隔">
              <InputNumber {...field} min={0} precision={0} unit="m" size="l" />
            </Form.Item>
          )}
        />
      </Form>
    </AntdLayout>
  );
}
