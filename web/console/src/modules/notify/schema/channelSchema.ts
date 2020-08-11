import { TYPES } from './schemaUtil';
import { resourceConfig } from '@config';

export const channelSchema = {
  properties: {
    apiVersion: {
      value: `${resourceConfig()['channel'].group}/${resourceConfig()['channel'].version}`
    },
    kind: {
      value: 'Channel'
    },
    metadata: {
      properties: {
        name: TYPES.string,
        namespace: TYPES.string
      }
    },
    spec: {
      type: 'pickOne',
      pick: 'smtp',
      properties: {
        displayName: { ...TYPES.string, required: true },
        smtp: {
          properties: {
            email: { ...TYPES.string, required: true },
            password: { ...TYPES.string, required: true },
            smtpHost: { ...TYPES.string, required: true },
            smtpPort: { ...TYPES.number, required: true },
            tls: TYPES.boolean
          }
        },
        tencentCloudSMS: {
          properties: {
            appKey: { ...TYPES.string, required: true },
            sdkAppID: { ...TYPES.string, required: true },
            extend: TYPES.string
          }
        },
        wechat: {
          properties: {
            appID: { ...TYPES.string, required: true },
            appSecret: { ...TYPES.string, required: true }
          }
        },
        webhook: {
          properties: {
            url: { ...TYPES.string, required: true },
            headers: { ...TYPES.string, placeholder: '自定义Header，仅支持Key:Value格式，中间用;号分割。eg param1:1;param2:2' }
          }
        }
      }
    }
  }
};
