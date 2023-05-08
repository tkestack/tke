import { getCustomConfig } from '@config';
import React, { useMemo } from 'react';

interface IPermissionProviderProps {
  children: React.ReactNode;
  value: string;
}

export function checkCustomVisible(value: string) {
  const path = value.split('.');
  const customConfig = getCustomConfig();

  const config = path.reduce((config, key) => config?.children?.[key], customConfig);

  return config?.visible ?? true;
}

export const PermissionProvider: React.FC<IPermissionProviderProps> = ({ value, children }) => {
  const isVisible = useMemo(() => checkCustomVisible(value), [value]);

  return <>{isVisible ? children : null}</>;
};
