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
import { ICComponter } from '../../models';
import { Justify, Button, Text, Bubble } from '@tencent/tea-component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { FormPanel } from '@tencent/ff-component';

export function ShowICComputerPanel({
  computer,
  canEdit,
  editTips,
  onEdit,
  onDelete
}: {
  computer: ICComponter;
  canEdit?: boolean;
  editTips?: React.ReactNode;
  onEdit: () => void;
  onDelete: () => void;
}) {
  //如果使用了footer，需要在下方留出足够的空间，避免重叠

  let ipList = computer.ipList.split(';');
  return (
    <FormPanel
      fixed
      isNeedCard={false}
      labelStyle={{
        minWidth: 460
      }}
      fieldStyle={{
        minWidth: 100
      }}
    >
      <FormPanel.Item
        label={
          <Text theme="strong" style={{ lineHeight: '26px' }}>
            {ipList.slice(0, 3).join(';')}
            {ipList.length > 3 ? '等' + ipList.length + '台机器' : ''}
          </Text>
        }
      >
        <Justify
          right={
            <React.Fragment>
              <Bubble content={canEdit ? null : t('请先完成待编辑项')}>
                <Button icon="pencil" disabled={!canEdit} onClick={onEdit} />
              </Bubble>
              <Button icon="close" onClick={onDelete} />
            </React.Fragment>
          }
        />
      </FormPanel.Item>
    </FormPanel>
  );
}
