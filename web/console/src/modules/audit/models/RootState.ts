import { RouteState } from '../../../../helpers';
import { FetcherState, FFListModel, OperationResult, RecordSet, WorkflowState } from '@tencent/ff-redux';
import { Audit, AuditFilter, AuditFilterConditionValues } from  '../models';


export interface RootState {
    /** 路由 */
    route?: RouteState;
    auditList?: FFListModel<Audit, AuditFilter>;
    auditFilterCondition?: OperationResult<AuditFilterConditionValues>;
}
