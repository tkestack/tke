import React, { useEffect, useState } from 'react';
import { Table, Justify, Button, Select, TagSearchBox, Text } from 'tea-component';
import { useRecoilValueLoadable, useRecoilState } from 'recoil';
import { namespaceListState, namespaceSelectionState } from '../../store/base';
import { router } from '@src/modules/cluster/router';

export const VMListActionPanel = ({ route, reFetch }) => {
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

            <Button icon="refresh" onClick={reFetch} />
          </>
        }
      />
    </Table.ActionPanel>
  );
};
