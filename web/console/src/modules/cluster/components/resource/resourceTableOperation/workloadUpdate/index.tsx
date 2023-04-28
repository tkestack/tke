import { t } from '@/tencent/tea-app/lib/i18n';
import { getParamByUrl } from '@helper';
import { router } from '@src/modules/cluster/router';
import { workloadApi } from '@src/webApi';
import { useRequest } from 'ahooks';
import React from 'react';
import { Breadcrumb, Form } from 'tea-component';
import { IWrokloadUpdatePanelProps, UpdateTypeEnum, updateType2text } from './constants';
import { ModifyNodeAffinityPanel } from './modifyNodeAffinityPanel';
import { ModifyStrategyPanel } from './modifyStrategyPanel';

export const WorkloadUpdatePanel = ({ kind, updateType, clusterVersion }: IWrokloadUpdatePanelProps) => {
  const clusterId = getParamByUrl('clusterId')!;
  const namespace = getParamByUrl('np')!;
  const resourceId = getParamByUrl('resourceIns')!;

  const { data: resource } = useRequest(
    () => {
      return workloadApi.fetchWorkloadResource({ namespace, clusterId, resourceId, kind, clusterVersion });
    },
    {
      ready: Boolean(namespace && resourceId && clusterId && kind)
    }
  );

  async function handleUpdate(data: any) {
    await workloadApi.updateWorkloadResource({
      clusterId,
      namespace,
      resourceId,
      kind,
      data,
      clusterVersion
    });

    goBackToListPanel();
  }

  function goBackToListPanel() {
    router.navigate({ sub: 'sub', mode: 'list', type: 'resource', resourceName: kind }, { clusterId, np: namespace });
  }

  const title = (
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
  );

  const baseInfo = (
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
  );

  return (
    <>
      {updateType === UpdateTypeEnum.ModifyStrategy && (
        <ModifyStrategyPanel
          kind={kind}
          resource={resource}
          title={title}
          baseInfo={baseInfo}
          onCancel={goBackToListPanel}
          onUpdate={handleUpdate}
        />
      )}

      {updateType === UpdateTypeEnum.ModifyNodeAffinity && (
        <ModifyNodeAffinityPanel
          kind={kind}
          resource={resource}
          title={title}
          baseInfo={baseInfo}
          onCancel={goBackToListPanel}
          onUpdate={handleUpdate}
        />
      )}
    </>
  );
};
