import { DisplayField } from '../../../src/modules/common/models';
import { dataFormatConfig } from './dataFormat';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

/** resource table 当中展示的数据 */
export const commonDisplayField: DisplayField = {
  check: {
    dataField: [],
    dataFormat: dataFormatConfig['checker'],
    width: '16px',
    headTitle: ' ',
    noExsitedValue: '-'
  },
  name: {
    dataField: ['metadata.name'],
    dataFormat: dataFormatConfig['text'],
    width: '20%',
    headTitle: t('名称'),
    noExsitedValue: '-',
    isLink: true, // 用于判断该值是否为链接
    isClip: true
  },
  labels: {
    dataField: ['metadata.labels'],
    dataFormat: dataFormatConfig['labels'],
    width: '15%',
    headTitle: 'Labels',
    noExsitedValue: t('无')
  },
  selector: {
    dataField: ['spec.selector.matchLabels'],
    dataFormat: dataFormatConfig['labels'],
    width: '20%',
    headTitle: 'Selector',
    noExsitedValue: t('无')
  }
};
