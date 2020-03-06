import { generateResourceInfo } from '../common';

/** event的相关配置 */
export const event = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'event',
    requestType: {
      list: 'events'
    },
    isRelevantToNamespace: true
  });
};

/** hpa的相关配置 */
export const hpa = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'hpa',
    requestType: {
      list: 'horizontalpodautoscalers'
    },
    isRelevantToNamespace: true
  });
};

/** cronhpa的相关配置 */
export const cronhpa = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'cronhpa',
    requestType: {
      list: 'cronhpas'
    },
    isRelevantToNamespace: true
  });
};

/** persistentEvent的相关配置 */
export const pe = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'pe',
    requestType: {
      list: 'persistentevents'
    }
  });
};

/** cluster的相关配置 */
export const cluster = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'cluster',
    requestType: {
      list: 'clusters'
    }
  });
};

export const clustercredential = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'clustercredential',
    requestType: {
      list: 'clustercredentials'
    }
  });
};

/** namespace的相关配置 */
export const namespace = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'ns',
    requestType: {
      list: 'namespaces'
    }
  });
};

/** 获取moduels的相关配置 */
export const moduleConfig = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'module',
    requestType: {
      list: 'sysinfo'
    }
  });
};

/** 登出的相关配置 */
export const logoutConfig = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'logout',
    requestType: {
      list: 'logout'
    }
  });
};

/** 获取info的相关配置 */
export const info = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'info',
    requestType: {
      list: 'tokens/info'
    }
  });
};
/** 获取info的相关配置 */
export const machines = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'machines',
    requestType: {
      list: 'machines'
    }
  });
};

/** 获取project的相关配置 */
export const projects = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'projects',
    requestType: {
      list: 'projects'
    }
  });
};

export const portal = (k8sVersion: string) => {
  // apiVersion的配置
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'portal',
    requestType: {
      list: 'portal'
    }
  });
};

/** 获取project的相关配置 */
export const platforms = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'platforms',
    requestType: {
      list: 'platforms'
    }
  });
};

export const namespaces = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'namespaces',
    requestType: {
      list: 'namespaces'
    }
  });
};

/** localidentities的配置 */
export const localidentity = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'localidentity',
    requestType: {
      list: 'localidentities'
    }
  });
};

/** policy的配置 */
export const policy = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'policy',
    requestType: {
      list: 'policies'
    }
  });
};

/** users的配置 */
export const user = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'user',
    requestType: {
      list: 'users'
    }
  });
};

/** roles的配置 */
export const role = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'role',
    requestType: {
      list: 'roles'
    }
  });
};

/** localgroups的配置 */
export const localgroup = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'localgroup',
    requestType: {
      list: 'localgroups'
    }
  });
};

/** groups的配置 */
export const group = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'group',
    requestType: {
      list: 'groups'
    }
  });
};

/** category的配置 */
export const category = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'category',
    requestType: {
      list: 'categories'
    }
  });
};

/** apikey的配置 */
export const apiKey = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'apiKey',
    requestType: {
      list: 'apikeys'
    }
  });
};

/** helm的配置 */
export const helm = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'helm',
    requestType: {
      list: 'helms'
    }
  });
};

/**
 *
 * logcs的配置
 */
export const logcs = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    isRelevantToNamespace: true,
    resourceName: 'logcs',
    requestType: {
      list: 'logcollector',
      addon: true
    }
  });
};

export const prometheus = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'prometheus',
    requestType: {
      list: 'prometheuses'
    }
  });
};
