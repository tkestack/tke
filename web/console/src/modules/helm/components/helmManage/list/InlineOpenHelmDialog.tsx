/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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
import { RootProps } from '../../HelmApp';
import { Button, Modal } from '@tea/component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
export class InlineOpenHelmDialog extends React.Component<RootProps, {}> {
  render() {
    let { actions } = this.props;
    const cancel = () => {};

    const submit = () => {};
    return (
      <Modal visible={true} caption={t('确认开通Helm应用管理？')} onClose={cancel} disableEscape={true}>
        <Modal.Body>
          <div style={{ fontSize: '14px', lineHeight: '20px' }}>
            <div className="docker-dialog jiqun">
              <p>
                {t(
                  '开通Helm应用，将在集群内安装Helm tiller组件，占用集群计算资源，目前您所选的集群尚未开通，是否立即开通？'
                )}
              </p>
            </div>
          </div>
        </Modal.Body>
        <Modal.Footer>
          <Button type="primary" onClick={submit}>
            {t('确定')}
          </Button>
          <Button onClick={cancel}>{t('取消')}</Button>
        </Modal.Footer>
      </Modal>
    );
  }
}
