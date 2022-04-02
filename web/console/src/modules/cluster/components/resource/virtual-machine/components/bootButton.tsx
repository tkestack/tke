import React from 'react';
import { ActionButton } from './actionButton';
import { virtualMachineAPI } from '@src/webApi';

export const BootButton = ({ type, clusterId, namespace, name, status, onSuccess = () => {} }) => {
  function boot() {
    return virtualMachineAPI.setVMRunningStatus({ clusterId, namespace, name }, true);
  }

  return (
    <ActionButton
      type={type}
      title={`你确定要启动虚拟机“${name}”吗？`}
      confirm={boot}
      disabled={status !== 'Stopped'}
      onSuccess={onSuccess}
    >
      开机
    </ActionButton>
  );
};
