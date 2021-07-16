import moment from 'moment';
import { ModelType } from 'core/model';
import { Period } from './helper';
import { OneDayMillisecond, QUERY } from './constants';
import { request } from '../tce/request';

/**
 * 定义 chart panel 对外统一接口与业务流程方法
 */

// 基本的名值类型
export interface NameValueType {
  name: string;
  value: string;
}

export namespace CHART_PANEL {
  interface StoreKeyContent {
    groupBy: NameValueType,
    values: Array<{name: string}>
  }

  /**
   * 指标参数类型
   */
  export interface MetricType {
    expr: string;
    thousands: number;
    alias: string;
    unit?: string;
    storeKey?: string;                      // 传递key则缓存上次某指标图表搜索条件
    defaultGroupBy?: StoreKeyContent;       // 初始化时默认的维度条件
    defaultConditions?: StoreKeyContent;    // 初始化时默认的查询条件
    chartType?: any;
    tooltipLabels?: any;                    // tooltip显示时数据转化格式
    valueTransform?: any;
    scale?:  Array<any>;                    // Y轴显示刻度
    colorTheme?: string,
    colors?: Array<string>,
  }

  /**
   * 表格数据请求参数，统一表格数据结构
   * 输入端需保证 fields 唯一，每个 field 表示一个 chart
   */
  export interface TableType {
    table: string;
    fields: Array<MetricType>;
    conditions?: Array<Array<string>>; // [["id", "=", "value"]]
    groupBy?: Array<NameValueType>;
  }

  /**
   * 图表折线显示所需传入参数
   */
  export interface ChartDataType {
    min?: number;
    max?: number;
    unit?: string;
    title: string;
    field: any;
    labels: Array<string>;
    tooltipLabels: any;
    lines: Array<ModelType>;
  }


  /**
   * 用于groupBY中的值产生图表中折线的legend或instances的key
   * 使图表的折线的legend 跟 instance 对应，用于instances勾选功能
   */
  export function GenerateRowKey(row = []) {
    return row.length > 0 ? row.join(" | ") : "";
  }

  /**
   * 生成图表中显示折线的数据对象
   * @param legend
   * @param data
   * @returns {ModelType}
   * @constructor
   */
  export function GenerateLine(legend, data, disable = false) {
    return {
      legend: legend,
      disable: disable, // 该折线是否置灰色
      data: data
    } as ModelType;
  }

  // 生成时间粒度
  export function GeneratePeriodOptions(startTime, endTime): Array<{value: string, text: string}> {
    let periods = [60, 300, 3600, 86400];
    function periodUnit(period): string {
      const unit = ["秒", "分钟", "小时"];
      let temp = period / 60;

      let unitRound = 1;
      while (unitRound < 3 && temp >= 60) {
        temp = temp / 60;
        unitRound += 1;
      }
      if (temp < 1) {
        return `${period}${unit[0]}`;
      }
      return `${temp}${unit[unitRound]}`;
    }
    const period = Period(startTime, endTime);
    const l = moment(endTime).diff(moment(startTime));
    return periods.filter(d => d>=period && d < l).map(item => {
      return {
        value: item as any,
        text: periodUnit(item)
      }
    });
  }

