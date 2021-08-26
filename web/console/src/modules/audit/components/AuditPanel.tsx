/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Card, Form, Select, Text, DatePicker, Table, Button } from 'tea-component';
import { TablePanel, TablePanelColumnProps } from '@tencent/ff-component';
import { bindActionCreators, insertCSS } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { useModal } from '../../common/utils';
import { emptyTips, LinkButton } from '../../common/components';
import { AuditDetailsDialog } from './AuditDetailsDialog';
import { AuditSettingDialog } from './AuditSettingDialog';
import { allActions } from '../actions';
import { Audit } from '../models';
import { dateFormatter } from '@helper/dateFormatter';
import { downloadCsv } from '@helper';
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
  const [settingDialogVisible, showSettingDialog] = useState(false);

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

  const columns: TablePanelColumnProps<Audit>[] = [
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

  const style = { marginRight: 5 };

  /** 导出数据 */
  const handleDownload = (resourceList: Audit[]) => {
    const tableColumns = [
      {
        key: 'stageTimestamp',
        header: '时间',
        render: item => dateFormatter(new Date(item.stageTimestamp), 'YYYY-MM-DD HH:mm:ss')
      },
      {
        key: 'userName',
        header: '操作人',
        render: item => item.userName
      },
      {
        key: 'verb',
        header: '操作类型',
        render: item => item.verb
      },
      {
        key: 'clusterName',
        header: '集群',
        render: item => item.clusterName
      },
      {
        key: 'namespace',
        header: '命名空间',
        render: item => item.namespace
      },
      {
        key: 'resource',
        header: '资源类型',
        render: item => item.resource
      },
      {
        key: 'name',
        header: '操作对象',
        render: item => item.name
      },
      {
        key: 'status',
        header: '操作结果',
        render: item => item.status
      },
      {
        key: 'requestURI',
        header: '请求URI',
        render: item => item.requestURI
      },
      {
        key: 'sourceIPs',
        header: '源地址',
        render: item => item.sourceIPs
      },
      {
        key: 'code',
        header: 'HTTP Code',
        render: item => item.code
      }
    ];
    function getHeader() {
      // 生成csv需要的表头
      return tableColumns.map(aColumn => aColumn.header);
    }
    function getBody() {
      // 生成csv需要的内容
      return resourceList.reduce((accu, anAudit) => {
        let result;
        const _row = tableColumns.map(acolumn => acolumn.render(anAudit));
        result = [...accu, _row];
        return result;
      }, []);
    }
    downloadCsv(getBody(), getHeader(), `tke_audit_${new Date().getTime()}.csv`);
  };

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
          <Form.Item label={t('时间')} align="middle">
            <RangePicker
              showTime
              onChange={value => {
                setStartTime(new Date(value[0].format()).getTime());
                setEndTime(new Date(value[1].format()).getTime());
              }}
            />
            <Button
              icon="setting"
              title="设置"
              style={style}
              onClick={() => {
                showSettingDialog(!settingDialogVisible);
              }}
            />
            <Button
              icon="download"
              title="下载"
              style={style}
              onClick={() => {
                handleDownload(auditList.list.data.records);
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
        <AuditSettingDialog isShowing={settingDialogVisible} toggle={() => showSettingDialog(!settingDialogVisible)} />
      </Card>
    </>
  );
};
