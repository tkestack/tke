/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */
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
