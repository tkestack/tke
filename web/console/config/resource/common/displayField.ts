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

import { DisplayField } from '../../../src/modules/common/models';
import { dataFormatConfig } from './dataFormat';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

/** resource table 当中展示的数据 */
export const commonDisplayField: DisplayField = {
  check: {
    dataField: [],
    dataFormat: dataFormatConfig['checker'],
    width: '16px',
    headTitle: ' ',
    noExsitedValue: '-'
  },
  name: {
    dataField: ['metadata.name'],
    dataFormat: dataFormatConfig['text'],
    width: '20%',
    headTitle: t('名称'),
    noExsitedValue: '-',
    isLink: true, // 用于判断该值是否为链接
    isClip: true
  },
  labels: {
    dataField: ['metadata.labels'],
    dataFormat: dataFormatConfig['labels'],
    width: '15%',
    headTitle: 'Labels',
    noExsitedValue: t('无')
  },
  selector: {
    dataField: ['spec.selector.matchLabels'],
    dataFormat: dataFormatConfig['labels'],
    width: '20%',
    headTitle: 'Selector',
    noExsitedValue: t('无')
  }
};
