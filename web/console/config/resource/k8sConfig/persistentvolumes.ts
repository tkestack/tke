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
import { defaulNotExistedValue, dataFormatConfig, commonActionField, generateResourceInfo } from '../common';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

const displayField: DisplayField = {
  name: {
    dataField: ['metadata.name'],
    dataFormat: dataFormatConfig['text'],
    width: '10%',
    headTitle: t('名称'),
    noExsitedValue: defaulNotExistedValue,
    isLink: true,
    isClip: true
  },
  status: {
    dataField: ['status.phase'],
    dataFormat: dataFormatConfig['status'],
    width: '10%',
    headTitle: t('状态'),
    noExsitedValue: defaulNotExistedValue
  },
  accessModes: {
    dataField: ['spec.accessModes.0'],
    dataFormat: dataFormatConfig['mapText'],
    mapTextConfig: {
      ReadWriteOnce: 'RWO',
      ReadOnlyMany: 'ROX',
      ReadWriteMany: 'RWX'
    },
    width: '10%',
    headTitle: t('访问权限'),
    noExsitedValue: defaulNotExistedValue
  },
  persistentVolumeReclaimPolicy: {
    dataField: ['spec.persistentVolumeReclaimPolicy'],
    dataFormat: dataFormatConfig['text'],
    width: '10%',
    headTitle: t('回收策略'),
    noExsitedValue: defaulNotExistedValue
  },
  pvc: {
    dataField: ['spec.claimRef.name'],
    dataFormat: dataFormatConfig['text'],
    width: '10%',
    headTitle: 'PVC',
    noExsitedValue: defaulNotExistedValue,
    isClip: true
  },
  storageClassName: {
    dataField: ['spec.storageClassName'],
    dataFormat: dataFormatConfig['text'],
    width: '10%',
    headTitle: 'StorageClass',
    noExsitedValue: defaulNotExistedValue,
    isClip: true
  },
  creationTimestamp: {
    dataField: ['metadata.creationTimestamp'],
    dataFormat: dataFormatConfig['time'],
    width: '13%',
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

/** 自定义tabllist */
const tabList = [
  {
    id: 'info',
    label: t('详情')
  },
  {
    id: 'event',
    label: t('事件')
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
        labels: {
          dataField: ['labels'],
          dataFormat: dataFormatConfig['labels'],
          label: 'Labels',
          noExsitedValue: defaulNotExistedValue,
          order: '5'
        },
        createTime: {
          dataField: ['creationTimestamp'],
          dataFormat: dataFormatConfig['time'],
          label: t('创建时间'),
          noExsitedValue: defaulNotExistedValue,
          order: '30'
        }
      }
    },
    spec: {
      dataField: ['spec'],
      displayField: {
        accessModes: {
          dataField: ['accessModes.0'],
          dataFormat: dataFormatConfig['text'],
          label: t('访问权限'),
          noExsitedValue: defaulNotExistedValue,
          order: '15'
        },
        pvc: {
          dataField: ['claimRef.name'],
          dataFormat: dataFormatConfig['text'],
          label: 'PVC',
          noExsitedValue: defaulNotExistedValue,
          order: '16'
        },
        storageClassName: {
          dataField: ['storageClassName'],
          dataFormat: dataFormatConfig['text'],
          label: 'StorageClass',
          noExsitedValue: defaulNotExistedValue,
          order: '20'
        },
        storage: {
          dataField: ['capacity.storage'],
          dataFormat: dataFormatConfig['text'],
          label: 'Storage',
          noExsitedValue: defaulNotExistedValue,
          order: '20'
        },
        persistentVolumeReclaimPolicy: {
          dataField: ['persistentVolumeReclaimPolicy'],
          dataFormat: dataFormatConfig['text'],
          label: t('回收策略'),
          noExsitedValue: defaulNotExistedValue,
          order: '25'
        }
      }
    },
    status: {
      dataField: ['status'],
      displayField: {
        phase: {
          dataField: ['phase'],
          dataFormat: dataFormatConfig['status'],
          label: t('状态'),
          noExsitedValue: defaulNotExistedValue,
          order: '10'
        }
      }
    }
  }
};

const detailField: DetailField = {
  tabList,
  detailInfo: Object.assign({}, detailBasicInfo)
};

/** persistentvolumes 的配置 */
export const pv = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'pv',
    requestType: {
      list: 'persistentvolumes'
    },
    displayField,
    actionField,
    detailField
  });
};
