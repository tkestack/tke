import { DetailField, DisplayField, DetailInfo } from '../../../src/modules/common/models';
import { commonActionField, defaulNotExistedValue, dataFormatConfig, generateResourceInfo } from '../common';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

const displayField: DisplayField = {
  name: {
    dataField: ['metadata.name'],
    dataFormat: dataFormatConfig['text'],
    width: '20%',
    headTitle: t('名称'),
    noExsitedValue: defaulNotExistedValue,
    isLink: true, // 用于判断该值是否为链接
    isClip: true
  },
  status: {
    dataField: ['status.phase'],
    dataFormat: dataFormatConfig['status'],
    width: '15%',
    headTitle: t('状态'),
    noExsitedValue: defaulNotExistedValue
  },
  creationTimestamp: {
    dataField: ['metadata.creationTimestamp'],
    dataFormat: dataFormatConfig['time'],
    width: '25%',
    headTitle: t('创建时间'),
    noExsitedValue: defaulNotExistedValue
  },
  hard: {
    dataField: ['spec.hard'],
    dataFormat: dataFormatConfig['resourceLimit'],
    width: '25%',
    headTitle: t('资源限制'),
    noExsitedValue: defaulNotExistedValue
  },
  used: {
    dataField: ['status.used'],
    dataFormat: dataFormatConfig['resourceLimit'],
    width: '25%',
    headTitle: t('已使用'),
    noExsitedValue: defaulNotExistedValue
  },
  operator: {
    dataField: [''],
    dataFormat: dataFormatConfig['operator'],
    width: '10%',
    headTitle: t('操作'),
    operatorList: []
  }
};

/** resrouce action当中的配置 */
const actionField = Object.assign({}, commonActionField, {
  create: {
    isAvailable: false
  }
});

/** 自定义tabList */
const tabList = [
  {
    id: 'nsInfo',
    label: t('详情')
  },
  {
    id: 'yaml',
    label: 'YAML'
  }
];

/** 自定义配置详情的展示 */
const detailBasicInfo: DetailInfo = {};

/** 详情页面的相关配置 */
const detailField: DetailField = {
  tabList,
  detailInfo: Object.assign({}, detailBasicInfo)
};

/** namespace的配置 */
export const np = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'np',
    requestType: {
      list: 'namespaces'
    },
    displayField,
    actionField,
    detailField
  });
};
