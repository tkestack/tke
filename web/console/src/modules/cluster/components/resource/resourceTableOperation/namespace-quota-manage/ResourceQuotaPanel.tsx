import React, { useState, useEffect } from 'react';
import { Button, Table, Collapse, Justify, Alert, ExternalLink } from '@tencent/tea-component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import {
  getColumnByResource,
  computeResourceLimitRecords,
  storageResourceLimitRecord,
  othersResourceLimitRecord
} from './ResourceQuotaConfig';
import { Resource } from '@src/modules/common';
import { useRequest } from 'ahooks';
import { fetchNamespaceResourcequotas, modifyNamespaceResourceQuota } from '@src/webApi/namespace';

interface IResourceQuotaPanelProps {
  clusterId: string;
  name: string;
}

export function ResourceQuotaPanel({ clusterId, name }: IResourceQuotaPanelProps) {
  const [isEditMode, setIsEditMode] = useState(false);
  const [resource, setResource] = useState<Resource>(null);

  const { data: resourceQuota, refresh } = useRequest(
    async () => {
      const rsp = await fetchNamespaceResourcequotas({ clusterId, name });

      return rsp?.items?.at(0);
    },
    {
      ready: !!clusterId && !!name,
      onSuccess(data) {
        setResource(data || {});
      }
    }
  );

  function handleResourceChange({ path, value }: { path: string[]; value: string }) {
    setResource(pre => {
      const copedResource = JSON.parse(JSON.stringify(pre));

      const lastIndex = path.length - 1;

      return path.reduce((rsPart, k, index) => {
        if (index === lastIndex) {
          rsPart[k] = value;

          return copedResource;
        } else {
          rsPart[k] = rsPart[k] ?? {};

          return rsPart[k];
        }
      }, copedResource);
    });
  }

  function handleOK() {
    // 需要先做校验
    if (
      computeResourceLimitRecords.some(
        ({ rules, ...others }) =>
          rules?.map(rule => rule({ ...others, resource }))?.some(({ status }) => status === 'error') ?? false
      )
    ) {
      return;
    }

    modifyNamespaceResourceQuota({
      clusterId,
      name,
      resource,
      isCreate: !resourceQuota
    });

    setIsEditMode(false);
  }

  function handleCancel() {
    setResource(resourceQuota);

    setIsEditMode(false);
  }

  return (
    <Collapse defaultActiveIds={['computeResourceLimit', 'storageResourceLimit']}>
      <Alert style={{ marginTop: '10px', marginBottom: '10px' }}>
        <Trans>
          对命名空间设置CPU/内存 limit/request配额后，创建工作负载時，必须指定CPU/内存
          limit/request，或为命名空间配置LimitRange，更多请参考
          <ExternalLink href={'https://kubernetes.io/docs/concepts/policy/resource-quotas/'}>
            ResourceQuotas
          </ExternalLink>
        </Trans>
      </Alert>

      <Justify
        left={
          !isEditMode ? (
            <Button
              type="primary"
              style={{ marginBottom: '10px' }}
              onClick={() => {
                setIsEditMode(true);
              }}
            >
              {t('编辑配额')}
            </Button>
          ) : (
            <>
              <Button type="primary" style={{ marginBottom: '10px' }} onClick={handleOK}>
                {t('确定')}
              </Button>
              <Button type="weak" style={{ marginBottom: '10px', marginLeft: '5px' }} onClick={handleCancel}>
                {t('取消')}
              </Button>
            </>
          )
        }
        right={<Button icon="refresh" onClick={refresh} />}
      />

      <Collapse.Panel id="computeResourceLimit" style={{ marginBottom: '20px' }} title={<h2>{t('计算资源限制')}</h2>}>
        <Table
          columns={getColumnByResource({ resource, isEditMode, handleResourceChange })}
          records={computeResourceLimitRecords}
        />
      </Collapse.Panel>

      <Collapse.Panel id="storageResourceLimit" style={{ marginBottom: '20px' }} title={<h2>{t('存储资源限制')}</h2>}>
        <Table
          columns={getColumnByResource({ resource, isEditMode, handleResourceChange })}
          records={storageResourceLimitRecord}
        />
      </Collapse.Panel>

      <Collapse.Panel id="othersResourceLimit" style={{ marginBottom: '20px' }} title={<h2>{t('其他资源限制')}</h2>}>
        <Table
          columns={getColumnByResource({ resource, isEditMode, handleResourceChange })}
          records={othersResourceLimitRecord}
        />
      </Collapse.Panel>
    </Collapse>
  );
}
