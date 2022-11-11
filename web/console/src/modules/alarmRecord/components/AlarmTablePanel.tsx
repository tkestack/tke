import React, { useState } from 'react';
import { Justify, Table, SearchBox, TableColumn, Text, Bubble, Pagination } from 'tea-component';
import { useFetch } from '@src/modules/common/hooks/useFetch';
import { fetchAlarmList } from '@src/webApi/alarm';
import { t } from '@/tencent/tea-app/lib/i18n';
import { dateFormatter } from '@helper/dateFormatter';

const { filterable, autotip } = Table?.addons;
const ALL_VALUE = '';

const defaultPageSize = 10;
const formatManager = managers => {
  if (managers) {
    return managers.map((m, index) => {
      return (
        <p key={index} className="text-overflow">
          {m}
        </p>
      );
    });
  }
};

export const AlarmTablePanel = ({ clusterId }) => {
  const columns: TableColumn[] = [
    {
      key: 'metadata.creationTimestamp',
      header: t('发生时间'),
      render: item => (
        <Text parent="div">{dateFormatter(new Date(item.metadata.creationTimestamp), 'YYYY-MM-DD HH:mm:ss')}</Text>
      )
    },

    {
      key: 'spec.alarmPolicyName',
      header: t('告警策略'),
      render: item => <Text parent="div">{item.spec.alarmPolicyName || '-'}</Text>
    },

    {
      key: 'spec.alarmPolicyType',
      header: t('策略类型'),
      render: item => <Text parent="div">{item.spec.alarmPolicyType || '-'}</Text>
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
            <Text parent="div">{showContent || '-'}</Text>
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
      render: item => <Text parent="div">{item.spec.receiverChannelName || '-'}</Text>
    },

    {
      key: 'spec.receiverName',
      header: t('接受组'),
      render: item => {
        const members = item.spec.receiverName ? item.spec.receiverName.split(',') : [];
        return (
          <Bubble placement="left" content={formatManager(members) || null}>
            <span className="text">
              {formatManager(members ? members.slice(0, 1) : [])}
              <Text parent="div" overflow>
                {members && members.length > 1 ? '...' : ''}
              </Text>
            </span>
          </Bubble>
        );
      }
    }
  ];

  const [query, setQuery] = useState('');

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
      defaultPageSize
    }
  );

  return (
    <>
      <Table.ActionPanel>
        <Justify right={<SearchBox onSearch={value => setQuery(value)} onClear={() => setQuery('')} />} />
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
            isError: status === 'error',
            emptyText: '暂无数据'
          })
        ]}
      />

      <Pagination
        recordCount={paging?.totalCount ?? 0}
        stateText={<Text>{`第${paging.pageIndex}页`}</Text>}
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
