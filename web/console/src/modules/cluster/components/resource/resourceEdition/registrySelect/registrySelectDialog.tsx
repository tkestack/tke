import React, { useMemo, useState } from 'react';
import { Button, Modal, Table, SearchBox, TableColumn, Text } from 'tea-component';
import { fetchRepositoryList } from '@src/webApi/registry';
import { useRequest } from 'ahooks';
import { fetchDockerRegUrl } from '@src/modules/registry/WebAPI';

const { filterable, scrollable, radioable } = Table.addons;

const ALL_VALUE = '__ALL__';

export const RegistrySelectDialog = ({ onConfirm }) => {
  const [visible, setVisible] = useState(false);

  const [selectedRepo, setSelectedRepo] = useState(null);

  const [searchValue, setSearchValue] = useState('');

  const [visibilityType, setVisibilityType] = useState(ALL_VALUE);

  const [selectedNamespaceList, setSelectedNamespaceList] = useState([]);

  const { data: repoBaseUrl } = useRequest(async () => {
    const res = await fetchDockerRegUrl();

    console.log('repoBaseUrl---->', res);

    return res;
  });

  const { data: _repositoryList } = useRequest(async () => {
    const res = await fetchRepositoryList();

    return res?.items ?? [];
  });

  const repositoryList = useMemo(() => {
    return (
      _repositoryList
        ?.filter(repo => repo?.spec?.name?.includes(searchValue))
        ?.filter(repo => visibilityType === ALL_VALUE || repo?.spec?.visibility === visibilityType)
        ?.filter(
          repo => selectedNamespaceList?.length === 0 || selectedNamespaceList?.includes(repo?.spec?.namespaceName)
        ) ?? []
    );
  }, [_repositoryList, searchValue, visibilityType, selectedNamespaceList]);

  const namespaceOptions = useMemo(() => {
    return _repositoryList?.map(repo => ({ value: repo?.spec?.namespaceName })) ?? [];
  }, [_repositoryList]);

  const columns: TableColumn[] = [
    {
      key: 'spec.name',
      header: '名称'
    },

    {
      key: 'spec.visibility',
      header: '类型'
    },

    {
      key: 'spec.namespaceName',
      header: '命名空间'
    },

    {
      key: 'spec.resourceVersion',
      header: '仓库地址',
      render(repo) {
        return <Text overflow copyable>{`${repoBaseUrl}/${repo?.spec?.namespaceName}/${repo?.spec?.name}`}</Text>;
      }
    }
  ];

  return (
    <>
      <Button type="link" onClick={() => setVisible(true)}>
        选择镜像
      </Button>

      <Modal caption="选择镜像" visible={visible} onClose={() => setVisible(false)}>
        <Modal.Body>
          <Table.ActionPanel>
            <SearchBox value={searchValue} onChange={value => setSearchValue(value)} />
          </Table.ActionPanel>

          <Table
            columns={columns}
            records={repositoryList}
            recordKey="metadata.name"
            addons={[
              filterable({
                type: 'single',
                column: 'spec.visibility',
                value: visibilityType,
                onChange: value => setVisibilityType(value),
                all: {
                  value: ALL_VALUE,
                  text: '全部'
                },
                // 选项列表
                options: [
                  { value: 'Public', text: '公有' },
                  { value: 'Private', text: '私有' }
                ]
              }),

              filterable({
                type: 'multiple',
                column: 'spec.namespaceName',
                value: selectedNamespaceList,
                onChange: value => {
                  setSelectedNamespaceList(value);
                },
                all: {
                  value: ALL_VALUE,
                  text: '全部'
                },
                options: namespaceOptions
              }),

              scrollable({
                maxHeight: 500
              }),

              radioable({
                value: selectedRepo,
                onChange(key) {
                  setSelectedRepo(key);
                },
                rowSelect: true
              })
            ]}
          />
        </Modal.Body>

        <Modal.Footer>
          <Button
            disabled={!selectedRepo}
            type="primary"
            onClick={() => {
              const repo = _repositoryList?.find(item => item?.metadata?.name === selectedRepo);

              const registry = `${repoBaseUrl}/${repo?.spec?.namespaceName}/${repo?.spec?.name}`;
              onConfirm({ registry, tags: repo?.status?.tags ?? null });
              setVisible(false);
            }}
          >
            确认
          </Button>

          <Button onClick={() => setVisible(false)}>取消</Button>
        </Modal.Footer>
      </Modal>
    </>
  );
};
