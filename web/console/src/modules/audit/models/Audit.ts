import { Identifiable } from '@tencent/ff-redux';

export interface Audit extends Identifiable {
    auditID: string; // id 无实际意义
    stage: string; // 事件生成阶段，无需展示
    requestURI: string; // 请求uri
    verb: string; // 请求动作，包括创建，更新，删除
    userName: string; // 操作用户
    userAgent: string;
    resource: string; // 操作资源，包括cluster，各种addon，以及k8s资源
    namespace: string; // 资源所属namespace
    name: string; // 操作资源名称
    uid: string; // 资源uid，无需展示
    apiGroup: string; // 资源所属group
    apiVersion: string; // 无需展示
    status: string; // 操作事件结果 状态
    message: string; // 操作事件结果 描述
    reason: string; // 操作事件结果 失败原因
    details: string; // 操作事件结果 失败详情，是一个json字符串
    code: number; // 操作事件结果 对应http code
    requestObject: string; // 请求body
    responseObject: string; // 响应body
    requestReceivedTimestamp: number; // 请求时间
    stageTimestamp: number; // 事件生成时间
    clusterName: string; // 操作的集群
    sourceIPs: string; // 源地址
}

export interface AuditFilter {
    cluster?: string;
    namespace?: string;
    resource?: string;
    user?: string; // 名字
    query?: string; // 提供关键字模糊查询
    startTime?: number | ''; // 开始时间
    endTime?: number | ''; // 结束时间
}

export interface AuditFilterConditionValues {
    clusterName: string[]; // 可查询的集群列表
    namespace: string[]; // 可查询的命名空间列表
    resource: string[]; // 可查询的资源类型列表
    userName: string[]; // 可查询的操作用户
}

