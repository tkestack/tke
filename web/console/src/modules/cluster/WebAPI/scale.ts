import { Modal } from '@tea/component';
import { Method, reduceK8sRestfulPath, reduceNetworkRequest } from '@helper/index';
import { RequestParams, Resource, ResourceInfo } from '../../common/models';
import { QueryState, RecordSet } from '@tencent/ff-redux/src';
import { ResourceFilter } from '@src/modules/cluster/models';
import { resourceConfig } from '@config/resourceConfig';
import { uuid } from '@tencent/ff-redux/libs/qcloud-lib';
import { cluster } from '@config/resource/k8sConfig';
import { cutNsStartClusterId } from '@helper';

const _ = require('lodash');
const { get, isEmpty } = require('lodash');

function alertError(error, url) {
  let message = `错误码：${error.response.status}，错误描述：${error.response.statusText}`;
  let description = `请求路径: ${url} `;
  // if (error.response.status === 500) {
  if (error.response.data) {
    // 内部异常的response可能是文本也可能是错误对象
    description += `错误消息：${_(error.response.data).value()}`;
  }
  Modal.error({
    message,
    description
  });
}

/**
 * 根据Project查询Namespace列表
 * @param projectId
 */
export async function fetchProjectNamespaceList({ projectId }: { projectId?: string }) {
  let NamespaceResourceInfo: ResourceInfo = resourceConfig().namespaces;
  let url = reduceK8sRestfulPath({
    resourceInfo: NamespaceResourceInfo,
    specificName: projectId,
    extraResource: 'namespaces'
  });

  let params: RequestParams = {
    method: Method.get,
    url
  };

  let namespaceList = [],
    total = 0;
  try {
    let response = await reduceNetworkRequest(params);
    if (response.code === 0) {
      let list = response.data;
      total = list.items.length;
      namespaceList = list.items.map(item => {
        return {
          ...item,
          id: uuid(),
          value: item.metadata.name,
          text: `${item.spec.namespace}(${item.spec.clusterName})`
        };
      });
    }
  } catch (error) {
    alertError(error, url);
  }
  const result: RecordSet<Resource> = {
    recordCount: total,
    records: namespaceList
  };
  return result;
}

export async function fetchNamespaceList({ clusterId }: { clusterId?: string }) {
  let url = 'api/v1/namespaces';
  let params: RequestParams = {
    method: Method.get,
    url
  };

  let namespaceList = [],
    total = 0;
  try {
    let response = await reduceNetworkRequest(params, clusterId);
    if (response.code === 0) {
      let list = response.data;
      total = list.items.length;
      namespaceList = list.items.map(item => {
        return {
          ...item,
          id: uuid(),
          value: item.metadata.name,
          text: item.metadata.name
        };
      });
    }
  } catch (error) {
    alertError(error, url);
  }

  const result: RecordSet<Resource> = {
    recordCount: total,
    records: namespaceList
  };
  return result;
}

/**
 * 获取全部HPA列表数据
 */
export async function getHPAList({ namespace, clusterId }: { namespace: string; clusterId: string }) {
  let url = `/apis/autoscaling/v2beta1/namespaces/${namespace}/horizontalpodautoscalers`;
  let params: RequestParams = {
    method: Method.get,
    url
  };

  try {
    let response = await reduceNetworkRequest(params, clusterId);
    let HPAList = [],
      total = 0;
    if (response.code === 0) {
      let list = response.data;
      total = list.items.length;
      HPAList = list.items.map(item => {
        return {
          ...item,
          id: uuid()
        };
      });
    }

    const result: RecordSet<Resource> = {
      recordCount: total,
      records: HPAList
    };
    return result;
  } catch (error) {
    alertError(error, url);
  }
}

/**
 * 删除HPA
 */
