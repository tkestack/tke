import { TYPES } from './schemaUtil';
import { resourceConfig } from '@config';

export const receiverGroupSchema = {
  properties: {
    apiVersion: {
      value: `${resourceConfig()['receiverGroup'].group}/${resourceConfig()['receiverGroup'].version}`
    },
    kind: {
      value: 'ReceiverGroup'
    },
    metadata: {
      properties: {
        name: TYPES.string,
        namespace: TYPES.string
      }
    },
    spec: {
      properties: {
        displayName: { ...TYPES.string, required: true },
        receivers: { ...TYPES.stringArray, required: true }
      }
    }
  }
};
