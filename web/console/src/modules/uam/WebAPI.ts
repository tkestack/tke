import { QueryState, RecordSet } from '@tencent/ff-redux';
import {
    Method, operationResult, reduceK8sQueryString, reduceK8sRestfulPath, reduceNetworkRequest
} from '../../../helpers';
import {
  User,
  UserFilter,
  Strategy,
  StrategyFilter,
  Category,
  Role,
  RoleFilter,
  RoleInfoFilter,
  RolePlain,
  RoleAssociation,
  Group,
  GroupFilter,
  GroupInfoFilter,
  GroupPlain,
  GroupAssociation,
  UserPlain,
  CommonUserFilter,
  CommonUserAssociation,
  PolicyInfoFilter,
  PolicyPlain,
  PolicyFilter,
  PolicyAssociation,
} from './models';
import { ResourceInfo, RequestParams } from '../common/models';
import { resourceConfig } from '../../../config';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { METHODS } from 'http';

const tips = seajs.require('tips');

class RequestResult {
  data: any;
  error: any;
}
const SEND = async (
  url: string,
  method: string,
  bodyData: any,
  tipErr: boolean = true
) => {
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
  let { isPolicyUser = false } = filter;
  const queryObj = !search
    ? {}
    : {
        fieldSelector: {
          keyword: search || ''
        }
      };

  try {
    const resourceInfo: ResourceInfo = isPolicyUser ? resourceConfig()['user'] : resourceConfig()['localidentity'];
    const url = reduceK8sRestfulPath({ resourceInfo });
    const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });
    const response = await reduceNetworkRequest({
      method: Method.get,
      url: url + queryString
    });

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
  try {
    const resourceInfo: ResourceInfo = resourceConfig()['localidentity'];
    const url = reduceK8sRestfulPath({ resourceInfo });
    const response = await reduceNetworkRequest({
      method: Method.post,
      url,
      data: userInfo
    });
    if (response.code === 0) {
      tips.success(t('添加成功'), 2000);
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
      tips.success(t('修改成功'), 2000);
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
 * 集群列表的查询
 * @param query 列表查询条件参数
 */
export async function fetchStrategyList(query: QueryState<StrategyFilter>) {
  let strategys: Strategy[] = [];
  let recordCount = 0;
  const { search, paging } = query;
  const queryObj = {
    keyword: search || '',
    page: paging.pageIndex - 1,
    page_size: paging.pageSize
  };
  try {
    const resourceInfo: ResourceInfo = resourceConfig()['policy'];
    const url = reduceK8sRestfulPath({ resourceInfo });
    const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });
    const response = await reduceNetworkRequest({
      method: 'GET',
      url: url + queryString
    });
    if (response.code === 0) {
      if (response.data.items) {
        strategys = response.data.items;
        recordCount = response.data.total;
      } else {
        strategys = [];
      }
    }
  } catch (error) {
    tips.error(error.response.data.message, 2000);
  }
  const result: RecordSet<Strategy> = {
    recordCount: strategys.length,
    records: strategys
  };

  return result;
}

/**
 * 增加策略
 * @param strategy 策略
 */
export async function addStrategy([strategy]) {
  try {
    const resourceInfo: ResourceInfo = resourceConfig()['policy'];
    const url = reduceK8sRestfulPath({ resourceInfo });
    const response = await reduceNetworkRequest({
      method: 'POST',
      url,
      data: strategy
    });
    if (response.code === 0) {
      tips.success(t('添加成功'), 2000);
      return operationResult(strategy);
    } else {
      return operationResult(strategy, response);
    }
  } catch (error) {
    tips.error(error.response.data.message, 2000);
    return operationResult(strategy, error.response);
  }
}

/**
 * 删除策略
 * @param id 策略id
 */
export async function removeStrategy([id]) {
  try {
    const resourceInfo: ResourceInfo = resourceConfig()['policy'];
    const url = reduceK8sRestfulPath({ resourceInfo, specificName: id });
    const response = await reduceNetworkRequest({
      method: 'DELETE',
      url
      // url: `/api/v1/polices/${id}`
    });
    if (response.code === 0) {
      tips.success(t('删除成功'), 2000);
      return operationResult(id);
    } else {
      return operationResult(id, response);
    }
  } catch (error) {
    tips.error(error.response.data.message, 2000);
    return operationResult(id, error.response);
  }
}

