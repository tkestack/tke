import { RouteState } from '../../../../helpers';
import { FetcherState, FFListModel, OperationResult, RecordSet, WorkflowState } from '@tencent/ff-redux';
import { AlarmRecord, AlarmRecordFilter } from  '../models';


export interface RootState {
    /** 路由 */
    route?: RouteState;
    alarmRecord?: FFListModel<AlarmRecord, AlarmRecordFilter>;
}
