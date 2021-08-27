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
import * as React from 'react';

import { Bubble, Badge } from '@tea/component';
import { t } from '@tencent/tea-app/lib/i18n';

export function YamlSearchHelperPanel(options: { isShow: boolean }) {
  let { isShow } = options;

  return (
    <Bubble
      visible={isShow}
      placement={'bottom'}
      content={
        <>
          <div style={{ display: 'flex', marginBottom: '3px' }}>
            <Badge style={{ width: '150px' }}>Ctrl-F / Cmd-F</Badge>
            <span>&nbsp;&nbsp;{t('开始搜索')}</span>
          </div>
          <div style={{ display: 'flex', marginBottom: '3px' }}>
            <Badge style={{ width: '150px' }}>Ctrl-G / Cmd-G</Badge>
            <span>&nbsp;&nbsp;{t('下一个')}</span>
          </div>
          <div style={{ display: 'flex', marginBottom: '3px' }}>
            <Badge style={{ width: '150px' }}>Shift-Ctrl-G / Shift-Cmd-G</Badge>
            <span>&nbsp;&nbsp;{t('上一个')}</span>
          </div>
          <div style={{ display: 'flex', marginBottom: '3px' }}>
            <Badge style={{ width: '150px' }}>Shift-Ctrl-F / Shift-Cmd-F</Badge>
            <span>&nbsp;&nbsp;{t('替换')}</span>
          </div>
          <div style={{ display: 'flex', marginBottom: '3px' }}>
            <Badge style={{ width: '150px' }}>Shift-Ctrl-R / Shift-Cmd-R</Badge>
            <span>&nbsp;&nbsp;{t('替换全部')}</span>
          </div>
          {/* <br />
          <Badge theme="success">Alt-F</Badge>
          {t('持久性搜索，对话框不会自动关闭，Enter键查找下一个，Shift-Enter键查找上一个')}
          <br />
          <Badge theme="success">Alt-G</Badge>
          {t('跳转到行')} */}
        </>
      }
    />
  );
}
