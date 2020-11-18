import React, { useCallback, useEffect } from 'react';
import { Card, Pagination, Table, TableProps } from '@tea/component';
import { PagingQuery } from '@tencent/ff-redux';
import { RootProps } from '../../ClusterApp';

interface PodTableProps extends RootProps, TableProps {}

export function PodTabel({
  subRoot: {
    resourceDetailState: {
      podList,
      podQuery: { paging, recordCount }
    }
  },
  actions,
  columns,
  addons
}: PodTableProps) {
  return (
    <Card>
      <Card.Body>
        <Table columns={columns} records={podList?.data?.records || []} recordKey="id" addons={addons} />
        <Pagination
          pageIndex={paging.pageIndex}
          pageSize={paging.pageSize}
          recordCount={recordCount}
          onPagingChange={actions.resourceDetail.pod.changePaging}
          stateText={`第${paging.pageIndex}页`}
          pageIndexVisible={false}
          endJumpVisible={false}
        />
      </Card.Body>
    </Card>
  );
}
