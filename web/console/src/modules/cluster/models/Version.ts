import { Identifiable } from '@tencent/ff-redux';

export interface Version extends Identifiable {
  name?: string;
  remark?: string;
  status?: string;
  version?: string;

  masterAutoEnable?: boolean;
  masterAutoDiableTips?: string;
  masterIndependentEnable?: boolean;
  masterIndependentDialbeTips?: string;
}
