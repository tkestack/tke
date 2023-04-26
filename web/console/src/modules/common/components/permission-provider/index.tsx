import { getCustomConfig } from '@config';
import React, { useMemo } from 'react';

interface IPermissionProviderProps {
  children: React.ReactElement;
  value: string;
}

function checkVisible(path: string[], config = getCustomConfig()) {
  if (path.length <= 0) return true;

  const [current, ...restPath] = path;

  const currentConfig = config?.children?.find(item => item.key === current);

  if (currentConfig?.data?.visible === false) return false;

  return checkVisible(restPath, currentConfig);
}

export const PermissionProvider: React.FC<IPermissionProviderProps> = ({ value, children }) => {
  const path = value.split('.');
  const isVisible = useMemo(() => checkVisible(path), [path]);

  return isVisible ? children : null;
};
