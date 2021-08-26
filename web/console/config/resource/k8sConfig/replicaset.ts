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

import { DetailField, DetailInfo } from '../../../src/modules/common/models';
import {
  commonDisplayField,
  commonActionField,
  defaulNotExistedValue,
  dataFormatConfig,
  commonDetailInfo,
  commonTabList,
  generateResourceInfo
} from '../common';

const displayField = Object.assign({}, commonDisplayField, {});

/** resource action当中的配置 */
const actionField = Object.assign({}, commonActionField);

/** 自定义basic info的配置 */
const detailBasicInfo: DetailInfo = {
  info: {
    metadata: {
      dataField: ['metadata'],
      displayField: {
        name: {
          dataField: ['name'],
          dataFormat: dataFormatConfig['text'],
          label: 'Name',
          tips: '',
          noExsitedValue: defaulNotExistedValue
        },
        namespace: {
          dataField: ['namespace'],
          dataFormat: dataFormatConfig['text'],
          label: 'Namespace',
          tips: '',
          noExsitedValue: defaulNotExistedValue
        },
        createdTime: {
          dataField: ['creationTimestamp'],
          dataFormat: dataFormatConfig['time'],
          label: 'Create Time',
          tips: '',
          noExsitedValue: defaulNotExistedValue
        },
        label: {
          dataField: ['labels'],
          dataFormat: dataFormatConfig['labels'],
          label: 'Labels',
          tips: '',
          noExsitedValue: defaulNotExistedValue
        }
      }
    },
    spec: {
      dataField: ['spec'],
      displayField: {
        selector: {
          dataField: ['selector.matchLabels'],
          dataFormat: dataFormatConfig['labels'],
          label: 'Selector',
          tips: '',
          noExsitedValue: defaulNotExistedValue
        },
        updateStrategy: {
          dataField: ['updateStrategy.type'],
          dataFormat: dataFormatConfig['text'],
          label: 'UpdateStrategy',
          tips: '',
          noExsitedValue: defaulNotExistedValue
        }
      }
    }
  }
};

/** 详情页面的相关配置 */
const detailField: DetailField = {
  tabList: [...commonTabList],
  detailInfo: Object.assign({}, commonDetailInfo(), detailBasicInfo)
};

/** replica set的配置 */
export const rs = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'rs',
    isRelevantToNamespace: true,
    requestType: {
      list: 'replicasets'
    },
    displayField,
    actionField,
    detailField
  });
};
