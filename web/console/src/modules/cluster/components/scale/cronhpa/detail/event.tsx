/**
 * hpa详情，事件列表tab
 */
import React, { useState, useEffect, useContext, useMemo, useCallback } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import * as classnames from 'classnames';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { isEmpty } from '@src/modules/common/utils';
import { router } from '@src/modules/cluster/router';
import {
  Layout,
  Text,
  Table,
  Bubble
} from '@tencent/tea-component';
import { Clip } from '../../../../../common/components';
import { dateFormat, dateFormatter } from '@helper';
import { fetchEventList } from '@src/modules/cluster/WebAPI/scale';

const { autotip } = Table.addons;

const Event = React.memo((props: {
  selectedHpa: any;
  refreshFlag: number;
}) => {
  const route = useSelector((state) => state.route);
  const { clusterId } = route.queries;
  const { selectedHpa, refreshFlag } = props;

  /**
   * 获取事件列表
   */
  const [eventData, setEventData] = useState();
  useEffect(() => {
    async function getEventList(namespace, clusterId, name, uid) {
      const eventData = await fetchEventList({ namespace, clusterId, name, uid });
      setEventData(eventData);
    }
    if (!isEmpty(selectedHpa)) {
      const { name, namespace, uid } = selectedHpa.metadata;
      getEventList(namespace, clusterId, name, uid);
    }
  }, [selectedHpa, refreshFlag]);

  /** 时间处理函数 */
  const reduceTime = useCallback((time: string) => {
    let [first, second] = dateFormatter(new Date(time), 'YYYY-MM-DD HH:mm:ss').split(' ');
    return (
      <React.Fragment>
        <Text>{`${first} ${second}`}</Text>
      </React.Fragment>
    );
  }, []);

  const columns = [
    {
      key: 'firstTime',
      header: t('首次出现时间'),
      width: '10%',
      render: x => reduceTime(x.firstTimestamp)
    },
    {
      key: 'lastTime',
      header: t('最后出现时间'),
      width: '10%',
      render: x => reduceTime(x.lastTimestamp)
    },
    {
      key: 'type',
      header: t('级别'),
      width: '8%',
      render: x => (
        <div>
          <p className={classnames('text-overflow', { 'text-danger': x.type === 'Warning' })}>{x.type}</p>
        </div>
      )
    },
    {
      key: 'resourceType',
      header: t('资源类型'),
      width: '8%',
      render: x => (
        <div>
          <p title={x.involvedObject.kind} className="text-overflow">
            {x.involvedObject.kind}
          </p>
        </div>
      )
    },
    {
      key: 'name',
      header: t('资源名称'),
      width: '12%',
      render: x => (
        <div>
          <span id={'eventName' + x.id} title={x.metadata.name} className="text-overflow m-width">
            {x.metadata.name}
          </span>
          <Clip target={'#eventName' + x.id} />
        </div>
      )
    },
    {
      key: 'content',
      header: t('内容'),
      width: '12%',
      render: x => (
        <Bubble placement="bottom" content={x.reason || null}>
          <Text parent="div" overflow>
            {x.reason}
          </Text>
        </Bubble>
      )
    },
    {
      key: 'desp',
      header: t('详细描述'),
      width: '15%',
      render: x => (
        <Bubble placement="bottom" content={x.message || null}>
          <Text parent="div" overflow>
            {x.message}
          </Text>
        </Bubble>
      )
    },
    {
      key: 'count',
      header: t('出现次数'),
      width: '6%',
      render: x => (
        <div>
          <Text parent="div" overflow>
            {x.count}
          </Text>
        </div>
      )
    }
  ];

  return !isEmpty(selectedHpa) && (
    <Table
      records={eventData ? eventData.records : []}
      recordKey="id"
      columns={columns}
      addons={[
        autotip({
          // isLoading: loading,
          emptyText: t('事件列表为空')
        }),
      ]}
    />
  );
});
export default Event;
