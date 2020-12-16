import Request from './request';
import { compare } from 'compare-versions';

function transObjToelector(obj: { [key: string]: string | number }) {
  return Object.entries(obj)
    .map(([key, value]) => `${key}${value && '='}${value}`)
    .join(',');
}

function sleep(time) {
  return new Promise(_ => setTimeout(_, time));
}

export const getK8sValidVersions = async (): Promise<Array<string>> => {
  const {
    data: { k8sValidVersions }
  } = await Request.get<any, { data: { k8sValidVersions: string } }>(
    '/api/v1/namespaces/kube-public/configmaps/cluster-info',
    {
      headers: {
        'X-TKE-ClusterName': 'global'
      }
    }
  );

  try {
    return JSON.parse(k8sValidVersions);
  } catch (error) {
    console.log('getK8sValidVersions error:', error, k8sValidVersions);
    return [];
  }
};

export const checkClusterIsNeedUpdate = async ({
  clusterName,
  clusterVersion
}: {
  clusterName: string;
  clusterVersion: string;
}) => {
  // 如果version不存在，说明集群在创建中
  if (!clusterVersion) {
    return {
      master: {
        isNeed: false,
        message: '集群创建中！'
      },

      worker: {
        isNeed: false,
        message: '集群创建中！'
      }
    };
  }
  // check machin
  const { items } = await Request.get<any, { items: Array<any> }>('/apis/platform.tkestack.io/v1/machines', {
    params: {
      labelSelector: transObjToelector({
        'platform.tkestack.io/need-upgrade': ''
      }),
      fieldSelector: transObjToelector({
        'spec.clusterName': clusterName
      })
    }
  });
  if (items.length > 0) {
    const master = {
      isNeed: false,
      message: '该集群有worker节点需要先完成升级'
    };

    let worker = {
      isNeed: true,
      message: ''
    };

    if (items.findIndex(({ status }) => status?.phase === 'Upgrading') > -1) {
      worker = {
        isNeed: false,
        message: '有节点正在升级中'
      };
    }

    return {
      master,

      worker
    };
  }

  // checkversion
  const canUpdateVersions = await getK8sValidVersions();
  if (canUpdateVersions.findIndex(v => compare(v, clusterVersion, '>')) < 0) {
    return {
      master: {
        isNeed: false,
        message: '没有可供升级的K8s版本'
      },

      worker: {
        isNeed: false,
        message: '无需升级'
      }
    };
  }

  return {
    master: {
      isNeed: true,
      message: ''
    },

    worker: {
      isNeed: false,
      message: '请先升级master'
    }
  };
};

interface UpdateClusterConfig {
  version: string;
  drainNodeBeforeUpgrade: boolean;
  maxUnready: number;
  autoMode: boolean;
  clusterName: string;
}
export const updateCluster = ({
  clusterName,
  version,
  drainNodeBeforeUpgrade,
  maxUnready,
  autoMode
}: UpdateClusterConfig) => {
  return Request.patch(`/apis/platform.tkestack.io/v1/clusters/${clusterName}`, {
    spec: {
      version,
      features: {
        upgrade: {
          mode: autoMode ? 'Auto' : 'Manual',
          strategy: autoMode
            ? {
                drainNodeBeforeUpgrade,
                maxUnready: maxUnready + '%'
              }
            : undefined
        }
      }
    }
  });
};

export const updateSingleWorker = async (mchineName: string) => {
  await Request.patch(`/apis/platform.tkestack.io/v1/machines/${mchineName}`, {
    status: {
      phase: 'Upgrading'
    }
  });
};

export const addUpgradeLabelToMachine = async (mchineName: string) => {
  await Request.patch(`/apis/platform.tkestack.io/v1/machines/${mchineName}`, {
    metadata: {
      labels: {
        'platform.tkestack.io/need-upgrade': 'willUpgrade'
      }
    }
  });
};

export const updateWorkers = async ({
  mchineNames,
  maxUnready,
  drainNodeBeforeUpgrade,
  clusterName
}: {
  mchineNames: Array<string>;
  maxUnready: number;
  clusterName: string;
  drainNodeBeforeUpgrade: boolean;
}) => {
  await Request.patch(`/apis/platform.tkestack.io/v1/clusters/${clusterName}`, {
    spec: {
      features: {
        upgrade: {
          strategy: {
            maxUnready: maxUnready + '%',
            drainNodeBeforeUpgrade
          }
        }
      }
    }
  });

  await Promise.all(mchineNames.map(n => addUpgradeLabelToMachine(n)));

  await updateSingleWorker(mchineNames[0]);
};

export const getNodes = async ({ clusterName, clusterVersion }: { clusterName: string; clusterVersion }) => {
  const machinesPromise = Request.get<any, { items: Array<any> }>('/apis/platform.tkestack.io/v1/machines', {
    params: {
      labelSelector: transObjToelector({
        'platform.tkestack.io/need-upgrade': ''
      }),
      fieldSelector: transObjToelector({
        'spec.clusterName': clusterName
      })
    }
  });

  const nodePromise = Request.get<any, { items: Array<any> }>('/api/v1/nodes', {
    headers: {
      'X-TKE-ClusterName': clusterName
    }
  });

  const [{ items: machines }, { items: nodes }] = await Promise.all([machinesPromise, nodePromise]);

  return nodes
    .map(({ metadata, status }) => ({
      key: metadata?.uid,
      name: metadata?.name,
      label: metadata?.labels,
      kubeletVersion: status.nodeInfo.kubeletVersion,
      clusterVersion,
      phase: status.conditions.find(item => item.type === 'Ready')?.status === 'True' ? 'Running' : 'Failed',
      machines: machines
        .filter(m => metadata?.labels?.['platform.tkestack.io/machine-ip'] === m?.spec?.ip)
        .map(m => m?.metadata?.name)
    }))
    .filter(({ machines }) => machines.length > 0);
};
