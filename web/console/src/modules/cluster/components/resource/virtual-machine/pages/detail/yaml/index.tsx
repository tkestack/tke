import React from 'react';
import { YamlEditorPanel } from '@src/modules/common/components';
import { Card } from 'tea-component';
import { virtualMachineAPI } from '@src/webApi';
import { useFetch } from '@src/modules/common/hooks';

export const YamlPanel = ({ clusterId, namespace, name }) => {
  const { data } = useFetch(
    async () => {
      const yaml = await virtualMachineAPI.fetchVMForYaml({ clusterId, namespace, name });

      return {
        data: yaml
      };
    },
    [clusterId, namespace, name],
    {
      fetchAble: !!(clusterId && namespace && name)
    }
  );

  return (
    <Card>
      <Card.Body>
        <YamlEditorPanel config={data ?? ''} readOnly={true} />
      </Card.Body>
    </Card>
  );
};
