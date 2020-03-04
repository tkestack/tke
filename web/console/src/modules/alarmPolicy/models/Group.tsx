import { Identifiable } from '@tencent/ff-redux';

export interface Group extends Identifiable {
  groupName: string;
  groupId: string;
  userInfo: GroupUser[];
}

export interface GroupUser {
  groupNames: string[];
  name: string;
  uid: number;
}

export interface GroupFilter {
  groupName?: string;
}
