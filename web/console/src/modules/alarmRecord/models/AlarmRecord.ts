import { Identifiable } from '@tencent/ff-redux';

export interface AlarmRecord extends Identifiable {
    /** metadata */
    metadata?: any;

    /** spec */
    spec?: any;

    /** status */
    status?: any;

    /** data */
    data?: any;

    /** other */
    [props: string]: any;
}

export interface AlarmRecordFilter {
    /** 告警策略名字 */
    alarmPolicyName?: string;
}

