import { VolumeModeEnum } from '@src/modules/cluster/components/resource/virtual-machine/constants';
import { Request, generateQueryString } from './request';
import { v4 as uuid } from 'uuid';

const IMAGE_NAMESPACE = 'kube-public';

const basePath = (clusterId: string) => `/apis/platform.tkestack.io/v1/clusters/${clusterId}/proxy?path=`;
const header = (clusterId: string) => ({
  headers: {
    'X-TKE-ClusterName': clusterId
  }
});

export function fetchVMList({ clusterId, namespace }, { limit = null, continueToken = null }) {
  const path = encodeURIComponent(
    `/apis/kubevirt.io/v1/namespaces/${namespace}/virtualmachines?${generateQueryString({
      limit,
      continue: continueToken
    })}`
  );
  return Request.get<any, { items: any[]; metadata: any }>(
    `/apis/platform.tkestack.io/v1/clusters/${clusterId}/proxy?path=${path}`,
    {
      headers: {
        'X-TKE-ClusterName': clusterId
      }
    }
  );
}

export function fetchVM({ clusterId, namespace, name }) {
  return Request.get<any, any>(
    `/apis/platform.tkestack.io/v1/clusters/${clusterId}/proxy?path=/apis/kubevirt.io/v1/namespaces/${namespace}/virtualmachines/${name}`,
    {
      headers: {
        'X-TKE-ClusterName': clusterId
      }
    }
  );
}

export function fetchVMForYaml({ clusterId, namespace, name }) {
  return Request.get<any, any>(
    `/apis/platform.tkestack.io/v1/clusters/${clusterId}/proxy?path=/apis/kubevirt.io/v1/namespaces/${namespace}/virtualmachines/${name}`,
    {
      headers: {
        'X-TKE-ClusterName': clusterId,
        Accept: 'application/yaml'
      }
    }
  );
}

export function fetchVMMirrorList(clusterId) {
  const labelSelector = encodeURIComponent(`tkestack.io/image-os-type`);
  return Request.get<any, { items: any }>(
    `/api/v1/namespaces/${IMAGE_NAMESPACE}/persistentvolumeclaims?labelSelector=${labelSelector}`,
    {
      headers: {
        'X-TKE-ClusterName': clusterId
      }
    }
  );
}

export async function fetchVMI({ clusterId, namespace, name }) {
  try {
    return await Request.get<any, any>(
      `/apis/platform.tkestack.io/v1/clusters/${clusterId}/proxy?path=/apis/kubevirt.io/v1/namespaces/${namespace}/virtualmachineinstances/${name}`,
      {
        headers: {
          'X-TKE-ClusterName': clusterId
        }
      }
    );
  } catch (error) {
    console.log('fetchVMI error:', error);

    return null;
  }
}

export async function fetchVMIList({ clusterId, namespace }) {
  return await Request.get<any, any>(
    `/apis/platform.tkestack.io/v1/clusters/${clusterId}/proxy?path=/apis/kubevirt.io/v1/namespaces/${namespace}/virtualmachineinstances`,
    {
      headers: {
        'X-TKE-ClusterName': clusterId
      }
    }
  );
}

export async function fetchVMListWithVMI({ clusterId, namespace }, { limit, continueToken, query }) {
  let items, metadata;

  if (query) {
    const vm = await fetchVM({ clusterId, namespace, name: query });
    items = vm ? [vm] : [];
  } else {
    const rsp = await fetchVMList({ clusterId, namespace }, { limit, continueToken });

    items = rsp?.items ?? [];
    metadata = rsp?.metadata;
  }

  return {
    items: await Promise.all(
      items?.map(async item => {
        const vmi = await fetchVMI({ clusterId, namespace, name: item.metadata.name });

        return {
          ...item,
          vmi
        };
      }) ?? []
    ),

    newContinueToken: metadata?.continue,

    restCount: metadata?.remainingItemCount ?? 0
  };
}

