export interface RecordSet<T, ExtendParamsT = any> {
  data?: ExtendParamsT;
  recordCount: number;
  records: T[];
}
