import { TYPES } from './schemaUtil';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { resourceConfig } from '@config';

export const receiverSchema = {
  properties: {
    apiVersion: {
      value: `${resourceConfig()['receiver'].group}/${resourceConfig()['receiver'].version}`
    },
    kind: {
      value: 'Receiver'
    },
    metadata: {
      properties: {
        name: TYPES.string,
        namespace: TYPES.string
      }
    },
    spec: {
      properties: {
        displayName: { ...TYPES.string, required: true, name: t('显示名称') },
        username: { ...TYPES.string, required: true, name: t('用户名') },
        identities: {
          properties: {
            mobile: { ...TYPES.string, required: true, name: t('移动电话') },
            email: { ...TYPES.string, required: true, name: t('电子邮件') },
            wechat_openid: { ...TYPES.string, required: true, name: t('微信OpenID') }
          }
        }
      }
    }
  }
};
