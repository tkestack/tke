import { atom, selector } from 'recoil';
import { virtualMachineAPI, storageClassAPI } from '@src/webApi';
import { clusterIdState } from './base';
import { v4 as uuidv4 } from 'uuid';
import { DiskInterface, DiskTypeEnum, VolumeModeEnum, ActionTypeEnum } from '../constants';
export { clusterIdState, namespaceListState, namespaceSelectionState } from './base';

const createKey = id => `cluster/virtual-machine/creation/${id}`;

export const mirrorListState = selector({
  key: createKey('mirrorList'),
  get: async ({ get }) => {
    const clusterId = get(clusterIdState);

    const { items } = await virtualMachineAPI.fetchVMMirrorList(clusterId);

    return (
      items?.map(({ metadata }) => ({
        text: metadata?.annotations?.['tkestack.io/image-display-name'],
        value: metadata?.name,
        tooltip: metadata?.annotations?.['tkestack.io/image-display-name']
      })) ?? []
    );
  }
});

export const diskListState = atom<DiskInterface[]>({
  key: createKey('diskList'),
  default: [
    {
      id: uuidv4(),
      name: 'rootfs',
      type: DiskTypeEnum.System,
      volumeMode: VolumeModeEnum.Filesystem,
      storageClass: null,
      scProvisioner: null,
      size: 50
    },

    {
      id: uuidv4(),
      name: 'datavolume1',
      type: DiskTypeEnum.Data,
      volumeMode: VolumeModeEnum.Filesystem,
      storageClass: null,
      scProvisioner: null,
      size: 50
    }
  ]
});

export const storageClassListState = selector({
  key: createKey('storageClassList'),
  get: async ({ get }) => {
    const clusterId = get(clusterIdState);

    if (!clusterId) return [];

    const { items } = await storageClassAPI.fetchStorageClassList(clusterId);

    return items?.map(({ metadata, provisioner }) => ({ value: metadata?.name, provisioner })) ?? [];
  }
});

export const diskListValidateState = selector({
  key: createKey('diskListValidate'),
  get: ({ get }) => {
    const diskList = get(diskListState);

    return diskList.map(item => valiData(item, diskValidate));
  }
});

const diskValidate = {
  name: {
    required: '磁盘名称必填!',
    maxLength: 63,
    pattern: /^[a-z]([-a-z0-9]*[a-z0-9])?$/
  },
  storageClass: {
    required: '存储类必选!'
  }
};

function valiData(data, rules) {
  return Object.entries(rules).reduce<{ [key: string]: { status: 'error' | null; message: string } }>(
    (all, [dataKey, rule]: [string, any]) => {
      const value = data[dataKey];

      let rsp = {
        status: null,
        message: null
      };

      if (rule.required && !value) {
        rsp = {
          status: 'error',
          message: rule.required
        };
      }

      if (rule.maxLength && value.length > rule.maxLength) {
        rsp = {
          status: 'error',
          message: `最大长度为${rule.maxLength}!`
        };
      }

      if (rule.pattern && !rule.pattern.test(value)) {
        rsp = {
          status: 'error',
          message: '格式不正确!'
        };
      }

      return {
        ...all,

        [dataKey]: rsp
      };
    },
    {}
  );
}
