import { t } from '@/tencent/tea-app/lib/i18n';
import { getParamByUrl } from '@helper';
import { router } from '@src/modules/cluster/router';
import { TeaFormLayout } from '@src/modules/common/layouts/TeaFormLayout';
import { workloadApi } from '@src/webApi';
import { useRequest } from 'ahooks';
import React, { useState } from 'react';
import { Breadcrumb, Button, Form, message } from 'tea-component';
import { IWrokloadUpdatePanelProps, UpdateTypeEnum, updateType2text } from './constants';
import { ModifyNodeAffinityPanel } from './modifyNodeAffinityPanel';
import { ModifyStrategyPanel } from './modifyStrategyPanel';

export const WorkloadUpdatePanel = ({ kind, updateType }: IWrokloadUpdatePanelProps) => {
  const clusterId = getParamByUrl('clusterId')!;
  const namespace = getParamByUrl('np')!;
  const resourceId = getParamByUrl('resourceIns')!;

  const [flag, setFlag] = useState(false);

  const { data: resource } = useRequest(
    () => {
      return workloadApi.fetchWorkloadResource({ namespace, clusterId, resourceId, kind });
    },
    {
      ready: Boolean(namespace && resourceId && clusterId && kind)
    }
  );

  function goBackToListPanel() {
    router.navigate({ sub: 'sub', mode: 'list', type: 'resource', resourceName: kind }, { clusterId, np: namespace });
  }

  async function onSubmit(data?: any) {
    if (data) {
      try {
        await workloadApi.updateWorkloadResource({
          clusterId,
          namespace,
          resourceId,
          kind,
          data
        });

        goBackToListPanel();
      } catch (error: any) {
        message.error({ content: error?.response?.data?.message ?? '请求失败！' });
      }
    }

    setFlag(false);
  }

  return (
    <TeaFormLayout
      title={
        <Breadcrumb>
          <Breadcrumb.Item>
            <a onClick={() => router.navigate({}, {})}>集群</a>
          </Breadcrumb.Item>

          <Breadcrumb.Item>
            <a onClick={goBackToListPanel}>{clusterId}</a>
          </Breadcrumb.Item>

          <Breadcrumb.Item>{`${kind}:${resourceId}(${namespace})`}</Breadcrumb.Item>

          <Breadcrumb.Item>{updateType2text[updateType]}</Breadcrumb.Item>
        </Breadcrumb>
      }
      footer={
        <>
          <Button type="primary" style={{ marginRight: 10 }} onClick={() => setFlag(true)}>
            {updateType2text[updateType]}
          </Button>

          <Button onClick={goBackToListPanel}>取消</Button>
        </>
      }
    >
      <>
        <Form>
          <Form.Title>基本信息</Form.Title>

          <Form.Item label={t('集群ID')}>
            <Form.Text>{clusterId}</Form.Text>
          </Form.Item>

          <Form.Item label={t('所在命名空间')}>
            <Form.Text>{namespace}</Form.Text>
          </Form.Item>

          <Form.Item label={t('资源名称')}>
            <Form.Text>{`${resourceId} (${kind})`}</Form.Text>
          </Form.Item>
        </Form>

        <hr />

        {updateType === UpdateTypeEnum.ModifyStrategy && (
          <ModifyStrategyPanel kind={kind} resource={resource} onSubmit={onSubmit} flag={flag} />
        )}

        {updateType === UpdateTypeEnum.ModifyNodeAffinity && (
          <ModifyNodeAffinityPanel flag={flag} onSubmit={onSubmit} />
        )}
      </>
    </TeaFormLayout>
  );
};
