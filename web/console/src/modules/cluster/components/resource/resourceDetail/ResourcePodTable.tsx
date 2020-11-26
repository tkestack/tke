import React, { useEffect } from 'react';
import { Card, Pagination, Table, TableProps } from '@tea/component';
import { RootProps } from '../../ClusterApp';
import { router } from '../../../router';
import { IsInNodeManageDetail } from './ResourceDetail';
import { ResourceFilter } from '@src/modules/common';

interface PodTableProps extends RootProps, TableProps {}

export function PodTabel({
  subRoot: {
    resourceDetailState: {
      podList,
      podQuery: { paging, recordCount },
      resourceDetailInfo: { selection }
    }
  },
  actions,
  columns,
  addons,
  route
}: PodTableProps) {
  useEffect(() => {
    const urlParams = router.resolve(route);
    const { type, resourceName } = urlParams;
    const isInNodeManage = IsInNodeManageDetail(type);

    if ((type === 'resource' || isInNodeManage) && resourceName !== 'cronjob') {
      const { rid, clusterId } = route.queries;
      let filter: ResourceFilter = {
        regionId: +rid,
        clusterId
      };

      if (!isInNodeManage) {
        if (!selection) return;

        filter = Object.assign(filter, {
          namespace: route.queries['np'],
          specificName: route.queries['resourceIns']
        });
      }
      // 进行pod列表的轮询拉取
      actions.resourceDetail.pod.poll(filter);
    }
  }, [selection, actions.resourceDetail.pod, route]);

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
          pageSizeVisible={podList?.data?.recordCount <= paging.pageSize}
        />
      </Card.Body>
    </Card>
  );
}
