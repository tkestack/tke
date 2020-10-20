import * as JsYAML from 'js-yaml';
import React, { useState } from 'react';

import { Button, List, Modal, ExternalLink } from '@tea/component';
import { t } from '@tencent/tea-app/lib/i18n';

import { YamlEditorPanel, YamlSearchHelperPanel } from '..';

export function YamlDialog(options: {
  onClose: () => void;
  yamlConfig: string;
  isShow: boolean;
  title: string | JSX.Element;
}) {
  let { onClose, yamlConfig, isShow, title } = options;
  const cancel = () => {
    onClose && onClose();
  };
  return (
    <Modal visible={isShow} caption={title} onClose={cancel} disableEscape={true} size={700}>
      <Modal.Body>
        <YamlEditorPanel config={yamlConfig} />
      </Modal.Body>
      <Modal.Footer>
        <Button onClick={cancel}>{t('关闭')}</Button>
      </Modal.Footer>
    </Modal>
  );
}
