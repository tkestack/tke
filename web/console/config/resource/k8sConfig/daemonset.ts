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
  commonDisplayField,
  dataFormatConfig,
  defaulNotExistedValue,
  generateResourceInfo,
  workloadCommonTabList
} from '../common';

/** displayField，列表展示的细节 */
const displayField = Object.assign({}, commonDisplayField, {
  runningReplicas: {
    dataField: ['status.numberReady', 'status.desiredNumberScheduled'],
    dataFormat: dataFormatConfig['replicas'],
    width: '20%',
    headTitle: t('运行/期望Pod数量'),
    noExsitedValue: '0'
  },
  operator: {
    dataField: [''],
    dataFormat: dataFormatConfig['operator'],
    width: '20%',
    headTitle: t('操作'),
    tips: '',
    operatorList: [
      {
        name: t('更新镜像'),
        actionType: 'modifyRegistry',
        isInMoreOp: false
      },
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
        name: t('设置更新策略'),
        actionType: 'modifyStrategy',
        isInMoreOp: true
      },
      {
        name: t('更新调度策略'),
        actionType: 'modifyNodeAffinity',
        isInMoreOp: true
      }
    ]
  }
});

/** resource action 当中的配置 */
const actionField = Object.assign({}, commonActionField);

/** 需要自定义Basic info的展示 */
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
          label: t('更新策略'),
          tips: '',
          noExsitedValue: defaulNotExistedValue
        },
        networkType: {
          dataField: ['template', 'metadata', 'annotations', 'k8s.v1.cni.cncf.io/networks'],
          dataFormat: dataFormatConfig['text'],
          label: t('网络模式'),
          noExsitedValue: '-'
        }
      }
    },
    status: {
      dataField: ['status'],
      displayField: {
        desiredNumberScheduled: {
          dataField: ['desiredNumberScheduled'],
          dataFormat: dataFormatConfig['text'],
          label: t('期望副本数'),
          tips: '',
          noExsitedValue: '0'
        },
        currentNumberScheduled: {
          dataField: ['desiredNumberScheduled'],
          dataFormat: dataFormatConfig['text'],
          label: t('运行副本数'),
          tips: '',
          noExsitedValue: '0'
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

/** daemon set 的配置 */
export const daemonset = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'daemonset',
    isRelevantToNamespace: true,
    requestType: {
      list: 'daemonsets'
    },
    displayField,
    actionField,
    detailField
  });
};
