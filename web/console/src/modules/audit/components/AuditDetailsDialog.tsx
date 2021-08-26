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
import { useDispatch, useSelector } from 'react-redux';

import { Modal, Button, Input } from '@tea/component';
import { bindActionCreators, insertCSS } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { allActions } from '../actions';
import { router } from '../router';
import { Audit } from '../models';
const { useState, useEffect } = React;

export const AuditDetailsDialog = (props: { isShowing: boolean; toggle: () => void; record: Audit }) => {
  const state = useSelector(state => state);
  const dispatch = useDispatch();
  const { actions } = bindActionCreators({ actions: allActions }, dispatch);
  const { isShowing, toggle, record } = props;
  return (
    <Modal visible={isShowing} caption={t('记录详情')} onClose={toggle}>
      <Modal.Body>
        <Input
          multiline
          size="full"
          value={JSON.stringify(record, null, 4)}
          readonly={true}
          style={{ height: '300px' }}
        />
      </Modal.Body>
      <Modal.Footer>
        <Button type="primary" onClick={toggle}>
          <Trans>确定</Trans>
        </Button>
      </Modal.Footer>
    </Modal>
  );
};
