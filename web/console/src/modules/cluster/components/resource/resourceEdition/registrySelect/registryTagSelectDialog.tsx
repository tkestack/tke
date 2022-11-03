import React, { useState } from 'react';
import { Button, Modal, Table, Checkbox, TableColumn, Text } from 'tea-component';

const { scrollable, radioable } = Table.addons;

export const RegistryTagSelectDialog = ({ tags, onConfirm }) => {
  const [visible, setVisible] = useState(false);

  const [selectedTag, setSelectedTag] = useState(null);

  const columns: TableColumn[] = [
    {
      key: 'name',
      header: '镜像版本'
    },

    {
      key: 'digest',
      header: '摘要(SHA256)',
      render({ digest }) {
        return <Text copyable>{digest}</Text>;
      }
    }
  ];

  return (
    <>
      <Button type="link" onClick={() => setVisible(true)}>
        选择镜像版本
      </Button>

      <Modal caption="选择镜像版本" visible={visible} onClose={() => setVisible(false)}>
        <Modal.Body>
          <Table
            columns={columns}
            records={tags}
            recordKey="name"
            addons={[
              scrollable({
                maxHeight: 500
              }),

              radioable({
                value: selectedTag,
                onChange(key) {
                  setSelectedTag(key);
                },
                rowSelect: true
              })
            ]}
          />
        </Modal.Body>

        <Modal.Footer>
          <Button
            disabled={!selectedTag}
            onClick={() => {
              onConfirm(selectedTag);
              setVisible(false);
            }}
            type="primary"
          >
            确认
          </Button>

          <Button onClick={() => setVisible(false)}>取消</Button>
        </Modal.Footer>
      </Modal>
    </>
  );
};
