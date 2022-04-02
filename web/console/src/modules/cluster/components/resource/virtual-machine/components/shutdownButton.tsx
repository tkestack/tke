import React from 'react';
import { ActionButton } from './actionButton';
import { virtualMachineAPI } from '@src/webApi';

export const ShutdownButton = ({ type, clusterId, namespace, name, status, onSuccess = () => {} }) => {
  function shutdown() {
    return virtualMachineAPI.setVMRunningStatus({ clusterId, namespace, name }, false);
  }

  return (
    <ActionButton
      type={type}
      title={`你确定要关闭虚拟机“${name}”吗？`}
      confirm={shutdown}
      disabled={status !== 'Running'}
      onSuccess={onSuccess}
    >
      关机
    </ActionButton>
  );
};
