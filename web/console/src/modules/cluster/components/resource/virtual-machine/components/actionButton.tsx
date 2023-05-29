import React, { useState } from 'react';
import { Button, PopConfirm, ButtonProps, Modal } from 'tea-component';

interface ActionButtonProps {
  type?: ButtonProps['type'];
  title: string;
  confirm?: () => Promise<any>;
  onSuccess?: () => void;
  children: React.ReactNode;
  disabled?: boolean;
  body?: React.ReactNode;
}

export const ActionButton = ({
  type = 'primary',
  onSuccess = () => {},
  title,
  confirm,
  children,
  disabled = false,
  body
}: ActionButtonProps) => {
  const [visible, setVisible] = useState(false);

  async function handleOk() {
    setVisible(false);

    await confirm();

    onSuccess();
  }

  return (
    <>
      <Button type={type} onClick={() => setVisible(true)} disabled={disabled}>
        {children}
      </Button>

      <Modal caption={title} visible={visible} onClose={() => setVisible(false)}>
        {body && <Modal.Body>{body}</Modal.Body>}

        <Modal.Footer>
          <Button type="primary" onClick={handleOk}>
            确定
          </Button>

          <Button onClick={() => setVisible(false)}>取消</Button>
        </Modal.Footer>
      </Modal>
    </>
  );
};
