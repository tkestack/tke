import React, { useState } from 'react';
import { Button, Modal, Form, Input } from 'tea-component';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { nameRule } from '@config/validateConfig';
import { z } from 'zod';
import { getReactHookFormStatusWithMessage } from '@helper';
import { virtualMachineAPI } from '@src/webApi';

export const CreateSnapshotButton = ({ clusterId, namespace, name, onSuccess = () => {}, disabled }) => {
  const [visible, setVisible] = useState(false);

  const { control, handleSubmit, reset } = useForm({
    mode: 'onBlur',
    defaultValues: {
      snapshotName: ''
    },
    resolver: zodResolver(z.object({ snapshotName: nameRule('快照名称') }))
  });

  function onCancel() {
    setVisible(false);

    reset();
  }

  async function onSubmit({ snapshotName }) {
    console.log('snapshot name', snapshotName);

    try {
      await virtualMachineAPI.createSnapshot({
        clusterId,
        namespace,
        vmName: name,
        name: snapshotName
      });

      onCancel();
      onSuccess();
    } catch (error) {
      console.log('create snapshot error:', error);
    }
  }

  return (
    <>
      <Button type="link" onClick={() => setVisible(true)} disabled={disabled}>
        新建快照
      </Button>

      <Modal caption="新建虚拟机快照" visible={visible} onClose={onCancel}>
        <Modal.Body>
          <Form>
            <Controller
              control={control}
              name="snapshotName"
              render={({ field, ...other }) => (
                <Form.Item
                  label="快照名称"
                  extra="快照名称不能超过63个字符"
                  {...getReactHookFormStatusWithMessage(other)}
                >
                  <Input {...field} />
                </Form.Item>
              )}
            />
          </Form>
        </Modal.Body>

        <Modal.Footer>
          <Button type="primary" onClick={handleSubmit(onSubmit)}>
            确认
          </Button>
          <Button onClick={onCancel}>取消</Button>
        </Modal.Footer>
      </Modal>
    </>
  );
};
