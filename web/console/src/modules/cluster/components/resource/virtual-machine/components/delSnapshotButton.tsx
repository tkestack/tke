import React from 'react';
import { ActionButton } from './actionButton';
import { virtualMachineAPI } from '@src/webApi';

export const DelSnapshotButton = ({ type, clusterId, namespace, name, onSuccess = () => {} }) => {
  function del() {
    return virtualMachineAPI.delSnapshot({ clusterId, namespace, name });
  }

  return (
    <ActionButton type={type} title={`你确定要删除快照“${name}”吗？`} confirm={del} onSuccess={onSuccess}>
      删除
    </ActionButton>
  );
};