export async function removeHPA({
  namespace,
  clusterId,
  name
}: {
  namespace: string;
  clusterId: string;
  name: string;
}) {
  const newNamespace = cutNsStartClusterId({ namespace, clusterId });
  let url = `/apis/autoscaling/v1/namespaces/${newNamespace}/horizontalpodautoscalers/${name}`;
  let params: RequestParams = {
    method: Method.delete,
    url
  };

  try {
    let response = await reduceNetworkRequest(params, clusterId);
    if (response.code === 0) {
      return true;
    } else {
      alertError({ response }, url);
      return false;
    }
  } catch (error) {
    alertError(error, url);
  }
}

/**
 * 创建HPA
 */
export async function createHPA({
  namespace,
  clusterId,
  hpaData
}: {
  namespace: string;
  clusterId: string;
  hpaData: any;
}) {
  const newNamespace = namespace.replace(new RegExp(`^${clusterId}-`), '');
  let url = `/apis/autoscaling/v2beta1/namespaces/${newNamespace}/horizontalpodautoscalers`;
  let params: RequestParams = {
    method: Method.post,
    url,
    data: hpaData
  };

  try {
    let response = await reduceNetworkRequest(params, clusterId);
    if (response.code === 0) {
      return true;
    }
  } catch (error) {
    alertError(error, url);
  }
}

/**
 * 更新HPA
 */
export async function modifyHPA({
  namespace,
  clusterId,
  name,
  hpaData
}: {
  namespace: string;
  clusterId: string;
  name: string;
  hpaData: any;
}) {
  const newNamespace = namespace.replace(new RegExp(`^${clusterId}-`), '');
  let url = `/apis/autoscaling/v2beta1/namespaces/${newNamespace}/horizontalpodautoscalers/${name}`;
  let params: RequestParams = {
    method: Method.put,
    url,
    data: hpaData
  };

  try {
    let response = await reduceNetworkRequest(params, clusterId);
    if (response.code === 0) {
      return true;
    }
  } catch (error) {
    alertError(error, url);
  }
}

/**
 * 获取YAML
 */
export async function fetchHPAYaml({
  namespace,
  clusterId,
  name
}: {
  namespace: string;
  clusterId: string;
  name: string;
}) {
  let url = `/apis/autoscaling/v1/namespaces/${namespace}/horizontalpodautoscalers/${name}`;
  const userDefinedHeader = {
    Accept: 'application/yaml'
  };
  let params: RequestParams = {
    method: Method.get,
    url,
    userDefinedHeader
  };

  try {
    let response = await reduceNetworkRequest(params, clusterId);
    let yamlList = response.code === 0 ? [response.data] : [];

    const result: RecordSet<string> = {
      recordCount: yamlList.length,
      records: yamlList
    };

    return result;
  } catch (error) {
    alertError(error, url);
  }
}

/**
 * 更新YAML
 */
export async function modifyHPAYaml({
  namespace,
  clusterId,
  name,
  yamlData
}: {
  namespace: string;
  clusterId: string;
  name: string;
  yamlData: any;
}) {
  let HPAList = [];
  // let url = '/apis/platform.tkestack.io/v1/clusters';
  let url = `/apis/autoscaling/v1/namespaces/${namespace}/horizontalpodautoscalers/${name}`;
  const userDefinedHeader = {
    Accept: 'application/json',
    'Content-Type': 'application/yaml'
  };
  let params: RequestParams = {
    method: Method.put,
    url,
    userDefinedHeader,
    data: yamlData
  };

  try {
    let response = await reduceNetworkRequest(params, clusterId);
    let yamlList = response.code === 0 ? [response.data] : [];
    if (response.code === 0) {
      return response;
    }
  } catch (error) {
    alertError(error, url);
  }
}

/**
 * 根据资源类型查询对应资源列表数据
 * @param projectId
 */
