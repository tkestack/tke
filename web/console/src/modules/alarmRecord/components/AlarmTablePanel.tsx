import React, { useState } from 'react';
import { Justify, Table, TableColumn, Text, Bubble, Pagination, TagSearchBox, Button } from 'tea-component';
import { useFetch } from '@src/modules/common/hooks/useFetch';
import { fetchAlarmList } from '@src/webApi/alarm';
import { t } from '@/tencent/tea-app/lib/i18n';
import { dateFormatter } from '@helper/dateFormatter';

const { filterable, autotip } = Table?.addons;
const ALL_VALUE = '';

export const AlarmTablePanel = ({ clusterId }) => {
  const columns: TableColumn[] = [
    {
      key: 'metadata.creationTimestamp',
      header: t('发生时间'),
      render: item => <Text>{dateFormatter(new Date(item.metadata.creationTimestamp), 'YYYY-MM-DD HH:mm:ss')}</Text>
    },

    {
      key: 'spec.alarmPolicyName',
      header: t('告警策略'),
      render: item => <Text copyable>{item.spec.alarmPolicyName || '-'}</Text>
    },

    {
      key: 'spec.alarmPolicyType',
      header: t('策略类型'),
      render: item => <Text>{item.spec.alarmPolicyType || '-'}</Text>
    },

    {
      key: 'spec.body',
      header: t('告警内容'),
      width: 400,
      render: item => {
        const content = item.spec.body;
        const showContent = content.length >= 250 ? content.substr(0, 250) + '...' : content;
        return (
          <Bubble placement="left" content={content || null}>
            <Text>{showContent || '-'}</Text>
          </Bubble>
        );
      }
    },

    {
      key: 'status.alertStatus',
      header: '告警状态',
      render: item => (
        <Text theme={item?.status?.alertStatus === 'resolved' ? 'success' : 'danger'}>
          {item?.status?.alertStatus === 'resolved' ? '已恢复' : '未恢复'}
        </Text>
      )
    },

    {
      key: 'spec.receiverChannelName',
      header: t('通知渠道'),
      render: item => <Text copyable>{item.spec.receiverChannelName || '-'}</Text>
    },

    {
      key: 'spec.receiverName',
      header: t('接收人'),
      render: item => {
        return (
          <Text overflow copyable>
            {item?.spec?.receiverName ?? '-'}
          </Text>
        );
      }
    }
  ];

  const [query, setQuery] = useState({});

  const [alertStatus, setAlertStatus] = useState(ALL_VALUE);

  const {
    data: alarmList,
    paging,
    status
  } = useFetch(
    async ({ paging, continueToken }) => {
      const rsp = await fetchAlarmList({ clusterId }, { limit: paging?.pageSize, continueToken, query, alertStatus });

      return {
        data: rsp?.items ?? [],
        continueToken: rsp?.metadata?.continue ?? null,
        totalCount: null
      };
    },
    [clusterId, query, alertStatus],
    {
      mode: 'continue',
      fetchAble: !!clusterId,
      polling: true,
      pollingDelay: 5 * 1000,
      needClearData: false,
      defaultPageSize: 10,
      onlyPollingPage1: true
    }
  );

  return (
    <>
      <Table.ActionPanel>
        <Justify
          right={
            <TagSearchBox
              hideHelp
              minWidth={360}
              style={{ maxWidth: 640 }}
              attributes={[
                {
                  type: 'input',
                  key: 'spec.alarmPolicyName',
                  name: t('策略名称')
                },

                {
                  type: 'single',
                  key: 'spec.alarmPolicyType',
                  name: t('策略类型'),
                  values: [
                    { key: 'cluster', name: '集群' },
                    { key: 'node', name: '节点' },
                    { key: 'pod', name: 'Pod' },
                    { key: 'virtualMachine', name: '虚拟机' }
                  ]
                },

                {
                  type: 'input',
                  key: 'spec.receiverChannelName',
                  name: t('通知渠道')
                },

                {
                  type: 'input',
                  key: 'spec.receiverName',
                  name: t('接收人')
                }
              ]}
              onSearchButtonClick={(_, tags) => {
                const query = tags.reduce((all, tag) => {
                  const value = tag?.values?.[0];

                  return {
                    ...all,
                    [tag?.attr?.key]: value?.key ?? value?.name
                  };
                }, {});

                setQuery(query);
              }}
              onClearButtonClick={() => setQuery({})}
            />
          }
        />
      </Table.ActionPanel>

      <Table
        columns={columns}
        records={alarmList}
        addons={[
          filterable({
            type: 'single',
            column: 'status.alertStatus',
            value: alertStatus,
            onChange: value => setAlertStatus(value),
            // 增加 "全部" 选项
            all: {
              value: ALL_VALUE,
              text: '全部'
            },
            // 选项列表
            options: [
              { value: 'firing', text: '未恢复' },
              { value: 'resolved', text: '已恢复' }
            ]
          }),

          autotip({
            isLoading: status === 'loading',
            isError: status === 'error' || status === 'expired',
            errorText:
              status === 'expired' ? (
                <Button type="link" onClick={() => paging.setPageIndex(1)}>
                  ContinueToken 过期，点击跳转到第一页
                </Button>
              ) : (
                '加载失败'
              ),

            emptyText: '暂无数据'
          })
        ]}
      />

      <Pagination
        recordCount={paging?.totalCount ?? 0}
        stateText={<Text>{`第${paging.pageIndex}页`}</Text>}
        pageIndexVisible={false}
        endJumpVisible={false}
        pageSize={paging.pageSize}
        pageIndex={paging.pageIndex}
        pageSizeVisible={true}
        onPagingChange={({ pageIndex, pageSize }) => {
          if (pageSize !== paging.pageSize) {
            paging.setPageSize(pageSize);
            return;
          }

          if (pageIndex > paging.pageIndex) paging.nextPageIndex();

          if (pageIndex < paging.pageIndex) paging.prePageIndex();
        }}
      />
    </>
  );
};
