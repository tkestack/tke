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
import { t } from '@tencent/tea-app/lib/i18n';
import { DetailField, DetailInfo } from '../../../src/modules/common/models';
import {
  commonActionField,
  commonDetailInfo,
  dataFormatConfig,
  defaulNotExistedValue,
  generateResourceInfo
} from '../common';

const displayField = Object.assign(
  {},
  {
    name: {
      dataField: ['metadata.name'],
      dataFormat: dataFormatConfig['text'],
      width: '10%',
      headTitle: t('名称'),
      noExsitedValue: '-',
      isLink: true, // 用于判断该值是否为链接
      isClip: true
    },
    schedule: {
      dataField: ['spec.schedule'],
      dataFormat: dataFormatConfig['text'],
      width: '10%',
      headTitle: t('执行策略'),
      noExsitedValue: '-'
    },
    parallelism: {
      dataField: ['spec.jobTemplate.spec.parallelism'],
      dataFormat: dataFormatConfig['text'],
      width: '8%',
      headTitle: t('并行度'),
      noExsitedValue: defaulNotExistedValue
    },
    completions: {
      dataField: ['spec.jobTemplate.spec.completions'],
      width: '8%',
      dataFormat: dataFormatConfig['text'],
      headTitle: t('重复次数'),
      noExsitedValue: defaulNotExistedValue
    },
    operator: {
      dataField: [''],
      dataFormat: dataFormatConfig['operator'],
      width: '15%',
      headTitle: t('操作'),
      tips: '',
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
        },
        {
          name: t('更新调度策略'),
          actionType: 'modifyNodeAffinity',
          isInMoreOp: true
        }
      ]
    }
  }
);

/** resource action 当中的配置 */
const actionField = Object.assign({}, commonActionField);

/** 自定义配置详情basic info 的展示 */
const detailBasicInfo: DetailInfo = {
  info: {
    metadata: {
      dataField: ['metadata'],
      displayField: {
        name: {
          dataField: ['name'],
          dataFormat: dataFormatConfig['text'],
          label: t('名称'),
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
        description: {
          dataField: ['annotations.description'],
          dataFormat: dataFormatConfig['text'],
          label: t('描述'),
          noExsitedValue: defaulNotExistedValue
        },
        createdTime: {
          dataField: ['creationTimestamp'],
          dataFormat: dataFormatConfig['time'],
          label: t('创建时间'),
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
        schedule: {
          dataField: ['schedule'],
          dataFormat: dataFormatConfig['text'],
          label: t('执行策略'),
          tips: '',
          noExsitedValue: defaulNotExistedValue
        },
        concurrencyPolicy: {
          dataField: ['concurrencyPolicy'],
          dataFormat: dataFormatConfig['text'],
          label: t('并发策略'),
          tips: '',
          noExsitedValue: defaulNotExistedValue
        },
        parallelism: {
          dataField: ['jobTemplate.spec.parallelism'],
          dataFormat: dataFormatConfig['text'],
          label: t('并行度'),
          tips: '',
          noExsitedValue: defaulNotExistedValue
        },
        completions: {
          dataField: ['jobTemplate.spec.completions'],
          dataFormat: dataFormatConfig['text'],
          label: t('重复次数'),
          tips: '',
          noExsitedValue: defaulNotExistedValue
        }
      }
    }
  }
};

const tabList = [
  {
    id: 'event',
    label: t('事件')
  },
  {
    id: 'info',
    label: t('详情')
  },
  {
    id: 'yaml',
    label: 'YAML'
  }
];

/** 详情页面的相关配置 */
const detailField: DetailField = {
  tabList,
  detailInfo: Object.assign({}, commonDetailInfo('cronjob'), detailBasicInfo)
};

/** cronJobs 的配置 */
export const cronjobs = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'cronjob',
    isRelevantToNamespace: true,
    requestType: {
      list: 'cronjobs'
    },
    displayField,
    actionField,
    detailField
  });
};
