import React from 'react';
import { ActionButton } from './actionButton';
import { virtualMachineAPI } from '@src/webApi';
import { Text, Alert } from 'tea-component';

export const RecoverySnapshotButton = ({
  type,
  clusterId,
  namespace,
  name,
  vmName,
  disabled,
  onSuccess = () => {}
}) => {
  function recovery() {
    return virtualMachineAPI.recoverySnapshot({ clusterId, namespace, name, vmName });
  }

  return (
    <ActionButton
      type={type}
      title={`恢复快照`}
      confirm={recovery}
      onSuccess={onSuccess}
      disabled={disabled}
      body={
        <>
          <Alert type="warning">请确保虚拟机处于关机状态！</Alert>
          <Text parent="p">
            您将要对虚拟机<Text theme="warning">{vmName}</Text>恢复快照<Text theme="warning">{name}</Text>,
            恢复操作将会覆盖当前状态下虚拟机数据
          </Text>
          <Text parent="p">您是否要继续执行？</Text>
        </>
      }
    >
      恢复
    </ActionButton>
  );
};
