import React from 'react';
import { Space, Form, InputNumber, Button, Checkbox, Input } from 'antd';
import { enablePromethus, EnablePromethusParams } from '@/src/webApi/promethus';
import { RootProps } from '../ClusterApp';
import { AntdLayout } from '@src/modules/common/layouts';

type LocalConfigType = Omit<EnablePromethusParams, 'clusterName'>;

export function ConfigPromethus({ route, actions }: RootProps) {
  const inputStyle = { width: '300px' };

  const initialConfig: LocalConfigType = {
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
  };

  async function submit(values) {
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
        <Space>
          <Button type="primary" htmlType="submit" form="promethusConfigForm">
            提交
          </Button>
          <Button onClick={cancelBack}>取消</Button>
        </Space>
      }
    >
      <Form
        labelAlign="left"
        labelCol={{ span: 3 }}
        size="middle"
        validateTrigger="onBlur"
        initialValues={initialConfig}
        onFinish={submit}
        id="promethusConfigForm"
      >
        <Form.Item label="Promethus CPU限制">
          <Space>
            <Form.Item noStyle name={['resources', 'limits', 'cpu']} rules={[{ type: 'number', min: 0 }]}>
              <InputNumber style={inputStyle} min={0} />
            </Form.Item>
            核
          </Space>
        </Form.Item>
        <Form.Item label="Promethus CPU预留">
          <Space>
            <Form.Item noStyle name={['resources', 'requests', 'cpu']} rules={[{ type: 'number', min: 0 }]}>
              <InputNumber style={inputStyle} min={0} />
            </Form.Item>
            核
          </Space>
        </Form.Item>
        <Form.Item label="Promethus 内存限制">
          <Space>
            <Form.Item noStyle name={['resources', 'limits', 'memory']} rules={[{ type: 'number', min: 4 }]}>
              <InputNumber style={inputStyle} min={4} />
            </Form.Item>
            Mi
          </Space>
        </Form.Item>
        <Form.Item label="Promethus 内存预留">
          <Space>
            <Form.Item noStyle name={['resources', 'requests', 'memory']} rules={[{ type: 'number', min: 4 }]}>
              <InputNumber style={inputStyle} min={4} />
            </Form.Item>
            Mi
          </Space>
        </Form.Item>
        <Form.Item label="Master节点上运行">
          <Space>
            <Form.Item noStyle name={['runOnMaster']} valuePropName="checked">
              <Checkbox />
            </Form.Item>
            runOnMaster
          </Space>
        </Form.Item>
        <Form.Item label="指定告警webhook地址" name="notifyWebhook" rules={[{ type: 'url' }]}>
          <Input style={inputStyle} />
        </Form.Item>
        <Form.Item label="重复告警的间隔">
          <Space>
            <Form.Item noStyle name={['alertRepeatInterval']} rules={[{ type: 'number', min: 0 }]}>
              <InputNumber style={inputStyle} min={0} />
            </Form.Item>
            m
          </Space>
        </Form.Item>
      </Form>
    </AntdLayout>
  );
}
