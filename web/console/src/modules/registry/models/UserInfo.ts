import { Identifiable } from '@tencent/ff-redux';

export interface UserInfo extends Identifiable {
  name: string;
  uid: string;
  groups?: string[];
  extra?: {
    displayname?: string;
    tenantid?: string;
  };
}
