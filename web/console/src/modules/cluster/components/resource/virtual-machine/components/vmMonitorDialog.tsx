import React, { useMemo, useState } from 'react';
import { Modal, Button } from 'tea-component';
import { ChartInstancesPanel } from '@tencent/tchart';
import { vmMonitorFields } from '../constants';

const _VmMonitorDialog = ({ clusterId, namespace, vmList }) => {
  const [visible, setVisible] = useState(false);

  const monitorProps = useMemo(() => {
    return {
      tables: [
        {
          fields: vmMonitorFields,
          table: 'vm',
          conditions: [
            ['tke_cluster_instance_id', '=', clusterId],

            ['namespace', '=', namespace]
          ]
        }
      ],

      groupBy: [{ value: 'name' }],

      instance: {
        columns: [{ key: 'name', name: 'vm 名称' }],
        list: vmList.map(vm => ({
          name: vm?.name,
          isChecked: []
        }))
      }
    };
  }, [clusterId, namespace, vmList]);

  return (
    <>
      <Button type="primary" onClick={() => setVisible(true)}>
        监控
      </Button>

      <Modal size={1050} caption="virtual machine 监控" visible={visible} onClose={() => setVisible(false)}>
        <Modal.Body>
          <ChartInstancesPanel {...monitorProps} />
        </Modal.Body>
      </Modal>
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
