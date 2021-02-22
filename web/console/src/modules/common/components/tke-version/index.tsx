import React, { useEffect, useState } from 'react';
import { getTkeStackVersion } from '@/src/webApi/tkestack';
import { Typography, Tooltip } from 'antd';

export const TkeVersion = () => {
  const [version, setVersion] = useState('');

  useEffect(() => {
    (async () => {
      const version = await getTkeStackVersion();
      setVersion(version);
    })();
  }, []);

  return (
    <Tooltip title={version}>
      <Typography.Text
        ellipsis
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
      </Typography.Text>
    </Tooltip>
  );
};
