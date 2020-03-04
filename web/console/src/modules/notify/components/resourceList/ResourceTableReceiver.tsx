import { ResourceTable } from './ResourceTable';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { TablePanelColumnProps } from '@tencent/ff-component';
import { Resource } from '../../../common';

export class ResourceTableReceiver extends ResourceTable {
  getColumns(): TablePanelColumnProps<Resource>[] {
    return [
      {
        key: 'username',
        header: t('username'),
        render: x => {
          return x.spec.username || '-';
        }
      },
      {
        key: 'mobile',
        header: t('移动号码'),
        render: x => {
          return (x.spec.identities && x.spec.identities.mobile) || '-';
        }
      },
      {
        key: 'email',
        header: t('电子邮件'),
        render: x => {
          return (x.spec.identities && x.spec.identities.email) || '-';
        }
      },
      {
        key: 'wechat_openid',
        header: t('微信OpenID'),
        render: x => {
          return (x.spec.identities && x.spec.identities.wechat_openid) || '-';
        }
      }
    ];
  }
}
