import React from 'react';
import { Form, Button, Text, Card, Table, TableColumn, Icon } from 'tea-component';
import { virtualMachineAPI, PVCAPI } from '@src/webApi';
import { useFetch } from '@src/modules/common/hooks';
import dayjs from 'dayjs';
import { BootButton, ShutdownButton, DelButton, VNCButton } from '../../../components';

export const VMInfoPanel = ({ clusterId, namespace, name }) => {
  const { data, reFetch } = useFetch(
    async () => {
      const { vm, vmi } = await virtualMachineAPI.fetchVMDetail({ clusterId, namespace, name });

      let realStatus = vm?.status?.printableStatus;
      if (realStatus === 'Running' && !vm?.status?.ready) {
        realStatus = 'Abnormal';
      }

      return {
        data: {
          description: vm?.metadata?.annotations?.description ?? '-',
          status: realStatus,
          createTime: vm?.metadata?.creationTimestamp,
          tags: Object.entries(vmi?.metadata?.labels ?? [])
            .map(([key, value]) => `${key}：${value}`)
            .join('、'),
          hardware: `${vm?.spec?.template?.spec?.domain?.cpu?.cores ?? '-'}核 / ${
            vm?.spec?.template?.spec?.domain?.resources?.requests?.memory ?? '-'
          }`,
          networkMode: '桥接',
          mirror: vm?.metadata?.annotations?.['tkestack.io/image-display-name'] ?? '-',
          ip: vmi?.status?.interfaces?.[0]?.ipAddress ?? '-',
          diskList: vm?.spec?.dataVolumeTemplates?.map((item, index) => ({
            name: item?.metadata?.name?.split('.')?.[0] ?? '-',
            pvcName: item?.metadata?.name,
            type: index === 0 ? '系统盘' : '数据盘',
            volumeMode: item?.spec?.pvc?.volumeMode ?? '-',
            storageClass: item?.spec?.storageClass ?? '-',
            size: item?.spec?.pvc?.resources?.requests?.storage ?? '-'
          }))
        }
      };
    },
    [clusterId, namespace, name],
    {
      fetchAble: !!(clusterId && namespace && name)
    }
  );

  const { data: scMap, status: scStatus } = useFetch(
    async () => {
      const pvcList: { pvcName: string; storageClassName: string }[] = await Promise.all(
        data?.diskList?.map(async ({ pvcName }) => {
          const pvc = await PVCAPI.fetchPVCInfo({ clusterId, namespace, name: pvcName });

          return {
            pvcName,
            storageClassName: pvc?.spec?.storageClassName
          };
        })
      );

      const scMap = pvcList.reduce((acc, cur) => ({ ...acc, [cur?.pvcName]: cur?.storageClassName }), {});

      return {
        data: scMap
      };
    },
    [clusterId, namespace, data],
    {
      fetchAble: !!(clusterId && namespace && data)
    }
  );

  const columns: TableColumn[] = [
    {
      key: 'name',
      header: '名称'
    },

    {
      key: 'type',
      header: '类型'
    },

    {
      key: 'volumeMode',
      header: '卷模式'
    },

    {
      key: 'storageClass',
      header: '存储类',
      render({ pvcName }) {
        return scStatus === 'loading' ? <Icon type="loading" /> : scMap?.[pvcName] ?? '-';
      }
    },

    {
      key: 'size',
      header: '容量'
    }
  ];

  return (
    <>
      <Card>
        <Card.Body
          operation={
            <>
              <VNCButton clusterId={clusterId} namespace={namespace} name={name} type="primary" status={data?.status} />

              <BootButton
                clusterId={clusterId}
                namespace={namespace}
                name={name}
                type="primary"
                status={data?.status}
                onSuccess={reFetch}
              />

              <ShutdownButton
                clusterId={clusterId}
                namespace={namespace}
                name={name}
                type="primary"
                status={data?.status}
                onSuccess={reFetch}
              />

              <DelButton
                clusterId={clusterId}
                namespace={namespace}
                name={name}
                type="primary"
                onSuccess={() => history.back()}
              />

              <Button icon="refresh" onClick={reFetch} />
            </>
          }
        >
          <Form readonly>
            <Form.Title>基本信息</Form.Title>

            <Form.Item label="名称">
              <Form.Text>{name}</Form.Text>
            </Form.Item>

            <Form.Item label="命名空间">
              <Form.Text>{namespace}</Form.Text>
            </Form.Item>

            <Form.Item label="描述">
              <Form.Text>{data?.description}</Form.Text>
            </Form.Item>

            <Form.Item label="状态">
              <Text reset theme={data?.status === 'Running' ? 'success' : 'danger'}>
                {data?.status}
              </Text>
            </Form.Item>

            <Form.Item label="创建时间">
              <Form.Text>{dayjs(data?.createTime).format('YYYY-MM-DD HH:mm:ss')}</Form.Text>
            </Form.Item>

            <Form.Item label="标签">
              <Form.Text>{data?.tags}</Form.Text>
            </Form.Item>

            <Form.Title style={{ marginTop: 20 }}>配置信息</Form.Title>

            <Form.Item label="资源规格">
              <Form.Text>{data?.hardware}</Form.Text>
            </Form.Item>

            <Form.Item label="网络模式">
              <Form.Text>{data?.networkMode}</Form.Text>
            </Form.Item>

            <Form.Item label="镜像">
              <Form.Text>{data?.mirror}</Form.Text>
            </Form.Item>

            <Form.Item label="IP地址">
              <Form.Text>{data?.ip}</Form.Text>
            </Form.Item>
          </Form>
        </Card.Body>
      </Card>

      <Card>
        <Card.Body>
          <Table columns={columns} records={data?.diskList ?? []} />
        </Card.Body>
      </Card>
    </>
  );
};
