import { ResourceInfo, DetailField, DisplayField, DetailInfo } from '../../../src/modules/common/models';
import { defaulNotExistedValue, commonActionField, dataFormatConfig, generateResourceInfo } from '../common';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

const displayField: DisplayField = {
  name: {
    dataField: ['metadata.name'],
    dataFormat: dataFormatConfig['text'],
    width: '10%',
    headTitle: t('名称'),
    noExsitedValue: defaulNotExistedValue,
    isLink: true, // 用于判断该值是否为链接
    isClip: true
  },
  hosts: {
    dataField: ['spec.hosts'],
    dataFormat: dataFormatConfig.array,
    width: '12%',
    headTitle: t('Hosts'),
    tips: '',
    noExsitedValue: defaulNotExistedValue
  },
  gateway: {
    dataField: ['spec.gateways'],
    dataFormat: dataFormatConfig.array,
    width: '15%',
    headTitle: 'Gateway',
    noExsitedValue: '-'
  },
  http: {
    dataField: ['spec.http.length'],
    dataFormat: dataFormatConfig.text,
    width: '15%',
    headTitle: t('HTTP路由'),
    noExsitedValue: t('0')
  },
  tcp: {
    dataField: ['spec.tcp.length'],
    dataFormat: dataFormatConfig.text,
    width: '15%',
    headTitle: t('TCP路由'),
    noExsitedValue: t('0')
  },
  tls: {
    dataField: ['spec.tls.length'],
    dataFormat: dataFormatConfig.text,
    width: '15%',
    headTitle: t('TLS路由'),
    noExsitedValue: t('0')
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
          order: '0'
        },
        namespace: {
          dataField: ['namespace'],
          dataFormat: dataFormatConfig['text'],
          label: 'Namespace',
          noExsitedValue: defaulNotExistedValue,
          order: '5'
        },
        description: {
          dataField: ['annotations.description'],
          dataFormat: dataFormatConfig['text'],
          label: t('描述'),
          noExsitedValue: defaulNotExistedValue,
          order: '6'
        },
        labels: {
          dataField: ['labels'],
          dataFormat: dataFormatConfig['labels'],
          label: 'Labels',
          noExsitedValue: defaulNotExistedValue,
          order: '10'
        },
        createTime: {
          dataField: ['creationTimestamp'],
          dataFormat: dataFormatConfig['time'],
          label: t('创建时间'),
          tips: '',
          noExsitedValue: defaulNotExistedValue,
          order: '15'
        }
      }
    },
    spec: {
      dataField: ['spec'],
      displayField: {
        selector: {
          dataField: ['selector'],
          dataFormat: dataFormatConfig['labels'],
          label: 'Selector',
          tips: '',
          noExsitedValue: defaulNotExistedValue,
          order: '20'
        },
        type: {
          dataField: ['type'],
          dataFormat: dataFormatConfig['text'],
          label: t('访问方式'),
          tips: '',
          noExsitedValue: defaulNotExistedValue,
          order: '25'
        },
        clusterIP: {
          dataField: ['clusterIP'],
          dataFormat: dataFormatConfig['ip'],
          extraInfo: 'Headless Service',
          label: t('集群IP'),
          tips: '',
          noExsitedValue: defaulNotExistedValue,
          order: '30'
        },
        ports: {
          dataField: ['ports'],
          dataFormat: dataFormatConfig['ports'],
          label: t('端口映射'),
          noExsitedValue: defaulNotExistedValue,
          order: '40'
        }
      }
    },
    status: {
      dataField: ['status'],
      displayField: {
        ingressIP: {
          dataField: ['loadBalancer.ingress.0.ip'],
          dataFormat: dataFormatConfig['ip'],
          label: t('负载均衡IP'),
          tips: '',
          noExsitedValue: defaulNotExistedValue,
          order: '35'
        }
      }
    }
  },
  advancedInfo: {
    spec: {
      dataField: ['spec'],
      displayField: {
        externalTrafficPolicy: {
          dataField: ['spec.externalTrafficPolicy'],
          dataFormat: dataFormatConfig['text'],
          label: 'externalTrafficPolicy',
          tips: '',
          noExsitedValue: defaulNotExistedValue
        },
        sessionAffinityConfig: {
          dataField: ['spec.sessionAffinityConfig.clientIP.timeoutSeconds'],
          dataFormat: dataFormatConfig['text'],
          label: 'sessionAffinityConfig',
          tips: '',
          noExsitedValue: defaulNotExistedValue
        },
        sessionAffinity: {
          dataField: ['spec.sessionAffinity'],
          dataFormat: dataFormatConfig['text'],
          label: 'sessionAffinity',
          tips: '',
          noExsitedValue: defaulNotExistedValue
        }
      }
    }
  }
};

/** 详情页面的相关配置 */
const detailField: DetailField = {
  detailInfo: Object.assign({}, detailBasicInfo)
};

/** service的配置 */
export const virtualService = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'virtualservice',
    requestType: {
      list: 'virtualservices'
    },
    isRelevantToNamespace: true,
    displayField,
    actionField,
    detailField
  });
};
