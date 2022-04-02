import React, { useEffect, useState } from 'react';
import { TeaFormLayout } from '@src/modules/common/layouts/TeaFormLayout';
import { Tabs, TabPanel, Breadcrumb, Button } from 'tea-component';
import { VMDetailTabEnum, VMDetailTabOptions } from '../../constants';
import { VMInfoPanel } from './info';
import { YamlPanel } from './yaml';
import { VMEventPanel } from './event';
import { router } from '@src/modules/cluster/router';

export const VMDetailPanel = () => {
  const [clusterId, setClusterId] = useState(null);
  const [namespace, setNamespace] = useState(null);
  const [vmName, setVmName] = useState(null);

  useEffect(() => {
    const searchParams = new URL(location.href).searchParams;
    const clusterId = searchParams.get('clusterId');
    const namespace = searchParams.get('np');
    const vmName = searchParams.get('resourceIns');

    if (clusterId) setClusterId(clusterId);
    if (namespace) setNamespace(namespace);
    if (vmName) setVmName(vmName);
  }, [location.href]);

  return (
    <TeaFormLayout
      title={
        <Breadcrumb>
          <Breadcrumb.Item>
            <a onClick={() => router.navigate({}, {})}>集群</a>
          </Breadcrumb.Item>

          <Breadcrumb.Item>
            <a
              onClick={() =>
                router.navigate(
                  {},
                  {},
                  `/tkestack/cluster/sub/list/resource/virtual-machine?clusterId=${clusterId}&np=${namespace}`
                )
              }
            >
              {clusterId}
            </a>
          </Breadcrumb.Item>

          <Breadcrumb.Item>
            <a>{`VirtualMachine:${vmName}(${namespace})`}</a>
          </Breadcrumb.Item>
        </Breadcrumb>
      }
      wrapCard={false}
    >
      <Tabs ceiling tabs={VMDetailTabOptions}>
        <TabPanel id={VMDetailTabEnum.Info}>
          <VMInfoPanel clusterId={clusterId} namespace={namespace} name={vmName} />
        </TabPanel>
        <TabPanel id={VMDetailTabEnum.Yaml}>
          <YamlPanel clusterId={clusterId} namespace={namespace} name={vmName} />
        </TabPanel>
        <TabPanel id={VMDetailTabEnum.Log}>
          <VMEventPanel clusterId={clusterId} namespace={namespace} name={vmName} />
        </TabPanel>
      </Tabs>
    </TeaFormLayout>
  );
};