/**
 * 获取策略
 * @param id 策略id
 */
export async function getStrategy(id: string) {
  try {
    const resourceInfo: ResourceInfo = resourceConfig()['policy'];
    const url = reduceK8sRestfulPath({ resourceInfo, specificName: id + '' });
    const response = await reduceNetworkRequest({
      method: 'GET',
      url
    });
    if (response.code === 0) {
      return operationResult(response.data);
    } else {
      return operationResult(id, response);
    }
  } catch (error) {
    tips.error(error.response.data.message, 2000);
    return operationResult(id, error.response);
  }
}

/**
 * 更新策略
 * @param strategy 新的策略
 */
export async function updateStrategy(strategy) {
  try {
    const resourceInfo: ResourceInfo = resourceConfig()['policy'];
    const url = reduceK8sRestfulPath({ resourceInfo, specificName: strategy.metadata.name });
    const response = await reduceNetworkRequest({
      method: 'PUT',
      url,
      data: {
        metadata: {
          name: strategy.metadata.name,
          resourceVersion: strategy.metadata.resourceVersion
        },
        spec: Object.assign({}, strategy.spec, {
          displayName: strategy.name ? strategy.name : strategy.spec.displayName,
          description: strategy.description ? strategy.description : strategy.spec.description,
          statement: strategy.statement ? strategy.statement : strategy.spec.statement,
        })
      }
    });
    if (response.code === 0) {
      tips.success(t('修改成功'), 2000);
      return operationResult(response.data);
    } else {
      return operationResult(strategy, response);
    }
  } catch (error) {
    tips.error(error.response.data.message, 2000);
    return operationResult(strategy, error.response);
  }
}

/**
 * 获取策略所属的服务列表
 */
export async function fetchCategoryList() {
  let categories: Category[] = [];
  try {
    const resourceInfo: ResourceInfo = resourceConfig()['category'];
    const url = reduceK8sRestfulPath({ resourceInfo });
    const response = await reduceNetworkRequest({
      method: 'GET',
      url
      // url: '/api/v1/categories/'
    });
    if (response.code === 0) {
      if (response.data.items) {
        categories = response.data.items;
      } else {
        categories = [];
      }
    }
  } catch (error) {
    tips.error(error.response.data.message, 2000);
  }
  const result: RecordSet<Category> = {
    recordCount: categories.length,
    records: categories
  };

  return result;
}

/**
 * 增加策略关联的用户
 * @param id 策略id字符串
 * @param userNames  用户名数组
 */
export async function associateUser([{ id, userNames }]) {
  try {
    const resourceInfo: ResourceInfo = resourceConfig()['policy'];
    const url = reduceK8sRestfulPath({ resourceInfo, specificName: id, extraResource: 'binding' });
    const response = await reduceNetworkRequest({
      method: Method.post,
      url,
      // url: `/api/v1/polices/${id}/users`,
      data: {
        users: userNames.map(item => {
          return {
            id: item
          };
        })
      }
    });
    if (response.code === 0) {
      tips.success(t('添加成功'), 2000);
      return operationResult({ id, userNames });
    } else {
      return operationResult({ id, userNames }, response);
    }
  } catch (error) {
    tips.error(error.response.data.message, 2000);
    return operationResult({ id, userNames }, error.response);
  }
}

/**
 * 获取策略关联的用户
 * @param params 策略id
 */
export async function fetchStrategyAssociatedUsers(id: string) {
  // buildQueryString(filter)
  let associatedUsers: User[] = [];
  try {
    const resourceInfo: ResourceInfo = resourceConfig()['policy'];
    const url = reduceK8sRestfulPath({ resourceInfo, specificName: id, extraResource: 'users' });
    const response = await reduceNetworkRequest({
      method: 'GET',
      url
    });
    if (response.code === 0) {
      if (response.data.items) {
        associatedUsers = response.data.items;
      } else {
        associatedUsers = [];
      }
    }
  } catch (error) {
    tips.error(error.response.data.message, 2000);
  }
  const result: RecordSet<User> = {
    recordCount: associatedUsers.length,
    records: associatedUsers
  };

  return result;
}

