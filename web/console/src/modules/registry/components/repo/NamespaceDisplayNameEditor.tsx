import React, { useState } from 'react';
import { Button, Modal, Form, TextArea } from 'tea-component';
import { registryApi } from '@src/webApi';

interface INamespaceDisplayNameEditorProps {
  value: string;
  name: string;
  onSuccess: () => void;
}

export const NamespaceDisplayNameEditor = ({ value, name, onSuccess }: INamespaceDisplayNameEditorProps) => {
  const [visible, _setVisible] = useState(false);

  function setVisible(visible: boolean) {
    _setVisible(visible);

    setDisplayName(value);
  }

  const [displayName, setDisplayName] = useState(value);

  async function handleSubmit() {
    await registryApi.modifyNamespaceDisplayName({ name, displayName });

    onSuccess();

    setVisible(false);
  }

  return (
    <>
      <Button type="icon" icon="pencil" onClick={() => setVisible(true)} />

      <Modal visible={visible} caption="编辑描述信息" onClose={() => setVisible(false)}>
        <Modal.Body>
          <Form>
            <Form.Item label="描述" required>
              <TextArea size="full" value={displayName} onChange={value => setDisplayName(value.trim())} />
            </Form.Item>
          </Form>
        </Modal.Body>

        <Modal.Footer>
          <Button type="primary" disabled={!displayName} onClick={handleSubmit}>
            确定
          </Button>

          <Button onClick={() => setVisible(false)}>取消</Button>
        </Modal.Footer>
      </Modal>
    </>
  );
};
