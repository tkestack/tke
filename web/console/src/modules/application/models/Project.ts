import { Identifiable } from '@tencent/ff-redux';

export interface Project extends Identifiable {
  metadata: {
    name: string;
  };
  spec: {
    displayName: string;
  };
}
