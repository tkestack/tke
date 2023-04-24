import { useFetch } from '@src/modules/common/hooks/useFetch';
import { virtualMachineAPI } from '@src/webApi';
import dayjs from 'dayjs';
import React, { useState } from 'react';
import { Alert, Bubble, Button, Justify, Switch, Table, TableColumn, Text } from 'tea-component';
import { v4 as uuidv4 } from 'uuid';

const { autotip } = Table.addons;

export const VMEventPanel = ({ clusterId, namespace, name }) => {
  const [polling, setPolling] = useState(true);

  const { data, reFetch, status } = useFetch(
    async () => {
      const { items } = await virtualMachineAPI.fetchEventList({ clusterId, namespace, name });

      return { data: items.map(_ => ({ ..._, id: uuidv4() })).reverse() };
    },
    [clusterId, namespace, name],
    {
      fetchAble: !!(clusterId && namespace && name),
      polling,
      needClearData: false
    }
  );

  const columns: TableColumn[] = [
    {
      key: 'firstTimestamp',
      header: '首次出现时间',
      render({ firstTimestamp }) {
        return firstTimestamp ? dayjs(firstTimestamp).format('YYYY-MM-DD HH:mm:ss') : '-';
      }
    },

    {
      key: 'lastTimestamp',
      header: '最后出现时间',
      render({ lastTimestamp }) {
        return lastTimestamp ? dayjs(lastTimestamp).format('YYYY-MM-DD HH:mm:ss') : '-';
      }
    },

    {
      key: 'type',
      header: '级别'
    },

    {
      key: 'involvedObject.kind',
      header: '资源类型'
    },

    {
      key: 'involvedObject.name',
      header: '资源名称'
    },

    {
      key: 'reason',
      header: '内容'
    },

    {
      key: 'message',
      header: '详细描述',
      render({ message }) {
        return (
          <Bubble content={message}>
            <Text overflow>{message}</Text>
          </Bubble>
        );
      }
    },
    {
      key: 'count',
      header: '出现次数'
    }
  ];

  return (
    <>
      <Alert>资源事件只保存最近1小时内发生的事件，请尽快查阅。</Alert>

      <Table.ActionPanel>
        <Justify
          left={<Button icon="refresh" onClick={reFetch} />}
          right={
            <Switch value={polling} onChange={value => setPolling(value)}>
              自动刷新
            </Switch>
          }
        />
      </Table.ActionPanel>

      <Table
        columns={columns}
        records={data || []}
        recordKey="id"
        addons={[
          autotip({
            isLoading: status === 'loading'
          })
        ]}
      />
    </>
  );
};
