import { virtualMachineAPI } from '@src/webApi';
import { useMemo } from 'react';

export const useVMActions = ({ clusterId, namespace, name, status, onSuccess = () => {} }) => {
  return useMemo(() => {
    async function boot() {
      try {
        await virtualMachineAPI.setVMRunningStatus({ clusterId, namespace, name }, true);

        onSuccess();
      } catch (error) {}
    }

    async function shutdown() {
      try {
        await virtualMachineAPI.setVMRunningStatus({ clusterId, namespace, name }, false);
        onSuccess();
      } catch (error) {}
    }

    async function deleteVM() {
      try {
        await virtualMachineAPI.deleteVM({ clusterId, namespace, name });
        onSuccess();
      } catch (error) {}
    }

    return {
      boot: {
        disabled: status !== 'Stopped',
        dispatch: boot
      },

      shutdown: {
        disabled: status !== 'Running',
        dispatch: shutdown
      },

      deleteVM: {
        disabled: false,
        dispatch: deleteVM
      }
    };
  }, [clusterId, namespace, name, status, onSuccess]);
};
