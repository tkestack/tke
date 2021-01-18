import React, { useEffect, useState } from 'react';
import { getTkeStackVersion } from '@/src/webApi/tkestack';

export const TkeVersion = () => {
  const [version, setVersion] = useState('');

  useEffect(() => {
    (async () => {
      const version = await getTkeStackVersion();
      setVersion(version);
    })();
  }, []);

  return <span style={{ fontSize: 12, color: '#fff', marginLeft: 5, transform: 'translateY(1px)' }}>{version}</span>;
};
