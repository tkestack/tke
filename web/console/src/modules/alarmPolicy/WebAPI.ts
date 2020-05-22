import { reduceNs } from '@helper';
import { OperationResult, QueryState, RecordSet, uuid } from '@tencent/ff-redux';

import { resourceConfig } from '../../../config';
// import * as regionConfig from '../../../config/region';
import { reduceK8sRestfulPath } from '../../../helpers';
import { reduceNetworkRequest } from '../../../helpers/reduceNetwork';
import { RequestParams, ResourceInfo } from '../common/models';
import { AlarmPolicyMetrics } from './constants/Config';
import { AlarmPolicy, AlarmPolicyFilter, NamespaceFilter, Resource, ResourceFilter } from './models';
import { AlarmPolicyEdition, AlarmPolicyOperator, MetricsObject } from './models/AlarmPolicy';
import { Namespace } from './models/Namespace';

/** RESTFUL风格的请求方法 */
const Method = {
  get: 'GET',
  post: 'POST',
  patch: 'PATCH',
  delete: 'DELETE',
  put: 'PUT'
};

function humanizeDuration4Time(initSecons: number) {
  let seconds = initSecons;

  if (seconds < 0 || seconds > 24 * 3600) {
    return '00:00:00';
  }

  let result = '';

  if (seconds > 3600) {
    let hours = Math.floor(seconds / 3600);
    result += hours >= 10 ? `${hours}:` : `0${hours}:`;
    seconds -= hours * 3600;
  } else {
    result += `00:`;
  }
  if (seconds > 60) {
    let minutes = Math.floor(seconds / 60);
    result += minutes >= 10 ? `${minutes}:` : `0${minutes}:`;
    seconds -= minutes * 60;
  } else {
    result += `00:`;
  }
  result += seconds >= 10 ? `${seconds}` : `0${seconds}`;
  return result;
}
// 返回标准操作结果
function operationResult<T>(target: T[] | T, error?: any): OperationResult<T>[] {
  if (target instanceof Array) {
    return target.map(x => ({ success: !error, target: x, error }));
  }
  return [{ success: !error, target: target as T, error }];
}

