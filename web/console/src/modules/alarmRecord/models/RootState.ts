import { RouteState } from '../../../../helpers';
import { FFListModel } from '@tencent/ff-redux';
import { AlarmRecord, AlarmRecordFilter, ClusterFilter } from  '../models';
import { Cluster } from '../../common';
export interface RootState {
    /** 路由 */
    route?: RouteState;
    alarmRecord?: FFListModel<AlarmRecord, AlarmRecordFilter>;
    cluster?: FFListModel<Cluster, ClusterFilter>;
}
