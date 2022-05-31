import React, { useMemo, useState } from 'react';
import { Button } from 'tea-component';
import { vmMonitorGroups } from '../constants';
import { MonitorPanel } from '@src/modules/common/components/monitorPanel';

const _VmMonitorDialog = ({ clusterId, namespace, vmList }) => {
  const [visible, setVisible] = useState(false);

  return (
    <>
      <Button type="primary" onClick={() => setVisible(true)}>
        监控
      </Button>

      <MonitorPanel
        title="Virtual Machine 监控"
        conditions={[
          {
            key: 'tke_cluster_instance_id',
            value: clusterId,
            expr: '='
          },

          {
            key: 'namespace',
            value: namespace,
            expr: '='
          }
        ]}
        groups={vmMonitorGroups}
        instanceType="vm"
        instanceList={vmList.map(({ name }) => name)}
        defaultSelectedInstances={[vmList?.[0]?.name]}
        visible={visible}
        onClose={() => setVisible(false)}
      />
    </>
  );
};

export const VmMonitorDialog = React.memo(_VmMonitorDialog, (preProps, currentProps) => {
  if (
    preProps.clusterId === currentProps.clusterId &&
    preProps.namespace === currentProps.namespace &&
    JSON.stringify(preProps.vmList) === JSON.stringify(currentProps.vmList)
  ) {
    return true;
  }

  return false;
});
