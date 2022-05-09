import React, { useMemo } from 'react';
import { Radio, Card, Form, Select, Text, Checkbox, Icon } from 'tea-component';
import { useFetch } from '@src/modules/common/hooks';
import { virtualMachineAPI } from '@src/webApi';
import { Namespace } from '@src/modules/alarmPolicy/models';

interface IEditAlarmPolicyVMProps {
  clusterId: string;
  namespaceList: Namespace[];
  type: string; // 'all' | 'part'
  setType: (type: string) => void;
  namespaceSelection: string;
  setNamespaceSelection: (namespace: string) => void;
  vmSelections: string[];
  setVmSelections: (vmSelections: string[]) => void;
}

export const EditAlarmPolicyVM = ({
  clusterId,
  namespaceList,
  type,
  setType,
  namespaceSelection,
  setNamespaceSelection,
  vmSelections,
  setVmSelections
}: IEditAlarmPolicyVMProps) => {
  const { data: vmList, status: vmListFetchStatus } = useFetch(
    async () => {
      const rsp = await virtualMachineAPI.fetchVMList({ clusterId, namespace: namespaceSelection }, {});

      return {
        data: rsp?.items?.map(item => item?.metadata?.name) ?? []
      };
    },
    [clusterId, namespaceSelection],
    {
      fetchAble: !!clusterId && !!namespaceSelection && namespaceSelection !== 'ALL'
    }
  );

  const namespaceOptions = useMemo(() => {
    return (
      namespaceList?.map(({ name, displayName }) => ({
        value: name,
        text: displayName
      })) ?? []
    );
  }, [namespaceList]);

  return (
    <Radio.Group
      value={type}
      onChange={(type: 'all' | 'part') => {
        if (type === 'all') setNamespaceSelection('all');

        if (type === 'part') setNamespaceSelection('default');

        setType(type);
      }}
    >
      <Card style={{ backgroundColor: '#f2f2f2', boxShadow: 'none', minWidth: 450 }}>
        <Card.Body>
          <Radio name="part">选择指定虚拟机</Radio>
          {type === 'part' && (
            <Form>
              <Form.Item label={<Text theme="text">Namespace</Text>}>
                <Select
                  appearance="button"
                  matchButtonWidth
                  size="m"
                  searchable
                  value={namespaceSelection}
                  options={namespaceOptions}
                  onChange={ns => setNamespaceSelection(ns)}
                />
              </Form.Item>

              <Form.Item label={<Text theme="text">虚拟机</Text>}>
                <Card bordered>
                  <Card.Body style={{ maxHeight: 150, overflow: 'auto' }}>
                    {vmListFetchStatus === 'loading' && <Icon type="loading" />}
                    <Checkbox.Group layout="column" value={vmSelections} onChange={values => setVmSelections(values)}>
                      {vmList?.map(name => (
                        <Checkbox key={name} name={name}>
                          {name}
                        </Checkbox>
                      )) ?? null}
                    </Checkbox.Group>
                  </Card.Body>
                </Card>
              </Form.Item>
            </Form>
          )}
        </Card.Body>
      </Card>

      <Card style={{ backgroundColor: '#f2f2f2', boxShadow: 'none', minWidth: 450 }}>
        <Card.Body>
          <Radio name="all">选择全部</Radio>

          {type === 'all' && (
            <Form>
              <Form.Item label={<Text theme="text">Namespace</Text>}>
                <Select
                  appearance="button"
                  matchButtonWidth
                  size="m"
                  searchable
                  value={namespaceSelection}
                  options={[
                    {
                      value: 'ALL',
                      text: 'ALL'
                    },
                    ...namespaceOptions
                  ]}
                  onChange={ns => setNamespaceSelection(ns)}
                />
              </Form.Item>

              <Form.Item label={<Text theme="text">虚拟机</Text>} extra="包括后续新增的虚拟机">
                <Select
                  disabled
                  appearance="button"
                  matchButtonWidth
                  size="m"
                  searchable
                  value="all"
                  options={[{ value: 'all', text: 'All' }]}
                />
              </Form.Item>
            </Form>
          )}
        </Card.Body>
      </Card>
    </Radio.Group>
  );
};