/**获取Alarm列表 */
export async function fetchAlarmPolicy(query: QueryState<AlarmPolicyFilter>) {
  let { paging } = query;
  let alarmPolicyList: AlarmPolicy[] = [];
  let resourceInfo: ResourceInfo = resourceConfig().alarmPolicy;
  let url = reduceK8sRestfulPath({
    resourceInfo: {
      ...resourceInfo,
      requestType: {
        list: `monitor/clusters/${query.filter.clusterId}/${resourceInfo.requestType.list}`
      }
    }
  });
  let params: RequestParams = {
    method: Method.get,
    url
  };
  if (paging) {
    let { pageIndex, pageSize } = paging;
    params['page'] = pageIndex;
    params['page_size'] = pageSize;

    //业务侧一次性拉取
    /// #if project
    params['page'] = 1;
    params['page_size'] = 999;
    /// #endif

    params.url += `?page=${params['page']}&page_size=${params['page_size']}`;
  }

  // if (search) {
  //   params['Filter'] = {
  //     AlarmPolicyName: search
  //   };
  // }
  let total = 0;
  let items = [];
  try {
    let response = await reduceNetworkRequest(params);
    if (response.code === 0) {
      items = response.data.data.alarmPolicies.map(item => {
        return Object.assign({}, item, { id: uuid() });
      });
      total = response.data.data.total;
    }
  } catch (error) {
    // 这里是搜索的时候，如果搜索不到的话，会报404的错误，只有在 resourceNotFound的时候，不把错误抛出去
    if (error.code !== 'ResourceNotFound' && (error.response && error.response.status) !== 404) {
      throw error;
    }
  }

  // let response = await sendCapiRequest('tke', 'DescribeAlarmPolicies', params, query.filter.regionId);

  alarmPolicyList = items.map(item => {
    let alarmPolicyMetricsConfig =
      (item.AlarmPolicySettings.AlarmPolicyType === 'cluster'
        ? AlarmPolicyMetrics['independentClusetr']
        : AlarmPolicyMetrics[item.AlarmPolicySettings.AlarmPolicyType]) || [];
    item.ShieldSettings = item.ShieldSettings || {};
    let temp = {
      id: item.AlarmPolicyId || item.AlarmPolicySettings.AlarmPolicyName,
      alarmPolicyId: item.AlarmPolicyId || item.AlarmPolicySettings.AlarmPolicyName,
      clusterId: query.filter.clusterId,
      alarmPolicyName: item.AlarmPolicySettings.AlarmPolicyName,
      alarmPolicyDescription: item.AlarmPolicySettings.AlarmPolicyDescription,
      alarmPolicyType: item.AlarmPolicySettings.AlarmPolicyType,
      statisticsPeriod: item.AlarmPolicySettings.statisticsPeriod,
      alarmMetrics: [] as MetricsObject[],
      alarmObjectWorkloadType: item.WorkloadType,
      alarmObjectNamespace: item.Namespace,
      alarmObjetcsType: item.AlarmPolicySettings.AlarmObjectsType,
      alarmObjetcs: [],
      shieldTimeStart:
        item.ShieldSettings.ShieldTimeStart && humanizeDuration4Time(item.ShieldSettings.ShieldTimeStart),
      shieldTimeEnd: item.ShieldSettings.ShieldTimeEnd && humanizeDuration4Time(item.ShieldSettings.ShieldTimeEnd),
      receiverGroups: item.NotifySettings.ReceiverGroups || [],
      enableShield: item.ShieldSettings.EnableShield,
      notifyWays: (item.NotifySettings.NotifyWay || []).map(notifyWay => ({
        id: uuid(),
        channel: notifyWay.ChannelName,
        template: notifyWay.TemplateName
      }))
      // phoneNotifyOrder: item.NotifySettings.PhoneNotifyOrder,
      // phoneCircleTimes: item.NotifySettings.PhoneCircleTimes,
      // phoneInnerInterval: item.NotifySettings.PhoneInnerInterval && item.NotifySettings.PhoneInnerInterval / 60,
      // phoneCircleInterval: item.NotifySettings.PhoneCircleInterval && item.NotifySettings.PhoneCircleInterval / 60,
      // phoneArriveNotice: item.NotifySettings.PhoneArriveNotice
    };
    if (item.AlarmPolicySettings.AlarmObjectsType === 'part') {
      temp.alarmObjetcs = item.AlarmPolicySettings.AlarmObjects ? item.AlarmPolicySettings.AlarmObjects.split(',') : [];
    }
    item.AlarmPolicySettings.AlarmMetrics.forEach(metric => {
      let finder = alarmPolicyMetricsConfig.find(config => config.metricName === metric.MetricName);
      let tempMetrics = {
        type: finder ? finder.type : 'percent',
        measurement: metric.Measurement,
        metricId: metric.MetricId,
        statisticsPeriod: metric.StatisticsPeriod / 60,
        metricName: metric.MetricName,
        metricDisplayName: finder.metricDisplayName || metric.MetricDisplayName,
        // argusPolicyName: metric.ArgusPolicyName,
        evaluatorType: metric.Evaluator.Type,
        evaluatorValue: metric.Evaluator.Value,
        continuePeriod: metric.ContinuePeriod,
        status: metric.Status,
        tip: finder ? finder.tip : '',
        unit: metric.Unit
        // metricType: metric.MetricType
      };
      if (metric.MetricName === 'k8s_pod_mem_no_cache_bytes' || metric.MetricName === 'k8s_pod_mem_usage_bytes') {
        tempMetrics.unit = 'MB';
        tempMetrics.evaluatorValue = +tempMetrics.evaluatorValue / 1024 / 1024;
      }
      temp['alarmMetrics'].push(tempMetrics);
    });
    return temp;
  });

  const result: RecordSet<AlarmPolicy> = {
    recordCount: total || alarmPolicyList.length,
    records: alarmPolicyList
  };
  return result;
}

