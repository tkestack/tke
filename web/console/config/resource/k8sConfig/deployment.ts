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
import { DetailField, DetailInfo, DisplayField } from '../../../src/modules/common/models';
import { cloneDeep } from '../../../src/modules/common/utils';
import {
  commonActionField,
  commonDetailInfo,
  commonDisplayField,
  dataFormatConfig,
  defaulNotExistedValue,
  generateResourceInfo,
  workloadCommonTabList
} from '../common';

/** resource table 当中展示的数据
 * commonDisplayField 使用公共的展示
 * {}，deployment自定义的数据
 */
const userDefinedDisplayField: DisplayField = {
  runningReplicas: {
    dataField: ['status.readyReplicas', 'status.replicas'],
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
        name: t('更新Pod数量'),
        actionType: 'modifyPod',
        isInMoreOp: false
      },
      {
        name: t('更新镜像'),
        actionType: 'modifyRegistry',
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
      },
      {
        name: t('编辑YAML'),
        actionType: 'modify',
        isInMoreOp: true
      },
      {
        name: t('删除'),
        actionType: 'delete',
        isInMoreOp: true
      }
    ]
  }
};

const displayField = Object.assign({}, commonDisplayField, userDefinedDisplayField);

/** resource action当中的配置 */
const actionField = Object.assign({}, commonActionField);

const detailBasicInfo: DetailInfo = {
  info: {
    metadata: {
      dataField: ['metadata'],
      displayField: {
        name: {
          dataField: ['name'],
          dataFormat: dataFormatConfig['text'],
          label: t('名称'),
          noExsitedValue: defaulNotExistedValue
        },
        namespace: {
          dataField: ['namespace'],
          dataFormat: dataFormatConfig['text'],
          label: 'Namespace',
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
          dataField: ['strategy.type'],
          dataFormat: dataFormatConfig['text'],
          label: t('更新策略'),
          tips: '',
          noExsitedValue: defaulNotExistedValue
        },
        replicas: {
          dataField: ['replicas'],
          dataFormat: dataFormatConfig['text'],
          label: t('副本数'),
          noExsitedValue: '0'
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
        readyReplicas: {
          dataField: ['readyReplicas'],
          dataFormat: dataFormatConfig['text'],
          label: t('运行副本数'),
          noExsitedValue: '0'
        }
      }
    }
  }
};

const tabList = cloneDeep(workloadCommonTabList);
tabList.splice(1, 0, {
  id: 'history',
  label: t('修订历史')
});
/** 详情页面的相关配置 */
const detailField: DetailField = {
  tabList,
  detailInfo: Object.assign({}, commonDetailInfo(), detailBasicInfo)
};

/** deployment的配置 */
export const deployment = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'deployment',
    requestType: {
      list: 'deployments'
    },
    isRelevantToNamespace: true,
    displayField,
    actionField,
    detailField
  });
};