/** 删除策略关联的用户 */
export async function removeAssociatedUser([{ id, userNames }]) {
  const data = {
    users: userNames.map(item => ({
      id: item
    }))
  };
  try {
    const resourceInfo: ResourceInfo = resourceConfig()['policy'];
    const url = reduceK8sRestfulPath({ resourceInfo, specificName: id, extraResource: 'unbinding' });
    const response = await reduceNetworkRequest({
      method: Method.post,
      url,
      data
    });
    if (response.code === 0) {
      tips.success('删除成功', 2000);
      return operationResult(data);
    } else {
      return operationResult(data, response);
    }
  } catch (error) {
    tips.error(error.response.data.message, 2000);
    return operationResult(data, error.response);
  }
}

/**
 * 角色列表的查询
 * @param query 列表查询条件参数
 */
export async function fetchRoleList(query: QueryState<RoleFilter>) {
  const { keyword, filter } = query;
  const queryObj = keyword ? {
    'fieldSelector=spec.displayName': keyword
  } : {};

  const resourceInfo: ResourceInfo = resourceConfig()['role'];
  const url = reduceK8sRestfulPath({ resourceInfo });
  const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });
  let rr: RequestResult = await GET(url + queryString);
  let roles: Role[] = (!rr.error && rr.data.items) ? rr.data.items : [];
  const result: RecordSet<Role> = {
    recordCount: roles.length,
    records: roles
  };
  return result;
}

/**
 * 角色的查询
 * @param filter 查询条件参数
 */
export async function fetchRole(filter: RoleInfoFilter) {
  const resourceInfo: ResourceInfo = resourceConfig()['role'];
  const url = reduceK8sRestfulPath({ resourceInfo, specificName: filter.name });
  let rr: RequestResult = await GET(url);
  return rr.data;
}

/**
 * 增加角色
 * @param roleInfo
 */
export async function addRole([roleInfo]) {
  const resourceInfo: ResourceInfo = resourceConfig()['role'];
  const url = reduceK8sRestfulPath({ resourceInfo });
  let rr: RequestResult = await POST(url, roleInfo);
  return operationResult(rr.data, rr.error);
}

/**
 * 修改角色
 * @param roleInfo
 */
export async function updateRole([roleInfo]) {
  const resourceInfo: ResourceInfo = resourceConfig()['role'];
  const url = reduceK8sRestfulPath({ resourceInfo, specificName: roleInfo.metadata.name });
  let rr: RequestResult = await PUT(url, roleInfo);
  return operationResult(rr.data, rr.error);
}

/**
 * 删除角色
 * @param role
 */
export async function deleteRole([role]: Role[]) {
  let resourceInfo: ResourceInfo = resourceConfig()['role'];
  const url = reduceK8sRestfulPath({ resourceInfo, specificName: role.metadata.name });
  let rr: RequestResult  = await DELETE(url);
  return operationResult(rr.data, rr.error);
}

/**
 * 角色列表的查询，不参杂其他场景参数
 * @param query 列表查询条件参数
 */
export async function fetchRolePlainList(query: QueryState<RoleFilter>) {
  const { keyword, filter } = query;
  const queryObj = keyword ? {
    'fieldSelector=spec.displayName': keyword
  } : {};

  const resourceInfo: ResourceInfo = resourceConfig()['role'];
  const url = reduceK8sRestfulPath({ resourceInfo });
  const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });
  let rr: RequestResult = await GET(url + queryString);
  let items: RolePlain[] = (!rr.error && rr.data.items) ? rr.data.items.map((i) => {
    return {
      id: i.metadata && i.metadata.name,
      name: i.metadata && i.metadata.name,
      displayName: i.spec && i.spec.displayName,
      description: i.spec && i.spec.description,
    };
  }) : [];
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
  const queryObj =  {};

  const resourceInfo: ResourceInfo = resourceConfig()[filter.resource];
  const url = reduceK8sRestfulPath({ resourceInfo, specificName: filter.resourceID, extraResource: 'roles' });
  const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });
  let rr: RequestResult = await GET(url + queryString);
  let items: RolePlain[] = (!rr.error && rr.data.items) ? rr.data.items.map((i) => {
    return {
      id: i.metadata && i.metadata.name,
      name: i.metadata && i.metadata.name,
      displayName: i.spec && i.spec.displayName,
      description: i.spec && i.spec.description,
    };
  }) : [];
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
 * 用户组列表的查询
 * @param query 列表查询条件参数
 */
