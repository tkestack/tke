import React from 'react';
import { ActionButton } from './actionButton';
import { virtualMachineAPI } from '@src/webApi';

export const DelButton = ({ type, clusterId, namespace, name, onSuccess = () => {} }) => {
  function del() {
    return virtualMachineAPI.deleteVM({ clusterId, namespace, name });
  }

  return (
    <ActionButton type={type} title={`你确定要删除虚拟机“${name}”吗？`} confirm={del} onSuccess={onSuccess}>
      删除
    </ActionButton>
  );
};
