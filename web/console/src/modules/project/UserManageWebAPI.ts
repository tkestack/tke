import { QueryState, RecordSet } from '@tencent/ff-redux';
import {
    Method,
    operationResult,
    reduceK8sQueryString,
    reduceK8sRestfulPath,
    reduceNetworkRequest
} from '../../../helpers';
import { ResourceInfo, RequestParams } from '../common/models';
import { resourceConfig } from '../../../config';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { METHODS } from 'http';
import {
    User,
    UserFilter,
    Strategy,
    StrategyFilter,
    GroupAssociation,
    GroupFilter, GroupPlain,
    PolicyAssociation,
    PolicyFilter,
    PolicyInfoFilter,
    PolicyPlain,
    RoleAssociation,
    RoleFilter,
    RolePlain
} from './models';

// @ts-ignore
const tips = seajs.require('tips');

class RequestResult {
    data: any;
    error: any;
}
const SEND = async (url: string, method: string, bodyData: any, tipErr: boolean = true) => {
    // 构建参数
    let params: RequestParams = {
        method: method,
        url,
        data: bodyData
    };
    let resp = new RequestResult();
    try {
        let response = await reduceNetworkRequest(params);
        if (response.code !== 0) {
            if (tipErr === true) {
                tips.error(response.message, 2000);
            }
            resp.data = bodyData;
            resp.error = response.message;
        } else {
            if (method !== Method.get) {
                tips.success('操作成功', 2000);
            }
            resp.data = response.data;
            resp.error = null;
        }
        return resp;
    } catch (error) {
        if (tipErr === true) {
            tips.error(error.response.data.message, 2000);
        }
        resp.data = bodyData;
        resp.error = error.response.data.message;
        return resp;
    }
};

const GET = async (url: string, tipErr: boolean = true) => {
    let response = await SEND(url, Method.get, null, tipErr);
    return response;
};
const DELETE = async (url: string, tipErr: boolean = true) => {
    let response = await SEND(url, Method.delete, null, tipErr);
    return response;
};
const POST = async (url: string, bodyData: any, tipErr: boolean = true) => {
    let response = await SEND(url, Method.post, JSON.stringify(bodyData), tipErr);
    return response;
};

const PUT = async (url: string, bodyData: any, tipErr: boolean = true) => {
    let response = await SEND(url, Method.put, JSON.stringify(bodyData), tipErr);
    return response;
};

const PATCH = async (url: string, bodyData: any, tipErr: boolean = true) => {
    let response = await SEND(url, Method.patch, JSON.stringify(bodyData), tipErr);
    return response;
};

/**
 * 用户列表的查询
 * @param query 列表查询条件参数
 */
export async function fetchUserList(query: QueryState<UserFilter>) {
    let users: User[] = [];
    const { search, filter } = query;
    const { projectId } = filter;
    // let { isPolicyUser = false } = filter;
    // const queryObj = !search ?
    //     {
    //         fieldSelector: {
    //             policy: true
    //         }
    //     }
    //     :
    //     {
    //         fieldSelector: {
    //             keyword: search || ''
    //         }
    //     };

    try {
        // const resourceInfo: ResourceInfo = isPolicyUser ? resourceConfig()['user'] : resourceConfig()['localidentity'];
        // const url = reduceK8sRestfulPath({ resourceInfo });
        // const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });
        const response = await reduceNetworkRequest({
            method: Method.get,
            url: `/apis/auth.tkestack.io/v1/projects/${projectId}/users`
            // url: url + queryString
        }, '', projectId);

        if (response.data.items) {
            users = response.data.items;
        } else {
            users = [];
        }
    } catch (error) {
        tips.error(error.response.data.message, 2000);
    }
    const result: RecordSet<User> = {
        recordCount: users.length,
        records: users
    };
    return result;
}

/**
 * 增加用户
 * @param [userInfo] 用户数据, 这里和actions.user.addUser.start([userInfo]);的对应
 */
