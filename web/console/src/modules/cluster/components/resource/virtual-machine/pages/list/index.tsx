import React, { useEffect, useState } from 'react';
import { Table, Button, TableColumn, Text, Pagination, Dropdown, List, Icon, Bubble } from 'tea-component';
import { VMListActionPanel } from './action-panel';
import { useFetch } from '@src/modules/common/hooks';
import * as dayjs from 'dayjs';
import { useSetRecoilState, useRecoilState, useRecoilValue } from 'recoil';
import { clusterIdState, namespaceSelectionState, vmSelectionState } from '../../store/base';
import { virtualMachineAPI } from '@src/webApi';
import { router } from '@src/modules/cluster/router';
import { BootButton, ShutdownButton, DelButton, VNCButton, CreateSnapshotButton } from '../../components';

const { autotip } = Table.addons;

const defaultPageSize = 10;

export const VMListPanel = ({ route }) => {
  const [clusterId, setClusterId] = useRecoilState(clusterIdState);
  const namespace = useRecoilValue(namespaceSelectionState);
  const setVMSelection = useSetRecoilState(vmSelectionState);

  const [query, setQuery] = useState('');

  const columns: TableColumn[] = [
    {
      key: 'name',
      header: '名称',
      render(vm) {
        return (
          <Button
            type="link"
            onClick={() => {
              setVMSelection(vm);
              router.navigate(
                Object.assign({}, router.resolve(route), { mode: 'detail', resourceName: 'virtual-machine' }),
                Object.assign({}, route.queries, { resourceIns: vm?.name, np: namespace })
              );
            }}
          >
            {vm?.name}
          </Button>
        );
      }
    },

    {
      key: 'status',
      header: '状态',
      render({ status, failureCondition }) {
        const theme = status === 'Running' ? 'success' : 'danger';

        const bubbleContent = Object.entries(failureCondition || {}).map(([key, value]) => (
          <Text parent="p">
            {key}: {value}
          </Text>
        ));

        return (
          <>
            <Text theme={theme}>{status}</Text>
            {failureCondition && (
              <Bubble content={<> {bubbleContent}</>}>
                <Icon style={{ marginLeft: 5 }} type="info" />
              </Bubble>
            )}
          </>
        );
      }
    },

    {
      key: 'mirror',
      header: '镜像'
    },

    {
      key: 'ip',
      header: 'IP地址'
    },

    {
      key: 'hardware',
      header: '资源规格'
    },

    {
      key: 'createTime',
      header: '创建时间',
      render({ createTime }) {
        return createTime ? dayjs(createTime).format('YYYY-MM-DD HH:mm:ss') : '-';
      }
    },

    {
      key: 'actions',
      header: '操作',
      render({ name, status, supportSnapshot }) {
        return (
          <>
            <VNCButton type="link" clusterId={clusterId} name={name} namespace={namespace} status={status} />
            <Dropdown button={<Button type="link">更多</Button>} destroyOnClose={false}>
              <List style={{ padding: 10 }}>
                <List.Item>
                  <BootButton
                    clusterId={clusterId}
                    name={name}
                    namespace={namespace}
                    type="link"
                    status={status}
                    onSuccess={reFetch}
                  />
                </List.Item>

                <List.Item>
                  <ShutdownButton
                    clusterId={clusterId}
                    name={name}
                    namespace={namespace}
                    type="link"
                    status={status}
                    onSuccess={reFetch}
                  />
                </List.Item>

                <List.Item>
                  <DelButton clusterId={clusterId} name={name} namespace={namespace} type="link" onSuccess={reFetch} />
                </List.Item>

                <List.Item>
                  <CreateSnapshotButton
                    clusterId={clusterId}
                    name={name}
                    namespace={namespace}
                    disabled={!supportSnapshot}
                    onSuccess={() => {
                      const urlParams = router.resolve(route);
                      router.navigate(Object.assign({}, urlParams, { mode: 'snapshot' }), route.queries);
                    }}
                  />
                </List.Item>
              </List>
            </Dropdown>
          </>
        );
      }
    }
  ];

  useEffect(() => {
    const url = new URL(location.href);
    const clusterId = url.searchParams.get('clusterId');

    if (clusterId) setClusterId(clusterId);
  }, [location.href]);

  const {
    data: vmList,
    status,
    reFetch,
    paging
  } = useFetch(
    async ({ paging, continueToken }) => {
      const { items, newContinueToken, restCount } = await virtualMachineAPI.fetchVMListWithVMI(
        { clusterId, namespace },
        { limit: paging?.pageSize, continueToken, query }
      );
      return {
        data:
          items.map(({ metadata, status, spec, vmi }) => {
            let realStatus = status?.printableStatus;
            if (realStatus === 'Running' && !status?.ready) {
              realStatus = 'Abnormal';
            }

            const failureCondition =
              realStatus === 'Stopped' ? status?.conditions?.find(({ type }) => type === 'Failure') : null;

            return {
              name: metadata?.name,
              status: realStatus,
              failureCondition,
              mirror: metadata?.annotations?.['tkestack.io/image-display-name'] ?? '-',
              ip: vmi?.status?.interfaces?.[0]?.ipAddress ?? '-',
              hardware: `${spec?.template?.spec?.domain?.cpu?.cores ?? '-'}核 / ${
                spec?.template?.spec?.domain?.resources?.requests?.memory ?? '-'
              }`,
              createTime: metadata?.creationTimestamp,

              id: metadata?.uid,

              supportSnapshot:
                metadata?.annotations?.['tkestack.io/support-snapshot'] === 'true' &&
                (realStatus === 'Running' || realStatus === 'Stopped')
            };
          }) ?? [],

        continueToken: newContinueToken,

        totalCount: (paging.pageIndex - 1) * paging.pageSize + items.length + restCount
      };
    },
    [clusterId, namespace, query],
    {
      mode: 'continue',
      defaultPageSize,
      fetchAble: !!(clusterId && namespace),
      polling: true,
      pollingDelay: 30 * 1000,
      needClearData: false
    }
  );

  return (
    <>
      <VMListActionPanel route={route} reFetch={reFetch} vmList={vmList ?? []} onQueryChange={setQuery} />
      <Table
        columns={columns}
        records={vmList ?? []}
        recordKey="id"
        addons={[autotip({ isLoading: status === 'loading', isError: status === 'error' })]}
      />
      <Pagination
        recordCount={paging?.totalCount ?? 0}
        pageIndexVisible={false}
        endJumpVisible={false}
        pageSize={defaultPageSize}
        pageSizeVisible={false}
        onPagingChange={({ pageIndex }) => {
          if (pageIndex > paging.pageIndex) paging.nextPageIndex();

          if (pageIndex < paging.pageIndex) paging.prePageIndex();
        }}
      />
    </>
  );
};