export function createVM({
  namespace,
  clusterId,
  vmOptions: { name, description, cpu, memory, mirror, diskList, networkMode }
}) {
  const body = {
    apiVersion: 'kubevirt.io/v1',
    kind: 'VirtualMachine',
    metadata: {
      name,
      annotations: {
        'kubevirt.io/latest-observed-api-version': 'v1',
        'kubevirt.io/storage-observed-api-version': 'v1alpha3',
        'tkestack.io/image-display-name': mirror.text,
        'tkestack.io/support-snapshot': diskList.every(item => item?.scProvisioner === 'rbd.csi.ceph.com')
          ? 'true'
          : 'false',
        description
      },
      labels: {
        'kubevirt.io/domain': name
      }
    },

    spec: {
      running: true,

      template: {
        metadata: {
          annotations: {
            'tkestack.io/image-display-name': mirror.text
          },

          labels: {
            'kubevirt.io/domain': name
          }
        },

        spec: {
          domain: {
            cpu: {
              cores: cpu
            },

            devices: {
              blockMultiQueue: true,
              disks: diskList.map((item, index) => ({
                disk: {
                  bus: 'virtio'
                },
                bootOrder: index + 1,
                name: item.name,
                cache: 'writeback'
              })),

              interfaces: [
                {
                  model: 'virtio',
                  name: 'default',
                  bridge: {}
                }
              ],

              inputs: [
                {
                  bus: 'usb',
                  name: 'tablet',
                  type: 'tablet'
                }
              ]
            },

            machine: {
              type: 'q35'
            },

            resources: {
              requests: {
                memory: `${memory}G`
              }
            }
          },

          networks: [
            {
              name: 'default',
              pod: {}
            }
          ],

          volumes: diskList.map(item => ({
            name: item.name,
            dataVolume: {
              name: `${item.name}.${name}`
            }
          }))
        }
      },

      dataVolumeTemplates: diskList.map((item, index) => ({
        metadata: {
          name: `${item.name}.${name}`
        },

        spec: {
          pvc: {
            accessModes: ['ReadWriteMany'],
            resources: {
              requests: {
                storage: `${item.size}Gi`
              }
            },

            volumeMode: ['loopdevice.csi.infra.tce.io', 'rbd.csi.ceph.com'].includes(item?.scProvisioner)
              ? VolumeModeEnum.Block
              : VolumeModeEnum.Filesystem,
            storageClassName: item.storageClass
          },

          source:
            index === 0
              ? {
                  pvc: {
                    name: mirror.value,
                    namespace: IMAGE_NAMESPACE
                  }
                }
              : {
                  blank: {}
                }
        }
      }))
    }
  };

  return Request.post(
    `/apis/platform.tkestack.io/v1/clusters/${clusterId}/proxy?path=/apis/kubevirt.io/v1/namespaces/${namespace}/virtualmachines`,
    body,
    {
      headers: {
        'X-TKE-ClusterName': clusterId
      }
    }
  );
}

export const fetchVMDetail = async ({ clusterId, namespace, name }) => {
  const [vm, vmi] = await Promise.all([
    fetchVM({ clusterId, namespace, name }),
    fetchVMI({ clusterId, namespace, name })
  ]);

  return {
    vm,
    vmi
  };
};

export const setVMRunningStatus = async ({ clusterId, namespace, name }, status: boolean) => {
  return Request.patch(
    `/apis/platform.tkestack.io/v1/clusters/${clusterId}/proxy?path=/apis/kubevirt.io/v1/namespaces/${namespace}/virtualmachines/${name}`,
    {
      spec: {
        running: status
      }
    },
    {
      headers: {
        'X-TKE-ClusterName': clusterId,
        'Content-Type': 'application/merge-patch+json'
      }
    }
  );
};

export const deleteVM = async ({ clusterId, namespace, name }) => {
  return Request.delete(
    `/apis/platform.tkestack.io/v1/clusters/${clusterId}/proxy?path=/apis/kubevirt.io/v1/namespaces/${namespace}/virtualmachines/${name}`,
    {
      headers: {
        'X-TKE-ClusterName': clusterId
      }
    }
  );
};

