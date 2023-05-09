import React, { useState } from 'react';
import { Table, TableColumn, Justify, SearchBox, Pagination, Button, Select, Text } from 'tea-component';
import { useFetch } from '@src/modules/common/hooks';
import { virtualMachineAPI, namespaceAPI } from '@src/webApi';
import { getParamByUrl } from '@helper';
import { useRequest } from 'ahooks';
import { router } from '@src/modules/cluster/router';
import dayjs from 'dayjs';
import { DelSnapshotButton, RecoverySnapshotButton } from '../../components';
import { TeaFormLayout } from '@src/modules/common/layouts/TeaFormLayout';

const { autotip } = Table.addons;

const defaultPageSize = 10;

export const SnapshotTablePanel = ({ route }) => {
  const clusterId = getParamByUrl('clusterId');
  const [namespace, setNamespace] = useState(() => getParamByUrl('np'));
  const [query, setQuery] = useState('');

  const { data: namespaceList = [] } = useRequest(
    async () => {
      const rsp = await namespaceAPI.fetchNamespaceList(clusterId);

      return rsp?.items?.map(item => ({ value: item?.metadata?.name })) ?? [];
    },
    {
      ready: Boolean(clusterId),
      refreshDeps: [clusterId]
    }
  );

  const {
    data = [],
    status,
    reFetch,
    paging
  } = useFetch(
    async ({ paging, continueToken }) => {
      const rsp = await virtualMachineAPI.fetchSnapshotList(
        { clusterId, namespace },
        { limit: paging?.pageSize, continueToken, query }
      );

      const items = rsp?.items ?? [];

      const newContinueToken = rsp?.metadata?.continue || null;

      const restCount = rsp?.metadata?.remainingItemCount ?? 0;

      return {
        data: items,
        continueToken: newContinueToken,
        totalCount: (paging.pageIndex - 1) * paging.pageSize + rsp?.items?.length + restCount
      };
    },
    [clusterId, namespace, query],
    {
      mode: 'continue',
      defaultPageSize,
      fetchAble: !!(clusterId && namespace)
    }
  );

  const columns: TableColumn[] = [
    {
      key: 'metadata.name',
      header: '快照名称',
      render(snapshot) {
        return <Text copyable>{snapshot?.metadata?.name}</Text>;
      }
    },

    {
      key: 'status.phase',
      header: '状态',
      render(snapshot) {
        const status = snapshot?.status?.phase;

        const theme = status === 'Succeeded' ? 'success' : 'danger';

        return <Text theme={theme}>{status}</Text>;
      }
    },

    {
      key: 'spec.source.name',
      header: '目标VM',
      render(snapshot) {
        return <Text copyable>{snapshot?.spec?.source?.name}</Text>;
      }
    },

    // {
    //   key: 'sdSize',
    //   header: '恢复磁盘大小'
    // },

    {
      key: 'metadata.creationTimestamp',
      header: '生成时间',
      render(snapshot) {
        const createTime = snapshot?.metadata?.creationTimestamp;

        return createTime ? dayjs(createTime).format('YYYY-MM-DD HH:mm:ss') : '-';
      }
    },

    {
      key: 'action',
      header: '操作',
      render(snapshot) {
        return (
          <>
            <DelSnapshotButton
              type="link"
              clusterId={clusterId}
              namespace={namespace}
              name={snapshot?.metadata?.name}
              onSuccess={reFetch}
            />

            <RecoverySnapshotButton
              type="link"
              disabled={snapshot?.status?.phase !== 'Succeeded'}
              clusterId={clusterId}
              namespace={namespace}
              name={snapshot?.metadata?.name}
              vmName={snapshot?.spec?.source?.name}
            />
          </>
        );
      }
    }
  ];

  return (
    <TeaFormLayout title="快照管理" wrapCard={false}>
      <Table.ActionPanel>
        <Justify
          right={
            <>
              <Text reset theme="label" verticalAlign="middle">
                命名空间
              </Text>

              <Select
                type="simulate"
                searchable
                appearence="button"
                size="s"
                style={{ width: '130px', marginRight: '5px' }}
                value={namespace}
                options={namespaceList}
                onChange={value => {
                  setNamespace(value);
                  const urlParams = router.resolve(route);
                  router.navigate(urlParams, Object.assign({}, route.queries, { np: value }));
                }}
              />

              <SearchBox onSearch={value => setQuery(value)} onClear={() => setQuery('')} />

              <Button type="icon" icon="refresh" onClick={reFetch} />
            </>
          }
        />
      </Table.ActionPanel>

      <Table
        columns={columns}
        records={data}
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
    </TeaFormLayout>
  );
};
