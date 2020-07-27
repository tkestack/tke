import { Identifiable } from '@tencent/ff-redux';
import { Validation } from '../../common/models';

export interface Cluster extends Identifiable {
  metadata?: {
    name?: string;
  };

  spec?: {
    displayName?: string;
    tenantID?: string;
    type?: string;
    version?: string;
  };
}
