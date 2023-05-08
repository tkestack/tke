import React from 'react';
import { Table, TableColumn, Justify, SearchBox } from 'tea-component';

export const SnapshotTablePanel = () => {
  const columns: TableColumn[] = [
    {
      key: 'name',
      header: '快照名称'
    },

    {
      key: 'status',
      header: '状态'
    },

    {
      key: 'vm',
      header: '目标VM'
    },

    {
      key: 'sdSize',
      header: '恢复磁盘大小'
    },

    {
      key: 'createTime',
      header: '生成时间'
    },

    {
      key: 'action',
      header: '操作'
    }
  ];

  return (
    <>
      <Table.ActionPanel>
        <Justify right={<SearchBox />} />
      </Table.ActionPanel>

      <Table columns={columns} />
    </>
  );
};