export async function fetchGroupList(query: QueryState<GroupFilter>) {
  const { keyword, filter } = query;
  const queryObj = keyword ? {
    'fieldSelector=spec.displayName': keyword
  } : {};

  const resourceInfo: ResourceInfo = resourceConfig()['localgroup'];
  const url = reduceK8sRestfulPath({ resourceInfo });
  const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });
  let rr: RequestResult = await GET(url + queryString);
  let groups: Group[] = (!rr.error && rr.data.items) ? rr.data.items : [];
  const result: RecordSet<Group> = {
    recordCount: groups.length,
    records: groups
  };
  return result;
}

/**
 * 用户组的查询
 * @param filter 查询条件参数
 */
export async function fetchGroup(filter: GroupInfoFilter) {
  const resourceInfo: ResourceInfo = resourceConfig()['localgroup'];
  const url = reduceK8sRestfulPath({ resourceInfo, specificName: filter.name });
  let rr: RequestResult = await GET(url);
  return rr.data;
}

/**
 * 增加用户组
 * @param groupInfo
 */
export async function addGroup([groupInfo]) {
  const resourceInfo: ResourceInfo = resourceConfig()['localgroup'];
  const url = reduceK8sRestfulPath({ resourceInfo });
  let rr: RequestResult = await POST(url, groupInfo);
  return operationResult(rr.data, rr.error);
}

/**
 * 修改用户组
 * @param groupInfo
 */
export async function updateGroup([groupInfo]) {
  const resourceInfo: ResourceInfo = resourceConfig()['localgroup'];
  const url = reduceK8sRestfulPath({ resourceInfo, specificName: groupInfo.metadata.name });
  let rr: RequestResult = await PUT(url, groupInfo);
  return operationResult(rr.data, rr.error);
}

/**
 * 删除用户组
 * @param group
 */
export async function deleteGroup([group]: Group[]) {
  let resourceInfo: ResourceInfo = resourceConfig()['localgroup'];
  const url = reduceK8sRestfulPath({ resourceInfo, specificName: group.metadata.name });
  let rr: RequestResult  = await DELETE(url);
  return operationResult(rr.data, rr.error);
}


/**
 * 用户组列表的查询，不参杂其他场景参数
 * @param query 列表查询条件参数
 */
export async function fetchGroupPlainList(query: QueryState<GroupFilter>) {
  const { search, filter } = query;
  const queryObj =  {
    'fieldSelector=keyword': search || ''
  };

  const resourceInfo: ResourceInfo = resourceConfig()['group'];
  const url = reduceK8sRestfulPath({ resourceInfo });
  const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });
  let rr: RequestResult = await GET(url + queryString);
  let items: GroupPlain[] = (!rr.error && rr.data.items) ? rr.data.items.map((i) => {
    return {
      id: i.metadata && i.metadata.name,
      name: i.metadata && i.metadata.name,
      displayName: i.spec && i.spec.displayName,
      description: i.spec && i.spec.description,
    };
  }) : [];
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
  const queryObj =  {};

  const resourceInfo: ResourceInfo = resourceConfig()[filter.resource];
  const url = reduceK8sRestfulPath({ resourceInfo, specificName: filter.resourceID, extraResource: 'groups' });
  const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });
  let rr: RequestResult = await GET(url + queryString);
  let items: GroupPlain[] = (!rr.error && rr.data.items) ? rr.data.items.map((i) => {
    return {
      id: i.metadata && i.metadata.name,
      name: i.metadata && i.metadata.name,
      displayName: i.spec && i.spec.displayName,
      description: i.spec && i.spec.description,
    };
  }) : [];
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

/**
 * 用户列表的查询，不跟localidentities混用，不参杂其他场景参数，如策略、角色
 * @param query 列表查询条件参数
 */
export async function fetchCommonUserList(query: QueryState<CommonUserFilter>) {
  const { search, filter } = query;
  const queryObj =  {
    'fieldSelector=keyword': search || ''
  };

  const resourceInfo: ResourceInfo = resourceConfig()['user'];
  const url = reduceK8sRestfulPath({ resourceInfo });
  const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });
  let rr: RequestResult = await GET(url + queryString);
  let users: UserPlain[] = (!rr.error && rr.data.items) ? rr.data.items.map((i) => {
    /** localgroup对应localidentity，role对应user，localidentity的spec.username等同于user的spec.name */
    return {
      id: i.metadata && i.metadata.name,
      name: i.spec && (i.spec.name ? i.spec.name : i.spec.username),
      displayName: i.spec && i.spec.displayName
    };
  }) : [];
  const result: RecordSet<UserPlain> = {
    recordCount: users.length,
    records: users
  };
  return result;
}

