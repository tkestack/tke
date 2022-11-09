import { t } from '@/tencent/tea-app/lib/i18n';
import React, { useState } from 'react';
import { Button, Modal, Tabs, TabPanel } from 'tea-component';
import { LimitRangePanel } from './LimitRangePanel';
import { ResourceQuotaPanel } from './ResourceQuotaPanel';

const tabs = [
  {
    id: 'resourceQuota',
    label: t('资源配额与限制')
  },

  {
    id: 'limitRange',
    label: t('编辑LimitRange')
  }
];

export const NamespaceQuotaManage = ({ name, clusterId }) => {
  const [visible, setVisible] = useState(false);

  return (
    <>
      <Button type="link" onClick={() => setVisible(true)}>
        配额管理
      </Button>

      <Modal size="l" visible={visible} onClose={() => setVisible(false)} caption={`Namespace: ${name}`}>
        <Tabs tabs={tabs} destroyInactiveTabPanel={false}>
          <TabPanel id="resourceQuota">
            <ResourceQuotaPanel clusterId={clusterId} name={name} />
          </TabPanel>

          <TabPanel id="limitRange">
            <LimitRangePanel clusterId={clusterId} name={name} />
          </TabPanel>
        </Tabs>
      </Modal>
    </>
  );
};
