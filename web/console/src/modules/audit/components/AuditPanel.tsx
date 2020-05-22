import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Card, Form, Select, Text, DatePicker, TableColumn, Table, Button } from '@tea/component';
import { TablePanel } from '@tencent/ff-component';
import { bindActionCreators, insertCSS } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { useModal } from '../../common/utils';
import { emptyTips, LinkButton } from '../../common/components';
import { AuditDetailsDialog } from './AuditDetailsDialog';
import { allActions } from '../actions';
import { router } from '../router';
import { Audit } from '../models';
import { dateFormatter } from '@helper/dateFormatter';
const { useState, useEffect } = React;
const { RangePicker } = DatePicker;
const { expandable } = Table.addons;

insertCSS(
  'auditPanelDatePicker',
  `.auditPanelFilter .tea-form__item { display: inline }
`
);
export const AuditPanel = () => {
  const state = useSelector(state => state);
  const dispatch = useDispatch();
  const { actions } = bindActionCreators({ actions: allActions }, dispatch);
  const { auditList, auditFilterCondition } = state;
  const auditFilterConditions =
    auditFilterCondition.data[0] && auditFilterCondition.data[0].success ? auditFilterCondition.data[0].target : [];
  const {
    clusterName: clusterArr = [],
    namespace: namespaceArr = [],
    resource: resourceArr = [],
    userName: userNameArr = []
  } = auditFilterConditions;
  const clusterOptions = [
    { value: '', text: t('全部'), tooltip: t('全部') },
    ...clusterArr.map(item => ({ value: item, text: item, tooltip: item }))
  ];
  const namespaceOptions = [
    { value: '', text: t('全部'), tooltip: t('全部') },
    ...namespaceArr.map(item => ({ value: item, text: item, tooltip: item }))
  ];
  const resourceOptions = [
    { value: '', text: t('全部'), tooltip: t('全部') },
    ...resourceArr.map(item => ({ value: item, text: item, tooltip: item }))
  ];
  const userNameOptions = [
    { value: '', text: t('全部'), tooltip: t('全部') },
    ...userNameArr.map(item => ({ value: item, text: item, tooltip: item }))
  ];

  const [cluster, setCluster] = useState('');
  const [namespace, setNamespace] = useState('');
  const [resource, setResource] = useState('');
  const [user, setUser] = useState('');
  const [startTime, setStartTime] = useState(0);
  const [endTime, setEndTime] = useState(0);
  const [expandedKeys, setExpandedKeys] = useState([]);
  const [record, setRecord] = useState(undefined);
  const { isShowing, toggle } = useModal(false);

  useEffect(() => {
    actions.audit.getAuditFilterCondition.fetch();
  }, []);

  useEffect(() => {
    actions.audit.applyFilter({
      cluster,
      namespace,
      resource,
      user,
      startTime: startTime ? startTime : '',
      endTime: endTime ? endTime : ''
    });
  }, [cluster, namespace, resource, user, startTime, endTime]);

  const columns: TableColumn<Audit>[] = [
    {
      key: 'stageTimestamp',
      header: t('时间'),
      render: item => <Text parent="div">{dateFormatter(new Date(item.stageTimestamp), 'YYYY-MM-DD HH:mm:ss')}</Text>
    },
    {
      key: 'userName',
      header: t('操作人'),
      render: item => <Text parent="div">{item.userName || '-'}</Text>
    },
    {
      key: 'verb',
      header: t('操作类型'),
      render: item => <Text parent="div">{item.verb || '-'}</Text>
    },
    {
      key: 'clusterName',
      header: t('集群'),
      render: item => <Text parent="div">{item.clusterName || '-'}</Text>
    },
    {
      key: 'namespace',
      header: t('命名空间'),
      render: item => <Text parent="div">{item.namespace || '-'}</Text>
    },
    {
      key: 'resource',
      header: t('资源类型'),
      render: item => <Text parent="div">{item.resource || '-'}</Text>
    },
    {
      key: 'name',
      header: t('操作对象'),
      render: item => <Text parent="div">{item.name || '-'}</Text>
    },
    {
      key: 'status',
      header: t('操作结果'),
      render: item => <Text parent="div">{item.status || '-'}</Text>
    }
  ];

  return (
    <>
      <Card>
        <Form className="auditPanelFilter" style={{ padding: '20px' }}>
          <Form.Item label={t('操作集群')}>
            <Select
              boxSizeSync
              size="m"
              type="simulate"
              appearence="button"
              options={clusterOptions}
              value={cluster}
              onChange={value => setCluster(value)}
              placeholder={t('请选择集群id')}
            />
          </Form.Item>
          <Form.Item label={t('命名空间')}>
            <Select
              boxSizeSync
              size="m"
              type="simulate"
              appearence="button"
              options={namespaceOptions}
              value={namespace}
              onChange={value => setNamespace(value)}
              placeholder={t('请选择命名空间')}
            />
          </Form.Item>
          <Form.Item label={t('操作人')}>
            <Select
              boxSizeSync
              size="m"
              type="simulate"
              appearence="button"
              options={userNameOptions}
              value={user}
              onChange={value => setUser(value)}
              placeholder={t('请选择操作人')}
            />
          </Form.Item>
          <Form.Item label={t('对象')}>
            <Select
              boxSizeSync
              size="m"
              type="simulate"
              appearence="button"
              options={resourceOptions}
              value={resource}
              onChange={value => setResource(value)}
              placeholder={t('请选择操作人')}
            />
          </Form.Item>
          <Form.Item label={t('时间')}>
            <RangePicker
              showTime
              onChange={value => {
                setStartTime(new Date(value[0].format()).getTime());
                setEndTime(new Date(value[1].format()).getTime());
              }}
            />
          </Form.Item>
        </Form>
      </Card>
      <Card>
        <TablePanel
          recordKey={record => {
            return record.auditID;
          }}
          columns={columns}
          model={auditList}
          action={actions.audit}
          rowDisabled={record => record.status['phase'] === 'Terminating'}
          emptyTips={emptyTips}
          isNeedPagination={true}
          bodyClassName={'tc-15-table-panel tc-15-table-fixed-body'}
          addons={[
            expandable({
              expandedKeys,
              onExpandedKeysChange: (keys, { event }) => {
                event.stopPropagation();
                setExpandedKeys(keys);
              },
              render(record) {
                return (
                  <div>
                    <p>
                      <Trans>请求URI：</Trans>
                      {record.requestURI}
                    </p>
                    <p>
                      <Trans>源地址：</Trans>
                      {record.sourceIPs}
                    </p>
                    <p>HTTP Code：{record.code}</p>
                    <p>
                      <Button
                        type="link"
                        onClick={() => {
                          setRecord(record);
                          toggle();
                        }}
                      >
                        <Trans>记录详情</Trans>
                      </Button>
                    </p>
                  </div>
                );
              }
            })
          ]}
        />
        <AuditDetailsDialog isShowing={isShowing} toggle={toggle} record={record} />
      </Card>
    </>
  );
};