export async function fetchResourceList({
  resourceType,
  namespace,
  clusterId
}: {
  resourceType: string;
  namespace: string;
  clusterId: string;
}) {
  const newNamespace = namespace.replace(new RegExp(`^${clusterId}-`), '');
  let url = `/apis/apps/v1/namespaces/${newNamespace}/${resourceType}`;
  if (resourceType === 'tapps') {
    url = `/apis/platform.tkestack.io/v1/clusters/${clusterId}/tapps?namespace=${newNamespace}`;
  }
  let params: RequestParams = {
    method: Method.get,
    url
  };

  try {
    let response = await reduceNetworkRequest(params, clusterId);
    let resourceList = [],
      total = 0;
    if (response.code === 0) {
      let list = response.data;
      total = list.items.length;
      resourceList = list.items.map(item => {
        return {
          ...item,
          id: uuid(),
          value: item.metadata.name,
          text: item.metadata.name
        };
      });
    }

    const result: RecordSet<Resource> = {
      recordCount: total,
      records: resourceList
    };
    return result;
  } catch (error) {
    alertError(error, url);
  }
}

/**
 * 获取HPA事件列表
 */
export async function fetchEventList({
  type,
  namespace,
  clusterId,
  name,
  uid
}: {
  type?: string;
  namespace: string;
  clusterId: string;
  name: string;
  uid: any;
}) {
  const newNamespace = namespace.replace(new RegExp(`^${clusterId}-`), '');
  let url = `/api/v1/namespaces/${newNamespace}/events?fieldSelector=involvedObject.namespace=${newNamespace},involvedObject.kind=HorizontalPodAutoscaler,involvedObject.uid=${uid},involvedObject.name=${name}`;
  if (type === 'cronhpa') {
    url = `/api/v1/namespaces/${newNamespace}/events?fieldSelector=involvedObject.namespace=${newNamespace},involvedObject.kind=CronHPA,involvedObject.uid=${uid},involvedObject.name=${name}`;
  }
  let params: RequestParams = {
    method: Method.get,
    url
  };

  try {
    let response = await reduceNetworkRequest(params, clusterId);
    let eventList = [],
      total = 0;
    if (response.code === 0) {
      let list = response.data;
      total = list.items.length;
      eventList = list.items.map(item => {
        return {
          ...item,
          id: uuid()
        };
      });
    }

    const result: RecordSet<Resource> = {
      recordCount: total,
      records: eventList
    };
    return result;
  } catch (error) {
    alertError(error, url);
  }
}

/**
 * 获取全部CronHPA列表数据
 */
export async function fetchCronHpaRecords({ namespace, clusterId }: { namespace: string; clusterId: string }) {
  // let url = '/apis/platform.tkestack.io/v1/clusters';
  const newNamespace = cutNsStartClusterId({ namespace, clusterId });
  let url = `/apis/platform.tkestack.io/v1/clusters/${clusterId}/cronhpas?namespace=${newNamespace}`;
  // let url = `/apis/autoscaling/v2beta1/namespaces/${namespace}/horizontalpodautoscalers`;
  let params: RequestParams = {
    method: Method.get,
    url
  };

  try {
    let response = await reduceNetworkRequest(params, clusterId);
    let cronHpaList = [],
      total = 0;
    if (response.code === 0) {
      let list = response.data;
      total = list.items.length;
      cronHpaList = list.items.map(item => {
        return {
          ...item,
          id: uuid()
        };
      });
    }

    const result: RecordSet<Resource> = {
      recordCount: total,
      records: cronHpaList
    };
    return result;
  } catch (error) {
    alertError(error, url);
  }
}

/**
 * 删除CronHPA
 */
export async function deleteCronHpa({
  namespace,
  clusterId,
  name
}: {
  namespace: string;
  clusterId: string;
  name: string;
}) {
  const newNamespace = cutNsStartClusterId({ namespace, clusterId });
  let url = `/apis/platform.tkestack.io/v1/clusters/${clusterId}/cronhpas?name=${name}&namespace=${newNamespace}`;
  let params: RequestParams = {
    method: Method.delete,
    url
  };

  try {
    let response = await reduceNetworkRequest(params, clusterId);
    if (response.code === 0) {
      return true;
    } else {
      alertError({ response }, url);
      return false;
    }
  } catch (error) {
    alertError(error, url);
  }
}

/**
 * 创建CronHPA
 */
