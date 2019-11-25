import { ResourceInfo, DetailField, DisplayField, DetailInfo } from '../../../src/modules/common/models';
import {
  defaulNotExistedValue,
  commonActionField,
  dataFormatConfig,
  apiVersion,
  commonTabList,
  workloadCommonTabList
} from '../common';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

const displayField: DisplayField = {
  name: {
    dataField: ['metadata.name'],
    dataFormat: dataFormatConfig['text'],
    width: '10%',
    headTitle: t('名称'),
    noExsitedValue: defaulNotExistedValue,
    isLink: true, // 判断该值是否为链接,
    isClip: true
  },
  type: {
    dataField: ['spec.lbDriver'],
    dataFormat: dataFormatConfig['text'],
    width: '10%',
    headTitle: t('Driver'),
    noExsitedValue: defaulNotExistedValue
  },
  lb: {
    dataField: ['status.lbInfo.lbID'],
    dataFormat: dataFormatConfig['text'],
    width: '10%',
    headTitle: '负载均衡对象',
    noExsitedValue: defaulNotExistedValue
  },
  backendService: {
    dataField: ['spec.backGroups'],
    dataFormat: dataFormatConfig['backendGroups'],
    width: '10%',
    headTitle: t('后端负载配置'),
    noExsitedValue: defaulNotExistedValue
  },
  creationTimestamp: {
    dataField: ['metadata.creationTimestamp'],
    dataFormat: dataFormatConfig['time'],
    width: '12%',
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
        name: t('更新后端负载'),
        actionType: 'updateBG',
        isInMoreOp: false
      },
      {
        name: t('添加后端负载'),
        actionType: 'createBG',
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
};

/** resrouce action当中的配置 */
const actionField = Object.assign({}, commonActionField);

/** 自定义配置详情的展示 */
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
          order: '10'
        },
        namespace: {
          dataField: ['namespace'],
          dataFormat: dataFormatConfig['text'],
          label: 'Namespace',
          noExsitedValue: defaulNotExistedValue,
          order: '5'
        },
        createTime: {
          dataField: ['creationTimestamp'],
          dataFormat: dataFormatConfig.time,
          label: t('创建时间'),
          noExsitedValue: defaulNotExistedValue,
          order: '60'
        }
      }
    },
    spec: {
      dataField: ['spec'],
      displayField: {
        config: {
          dataField: ['lbSpec'],
          dataFormat: dataFormatConfig.keyvalue,
          label: t('负载均衡配置'),
          noExsitedValue: defaulNotExistedValue,
          order: '40'
        },
        attribute: {
          dataField: ['attributes'],
          dataFormat: dataFormatConfig.keyvalue,
          label: t('负载均衡属性'),
          noExsitedValue: defaulNotExistedValue,
          order: '50'
        },
        lbDriver: {
          dataField: ['lbDriver'],
          dataFormat: dataFormatConfig.mapText,
          mapTextConfig: {
            'lbcf-clb-driver': t('腾讯云CLB')
          },
          label: t('类型'),
          noExsitedValue: defaulNotExistedValue,
          order: '30'
        }
      }
    }
  }
};

/** 自定义配置详情的展示 */
const backGroupInfo: DetailInfo = {
  backGroup: {
    backGroups: {
      dataField: ['spec.backGroups'],
      displayField: {
        name: {
          dataField: ['name'],
          dataFormat: dataFormatConfig['text'],
          label: t('名称'),
          noExsitedValue: defaulNotExistedValue
        },
        port: {
          dataField: ['port'],
          dataFormat: dataFormatConfig['lbcfBGPort'],
          label: t('端口&协议'),
          noExsitedValue: defaulNotExistedValue
        },
        labels: {
          dataField: ['labels'],
          dataFormat: dataFormatConfig['labels'],
          label: t('Label'),
          noExsitedValue: defaulNotExistedValue
        },
        backendRecords: {
          dataField: ['backendRecords'],
          dataFormat: dataFormatConfig['backendRecords'],
          label: t('backendRecords'),
          noExsitedValue: defaulNotExistedValue
        },
        operator: {
          dataField: ['name'],
          dataFormat: dataFormatConfig['operator'],
          label: t('操作'),
          noExsitedValue: defaulNotExistedValue
        }
      }
    }
  }
};

/** 详情页面的相关配置 */
const detailField: DetailField = {
  tabList: [...commonTabList],
  detailInfo: Object.assign({}, backGroupInfo, detailBasicInfo)
};

/** lbcf 的配置 */
export const lbcf = (k8sVersion: string) => {
  // apiVersion的配置
  const apiKind = apiVersion[k8sVersion].lbcf;
  let config: ResourceInfo = {
    headTitle: apiKind.headTitle,
    basicEntry: apiKind.basicEntry,
    group: apiKind.group,
    version: apiKind.version,
    namespaces: 'namespaces',
    requestType: {
      list: 'lbcflbs',
      addon: true,
      useDetailInfo: true,
      detailInfoList: {
        info: [{ text: '负载均衡', value: 'lbcf' }],
        yaml: [{ text: '负载均衡', value: 'lbcf' }, { text: '后端负载', value: 'lbcf_bg' }],
        event: [{ text: '负载均衡', value: 'lbcf' }, { text: '后端记录', value: 'lbcf_br' }]
      }
    },
    displayField,
    actionField,
    detailField
  };
  return config;
};

/** ingress 的配置 */
export const lbcf_bg = (k8sVersion: string) => {
  // apiVersion的配置
  const apiKind = apiVersion[k8sVersion].lbcf_bg;
  let config: ResourceInfo = {
    headTitle: apiKind.headTitle,
    basicEntry: apiKind.basicEntry,
    group: apiKind.group,
    version: apiKind.version,
    namespaces: 'namespaces',
    requestType: {
      list: 'lbcfbackendgroups',
      addon: true
    }
  };
  return config;
};

/** ingress 的配置 */
export const lbcf_br = (k8sVersion: string) => {
  // apiVersion的配置
  const apiKind = apiVersion[k8sVersion].lbcf_br;
  let config: ResourceInfo = {
    headTitle: apiKind.headTitle,
    basicEntry: apiKind.basicEntry,
    group: apiKind.group,
    version: apiKind.version,
    namespaces: 'namespaces',
    requestType: {
      list: 'lbcfbackendrecords',
      addon: true
    }
  };
  return config;
};
