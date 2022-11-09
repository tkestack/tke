import { Request } from './request';

export const fetchNamespaceList = (clusterId: string) => {
  return Request.get<any, { items: any }>('/api/v1/namespaces', {
    headers: {
      'X-TKE-ClusterName': clusterId
    }
  });
};

export const fetchNamespaceResourcequotas = ({ clusterId, name }) => {
  return Request.get<any, { items: any }>(`/api/v1/namespaces/${name}/resourcequotas`, {
    headers: {
      'X-TKE-ClusterName': clusterId
    }
  });
};

export const fetchNamespaceLimitranges = ({ clusterId, name }) => {
  return Request.get<any, { items: any }>(`/api/v1/namespaces/${name}/limitranges`, {
    headers: {
      'X-TKE-ClusterName': clusterId
    }
  });
};

export const modifyNamespaceLimitRange = ({ clusterId, name, resource, isCreate }) => {
  const resourceName = `${name}-limit-range`;

  const data = {
    kind: 'LimitRange',
    apiVersion: 'v1',
    metadata: {
      name: resourceName
    },

    spec: {
      limits: [
        {
          type: 'Container',
          ...resource?.spec?.limits?.[0]
        }
      ]
    }
  };

  if (isCreate) {
    return Request.post(`/api/v1/namespaces/${name}/limitranges`, data, {
      headers: {
        'X-TKE-ClusterName': clusterId
      }
    });
  }

  return Request.put(`/api/v1/namespaces/${name}/limitranges/${resourceName}`, data, {
    headers: {
      'X-TKE-ClusterName': clusterId
    }
  });
};

export const modifyNamespaceResourceQuota = ({ clusterId, name, resource, isCreate }) => {
  const resourceName = `${name}-resource-quota`;

  const data = {
    kind: 'ResourceQuota',
    apiVersion: 'v1',
    metadata: {
      name: resourceName
    },

    spec: {
      hard: {
        ...(resource.spec?.hard ?? {})
      }
    }
  };

  if (isCreate) {
    return Request.post(`/api/v1/namespaces/${name}/resourcequotas`, data, {
      headers: {
        'X-TKE-ClusterName': clusterId
      }
    });
  }

  return Request.put(`/api/v1/namespaces/${name}/resourcequotas/${resourceName}`, data, {
    headers: {
      'X-TKE-ClusterName': clusterId
    }
  });
};
