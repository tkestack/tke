import React from 'react';
import { Button, ButtonProps } from 'tea-component';

interface IVNCButtonProps {
  clusterId: string;
  name: string;
  namespace: string;
  type: ButtonProps['type'];
  status: string;
}

export const VNCButton = ({ clusterId, namespace, name, type, status }: IVNCButtonProps) => {
  const href = `/tkestack/vnc?clusterId=${clusterId}&namespace=${namespace}&name=${name}`;

  return (
    <Button type={type} disabled={status !== 'Running'} onClick={() => window.open(href)}>
      登录
    </Button>
  );
};
