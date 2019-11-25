import { TYPES } from './schemaUtil';
import { resourceConfig } from '@config';

export const templateSchema = {
  properties: {
    apiVersion: {
      value: `${resourceConfig()['template'].group}/${resourceConfig()['template'].version}`
    },
    kind: {
      value: 'Template'
    },
    metadata: {
      properties: {
        name: TYPES.string,
        namespace: TYPES.string
      }
    },
    spec: {
      type: 'pickOne',
      pick: 'text',
      properties: {
        displayName: { ...TYPES.string, required: true },
        text: {
          properties: {
            body: { ...TYPES.string, required: true },
            header: { ...TYPES.string, required: true }
          }
        },
        tencentCloudSMS: {
          properties: {
            body: { ...TYPES.string, required: true },
            sign: { ...TYPES.string, required: true },
            templateID: { ...TYPES.string, required: true }
          }
        },
        wechat: {
          properties: {
            templateID: { ...TYPES.string, required: true },
            miniProgramAppID: TYPES.string,
            miniProgramPagePath: TYPES.string,
            url: TYPES.string
          }
        }
      }
    }
  }
};
