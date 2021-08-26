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

import { DetailField, DisplayField, DetailInfo } from '../../../src/modules/common/models';
import {
  commonActionField,
  defaulNotExistedValue,
  dataFormatConfig,
  workloadCommonTabList,
  generateResourceInfo
} from '../common';
import { cloneDeep } from '../../../src/modules/common/utils';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

const displayField: DisplayField = {
  name: {
    dataField: ['metadata.name'],
    dataFormat: dataFormatConfig['text'],
    width: '10%',
    headTitle: 'Name',
    noExsitedValue: defaulNotExistedValue
  },
  InternalIP: {
    dataField: ['status.addresses.0.address'],
    dataFormat: dataFormatConfig['ip'],
    width: '10%',
    headTitle: 'IP',
    noExsitedValue: defaulNotExistedValue
  },
  podCIDR: {
    dataField: ['spec.podCIDR'],
    dataFormat: dataFormatConfig['ip'],
    width: '10%',
    headTitle: 'podCIDR',
    noExsitedValue: defaulNotExistedValue
  },
  labels: {
    dataField: ['metadata.labels'],
    dataFormat: dataFormatConfig['labels'],
    width: '10%',
    headTitle: 'Labels',
    noExsitedValue: defaulNotExistedValue
  },
  kubeletVersion: {
    dataField: ['status.nodeInfo.kubeletVersion'],
    dataFormat: dataFormatConfig['text'],
    width: '10%',
    headTitle: 'kubeletVersion',
    noExsitedValue: defaulNotExistedValue
  },
  os: {
    dataField: ['status.nodeInfo.osImage'],
    dataFormat: dataFormatConfig['text'],
    width: '10%',
    headTitle: 'osImage',
    noExsitedValue: defaulNotExistedValue
  }
  // operator: {
  //     dataField: [''],
  //     dataFormat: dataFormatConfig['operator'],
  //     width: '10%',
  //     headTitle: t('操作'),
  //     operatorList: [
  //         {
  //             name: t('驱逐'),
  //             actionType: 'drain',
  //             isInMoreOp: false
  //         }
  //     ]
  // }
};

/** resrouce action当中的配置 */
const actionField = Object.assign({}, commonActionField, {
  create: {
    isAvailable: false
  }
});

/** 自定义配置详情的展示 */
const detailBasicInfo: DetailInfo = {};

/** 自定义tabList */
let tabList: any[] = cloneDeep(workloadCommonTabList);
tabList.splice(2, 1);

/** 详情页面的相关配置 */
const detailField: DetailField = {
  tabList: [
    {
      id: 'pod',
      label: t('Pod管理')
    },
    {
      id: 'event',
      label: t('事件')
    },
    {
      id: 'nodeInfo',
      label: t('详情')
    },
    {
      id: 'yaml',
      label: 'YAML'
    }
  ],
  detailInfo: Object.assign({}, detailBasicInfo)
};

/** apiVersion的配置 */
export const node = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'node',
    requestType: {
      list: 'nodes'
    },
    displayField,
    actionField,
    detailField
  });
};