export const fetchEventList = async ({ namespace, clusterId, name }) => {
  const fieldSelector = encodeURIComponent(`involvedObject.name=${name},involvedObject.apiVersion=kubevirt.io/v1`);
  return Request.get<any, any>(`/api/v1/namespaces/${namespace}/events?fieldSelector=${fieldSelector}`, {
    headers: {
      'X-TKE-ClusterName': clusterId
    }
  });
};

export const checkVmEnable = async clusterId => {
  try {
    await Request.get(
      `/apis/platform.tkestack.io/v1/clusters/${clusterId}/proxy?path=/apis/apiextensions.k8s.io/v1/customresourcedefinitions/virtualmachines.kubevirt.io`,
      {
        headers: {
          'X-TKE-ClusterName': clusterId
        }
      }
    );

    return true;
  } catch (error) {
    return false;
  }
};

export const createSnapshot = async ({ clusterId, namespace, vmName, name }) => {
  return Request.post(
    `/apis/platform.tkestack.io/v1/clusters/${clusterId}/proxy?path=/apis/snapshot.kubevirt.io/v1alpha1/namespaces/${namespace}/virtualmachinesnapshots`,
    {
      apiVersion: 'snapshot.kubevirt.io/v1alpha1',
      kind: 'VirtualMachineSnapshot',
      metadata: {
        name: name
      },
      spec: {
        source: {
          apiGroup: 'kubevirt.io',
          kind: 'VirtualMachine',
          name: vmName
        }
      }
    },
    {
      headers: {
        'X-TKE-ClusterName': clusterId
      }
    }
  );
};

export function fetchSnapshotItem({ clusterId, namespace, name }) {
  return Request.get(
    `${basePath(clusterId)}/apis/snapshot.kubevirt.io/v1alpha1/namespaces/${namespace}/virtualmachinesnapshots/${name}`,
    header(clusterId)
  );
}

export async function fetchSnapshotList(
  { clusterId, namespace },
  { limit = null, continueToken = null, query = null }
) {
  if (query) {
    const item = await fetchSnapshotItem({ clusterId, namespace, name: query });

    return {
      items: item ? [item] : []
    };
  }

  const path = `/apis/snapshot.kubevirt.io/v1alpha1/namespaces/${namespace}/virtualmachinesnapshots?${generateQueryString(
    {
      limit,
      continue: continueToken
    }
  )}`;

  return Request.get<any, any>(`${basePath(clusterId)}${path}`, {
    headers: {
      'X-TKE-ClusterName': clusterId
    }
  });
}

export function delSnapshot({ clusterId, namespace, name }) {
  return Request.delete(
    `${basePath(clusterId)}/apis/snapshot.kubevirt.io/v1alpha1/namespaces/${namespace}/virtualmachinesnapshots/${name}`,
    {
      headers: {
        'X-TKE-ClusterName': clusterId
      }
    }
  );
}

export function recoverySnapshot({ clusterId, namespace, name, vmName }) {
  const body = {
    apiVersion: 'snapshot.kubevirt.io/v1alpha1',
    kind: 'VirtualMachineRestore',
    metadata: {
      name: `${name}-restore-${uuid()}`
    },
    spec: {
      target: {
        apiGroup: 'kubevirt.io',
        kind: 'VirtualMachine',
        name: vmName
      },
      virtualMachineSnapshotName: name
    }
  };

  return Request.post(
    `${basePath(clusterId)}/apis/snapshot.kubevirt.io/v1alpha1/namespaces/${namespace}/virtualmachinerestores`,
    body,
    header(clusterId)
  );
}

export function fetchRecoveryStoreList({ clusterId, namespace }) {
  return Request.get<any, any>(
    `${basePath(clusterId)}/apis/snapshot.kubevirt.io/v1alpha1/namespaces/${namespace}/virtualmachinerestores`,
    header(clusterId)
  );
}