function getAlarmPolicyParams(alarmPolicyEdition: AlarmPolicyEdition[], opreator: AlarmPolicyOperator, receiverGroups) {
  let params = {
    ClusterInstanceId: opreator.clusterId,
    AlarmPolicyId: alarmPolicyEdition[0].alarmPolicyId,
    AlarmPolicySettings: {
      AlarmPolicyName: alarmPolicyEdition[0].alarmPolicyName,
      AlarmPolicyDescription: alarmPolicyEdition[0].alarmPolicyDescription,
      AlarmPolicyType: alarmPolicyEdition[0].alarmPolicyType,
      StatisticsPeriod: alarmPolicyEdition[0].statisticsPeriod * 60,
      AlarmMetrics: [],
      AlarmObjects: alarmPolicyEdition[0].alarmObjects.join(','),
      AlarmObjectsType: alarmPolicyEdition[0].alarmObjectsType
    },
    NotifySettings: {
      ReceiverGroups: receiverGroups.selections.map(group => group.metadata.name),
      NotifyWay: alarmPolicyEdition[0].notifyWays.map(notifyWay => ({
        ChannelName: notifyWay.channel,
        TemplateName: notifyWay.template
      }))
    },
    ShieldSettings: {
      EnableShield: false
    }
  };
  if (alarmPolicyEdition[0].alarmPolicyType === 'pod') {
    if (alarmPolicyEdition[0].alarmObjectsType === 'all') {
      if (alarmPolicyEdition[0].alarmObjectNamespace !== 'ALL') {
        params['Namespace'] = reduceNs(alarmPolicyEdition[0].alarmObjectNamespace);
      }
      if (alarmPolicyEdition[0].alarmObjectWorkloadType !== 'ALL') {
        params['WorkloadType'] = alarmPolicyEdition[0].alarmObjectWorkloadType;
      }
    } else {
      params['Namespace'] = reduceNs(alarmPolicyEdition[0].alarmObjectNamespace);
      params['WorkloadType'] = alarmPolicyEdition[0].alarmObjectWorkloadType;
    }
  }

  alarmPolicyEdition[0].alarmMetrics.forEach(item => {
    if (item.enable) {
      let metrics = {
        Measurement: item.measurement,
        // StatisticsPeriod: item.statisticsPeriod * 60,
        MetricName: item.metricName,
        MetricDisplayName: item.metricDisplayName,

        Evaluator: {
          Type: item.evaluatorType,
          Value: item.evaluatorValue
        },
        ContinuePeriod: item.continuePeriod,
        Unit: item.unit
      };
      if (item.metricName === 'k8s_pod_mem_no_cache_bytes' || item.metricName === 'k8s_pod_mem_usage_bytes') {
        metrics.Unit = 'B';
        metrics.Evaluator.Value = +metrics.Evaluator.Value * 1024 * 1024 + '';
      }
      if (item.metricId) {
        metrics['MetricId'] = item.metricId;
      }
      params.AlarmPolicySettings['AlarmMetrics'].push(metrics);
    }
  });
  return params;
}

/**添加Alarm列表 */
export async function editAlarmPolicy(
  alarmPolicyEdition: AlarmPolicyEdition[],
  opreator: AlarmPolicyOperator,
  receiverGroup
) {
  let params = getAlarmPolicyParams(alarmPolicyEdition, opreator, receiverGroup);
  let clusterId = opreator.clusterId;

  let resourceInfo: ResourceInfo = resourceConfig().alarmPolicy;
  let url = reduceK8sRestfulPath({
    resourceInfo: {
      ...resourceInfo,
      requestType: {
        list: `monitor/clusters/${clusterId}/${resourceInfo.requestType.list}`
      }
    }
  });

  if (params.AlarmPolicyId) {
    url += '/' + params.AlarmPolicyId;
  }

  try {
    let response = await reduceNetworkRequest({
      method: params.AlarmPolicyId ? Method.put : Method.post,
      data: params,
      url
    });
    if (!(response.Response && response.Response.Error)) {
      // 更新缓存状态
      return operationResult(alarmPolicyEdition[0]);
    } else {
      return operationResult(alarmPolicyEdition[0], response.Response.Error);
    }
  } catch (error) {
    if (error.response && error.response.data && error.response.data.err) {
      return operationResult(alarmPolicyEdition[0], { message: error.response.data.err });
    }
    return operationResult(alarmPolicyEdition[0], error);
  }
}

/**删除告警设置列表 */
export async function deleteAlarmPolicy(alarmPolicys: AlarmPolicy[], opreator: AlarmPolicyOperator) {
  let clusterId = opreator.clusterId;
  let resourceInfo: ResourceInfo = resourceConfig().alarmPolicy;
  let url = reduceK8sRestfulPath({
    resourceInfo: {
      ...resourceInfo,
      requestType: {
        list: `monitor/clusters/${clusterId}/${resourceInfo.requestType.list}`
      }
    }
  });
  try {
    let response = await Promise.all(
      alarmPolicys.map(alarmPolicy =>
        reduceNetworkRequest({
          method: Method.delete,
          data: {},
          url: url + '/' + alarmPolicy.alarmPolicyId
        })
      )
    );

    if (!response[0].code) {
      // 更新缓存状态
      return operationResult(alarmPolicys);
    } else {
      return operationResult(alarmPolicys, response[0].message);
    }
  } catch (error) {
    return operationResult(alarmPolicys, error);
  }
}

