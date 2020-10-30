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