export async function addUser([userInfo]) {
    const { projectId, users, policies } = userInfo;
    try {
        const resourceInfo: ResourceInfo = resourceConfig()['localidentity'];
        // const url = reduceK8sRestfulPath({ resourceInfo });
        // console.log('addUser get userInfo is:', userInfo);
        const url = `/apis/auth.tkestack.io/v1/projects/${projectId}/users`;
        const response = await reduceNetworkRequest({
            method: Method.post,
            url,
            data: {
                users,
                policies
            }
        });
        if (response.code === 0) {
            tips.success(t('操作成功'), 2000);
            return operationResult(userInfo);
        } else {
            // 是否给tip得看具体返回的数据
            return operationResult(userInfo, response);
        }
    } catch (error) {
        tips.error(error.response.data.message, 2000);
        // 返回相关数据存储在redux中, 这里的error应该不用reduceNetworkWorkflow作数据处理
        return operationResult(userInfo, error.response);
    }
}

/**
 * 删除用户
 * @param name 用户名
 */
export async function removeUser([name]) {
    try {
        const resourceInfo: ResourceInfo = resourceConfig()['localidentity'];
        const url = reduceK8sRestfulPath({ resourceInfo, specificName: name });
        const response = await reduceNetworkRequest({
            method: 'DELETE',
            url
        });
        if (response.code === 0) {
            tips.success('删除成功', 2000);
            return operationResult(name);
        } else {
            return operationResult(name, response);
        }
    } catch (error) {
        tips.error(error.response.data.message, 2000);
        return operationResult(name, error.response);
    }
}

/**
 * 修改用户
 * @param [userInfo] 用户数据, 这里和actions.user.addUser.start([userInfo]);的对应
 */
export async function getUser(name: string) {
    try {
        const resourceInfo: ResourceInfo = resourceConfig()['localidentity'];
        const url = reduceK8sRestfulPath({ resourceInfo, specificName: name });
        const response = await reduceNetworkRequest({
            method: 'GET',
            url
        });
        if (response.code === 0) {
            return operationResult(response.data);
        } else {
            // 是否给tip得看具体返回的数据
            return operationResult(name, response);
        }
    } catch (error) {
        tips.error(error.response.data.message, 2000);

        // 返回相关数据存储在redux中, 这里的error应该不用reduceNetworkWorkflow作数据处理
        return operationResult(name, error.response);
    }
}

/**
 * 修改用户
 * @param [userInfo] 用户数据, 这里和actions.user.addUser.start([userInfo]);的对应
 */
export async function updateUser(user: User) {
    try {
        const resourceInfo: ResourceInfo = resourceConfig()['localidentity'];
        const url = reduceK8sRestfulPath({ resourceInfo, specificName: user.metadata.name });
        const response = await reduceNetworkRequest({
            method: Method.put,
            url,
            data: user
        });
        if (response.code === 0) {
            setTimeout(() => {
                tips.success(t('修改成功'), 2000);
            }, 1000);
            return operationResult(response.data);
        } else {
            // 是否给tip得看具体返回的数据
            return operationResult(user, response);
        }
    } catch (error) {
        tips.error(error.response.data.message, 2000);

        // 返回相关数据存储在redux中, 这里的error应该不用reduceNetworkWorkflow作数据处理
        return operationResult(user, error.response);
    }
}

/**
 * 用户组的查询
 * @param filter 查询条件参数
 */
export async function fetchPolicy(filter: PolicyInfoFilter) {
    const resourceInfo: ResourceInfo = resourceConfig()['policy'];
    const url = reduceK8sRestfulPath({ resourceInfo, specificName: filter.name });
    console.log('fetchPolicy url', url);
    let rr: RequestResult = await GET(url);
    return rr.data;
}

/**
 * 策略列表的查询，不参杂其他场景参数
 * @param query 列表查询条件参数
 */
export async function fetchPolicyPlainList(query: QueryState<PolicyFilter>) {
    const { search, filter, keyword } = query;
    const queryObj = {
        // 'fieldSelector=keyword': search || ''
    };
    console.log('fetchPolicyPlainList query is:', query);
    let queryString = '';
    if (filter.resource === 'platform') {
        queryString = '?fieldSelector=spec.scope!=project';
    } else if (filter.resource === 'project') {
        queryString = '?fieldSelector=spec.scope=project';
    }

    const resourceInfo: ResourceInfo = resourceConfig()['policy'];
    const url = reduceK8sRestfulPath({ resourceInfo });
    // const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });
    console.log('fetchPolicyPlainList url + queryString', url, queryString, 111, query);
    let rr: RequestResult = await GET(url + queryString);
    let items: PolicyPlain[] =
        !rr.error && rr.data.items
            ? rr.data.items.map(i => {
                return {
                    id: i.metadata && i.metadata.name,
                    name: i.metadata && i.metadata.name,
                    displayName: i.spec && i.spec.displayName,
                    category: i.spec && i.spec.category,
                    description: i.spec && i.spec.description,
                    tenantID: i.sepc && i.spec.tenantID
                };
            })
            : [];
    console.log('fetchPolicyPlainList items is:', items);
    const result: RecordSet<PolicyPlain> = {
        recordCount: items.length,
        records: items
    };
    return result;
}