export async function createCronHpa({
  namespace,
  clusterId,
  cronHpaData
}: {
  namespace: string;
  clusterId: string;
  cronHpaData: any;
}) {
  const newNamespace = cutNsStartClusterId({ namespace, clusterId });
  let url = `/apis/platform.tkestack.io/v1/clusters/${clusterId}/cronhpas?namespace=${newNamespace}`;
  let params: RequestParams = {
    method: Method.post,
    url,
    data: cronHpaData
  };

  try {
    let response = await reduceNetworkRequest(params, clusterId);
    if (response.code === 0) {
      return true;
    }
  } catch (error) {
    alertError(error, url);
  }
}

/**
 * 更新CronHPA
 */
export async function modifyCronHpa({
  namespace,
  clusterId,
  name,
  cronHpaData
}: {
  namespace: string;
  clusterId: string;
  name: string;
  cronHpaData: any;
}) {
  const newNamespace = cutNsStartClusterId({ namespace, clusterId });
  let url = `/apis/platform.tkestack.io/v1/clusters/${clusterId}/cronhpas?name=${name}&namespace=${newNamespace}`;
  // let url = `/apis/autoscaling/v2beta1/namespaces/${newNamespace}/horizontalpodautoscalers/${name}`;
  let params: RequestParams = {
    method: Method.put,
    url,
    data: cronHpaData
  };

  try {
    let response = await reduceNetworkRequest(params, clusterId);
    if (response.code === 0) {
      return true;
    }
  } catch (error) {
    alertError(error, url);
  }
}

/**
 * 获取YAML
 */
export async function fetchCronHpaYaml({
  namespace,
  clusterId,
  name
}: {
  namespace: string;
  clusterId: string;
  name: string;
}) {
  const newNamespace = cutNsStartClusterId({ namespace, clusterId });
  let url = `/apis/platform.tkestack.io/v1/clusters/${clusterId}/cronhpas?name=${name}&namespace=${newNamespace}`;
  // let url = `/apis/autoscaling/v1/namespaces/${namespace}/horizontalpodautoscalers/${name}`;
  const userDefinedHeader = {
    Accept: 'application/yaml'
  };
  let params: RequestParams = {
    method: Method.get,
    url,
    userDefinedHeader
  };

  try {
    let response = await reduceNetworkRequest(params, clusterId);
    let yamlList = response.code === 0 ? [response.data] : [];

    const result: RecordSet<string> = {
      recordCount: yamlList.length,
      records: yamlList
    };

    return result;
  } catch (error) {
    alertError(error, url);
  }
}

/**
 * 更新YAML
 */
export async function modifyCronHpaYaml({
  namespace,
  clusterId,
  name,
  yamlData
}: {
  namespace: string;
  clusterId: string;
  name: string;
  yamlData: any;
}) {
  const newNamespace = cutNsStartClusterId({ namespace, clusterId });
  let url = `/apis/platform.tkestack.io/v1/clusters/${clusterId}/cronhpas?name=${name}&namespace=${newNamespace}`;
  // let url = `/apis/autoscaling/v1/namespaces/${namespace}/horizontalpodautoscalers/${name}`;
  const userDefinedHeader = {
    Accept: 'application/json',
    'Content-Type': 'application/yaml'
  };
  let params: RequestParams = {
    method: Method.put,
    url,
    userDefinedHeader,
    data: yamlData
  };

  try {
    let response = await reduceNetworkRequest(params, clusterId);
    if (response.code === 0) {
      return response;
    }
  } catch (error) {
    alertError(error, url);
  }
}

export async function fetchAddons({ clusterId }: { clusterId: string }) {
  let url = `/apis/platform.tkestack.io/v1/clusters/${clusterId}/addons`;
  let params: RequestParams = {
    method: Method.get,
    url
  };

  try {
    let response = await reduceNetworkRequest(params);
    const addons = {};
    if (response.code === 0) {
      response.data.items.forEach(item => {
        addons[item.spec.type] = item;
      });
    }
    return addons;
  } catch (error) {
    alertError(error, url);
  }
}
