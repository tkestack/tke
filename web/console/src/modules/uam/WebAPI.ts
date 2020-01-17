import { RecordSet } from '@tencent/qcloud-lib';
import { QueryState } from '@tencent/qcloud-redux-query';
import {
  reduceNetworkRequest,
  operationResult,
  reduceK8sRestfulPath,
  reduceK8sQueryString,
  Method
} from '../../../helpers';
import { User, UserFilter, Strategy, StrategyFilter, Category } from './models';
import { ResourceInfo, RequestParams } from '../common/models';
import { resourceConfig } from '../../../config';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

const tips = seajs.require('tips');

/**
 * 用户列表的查询
 * @param query 列表查询条件参数
 */
export async function fetchUserList(query: QueryState<UserFilter>) {
  let users: User[] = [];
  const { search, filter } = query;
  let { isPolicyUser = false } = filter;
  const queryObj = filter.ifAll
    ? {}
    : {
        keyword: search || ''
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
          description: strategy.description ? strategy.description : strategy.spec.description,
          statement: strategy.statement ? strategy.statement : strategy.spec.statement
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
