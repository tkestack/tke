export interface RecordSet<T, ExtendParamsT = any> {
  data?: ExtendParamsT;
  recordCount: number;
  records: T[];
  continue?: boolean;
  continueToken?: string;
}
