import React, { useEffect, useState } from 'react';
import { Space, Button, Form, Select, Checkbox, InputNumber, Typography } from 'antd';
import { AntdLayout } from '@src/modules/common/layouts';
import { getK8sValidVersions } from '@src/webApi/cluster';
import { compare } from 'compare-versions';
import { RootProps } from '../ClusterApp';
import { updateCluster } from '@src/webApi/cluster';

export function ClusterUpdate({ route, actions }: RootProps) {
  const ItemStyle = () => ({
    width: 120
  });

  const defaultUpgradeConfig = {
    version: null,
    drainNodeBeforeUpgrade: true,
    maxUnready: 20,
    autoMode: true
  };

  const [showMaxPodUnready, setShowMaxPodUnready] = useState(true);

  const { clusterId, clusterVersion } = route.queries;
  const [_, clusterVersionSecondPart] = clusterVersion.split('.');

  function goBack() {
    history.back();
  }

  async function perform(values) {
    await updateCluster({ ...values, clusterName: clusterId });
    actions.cluster.applyPolling();
    goBack();
  }

  console.log(React.version);

  const [k8sValidVersions, setK8sValidVersions] = useState([]);

  useEffect(() => {
    async function fetchK8sValidVersions() {
      const versions = await getK8sValidVersions();
      setK8sValidVersions(versions);
    }

    fetchK8sValidVersions();
  }, []);

  return (
    <AntdLayout
      title="升级Master"
      footer={
        <Space>
          <Button type="primary" htmlType="submit" form="clusterUpdateConfigForm">
            提交
          </Button>
          <Button onClick={goBack}>取消</Button>
        </Space>
      }
    >
      <Typography.Title level={5} style={{ marginBottom: 35 }}>
        更新集群的K8S版本{clusterVersion}至：
      </Typography.Title>

      <Form
        id="clusterUpdateConfigForm"
        labelAlign="left"
        labelCol={{ span: 3 }}
        size="middle"
        validateTrigger="onBlur"
        initialValues={defaultUpgradeConfig}
        onFinish={perform}
      >
        <Form.Item
          label="升级目标版本"
          extra="注意：master升级支持一个次版本升级到下一个次版本，或者同样次版本的补丁版。"
          name={['version']}
          rules={[
            { required: true },
            {
              validator(_, value: string) {
                const [__, targetSecond] = value.split('.');
                return +targetSecond - +clusterVersionSecondPart === 1
                  ? Promise.resolve()
                  : Promise.reject('不支持直接升级到该版本！');
              }
            }
          ]}
        >
          <Select style={ItemStyle()}>
            {k8sValidVersions.map(v => (
              <Select.Option key={v} disabled={compare(clusterVersion, v, '>=')} value={v}>
                {v}
              </Select.Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          valuePropName="checked"
          name={['autoMode']}
          label="自动升级Worker"
          extra="注意：启用自定升级Worker，在升级完Master后，将自动升级集群下所有Worker节点。"
        >
          <Checkbox onChange={e => setShowMaxPodUnready(e.target.checked)}>启用自动升级</Checkbox>
        </Form.Item>

        {showMaxPodUnready && (
          <>
            <Form.Item
              valuePropName="checked"
              name={['drainNodeBeforeUpgrade']}
              label="驱逐节点"
              extra="若选择升级前驱逐节点，该节点所有pod将在升级前被驱逐，此时节点如有pod使用emptyDir类卷会导致驱逐失败而影响升级流程"
            >
              <Checkbox>驱逐节点</Checkbox>
            </Form.Item>

            <Form.Item
              label="最大不可用Pod占比"
              extra="注意如果节点过少，而设置比例过低，没有足够多的节点承载pod的迁移会导致升级卡死。如果业务对pod可用比例较高，请考虑选择升级前不驱逐节点。"
            >
              <Space>
                <Form.Item name={['maxUnready']} noStyle rules={[{ type: 'number', required: true, min: 0, max: 100 }]}>
                  <InputNumber style={ItemStyle()} min={0} max={100} />
                </Form.Item>
                %
              </Space>
            </Form.Item>
          </>
        )}
      </Form>
    </AntdLayout>
  );
}
