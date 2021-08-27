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
import * as React from 'react';

import { Button, Modal } from '@tea/component';
import { BaseReactProps } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

interface TipDialogProps extends BaseReactProps {
  /**是否显示 */
  isShow: boolean;

  /**提示框标题 */
  caption?: string;

  /**提示框主体 */
  body?: string | JSX.Element;

  /** footer的按钮 */
  footerButton?: JSX.Element;

  /**显示宽度 */
  width?: number;

  /**确定操作 */
  performAction?: (value?: any) => void;

  /**取消操作 */
  cancelAction?: (value?: any) => void;
}

export class TipDialog extends React.Component<TipDialogProps, {}> {
  render() {
    let { isShow, caption, body, width, performAction, cancelAction, children, footerButton } = this.props;
    const cancel = () => {
      cancelAction(false);
    };
    const perform = () => {
      performAction(false);
    };

    if (!isShow) {
      return <noscript />;
    }

    return (
      <Modal visible={true} caption={caption || t('提示')} onClose={cancel} size={width || 485} disableEscape={true}>
        <Modal.Body>{body || children}</Modal.Body>
        <Modal.Footer>
          {footerButton ? (
            footerButton
          ) : (
            <Button type="primary" onClick={perform}>
              {t('确定')}
            </Button>
          )}
        </Modal.Footer>
      </Modal>
    );
  }
}
