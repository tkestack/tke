import React, { useEffect, useState } from 'react';
import { Space, Button, Form, Select, Checkbox, InputNumber, Typography, Transfer, Table, Alert } from 'antd';
import { AntdLayout } from '@src/modules/common/layouts';
import { TransferProps } from 'antd/lib/transfer';
import { difference } from 'lodash';
import { RootProps } from '../ClusterApp';
import { getNodes, updateWorkers } from '@src/webApi/cluster';

export function WorkerUpdate({ route }: RootProps) {
  const ItemStyle = () => ({
    width: 120
  });

  const [nodes, setNodes] = useState([]);
  const [targetKeys, setTargatKeys] = useState([]);
  const [maxUnready, setMaxUnready] = useState(20);
  const [drainNodeBeforeUpgrade, setDrainNodeBeforeUpgrade] = useState(true);

  const { clusterId, clusterVersion } = route.queries;

  useEffect(() => {
    (async function () {
      console.log('start:---');
      const nodes = await getNodes({ clusterName: clusterId, clusterVersion });
      console.log('getMachines:', nodes);
      setNodes(nodes);
    })();
  }, [clusterId, clusterVersion]);

  const transferColumns = [
    {
      dataIndex: 'name',
      title: 'ID/名称'
    },

    {
      dataIndex: 'labels',
      title: 'label'
    },

    {
      dataIndex: 'kubeletVersion',
      title: 'Kubernetes版本'
    },

    {
      dataIndex: 'clusterVersion',
      title: '目标Kubernetes版本'
    }
  ];

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
        <Space>
          <Button type="primary" disabled={targetKeys.length <= 0} onClick={submit}>
            确定
          </Button>
          <Button onClick={goback}>取消</Button>
        </Space>
      }
    >
      <Form labelAlign="left" labelCol={{ span: 3 }} size="middle">
        <Form.Item label=" 升级说明">
          当前所选集群Master版本为{clusterVersion}，您可为您的节点Kubernetes版本升级到当前的最新版本。
        </Form.Item>

        <Form.Item label="选择节点">
          <TableTransfer
            titles={[`当前集群下有以下可升级节点`, `已选择${targetKeys.length}项`]}
            columns={transferColumns}
            dataSource={nodes}
            targetKeys={targetKeys}
            listStyle={{}}
            onChange={targetKeys => setTargatKeys(targetKeys)}
          />
        </Form.Item>

        <Form.Item
          label="驱逐节点"
          extra="若选择升级前驱逐节点，该节点所有pod将在升级前被驱逐，此时节点如有pod使用emptyDir类卷会导致驱逐失败而影响升级流程"
        >
          <Checkbox defaultChecked={drainNodeBeforeUpgrade} onChange={e => setDrainNodeBeforeUpgrade(e.target.checked)}>
            驱逐节点
          </Checkbox>
        </Form.Item>

        <Form.Item label="最大不可用Pod占比" extra="升级过程中不可用Pod数超过该占比将暂停升级">
          <Space>
            <InputNumber
              style={ItemStyle()}
              min={0}
              max={100}
              defaultValue={maxUnready}
              onChange={value => setMaxUnready(+value)}
            />
            %
          </Space>
        </Form.Item>
      </Form>
    </AntdLayout>
  );
}

interface TableTransferProps extends TransferProps<any> {
  columns: Array<any>;
}

function TableTransfer({ columns, ...restProps }: TableTransferProps) {
  return (
    <Transfer {...restProps} showSelectAll={false}>
      {({ filteredItems, onItemSelectAll, onItemSelect, selectedKeys: listSelectedKeys, disabled: listDisabled }) => {
        const rowSelection = {
          getCheckboxProps: item => ({ disabled: listDisabled || item.disabled }),
          onSelectAll(selected, selectedRows) {
            const treeSelectedKeys = selectedRows.filter(item => !item.disabled).map(({ key }) => key);
            const diffKeys = selected
              ? difference(treeSelectedKeys, listSelectedKeys)
              : difference(listSelectedKeys, treeSelectedKeys);
            onItemSelectAll(diffKeys, selected);
          },
          onSelect({ key }, selected) {
            onItemSelect(key, selected);
          },
          selectedRowKeys: listSelectedKeys
        };

        return (
          <Table
            rowSelection={rowSelection}
            columns={columns}
            dataSource={filteredItems}
            size="small"
            pagination={false}
            style={{ pointerEvents: listDisabled ? 'none' : null }}
            onRow={({ key, disabled: itemDisabled }) => ({
              onClick: () => {
                if (itemDisabled || listDisabled) return;
                onItemSelect(key, !listSelectedKeys.includes(key));
              }
            })}
          />
        );
      }}
    </Transfer>
  );
}
