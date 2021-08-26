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

import { useState, useEffect, useRef, useCallback } from 'react';
import { fetchNamespaceList, fetchProjectNamespaceList } from '@src/modules/cluster/WebAPI/scale';

/**
 * 业务侧或者平台侧获取命名空间数据
 */
export const useNamespaces = ({ projectId, clusterId }) => {
  const [namespaces, setNamespaces] = useState({
    recordCount: 0,
    records: []
  });
  useEffect(() => {

    // 平台侧
    async function getPlatformNamespaces({ clusterId }) {
      const namespaceFetchResult = await fetchNamespaceList({ clusterId });
      setNamespaces(namespaceFetchResult);
    }

    // 业务侧
    async function getBusinessNamespaces({ projectId }) {
      const namespaceFetchResult = await fetchProjectNamespaceList({ projectId });
      setNamespaces(namespaceFetchResult);
    }

    if (projectId) {
      getBusinessNamespaces({ projectId });
    } else if (clusterId) {
      getPlatformNamespaces({ clusterId });
    }
  }, [projectId, clusterId]);

  return namespaces;
};

