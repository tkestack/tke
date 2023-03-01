import React, { useMemo } from 'react';
import { Table, TableColumn, Input, Select, InputNumber, Button, Form } from 'tea-component';
import { diskListState, storageClassListState, diskListValidateState } from '../../store/creation';
import { useRecoilState, useRecoilValueLoadable, useRecoilValue } from 'recoil';
import { DiskTypeEnum, VolumeModeOptions, VolumeModeEnum } from '../../constants';
import { v4 as uuidv4 } from 'uuid';

export const DiskPanel = () => {
  const [diskList, setDiskList] = useRecoilState(diskListState);
  const storageClassListLoadable = useRecoilValueLoadable(storageClassListState);
  const diskListValidate = useRecoilValue(diskListValidateState);

  function delDiskItem(id) {
    const newDiskList = diskList.filter(item => item.id !== id);
    setDiskList(newDiskList);
  }

  function addDiskItem() {
    const newDiskList = [
      ...diskList,
      {
        id: uuidv4(),
        name: '',
        type: DiskTypeEnum.Data,
        volumeMode: VolumeModeEnum.Filesystem,
        storageClass: null,
        scProvisioner: null,
        size: 50
      }
    ];
    setDiskList(newDiskList);
  }

  function modifyDiskItem(newItem) {
    const newDiskList = diskList?.map(item => {
      return item.id === newItem.id
        ? {
            ...item,
            ...newItem
          }
        : item;
    });

    setDiskList(newDiskList);
  }

  const columns: TableColumn[] = [
    {
      key: 'name',
      header: '磁盘名称',
      render({ name, id }, _, index) {
        return (
          <Form.Control {...diskListValidate?.[index]?.name}>
            <Input size="s" value={name} onChange={name => modifyDiskItem({ name, id })} />
          </Form.Control>
        );
      }
    },

    {
      key: 'type',
      header: '磁盘类型',
      render({ type }) {
        return type === DiskTypeEnum.System ? '系统盘' : '数据盘';
      }
    },

    // {
    //   key: 'volumeMode',
    //   header: '卷模式',
    //   render({ volumeMode, id }) {
    //     return (
    //       <Form.Control>
    //         <Select
    //           size="s"
    //           value={volumeMode}
    //           options={VolumeModeOptions}
    //           onChange={volumeMode => modifyDiskItem({ volumeMode, id })}
    //         />
    //       </Form.Control>
    //     );
    //   }
    // },

    {
      key: 'storageClass',
      header: '存储类',
      render({ storageClass, id }, _, index) {
        return (
          <Form.Control {...diskListValidate?.[index]?.storageClass}>
            <Select
              type="simulate"
              searchable
              appearence="button"
              size="s"
              value={storageClass}
              options={storageClassListLoadable?.state === 'hasValue' ? storageClassListLoadable?.contents : []}
              onChange={storageClass =>
                modifyDiskItem({
                  storageClass,
                  scProvisioner:
                    storageClassListLoadable?.contents?.find(({ value }) => value === storageClass)?.provisioner ??
                    null,
                  id
                })
              }
            />
          </Form.Control>
        );
      }
    },

    {
      key: 'size',
      header: '容量（Gi）',
      render({ size, id }) {
        return (
          <Form.Control>
            <InputNumber
              hideButton
              value={size}
              min={1}
              step={1}
              onChange={size =>
                modifyDiskItem({
                  size,
                  id
                })
              }
            />
          </Form.Control>
        );
      }
    },

    {
      key: 'actions',
      header: '操作',
      render({ type, id }) {
        return (
          <Button type="link" disabled={type === DiskTypeEnum.System} onClick={() => delDiskItem(id)}>
            删除
          </Button>
        );
      }
    }
  ];

  return (
    <>
      <Table columns={columns} records={diskList} recordKey="id" />
      <Button style={{ width: '100%', margin: '20px 0' }} onClick={addDiskItem}>
        +新增
      </Button>
    </>
  );
};
