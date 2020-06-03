import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Text, TableColumn, Table, Justify, SearchBox, Bubble } from '@tea/component';
import { TablePanel } from '@tencent/ff-component';
import { bindActionCreators, insertCSS } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { useModal } from '../../common/utils';
import { emptyTips, LinkButton } from '../../common/components';
import { allActions } from '../actions';
import { AlarmRecord } from '../models';
import { dateFormatter } from '@helper/dateFormatter';
const { useState, useEffect } = React;

export const AlarmRecordPanel = () => {
  const state = useSelector(state => state);
  const dispatch = useDispatch();
  const { actions } = bindActionCreators({ actions: allActions }, dispatch);
  const { alarmRecord, route } = state;
  const selectedClusterId = route.queries.clusterId;

  useEffect(() => {
    if (selectedClusterId) {
      actions.alarmRecord.applyFilter({ clusterID: selectedClusterId });
    } else {
      actions.alarmRecord.clear();
    }
  }, [selectedClusterId]);

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

  const columns: TableColumn<AlarmRecord>[] = [
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

  return (
    <>
      <Table.ActionPanel>
        <Justify
          left={<React.Fragment />}
          right={
            <SearchBox
              value={alarmRecord.query.keyword || ''}
              onChange={actions.alarmRecord.changeKeyword}
              onSearch={actions.alarmRecord.performSearch}
              onClear={() => {
                actions.alarmRecord.performSearch('');
              }}
              placeholder={t('请输入告警策略名')}
            />
          }
        />
      </Table.ActionPanel>
      <TablePanel
        recordKey={record => {
          return record.id;
        }}
        columns={columns}
        model={alarmRecord}
        action={actions.alarmRecord}
        rowDisabled={record => record.status['phase'] === 'Terminating'}
        emptyTips={emptyTips}
        isNeedContinuePagination={true}
        bodyClassName={'tc-15-table-panel tc-15-table-fixed-body'}
      />
    </>
  );
};
