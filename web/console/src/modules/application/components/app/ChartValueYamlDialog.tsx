import * as JsYAML from 'js-yaml';
import * as React from 'react';

import { Button, List, Modal, Alert } from '@tea/component';
import { t } from '@tencent/tea-app/lib/i18n';

import { YamlEditorPanel } from '../../../common/components';

export function ChartValueYamlDialog(options: {
  onChange: (value) => void;
  onClose: () => void;
  yamlConfig: string;
  isShow: boolean;
}) {
  let { onChange, onClose, yamlConfig, isShow } = options;
  let [validator, setvalidator] = React.useState({
    result: 0,
    message: ''
  });

  const cancel = () => {
    onClose && onClose();
  };
  const save = () => {
    try {
      JsYAML.safeLoad(yamlConfig);
      setvalidator({ result: 0, message: '' });
      onChange && onChange(yamlConfig);
      onClose && onClose();
    } catch (error) {
      setvalidator({ result: 2, message: t('Yaml格式错误') });
    }
  };

  const _handleForInputEditor = (config: string) => {
    onChange && onChange(config);
  };
  return (
    <Modal visible={isShow} caption={t('参数')} onClose={cancel} disableEscape={true} size={700}>
      <Modal.Body>
        <YamlEditorPanel config={yamlConfig} handleInputForEditor={_handleForInputEditor} />
      </Modal.Body>
      <Modal.Footer>
        <Button type="primary" onClick={save}>
          {validator.result === 2 ? t('重试') : t('保存')}
        </Button>
        <Button onClick={cancel}>{t('取消')}</Button>
        {validator.result === 2 && (
          <Alert type="error" style={{ marginTop: 8 }}>
            {validator.message}
          </Alert>
        )}
      </Modal.Footer>
    </Modal>
  );
}
