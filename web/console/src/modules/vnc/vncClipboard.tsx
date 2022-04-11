import React, { useState, useMemo } from 'react';
import { Button, Modal, Alert, Form, Input, FormItemProps } from 'tea-component';

interface IVncClipboardProps {
  onConfirm: (content: string) => void;
}

export const VncClipboard = ({ onConfirm }: IVncClipboardProps) => {
  const [visible, setVisible] = useState(false);
  const [content, setContent] = useState('');

  const contentValidate = useMemo<{ message: string; status: FormItemProps['status'] }>(() => {
    const rule = /^[^\u4e00-\u9fa5]+$/;

    if (content && !rule.test(content)) {
      return {
        message: '暂时不支持汉字！',
        status: 'error'
      };
    }

    return {
      message: null,
      status: null
    };
  }, [content]);

  function handleConfirm() {
    if (!content) {
      return;
    }

    setVisible(false);
    onConfirm(content);
    setContent('');
  }

  return (
    <>
      <Button type="link" style={{ textDecoration: 'underline #fff', color: '#fff' }} onClick={() => setVisible(true)}>
        这里
      </Button>

      <Modal caption="粘贴命令" visible={visible} onClose={() => setVisible(false)}>
        <Modal.Body>
          <Alert>将内容粘贴至文本框,暂不支持中文等非标准键盘值特殊字符</Alert>

          <Form>
            <Form.Item required label="命令文本内容" {...contentValidate}>
              <Input.TextArea size="full" value={content} onChange={value => setContent(value)} />
            </Form.Item>
          </Form>
        </Modal.Body>

        <Modal.Footer>
          <Button type="primary" disabled={!content || contentValidate?.status === 'error'} onClick={handleConfirm}>
            确定
          </Button>

          <Button onClick={() => setVisible(false)}>取消</Button>
        </Modal.Footer>
      </Modal>
    </>
  );
};
