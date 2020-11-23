import React, { useEffect, useState } from 'react';
import { Space, Button, Form, Select, Checkbox, InputNumber, Typography } from 'antd';
import { AntdLayout } from '@src/modules/common/layouts';
import { getK8sValidVersions } from '@src/webApi/cluster';
import compareVersion from 'compare-versions';

export function ClusterUpdate() {
  const ItemStyle = () => ({
    width: 120
  });

  console.log(React.version);

  const [k8sValidVersions, setK8sValidVersions] = useState([]);

  useEffect(() => {
    async function fetchK8sValidVersions() {
      const { k8sValidVersions } = await getK8sValidVersions();
      setK8sValidVersions(k8sValidVersions);
    }

    fetchK8sValidVersions();
  }, []);

  return (
    <AntdLayout
      title="升级Master"
      footer={
        <Space>
          <Button type="primary" htmlType="submit" form="promethusConfigForm">
            提交
          </Button>
          <Button>取消</Button>
        </Space>
      }
    >
      <Typography.Title level={5} style={{ marginBottom: 35 }}>
        更新集群的K8S版本1.14.3至：
      </Typography.Title>

      <Form labelAlign="left" labelCol={{ span: 3 }} size="middle">
        <Form.Item
          label="升级目标版本"
          extra="注意：master升级支持一个此版本升级到下一个次版本，或者同样次版本的补丁版。"
        >
          <Select style={ItemStyle()}>
            {k8sValidVersions.map(v => (
              <Select.Option value={v}>{v}</Select.Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          label="自动升级Worker"
          extra="注意：启用自定升级Worker，在升级完Master后，将自动升级集群下所有Worker节点。"
        >
          <Checkbox>启用自动升级</Checkbox>
        </Form.Item>
        <Form.Item label="最大不可用Pod占比" extra="升级过程中不可以Pod数超过该占比将暂停升级">
          <Space>
            <InputNumber style={ItemStyle()} min={0} max={100} /> %
          </Space>
        </Form.Item>
      </Form>
    </AntdLayout>
  );
}
