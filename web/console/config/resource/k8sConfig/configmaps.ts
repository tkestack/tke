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
import { DetailField, DisplayField, DetailInfo } from '../../../src/modules/common/models';
import { defaulNotExistedValue, dataFormatConfig, commonActionField, generateResourceInfo } from '../common';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

const displayField: DisplayField = {
  name: {
    dataField: ['metadata.name'],
    dataFormat: dataFormatConfig['text'],
    width: '25%',
    headTitle: t('名称'),
    noExsitedValue: defaulNotExistedValue,
    isLink: true,
    isClip: true
  },
  labels: {
    dataField: ['metadata.labels'],
    dataFormat: dataFormatConfig['labels'],
    width: '25%',
    headTitle: 'Labels',
    noExsitedValue: defaulNotExistedValue
  },
  creationTimestamp: {
    dataField: ['metadata.creationTimestamp'],
    dataFormat: dataFormatConfig['time'],
    width: '20%',
    headTitle: t('创建时间'),
    noExsitedValue: defaulNotExistedValue
  },
  operator: {
    dataField: [''],
    dataFormat: dataFormatConfig['operator'],
    width: '15%',
    headTitle: t('操作'),
    operatorList: [
      {
        name: t('编辑YAML'),
        actionType: 'modify',
        isInMoreOp: false
      },
      {
        name: t('删除'),
        actionType: 'delete',
        isInMoreOp: false
      }
    ]
  }
};

/** resource action当中的配置 */
const actionField = Object.assign({}, commonActionField);

/** 自定义tablist */
const tabList = [
  {
    id: 'info',
    label: t('详情')
  },
  {
    id: 'yaml',
    label: 'YAML'
  }
];

const detailBasicInfo: DetailInfo = {
  info: {
    metadata: {
      dataField: ['metadata'],
      displayField: {
        name: {
          dataField: ['name'],
          dataFormat: dataFormatConfig['text'],
          label: t('名称'),
          noExsitedValue: defaulNotExistedValue,
          order: '0'
        },
        namespace: {
          dataField: ['namespace'],
          dataFormat: dataFormatConfig['text'],
          label: 'Namespace',
          noExsitedValue: defaulNotExistedValue,
          order: '5'
        },
        labels: {
          dataField: ['labels'],
          dataFormat: dataFormatConfig['labels'],
          label: 'Labels',
          noExsitedValue: defaulNotExistedValue,
          order: '10'
        },
        createTime: {
          dataField: ['creationTimestamp'],
          dataFormat: dataFormatConfig['time'],
          label: t('创建时间'),
          tips: '',
          noExsitedValue: defaulNotExistedValue,
          order: '15'
        }
      }
    }
  }
};

const detailField: DetailField = {
  tabList,
  detailInfo: Object.assign({}, detailBasicInfo)
};

/** configmaps的配置 */
export const configmap = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'configmap',
    isRelevantToNamespace: true,
    requestType: {
      list: 'configmaps'
    },
    displayField,
    actionField,
    detailField
  });
};
