import { ResourceTable } from './ResourceTable';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { TablePanelColumnProps } from '@tencent/ff-component';
import { Resource } from '../../../common';

export class ResourceTableChannel extends ResourceTable {
  getColumns(): TablePanelColumnProps<Resource>[] {
    return [
      {
        key: 'type',
        header: t('类型'),
        render: x => {
          if (x.spec.smtp || x.spec.text) {
            return t('邮件');
          }
          if (x.spec.tencentCloudSMS) {
            return t('短信');
          }
          if (x.spec.wechat) {
            return t('微信公众号');
          }
          if (x.spec.webhook) {
            return 'webhook';
          }

          return '-';
        }
      }
    ];
  }
}
