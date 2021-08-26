/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

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
