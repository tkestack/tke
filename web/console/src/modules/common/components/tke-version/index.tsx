import React, { useEffect, useState } from 'react';
import { getTkeStackVersion } from '@/src/webApi/tkestack';
import { Text } from 'tea-component';

export const TkeVersion = () => {
  const [version, setVersion] = useState('');

  useEffect(() => {
    (async () => {
      const version = await getTkeStackVersion();
      setVersion(version);
    })();
  }, []);

  return (
    <Text
      tooltip={version}
      overflow
      style={{
        position: 'absolute',
        bottom: 0,
        left: 0,
        color: '#fff',
        paddingLeft: 10,
        paddingBottom: 10,
        maxWidth: '100%'
      }}
    >
      version: {version}
    </Text>
  );
};
