import React, { useEffect, useState } from 'react';
import { Table, Justify, Button, Select, TagSearchBox, Text } from 'tea-component';
import { useRecoilValueLoadable, useRecoilState, useRecoilValue } from 'recoil';
import { namespaceListState, namespaceSelectionState, clusterIdState } from '../../store/base';
import { router } from '@src/modules/cluster/router';
import { VmMonitorDialog } from '../../components';

export const VMListActionPanel = ({ route, reFetch, vmList, onQueryChange }) => {
  const namespaceListLoadable = useRecoilValueLoadable(namespaceListState);
  const [namespaceSelection, setNamespaceSelection] = useRecoilState(namespaceSelectionState);

  const clusterId = useRecoilValue(clusterIdState);

  return (
    <Table.ActionPanel>
      <Justify
        left={
          <>
            <Button
              type="primary"
              onClick={() => {
                const urlParams = router.resolve(route);
                router.navigate(Object.assign({}, urlParams, { mode: 'create' }), route.queries);
              }}
            >
              新建
            </Button>

            <VmMonitorDialog clusterId={clusterId} namespace={namespaceSelection} vmList={vmList} />

            <Button
              type="primary"
              onClick={() => {
                const urlParams = router.resolve(route);
                router.navigate(Object.assign({}, urlParams, { mode: 'snapshot' }), route.queries);
              }}
            >
              快照
            </Button>
          </>
        }
        right={
          <>
            <Text reset theme="label" verticalAlign="middle">
              命名空间
            </Text>
            <Select
              type="simulate"
              searchable
              appearence="button"
              size="s"
              style={{ width: '130px', marginRight: '5px' }}
              value={namespaceSelection}
              options={
                namespaceListLoadable?.state === 'hasValue'
                  ? namespaceListLoadable?.contents?.map(value => ({ value }))
                  : []
              }
              onChange={value => {
                setNamespaceSelection(value);
                const urlParams = router.resolve(route);
                router.navigate(urlParams, Object.assign({}, route.queries, { np: value }));
              }}
            />

            <TagSearchBox
              tips=""
              style={{ width: 300 }}
              attributes={[
                {
                  type: 'input',
                  key: 'name',
                  name: '虚拟机名称'
                }
              ]}
              onChange={tags => {
                const name = tags?.find(item => item?.attr?.key === 'name')?.values?.[0]?.name ?? '';

                onQueryChange(name);
              }}
            />

            <Button icon="refresh" onClick={reFetch} />
          </>
        }
      />
    </Table.ActionPanel>
  );
};