/**
 * 已经绑定的策略列表的查询
 * @param query 列表查询条件参数
 */
export async function fetchPolicyAssociatedList(query: QueryState<PolicyFilter>) {
    const { search, filter } = query;
    const queryObj = {};

    const resourceInfo: ResourceInfo = resourceConfig()[filter.resource];
    const url = reduceK8sRestfulPath({ resourceInfo, specificName: filter.resourceID, extraResource: 'policies' });
    const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });
    let rr: RequestResult = await GET(url + queryString);
    let items: PolicyPlain[] =
        !rr.error && rr.data.items
            ? rr.data.items.map(i => {
                return {
                    id: i.metadata && i.metadata.name,
                    name: i.metadata && i.metadata.name,
                    displayName: i.spec && i.spec.displayName,
                    category: i.spec && i.spec.category,
                    description: i.spec && i.spec.description
                };
            })
            : [];
    const result: RecordSet<PolicyPlain> = {
        recordCount: items.length,
        records: items
    };
    return result;
}

/**
 * 关联策略
 * @param param0
 * @param params
 */
export async function associatePolicy([policy]: PolicyAssociation[], params: PolicyFilter) {
    const resourceInfo: ResourceInfo = resourceConfig()[params.resource];
    const url = reduceK8sRestfulPath({ resourceInfo, specificName: params.resourceID, extraResource: 'policybinding' });
    let rr: RequestResult = await POST(url, {
        policies: policy.addPolicies.map(p => {
            return p.name;
        })
    });
    return operationResult(rr.data, rr.error);
}

/**
 * 解关联策略
 * @param param0
 * @param params
 */
export async function disassociatePolicy([policy]: PolicyAssociation[], params: PolicyFilter) {
    const resourceInfo: ResourceInfo = resourceConfig()[params.resource];
    const url = reduceK8sRestfulPath({ resourceInfo, specificName: params.resourceID, extraResource: 'policyunbinding' });
    let rr: RequestResult = await POST(url, {
        policies: policy.removePolicies.map(p => {
            return p.name;
        })
    });
    return operationResult(rr.data, rr.error);
}


/**
 * 角色列表的查询，不参杂其他场景参数
 * @param query 列表查询条件参数
 */
export async function fetchRolePlainList(query: QueryState<RoleFilter>) {
    const { keyword, filter } = query;
    const queryObj = keyword
        ? {
            'fieldSelector=spec.displayName': keyword
        }
        : {};

    const resourceInfo: ResourceInfo = resourceConfig()['role'];
    const url = reduceK8sRestfulPath({ resourceInfo });
    const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });
    let rr: RequestResult = await GET(url + queryString);
    let items: RolePlain[] =
        !rr.error && rr.data.items
            ? rr.data.items.map(i => {
                return {
                    id: i.metadata && i.metadata.name,
                    name: i.metadata && i.metadata.name,
                    displayName: i.spec && i.spec.displayName,
                    description: i.spec && i.spec.description
                };
            })
            : [];
    const result: RecordSet<RolePlain> = {
        recordCount: items.length,
        records: items
    };
    return result;
}

/**
 * 已经绑定的角色列表的查询
 * @param query 列表查询条件参数
 */
