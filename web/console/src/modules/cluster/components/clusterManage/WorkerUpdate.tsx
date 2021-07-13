import React, { useEffect, useState } from 'react';
import { AntdLayout } from '@src/modules/common/layouts';
import { RootProps } from '../ClusterApp';
import { getNodes, updateWorkers } from '@src/webApi/cluster';
import { Button, Form, Text, Checkbox, InputNumber, Transfer, Table } from 'tea-component';

const { selectable, removeable, scrollable } = Table.addons;

export function WorkerUpdate({ route }: RootProps) {
  const { clusterId, clusterVersion } = route.queries;

  const [nodes, setNodes] = useState([]);
  const [targetKeys, setTargetKeys] = useState([]);
  const [maxUnready, setMaxUnready] = useState(20);
  const [drainNodeBeforeUpgrade, setDrainNodeBeforeUpgrade] = useState(true);

  useEffect(() => {
    (async function () {
      console.log('start:---');
      const nodes = await getNodes({ clusterName: clusterId, clusterVersion });
      console.log('getMachines:', nodes);
      setNodes(nodes.map(node => ({ ...node, disabled: node.phase !== 'Running' })));
    })();
  }, [clusterId, clusterVersion]);

  async function submit() {
    const mchineNames = [
      ...new Set(
        nodes
          .filter(({ key }) => targetKeys.includes(key))
          .map(({ machines }) => machines as string[])
          .reduce((all, current) => [...all, ...current], [])
      )
    ];

    await updateWorkers({
      mchineNames,
      maxUnready,
      drainNodeBeforeUpgrade,
      clusterName: clusterId
    });
    goback();
  }

  function goback() {
    history.back();
  }

  return (
    <AntdLayout
      title="升级Worker"
      footer={
        <>
          <Button type="primary" style={{ marginRight: 10 }} disabled={targetKeys.length <= 0} onClick={submit}>
            确定
          </Button>
          <Button onClick={goback}>取消</Button>
        </>
      }
    >
      <Form>
        <Form.Item label="升级说明">
          <Text reset>
            当前所选集群Master版本为{clusterVersion}，您可为您的节点Kubernetes版本升级到当前的最新版本。
          </Text>
        </Form.Item>

        <Form.Item label="选择节点">
          <Transfer
            leftCell={
              <Transfer.Cell title="当前集群下有以下可升级节点">
                <SourceTable dataSource={nodes} targetKeys={targetKeys} onChange={setTargetKeys} />
              </Transfer.Cell>
            }
            rightCell={
              <Transfer.Cell title={`已选择${targetKeys.length}项`}>
                <TargetTable
                  dataSource={nodes.filter(({ key }) => targetKeys.includes(key))}
                  onRemove={key => setTargetKeys(pre => pre.filter(k => k !== key))}
                />
              </Transfer.Cell>
            }
          />
        </Form.Item>

        <Form.Item
          label="驱逐节点"
          extra="若选择升级前驱逐节点，该节点所有pod将在升级前被驱逐，此时节点如有pod使用emptyDir类卷会导致驱逐失败而影响升级流程"
        >
          <Checkbox value={drainNodeBeforeUpgrade} onChange={setDrainNodeBeforeUpgrade}>
            驱逐节点
          </Checkbox>
        </Form.Item>

        <Form.Item label="最大不可用Pod占比" extra="升级过程中不可用Pod数超过该占比将暂停升级">
          <InputNumber value={maxUnready} onChange={setMaxUnready} min={0} max={100} />
        </Form.Item>
      </Form>
    </AntdLayout>
  );
}

const transferColumns = [
  {
    key: 'name',
    header: 'ID/名称'
  },

  {
    key: 'phase',
    header: '状态',
    render(phase) {
      return <Text theme={phase === 'Running' ? 'success' : 'danger'}>{phase}</Text>;
    }
  },

  {
    key: 'kubeletVersion',
    header: 'Kubernetes版本'
  },

  {
    key: 'clusterVersion',
    header: '目标Kubernetes版本'
  }
];

function SourceTable({ dataSource, targetKeys, onChange }) {
  return (
    <Table
      records={dataSource}
      columns={transferColumns}
      recordKey="key"
      addons={[
        scrollable({
          maxHeight: 310,
          onScrollBottom: () => console.log('到达底部')
        }),
        selectable({
          value: targetKeys,
          onChange,
          rowSelect: true
        })
      ]}
    />
  );
}

function TargetTable({ dataSource, onRemove }) {
  return <Table records={dataSource} recordKey="name" columns={transferColumns} addons={[removeable({ onRemove })]} />;
}
