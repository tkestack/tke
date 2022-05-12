import React from 'react';
import { Table, Justify, Button, Select, TagSearchBox, Text } from 'tea-component';
import { useRecoilValueLoadable, useRecoilState } from 'recoil';
import { namespaceListState, namespaceSelectionState } from '../../store/base';
import { router } from '@src/modules/cluster/router';

export const VMListActionPanel = ({ route, reFetch, onQueryChange }) => {
  const namespaceListLoadable = useRecoilValueLoadable(namespaceListState);
  const [namespaceSelection, setNamespaceSelection] = useRecoilState(namespaceSelectionState);

  return (
    <Table.ActionPanel>
      <Justify
        left={
          <Button
            type="primary"
            onClick={() => {
              const urlParams = router.resolve(route);
              router.navigate(Object.assign({}, urlParams, { mode: 'create' }), route.queries);
            }}
          >
            新建
          </Button>
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
              onChange={value => setNamespaceSelection(value)}
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