/**
 * 已经绑定的用户列表的查询
 * @param query 列表查询条件参数
 */
export async function fetchCommonUserAssociatedList(query: QueryState<CommonUserFilter>) {
  const { search, filter } = query;
  const queryObj =  {};

  const resourceInfo: ResourceInfo = resourceConfig()[filter.resource];
  const url = reduceK8sRestfulPath({ resourceInfo, specificName: filter.resourceID, extraResource: 'users' });
  const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });
  let rr: RequestResult = await GET(url + queryString);
  let users: UserPlain[] = (!rr.error && rr.data.items) ? rr.data.items.map((i) => {
    /** localgroup对应localidentity，role对应user，localidentity的spec.username等同于user的spec.name */
    return {
      id: i.metadata && i.metadata.name,
      name: i.spec && (i.spec.name ? i.spec.name : i.spec.username),
      displayName: i.spec && i.spec.displayName
    };
  }) : [];
  const result: RecordSet<UserPlain> = {
    recordCount: users.length,
    records: users
  };
  return result;
}

/**
 * 关联用户
 * @param param0
 * @param params
 */
export async function commonAssociateUser([user]: CommonUserAssociation[], params: CommonUserFilter) {
  const resourceInfo: ResourceInfo = resourceConfig()[params.resource];
  const url = reduceK8sRestfulPath({ resourceInfo, specificName: params.resourceID, extraResource: 'binding' });
  let rr: RequestResult = await POST(url, { users: user.addUsers });
  return operationResult(rr.data, rr.error);
}

/**
 * 解关联用户
 * @param param0
 * @param params
 */
export async function commonDisassociateUser([user]: CommonUserAssociation[], params: CommonUserFilter) {
  const resourceInfo: ResourceInfo = resourceConfig()[params.resource];
  const url = reduceK8sRestfulPath({ resourceInfo, specificName: params.resourceID, extraResource: 'unbinding' });
  let rr: RequestResult = await POST(url, { users: user.removeUsers });
  return operationResult(rr.data, rr.error);
}

/**
 * 用户组的查询
 * @param filter 查询条件参数
 */
export async function fetchPolicy(filter: PolicyInfoFilter) {
  const resourceInfo: ResourceInfo = resourceConfig()['policy'];
  const url = reduceK8sRestfulPath({ resourceInfo, specificName: filter.name });
  let rr: RequestResult = await GET(url);
  return rr.data;
}

/**
 * 策略列表的查询，不参杂其他场景参数
 * @param query 列表查询条件参数
 */
export async function fetchPolicyPlainList(query: QueryState<PolicyFilter>) {
  const { search, filter } = query;
  const queryObj =  {
    // 'fieldSelector=keyword': search || ''
  };

  const resourceInfo: ResourceInfo = resourceConfig()['policy'];
  const url = reduceK8sRestfulPath({ resourceInfo });
  const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });
  let rr: RequestResult = await GET(url + queryString);
  let items: PolicyPlain[] = (!rr.error && rr.data.items) ? rr.data.items.map((i) => {
    return {
      id: i.metadata && i.metadata.name,
      name: i.metadata && i.metadata.name,
      displayName: i.spec && i.spec.displayName,
      category: i.spec && i.spec.category,
      description: i.spec && i.spec.description,
    };
  }) : [];
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
  const queryObj =  {};

  const resourceInfo: ResourceInfo = resourceConfig()[filter.resource];
  const url = reduceK8sRestfulPath({ resourceInfo, specificName: filter.resourceID, extraResource: 'policies' });
  const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });
  let rr: RequestResult = await GET(url + queryString);
  let items: PolicyPlain[] = (!rr.error && rr.data.items) ? rr.data.items.map((i) => {
    return {
      id: i.metadata && i.metadata.name,
      name: i.metadata && i.metadata.name,
      displayName: i.spec && i.spec.displayName,
      category: i.spec && i.spec.category,
      description: i.spec && i.spec.description,
    };
  }) : [];
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
  let rr: RequestResult = await POST(url, { policies: policy.addPolicies.map((p) => { return p.name }) });
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
  let rr: RequestResult = await POST(url, { policies: policy.removePolicies.map((p) => { return p.name }) });
  return operationResult(rr.data, rr.error);
}
