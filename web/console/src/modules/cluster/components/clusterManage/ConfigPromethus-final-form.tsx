import React from 'react';
import { enablePromethus, EnablePromethusParams } from '@/src/webApi/promethus';
import { RootProps } from '../ClusterApp';
import { AntdLayout } from '@src/modules/common/layouts';
import { Button, Form, InputNumber, Checkbox, Input } from 'tea-component';
import { useForm, useField } from 'react-final-form-hooks';

type LocalConfigType = Omit<EnablePromethusParams, 'clusterName'>;

function getStatus(meta, validating) {
  if (meta.active && validating) {
    return 'validating';
  }
  if (!meta.touched) {
    return null;
  }
  return meta.error ? 'error' : 'success';
}

export function ConfigPromethus({ route, actions }: RootProps) {
  async function submit(values) {
    await enablePromethus({ clusterName: route.queries.clusterId, ...values });
    actions.cluster.applyFilter({});
    cancelBack();
  }

  function cancelBack() {
    history.back();
  }

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

  const { form, handleSubmit, validating } = useForm({
    onSubmit: submit,

    initialValues: {
      cpuLimit: 4,
      memoryLimit: 8096,
      cpuRequest: 0.1,
      memoryRequest: 128,
      runOnMaster: false,
      notifyWebhook: '',
      alertRepeatInterval: 20
    },

    validate({ notifyWebhook }) {
      return {
        notifyWebhook: notifyWebhook ? undefined : '不能为空'
      };
    }
  });

  const notifyWebhook = useField('notifyWebhook', form);

  return (
    <AntdLayout
      title="配置告警"
      footer={
        <>
          <Button type="primary" style={{ marginRight: 10 }}>
            提交
          </Button>
          <Button>取消</Button>
        </>
      }
    >
      <Form>
        <Form.Item label="Promethus CPU限制">
          <InputNumber unit="核" min={0} step={0.01} precision={2} />
        </Form.Item>

        <Form.Item label="Promethus CPU预留">
          <InputNumber unit="核" min={0} step={0.01} precision={2} />
        </Form.Item>

        <Form.Item label="Promethus 内存限制">
          <InputNumber unit="Mi" min={4} precision={0} />
        </Form.Item>

        <Form.Item label="Promethus 内存预留">
          <InputNumber unit="Mi" min={4} precision={0} />
        </Form.Item>

        <Form.Item label="Master节点上运行">
          <Checkbox>runOnMaster</Checkbox>
        </Form.Item>

        <Form.Item
          label="指定告警webhook地址"
          status={getStatus(notifyWebhook.meta, validating)}
          message={getStatus(notifyWebhook.meta, validating) === 'error' ? notifyWebhook.meta.error : ''}
        >
          <Input {...notifyWebhook.input} />
        </Form.Item>

        <Form.Item label="重复告警的间隔">
          <InputNumber min={0} precision={0} unit="m" />
        </Form.Item>
      </Form>
    </AntdLayout>
  );
}
