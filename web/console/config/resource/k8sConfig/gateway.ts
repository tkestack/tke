import { ResourceInfo, DetailField, DisplayField, DetailInfo } from '../../../src/modules/common/models';
import { commonActionField, defaulNotExistedValue, apiVersion, dataFormatConfig } from '../common';
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
  servers: {
    dataField: ['servers'],
    dataFormat: dataFormatConfig['array'],
    width: '15%',
    headTitle: t('端口配置'),
    noExsitedValue: defaulNotExistedValue
  },
  selector: {
    dataField: ['selecor'],
    dataFormat: dataFormatConfig['text'],
    width: '25%',
    headTitle: t('类型'),
    noExsitedValue: defaulNotExistedValue
  },
  ip: {
    dataField: ['ip'],
    dataFormat: dataFormatConfig['text'],
    width: '25%',
    headTitle: t('IP'),
    noExsitedValue: '-'
  },
  creationTimestamp: {
    dataField: ['metadata.creationTimestamp'],
    dataFormat: dataFormatConfig['time'],
    width: '25%',
    headTitle: t('创建时间'),
    noExsitedValue: defaulNotExistedValue
  },
  operator: {
    dataField: [''],
    dataFormat: dataFormatConfig['operator'],
    width: '10%',
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

/** resrouce action当中的配置 */
const actionField = Object.assign({}, commonActionField);

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
export const gateway = (k8sVersion: string) => {
  // apiVersion的配置
  const apiKind = apiVersion[k8sVersion].gateway;
  let config: ResourceInfo = {
    headTitle: apiKind.headTitle,
    basicEntry: apiKind.basicEntry,
    group: apiKind.group,
    version: apiKind.version,
    namespaces: 'namespaces',
    requestType: {
      list: 'gateways'
    },
    displayField,
    actionField,
    detailField
  };
  return config;
};
