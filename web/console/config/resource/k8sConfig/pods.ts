import { DetailField, DetailInfo } from '../../../src/modules/common/models';
import {
  commonDisplayField,
  commonActionField,
  defaulNotExistedValue,
  dataFormatConfig,
  commonDetailInfo,
  commonTabList,
  generateResourceInfo
} from '../common';

const displayField = Object.assign({}, commonDisplayField, {});

/** resource action当中的配置 */
const actionField = Object.assign({}, commonActionField);

/** 自定义配置详情 basic info的展示 */
const detailBasicInfo: DetailInfo = {
  info: {
    metadata: {
      dataField: ['metadata'],
      displayField: {
        name: {
          dataField: ['name'],
          dataFormat: dataFormatConfig['text'],
          label: 'Name',
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
        createdTime: {
          dataField: ['creationTimestamp'],
          dataFormat: dataFormatConfig['time'],
          label: 'Create Time',
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
    }
  }
};

/** 详情页面的相关配置 */
const detailField: DetailField = {
  tabList: [...commonTabList],
  detailInfo: Object.assign({}, commonDetailInfo('pod'), detailBasicInfo)
};

/** pods的配置 */
export const pods = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'pods',
    isRelevantToNamespace: true,
    requestType: {
      list: 'pods'
    },
    displayField,
    actionField,
    detailField
  });
};