/**
 * namespace列表的查询
 * @param query namespace列表的查询
 * @param namespaceInfo 当前namespace查询api的配置
 */
export async function fetchNamespaceList(query: QueryState<NamespaceFilter>, namespaceInfo: ResourceInfo) {
  let { filter, search } = query;
  let { clusterId, regionId } = filter;
  let namespaceList = [];

  let k8sUrl = `/${namespaceInfo.basicEntry}/${namespaceInfo.version}/${namespaceInfo.requestType['list']}`;
  let url = k8sUrl;

  if (search) {
    url = url + '/' + search;
  }

  /** 构建参数 */
  let params: RequestParams = {
    method: Method.get,
    url,
    apiParams: {
      module: 'tke',
      interfaceName: 'ForwardRequest',
      regionId: +regionId || 1,
      restParams: {
        Method: Method.get,
        Path: url,
        Version: '2018-05-25'
      },
      opts: {
        tipErr: false
      }
    }
  };

  try {
    let response = await reduceNetworkRequest(params, clusterId);

    if (response.code === 0) {
      let list = JSON.parse(response.data.ResponseBody);
      if (list.items) {
        namespaceList = list.items.map(item => {
          return {
            id: uuid(),
            name: item.metadata.name
          };
        });
      } else {
        namespaceList.push({
          id: uuid(),
          name: list.metadata.name
        });
      }
    }
  } catch (error) {
    // 这里是搜索的时候，如果搜索不到的话，会报404的错误，只有在 resourceNotFound的时候，不把错误抛出去
    if (error.code !== 'ResourceNotFound') {
      throw error;
    }
  }

  const result: RecordSet<Namespace> = {
    recordCount: namespaceList.length,
    records: namespaceList
  };

  return result;
}

/**
 * 获取资源的具体的 yaml文件
 * @param resourceIns: Resource[]   当前需要请求的具体资源数据
 * @param resourceInfo: ResouceInfo 当前请求数据url的基本配置
 */
export async function fetchUserPortal(resourceInfo: ResourceInfo) {
  let url = reduceK8sRestfulPath({ resourceInfo });

  // 构建参数
  let params: RequestParams = {
    method: Method.get,
    url
  };

  let response = await reduceNetworkRequest(params);
  return response.data;
}

/**
 * Namespace查询
 * @param query Namespace查询的一些过滤条件
 */
export async function fetchProjectNamespaceList(query: QueryState<ResourceFilter>) {
  let { filter } = query;
  let NamespaceResourceInfo: ResourceInfo = resourceConfig().namespaces;
  let url = reduceK8sRestfulPath({
    resourceInfo: NamespaceResourceInfo,
    specificName: filter.specificName,
    extraResource: 'namespaces'
  });
  /** 构建参数 */
  let method = 'GET';
  let params: RequestParams = {
    method,
    url
  };

  let response = await reduceNetworkRequest(params);
  let namespaceList = [],
    total = 0;
  if (response.code === 0) {
    let list = response.data;
    total = list.items.length;
    namespaceList = list.items.map(item => {
      return Object.assign({}, item, { id: uuid(), name: item.metadata.name });
    });
  }

  const result: RecordSet<Resource> = {
    recordCount: total,
    records: namespaceList
  };

  return result;
}

export async function fetchPrometheuses() {
  let resourceInfo: ResourceInfo = resourceConfig().prometheus;
  let url = reduceK8sRestfulPath({
    resourceInfo
  });
  let params: RequestParams = {
    method: Method.get,
    url
  };
  let records = [];
  try {
    let response = await reduceNetworkRequest(params);
    if (response.code === 0) {
      records = response.data.items.map(item => {
        return Object.assign({}, item, { id: uuid() });
      });
    }
  } catch (error) {
    // 这里是搜索的时候，如果搜索不到的话，会报404的错误，只有在 resourceNotFound的时候，不把错误抛出去
    if (error.code !== 'ResourceNotFound') {
      throw error;
    }
  }
  return {
    records,
    recordCount: records.length
  };
}