export async function fetchRoleAssociatedList(query: QueryState<RoleFilter>) {
    const { search, filter } = query;
    const queryObj = {};

    const resourceInfo: ResourceInfo = resourceConfig()[filter.resource];
    const url = reduceK8sRestfulPath({ resourceInfo, specificName: filter.resourceID, extraResource: 'roles' });
    const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });
    let rr: RequestResult = await GET(url + queryString);
    let items: RolePlain[] =
        !rr.error && rr.data.items
            ? rr.data.items.map(i => {
                return {
                    id: i.metadata && i.metadata.name,
                    name: i.metadata && i.metadata.name,
                    displayName: i.spec && i.spec.displayName,
                    description: i.spec && i.spec.description
                };
            })
            : [];
    const result: RecordSet<RolePlain> = {
        recordCount: items.length,
        records: items
    };
    return result;
}


/**
 * 关联角色
 * @param param0
 * @param params
 */
export async function associateRole([role]: RoleAssociation[], params: RoleFilter) {
    const resourceInfo: ResourceInfo = resourceConfig()[params.resource];
    const url = reduceK8sRestfulPath({ resourceInfo, specificName: params.resourceID, extraResource: 'binding' });
    let rr: RequestResult = await POST(url, { roles: role.addRoles });
    return operationResult(rr.data, rr.error);
}

/**
 * 解关联角色
 * @param param0
 * @param params
 */
export async function disassociateRole([role]: RoleAssociation[], params: RoleFilter) {
    const resourceInfo: ResourceInfo = resourceConfig()[params.resource];
    const url = reduceK8sRestfulPath({ resourceInfo, specificName: params.resourceID, extraResource: 'unbinding' });
    let rr: RequestResult = await POST(url, { roles: role.removeRoles });
    return operationResult(rr.data, rr.error);
}


/**
 * 用户组列表的查询，不参杂其他场景参数
 * @param query 列表查询条件参数
 */
export async function fetchGroupPlainList(query: QueryState<GroupFilter>) {
    const { search, filter } = query;
    const queryObj = {
        'fieldSelector=keyword': search || ''
    };

    const resourceInfo: ResourceInfo = resourceConfig()['group'];
    const url = reduceK8sRestfulPath({ resourceInfo });
    const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });
    let rr: RequestResult = await GET(url + queryString);
    let items: GroupPlain[] =
        !rr.error && rr.data.items
            ? rr.data.items.map(i => {
                return {
                    id: i.metadata && i.metadata.name,
                    name: i.metadata && i.metadata.name,
                    displayName: i.spec && i.spec.displayName,
                    description: i.spec && i.spec.description
                };
            })
            : [];
    const result: RecordSet<GroupPlain> = {
        recordCount: items.length,
        records: items
    };
    return result;
}

/**
 * 已经绑定的用户组列表的查询
 * @param query 列表查询条件参数
 */
export async function fetchGroupAssociatedList(query: QueryState<GroupFilter>) {
    const { search, filter } = query;
    const queryObj = {};

    const resourceInfo: ResourceInfo = resourceConfig()[filter.resource];
    const url = reduceK8sRestfulPath({ resourceInfo, specificName: filter.resourceID, extraResource: 'groups' });
    const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });
    let rr: RequestResult = await GET(url + queryString);
    let items: GroupPlain[] =
        !rr.error && rr.data.items
            ? rr.data.items.map(i => {
                return {
                    id: i.metadata && i.metadata.name,
                    name: i.metadata && i.metadata.name,
                    displayName: i.spec && i.spec.displayName,
                    description: i.spec && i.spec.description
                };
            })
            : [];
    const result: RecordSet<GroupPlain> = {
        recordCount: items.length,
        records: items
    };
    return result;
}

/**
 * 关联用户组
 * @param param0
 * @param params
 */
export async function associateGroup([group]: GroupAssociation[], params: GroupFilter) {
    const resourceInfo: ResourceInfo = resourceConfig()[params.resource];
    const url = reduceK8sRestfulPath({ resourceInfo, specificName: params.resourceID, extraResource: 'binding' });
    let rr: RequestResult = await POST(url, { groups: group.addGroups });
    return operationResult(rr.data, rr.error);
}

/**
 * 解关联用户组
 * @param param0
 * @param params
 */
export async function disassociateGroup([group]: GroupAssociation[], params: GroupFilter) {
    const resourceInfo: ResourceInfo = resourceConfig()[params.resource];
    const url = reduceK8sRestfulPath({ resourceInfo, specificName: params.resourceID, extraResource: 'unbinding' });
    let rr: RequestResult = await POST(url, { groups: group.removeGroups });
    return operationResult(rr.data, rr.error);
}

