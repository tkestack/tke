import { DetailField, DetailInfo } from '../../../src/modules/common/models';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import {
  commonDisplayField,
  defaulNotExistedValue,
  commonActionField,
  commonDetailInfo,
  dataFormatConfig,
  workloadCommonTabList,
  generateResourceInfo
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
