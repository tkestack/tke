import { atom, selector } from 'recoil';
import { namespaceAPI, virtualMachineAPI } from '@src/webApi';

const createKey = id => `cluster/virtual-machine/base/${id}`;

export const clusterIdState = atom({
  key: createKey('clusterId'),
  default: null
});

export const namespaceListState = selector({
  key: createKey('namespaceList'),
  get: async ({ get }) => {
    const clusterId = get(clusterIdState);

    if (!clusterId) return [];

    const { items } = await namespaceAPI.fetchNamespaceList(clusterId);

    return items?.map(item => item?.metadata?.name) ?? [];
  }
});

export const namespaceSelectionState = atom({
  key: createKey('namespaceSelection'),
  default: 'default'
});

export const vmSelectionState = atom({
  key: createKey('vmSelection'),
  default: null
});