  /**
   * 发起查询数据请求
   * @param {string} table
   * @param {Date} startTime
   * @param {Date} endTime
   * @param {Array<any>} fields
   * @param {Array<string>} dimensions
   * @param {Array<Array<any>>} conditions
   * @param {string} period
   */
  export async function RequestData(options: {table: string, startTime: Date, endTime: Date, fields: Array<any>, dimensions: Array<string>, conditions: Array<Array<any>>, period: string}) {
    const params = {
      table: options.table,
      startTime: options.startTime.getTime(),
      endTime: options.endTime.getTime(),
      fields: [...options.fields.map(item => `${item.expr}`)],
      conditions: options.conditions,
      orderBy: "timestamp",
      groupBy: [`timestamp(${options.period}s)`, ...options.dimensions],
      order: "asc",
      limit: QUERY.Limit
    };
    try {
      const res = await request({
        data: {
          Path: "/front/v1/get/query",
          RequestBody: params
        }
      });
      // 没有数据，抛出异常
      if (!res.hasOwnProperty("columns") || !res.hasOwnProperty("data")) {
        throw new Error();
      }
      // 更新图表数据
      const {columns, data} = res as any;
      if (!Array.isArray(columns) || !Array.isArray(data)) {
        throw new Error();
      }
      if (data.length === QUERY.Limit) {
        /**
         * 返回数据长度与请求限制一致，则说明对象显示数据点不够，需要分段发起请求
         */
          // 获取要显示的对象
        const instances = AggregateDataByDimensions(options.dimensions, columns, data);
        // 每个对象一天显示点数
        const instancePoints = 86400 / parseInt(options.period);
        // 显示的对象一天显示点数
        const sumInstancesPoints = Object.keys(instances).length * instancePoints;
        const daysNum = Math.ceil((params.endTime - params.startTime) / OneDayMillisecond);
        const requestNum = Math.ceil(sumInstancesPoints * daysNum / QUERY.Limit);
        const timeInterval = Math.floor((params.endTime - params.startTime) / requestNum);
        // 生成请求数组
        let requestPromises = [];
        for (let i = 1; i <= requestNum; i++) {
          const paramsTemp = {
            ...params,
            startTime: params.startTime + timeInterval * (i - 1),
            endTime: params.startTime + timeInterval * i
          };
          requestPromises.push(request({
            data: {
              Path: "/front/v1/get/query",
              RequestBody: paramsTemp
            }
          }));
        }
        const resAll = await Promise.all(requestPromises);
        // 将data转化为一维数组
        return {columns, data: [].concat(...resAll.map(item => item.data))};
      }
      return res;
    }
    catch (error) {
      throw error;
    }
  }

  /**
   * 根据维度对请求返回数据进行聚合
   * @param dimensions
   * @param columns
   * @param data
   * @returns {{}}
   */
  export function AggregateDataByDimensions(dimensions: Array<string>, columns: Array<string>, data) {
    columns = columns || [];
    data = data || [];
    /**
     * 根据 groupBy 条件对数据做聚合，已 groupBy 为key进行存储
     */
    let dimensionIndex = []; // 记录是 groupBy 数据的下标
    let dimensionData = {};  // 记录 groupBy 联合 key 集合数据
    dimensions.forEach(item => {
      const index = columns.indexOf(item);
      if (index >= 0) {
        dimensionIndex.push(index);
      }
    });

    data.forEach(row => {
      // 使用每一行数据中groupBy字段对应值的字符串合并作为该行的唯一key
      const rowKey = CHART_PANEL.GenerateRowKey(dimensionIndex.map(index => row[index]));
      if (dimensionData[rowKey]) {
        dimensionData[rowKey].push(row);
      }
      else {
        dimensionData[rowKey] = [row];
      }
    });
    return dimensionData;
  }

  /**
   *
   * @param {Date} startTime
   * @param {Date} endTime
   */
  export function OffsetTimeSeries(startTime: Date, endTime: Date, period: number, labels: Array<number>) {
    labels = [labels[0] || startTime.getTime()];
    const periodNum = Math.floor((endTime.getTime() - startTime.getTime()) / period / 1000);
    if (periodNum - labels.length > 0) {
      const periodMillisecond = period * 1000;
      const startTimestamp = startTime.getTime();
      const endTimestamp = endTime.getTime();
      let cycleIndex = 0;
      while (true) {
        if (cycleIndex > 15000) {
          // 安全机制，防止死循环，相当于15000分钟的跨度
          break;
        }
        cycleIndex += 1;
        const start = labels[0];
        const end = labels[labels.length - 1];
        if (start - startTimestamp >= periodMillisecond) {
          labels.unshift(start - periodMillisecond);
        }
        if (endTimestamp - end >= periodMillisecond) {
          labels.push(end + periodMillisecond);
        }
        if (periodNum < labels.length
          || (labels[0] - startTimestamp < periodMillisecond && endTimestamp - labels[labels.length - 1] < periodMillisecond)) {
          break;
        }
      }
    }
    return labels;
  }

  /**
   * 根据指标表达式查找在数据中的下标值
   * @param {Array<string>} columns
   * @param {string} metricExpr
   */
  export function FindMetricIndex(columns: Array<string>, metricExpr: string) {
    const match = /(\w+)?\((\w+)\)/.exec(metricExpr);
    let fieldIndex = -1;  // field 对应的在 data 中的下标
    if (match && match.length === 3) {
      fieldIndex = columns.indexOf(`${match[2]}_${match[1]}`);
      if (fieldIndex === -1) {
        // 兼容influxdb没拼接聚合方式在 k8s_workload_pod_restart_total
        fieldIndex = columns.indexOf(match[2]);
      }
    }
    else {
      fieldIndex = columns.findIndex(item => item === metricExpr);
    }

    return fieldIndex;
  }

}