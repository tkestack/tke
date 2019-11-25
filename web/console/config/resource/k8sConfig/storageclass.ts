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
  provisioner: {
    dataField: ['provisioner'],
    dataFormat: dataFormatConfig['text'],
    width: '18%',
    headTitle: t('来源')
  },
  // type: {
  //   dataField: ['parameters.type'],
  //   dataFormat: dataFormatConfig['mapText'],
  //   mapTextConfig: {
  //     CLOUD_BASIC: t('普通云硬盘'),
  //     CLOUD_PREMIUM: t('高性能云硬盘'),
  //     CLOUD_SSD: t('SSD云硬盘'),
  //     cbs: t('普通云硬盘')
  //   },
  //   width: '10%',
  //   headTitle: t('云盘类型'),
  //   noExsitedValue: defaulNotExistedValue
  // },
  // paymode: {
  //   dataField: ['parameters.paymode'],
  //   dataFormat: dataFormatConfig['mapText'],
  //   mapTextConfig: {
  //     POSTPAID: t('按量计费'),
  //     PREPAID: t('包年包月')
  //   },
  //   width: '10%',
  //   headTitle: t('计费模式'),
  //   noExsitedValue: t('按量计费')
  // },
  reclaimPolicy: {
    dataField: ['reclaimPolicy'],
    dataFormat: dataFormatConfig['text'],
    width: '10%',
    headTitle: t('回收策略'),
    noExsitedValue: 'Delete'
  },
  creationTimestamp: {
    dataField: ['metadata.creationTimestamp'],
    dataFormat: dataFormatConfig['time'],
    width: '10%',
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
        name: t('删除'),
        actionType: 'delete',
        isInMoreOp: false
      }
    ]
  }
};

/** resource action当中的配置 */
const actionField = Object.assign({}, commonActionField);

/** 自定义tabList */
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
        creationTimestamp: {
          dataField: ['creationTimestamp'],
          dataFormat: dataFormatConfig['time'],
          label: t('创建时间'),
          noExsitedValue: defaulNotExistedValue,
          order: '25'
        },
        isDefaul: {
          dataField: ['annotations', 'storageclass.beta.kubernetes.io/is-default-class'],
          dataFormat: dataFormatConfig['text'],
          label: 'Default Class',
          noExsitedValue: 'false',
          order: '20'
        }
      }
    },
    // parameters: {
    //   dataField: ['parameters'],
    //   displayField: {
    //     type: {
    //       dataField: ['type'],
    //       dataFormat: dataFormatConfig['mapText'],
    //       mapTextConfig: {
    //         CLOUD_BASIC: t('普通云硬盘'),
    //         CLOUD_PREMIUM: t('高性能云硬盘'),
    //         CLOUD_SSD: t('SSD云硬盘'),
    //         cbs: t('普通云硬盘')
    //       },
    //       label: t('云盘类型'),
    //       noExsitedValue: defaulNotExistedValue,
    //       order: '5'
    //     }
    //   }
    // },
    provisioner: {
      dataField: ['provisioner'],
      displayField: {
        provisioner: {
          dataField: [''],
          dataFormat: dataFormatConfig['text'],
          label: t('来源'),
          noExsitedValue: defaulNotExistedValue,
          order: '10'
        }
      }
    },
    reclaimPolicy: {
      dataField: ['reclaimPolicy'],
      displayField: {
        reclaimPolicy: {
          dataField: [''],
          dataFormat: dataFormatConfig['text'],
          label: t('回收策略'),
          noExsitedValue: 'Delete',
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

/** storage classes 的配置 */
export const sc = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'sc',
    requestType: {
      list: 'storageclasses'
    },
    displayField,
    actionField,
    detailField
  });
};
