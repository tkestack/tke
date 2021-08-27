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
import {
  defaulNotExistedValue,
  commonActionField,
  commonDetailInfo,
  workloadCommonTabList,
  dataFormatConfig,
  generateResourceInfo
} from '../common';
import { DetailField, DetailInfo } from '../../../src/modules/common/models';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

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
    labels: {
      dataField: ['metadata.labels'],
      dataFormat: dataFormatConfig['labels'],
      width: '10%',
      headTitle: 'Labels',
      noExsitedValue: t('无')
    },
    selector: {
      dataField: ['spec.selector.matchLabels'],
      dataFormat: dataFormatConfig['labels'],
      width: '10%',
      headTitle: 'Selector',
      noExsitedValue: t('无')
    },
    parallelism: {
      dataField: ['spec.parallelism'],
      dataFormat: dataFormatConfig['text'],
      width: '8%',
      headTitle: t('并行度'),
      noExsitedValue: defaulNotExistedValue
    },
    completions: {
      dataField: ['spec.completions'],
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
        }
      ]
    }
  }
);

/** resource action 当中的配置 */
const actionField = Object.assign({}, commonActionField);

/** 自定义配置详情basic info的展示 */
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
        parallelism: {
          dataField: ['parallelism'],
          dataFormat: dataFormatConfig['text'],
          label: t('并行度'),
          tips: '',
          noExsitedValue: defaulNotExistedValue
        },
        completions: {
          dataField: ['completions'],
          dataFormat: dataFormatConfig['text'],
          label: t('重复次数'),
          tips: '',
          noExsitedValue: defaulNotExistedValue
        }
      }
    },
    status: {
      dataField: ['status'],
      displayField: {
        startTime: {
          dataField: ['startTime'],
          dataFormat: dataFormatConfig['time'],
          label: t('启动时间'),
          tips: '',
          noExsitedValue: defaulNotExistedValue
        }
      }
    }
  }
};

/** 详情页面的相关配置 */
const detailField: DetailField = {
  tabList: [...workloadCommonTabList],
  detailInfo: Object.assign({}, commonDetailInfo(), detailBasicInfo)
};

/** jobs 的配置 */
export const jobs = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'job',
    isRelevantToNamespace: true,
    requestType: {
      list: 'jobs'
    },
    displayField,
    actionField,
    detailField
  });
};
