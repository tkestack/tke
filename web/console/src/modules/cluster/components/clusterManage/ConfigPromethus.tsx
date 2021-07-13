import React from 'react';
import { enablePromethus, EnablePromethusParams } from '@/src/webApi/promethus';
import { RootProps } from '../ClusterApp';
import { AntdLayout } from '@src/modules/common/layouts';
import { Button, Form, InputNumber, Checkbox, Input } from 'tea-component';
import { useForm, Controller } from 'react-hook-form';

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
    notifyWebhook: '1234',
    alertRepeatInterval: 20
  });

  const { control, handleSubmit } = useForm<LocalConfigType>({
    defaultValues: initialConfig()
  });

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

        <Controller
          name="notifyWebhook"
          control={control}
          render={() => (
            <Form.Item>
              <Input />
            </Form.Item>
          )}
        />

        <Form.Item label="重复告警的间隔">
          <InputNumber min={0} precision={0} unit="m" />
        </Form.Item>
      </Form>
    </AntdLayout>
  );
}
